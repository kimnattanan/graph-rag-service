package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/user/config"
	"github.com/kimnattanan/graph-rag-service/internal/user/domain/user"
)

func startupPostgres(ctx context.Context, pool *pgxpool.Pool) error {
	statements := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			email TEXT NOT NULL UNIQUE,
			username TEXT NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL CHECK (role IN ('user', 'admin')),
			created_at TIMESTAMPTZ NOT NULL,
			updated_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
			expires_at TIMESTAMPTZ NOT NULL,
			created_at TIMESTAMPTZ NOT NULL
		)`,
		`CREATE INDEX IF NOT EXISTS sessions_user_id_idx ON sessions (user_id)`,
	}
	for _, statement := range statements {
		if _, err := pool.Exec(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

func ensureAdmin(ctx context.Context, repo user.Repository, cfg config.Admin) error {
	if cfg.Email == "" && cfg.Username == "" && cfg.Password == "" {
		return nil
	}
	if cfg.Email == "" || cfg.Username == "" || cfg.Password == "" {
		return fmt.Errorf("ADMIN_EMAIL, ADMIN_USERNAME, and ADMIN_PASSWORD must all be set")
	}

	email, err := user.NormalizeEmail(cfg.Email)
	if err != nil {
		return err
	}
	_, err = repo.GetUserByEmail(ctx, email)
	if err == nil {
		logrus.Info("admin user already exists")
		return nil
	}
	if !errors.Is(err, user.ErrNotFound) {
		return err
	}

	account, err := user.NewAdmin(uuid.NewString(), email, cfg.Username, cfg.Password)
	if err != nil {
		return err
	}
	if err := repo.AddUser(ctx, account); err != nil {
		return err
	}
	logrus.WithField("email", account.Email()).Info("created admin user")
	return nil
}
