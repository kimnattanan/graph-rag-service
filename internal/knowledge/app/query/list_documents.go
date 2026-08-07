package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/sirupsen/logrus"
)

type ListDocuments struct{}

type ListDocumentsHandler decorator.QueryHandler[ListDocuments, []Document]

type listDocumentsHandler struct {
	readModel ListDocumentsReadModel
}

func NewListDocumentsHandler(
	readModel ListDocumentsReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) listDocumentsHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator
}

type ListDocumentsReadModel interface {
	ListDocuments(ctx context.Context) ([]Document, error)
}

func (h listDocumentsHandler) Handle(ctx context.Context, _ ListDocuments) (docs []Document, err error) {
	return h.readModel.ListDocuments(ctx)
}
