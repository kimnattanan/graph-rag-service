package query

import (
	"time"

	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
)

type Document struct {
	ID          string
	Title       string
	Content     string
	Tags        []string
	IndexStatus document.IndexStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DocumentSummary struct {
	ID          string
	Title       string
	Tags        []string
	IndexStatus document.IndexStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DocumentSummaryList struct {
	Items []DocumentSummary
	Total int
}

type DocumentIndexStatus struct {
	DocumentID   string
	Status       document.IndexStatus
	StartedAt    *time.Time
	FinishedAt   *time.Time
	ErrorMessage *string
}

type RetrieveResult struct {
	Chunks []RetrievedChunk
}

type RetrievedChunk struct {
	DocumentID    string
	DocumentTitle string
	ChunkID       *string
	Text          string
	Score         float64
	GraphPath     *[]string
}
