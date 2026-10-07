package query

import (
	"context"
	"errors"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/common/auth"
	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/user/domain/user"
)

type Authenticate struct {
	Token string
}

func (q Authenticate) GoString() string {
	return "query.Authenticate{Token:\"***\"}"
}

type AuthenticateHandler decorator.QueryHandler[Authenticate, auth.User]

type authenticateHandler struct {
	users    user.Repository
	sessions user.SessionRepository
	secret   string
}

func NewAuthenticateHandler(
	users user.Repository,
	sessions user.SessionRepository,
	secret string,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) AuthenticateHandler {
	if users == nil {
		panic("nil users")
	}
	if sessions == nil {
		panic("nil sessions")
	}
	if secret == "" {
		panic("empty jwt secret")
	}

	return decorator.ApplyQueryDecorators(
		authenticateHandler{users: users, sessions: sessions, secret: secret},
		logger,
		metricsClient,
	)
}

func (h authenticateHandler) Handle(ctx context.Context, q Authenticate) (auth.User, error) {
	principal, err := auth.Parse(h.secret, q.Token)
	if err != nil {
		return auth.User{}, err
	}

	session, err := h.sessions.GetSession(ctx, principal.SessionID)
	if err != nil {
		if errors.Is(err, user.ErrSessionNotFound) {
			return auth.User{}, sessionInvalid()
		}
		return auth.User{}, err
	}
	if session.UserID() != principal.ID || session.Expired(time.Now().UTC()) {
		if session.Expired(time.Now().UTC()) {
			if err := h.sessions.DeleteSession(ctx, session.ID()); err != nil {
				return auth.User{}, err
			}
		}
		return auth.User{}, sessionInvalid()
	}

	account, err := h.users.GetUserByID(ctx, principal.ID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return auth.User{}, sessionInvalid()
		}
		return auth.User{}, err
	}

	return auth.User{
		ID:          account.ID(),
		Email:       account.Email(),
		Username:    account.Username(),
		Role:        string(account.Role()),
		Permissions: account.Permissions(),
		SessionID:   session.ID(),
	}, nil
}

func sessionInvalid() error {
	return commonerrors.NewAuthorizationError("session is no longer valid", "session-invalid")
}
