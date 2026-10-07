package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/sirupsen/logrus"
)

type ListConversations struct {
	UserID string
	Limit  *int
	Offset *int
}

type ListConversationsHandler decorator.QueryHandler[ListConversations, ConversationSummaryList]

type listConversationsHandler struct {
	readModel ListConversationsReadModel
}

func NewListConversationsHandler(
	readModel ListConversationsReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) ListConversationsHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators(
		listConversationsHandler{readModel},
		logger,
		metricsClient,
	)
}

type ListConversationsReadModel interface {
	ListConversations(ctx context.Context, userID string, limit int, offset int) (*ConversationSummaryList, error)
}

func (h listConversationsHandler) Handle(ctx context.Context, query ListConversations) (ConversationSummaryList, error) {
	limit := 20
	if query.Limit != nil {
		limit = *query.Limit
	}
	offset := 0
	if query.Offset != nil {
		offset = *query.Offset
	}
	if limit < 1 || limit > 100 {
		return ConversationSummaryList{}, commonerrors.NewIncorrectInputError(
			"Limit must be between 1 and 100",
			"limit-invalid",
		)
	}
	if offset < 0 {
		return ConversationSummaryList{}, commonerrors.NewIncorrectInputError(
			"Offset must be greater than 0",
			"offset-invalid",
		)
	}
	conversations, err := h.readModel.ListConversations(ctx, query.UserID, limit, offset)
	if err != nil {
		return ConversationSummaryList{}, err
	}
	return *conversations, nil
}
