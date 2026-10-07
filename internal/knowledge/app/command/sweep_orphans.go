package command

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/indexing"
	"github.com/sirupsen/logrus"
)

type SweepOrphans struct {
}

type SweepOrphansHandler decorator.CommandHandler[SweepOrphans]

type sweepOrphansHandler struct {
	repo indexing.Repository
}

func NewSweepOrphansHandler(
	repo indexing.Repository,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) SweepOrphansHandler {
	if repo == nil {
		panic("nil repo service")
	}

	return decorator.ApplyCommandDecorators(
		sweepOrphansHandler{repo: repo},
		logger,
		metricsClient,
	)
}

func (h sweepOrphansHandler) Handle(ctx context.Context, cmd SweepOrphans) error {
	return h.repo.DeleteOrphans(ctx)
}
