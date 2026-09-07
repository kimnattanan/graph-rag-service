package service

import (
	"context"
	"fmt"
	"log"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/common/metrics"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/adapters"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/config"
)

func NewApplication(ctx context.Context, cfg *config.Config) (app.Application, func()) {
	dbUri := fmt.Sprintf("bolt://%s:%s", cfg.Memgraph.Host, cfg.Memgraph.Port)
	memgraphDriver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(cfg.Memgraph.User, cfg.Memgraph.Password, ""))
	if err != nil {
		panic(err)
	}
	err = memgraphDriver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	knowledgeRepository := adapters.NewKnowledgeRepository(memgraphDriver)

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	return app.Application{
			Commands: app.Commands{},
			Queries: app.Queries{
				ListDocuments: query.NewListDocumentsHandler(knowledgeRepository, logger, metricsClient),
			},
		}, func() {
			err := memgraphDriver.Close(ctx)
			if err != nil {
				log.Printf("error closing memgraph driver: %v \n", err)
			}
		}
}
