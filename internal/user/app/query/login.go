package query

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/common/auth"
	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/user/domain/user"
)

type Login struct {
	Email    string
	Password string
}

func (q Login) GoString() string {
	return fmt.Sprintf("query.Login{Email:%q, Password:\"***\"}", q.Email)
}

type LoginResult struct {
	AccessToken string
	ExpiresIn   int
	TokenType   string
	User        User
}

type LoginHandler decorator.QueryHandler[Login, LoginResult]

type loginHandler struct {
	users    user.Repository
	sessions user.SessionRepository
	secret   string
	ttl      time.Duration
}

func NewLoginHandler(
	users user.Repository,
	sessions user.SessionRepository,
	secret string,
	ttl time.Duration,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) LoginHandler {
	if users == nil {
		panic("nil users")
	}
	if sessions == nil {
		panic("nil sessions")
	}
	if secret == "" {
		panic("empty jwt secret")
	}
	if ttl <= 0 {
		panic("non-positive token ttl")
	}

	return decorator.ApplyQueryDecorators(
		loginHandler{
			users:    users,
			sessions: sessions,
			secret:   secret,
			ttl:      ttl,
		},
		logger,
		metricsClient,
	)
}

func (h loginHandler) Handle(ctx context.Context, q Login) (LoginResult, error) {
	account, err := h.users.GetUserByEmail(ctx, q.Email)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) || errors.Is(err, user.ErrInvalidEmail) {
			return LoginResult{}, user.ErrInvalidCredentials
		}
		return LoginResult{}, err
	}
	if err := account.CheckPassword(q.Password); err != nil {
		return LoginResult{}, err
	}

	session, err := user.NewSession(uuid.NewString(), account.ID(), h.ttl)
	if err != nil {
		return LoginResult{}, err
	}
	if err := h.sessions.AddSession(ctx, session); err != nil {
		return LoginResult{}, err
	}

	principal := auth.User{
		ID:          account.ID(),
		Email:       account.Email(),
		Username:    account.Username(),
		Role:        string(account.Role()),
		Permissions: account.Permissions(),
	}
	token, expiresIn, err := auth.Issue(h.secret, principal, session.ID(), h.ttl)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		AccessToken: token,
		ExpiresIn:   expiresIn,
		TokenType:   auth.TokenTypeBearer,
		User:        UserFromDomain(account),
	}, nil
}
