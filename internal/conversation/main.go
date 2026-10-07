package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc"

	commonConfig "github.com/kimnattanan/graph-rag-service/internal/common/config"
	conversationpb "github.com/kimnattanan/graph-rag-service/internal/common/genproto/conversation"
	"github.com/kimnattanan/graph-rag-service/internal/common/logs"
	"github.com/kimnattanan/graph-rag-service/internal/common/server"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/config"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/ports"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/service"
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
			return ports.HandlerFromMux(
				ports.NewHttpServer(application),
				router,
			)
		})
	case "grpc":
		server.RunGRPCServer(&cfg.Common, func(grpcServer *grpc.Server) {
			svc := ports.NewGrpcServer(application)
			conversationpb.RegisterConversationServiceServer(grpcServer, svc)
		})
	default:
		panic(fmt.Sprintf("server type '%s' is not supported", serverType))
	}
}
