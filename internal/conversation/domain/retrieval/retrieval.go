package retrieval

import "context"

// Chunk is a knowledge passage returned by retrieval.
type Chunk struct {
	DocumentID    string
	DocumentTitle string
	ChunkID       string
	Text          string
	Score         float64
}

// Retriever loads knowledge chunks for a user question.
type Retriever interface {
	Retrieve(ctx context.Context, query string, topK int, tags []string) ([]Chunk, error)
}

// Completer generates an assistant reply from the prompt messages.
type Completer interface {
	Complete(ctx context.Context, prompt string) (string, error)
}
