package command

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
	"github.com/sirupsen/logrus"
)

type CreateDocument struct {
	DocumentID string
	Title      string
	Content    string
	Tags       []string
}

type CreateDocumentHandler decorator.CommandHandler[CreateDocument]

type createDocumentHandler struct {
	repo document.Repository
}

func NewCreateDocumentHandler(
	repo document.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) CreateDocumentHandler {
	if repo == nil {
		panic("nil repo service")
	}

	return decorator.ApplyCommandDecorators(
		createDocumentHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h createDocumentHandler) Handle(ctx context.Context, cmd CreateDocument) error {
	doc, err := document.NewDocument(cmd.DocumentID, cmd.Title, cmd.Content, cmd.Tags)
	if err != nil {
		return err
	}
	return h.repo.AddDocument(ctx, doc)
}
