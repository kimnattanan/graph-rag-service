package indexing

import "context"

type Repository interface {
	ReplaceChunks(ctx context.Context, documentID string, chunks []ExtractResult) error
	DeleteOrphans(ctx context.Context) error
}

type Extractor interface {
	Extract(ctx context.Context, content string) ([]ExtractResult, error)
}
