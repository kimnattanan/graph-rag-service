package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/conversation"
)

var (
	errConversationNotFound      = commonerrors.NewNotFoundError("conversation not found", "conversation-not-found")
	errMessageNotFound           = commonerrors.NewNotFoundError("message not found", "message-not-found")
	errConversationAlreadyExists = commonerrors.NewIncorrectInputError("conversation already exists", "conversation-already-exists")
)

type queryer interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type sourceJSON struct {
	DocumentID    string  `json:"document_id"`
	DocumentTitle string  `json:"document_title"`
	ChunkID       string  `json:"chunk_id"`
	Text          string  `json:"text"`
	Score         float64 `json:"score"`
}

func loadConversation(ctx context.Context, q queryer, conversationID string, userID string, lock bool) (*conversation.Conversation, error) {
	sql := `
		SELECT id, user_id, title, created_at, updated_at
		FROM conversations
		WHERE id = $1`
	args := []any{conversationID}
	if userID != "" {
		sql += ` AND user_id = $2`
		args = append(args, userID)
	}
	if lock {
		sql += ` FOR UPDATE`
	}

	var (
		id        string
		ownerID   string
		title     string
		createdAt time.Time
		updatedAt time.Time
	)
	err := q.QueryRow(ctx, sql, args...).Scan(&id, &ownerID, &title, &createdAt, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, errConversationNotFound
	}
	if err != nil {
		return nil, err
	}

	messages, err := loadMessages(ctx, q, id)
	if err != nil {
		return nil, err
	}
	return conversation.UnmarshalConversationFromDatabase(id, ownerID, title, createdAt, updatedAt, messages)
}

func loadMessages(ctx context.Context, q queryer, conversationID string) ([]conversation.Message, error) {
	rows, err := q.Query(ctx, `
		SELECT id, conversation_id, role, content, sources, created_at
		FROM messages
		WHERE conversation_id = $1
		ORDER BY created_at ASC, id ASC
	`, conversationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]conversation.Message, 0)
	for rows.Next() {
		message, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return messages, nil
}

func scanMessage(row pgx.Row) (conversation.Message, error) {
	var (
		id             string
		conversationID string
		role           string
		content        string
		sources        []sourceJSON
		createdAt      time.Time
	)
	if err := row.Scan(&id, &conversationID, &role, &content, &sources, &createdAt); err != nil {
		return conversation.Message{}, err
	}
	return messageFromRow(id, conversationID, role, content, sources, createdAt)
}

func messageFromRow(id string, conversationID string, role string, content string, sources []sourceJSON, createdAt time.Time) (conversation.Message, error) {
	domainSources := make([]conversation.Source, len(sources))
	for i, source := range sources {
		domainSources[i] = conversation.NewSource(
			source.DocumentID,
			source.DocumentTitle,
			source.ChunkID,
			source.Text,
			source.Score,
		)
	}
	message, err := conversation.UnmarshalMessageFromDatabase(
		id,
		conversationID,
		conversation.Role(role),
		content,
		domainSources,
		createdAt,
	)
	if err != nil {
		return conversation.Message{}, err
	}
	return *message, nil
}

func sourcesFromMessage(message conversation.Message) []sourceJSON {
	domainSources := message.Sources()
	sources := make([]sourceJSON, len(domainSources))
	for i, source := range domainSources {
		sources[i] = sourceJSON{
			DocumentID:    source.DocumentID(),
			DocumentTitle: source.DocumentTitle(),
			ChunkID:       source.ChunkID(),
			Text:          source.Text(),
			Score:         source.Score(),
		}
	}
	return sources
}

func insertMessages(ctx context.Context, q queryer, messages []conversation.Message) error {
	for _, message := range messages {
		if _, err := q.Exec(ctx, `
			INSERT INTO messages (id, conversation_id, role, content, sources, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, message.ID(), message.ConversationID(), string(message.Role()), message.Content(), sourcesFromMessage(message), message.CreatedAt()); err != nil {
			return err
		}
	}
	return nil
}

func saveConversation(ctx context.Context, q queryer, conv *conversation.Conversation) error {
	tag, err := q.Exec(ctx, `
		UPDATE conversations
		SET title = $2, updated_at = $3
		WHERE id = $1
	`, conv.ID(), conv.Title(), conv.UpdatedAt())
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return errConversationNotFound
	}

	messages := conv.Messages()
	ids := make([]string, 0, len(messages))
	for _, message := range messages {
		if _, err := q.Exec(ctx, `
			INSERT INTO messages (id, conversation_id, role, content, sources, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			ON CONFLICT (id) DO UPDATE
			SET conversation_id = EXCLUDED.conversation_id,
			    role = EXCLUDED.role,
			    content = EXCLUDED.content,
			    sources = EXCLUDED.sources,
			    created_at = EXCLUDED.created_at
		`, message.ID(), message.ConversationID(), string(message.Role()), message.Content(), sourcesFromMessage(message), message.CreatedAt()); err != nil {
			return err
		}
		ids = append(ids, message.ID())
	}

	_, err = q.Exec(ctx, `
		DELETE FROM messages
		WHERE conversation_id = $1
		  AND id <> ALL($2::text[])
	`, conv.ID(), ids)
	return err
}

func ensureConversation(ctx context.Context, q queryer, userID string, conversationID string) error {
	var exists bool
	err := q.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1 FROM conversations WHERE id = $1 AND user_id = $2
		)
	`, conversationID, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errConversationNotFound
	}
	return nil
}

func toQueryMessage(message conversation.Message) query.Message {
	domainSources := message.Sources()
	sources := make([]query.Source, len(domainSources))
	for i, source := range domainSources {
		sources[i] = query.Source{
			DocumentID:    source.DocumentID(),
			DocumentTitle: source.DocumentTitle(),
			ChunkID:       source.ChunkID(),
			Text:          source.Text(),
			Score:         source.Score(),
		}
	}
	return query.Message{
		ID:             message.ID(),
		ConversationID: message.ConversationID(),
		Role:           message.Role(),
		Content:        message.Content(),
		Sources:        sources,
		CreatedAt:      message.CreatedAt(),
	}
}

func isUniqueViolation(err error, constraint string) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == constraint
}
