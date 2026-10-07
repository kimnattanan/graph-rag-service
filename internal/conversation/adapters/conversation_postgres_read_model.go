package adapters

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/conversation"
)

var (
	_ query.ListConversationsReadModel = (*ConversationPostgresReadModel)(nil)
	_ query.GetConversationReadModel   = (*ConversationPostgresReadModel)(nil)
	_ query.ListMessagesReadModel      = (*ConversationPostgresReadModel)(nil)
	_ query.GetMessageReadModel        = (*ConversationPostgresReadModel)(nil)
)

type ConversationPostgresReadModel struct {
	pool *pgxpool.Pool
}

func NewConversationPostgresReadModel(pool *pgxpool.Pool) *ConversationPostgresReadModel {
	return &ConversationPostgresReadModel{pool: pool}
}

func (r *ConversationPostgresReadModel) ListConversations(ctx context.Context, userID string, limit int, offset int) (*query.ConversationSummaryList, error) {
	return nil, nil
}

func (r *ConversationPostgresReadModel) GetConversation(ctx context.Context, userID string, conversationID string) (*conversation.Conversation, error) {
	return nil, nil
}

func (r *ConversationPostgresReadModel) ListMessages(ctx context.Context, userID string, conversationID string, limit int, offset int) (*query.MessageList, error) {
	return nil, nil
}

func (r *ConversationPostgresReadModel) GetMessage(ctx context.Context, userID string, conversationID string, messageID string) (*query.Message, error) {
	return nil, nil
}
