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

func (r *ConversationPostgresRepository) AddConversation(ctx context.Context, conv *conversation.Conversation) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `
		INSERT INTO conversations (id, user_id, title, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, conv.ID(), conv.UserID(), conv.Title(), conv.CreatedAt(), conv.UpdatedAt())
	if isUniqueViolation(err, "conversations_pkey") {
		return errConversationAlreadyExists
	}
	if err != nil {
		return err
	}
	if err := insertMessages(ctx, tx, conv.Messages()); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *ConversationPostgresRepository) UpdateConversation(ctx context.Context, conversationID string, updateFn func(conv *conversation.Conversation) error) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	conv, err := loadConversation(ctx, tx, conversationID, "", true)
	if err != nil {
		return err
	}
	if err := updateFn(conv); err != nil {
		return err
	}
	if err := saveConversation(ctx, tx, conv); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *ConversationPostgresRepository) DeleteConversation(ctx context.Context, conversationID string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM conversations WHERE id = $1`, conversationID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errConversationNotFound
	}
	return nil
}

func (r *ConversationPostgresRepository) GetConversation(ctx context.Context, conversationID string) (*conversation.Conversation, error) {
	return loadConversation(ctx, r.pool, conversationID, "", false)
}
