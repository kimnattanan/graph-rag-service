package adapters

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/conversation/config"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/retrieval"
)

var _ retrieval.Completer = (*LLMCompleter)(nil)

type LLMCompleter struct {
	cfg config.LLM
}

func NewLLMCompleter(cfg config.LLM) *LLMCompleter {
	return &LLMCompleter{cfg: cfg}
}

func (c *LLMCompleter) Complete(ctx context.Context, messages []retrieval.ChatMessage) (string, error) {
	return "", nil
}
