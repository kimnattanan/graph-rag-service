package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
)

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func ApplyCommandDecorators[C any](handler CommandHandler[C], logger *logrus.Entry, metricsClient MetricsClient) CommandHandler[C] {
	return commandLoggingDecorator[C]{
		base: commandMetricsDecorator[C]{
			base:   handler,
			client: metricsClient,
		},
		logger: logger,
	}
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
