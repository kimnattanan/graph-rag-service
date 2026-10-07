package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/sirupsen/logrus"
)

type GetMessage struct {
	ConversationID string
	MessageID      string
	UserID         string
}

type GetMessageHandler decorator.QueryHandler[GetMessage, Message]

type getMessageHandler struct {
	readModel GetMessageReadModel
}

func NewGetMessageHandler(
	readModel GetMessageReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetMessageHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators(
		getMessageHandler{readModel},
		logger,
		metricsClient,
	)
}

type GetMessageReadModel interface {
	GetMessage(ctx context.Context, userID string, conversationID string, messageID string) (*Message, error)
}

func (h getMessageHandler) Handle(ctx context.Context, query GetMessage) (Message, error) {
	return Message{}, nil
}
