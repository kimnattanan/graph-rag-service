package command

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/conversation"
	"github.com/sirupsen/logrus"
)

type DeleteConversation struct {
	ConversationID string
	UserID         string
}

type DeleteConversationHandler decorator.CommandHandler[DeleteConversation]

type deleteConversationHandler struct {
	repo conversation.Repository
}

func NewDeleteConversationHandler(
	repo conversation.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) DeleteConversationHandler {
	if repo == nil {
		panic("nil repo service")
	}

	return decorator.ApplyCommandDecorators(
		deleteConversationHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h deleteConversationHandler) Handle(ctx context.Context, cmd DeleteConversation) error {
	return nil
}
