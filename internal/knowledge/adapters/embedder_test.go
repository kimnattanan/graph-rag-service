package adapters

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kimnattanan/graph-rag-service/internal/knowledge/config"
)

func TestEmbedderEmbed(t *testing.T) {
	var requests []embeddingRequest
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/embeddings" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.Method != http.MethodPost {
			t.Errorf("method = %s", r.Method)
		}
		gotAuth = r.Header.Get("Authorization")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		var got embeddingRequest
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}
		requests = append(requests, got)

		data := make([]embeddingData, len(got.Input))
		for i := len(got.Input) - 1; i >= 0; i-- {
			data[len(got.Input)-1-i] = embeddingData{
				Index:     i,
				Embedding: []float64{float64(i), float64(len(got.Input))},
			}
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(embeddingResponse{Data: data}); err != nil {
			t.Fatal(err)
		}
	}))
	defer server.Close()

	embedder := NewEmbedder(config.Embedder{
		URL:    server.URL + "/v1/embeddings/",
		APIKey: "ollama",
		Model:  "all-minilm",
	})
	embeddings, err := embedder.Embed(context.Background(), []string{"chunk a", "chunk b", "EntityX"})
	if err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer ollama" {
		t.Fatalf("authorization = %q", gotAuth)
	}
	if len(requests) != 1 {
		t.Fatalf("requests = %d", len(requests))
	}
	got := requests[0]
	if got.Model != "all-minilm" || len(got.Input) != 3 || got.Input[0] != "chunk a" || got.Input[1] != "chunk b" || got.Input[2] != "EntityX" {
		t.Fatalf("request = %#v", got)
	}
	if len(embeddings) != 3 || embeddings[0][0] != 0 || embeddings[1][0] != 1 || embeddings[2][0] != 2 || embeddings[2][1] != 3 {
		t.Fatalf("embeddings = %#v", embeddings)
	}
}

func TestEmbedderEmbedAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":{"message":"model not found"}}`, http.StatusNotFound)
	}))
	defer server.Close()

	embedder := NewEmbedder(config.Embedder{
		URL:   server.URL + "/v1/embeddings",
		Model: "missing",
	})
	_, err := embedder.Embed(context.Background(), []string{"hello"})
	if err == nil || !strings.Contains(err.Error(), "model not found") {
		t.Fatalf("err = %v", err)
	}
}

func TestEmbedderEmbedEmptyInputs(t *testing.T) {
	embedder := NewEmbedder(config.Embedder{
		URL:   "http://embedder.invalid/v1/embeddings",
		Model: "all-minilm",
	})
	embeddings, err := embedder.Embed(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	if embeddings != nil {
		t.Fatalf("embeddings = %#v", embeddings)
	}
}
