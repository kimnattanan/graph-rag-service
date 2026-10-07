package service

import (
	"context"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

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
