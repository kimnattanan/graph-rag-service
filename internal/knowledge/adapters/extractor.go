package adapters

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"net/http"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/kimnattanan/graph-rag-service/internal/knowledge/config"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/indexing"
)

const (
	// all-minilm truncates near 256 tokens. 500 runes stays inside that for
	// typical prose, including denser scripts than English.
	maxChunkRunes     = 500
	chunkOverlapRunes = 80
)

var _ indexing.Extractor = (*Extractor)(nil)

type Extractor struct {
	cfg    config.LLM
	client *http.Client
}

func NewExtractor(cfg config.LLM) *Extractor {
	return &Extractor{
		cfg:    cfg,
		client: &http.Client{},
	}
}

/*
markdown -> sections -> paragraph/sentence/window -> chunks
*/
func (e *Extractor) Extract(ctx context.Context, content string) ([]indexing.ExtractResult, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	chunks := chunkMarkdown(content, maxChunkRunes, chunkOverlapRunes)
	if len(chunks) == 0 {
		return nil, fmt.Errorf("extractor: content is empty")
	}

	entities, err := e.extractEntities(ctx, chunks)
	if err != nil {
		return nil, err
	}

	results := make([]indexing.ExtractResult, len(chunks))
	for i, chunk := range chunks {
		results[i] = indexing.NewExtractResult(chunk, entities[i])
	}
	return results, nil
}

type entityChatRequest struct {
	Model    string              `json:"model"`
	Messages []entityChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
}

type entityChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type entityChatResponse struct {
	Choices []struct {
		Message entityChatMessage `json:"message"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (e *Extractor) extractEntities(ctx context.Context, chunks []string) ([][]string, error) {
	if e.cfg.BaseURL == "" {
		return nil, fmt.Errorf("extractor: llm base url is empty")
	}
	if e.cfg.Model == "" {
		return nil, fmt.Errorf("extractor: llm model is empty")
	}

	bodyStruct := entityChatRequest{
		Model: e.cfg.Model,
		Messages: []entityChatMessage{{
			Role:    "user",
			Content: entityPrompt(chunks),
		}},
	}
	fmt.Printf("REQUEST: %#v\n", bodyStruct)
	body, err := json.Marshal(bodyStruct)
	if err != nil {
		return nil, fmt.Errorf("extractor: encode request: %w", err)
	}

	endpoint := strings.TrimRight(e.cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("extractor: create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	if e.cfg.APIKey != "" {
		req.Header.Set("Authorization", "Bearer "+e.cfg.APIKey)
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("extractor: request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("extractor: read response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("extractor: %s: %s", resp.Status, entityErrorMessage(respBody))
	}

	var decoded entityChatResponse
	if err := json.Unmarshal(respBody, &decoded); err != nil {
		return nil, fmt.Errorf("extractor: decode response: %w", err)
	}

	// mock
	// decoded := entityChatResponseMock()

	if len(decoded.Choices) == 0 {
		return nil, fmt.Errorf("extractor: empty choices")
	}
	fmt.Printf("RESPONSE: %#v\n", decoded)
	return parseEntityResponse(decoded.Choices[0].Message.Content, len(chunks))
}

func entityPrompt(chunks []string) string {
	var b strings.Builder
	b.WriteString("Extract the named entities from each chunk. Named entities are people, organizations, products, systems, and specific concepts written in the chunk.\n")
	b.WriteString("Return JSON only, with this shape: {\"entities\":[[\"Name\"],[]]}\n")
	b.WriteString("The entities array must contain exactly ")
	b.WriteString(strconv.Itoa(len(chunks)))
	b.WriteString(" lists, in chunk order. Use an empty list when a chunk has no entities. Do not add entities that are not written in the chunk. Keep each name as written.\n\n")
	for i, chunk := range chunks {
		fmt.Fprintf(&b, "Chunk %d:\n%s\n\n", i+1, chunk)
	}
	return b.String()
}

func parseEntityResponse(raw string, chunkCount int) ([][]string, error) {
	payload := jsonPayload(raw)
	var lists [][]string
	var wrapped struct {
		Entities [][]string `json:"entities"`
	}
	if err := json.Unmarshal([]byte(payload), &wrapped); err == nil && wrapped.Entities != nil {
		lists = wrapped.Entities
	} else if err := json.Unmarshal([]byte(payload), &lists); err != nil {
		return nil, fmt.Errorf("extractor: decode entities: %w", err)
	}
	fmt.Printf("PARSED: %#v\n", lists)
	if len(lists) != chunkCount {
		return nil, fmt.Errorf("extractor: expected %d entity lists, got %d", chunkCount, len(lists))
	}
	return normalizeEntityLists(lists), nil
}

func jsonPayload(raw string) string {
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "```") {
		raw = strings.TrimPrefix(raw, "```")
		if nl := strings.IndexByte(raw, '\n'); nl >= 0 {
			raw = raw[nl+1:]
		}
		raw = strings.TrimSuffix(strings.TrimSpace(raw), "```")
		raw = strings.TrimSpace(raw)
	}
	startObj := strings.Index(raw, "{")
	startArr := strings.Index(raw, "[")
	if startObj == -1 && startArr == -1 {
		return raw
	}
	start := startObj
	closer := byte('}')
	if startArr != -1 && (startObj == -1 || startArr < startObj) {
		start = startArr
		closer = ']'
	}
	end := strings.LastIndexByte(raw, closer)
	if end < start {
		return raw[start:]
	}
	return raw[start : end+1]
}

