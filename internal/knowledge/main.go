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
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/config"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/ports"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/service"
)

func main() {
	cfg := config.Config{}
	if err := commonConfig.NewConfig(&cfg); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	logs.Init(&cfg.Common)

	ctx := context.Background()

	application := service.NewApplication(ctx, &cfg)

	serverType := strings.ToLower(cfg.App.ServerToRun)
	switch serverType {
	case "http":
		// go loadFixtures(application)

		server.RunHTTPServer(&cfg.Common, func(router chi.Router) http.Handler {
			return ports.HandlerFromMux(
				ports.NewHttpServer(application),
				router,
			)
		})
	// case "grpc":
	// 	server.RunGRPCServer(func(server *grpc.Server) {
	// 		svc := ports.NewGrpcServer(application)
	// 		trainer.RegisterTrainerServiceServer(server, svc)
	// 	})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
