package adapters

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/kimnattanan/graph-rag-service/internal/knowledge/config"
)

func TestExtractEntitiesInOneCall(t *testing.T) {
	var calls int
	var gotAuth string
	var gotReq entityChatRequest
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
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
		if err := json.Unmarshal(body, &gotReq); err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{
			"choices":[{"message":{"role":"assistant","content":"{\"entities\":[[\" Memgraph \",\"vector index\",\"memgraph\"],[\"Memgraph\"]]}"}}]
		}`)
	}))
	defer server.Close()

	extractor := NewExtractor(config.LLM{
		BaseURL: server.URL + "/v1",
		APIKey:  "test-key",
		Model:   "llama",
	})
	got, err := extractor.Extract(context.Background(), "# Alpha\r\n\r\nMemgraph stores a vector index.\r\n\r\n# Beta\r\n\r\nMemgraph is the database.")
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("calls = %d", calls)
	}
	if gotAuth != "Bearer test-key" {
		t.Fatalf("auth = %s", gotAuth)
	}
	if gotReq.Model != "llama" || len(gotReq.Messages) != 1 || !strings.Contains(gotReq.Messages[0].Content, "Chunk 1:") || !strings.Contains(gotReq.Messages[0].Content, "Chunk 2:") {
		t.Fatalf("request = %#v", gotReq)
	}
	if len(got) != 2 {
		t.Fatalf("len = %d", len(got))
	}
	if got[0].Content() != "# Alpha\n\nMemgraph stores a vector index." {
		t.Fatalf("first content = %q", got[0].Content())
	}
	if strings.Join(got[0].Entities(), "|") != "Memgraph|vector index" {
		t.Fatalf("first entities = %#v", got[0].Entities())
	}
	if strings.Join(got[1].Entities(), "|") != "Memgraph" {
		t.Fatalf("second entities = %#v", got[1].Entities())
	}
}

func TestParseEntityResponse(t *testing.T) {
	got, err := parseEntityResponse("```json\n{\"entities\":[[\" Memgraph  \",\"memgraph\"],[]]}\n```", 2)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got[0], "|") != "Memgraph" || len(got[1]) != 0 {
		t.Fatalf("got %#v", got)
	}

	got, err = parseEntityResponse(`[["Alpha"],["Beta"]]`, 2)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Join(got[0], "|") != "Alpha" || strings.Join(got[1], "|") != "Beta" {
		t.Fatalf("got %#v", got)
	}

	if _, err := parseEntityResponse(`{"entities":[[]]}`, 2); err == nil {
		t.Fatal("expected count mismatch")
	}
}

func TestExtractEmptyContent(t *testing.T) {
	_, err := NewExtractor(config.LLM{}).Extract(context.Background(), " \n\t ")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExtractRequiresLLMConfig(t *testing.T) {
	_, err := NewExtractor(config.LLM{}).Extract(context.Background(), "hello")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestExtractCancelledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := NewExtractor(config.LLM{}).Extract(ctx, "hello")
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestChunkMarkdownPacksParagraphs(t *testing.T) {
	content := "aaaa\n\nbbbb\n\ncccc"
	got := chunkMarkdown(content, 14, 0)
	want := []string{"aaaa\n\nbbbb", "cccc"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Fatalf("got %#v", got)
	}
}

func TestChunkMarkdownRepeatsHeading(t *testing.T) {
	body := strings.Repeat("word ", 30)
	got := chunkMarkdown("# Title\n\n"+body, 80, 10)
	if len(got) < 2 {
		t.Fatalf("got %d chunks: %#v", len(got), got)
	}
	for _, chunk := range got {
		if !strings.HasPrefix(chunk, "# Title\n\n") {
			t.Fatalf("missing heading: %q", chunk)
		}
		if n := utf8.RuneCountInString(chunk); n > 80 {
			t.Fatalf("len %d: %q", n, chunk)
		}
	}
}

func TestChunkMarkdownOverlapsLongText(t *testing.T) {
	got := chunkMarkdown(strings.Repeat("a", 250), 100, 20)
	if len(got) < 3 {
		t.Fatalf("got %d chunks", len(got))
	}
	if !strings.HasSuffix(got[0], got[1][:20]) {
		t.Fatalf("missing overlap\n%q\n%q", got[0], got[1])
	}
	for _, chunk := range got {
		if n := utf8.RuneCountInString(chunk); n > 100 {
			t.Fatalf("len %d", n)
		}
	}
}

func TestChunkMarkdownKeepsFenceIntact(t *testing.T) {
	content := "# Title\n\n```\n# not a heading\n```\n\nAfter."
	got := chunkMarkdown(content, 500, 20)
	if len(got) != 1 {
		t.Fatalf("got %#v", got)
	}
	if !strings.Contains(got[0], "# not a heading") || !strings.HasPrefix(got[0], "# Title\n\n") {
		t.Fatalf("got %q", got[0])
	}
}
