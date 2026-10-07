package ports

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/common/server/httperr"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type HttpServer struct {
	app app.Application
}

func NewHttpServer(app app.Application) HttpServer {
	return HttpServer{app}
}

func (h HttpServer) GetDocument(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {
	appDoc, err := h.app.Queries.GetDocument.Handle(r.Context(), query.GetDocument{
		DocumentID: documentId.String(),
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	doc, err := appDocumentToResponse(appDoc)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, doc)
}


func (h HttpServer) GetDocumentIndexStatus(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {
	status, err := h.app.Queries.GetDocumentIndexStatus.Handle(r.Context(), query.GetDocumentIndexStatus{
		DocumentID: documentId.String(),
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	indexJob, err := appDocumentIndexStatusToResponse(status)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, indexJob)
}

func (h HttpServer) ListDocuments(w http.ResponseWriter, r *http.Request, params ListDocumentsParams) {
	var indexStatus *document.IndexStatus
	if params.IndexStatus != nil {
		parsedIndexStatus, err := parseIndexStatus(*params.IndexStatus)
		if err != nil {
			httperr.RespondWithSlugError(err, w, r)
			return
		}
		indexStatus = &parsedIndexStatus
	}
	documentSummaryList, err := h.app.Queries.ListDocuments.Handle(r.Context(), query.ListDocuments{
		Limit:       params.Limit,
		Offset:      params.Offset,
		Tag:         params.Tag,
		IndexStatus: indexStatus,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	documentSummaries := make([]DocumentSummary, len(documentSummaryList.Items))
	for i, documentSummary := range documentSummaryList.Items {
		if documentSummaries[i], err = appDocumentSummaryToResponse(documentSummary); err != nil {
			httperr.RespondWithSlugError(err, w, r)
			return
		}
	}

	render.Respond(w, r, DocumentList{
		Items: documentSummaries,
		Total: documentSummaryList.Total,
	})
}

func (h HttpServer) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var body CreateDocumentRequest
	if err := render.Decode(r, &body); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	var tags []string
	if body.Tags != nil {
		tags = *body.Tags
	}

	err := h.app.Commands.CreateDocument.Handle(r.Context(), command.CreateDocument{
		DocumentID: body.Id.String(),
		Title:      body.Title,
		Content:    body.Content,
		Tags:       tags,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func (h HttpServer) UpdateDocument(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {
	var body UpdateDocumentRequest
	if err := render.Decode(r, &body); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	err := h.app.Commands.UpdateDocument.Handle(r.Context(), command.UpdateDocument{
		DocumentID: documentId.String(),
		Title:      body.Title,
		Content:    body.Content,
		Tags:       body.Tags,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) DeleteDocument(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {
	err := h.app.Commands.DeleteDocument.Handle(r.Context(), command.DeleteDocument{
		DocumentID: documentId.String(),
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) ReindexDocument(w http.ResponseWriter, r *http.Request, documentId openapi_types.UUID) {
	err := h.app.Commands.ReindexDocument.Handle(r.Context(), command.ReindexDocument{
		DocumentID: documentId.String(),
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusAccepted)
}

func (h HttpServer) Retrieve(w http.ResponseWriter, r *http.Request) {
	var body RetrieveRequest
	if err := render.Decode(r, &body); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	result, err := h.app.Queries.Retrieve.Handle(r.Context(), query.Retrieve{
		QueryStr: body.Query,
		TopK:     body.TopK,
		Tags:     body.Tags,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	chunks, err := appRetrievedChunksToResponse(result.Chunks)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, RetrieveResult{
		Query:  body.Query,
		Chunks: chunks,
	})
}

// Mapper functions

func parseDomainIndexStatus(status document.IndexStatus) (IndexStatus, error) {
	switch status {
	case document.IndexStatusPending:
		return Pending, nil
	case document.IndexStatusIndexing:
		return Indexing, nil
	case document.IndexStatusCompleted:
		return Completed, nil
	case document.IndexStatusFailed:
		return Failed, nil
	default:
		return "", commonerrors.NewIncorrectInputError("Invalid index status", "index-status-invalid")
	}
}

func parseIndexStatus(status IndexStatus) (document.IndexStatus, error) {
	switch status {
	case Pending:
		return document.IndexStatusPending, nil
	case Indexing:
		return document.IndexStatusIndexing, nil
	case Completed:
		return document.IndexStatusCompleted, nil
	case Failed:
		return document.IndexStatusFailed, nil
	default:
		return 0, commonerrors.NewIncorrectInputError("Invalid index status", "index-status-invalid")
	}
}

func appDocumentToResponse(d query.Document) (Document, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil {
		return Document{}, commonerrors.NewIncorrectInputError("Invalid document ID", "document-id-invalid")
	}
	indexStatus, err := parseDomainIndexStatus(d.IndexStatus)
	if err != nil {
		return Document{}, err
	}
	return Document{
		Id:          id,
		Title:       d.Title,
		Content:     d.Content,
		Tags:        d.Tags,
		IndexStatus: indexStatus,
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}, nil
}

func appDocumentIndexStatusToResponse(s query.DocumentIndexStatus) (IndexJob, error) {
	id, err := uuid.Parse(s.DocumentID)
	if err != nil {
		return IndexJob{}, commonerrors.NewIncorrectInputError("Invalid document ID", "document-id-invalid")
	}
	status, err := parseDomainIndexStatus(s.Status)
	if err != nil {
		return IndexJob{}, err
	}
	return IndexJob{
		DocumentId:   id,
		Status:       status,
		StartedAt:    s.StartedAt,
		FinishedAt:   s.FinishedAt,
		ErrorMessage: s.ErrorMessage,
	}, nil
}

func appRetrievedChunksToResponse(chunks []query.RetrievedChunk) ([]RetrievedChunk, error) {
	response := make([]RetrievedChunk, len(chunks))
	for i, chunk := range chunks {
		docID, err := uuid.Parse(chunk.DocumentID)
		if err != nil {
			return nil, commonerrors.NewIncorrectInputError("Invalid document ID", "document-id-invalid")
		}
		response[i] = RetrievedChunk{
			DocumentId:    docID,
			DocumentTitle: chunk.DocumentTitle,
			ChunkId:       chunk.ChunkID,
			Text:          chunk.Text,
			Score:         float32(chunk.Score),
			GraphPath:     chunk.GraphPath,
		}
	}
	return response, nil
}

func appDocumentSummaryToResponse(s query.DocumentSummary) (DocumentSummary, error) {
	id, err := uuid.Parse(s.ID)
	if err != nil {
		return DocumentSummary{}, commonerrors.NewIncorrectInputError("Invalid document ID", "document-id-invalid")
	}
	indexStatus, err := parseDomainIndexStatus(s.IndexStatus)
	if err != nil {
		return DocumentSummary{}, err
	}
	return DocumentSummary{
		Id:          id,
		Title:       s.Title,
		Tags:        s.Tags,
		IndexStatus: indexStatus,
		CreatedAt:   s.CreatedAt,
		UpdatedAt:   s.UpdatedAt,
	}, nil
}
