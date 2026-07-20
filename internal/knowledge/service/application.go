package service

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app"
)

func NewApplication(ctx context.Context) app.Application {
	return app.Application{
		Commands: app.Commands{},
		Queries:  app.Queries{},
	}
}
