package conversation

import "context"

type Repository interface {
	AddConversation(ctx context.Context, conv *Conversation) error
	UpdateConversation(ctx context.Context, convID string, updateFn func(conv *Conversation) error) error
	DeleteConversation(ctx context.Context, convID string) error
	GetConversation(ctx context.Context, convID string) (*Conversation, error)
}
