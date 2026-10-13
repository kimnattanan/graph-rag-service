package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kimnattanan/graph-rag-service/internal/conversation/config"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/retrieval"
)

var _ retrieval.Completer = (*LLMCompleter)(nil)

type LLMCompleter struct {
	cfg    config.LLM
	client *http.Client
}

func NewLLMCompleter(cfg config.LLM) *LLMCompleter {
	return &LLMCompleter{
		cfg:    cfg,
		client: &http.Client{},
	}
}

type chatCompletionRequest struct {
	Model               string        `json:"model"`
	Messages            []chatMessage `json:"messages"`
	Stream              bool          `json:"stream"`
	MaxCompletionTokens int           `json:"max_completion_tokens"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatCompletionResponse struct {
	Choices []struct {
		Message      chatMessage `json:"message"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (c *LLMCompleter) Complete(ctx context.Context, prompt string) (string, error) {
	// fmt.Println("prompt:", prompt)

	if c.cfg.BaseURL == "" {
		return "", fmt.Errorf("llm: base url is empty")
	}
	if c.cfg.Model == "" {
		return "", fmt.Errorf("llm: model is empty")
	}

	body, err := json.Marshal(chatCompletionRequest{
		Model: c.cfg.Model,
		Messages: []chatMessage{{
			Role:    "user",
			Content: prompt,
		}},
		MaxCompletionTokens: 4096,
	})
	if err != nil {
		return "", fmt.Errorf("llm: encode request: %w", err)
	}

	endpoint := strings.TrimRight(c.cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("llm: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if c.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+c.cfg.APIKey)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("llm: request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("llm: read response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("llm: %s: %s", resp.Status, chatErrorMessage(respBody))
	}

	var decoded chatCompletionResponse
	if err := json.Unmarshal(respBody, &decoded); err != nil {
		return "", fmt.Errorf("llm: decode response: %w", err)
	}
	if len(decoded.Choices) == 0 {
		return "", fmt.Errorf("llm: empty choices")
	}
	if decoded.Choices[0].FinishReason == "length" {
		return "", fmt.Errorf("llm: response is too long")
	}
	content := strings.TrimSpace(decoded.Choices[0].Message.Content)
	if content == "" {
		return "", fmt.Errorf("llm: empty completion")
	}
	return content, nil
}

func chatErrorMessage(body []byte) string {
	var decoded chatCompletionResponse
	if err := json.Unmarshal(body, &decoded); err == nil && decoded.Error != nil && decoded.Error.Message != "" {
		return decoded.Error.Message
	}
	message := strings.TrimSpace(string(body))
	if message == "" {
		return "request failed"
	}
	return message
}
