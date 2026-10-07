package service

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/common/metrics"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/adapters"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/config"
)

func NewApplication(ctx context.Context, cfg *config.Config) (app.Application, func()) {
	pool, err := pgxpool.New(ctx, cfg.Postgres.URL())
	if err != nil {
		panic(err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		panic(fmt.Errorf("postgres: %w", err))
	}

	conversationRepository := adapters.NewConversationPostgresRepository(pool)
	conversationReadModel := adapters.NewConversationPostgresReadModel(pool)
	retriever := adapters.NewKnowledgeRetriever(cfg.Knowledge.GRPCAddress)
	completer := adapters.NewLLMCompleter(cfg.LLM)

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	return app.Application{
			Commands: app.Commands{
				CreateConversation: command.NewCreateConversationHandler(conversationRepository, logger, metricsClient),
				DeleteConversation: command.NewDeleteConversationHandler(conversationRepository, logger, metricsClient),
				SendMessage:        command.NewSendMessageHandler(conversationRepository, retriever, completer, logger, metricsClient),
			},
			Queries: app.Queries{
				ListConversations: query.NewListConversationsHandler(conversationReadModel, logger, metricsClient),
				GetConversation:   query.NewGetConversationHandler(conversationReadModel, logger, metricsClient),
				ListMessages:      query.NewListMessagesHandler(conversationReadModel, logger, metricsClient),
				GetMessage:        query.NewGetMessageHandler(conversationReadModel, logger, metricsClient),
			},
		}, func() {
			pool.Close()
		}
}
