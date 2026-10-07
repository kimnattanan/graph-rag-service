package ports

import (
	"net/http"

	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/kimnattanan/graph-rag-service/internal/common/auth"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/common/server/httperr"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/conversation"
	openapi_types "github.com/oapi-codegen/runtime/types"
)

type HttpServer struct {
	app       app.Application
	jwtSecret string
}

func NewHttpServer(application app.Application, jwtSecret string) HttpServer {
	return HttpServer{app: application, jwtSecret: jwtSecret}
}

func (h HttpServer) ListConversations(w http.ResponseWriter, r *http.Request, params ListConversationsParams) {
	user, err := h.authorize(r)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	conversationSummaryList, err := h.app.Queries.ListConversations.Handle(r.Context(), query.ListConversations{
		UserID: user.ID,
		Limit:  params.Limit,
		Offset: params.Offset,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	summaries := make([]ConversationSummary, len(conversationSummaryList.Items))
	for i, summary := range conversationSummaryList.Items {
		summaries[i], err = appConversationSummaryToResponse(summary)
		if err != nil {
			httperr.RespondWithSlugError(err, w, r)
			return
		}
	}

	render.Respond(w, r, ConversationList{
		Items: summaries,
		Total: conversationSummaryList.Total,
	})
}

func (h HttpServer) CreateConversation(w http.ResponseWriter, r *http.Request) {
	user, err := h.authorize(r)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	var body CreateConversationRequest
	if err := render.Decode(r, &body); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	title := ""
	if body.Title != nil {
		title = *body.Title
	}

	err = h.app.Commands.CreateConversation.Handle(r.Context(), command.CreateConversation{
		ConversationID: body.Id.String(),
		UserID:         user.ID,
		Title:          title,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) DeleteConversation(w http.ResponseWriter, r *http.Request, conversationId openapi_types.UUID) {
	user, err := h.authorize(r)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	err = h.app.Commands.DeleteConversation.Handle(r.Context(), command.DeleteConversation{
		ConversationID: conversationId.String(),
		UserID:         user.ID,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) GetConversation(w http.ResponseWriter, r *http.Request, conversationId openapi_types.UUID) {
	user, err := h.authorize(r)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	appConversation, err := h.app.Queries.GetConversation.Handle(r.Context(), query.GetConversation{
		ConversationID: conversationId.String(),
		UserID:         user.ID,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	detail, err := appConversationToResponse(appConversation)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, detail)
}

func (h HttpServer) ListMessages(w http.ResponseWriter, r *http.Request, conversationId openapi_types.UUID, params ListMessagesParams) {
	user, err := h.authorize(r)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	messageList, err := h.app.Queries.ListMessages.Handle(r.Context(), query.ListMessages{
		ConversationID: conversationId.String(),
		UserID:         user.ID,
		Limit:          params.Limit,
		Offset:         params.Offset,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	messages := make([]Message, len(messageList.Items))
	for i, message := range messageList.Items {
		messages[i], err = appMessageToResponse(message)
		if err != nil {
			httperr.RespondWithSlugError(err, w, r)
			return
		}
	}

	render.Respond(w, r, MessageList{
		Items: messages,
		Total: messageList.Total,
	})
}

func (h HttpServer) SendMessage(w http.ResponseWriter, r *http.Request, conversationId openapi_types.UUID) {
	user, err := h.authorize(r)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	var body SendMessageRequest
	if err := render.Decode(r, &body); err != nil {
		httperr.BadRequest("invalid-request", err, w, r)
		return
	}

	err = h.app.Commands.SendMessage.Handle(r.Context(), command.SendMessage{
		ConversationID:  conversationId.String(),
		UserID:          user.ID,
		MessageID:       body.Id.String(),
		Content:         body.Content,
		TopK:            body.TopK,
		Tags:            body.Tags,
		HistoryCapacity: body.HistoryCapacity,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h HttpServer) GetMessage(w http.ResponseWriter, r *http.Request, conversationId openapi_types.UUID, messageId openapi_types.UUID) {
	user, err := h.authorize(r)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	appMessage, err := h.app.Queries.GetMessage.Handle(r.Context(), query.GetMessage{
		ConversationID: conversationId.String(),
		MessageID:      messageId.String(),
		UserID:         user.ID,
	})
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	message, err := appMessageToResponse(appMessage)
	if err != nil {
		httperr.RespondWithSlugError(err, w, r)
		return
	}

	render.Respond(w, r, message)
}

func (h HttpServer) authorize(r *http.Request) (auth.User, error) {
	token, err := auth.BearerToken(r)
	if err != nil {
		return auth.User{}, err
	}

	user, err := auth.Parse(h.jwtSecret, token)
	if err != nil {
		return auth.User{}, err
	}
	if !auth.HasPermission(user, auth.PermissionConversationAsk) {
		return auth.User{}, commonerrors.NewAuthorizationError("missing conversation:ask permission", "forbidden")
	}
	return user, nil
}

func appConversationSummaryToResponse(summary query.ConversationSummary) (ConversationSummary, error) {
	id, err := uuid.Parse(summary.ID)
	if err != nil {
		return ConversationSummary{}, commonerrors.NewIncorrectInputError("invalid conversation id", "conversation-id-invalid")
	}
	return ConversationSummary{
		Id:           id,
		Title:        summary.Title,
		CreatedAt:    summary.CreatedAt,
		UpdatedAt:    summary.UpdatedAt,
		MessageCount: summary.MessageCount,
	}, nil
}

func appConversationToResponse(conversation query.Conversation) (ConversationDetail, error) {
	id, err := uuid.Parse(conversation.ID)
	if err != nil {
		return ConversationDetail{}, commonerrors.NewIncorrectInputError("invalid conversation id", "conversation-id-invalid")
	}

	messages := make([]Message, len(conversation.Messages))
	for i, message := range conversation.Messages {
		messages[i], err = appMessageToResponse(message)
		if err != nil {
			return ConversationDetail{}, err
		}
	}

	return ConversationDetail{
		Id:        id,
		Title:     conversation.Title,
		CreatedAt: conversation.CreatedAt,
		UpdatedAt: conversation.UpdatedAt,
		Messages:  messages,
	}, nil
}

func appMessageToResponse(message query.Message) (Message, error) {
	id, err := uuid.Parse(message.ID)
	if err != nil {
		return Message{}, commonerrors.NewIncorrectInputError("invalid message id", "message-id-invalid")
	}
	conversationID, err := uuid.Parse(message.ConversationID)
	if err != nil {
		return Message{}, commonerrors.NewIncorrectInputError("invalid conversation id", "conversation-id-invalid")
	}
	role, err := messageRoleToResponse(message.Role)
	if err != nil {
		return Message{}, err
	}

	response := Message{
		Id:             id,
		ConversationId: conversationID,
		Role:           role,
		Content:        message.Content,
		CreatedAt:      message.CreatedAt,
	}
	if len(message.Sources) == 0 {
		return response, nil
	}

	sources, err := appSourcesToResponse(message.Sources)
	if err != nil {
		return Message{}, err
	}
	response.Sources = &sources
	return response, nil
}

func appSourcesToResponse(sources []query.Source) ([]Source, error) {
	response := make([]Source, len(sources))
	for i, source := range sources {
		documentID, err := uuid.Parse(source.DocumentID)
		if err != nil {
			return nil, commonerrors.NewIncorrectInputError("invalid document id", "document-id-invalid")
		}
		mapped := Source{
			DocumentId:    documentID,
			DocumentTitle: source.DocumentTitle,
			Text:          source.Text,
			Score:         float32(source.Score),
		}
		if source.ChunkID != "" {
			chunkID := source.ChunkID
			mapped.ChunkId = &chunkID
		}
		response[i] = mapped
	}
	return response, nil
}

func messageRoleToResponse(role conversation.Role) (MessageRole, error) {
	switch role {
	case conversation.RoleUser:
		return User, nil
	case conversation.RoleAssistant:
		return Assistant, nil
	default:
		return "", commonerrors.NewIncorrectInputError("invalid message role", "invalid-message-role")
	}
}
