package command

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/indexing"
	"github.com/sirupsen/logrus"
)

type IndexNextDocument struct {
	// NoPending is set when the queue has no document to claim.
	// An empty queue is not a failure, so Handle returns nil and the worker stops on this flag.
	NoPending *bool
}

type IndexNextDocumentHandler decorator.CommandHandler[IndexNextDocument]

type indexNextDocumentHandler struct {
	repo      document.Repository
	extractor indexing.Extractor
	indexRepo indexing.Repository
}

func NewIndexNextDocumentHandler(
	repo document.Repository,
	extractor indexing.Extractor,
	indexRepo indexing.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) IndexNextDocumentHandler {
	if repo == nil {
		panic("nil repo service")
	}
	if extractor == nil {
		panic("nil extractor service")
	}
	if indexRepo == nil {
		panic("nil index repo service")
	}

	return decorator.ApplyCommandDecorators(
		indexNextDocumentHandler{repo: repo, extractor: extractor, indexRepo: indexRepo},
		logger,
		metricsClient,
	)
}

func (h indexNextDocumentHandler) Handle(ctx context.Context, cmd IndexNextDocument) error {
	doc, err := h.repo.ClaimNextPending(ctx)
	if err != nil {
		return err
	}
	if doc == nil {
		if cmd.NoPending != nil {
			*cmd.NoPending = true
		}
		return nil
	}
	results, err := h.extractor.Extract(ctx, doc.Content())
	if err != nil {
		if updateErr := h.repo.UpdateDocument(ctx, doc.ID(), func(doc *document.Document) error {
			return doc.MarkAsFailed(err.Error())
		}); updateErr != nil {
			return updateErr
		}
		return err
	}
	if err := h.indexRepo.ReplaceChunks(ctx, doc.ID(), results); err != nil {
		if updateErr := h.repo.UpdateDocument(ctx, doc.ID(), func(doc *document.Document) error {
			return doc.MarkAsFailed(err.Error())
		}); updateErr != nil {
			return updateErr
		}
		return err
	}
	if updateErr := h.repo.UpdateDocument(ctx, doc.ID(), func(doc *document.Document) error {
		return doc.MarkAsCompleted()
	}); updateErr != nil {
		return updateErr
	}
	return nil
}
