package ports

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/common/auth"
	"github.com/kimnattanan/graph-rag-service/internal/common/logs"
	"github.com/kimnattanan/graph-rag-service/internal/common/metrics"
	"github.com/kimnattanan/graph-rag-service/internal/user/app"
	"github.com/kimnattanan/graph-rag-service/internal/user/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/user/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/user/domain/user"
)

func TestAccountFlow(t *testing.T) {
	server := httptest.NewServer(newTestHandler(t))
	t.Cleanup(server.Close)

	register(t, server.URL, "ada@example.com", "ada", "correct-horse")

	duplicate := request(t, http.MethodPost, server.URL+"/auth/register", "", map[string]string{
		"email":    "Ada@Example.com",
		"username": "ada",
		"password": "correct-horse",
	})
	if duplicate.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate status = %d, body = %s", duplicate.StatusCode, duplicate.Body)
	}

	loginBody := login(t, server.URL, "ada@example.com", "correct-horse")
	if loginBody.User.Role != UserRoleUser {
		t.Fatalf("role = %s", loginBody.User.Role)
	}
	if len(loginBody.User.Permissions) != 1 || loginBody.User.Permissions[0] != ConversationAsk {
		t.Fatalf("permissions = %#v", loginBody.User.Permissions)
	}
	if loginBody.User.Username != "ada" {
		t.Fatalf("username = %s", loginBody.User.Username)
	}
	principal, err := auth.Parse("test-secret", loginBody.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if principal.Role != auth.RoleUser || !auth.HasPermission(principal, auth.PermissionConversationAsk) {
		t.Fatalf("token principal = %#v", principal)
	}
	if auth.HasPermission(principal, auth.PermissionKnowledgeWrite) {
		t.Fatal("user token includes knowledge:write")
	}

	me := request(t, http.MethodGet, server.URL+"/users/me", loginBody.AccessToken, nil)
	if me.StatusCode != http.StatusOK {
		t.Fatalf("me status = %d, body = %s", me.StatusCode, me.Body)
	}

	missing := request(t, http.MethodGet, server.URL+"/users/me", "", nil)
	if missing.StatusCode != http.StatusUnauthorized {
		t.Fatalf("missing token status = %d", missing.StatusCode)
	}

	logout := request(t, http.MethodPost, server.URL+"/auth/logout", loginBody.AccessToken, nil)
	if logout.StatusCode != http.StatusNoContent {
		t.Fatalf("logout status = %d, body = %s", logout.StatusCode, logout.Body)
	}

	afterLogout := request(t, http.MethodGet, server.URL+"/users/me", loginBody.AccessToken, nil)
	if afterLogout.StatusCode != http.StatusUnauthorized {
		t.Fatalf("after logout status = %d, body = %s", afterLogout.StatusCode, afterLogout.Body)
	}

	wrong := request(t, http.MethodPost, server.URL+"/auth/login", "", map[string]string{
		"email":    "ada@example.com",
		"password": "wrong-password",
	})
	if wrong.StatusCode != http.StatusUnauthorized {
		t.Fatalf("wrong password status = %d", wrong.StatusCode)
	}

	again := login(t, server.URL, "ada@example.com", "correct-horse")
	deleted := request(t, http.MethodDelete, server.URL+"/users/me", again.AccessToken, nil)
	if deleted.StatusCode != http.StatusNoContent {
		t.Fatalf("delete status = %d, body = %s", deleted.StatusCode, deleted.Body)
	}

	gone := request(t, http.MethodPost, server.URL+"/auth/login", "", map[string]string{
		"email":    "ada@example.com",
		"password": "correct-horse",
	})
	if gone.StatusCode != http.StatusUnauthorized {
		t.Fatalf("deleted login status = %d, body = %s", gone.StatusCode, gone.Body)
	}
}

func newTestHandler(t *testing.T) http.Handler {
	t.Helper()
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	entry := logrus.NewEntry(logger)
	router := chi.NewRouter()
	router.Use(logs.NewStructuredLogger(logger))
	return HandlerFromMux(newTestServer(entry), router)
}

