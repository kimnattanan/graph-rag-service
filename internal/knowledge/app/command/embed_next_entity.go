package command

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/indexing"
	"github.com/sirupsen/logrus"
)

type EmbedNextEntity struct {
	NoPending *bool
}

type EmbedNextEntityHandler decorator.CommandHandler[EmbedNextEntity]

type embedNextEntityHandler struct {
	repo     indexing.Repository
	embedder indexing.Embedder
}

func NewEmbedNextEntityHandler(
	repo indexing.Repository,
	embedder indexing.Embedder,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) EmbedNextEntityHandler {
	if repo == nil {
		panic("nil index repo service")
	}
	if embedder == nil {
		panic("nil embedder service")
	}

	return decorator.ApplyCommandDecorators(
		embedNextEntityHandler{repo: repo, embedder: embedder},
		logger,
		metricsClient,
	)
}

func (h embedNextEntityHandler) Handle(ctx context.Context, cmd EmbedNextEntity) error {
	entity, err := h.repo.ClaimNextPendingEntity(ctx)
	if err != nil {
		return err
	}
	if entity == nil {
		if cmd.NoPending != nil {
			*cmd.NoPending = true
		}
		return nil
	}
	embedding, err := h.embedder.Embed(ctx, []string{entity.Name()})
	if err != nil {
		return err
	}
	if len(embedding) == 0 {
		return commonerrors.NewSlugError("embedding is empty", "embedding-empty")
	}
	return h.repo.EmbedEntity(ctx, entity.Name(), embedding[0])
}
