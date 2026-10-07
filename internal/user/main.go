package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	commonConfig "github.com/kimnattanan/graph-rag-service/internal/common/config"
	"github.com/kimnattanan/graph-rag-service/internal/common/logs"
	"github.com/kimnattanan/graph-rag-service/internal/common/server"
	"github.com/kimnattanan/graph-rag-service/internal/user/config"
	"github.com/kimnattanan/graph-rag-service/internal/user/ports"
	"github.com/kimnattanan/graph-rag-service/internal/user/service"
)

func main() {
	cfg := config.Config{}
	if err := commonConfig.NewConfig(&cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logs.Init(&cfg.Common)

	ctx := context.Background()
	application, cleanup := service.NewApplication(ctx, &cfg)
	defer cleanup()

	serverType := strings.ToLower(cfg.App.ServerToRun)
	switch serverType {
	case "http":
		server.RunHTTPServer(&cfg.Common, func(router chi.Router) http.Handler {
			return ports.HandlerFromMux(ports.NewHttpServer(application), router)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
