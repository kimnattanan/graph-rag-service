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

// ChatMessage is one turn passed to the language model.
type ChatMessage struct {
	Role    string
	Content string
}

// Completer generates an assistant reply from the prompt messages.
type Completer interface {
	Complete(ctx context.Context, messages []ChatMessage) (string, error)
}
