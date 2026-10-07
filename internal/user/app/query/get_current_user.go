package query

import (
	"context"

	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/user/domain/user"
)

type GetCurrentUser struct {
	UserID string
}

type GetCurrentUserHandler decorator.QueryHandler[GetCurrentUser, User]

type getCurrentUserHandler struct {
	repo user.Repository
}

func NewGetCurrentUserHandler(
	repo user.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetCurrentUserHandler {
	if repo == nil {
		panic("nil repo")
	}

	return decorator.ApplyQueryDecorators(
		getCurrentUserHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h getCurrentUserHandler) Handle(ctx context.Context, q GetCurrentUser) (User, error) {
	if q.UserID == "" {
		return User{}, user.ErrEmptyUserID
	}
	account, err := h.repo.GetUserByID(ctx, q.UserID)
	if err != nil {
		return User{}, err
	}
	return UserFromDomain(account), nil
}
