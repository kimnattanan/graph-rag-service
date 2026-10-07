package adapters

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/indexing"
)

var _ indexing.Extractor = (*Extractor)(nil)

type Extractor struct{}

func NewExtractor() *Extractor {
	return &Extractor{}
}

func (e *Extractor) Extract(ctx context.Context, content string) ([]indexing.ExtractResult, error) {
	return []indexing.ExtractResult{
		indexing.NewExtractResult("Content A", []string{"EntityX", "EntityY"}),
		indexing.NewExtractResult("Content B", []string{"EntityZ", "EntityY"}),
	}, nil
}
