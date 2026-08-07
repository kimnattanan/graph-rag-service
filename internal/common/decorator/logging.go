package decorator

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type queryLoggingDecorator[Q any, R any] struct {
	base   QueryHandler[Q, R]
	logger *logrus.Entry
}

func (d queryLoggingDecorator[Q, R]) Handle(ctx context.Context, query Q) (result R, err error) {
	logger := d.logger.WithFields(logrus.Fields{
		"query":      generateActionName(query),
		"query_body": fmt.Sprintf("%#v", query),
	})

	logger.Debug("Executing query")
	defer func() {
		if err == nil {
			logger.Info("Query executed successfully")
		} else {
			logger.WithError(err).Error("Failed to execute query")
		}
	}()

	return d.base.Handle(ctx, query)
}
