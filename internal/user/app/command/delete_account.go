package command

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/user/domain/user"
)

type DeleteAccount struct {
	UserID string
}

type DeleteAccountHandler decorator.CommandHandler[DeleteAccount]

type deleteAccountHandler struct {
	repo user.Repository
}

func NewDeleteAccountHandler(
	repo user.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) DeleteAccountHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyCommandDecorators(
		deleteAccountHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h deleteAccountHandler) Handle(ctx context.Context, cmd DeleteAccount) error {
	if cmd.UserID == "" {
		return user.ErrEmptyUserID
	}
	return h.repo.DeleteUser(ctx, cmd.UserID)
}
