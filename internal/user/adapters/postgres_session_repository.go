package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kimnattanan/graph-rag-service/internal/user/domain/user"
)

type PostgresSessionRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresSessionRepository(pool *pgxpool.Pool) *PostgresSessionRepository {
	if pool == nil {
		panic("nil pool")
	}
	return &PostgresSessionRepository{pool: pool}
}

func (r *PostgresSessionRepository) AddSession(ctx context.Context, session user.Session) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO sessions (id, user_id, expires_at, created_at)
		VALUES ($1, $2, $3, $4)
	`, session.ID(), session.UserID(), session.ExpiresAt(), session.CreatedAt())
	return err
}

func (r *PostgresSessionRepository) GetSession(ctx context.Context, sessionID string) (user.Session, error) {
	var (
		id        string
		userID    string
		expiresAt time.Time
		createdAt time.Time
	)
	err := r.pool.QueryRow(ctx, `
		SELECT id, user_id, expires_at, created_at
		FROM sessions
		WHERE id = $1
	`, sessionID).Scan(&id, &userID, &expiresAt, &createdAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return user.Session{}, user.ErrSessionNotFound
	}
	if err != nil {
		return user.Session{}, err
	}
	return user.UnmarshalSessionFromDatabase(id, userID, expiresAt, createdAt)
}

func (r *PostgresSessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, sessionID)
	return err
}
