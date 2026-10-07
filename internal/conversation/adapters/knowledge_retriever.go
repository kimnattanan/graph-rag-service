package adapters

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/retrieval"
)

var _ retrieval.Retriever = (*KnowledgeRetriever)(nil)

type KnowledgeRetriever struct {
	address string
}

func NewKnowledgeRetriever(address string) *KnowledgeRetriever {
	return &KnowledgeRetriever{address: address}
}

func (r *KnowledgeRetriever) Retrieve(ctx context.Context, query string, topK int, tags []string) ([]retrieval.Chunk, error) {
	return nil, nil
}
