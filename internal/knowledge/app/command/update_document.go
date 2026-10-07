package command

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
	"github.com/sirupsen/logrus"
)

type UpdateDocument struct {
	DocumentID string
	Title      *string
	Content    *string
	Tags       *[]string
}

type UpdateDocumentHandler decorator.CommandHandler[UpdateDocument]

type updateDocumentHandler struct {
	repo document.Repository
}

func NewUpdateDocumentHandler(
	repo document.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) UpdateDocumentHandler {
	if repo == nil {
		panic("nil repo service")
	}

	return decorator.ApplyCommandDecorators(
		updateDocumentHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h updateDocumentHandler) Handle(ctx context.Context, cmd UpdateDocument) error {
	return h.repo.UpdateDocument(ctx, cmd.DocumentID, func(doc *document.Document) error {
		if cmd.Title != nil {
			if err := doc.UpdateTitle(*cmd.Title); err != nil {
				return err
			}
		}
		if cmd.Content != nil {
			if err := doc.MarkAsPending(); err != nil {
				return err
			}
			if err := doc.UpdateContent(*cmd.Content); err != nil {
				return err
			}
		}
		if cmd.Tags != nil {
			doc.UpdateTags(*cmd.Tags)
		}
		return nil
	})
}
