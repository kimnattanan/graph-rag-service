package query

import "time"

type IndexStatus int

const (
	IndexStatusPending IndexStatus = iota
	IndexStatusCompleted
	IndexStatusFailed
)

type Document struct {
	ID          string
	Title       string
	Content     string
	Tags        []string
	IndexStatus IndexStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type DocumentSummary struct {
	ID          string
	Title       string
	Tags        []string
	IndexStatus IndexStatus
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
