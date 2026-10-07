package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kimnattanan/graph-rag-service/internal/user/domain/user"
)

type PostgresUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgresUserRepository(pool *pgxpool.Pool) *PostgresUserRepository {
	if pool == nil {
		panic("nil pool")
	}
	return &PostgresUserRepository{pool: pool}
}

func (r *PostgresUserRepository) AddUser(ctx context.Context, account *user.User) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO users (id, email, username, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, account.ID(), account.Email(), account.Username(), account.PasswordHash(), string(account.Role()), account.CreatedAt(), account.UpdatedAt())
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		return user.ErrEmailAlreadyExists
	}
	return err
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	normalized, err := user.NormalizeEmail(email)
	if err != nil {
		return nil, err
	}
	return r.getUser(ctx, `SELECT id, email, username, password_hash, role, created_at, updated_at FROM users WHERE email = $1`, normalized)
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id string) (*user.User, error) {
	return r.getUser(ctx, `SELECT id, email, username, password_hash, role, created_at, updated_at FROM users WHERE id = $1`, id)
}

func (r *PostgresUserRepository) getUser(ctx context.Context, sql string, arg any) (*user.User, error) {
	var (
		id           string
		email        string
		username     string
		passwordHash string
		role         string
		createdAt    time.Time
		updatedAt    time.Time
	)
	err := r.pool.QueryRow(ctx, sql, arg).Scan(&id, &email, &username, &passwordHash, &role, &createdAt, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, user.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user.UnmarshalUserFromDatabase(id, email, username, passwordHash, role, createdAt, updatedAt)
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, id string) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return user.ErrNotFound
	}
	return nil
}
