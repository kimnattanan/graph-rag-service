package service

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/config"
)

func NewApplication(ctx context.Context, cfg *config.Config) app.Application {
	dbUri := fmt.Sprintf("bolt://%s:%s", cfg.Memgraph.Host, cfg.Memgraph.Port)
	driver, err := neo4j.NewDriverWithContext(dbUri, neo4j.BasicAuth(cfg.Memgraph.User, cfg.Memgraph.Port, ""))
	if err != nil {
		panic(err)
	}
	defer driver.Close(ctx)

	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	return app.Application{
		Commands: app.Commands{},
		Queries:  app.Queries{},
	}
}