func normalizeEntityLists(lists [][]string) [][]string {
	canonical := make(map[string]string, len(lists))
	out := make([][]string, len(lists))
	for i, list := range lists {
		seen := make(map[string]struct{}, len(list))
		for _, name := range list {
			name = strings.Join(strings.Fields(name), " ")
			if name == "" {
				continue
			}
			key := strings.ToLower(name)
			if chosen, ok := canonical[key]; ok {
				name = chosen
			} else {
				canonical[key] = name
			}
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			out[i] = append(out[i], name)
		}
		if out[i] == nil {
			out[i] = []string{}
		}
	}
	return out
}

func entityErrorMessage(body []byte) string {
	var decoded entityChatResponse
	if err := json.Unmarshal(body, &decoded); err == nil && decoded.Error != nil && decoded.Error.Message != "" {
		return decoded.Error.Message
	}
	message := strings.TrimSpace(string(body))
	if message == "" {
		return "request failed"
	}
	return message
}

type markdownSection struct {
	heading string
	body    string
}

func chunkMarkdown(content string, maxRunes, overlap int) []string {
	content = strings.ReplaceAll(content, "\r\n", "\n")
	content = strings.ReplaceAll(content, "\r", "\n")
	content = strings.TrimSpace(content)
	if content == "" || maxRunes < 1 {
		return nil
	}
	if overlap < 0 {
		overlap = 0
	}

	var chunks []string
	for _, section := range splitHeadingSections(content) {
		chunks = append(chunks, chunkSection(section, maxRunes, overlap)...)
	}
	return chunks
}

func chunkSection(section markdownSection, maxRunes, overlap int) []string {
	heading := section.heading
	body := section.body
	if heading == "" {
		return chunkText(body, maxRunes, overlap)
	}
	if body == "" {
		return chunkText(heading, maxRunes, overlap)
	}

	joined := heading + "\n\n" + body
	if utf8.RuneCountInString(joined) <= maxRunes {
		return []string{joined}
	}

	// budget for body content
	budget := maxRunes - utf8.RuneCountInString(heading) - 2
	if budget < maxRunes/4 {
		return chunkText(joined, maxRunes, overlap)
	}

	pieces := chunkText(body, budget, overlap)
	chunks := make([]string, 0, len(pieces))
	for _, piece := range pieces {
		chunks = append(chunks, heading+"\n\n"+piece)
	}
	return chunks
}

func chunkText(text string, maxRunes, overlap int) []string {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil
	}
	if utf8.RuneCountInString(text) <= maxRunes {
		return []string{text}
	}

	paragraphs := splitParagraphs(text)
	if len(paragraphs) > 1 {
		return packParts(paragraphs, "\n\n", maxRunes, overlap)
	}

	sentences := splitSentences(text)
	if len(sentences) > 1 {
		return packParts(sentences, " ", maxRunes, overlap)
	}

	return window(text, maxRunes, overlap)
}

func splitHeadingSections(content string) []markdownSection {
	var sections []markdownSection
	var heading string
	var body []string
	fence := ""

	flush := func() {
		b := strings.TrimSpace(strings.Join(body, "\n"))
		h := strings.TrimSpace(heading)
		body = nil
		if h == "" && b == "" {
			return
		}
		sections = append(sections, markdownSection{heading: h, body: b})
	}

	for _, line := range strings.Split(content, "\n") {
		if fence != "" {
			body = append(body, line)
			if isClosingFence(line, fence) {
				fence = ""
			}
			continue
		}
		if marker, ok := openingFence(line); ok {
			body = append(body, line)
			fence = marker
			continue
		}
		if isATXHeading(line) {
			flush()
			heading = strings.TrimSpace(line)
			continue
		}
		body = append(body, line)
	}
	flush()
	return sections
}

func isATXHeading(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" || trimmed[0] != '#' {
		return false
	}
	marks := 0
	for _, r := range trimmed {
		if r != '#' {
			break
		}
		marks++
	}
	if marks == 0 || marks > 6 || marks == utf8.RuneCountInString(trimmed) {
		return false
	}
	next, _ := utf8.DecodeRuneInString(trimmed[marks:])
	return next == ' ' || next == '\t'
}

