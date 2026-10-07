package query

import (
	"context"

	"github.com/kimnattanan/graph-rag-service/internal/common/decorator"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
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
}

func NewRetrieveHandler(
	readModel RetrieveReadModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) RetrieveHandler {
	if readModel == nil {
		panic("nil readModel")
	}

	return decorator.ApplyQueryDecorators(
		retrieveHandler{readModel},
		logger,
		metricsClient,
	)
}

type RetrieveReadModel interface {
	Retrieve(ctx context.Context, queryStr string, topK int, tags *[]string) (RetrieveResult, error)
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
	return h.readModel.Retrieve(ctx, query.QueryStr, topK, query.Tags)
}
