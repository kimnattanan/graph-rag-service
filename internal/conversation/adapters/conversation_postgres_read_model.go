package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
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
	var total int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM conversations WHERE user_id = $1
	`, userID).Scan(&total)
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT
			c.id,
			c.title,
			c.created_at,
			c.updated_at,
			(SELECT COUNT(*) FROM messages m WHERE m.conversation_id = c.id) AS message_count
		FROM conversations c
		WHERE c.user_id = $1
		ORDER BY c.updated_at DESC, c.id DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]query.ConversationSummary, 0)
	for rows.Next() {
		var item query.ConversationSummary
		if err := rows.Scan(&item.ID, &item.Title, &item.CreatedAt, &item.UpdatedAt, &item.MessageCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &query.ConversationSummaryList{Items: items, Total: total}, nil
}

func (r *ConversationPostgresReadModel) GetConversation(ctx context.Context, userID string, conversationID string) (*conversation.Conversation, error) {
	return loadConversation(ctx, r.pool, conversationID, userID, false)
}

func (r *ConversationPostgresReadModel) ListMessages(ctx context.Context, userID string, conversationID string, limit int, offset int) (*query.MessageList, error) {
	if err := ensureConversation(ctx, r.pool, userID, conversationID); err != nil {
		return nil, err
	}

	var total int
	err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM messages WHERE conversation_id = $1
	`, conversationID).Scan(&total)
	if err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `
		SELECT id, conversation_id, role, content, sources, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC, id ASC
		LIMIT $2 OFFSET $3
	`, conversationID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]query.Message, 0)
	for rows.Next() {
		message, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, toQueryMessage(message))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return &query.MessageList{Items: items, Total: total}, nil
}

func (r *ConversationPostgresReadModel) GetMessage(ctx context.Context, userID string, conversationID string, messageID string) (*query.Message, error) {
	if err := ensureConversation(ctx, r.pool, userID, conversationID); err != nil {
		return nil, err
	}

	var (
		id                    string
		messageConversationID string
		role                  string
		content               string
		sources               []sourceJSON
		createdAt             time.Time
	)
	err := r.pool.QueryRow(ctx, `
		SELECT id, conversation_id, role, content, sources, created_at
		FROM messages
		WHERE conversation_id = $1 AND id = $2
	`, conversationID, messageID).Scan(&id, &messageConversationID, &role, &content, &sources, &createdAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errMessageNotFound
	}
	if err != nil {
		return nil, err
	}

	message, err := messageFromRow(id, messageConversationID, role, content, sources, createdAt)
	if err != nil {
		return nil, err
	}
	mapped := toQueryMessage(message)
	return &mapped, nil
}
