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

func (r *IndexMemgraphRepository) ReplaceChunks(ctx context.Context, documentID string, chunks []indexing.ExtractResult, chunkEmbeddings [][]float64) error {
	chunkParams, err := buildChunkParams(chunks, chunkEmbeddings)
	if err != nil {
		return err
	}
	return r.executeWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (d:Document {id: $id})
			OPTIONAL MATCH (d)-[:HAS_CHUNK]->(c:Chunk)
			WITH d, collect(c) AS oldChunks
			FOREACH (c IN oldChunks | DETACH DELETE c)
			FOREACH (chunk IN $chunks |
				CREATE (newChunk:Chunk {id: chunk.id, text: chunk.text, embedding: chunk.embedding})
				MERGE (d)-[:HAS_CHUNK]->(newChunk)
				FOREACH (name IN chunk.entities |
					MERGE (e:Entity {name: name})
					ON CREATE SET e.status = $pendingStatus
					MERGE (newChunk)-[:MENTIONS]->(e)
				)
			)
			RETURN d.id AS id
		`, map[string]any{
			"id":            documentID,
			"chunks":        chunkParams,
			"pendingStatus": int64(indexing.EntityStatusPending),
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

func (r *IndexMemgraphRepository) ClaimNextPendingEntity(ctx context.Context) (*indexing.Entity, error) {
	var claimed *indexing.Entity
	err := r.executeWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// get one pending entity and mark it as embedding
		result, err := tx.Run(ctx, `
			MATCH (e:Entity)
			WHERE e.status = $pending OR e.status IS NULL
			WITH e
			ORDER BY e.name ASC
			LIMIT 1
			SET e.status = $embedding
			RETURN e.name AS name
		`, map[string]any{
			"pending":   int64(indexing.EntityStatusPending),
			"embedding": int64(indexing.EntityStatusEmbedding),
		})
		if err != nil {
			return nil, err
		}
		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}
		if len(records) == 0 {
			return nil, nil
		}
		name, isNil, err := neo4j.GetRecordValue[string](records[0], "name")
		if err != nil {
			return nil, err
		}
		if isNil {
			return nil, commonerrors.NewSlugError("entity name is empty", "entity-name-empty")
		}
		claimed = indexing.UnmarshalEntityFromDatabase(name, indexing.EntityStatusEmbedding)
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return claimed, nil
}

func (r *IndexMemgraphRepository) EmbedEntity(ctx context.Context, name string, embedding []float64) error {
	return r.executeWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (e:Entity {name: $name})
			SET e.embedding = $embedding,
				e.status = $embeddedStatus
			RETURN e.name AS name
		`, map[string]any{
			"name":           name,
			"embedding":      embedding,
			"embeddedStatus": int64(indexing.EntityStatusEmbedded),
		})
		if err != nil {
			r.setEntityStatus(ctx, name, indexing.EntityStatusPending)
			return nil, err
		}
		records, err := result.Collect(ctx)
		if err != nil {
			r.setEntityStatus(ctx, name, indexing.EntityStatusPending)
			return nil, err
		}
		if len(records) == 0 {
			return nil, commonerrors.NewNotFoundError("entity not found", "entity-not-found")
		}
		return nil, nil
	})
}

func (r *IndexMemgraphRepository) setEntityStatus(ctx context.Context, name string, status indexing.EntityStatus) error {
	return r.executeWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (e:Entity {name: $name})
			SET e.status = $status
			RETURN e.name AS name
		`, map[string]any{
			"name":   name,
			"status": int64(status),
		})
		if err != nil {
			return nil, err
		}
		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}
		if len(records) == 0 {
			return nil, commonerrors.NewNotFoundError("entity not found", "entity-not-found")
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

func buildChunkParams(chunks []indexing.ExtractResult, embeddings [][]float64) ([]map[string]any, error) {
	if len(chunks) != len(embeddings) {
		return nil, commonerrors.NewSlugError("chunkParams: chunks and embeddings must have the same length", "chunk-params-mismatch")
	}
	params := make([]map[string]any, 0, len(chunks))
	for i, chunk := range chunks {
		entities := chunk.Entities()
		if entities == nil {
			entities = []string{}
		}
		params = append(params, map[string]any{
			"id":        uuid.NewString(),
			"text":      chunk.Content(),
			"entities":  entities,
			"embedding": embeddings[i],
		})
	}
	return params, nil
}