func openingFence(line string) (string, bool) {
	trimmed := strings.TrimSpace(line)
	switch {
	case strings.HasPrefix(trimmed, "```"):
		return "`", true
	case strings.HasPrefix(trimmed, "~~~"):
		return "~", true
	default:
		return "", false
	}
}

func isClosingFence(line, marker string) bool {
	trimmed := strings.TrimSpace(line)
	if len(trimmed) < 3 {
		return false
	}
	for _, r := range trimmed {
		if string(r) != marker {
			return false
		}
	}
	return true
}

func splitParagraphs(text string) []string {
	raw := strings.Split(text, "\n\n")
	parts := make([]string, 0, len(raw))
	for _, part := range raw {
		part = strings.TrimSpace(part)
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func splitSentences(text string) []string {
	runes := []rune(text)
	var parts []string
	start := 0
	for i := 0; i < len(runes); i++ {
		if !isSentenceEnd(runes[i]) {
			continue
		}
		if runes[i] == '.' && i > 0 && i+1 < len(runes) && unicode.IsDigit(runes[i-1]) && unicode.IsDigit(runes[i+1]) {
			continue
		}
		end := i + 1
		for end < len(runes) && (runes[end] == ' ' || runes[end] == '\n' || runes[end] == '\t') {
			end++
		}
		if end == i+1 && end < len(runes) {
			continue
		}
		part := strings.TrimSpace(string(runes[start : i+1]))
		if part != "" {
			parts = append(parts, part)
		}
		start = end
		i = end - 1
	}
	if start < len(runes) {
		part := strings.TrimSpace(string(runes[start:]))
		if part != "" {
			parts = append(parts, part)
		}
	}
	return parts
}

func isSentenceEnd(r rune) bool {
	return r == '.' || r == '!' || r == '?' || r == '。'
}

func packParts(parts []string, sep string, maxRunes, overlap int) []string {
	var chunks []string
	var current []string
	currentLen := 0
	sepLen := utf8.RuneCountInString(sep)

	emit := func(keepOverlap bool) {
		if len(current) == 0 {
			return
		}
		chunks = append(chunks, strings.Join(current, sep))
		if !keepOverlap || overlap <= 0 {
			current = nil
			currentLen = 0
			return
		}

		kept := make([]string, 0, len(current))
		keptLen := 0
		for i := len(current) - 1; i >= 0; i-- {
			partLen := utf8.RuneCountInString(current[i])
			extra := partLen
			if len(kept) > 0 {
				extra += sepLen
			}
			if keptLen+extra > overlap && len(kept) > 0 {
				break
			}
			if partLen > overlap {
				break
			}
			kept = append([]string{current[i]}, kept...)
			keptLen += extra
		}
		current = kept
		currentLen = keptLen
	}

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		partLen := utf8.RuneCountInString(part)
		if partLen > maxRunes {
			emit(false)
			chunks = append(chunks, window(part, maxRunes, overlap)...)
			continue
		}

		add := partLen
		if len(current) > 0 {
			add += sepLen
		}
		if len(current) > 0 && currentLen+add > maxRunes {
			emit(true)
			add = partLen
			if len(current) > 0 {
				add += sepLen
			}
			if len(current) > 0 && currentLen+add > maxRunes {
				current = nil
				currentLen = 0
				add = partLen
			}
		}
		current = append(current, part)
		currentLen += add
	}
	emit(false)
	return chunks
}

func window(text string, maxRunes, overlap int) []string {
	runes := []rune(strings.TrimSpace(text))
	if len(runes) == 0 {
		return nil
	}
	if len(runes) <= maxRunes {
		return []string{string(runes)}
	}
	if overlap < 0 {
		overlap = 0
	}
	if overlap >= maxRunes {
		overlap = maxRunes / 5
	}

	var chunks []string
	start := 0
	for start < len(runes) {
		end := start + maxRunes
		if end > len(runes) {
			end = len(runes)
		} else if br := lastBreak(runes[start:end]); br >= 0 {
			end = start + br
		}

		piece := strings.TrimSpace(string(runes[start:end]))
		if piece != "" {
			chunks = append(chunks, piece)
		}
		if end >= len(runes) {
			break
		}

		next := end - overlap
		if next <= start {
			next = end
			if next <= start {
				next = start + 1
			}
		}
		start = next
	}
	return chunks
}

func lastBreak(runes []rune) int {
	min := len(runes) / 2
	for i := len(runes) - 1; i >= min; i-- {
		if runes[i] == ' ' || runes[i] == '\n' || runes[i] == '\t' {
			return i
		}
	}
	return -1
}
