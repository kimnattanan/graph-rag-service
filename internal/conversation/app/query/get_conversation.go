package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/conversation"
	"github.com/sirupsen/logrus"
)

type GetConversation struct {
	ConversationID string
	UserID         string
}

type GetConversationHandler decorator.QueryHandler[GetConversation, Conversation]

type getConversationHandler struct {
	readModel GetConversationReadModel
}

func NewGetConversationHandler(
	readModel GetConversationReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetConversationHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators(
		getConversationHandler{readModel},
		logger,
		metricsClient,
	)
}

type GetConversationReadModel interface {
	GetConversation(ctx context.Context, userID string, conversationID string) (*conversation.Conversation, error)
}

func (h getConversationHandler) Handle(ctx context.Context, query GetConversation) (Conversation, error) {
	conversation, err := h.readModel.GetConversation(ctx, query.UserID, query.ConversationID)
	if err != nil {
		return Conversation{}, err
	}
	messages := make([]Message, len(conversation.Messages()))
	for i, message := range conversation.Messages() {
		messages[i] = Message{
			ID:             message.ID(),
			ConversationID: message.ConversationID(),
		}
	}
	return Conversation{
		ID:        conversation.ID(),
		UserID:    conversation.UserID(),
		Title:     conversation.Title(),
		CreatedAt: conversation.CreatedAt(),
		UpdatedAt: conversation.UpdatedAt(),
		Messages:  messages,
	}, nil
}
