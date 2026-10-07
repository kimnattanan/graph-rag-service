package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
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
	return ConversationSummaryList{}, nil
}
