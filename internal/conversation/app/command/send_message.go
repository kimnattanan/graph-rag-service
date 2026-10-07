package command

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
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
	// create user message
	conv, err := h.repo.GetConversation(ctx, cmd.ConversationID)
	if err != nil {
		return err
	}
	if conv.UserID() != cmd.UserID {
		return conversation.ErrConversationNotOwnedByUser
	}
	msg, err := conversation.NewMessage(
		cmd.MessageID,
		cmd.ConversationID,
		conversation.RoleUser,
		cmd.Content,
		[]conversation.Source{},
	)
	if err != nil {
		return err
	}
	if err := conv.AddMessage(*msg); err != nil {
		return err
	}

	// retrieve knowledge chunks and generate prompt
	topK := 5
	if cmd.TopK != nil {
		topK = *cmd.TopK
	}
	if topK < 1 || topK > 50 {
		return commonerrors.NewIncorrectInputError(
			"TopK must be between 1 and 50",
			"topk-invalid",
		)
	}
	tags := []string{}
	if cmd.Tags != nil {
		tags = *cmd.Tags
	}
	historyCapacity := 3
	if cmd.HistoryCapacity != nil {
		historyCapacity = *cmd.HistoryCapacity
	}
	if historyCapacity < 1 || historyCapacity > 100 {
		return commonerrors.NewIncorrectInputError(
			"HistoryCapacity must be between 1 and 100",
			"historycapacity-invalid",
		)
	}
	historyMessages := h.combineHistoryMessages(conv.Messages(), historyCapacity)
	chunks, err := h.retriever.Retrieve(ctx, historyMessages, topK, tags)
	if err != nil {
		return err
	}
	prompt := h.generatePrompt(historyMessages, chunks)
	sources := make([]conversation.Source, len(chunks))
	for i, chunk := range chunks {
		sources[i] = conversation.NewSource(
			chunk.DocumentID,
			chunk.DocumentTitle,
			chunk.ChunkID,
			chunk.Text,
			chunk.Score,
		)
	}

	// create assistant message
	reply, err := h.completer.Complete(ctx, prompt)
	if err != nil {
		return err
	}
	replyMsg, err := conversation.NewMessage(
		uuid.NewString(),
		cmd.ConversationID,
		conversation.RoleAssistant,
		reply,
		sources,
	)
	if err != nil {
		return err
	}
	return h.repo.UpdateConversation(ctx, cmd.ConversationID, func(conv *conversation.Conversation) error {
		if err := conv.AddMessage(*msg); err != nil {
			return err
		}
		if err := conv.AddMessage(*replyMsg); err != nil {
			return err
		}
		return nil
	})
}

func (h sendMessageHandler) combineHistoryMessages(messages []conversation.Message, capacity int) string {
	query := ""
	for i := max(len(messages)-capacity, 0); i < len(messages); i++ {
		msg := messages[i]
		query += fmt.Sprintf("[%s]: %s\n\n", msg.Role(), msg.Content())
	}
	return query
}

const unknownAnswer = "I don't know."

func (h sendMessageHandler) generatePrompt(historyMessages string, chunks []retrieval.Chunk) string {
	var b strings.Builder
	b.WriteString(`Answer the user's latest question using only the knowledge passages and the conversation below.

Rules:
- Use only facts stated in the knowledge passages.
- Conversation history is context for what the user is asking. It is not a source of facts.
- If the passages do not contain the answer, reply with exactly: ` + unknownAnswer + `
- Do not guess, and do not add facts from outside the passages.
- Answer in the same language as the user's latest question.

Knowledge:
`)
	if len(chunks) == 0 {
		b.WriteString("(none)\n")
	}
	for i, chunk := range chunks {
		fmt.Fprintf(&b, "[%d] %s\n%s\n\n", i+1, chunk.DocumentTitle, chunk.Text)
	}
	b.WriteString("Conversation:\n")
	b.WriteString(historyMessages)
	return b.String()
}
