package ports

import (
	"context"
	stderrors "errors"
	"strings"

	"github.com/kimnattanan/graph-rag-service/internal/common/auth"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	conversationpb "github.com/kimnattanan/graph-rag-service/internal/common/genproto/conversation"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/domain/conversation"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	conversationpb.UnimplementedConversationServiceServer
	app       app.Application
	jwtSecret string
}

func NewGrpcServer(application app.Application, jwtSecret string) GrpcServer {
	return GrpcServer{app: application, jwtSecret: jwtSecret}
}

func (g GrpcServer) ListConversations(ctx context.Context, request *conversationpb.ListConversationsRequest) (*conversationpb.ListConversationsResponse, error) {
	user, err := g.authorize(ctx)
	if err != nil {
		return nil, grpcStatus(err)
	}

	offset := int(request.GetOffset())
	conversationSummaryList, err := g.app.Queries.ListConversations.Handle(ctx, query.ListConversations{
		UserID: user.ID,
		Limit:  optionalInt(request.GetLimit()),
		Offset: &offset,
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	conversations := make([]*conversationpb.ConversationSummary, len(conversationSummaryList.Items))
	for i, summary := range conversationSummaryList.Items {
		conversations[i] = appConversationSummaryToProto(summary)
	}

	return &conversationpb.ListConversationsResponse{
		Conversations: conversations,
		Total:         int32(conversationSummaryList.Total),
	}, nil
}

func (g GrpcServer) CreateConversation(ctx context.Context, request *conversationpb.CreateConversationRequest) (*emptypb.Empty, error) {
	user, err := g.authorize(ctx)
	if err != nil {
		return nil, grpcStatus(err)
	}

	err = g.app.Commands.CreateConversation.Handle(ctx, command.CreateConversation{
		ConversationID: request.GetId(),
		UserID:         user.ID,
		Title:          request.GetTitle(),
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	return &emptypb.Empty{}, nil
}

func (g GrpcServer) GetConversation(ctx context.Context, request *conversationpb.GetConversationRequest) (*conversationpb.Conversation, error) {
	user, err := g.authorize(ctx)
	if err != nil {
		return nil, grpcStatus(err)
	}

	appConversation, err := g.app.Queries.GetConversation.Handle(ctx, query.GetConversation{
		ConversationID: request.GetConversationId(),
		UserID:         user.ID,
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	conversation, err := appConversationToProto(appConversation)
	if err != nil {
		return nil, grpcStatus(err)
	}

	return conversation, nil
}

func (g GrpcServer) DeleteConversation(ctx context.Context, request *conversationpb.DeleteConversationRequest) (*emptypb.Empty, error) {
	user, err := g.authorize(ctx)
	if err != nil {
		return nil, grpcStatus(err)
	}

	err = g.app.Commands.DeleteConversation.Handle(ctx, command.DeleteConversation{
		ConversationID: request.GetConversationId(),
		UserID:         user.ID,
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	return &emptypb.Empty{}, nil
}

func (g GrpcServer) ListMessages(ctx context.Context, request *conversationpb.ListMessagesRequest) (*conversationpb.ListMessagesResponse, error) {
	user, err := g.authorize(ctx)
	if err != nil {
		return nil, grpcStatus(err)
	}

	offset := int(request.GetOffset())
	messageList, err := g.app.Queries.ListMessages.Handle(ctx, query.ListMessages{
		ConversationID: request.GetConversationId(),
		UserID:         user.ID,
		Limit:          optionalInt(request.GetLimit()),
		Offset:         &offset,
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	messages := make([]*conversationpb.Message, len(messageList.Items))
	for i, message := range messageList.Items {
		messages[i], err = appMessageToProto(message)
		if err != nil {
			return nil, grpcStatus(err)
		}
	}

	return &conversationpb.ListMessagesResponse{
		Messages: messages,
		Total:    int32(messageList.Total),
	}, nil
}

func (g GrpcServer) SendMessage(ctx context.Context, request *conversationpb.SendMessageRequest) (*emptypb.Empty, error) {
	user, err := g.authorize(ctx)
	if err != nil {
		return nil, grpcStatus(err)
	}

	var tags *[]string
	if tagsValue := request.GetTags(); tagsValue != nil {
		tags = &tagsValue
	}

	err = g.app.Commands.SendMessage.Handle(ctx, command.SendMessage{
		ConversationID:  request.GetConversationId(),
		UserID:          user.ID,
		MessageID:       request.GetId(),
		Content:         request.GetContent(),
		TopK:            optionalInt(request.GetTopK()),
		Tags:            tags,
		HistoryCapacity: optionalInt(request.GetHistoryCapacity()),
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	return &emptypb.Empty{}, nil
}

func (g GrpcServer) GetMessage(ctx context.Context, request *conversationpb.GetMessageRequest) (*conversationpb.Message, error) {
	user, err := g.authorize(ctx)
	if err != nil {
		return nil, grpcStatus(err)
	}

	appMessage, err := g.app.Queries.GetMessage.Handle(ctx, query.GetMessage{
		ConversationID: request.GetConversationId(),
		MessageID:      request.GetMessageId(),
		UserID:         user.ID,
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	message, err := appMessageToProto(appMessage)
	if err != nil {
		return nil, grpcStatus(err)
	}

	return message, nil
}

func (g GrpcServer) authorize(ctx context.Context) (auth.User, error) {
	token, err := bearerTokenFromContext(ctx)
	if err != nil {
		return auth.User{}, err
	}

	user, err := auth.Parse(g.jwtSecret, token)
	if err != nil {
		return auth.User{}, err
	}
	if !auth.HasPermission(user, auth.PermissionConversationAsk) {
		return auth.User{}, commonerrors.NewAuthorizationError("missing conversation:ask permission", "forbidden")
	}
	return user, nil
}

func bearerTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", commonerrors.NewAuthorizationError("missing bearer token", "unauthenticated")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return "", commonerrors.NewAuthorizationError("missing bearer token", "unauthenticated")
	}

	scheme, token, ok := strings.Cut(values[0], " ")
	if !ok || !strings.EqualFold(scheme, "bearer") || strings.TrimSpace(token) == "" {
		return "", commonerrors.NewAuthorizationError("missing bearer token", "unauthenticated")
	}
	return strings.TrimSpace(token), nil
}

func optionalInt(value int32) *int {
	if value == 0 {
		return nil
	}
	converted := int(value)
	return &converted
}

func appConversationSummaryToProto(summary query.ConversationSummary) *conversationpb.ConversationSummary {
	return &conversationpb.ConversationSummary{
		Id:           summary.ID,
		Title:        summary.Title,
		CreatedAt:    timestamppb.New(summary.CreatedAt),
		UpdatedAt:    timestamppb.New(summary.UpdatedAt),
		MessageCount: int32(summary.MessageCount),
	}
}

func appConversationToProto(appConversation query.Conversation) (*conversationpb.Conversation, error) {
	messages := make([]*conversationpb.Message, len(appConversation.Messages))
	for i, message := range appConversation.Messages {
		mapped, err := appMessageToProto(message)
		if err != nil {
			return nil, err
		}
		messages[i] = mapped
	}

	return &conversationpb.Conversation{
		Id:        appConversation.ID,
		Title:     appConversation.Title,
		CreatedAt: timestamppb.New(appConversation.CreatedAt),
		UpdatedAt: timestamppb.New(appConversation.UpdatedAt),
		Messages:  messages,
	}, nil
}

func appMessageToProto(message query.Message) (*conversationpb.Message, error) {
	role, err := messageRoleToProto(message.Role)
	if err != nil {
		return nil, err
	}

	sources := make([]*conversationpb.Source, len(message.Sources))
	for i, source := range message.Sources {
		sources[i] = &conversationpb.Source{
			DocumentId:    source.DocumentID,
			DocumentTitle: source.DocumentTitle,
			ChunkId:       source.ChunkID,
			Text:          source.Text,
			Score:         float32(source.Score),
		}
	}

	return &conversationpb.Message{
		Id:             message.ID,
		ConversationId: message.ConversationID,
		Role:           role,
		Content:        message.Content,
		Sources:        sources,
		CreatedAt:      timestamppb.New(message.CreatedAt),
	}, nil
}

func messageRoleToProto(role conversation.Role) (conversationpb.MessageRole, error) {
	switch role {
	case conversation.RoleUser:
		return conversationpb.MessageRole_MESSAGE_ROLE_USER, nil
	case conversation.RoleAssistant:
		return conversationpb.MessageRole_MESSAGE_ROLE_ASSISTANT, nil
	default:
		return conversationpb.MessageRole_MESSAGE_ROLE_UNSPECIFIED, commonerrors.NewIncorrectInputError("invalid message role", "invalid-message-role")
	}
}

func grpcStatus(err error) error {
	var slugError commonerrors.SlugError
	if !stderrors.As(err, &slugError) {
		return status.Error(codes.Internal, "internal server error")
	}

	code := codes.Internal
	switch slugError.ErrorType() {
	case commonerrors.ErrorTypeAuthorization:
		code = codes.Unauthenticated
	case commonerrors.ErrorTypeIncorrectInput:
		code = codes.InvalidArgument
	case commonerrors.ErrorTypeNotFound:
		code = codes.NotFound
	case commonerrors.ErrorTypeConflict:
		code = codes.AlreadyExists
	}

	return status.Error(code, slugError.Error())
}
