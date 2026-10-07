package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
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
	limit := 50
	if query.Limit != nil {
		limit = *query.Limit
	}
	offset := 0
	if query.Offset != nil {
		offset = *query.Offset
	}
	if limit < 1 || limit > 100 {
		return MessageList{}, commonerrors.NewIncorrectInputError(
			"Limit must be between 1 and 100",
			"limit-invalid",
		)
	}
	if offset < 0 {
		return MessageList{}, commonerrors.NewIncorrectInputError(
			"Offset must be greater than 0",
			"offset-invalid",
		)
	}
	messages, err := h.readModel.ListMessages(ctx, query.UserID, query.ConversationID, limit, offset)
	if err != nil {
		return MessageList{}, err
	}
	return *messages, nil
}
