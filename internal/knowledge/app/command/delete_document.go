package command

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
	"github.com/sirupsen/logrus"
)

type DeleteDocument struct {
	DocumentID string
}

type DeleteDocumentHandler decorator.CommandHandler[DeleteDocument]

type deleteDocumentHandler struct {
	repo document.Repository
}

func NewDeleteDocumentHandler(
	repo document.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) DeleteDocumentHandler {
	if repo == nil {
		panic("nil repo service")
	}

	return decorator.ApplyCommandDecorators(
		deleteDocumentHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h deleteDocumentHandler) Handle(ctx context.Context, cmd DeleteDocument) error {
	return h.repo.DeleteDocument(ctx, cmd.DocumentID)
}
