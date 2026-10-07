package adapters

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/conversation"
)

var _ conversation.Repository = (*ConversationPostgresRepository)(nil)

type ConversationPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewConversationPostgresRepository(pool *pgxpool.Pool) *ConversationPostgresRepository {
	return &ConversationPostgresRepository{pool: pool}
}

func (r *ConversationPostgresRepository) AddConversation(ctx context.Context, conversation *conversation.Conversation) error {
	return nil
}

func (r *ConversationPostgresRepository) UpdateConversation(ctx context.Context, conversationID string, updateFn func(conversation *conversation.Conversation) error) error {
	return nil
}

func (r *ConversationPostgresRepository) DeleteConversation(ctx context.Context, conversationID string) error {
	return nil
}
