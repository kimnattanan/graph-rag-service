package service

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kimnattanan/graph-rag-service/internal/common/auth"
	"github.com/kimnattanan/graph-rag-service/internal/user/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/user/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/user/config"
	"github.com/kimnattanan/graph-rag-service/internal/user/domain/user"
)

func TestPostgresAccountFlow(t *testing.T) {
	cfg := config.Config{
		Postgres: config.Postgres{
			Host:     "localhost",
			Port:     "5432",
			User:     "postgres",
			Password: "postgres",
			Database: "users",
			SSLMode:  "disable",
		},
		Auth: config.Auth{
			JWTSecret:       "test-secret",
			TokenTTLSeconds: int(time.Hour / time.Second),
		},
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.Postgres.URL())
	if err != nil {
		t.Fatal(err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Skipf("postgres is not available: %v", err)
	}
	pool.Close()

	application, cleanup := NewApplication(ctx, &cfg)
	t.Cleanup(cleanup)

	email := uuid.NewString() + "@example.com"
	if err := application.Commands.Register.Handle(ctx, command.Register{
		Email:    email,
		Username: "ada",
		Password: "correct-horse",
	}); err != nil {
		t.Fatal(err)
	}

	loggedIn, err := application.Queries.Login.Handle(ctx, query.Login{
		Email:    email,
		Password: "correct-horse",
	})
	if err != nil {
		t.Fatal(err)
	}
	if loggedIn.User.Role != string(user.RoleUser) {
		t.Fatalf("role = %s", loggedIn.User.Role)
	}

	principal, err := application.Queries.Authenticate.Handle(ctx, query.Authenticate{Token: loggedIn.AccessToken})
	if err != nil {
		t.Fatal(err)
	}
	if principal.Role != auth.RoleUser || !auth.HasPermission(principal, auth.PermissionConversationAsk) {
		t.Fatalf("principal = %#v", principal)
	}

	if err := application.Commands.Logout.Handle(ctx, command.Logout{Token: loggedIn.AccessToken}); err != nil {
		t.Fatal(err)
	}
	if _, err := application.Queries.Authenticate.Handle(ctx, query.Authenticate{Token: loggedIn.AccessToken}); err == nil {
		t.Fatal("expected logged out token to be rejected")
	}

	loggedIn, err = application.Queries.Login.Handle(ctx, query.Login{
		Email:    email,
		Password: "correct-horse",
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := application.Commands.DeleteAccount.Handle(ctx, command.DeleteAccount{UserID: loggedIn.User.ID}); err != nil {
		t.Fatal(err)
	}
	if _, err := application.Queries.Login.Handle(ctx, query.Login{
		Email:    email,
		Password: "correct-horse",
	}); err != user.ErrInvalidCredentials {
		t.Fatalf("login after delete = %v", err)
	}
}
