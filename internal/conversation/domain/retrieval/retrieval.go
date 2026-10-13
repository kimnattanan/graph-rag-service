package retrieval

import "context"

type Chunk struct {
	DocumentID    string
	DocumentTitle string
	ChunkID       string
	Text          string
	Score         float64
}

type Retriever interface {
	Retrieve(ctx context.Context, query string, topK int, tags []string) ([]Chunk, error)
}

type Completer interface {
	Complete(ctx context.Context, prompt string) (string, error)
}
