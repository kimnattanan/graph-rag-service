package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kimnattanan/graph-rag-service/internal/knowledge/config"
)

type Embedder struct {
	cfg    config.Embedder
	client *http.Client
}

func NewEmbedder(cfg config.Embedder) *Embedder {
	return &Embedder{
		cfg:    cfg,
		client: &http.Client{},
	}
}

type embeddingRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
}

type embeddingData struct {
	Index     int       `json:"index"`
	Embedding []float64 `json:"embedding"`
}

type embeddingResponse struct {
	Data  []embeddingData `json:"data"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *Embedder) Embed(ctx context.Context, texts []string) ([][]float64, error) {
	if e.cfg.URL == "" {
		return nil, fmt.Errorf("embedder: url is empty")
	}
	if e.cfg.Model == "" {
		return nil, fmt.Errorf("embedder: model is empty")
	}
	if len(texts) == 0 {
		return nil, nil
	}

	embeddings, err := e.embed(ctx, texts)
	if err != nil {
		return nil, err
	}
	return embeddings, nil
}

func (e *Embedder) embed(ctx context.Context, inputs []string) ([][]float64, error) {
	if len(inputs) == 0 {
		return nil, nil
	}

	body, err := json.Marshal(embeddingRequest{
		Model: e.cfg.Model,
		Input: inputs,
	})
	if err != nil {
		return nil, fmt.Errorf("embedder: encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(e.cfg.URL, "/"), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("embedder: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if e.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.cfg.APIKey)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("embedder: request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("embedder: read response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("embedder: %s: %s", resp.Status, embeddingErrorMessage(respBody))
	}

	var decoded embeddingResponse
	if err := json.Unmarshal(respBody, &decoded); err != nil {
		return nil, fmt.Errorf("embedder: decode response: %w", err)
	}
	embeddings, err := orderEmbeddings(inputs, decoded.Data)
	if err != nil {
		return nil, err
	}
	return embeddings, nil
}

func orderEmbeddings(inputs []string, data []embeddingData) ([][]float64, error) {
	if len(data) != len(inputs) {
		return nil, fmt.Errorf("embedder: expected %d embeddings, got %d", len(inputs), len(data))
	}
	out := make([][]float64, len(inputs))
	seen := make([]bool, len(inputs))
	for _, item := range data {
		if item.Index < 0 || item.Index >= len(inputs) {
			return nil, fmt.Errorf("embedder: embedding index %d out of range", item.Index)
		}
		if seen[item.Index] {
			return nil, fmt.Errorf("embedder: duplicate embedding index %d", item.Index)
		}
		if len(item.Embedding) == 0 {
			return nil, fmt.Errorf("embedder: empty embedding at index %d", item.Index)
		}
		seen[item.Index] = true
		out[item.Index] = item.Embedding
	}
	return out, nil
}

func embeddingErrorMessage(body []byte) string {
	var decoded embeddingResponse
	if err := json.Unmarshal(body, &decoded); err == nil && decoded.Error != nil && decoded.Error.Message != "" {
		return decoded.Error.Message
	}
	message := strings.TrimSpace(string(body))
	if message == "" {
		return "request failed"
	}
	return message
}
