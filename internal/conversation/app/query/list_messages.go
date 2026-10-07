package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/sirupsen/logrus"
)

type ListMessages struct {
	ConversationID string
	UserID         string
	Limit          *int
	Offset         *int
}

type ListMessagesHandler decorator.QueryHandler[ListMessages, MessageList]

type listMessagesHandler struct {
	readModel ListMessagesReadModel
}

func NewListMessagesHandler(
	readModel ListMessagesReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) ListMessagesHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators(
		listMessagesHandler{readModel},
		logger,
		metricsClient,
	)
}

type ListMessagesReadModel interface {
	ListMessages(ctx context.Context, userID string, conversationID string, limit int, offset int) (*MessageList, error)
}

func (h listMessagesHandler) Handle(ctx context.Context, query ListMessages) (MessageList, error) {
	return MessageList{}, nil
}
