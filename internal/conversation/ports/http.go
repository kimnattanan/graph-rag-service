package ports

import (
	"net/http"

	"github.com/kimnattanan/graph-rag-service/internal/conversation/app"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(application app.Application) HttpServer {
	return HttpServer{app: application}
}

func (h HttpServer) ListConversations(w http.ResponseWriter, r *http.Request, params ListConversationsParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h HttpServer) CreateConversation(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h HttpServer) DeleteConversation(w http.ResponseWriter, r *http.Request, conversationId openapi_types.UUID) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h HttpServer) GetConversation(w http.ResponseWriter, r *http.Request, conversationId openapi_types.UUID) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h HttpServer) ListMessages(w http.ResponseWriter, r *http.Request, conversationId openapi_types.UUID, params ListMessagesParams) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h HttpServer) SendMessage(w http.ResponseWriter, r *http.Request, conversationId openapi_types.UUID) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (h HttpServer) GetMessage(w http.ResponseWriter, r *http.Request, conversationId openapi_types.UUID, messageId openapi_types.UUID) {
	w.WriteHeader(http.StatusNotImplemented)
}
