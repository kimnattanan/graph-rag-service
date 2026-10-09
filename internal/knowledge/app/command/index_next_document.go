package command

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/indexing"
	"github.com/sirupsen/logrus"
)

type IndexNextDocument struct {
	NoPending *bool
}

type IndexNextDocumentHandler decorator.CommandHandler[IndexNextDocument]

type indexNextDocumentHandler struct {
	repo      document.Repository
	extractor indexing.Extractor
	embedder  indexing.Embedder
	indexRepo indexing.Repository
}

func NewIndexNextDocumentHandler(
	repo document.Repository,
	extractor indexing.Extractor,
	embedder indexing.Embedder,
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
		indexNextDocumentHandler{repo: repo, extractor: extractor, embedder: embedder, indexRepo: indexRepo},
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
	chunkContents := make([]string, 0, len(results))
	for _, result := range results {
		chunkContents = append(chunkContents, result.Content())
	}
	chunkEmbeddings, err := h.embedder.Embed(ctx, chunkContents)
	if err != nil {
		return err
	}
	if err := h.indexRepo.ReplaceChunks(ctx, doc.ID(), results, chunkEmbeddings); err != nil {
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
