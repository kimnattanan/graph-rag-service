package command

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/conversation"
	"github.com/sirupsen/logrus"
)

type CreateConversation struct {
	ConversationID string
	UserID         string
	Title          string
}

type CreateConversationHandler decorator.CommandHandler[CreateConversation]

type createConversationHandler struct {
	repo conversation.Repository
}

func NewCreateConversationHandler(
	repo conversation.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CreateConversationHandler {
	if repo == nil {
		panic("nil repo service")
	}

	return decorator.ApplyCommandDecorators(
		createConversationHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h createConversationHandler) Handle(ctx context.Context, cmd CreateConversation) error {
	return nil
}
