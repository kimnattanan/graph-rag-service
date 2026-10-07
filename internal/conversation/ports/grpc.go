package ports

import (
	"context"

	conversationpb "github.com/kimnattanan/graph-rag-service/internal/common/genproto/conversation"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GrpcServer struct {
	conversationpb.UnimplementedConversationServiceServer
	app app.Application
}

func NewGrpcServer(application app.Application) GrpcServer {
	return GrpcServer{app: application}
}

func (g GrpcServer) ListConversations(ctx context.Context, request *conversationpb.ListConversationsRequest) (*conversationpb.ListConversationsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (g GrpcServer) CreateConversation(ctx context.Context, request *conversationpb.CreateConversationRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (g GrpcServer) GetConversation(ctx context.Context, request *conversationpb.GetConversationRequest) (*conversationpb.Conversation, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (g GrpcServer) DeleteConversation(ctx context.Context, request *conversationpb.DeleteConversationRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (g GrpcServer) ListMessages(ctx context.Context, request *conversationpb.ListMessagesRequest) (*conversationpb.ListMessagesResponse, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (g GrpcServer) SendMessage(ctx context.Context, request *conversationpb.SendMessageRequest) (*emptypb.Empty, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}

func (g GrpcServer) GetMessage(ctx context.Context, request *conversationpb.GetMessageRequest) (*conversationpb.Message, error) {
	return nil, status.Error(codes.Unimplemented, "not implemented")
}
