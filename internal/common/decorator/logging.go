package decorator

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"
)

type commandLoggingDecorator[C any] struct {
	base   CommandHandler[C]
	logger *logrus.Entry
}

func (d commandLoggingDecorator[C]) Handle(ctx context.Context, cmd C) (err error) {
	logger := d.logger.WithFields(logrus.Fields{
		"command":      generateActionName(cmd),
		"command_body": fmt.Sprintf("%#v", cmd),
	})

	logger.Debug("Executing command")
	defer func() {
		if err == nil {
			logger.Info("Command executed successfully")
		} else {
			logger.WithError(err).Error("Failed to execute command")
		}
	}()

	return d.base.Handle(ctx, cmd)
}

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
