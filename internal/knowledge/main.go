package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	commonConfig "github.com/kimnattanan/graph-rag-service/internal/common/config"
	"github.com/kimnattanan/graph-rag-service/internal/common/server"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/config"
)

func main() {
	// ...
	cfg := &config.Config{}
	if err := commonConfig.NewConfig(&cfg.Common); err != nil {
		log.Fatalf("Failed to load common config: %v", err)
	}
	server.RunHTTPServer(&cfg.Common, func(router chi.Router) http.Handler {
		return HandlerFromMux(HttpServer{firebaseDB, trainerClient, usersClient}, router)
	})
}
