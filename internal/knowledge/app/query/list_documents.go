package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
	"github.com/sirupsen/logrus"
)

type ListDocuments struct {
	Limit       *int
	Offset      *int
	Tag         *string
	IndexStatus *document.IndexStatus
}

type ListDocumentsHandler decorator.QueryHandler[ListDocuments, DocumentSummaryList]

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
	ListDocuments(ctx context.Context, limit int, offset int, tag *string, indexStatus *document.IndexStatus) (*DocumentSummaryList, error)
}

func (h listDocumentsHandler) Handle(ctx context.Context, query ListDocuments) (DocumentSummaryList, error) {
	limit := 20
	offset := 0
	if query.Limit != nil {
		limit = *query.Limit
	}
	if query.Offset != nil {
		offset = *query.Offset
	}
	if limit < 1 || limit > 100 {
		return DocumentSummaryList{}, commonerrors.NewIncorrectInputError(
			"Limit must be between 1 and 100",
			"limit-invalid",
		)
	}
	if offset < 0 {
		return DocumentSummaryList{}, commonerrors.NewIncorrectInputError(
			"Offset must be greater than 0",
			"offset-invalid",
		)
	}
	docs, err := h.readModel.ListDocuments(ctx, limit, offset, query.Tag, query.IndexStatus)
	if err != nil {
		return DocumentSummaryList{}, err
	}
	return *docs, nil
}
