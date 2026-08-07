package decorator

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type MetricsClient interface {
	Inc(key string, value int)
}

type queryMetricsDecorator[Q any, R any] struct {
	base   QueryHandler[Q, R]
	client MetricsClient
}

func (d queryMetricsDecorator[Q, R]) Handle(ctx context.Context, query Q) (result R, err error) {
	start := time.Now()

	actionName := strings.ToLower(generateActionName(query))

	defer func() {
		end := time.Since(start)

		d.client.Inc(fmt.Sprintf("querys.%s.duration", actionName), int(end.Seconds()))

		if err == nil {
			d.client.Inc(fmt.Sprintf("querys.%s.success", actionName), 1)
		} else {
			d.client.Inc(fmt.Sprintf("querys.%s.failure", actionName), 1)
		}
	}()

	return d.base.Handle(ctx, query)
}
