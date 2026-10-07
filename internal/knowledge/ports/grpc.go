package ports

import (
	"context"
	stderrors "errors"
	"time"

	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/common/genproto/knowledge"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type GrpcServer struct {
	knowledge.UnimplementedKnowledgeServiceServer
	app       app.Application
}

func NewGrpcServer(application app.Application) GrpcServer {
	return GrpcServer{app: application}
}

func (g GrpcServer) ListDocuments(ctx context.Context, request *knowledge.ListDocumentsRequest) (*knowledge.ListDocumentsResponse, error) {
	indexStatus, err := listIndexStatusFilter(request)
	if err != nil {
		return nil, grpcStatus(err)
	}

	var limit *int
	if request.GetLimit() != 0 {
		value := int(request.GetLimit())
		limit = &value
	}
	offset := int(request.GetOffset())

	var tag *string
	if request.GetTag() != "" {
		value := request.GetTag()
		tag = &value
	}

	documentSummaryList, err := g.app.Queries.ListDocuments.Handle(ctx, query.ListDocuments{
		Limit:       limit,
		Offset:      &offset,
		Tag:         tag,
		IndexStatus: indexStatus,
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	documents := make([]*knowledge.DocumentSummary, len(documentSummaryList.Items))
	for i, documentSummary := range documentSummaryList.Items {
		documents[i], err = appDocumentSummaryToProto(documentSummary)
		if err != nil {
			return nil, grpcStatus(err)
		}
	}

	return &knowledge.ListDocumentsResponse{
		Documents: documents,
		Total:     int32(documentSummaryList.Total),
	}, nil
}

func (g GrpcServer) CreateDocument(ctx context.Context, request *knowledge.CreateDocumentRequest) (*emptypb.Empty, error) {
	err := g.app.Commands.CreateDocument.Handle(ctx, command.CreateDocument{
		DocumentID: request.GetId(),
		Title:      request.GetTitle(),
		Content:    request.GetContent(),
		Tags:       request.GetTags(),
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	return &emptypb.Empty{}, nil
}

func (g GrpcServer) GetDocument(ctx context.Context, request *knowledge.GetDocumentRequest) (*knowledge.Document, error) {
	appDoc, err := g.app.Queries.GetDocument.Handle(ctx, query.GetDocument{
		DocumentID: request.GetDocumentId(),
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	doc, err := appDocumentToProto(appDoc)
	if err != nil {
		return nil, grpcStatus(err)
	}

	return doc, nil
}

func (g GrpcServer) UpdateDocument(ctx context.Context, request *knowledge.UpdateDocumentRequest) (*emptypb.Empty, error) {
	var title *string
	if request.GetTitle() != "" {
		value := request.GetTitle()
		title = &value
	}
	var content *string
	if request.GetContent() != "" {
		value := request.GetContent()
		content = &value
	}
	var tags *[]string
	if tagsValue := request.GetTags(); tagsValue != nil {
		tags = &tagsValue
	}

	err := g.app.Commands.UpdateDocument.Handle(ctx, command.UpdateDocument{
		DocumentID: request.GetDocumentId(),
		Title:      title,
		Content:    content,
		Tags:       tags,
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	return &emptypb.Empty{}, nil
}

func (g GrpcServer) DeleteDocument(ctx context.Context, request *knowledge.DeleteDocumentRequest) (*emptypb.Empty, error) {
	err := g.app.Commands.DeleteDocument.Handle(ctx, command.DeleteDocument{
		DocumentID: request.GetDocumentId(),
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	return &emptypb.Empty{}, nil
}

func (g GrpcServer) ReindexDocument(ctx context.Context, request *knowledge.ReindexDocumentRequest) (*emptypb.Empty, error) {
	err := g.app.Commands.ReindexDocument.Handle(ctx, command.ReindexDocument{
		DocumentID: request.GetDocumentId(),
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	return &emptypb.Empty{}, nil
}

func (g GrpcServer) GetDocumentIndexStatus(ctx context.Context, request *knowledge.GetDocumentIndexStatusRequest) (*knowledge.GetDocumentIndexStatusResponse, error) {
	indexStatus, err := g.app.Queries.GetDocumentIndexStatus.Handle(ctx, query.GetDocumentIndexStatus{
		DocumentID: request.GetDocumentId(),
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	job, err := appDocumentIndexStatusToProto(indexStatus)
	if err != nil {
		return nil, grpcStatus(err)
	}

	return &knowledge.GetDocumentIndexStatusResponse{Job: job}, nil
}

func (g GrpcServer) Retrieve(ctx context.Context, request *knowledge.RetrieveRequest) (*knowledge.RetrieveResponse, error) {
	var topK *int
	if request.GetTopK() != 0 {
		value := int(request.GetTopK())
		topK = &value
	}
	var tags *[]string
	if tagsValue := request.GetTags(); tagsValue != nil {
		tags = &tagsValue
	}

	result, err := g.app.Queries.Retrieve.Handle(ctx, query.Retrieve{
		QueryStr: request.GetQuery(),
		TopK:     topK,
		Tags:     tags,
	})
	if err != nil {
		return nil, grpcStatus(err)
	}

	return &knowledge.RetrieveResponse{
		Query:  request.GetQuery(),
		Chunks: appRetrievedChunksToProto(result.Chunks),
	}, nil
}

func listIndexStatusFilter(request *knowledge.ListDocumentsRequest) (*document.IndexStatus, error) {
	field := request.ProtoReflect().Descriptor().Fields().ByName("index_status")
	// Proto3 does not track presence for the zero enum value, so an omitted
	// status and an explicit pending status both look unset.
	if field == nil || !request.ProtoReflect().Has(field) {
		return nil, nil
	}

	indexStatus, err := protoIndexStatusToDomain(request.GetIndexStatus())
	if err != nil {
		return nil, err
	}
	return &indexStatus, nil
}

func protoIndexStatusToDomain(indexStatus knowledge.IndexStatus) (document.IndexStatus, error) {
	switch indexStatus {
	case knowledge.IndexStatus_INDEX_STATUS_PENDING:
		return document.IndexStatusPending, nil
	case knowledge.IndexStatus_INDEX_STATUS_INDEXING:
		return document.IndexStatusIndexing, nil
	case knowledge.IndexStatus_INDEX_STATUS_COMPLETED:
		return document.IndexStatusCompleted, nil
	case knowledge.IndexStatus_INDEX_STATUS_FAILED:
		return document.IndexStatusFailed, nil
	default:
		return 0, commonerrors.NewIncorrectInputError("Invalid index status", "index-status-invalid")
	}
}

func domainIndexStatusToProto(indexStatus document.IndexStatus) (knowledge.IndexStatus, error) {
	switch indexStatus {
	case document.IndexStatusPending:
		return knowledge.IndexStatus_INDEX_STATUS_PENDING, nil
	case document.IndexStatusIndexing:
		return knowledge.IndexStatus_INDEX_STATUS_INDEXING, nil
	case document.IndexStatusCompleted:
		return knowledge.IndexStatus_INDEX_STATUS_COMPLETED, nil
	case document.IndexStatusFailed:
		return knowledge.IndexStatus_INDEX_STATUS_FAILED, nil
	default:
		return 0, commonerrors.NewIncorrectInputError("Invalid index status", "index-status-invalid")
	}
}

func appDocumentToProto(d query.Document) (*knowledge.Document, error) {
	indexStatus, err := domainIndexStatusToProto(d.IndexStatus)
	if err != nil {
		return nil, err
	}
	return &knowledge.Document{
		Id:          d.ID,
		Title:       d.Title,
		Content:     d.Content,
		Tags:        d.Tags,
		IndexStatus: indexStatus,
		CreatedAt:   timestamppb.New(d.CreatedAt),
		UpdatedAt:   timestamppb.New(d.UpdatedAt),
	}, nil
}

func appDocumentSummaryToProto(s query.DocumentSummary) (*knowledge.DocumentSummary, error) {
	indexStatus, err := domainIndexStatusToProto(s.IndexStatus)
	if err != nil {
		return nil, err
	}
	return &knowledge.DocumentSummary{
		Id:          s.ID,
		Title:       s.Title,
		Tags:        s.Tags,
		IndexStatus: indexStatus,
		CreatedAt:   timestamppb.New(s.CreatedAt),
		UpdatedAt:   timestamppb.New(s.UpdatedAt),
	}, nil
}

func appDocumentIndexStatusToProto(s query.DocumentIndexStatus) (*knowledge.IndexJob, error) {
	indexStatus, err := domainIndexStatusToProto(s.Status)
	if err != nil {
		return nil, err
	}
	errorMessage := ""
	if s.ErrorMessage != nil {
		errorMessage = *s.ErrorMessage
	}
	return &knowledge.IndexJob{
		DocumentId:   s.DocumentID,
		Status:       indexStatus,
		StartedAt:    timePtrToProto(s.StartedAt),
		FinishedAt:   timePtrToProto(s.FinishedAt),
		ErrorMessage: errorMessage,
	}, nil
}

func appRetrievedChunksToProto(chunks []query.RetrievedChunk) []*knowledge.RetrievedChunk {
	response := make([]*knowledge.RetrievedChunk, len(chunks))
	for i, chunk := range chunks {
		chunkID := ""
		if chunk.ChunkID != nil {
			chunkID = *chunk.ChunkID
		}
		var graphPath []string
		if chunk.GraphPath != nil {
			graphPath = *chunk.GraphPath
		}
		response[i] = &knowledge.RetrievedChunk{
			DocumentId:    chunk.DocumentID,
			DocumentTitle: chunk.DocumentTitle,
			ChunkId:       chunkID,
			Text:          chunk.Text,
			Score:         float32(chunk.Score),
			GraphPath:     graphPath,
		}
	}
	return response
}

func timePtrToProto(t *time.Time) *timestamppb.Timestamp {
	if t == nil || t.IsZero() {
		return nil
	}
	return timestamppb.New(*t)
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
	}

	return status.Error(code, slugError.Error())
}
