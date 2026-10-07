package conversation

import "context"

type Repository interface {
	AddConversation(ctx context.Context, conversation *Conversation) error
	UpdateConversation(ctx context.Context, conversationID string, updateFn func(conversation *Conversation) error) error
	DeleteConversation(ctx context.Context, conversationID string) error
}
