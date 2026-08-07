package decorator

import (
	"context"

	"github.com/sirupsen/logrus"
)

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}

func ApplyQueryDecorators[Q any, R any](handler QueryHandler[Q, R], logger *logrus.Entry, metricsClient MetricsClient) QueryHandler[Q, R] {
	return queryLoggingDecorator[Q, R]{
		base: queryMetricsDecorator[Q, R]{
			base:   handler,
			client: metricsClient,
		},
		logger: logger,
	}
}
