package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
	"github.com/sirupsen/logrus"
)

type GetDocument struct {
	DocumentID string
}

type GetDocumentHandler decorator.QueryHandler[GetDocument, Document]

type getDocumentHandler struct {
	readModel GetDocumentReadModel
}

func NewGetDocumentHandler(
	readModel GetDocumentReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetDocumentHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators(
		getDocumentHandler{readModel},
		logger,
		metricsClient,
	)
}

type GetDocumentReadModel interface {
	GetDocument(ctx context.Context, documentID string) (*document.Document, error)
}

func (h getDocumentHandler) Handle(ctx context.Context, query GetDocument) (Document, error) {
	doc, err := h.readModel.GetDocument(ctx, query.DocumentID)
	if err != nil {
		return Document{}, err
	}
	indexState := doc.IndexState()
	return Document{
		ID: doc.ID(),
		Title: doc.Title(),
		Content: doc.Content(),
		Tags: doc.Tags(),
		IndexStatus: indexState.Status(),
		CreatedAt: doc.CreatedAt(),
		UpdatedAt: doc.UpdatedAt(),
	}, nil
}
