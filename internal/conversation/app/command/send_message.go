package command

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/conversation"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/retrieval"
	"github.com/sirupsen/logrus"
)

type SendMessage struct {
	ConversationID  string
	UserID          string
	MessageID       string
	Content         string
	TopK            *int
	Tags            *[]string
	HistoryCapacity *int
}

type SendMessageHandler decorator.CommandHandler[SendMessage]

type sendMessageHandler struct {
	repo      conversation.Repository
	retriever retrieval.Retriever
	completer retrieval.Completer
}

func NewSendMessageHandler(
	repo conversation.Repository,
	retriever retrieval.Retriever,
	completer retrieval.Completer,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) SendMessageHandler {
	if repo == nil {
		panic("nil repo service")
	}
	if retriever == nil {
		panic("nil retriever service")
	}
	if completer == nil {
		panic("nil completer service")
	}

	return decorator.ApplyCommandDecorators(
		sendMessageHandler{repo: repo, retriever: retriever, completer: completer},
		logger,
		metricsClient,
	)
}

func (h sendMessageHandler) Handle(ctx context.Context, cmd SendMessage) error {
	return nil
}