func newTestServer(entry *logrus.Entry) HttpServer {
	metricsClient := metrics.NoOp{}
	users := newMemUsers()
	sessions := newMemSessions()
	secret := "test-secret"
	ttl := time.Hour

	return NewHttpServer(app.Application{
		Commands: app.Commands{
			Register:      command.NewRegisterHandler(users, entry, metricsClient),
			Logout:        command.NewLogoutHandler(sessions, secret, entry, metricsClient),
			DeleteAccount: command.NewDeleteAccountHandler(users, entry, metricsClient),
		},
		Queries: app.Queries{
			Login:          query.NewLoginHandler(users, sessions, secret, ttl, entry, metricsClient),
			Authenticate:   query.NewAuthenticateHandler(users, sessions, secret, entry, metricsClient),
			GetCurrentUser: query.NewGetCurrentUserHandler(users, entry, metricsClient),
		},
	})
}

func register(t *testing.T, baseURL string, email string, username string, password string) {
	t.Helper()
	response := request(t, http.MethodPost, baseURL+"/auth/register", "", map[string]string{
		"email":    email,
		"username": username,
		"password": password,
	})
	if response.StatusCode != http.StatusNoContent {
		t.Fatalf("register status = %d, body = %s", response.StatusCode, response.Body)
	}
}

func login(t *testing.T, baseURL string, email string, password string) AuthResult {
	t.Helper()
	response := request(t, http.MethodPost, baseURL+"/auth/login", "", map[string]string{
		"email":    email,
		"password": password,
	})
	if response.StatusCode != http.StatusOK {
		t.Fatalf("login status = %d, body = %s", response.StatusCode, response.Body)
	}
	var result AuthResult
	if err := json.Unmarshal([]byte(response.Body), &result); err != nil {
		t.Fatal(err)
	}
	if result.AccessToken == "" || result.TokenType != auth.TokenTypeBearer {
		t.Fatalf("auth result = %#v", result)
	}
	return result
}

type recordedResponse struct {
	StatusCode int
	Body       string
}

func request(t *testing.T, method string, url string, token string, body any) recordedResponse {
	t.Helper()
	var reader io.Reader
	if body != nil {
		payload, err := json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
		reader = bytes.NewReader(payload)
	}
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		t.Fatal(err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	payload, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	return recordedResponse{StatusCode: response.StatusCode, Body: string(payload)}
}

type memUsers struct {
	mu      sync.Mutex
	byID    map[string]*user.User
	byEmail map[string]string
}

func newMemUsers() *memUsers {
	return &memUsers{
		byID:    map[string]*user.User{},
		byEmail: map[string]string{},
	}
}

func (m *memUsers) AddUser(_ context.Context, account *user.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.byEmail[account.Email()]; ok {
		return user.ErrEmailAlreadyExists
	}
	m.byID[account.ID()] = account
	m.byEmail[account.Email()] = account.ID()
	return nil
}

func (m *memUsers) GetUserByEmail(_ context.Context, email string) (*user.User, error) {
	normalized, err := user.NormalizeEmail(email)
	if err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	id, ok := m.byEmail[normalized]
	if !ok {
		return nil, user.ErrNotFound
	}
	return m.byID[id], nil
}

func (m *memUsers) GetUserByID(_ context.Context, id string) (*user.User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	account, ok := m.byID[id]
	if !ok {
		return nil, user.ErrNotFound
	}
	return account, nil
}

func (m *memUsers) DeleteUser(_ context.Context, id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	account, ok := m.byID[id]
	if !ok {
		return user.ErrNotFound
	}
	delete(m.byEmail, account.Email())
	delete(m.byID, id)
	return nil
}

type memSessions struct {
	mu   sync.Mutex
	byID map[string]user.Session
}

func newMemSessions() *memSessions {
	return &memSessions{byID: map[string]user.Session{}}
}

func (m *memSessions) AddSession(_ context.Context, session user.Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.byID[session.ID()] = session
	return nil
}

func (m *memSessions) GetSession(_ context.Context, sessionID string) (user.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	session, ok := m.byID[sessionID]
	if !ok {
		return user.Session{}, user.ErrSessionNotFound
	}
	return session, nil
}

func (m *memSessions) DeleteSession(_ context.Context, sessionID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.byID, sessionID)
	return nil
}
