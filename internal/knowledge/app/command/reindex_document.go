package command

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
	"github.com/sirupsen/logrus"
)

type ReindexDocument struct {
	DocumentID string
}

type ReindexDocumentHandler decorator.CommandHandler[ReindexDocument]

type reindexDocumentHandler struct {
	repo document.Repository
}

func NewReindexDocumentHandler(
	repo document.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) ReindexDocumentHandler {
	if repo == nil {
		panic("nil repo service")
	}

	return decorator.ApplyCommandDecorators(
		reindexDocumentHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h reindexDocumentHandler) Handle(ctx context.Context, cmd ReindexDocument) error {
	return h.repo.UpdateDocument(ctx, cmd.DocumentID, func(doc *document.Document) error {
		if err := doc.MarkAsPending(); err != nil {
			return err
		}
		return nil
	})
}
