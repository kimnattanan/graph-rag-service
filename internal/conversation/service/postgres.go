package service

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func startupPostgres(ctx context.Context, pool *pgxpool.Pool) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS conversations (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			title TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS conversations_user_id_updated_at_idx
			ON conversations (user_id, updated_at DESC, id DESC)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			conversation_id TEXT NOT NULL REFERENCES conversations (id) ON DELETE CASCADE,
			role TEXT NOT NULL CHECK (role IN ('User', 'Assistant')),
			content TEXT NOT NULL,
			sources JSONB NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS messages_conversation_id_created_at_idx
			ON messages (conversation_id, created_at ASC, id ASC)`,
	}
	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}
