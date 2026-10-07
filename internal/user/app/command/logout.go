package command

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/common/auth"
	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/user/domain/user"
)

type Logout struct {
	Token string
}

func (c Logout) GoString() string {
	return "command.Logout{Token:\"***\"}"
}

type LogoutHandler decorator.CommandHandler[Logout]

type logoutHandler struct {
	sessions user.SessionRepository
	secret   string
}

func NewLogoutHandler(
	sessions user.SessionRepository,
	secret string,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) LogoutHandler {
	if sessions == nil {
		panic("nil sessions")
	}
	if secret == "" {
		panic("empty jwt secret")
	}

	return decorator.ApplyCommandDecorators(
		logoutHandler{sessions: sessions, secret: secret},
		logger,
		metricsClient,
	)
}

func (h logoutHandler) Handle(ctx context.Context, cmd Logout) error {
	principal, err := auth.Parse(h.secret, cmd.Token)
	if err != nil {
		return err
	}
	return h.sessions.DeleteSession(ctx, principal.SessionID)
}
