package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/sirupsen/logrus"
)

type ListDocuments struct{}

type ListDocumentsHandler decorator.QueryHandler[ListDocuments, []*DocumentSummary]

type listDocumentsHandler struct {
	readModel ListDocumentsReadModel
}

func NewListDocumentsHandler(
	readModel ListDocumentsReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) ListDocumentsHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators(
		listDocumentsHandler{readModel},
		logger,
		metricsClient,
	)
}

type ListDocumentsReadModel interface {
	ListDocuments(ctx context.Context) ([]*DocumentSummary, error)
}

func (h listDocumentsHandler) Handle(ctx context.Context, _ ListDocuments) (docs []*DocumentSummary, err error) {
	return h.readModel.ListDocuments(ctx)
}
