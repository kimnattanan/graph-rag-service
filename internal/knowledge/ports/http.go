package ports

import (
	"net/http"

	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app}
}

func (h HttpServer) ListDocuments(w http.ResponseWriter, r *http.Request, params ListDocumentsParams) {

}

func (h HttpServer) CreateDocument(w http.ResponseWriter, r *http.Request) {

}

func (h HttpServer) DeleteDocument(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {

}

func (h HttpServer) GetDocument(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {

}

func (h HttpServer) UpdateDocument(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {

}

func (h HttpServer) GetDocumentIndexStatus(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {

}

func (h HttpServer) ReindexDocument(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {

}

func (h HttpServer) Retrieve(w http.ResponseWriter, r *http.Request) {

}
