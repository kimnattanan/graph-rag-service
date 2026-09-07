package ports

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/common/server/httperr"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/query"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app}
}

func (h HttpServer) ListDocuments(w http.ResponseWriter, r *http.Request, params ListDocumentsParams) {
	appDocs, err := h.app.Queries.ListDocuments.Handle(r.Context(), query.ListDocuments{})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	documentSummaries := make([]DocumentSummary, len(appDocs))
	for i, appDoc := range appDocs {
		summary, err := appDocumentSummaryToResponse(appDoc)
		if err != nil {
			httperr.RespondWithSlugError(err, w, r)
			return
		}
		documentSummaries[i] = *summary
	}

	render.Respond(w, r, DocumentList{
		Items: documentSummaries,
		Total: len(documentSummaries),
	})
}

func (h HttpServer) CreateDocument(w http.ResponseWriter, r *http.Request) {
	httperr.RespondWithSlugError(errors.NewSlugError("Not implemented", "not-implemented"), w, r)
}

func (h HttpServer) DeleteDocument(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {
	httperr.RespondWithSlugError(errors.NewSlugError("Not implemented", "not-implemented"), w, r)
}

func (h HttpServer) GetDocument(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {
	httperr.RespondWithSlugError(errors.NewSlugError("Not implemented", "not-implemented"), w, r)
}

func (h HttpServer) UpdateDocument(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {
	httperr.RespondWithSlugError(errors.NewSlugError("Not implemented", "not-implemented"), w, r)
}

func (h HttpServer) GetDocumentIndexStatus(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {
	httperr.RespondWithSlugError(errors.NewSlugError("Not implemented", "not-implemented"), w, r)
}

func (h HttpServer) ReindexDocument(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {
	httperr.RespondWithSlugError(errors.NewSlugError("Not implemented", "not-implemented"), w, r)
}

func (h HttpServer) Retrieve(w http.ResponseWriter, r *http.Request) {
	httperr.RespondWithSlugError(errors.NewSlugError("Not implemented", "not-implemented"), w, r)
}

// Mapper functions

func parseIndexStatus(status query.IndexStatus) (IndexStatus, error) {
	switch status {
	case query.IndexStatusCompleted:
		return Completed, nil
	case query.IndexStatusFailed:
		return Failed, nil
	case query.IndexStatusPending:
		return Pending, nil
	default:
		return "", errors.NewIncorrectInputError("Invalid index status", "index-status-invalid")
	}
}

func appDocumentToResponse(d *query.Document) (*Document, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil {
		return nil, errors.NewIncorrectInputError("Invalid document ID", "document-id-invalid")
	}
	indexStatus, err := parseIndexStatus(d.IndexStatus)
	if err != nil {
		return nil, err
	}
	return &Document{
		Id:          id,
		Title:       d.Title,
		Content:     d.Content,
		Tags:        d.Tags,
		IndexStatus: indexStatus,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}, nil
}

func appDocumentSummaryToResponse(s *query.DocumentSummary) (*DocumentSummary, error) {
	id, err := uuid.Parse(s.ID)
	if err != nil {
		return nil, errors.NewIncorrectInputError("Invalid document ID", "document-id-invalid")
	}
	indexStatus, err := parseIndexStatus(s.IndexStatus)
	if err != nil {
		return nil, err
	}
	return &DocumentSummary{
		Id:          id,
		Title:       s.Title,
		Tags:        s.Tags,
		IndexStatus: indexStatus,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}, nil
}
