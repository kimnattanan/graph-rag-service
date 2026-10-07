package query

import (
	"time"

	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/conversation"
)

type Conversation struct {
	ID        string
	UserID    string
	Title     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Messages  []Message
}

type ConversationSummary struct {
	ID           string
	Title        string
	CreatedAt    time.Time
	UpdatedAt    time.Time
	MessageCount int
}

type ConversationSummaryList struct {
	Items []ConversationSummary
	Total int
}

type Message struct {
	ID             string
	ConversationID string
	Role           conversation.Role
	Content        string
	Sources        []Source
	CreatedAt      time.Time
}

type MessageList struct {
	Items []Message
	Total int
}

type Source struct {
	DocumentID    string
	DocumentTitle string
	ChunkID       string
	Text          string
	Score         float64
}
