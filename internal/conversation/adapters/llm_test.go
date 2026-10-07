package adapters

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kimnattanan/graph-rag-service/internal/conversation/config"
)

func TestLLMCompleterComplete(t *testing.T) {
	var got chatCompletionRequest
	var gotAuth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
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
		if err := json.Unmarshal(body, &got); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"  Paris  "}}]}`))
	}))
	defer server.Close()

	completer := NewLLMCompleter(config.LLM{
		BaseURL: server.URL + "/v1",
		APIKey:  "ollama",
		Model:   "llama3.1",
	})
	reply, err := completer.Complete(context.Background(), "Where is the capital?")
	if err != nil {
		t.Fatal(err)
	}
	if reply != "Paris" {
		t.Fatalf("reply = %q", reply)
	}
	if gotAuth != "Bearer ollama" {
		t.Fatalf("authorization = %q", gotAuth)
	}
	if got.Model != "llama3.1" {
		t.Fatalf("model = %q", got.Model)
	}
	if len(got.Messages) != 1 || got.Messages[0].Role != "user" || got.Messages[0].Content != "Where is the capital?" {
		t.Fatalf("messages = %#v", got.Messages)
	}
}

func TestLLMCompleterCompleteAPIError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, `{"error":{"message":"model not found"}}`, http.StatusNotFound)
	}))
	defer server.Close()

	completer := NewLLMCompleter(config.LLM{
		BaseURL: server.URL + "/v1/",
		Model:   "missing",
	})
	_, err := completer.Complete(context.Background(), "hello")
	if err == nil || !strings.Contains(err.Error(), "model not found") {
		t.Fatalf("err = %v", err)
	}
}
