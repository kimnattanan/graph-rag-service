package service

import (
	"context"
	"log"

	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/common/metrics"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/adapters"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/config"
)

func NewApplication(ctx context.Context, cfg *config.Config) (app.Application, func()) {
	memgraphDriver, err := connectMemgraph(ctx, cfg)
	if err != nil {
		panic(err)
	}
	if err := startupMemgraph(ctx, memgraphDriver, cfg); err != nil {
		_ = memgraphDriver.Close(ctx)
		panic(err)
	}

	documentMemgraphReadModel := adapters.NewDocumentMemgraphReadModel(memgraphDriver)
	documentMemgraphRepository := adapters.NewDocumentMemgraphRepository(memgraphDriver)
	indexMemgraphRepository := adapters.NewIndexMemgraphRepository(memgraphDriver)
	extractor := adapters.NewExtractor(cfg.LLM)
	embedder := adapters.NewEmbedder(cfg.Embedder)

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	return app.Application{
			Commands: app.Commands{
				CreateDocument:  command.NewCreateDocumentHandler(documentMemgraphRepository, logger, metricsClient),
				UpdateDocument:  command.NewUpdateDocumentHandler(documentMemgraphRepository, logger, metricsClient),
				DeleteDocument:  command.NewDeleteDocumentHandler(documentMemgraphRepository, logger, metricsClient),
				ReindexDocument: command.NewReindexDocumentHandler(documentMemgraphRepository, logger, metricsClient),
				IndexNextDocument: command.NewIndexNextDocumentHandler(
					documentMemgraphRepository,
					extractor,
					embedder,
					indexMemgraphRepository,
					logger,
					metricsClient,
				),
				EmbedNextEntity: command.NewEmbedNextEntityHandler(indexMemgraphRepository, embedder, logger, metricsClient),
				SweepOrphans: command.NewSweepOrphansHandler(indexMemgraphRepository, logger, metricsClient),
			},
			Queries: app.Queries{
				ListDocuments:          query.NewListDocumentsHandler(documentMemgraphReadModel, logger, metricsClient),
				GetDocument:            query.NewGetDocumentHandler(documentMemgraphReadModel, logger, metricsClient),
				GetDocumentIndexStatus: query.NewGetDocumentIndexStatusHandler(documentMemgraphReadModel, logger, metricsClient),
				Retrieve:               query.NewRetrieveHandler(documentMemgraphReadModel, embedder, logger, metricsClient),
			},
		}, func() {
			err := memgraphDriver.Close(ctx)
			if err != nil {
				log.Printf("error closing memgraph driver: %v \n", err)
			}
		}
}
