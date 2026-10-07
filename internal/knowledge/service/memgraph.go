package service

import (
	"context"
	"fmt"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"github.com/sirupsen/logrus"

	"github.com/kimnattanan/graph-rag-service/internal/knowledge/config"
)

const (
	memgraphConnectAttempts = 5
	memgraphConnectDelay    = 2 * time.Second
)

func connectMemgraph(ctx context.Context, cfg *config.Config) (neo4j.DriverWithContext, error) {
	uri := fmt.Sprintf("bolt://%s:%s", cfg.Memgraph.Host, cfg.Memgraph.Port)
	auth := neo4j.BasicAuth(cfg.Memgraph.User, cfg.Memgraph.Password, "")

	var lastErr error
	for attempt := 1; attempt <= memgraphConnectAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		driver, err := neo4j.NewDriverWithContext(uri, auth)
		if err == nil {
			err = driver.VerifyConnectivity(ctx)
			if err == nil {
				return driver, nil
			}
			if closeErr := driver.Close(ctx); closeErr != nil {
				logrus.WithError(closeErr).Warn("failed to close memgraph driver after connect error")
			}
		}

		lastErr = err
		logrus.WithError(err).Warnf("memgraph connection attempt %d/%d failed", attempt, memgraphConnectAttempts)
		if attempt == memgraphConnectAttempts {
			break
		}
		if err := waitMemgraphRetry(ctx); err != nil {
			return nil, err
		}
	}
	return nil, fmt.Errorf("memgraph: failed to connect after %d attempts: %w", memgraphConnectAttempts, lastErr)
}

func waitMemgraphRetry(ctx context.Context) error {
	timer := time.NewTimer(memgraphConnectDelay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func startupMemgraph(ctx context.Context, memgraphDriver neo4j.DriverWithContext) error {
	session := memgraphDriver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	cyphers := []string{
		`CREATE CONSTRAINT ON (d:Document) ASSERT d.id IS UNIQUE`,
		`CREATE CONSTRAINT ON (e:Entity) ASSERT e.value IS UNIQUE`,
		`CREATE CONSTRAINT ON (t:Tag) ASSERT t.value IS UNIQUE`,
	}
	for _, cypher := range cyphers {
		result, err := session.Run(ctx, cypher, nil)
		if err != nil {
			return err
		}
		if _, err := result.Consume(ctx); err != nil {
			return err
		}
	}
	return nil
}
