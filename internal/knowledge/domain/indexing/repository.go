package indexing

import "context"

type Repository interface {
	ReplaceChunks(ctx context.Context, documentID string, chunks []ExtractResult, chunkEmbeddings [][]float64) error
	ClaimNextPendingEntity(ctx context.Context) (*Entity, error)
	EmbedEntity(ctx context.Context, entityName string, embedding []float64) error
	DeleteOrphans(ctx context.Context) error
}

type Extractor interface {
	Extract(ctx context.Context, content string) ([]ExtractResult, error)
}

type Embedder interface {
	Embed(ctx context.Context, texts []string) ([][]float64, error)
}
