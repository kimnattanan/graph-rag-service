package command

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/user/domain/user"
)

type Register struct {
	Email    string
	Username string
	Password string
}

func (c Register) GoString() string {
	return fmt.Sprintf("command.Register{Email:%q, Username:%q, Password:\"***\"}", c.Email, c.Username)
}

type RegisterHandler decorator.CommandHandler[Register]

type registerHandler struct {
	repo user.Repository
}

func NewRegisterHandler(
	repo user.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) RegisterHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators(
		registerHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h registerHandler) Handle(ctx context.Context, cmd Register) error {
	account, err := user.NewUser(uuid.NewString(), cmd.Email, cmd.Username, cmd.Password)
	if err != nil {
		return err
	}
	return h.repo.AddUser(ctx, account)
}
