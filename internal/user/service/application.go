package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/common/metrics"
	"github.com/kimnattanan/graph-rag-service/internal/user/adapters"
	"github.com/kimnattanan/graph-rag-service/internal/user/app"
	"github.com/kimnattanan/graph-rag-service/internal/user/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/user/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/user/config"
)

func NewApplication(ctx context.Context, cfg *config.Config) (app.Application, func()) {
	if cfg.Auth.TokenTTLSeconds <= 0 {
		panic("JWT_TTL_SECONDS must be positive")
	}

	pool, err := pgxpool.New(ctx, cfg.Postgres.URL())
	if err != nil {
		panic(err)
	}
	if err := pool.Ping(ctx); err != nil {
		panic(fmt.Errorf("postgres: %w", err))
	}
	if err := startupPostgres(ctx, pool); err != nil {
		panic(err)
	}

	userRepository := adapters.NewPostgresUserRepository(pool)
	sessionRepository := adapters.NewPostgresSessionRepository(pool)
	if err := ensureAdmin(ctx, userRepository, cfg.Admin); err != nil {
		panic(err)
	}

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}
	secret := cfg.Auth.JWTSecret
	ttl := cfg.Auth.TokenTTL()

	return app.Application{
			Commands: app.Commands{
				Register:      command.NewRegisterHandler(userRepository, logger, metricsClient),
				Logout:        command.NewLogoutHandler(sessionRepository, secret, logger, metricsClient),
				DeleteAccount: command.NewDeleteAccountHandler(userRepository, logger, metricsClient),
			},
			Queries: app.Queries{
				Login:          query.NewLoginHandler(userRepository, sessionRepository, secret, ttl, logger, metricsClient),
				Authenticate:   query.NewAuthenticateHandler(userRepository, sessionRepository, secret, logger, metricsClient),
				GetCurrentUser: query.NewGetCurrentUserHandler(userRepository, logger, metricsClient),
			},
		}, func() {
			pool.Close()
		}
}
