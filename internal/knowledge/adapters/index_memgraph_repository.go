package adapters

import (
	"context"

	"github.com/google/uuid"
	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/indexing"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var _ indexing.Repository = (*IndexMemgraphRepository)(nil)

type IndexMemgraphRepository struct {
	memgraphDriver neo4j.DriverWithContext
}

func NewIndexMemgraphRepository(memgraphDriver neo4j.DriverWithContext) *IndexMemgraphRepository {
	return &IndexMemgraphRepository{
		memgraphDriver: memgraphDriver,
	}
}

func (r *IndexMemgraphRepository) ReplaceChunks(ctx context.Context, documentID string, chunks []indexing.ExtractResult) error {
	return r.executeWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (d:Document {id: $id})
			OPTIONAL MATCH (d)-[:HAS_CHUNK]->(c:Chunk)
			WITH d, collect(c) AS oldChunks
			FOREACH (c IN oldChunks | DETACH DELETE c)
			FOREACH (chunk IN $chunks |
				CREATE (newChunk:Chunk {id: chunk.id, text: chunk.text})
				MERGE (d)-[:HAS_CHUNK]->(newChunk)
				FOREACH (name IN chunk.entities |
					MERGE (e:Entity {value: name})
					MERGE (newChunk)-[:MENTIONS]->(e)
				)
			)
			RETURN d.id AS id
		`, map[string]any{
			"id":     documentID,
			"chunks": chunkParams(chunks),
		})
		if err != nil {
			return nil, err
		}
		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}
		if len(records) == 0 {
			return nil, commonerrors.NewNotFoundError("document not found", "document-not-found")
		}
		return nil, nil
	})
}

func (r *IndexMemgraphRepository) DeleteOrphans(ctx context.Context) error {
	// order matters
	// 1. delete chunks with no HAS_CHUNK relationship
	// 2. delete entities with no MENTIONS relationship
	// 3. delete tags with no TAGGED_AS relationship
	return r.executeWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		for _, cypher := range []string{
			`
				MATCH (c:Chunk)
				WHERE NOT (()-[:HAS_CHUNK]->(c))
				DETACH DELETE c
			`,
			`
				MATCH (e:Entity)
				WHERE NOT (()-[:MENTIONS]->(e))
				DETACH DELETE e
			`,
			`
				MATCH (t:Tag)
				WHERE NOT (()-[:TAGGED_AS]->(t))
				DETACH DELETE t
			`,
		} {
			result, err := tx.Run(ctx, cypher, nil)
			if err != nil {
				return nil, err
			}
			if _, err := result.Consume(ctx); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
}

func (r *IndexMemgraphRepository) executeWrite(ctx context.Context, work neo4j.ManagedTransactionWork) error {
	session := r.memgraphDriver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, work)
	return err
}

func chunkParams(chunks []indexing.ExtractResult) []map[string]any {
	params := make([]map[string]any, 0, len(chunks))
	for _, chunk := range chunks {
		entities := chunk.Entities()
		if entities == nil {
			entities = []string{}
		}
		params = append(params, map[string]any{
			"id":       uuid.NewString(),
			"text":     chunk.Content,
			"entities": entities,
		})
	}
	return params
}
