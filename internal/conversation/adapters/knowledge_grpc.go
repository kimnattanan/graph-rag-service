package adapters

import (
	"context"

	knowledgepb "github.com/kimnattanan/graph-rag-service/internal/common/genproto/knowledge"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/retrieval"
)

var _ retrieval.Retriever = (*KnowledgeGrpc)(nil)

type KnowledgeGrpc struct {
	client knowledgepb.KnowledgeServiceClient
}

func NewKnowledgeGrpc(client knowledgepb.KnowledgeServiceClient) *KnowledgeGrpc {
	return &KnowledgeGrpc{client: client}
}

func (r *KnowledgeGrpc) Retrieve(ctx context.Context, query string, topK int, tags []string) ([]retrieval.Chunk, error) {
	resp, err := r.client.Retrieve(ctx, &knowledgepb.RetrieveRequest{
		Query: query,
		TopK:  int32(topK),
		Tags:  tags,
	})
	if err != nil {
		return nil, err
	}
	chunks := make([]retrieval.Chunk, len(resp.Chunks))
	for i, chunk := range resp.Chunks {
		chunks[i] = retrieval.Chunk{
			DocumentID:    chunk.DocumentId,
			DocumentTitle: chunk.DocumentTitle,
			ChunkID:       chunk.ChunkId,
			Text:          chunk.Text,
			Score:         float64(chunk.Score),
		}
	}
	return chunks, nil
}
