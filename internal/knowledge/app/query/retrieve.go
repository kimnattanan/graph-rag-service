package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/indexing"
	"github.com/sirupsen/logrus"
)

type Retrieve struct {
	QueryStr string
	TopK     *int
	Tags     *[]string
}

type RetrieveHandler decorator.QueryHandler[Retrieve, RetrieveResult]

type retrieveHandler struct {
	readModel RetrieveReadModel
	embedder  indexing.Embedder
}

func NewRetrieveHandler(
	readModel RetrieveReadModel,
	embedder indexing.Embedder,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) RetrieveHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators(
		retrieveHandler{readModel, embedder},
		logger,
		metricsClient,
	)
}

type RetrieveReadModel interface {
	Retrieve(ctx context.Context, embedding []float64, topK int, tags *[]string) (RetrieveResult, error)
}

func (h retrieveHandler) Handle(ctx context.Context, query Retrieve) (RetrieveResult, error) {
	topK := 5
	if query.TopK != nil {
		topK = *query.TopK
	}
	if topK < 1 || topK > 50 {
		return RetrieveResult{}, commonerrors.NewIncorrectInputError(
			"TopK must be between 1 and 50",
			"topk-invalid",
		)
	}
	embedding, err := h.embedder.Embed(ctx, []string{query.QueryStr})
	if err != nil {
		return RetrieveResult{}, err
	}
	if len(embedding) == 0 {
		return RetrieveResult{}, commonerrors.NewSlugError(
			"Embedding is empty",
			"embedding-empty",
		)
	}
	return h.readModel.Retrieve(ctx, embedding[0], topK, query.Tags)
}
