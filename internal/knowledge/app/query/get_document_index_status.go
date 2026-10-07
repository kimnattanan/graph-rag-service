package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/sirupsen/logrus"
)

type GetDocumentIndexStatus struct {
	DocumentID string
}

type GetDocumentIndexStatusHandler decorator.QueryHandler[GetDocumentIndexStatus, DocumentIndexStatus]

type getDocumentIndexStatusHandler struct {
	readModel GetDocumentReadModel
}

func NewGetDocumentIndexStatusHandler(
	readModel GetDocumentReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) GetDocumentIndexStatusHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators(
		getDocumentIndexStatusHandler{readModel},
		logger,
		metricsClient,
	)
}

func (h getDocumentIndexStatusHandler) Handle(ctx context.Context, query GetDocumentIndexStatus) (DocumentIndexStatus, error) {
	doc, err := h.readModel.GetDocument(ctx, query.DocumentID)
	if err != nil {
		return DocumentIndexStatus{}, err
	}
	indexState := doc.IndexState()
	startedAtValue := indexState.StartedAt()
	finishedAtValue := indexState.FinishedAt()
	errorMessageValue := indexState.ErrorMessage()
	startedAt := &startedAtValue
	finishedAt := &finishedAtValue
	errorMessage := &errorMessageValue
	if startedAtValue.IsZero() {
		startedAt = nil
	}
	if finishedAtValue.IsZero() {
		finishedAt = nil
	}
	if errorMessageValue == "" {
		errorMessage = nil
	}
	return DocumentIndexStatus{
		DocumentID:   doc.ID(),
		Status:       indexState.Status(),
		StartedAt:    startedAt,
		FinishedAt:   finishedAt,
		ErrorMessage: errorMessage,
	}, nil
}
