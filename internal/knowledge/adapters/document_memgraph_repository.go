package adapters

import (
	"context"
	"time"

	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var _ document.Repository = (*DocumentMemgraphRepository)(nil)

type DocumentMemgraphRepository struct {
	memgraphDriver neo4j.DriverWithContext
}

func NewDocumentMemgraphRepository(memgraphDriver neo4j.DriverWithContext) *DocumentMemgraphRepository {
	return &DocumentMemgraphRepository{
		memgraphDriver: memgraphDriver,
	}
}

func (r *DocumentMemgraphRepository) AddDocument(ctx context.Context, doc *document.Document) error {
	return r.executeWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Create only when the id is free. The unique constraint on Document.id also rejects a concurrent insert that passes this check.
		result, err := tx.Run(ctx, `
			OPTIONAL MATCH (existing:Document {id: $id})
			WITH existing
			WHERE existing IS NULL
			CREATE (d:Document {
				id: $id,
				title: $title,
				content: $content,
				index_status: $indexStatus,
				created_at: $createdAt,
				updated_at: $updatedAt
			})
			SET d.index_started_at = $indexStartedAt,
			    d.index_finished_at = $indexFinishedAt,
			    d.index_error_message = $indexErrorMessage
			FOREACH (tag IN $tags |
				MERGE (t:Tag {value: tag})
				MERGE (d)-[:TAGGED_AS]->(t)
			)
			RETURN d.id AS id
		`, documentParams(doc))
		if err != nil {
			return nil, err
		}
		records, err := result.Collect(ctx)
		if err != nil {
			return nil, err
		}
		if len(records) == 0 {
			return nil, commonerrors.NewIncorrectInputError("document already exists", "document-already-exists")
		}
		return nil, nil
	})
}

func (r *DocumentMemgraphRepository) UpdateDocument(ctx context.Context, docID string, updateFn func(doc *document.Document) error) error {
	return r.executeWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		doc, err := documentByID(ctx, tx, docID)
		if err != nil {
			return nil, err
		}
		indexState := doc.IndexState()
		wasPending := indexState.IsPending()
		contentBefore := doc.Content()
		if err := updateFn(doc); err != nil {
			return nil, err
		}
		if err := saveDocument(ctx, tx, doc); err != nil {
			return nil, err
		}
		// delete chunks if reindexing is needed
		indexState = doc.IndexState()
		if doc.Content() != contentBefore || (indexState.IsPending() && !wasPending) {
			if err := deleteDocumentChunks(ctx, tx, doc.ID()); err != nil {
				return nil, err
			}
		}
		return nil, nil
	})
}

func (r *DocumentMemgraphRepository) DeleteDocument(ctx context.Context, docID string) error {
	return r.executeWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		result, err := tx.Run(ctx, `
			MATCH (d:Document {id: $id})
			OPTIONAL MATCH (d)-[:HAS_CHUNK]->(c:Chunk)
			WITH d, collect(c) AS chunks
			FOREACH (c IN chunks | DETACH DELETE c)
			DETACH DELETE d
			RETURN true AS deleted
		`, map[string]any{"id": docID})
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

func (r *DocumentMemgraphRepository) ClaimNextPending(ctx context.Context) (*document.Document, error) {
	var claimed *document.Document
	err := r.executeWrite(ctx, func(tx neo4j.ManagedTransaction) (any, error) {
		// Oldest pending document first.
		// The status flip happens in this write, so a concurrent claim conflicts instead of taking the same document.
		result, err := tx.Run(ctx, `
			MATCH (d:Document)
			WHERE d.index_status = $pending
			WITH d
			ORDER BY d.created_at ASC, d.id ASC
			LIMIT 1
			OPTIONAL MATCH (d)-[:TAGGED_AS]->(t:Tag)
			RETURN d, collect(t.value) AS tags
		`, map[string]any{
			"pending": int64(document.IndexStatusPending),
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
		doc, err := domainDocumentFromRecord(records[0])
		if err != nil {
			return nil, err
		}
		if err := doc.MarkAsIndexing(); err != nil {
			return nil, err
		}
		if err := saveDocument(ctx, tx, doc); err != nil {
			return nil, err
		}
		claimed = doc
		return nil, nil
	})
	if err != nil {
		return nil, err
	}
	return claimed, nil
}

func (r *DocumentMemgraphRepository) executeWrite(ctx context.Context, work neo4j.ManagedTransactionWork) error {
	session := r.memgraphDriver.NewSession(ctx, neo4j.SessionConfig{})
	defer session.Close(ctx)
	_, err := session.ExecuteWrite(ctx, work)
	return err
}

func documentByID(ctx context.Context, tx neo4j.ManagedTransaction, documentID string) (*document.Document, error) {
	result, err := tx.Run(ctx, matchDocumentWithTagsCypher, map[string]any{
		"id": documentID,
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
	return domainDocumentFromRecord(records[0])
}

func saveDocument(ctx context.Context, tx neo4j.ManagedTransaction, doc *document.Document) error {
	result, err := tx.Run(ctx, `
		MATCH (d:Document {id: $id})
		SET d.title = $title,
		    d.content = $content,
		    d.index_status = $indexStatus,
		    d.created_at = $createdAt,
		    d.updated_at = $updatedAt,
		    d.index_started_at = $indexStartedAt,
		    d.index_finished_at = $indexFinishedAt,
		    d.index_error_message = $indexErrorMessage
		WITH d
		OPTIONAL MATCH (d)-[oldTag:TAGGED_AS]->(:Tag)
		WITH d, collect(oldTag) AS oldTags
		FOREACH (oldTag IN oldTags | DELETE oldTag)
		FOREACH (tag IN $tags |
			MERGE (t:Tag {value: tag})
			MERGE (d)-[:TAGGED_AS]->(t)
		)
		RETURN d.id AS id
	`, documentParams(doc))
	if err != nil {
		return err
	}
	records, err := result.Collect(ctx)
	if err != nil {
		return err
	}
	if len(records) == 0 {
		return commonerrors.NewNotFoundError("document not found", "document-not-found")
	}
	return nil
}

func deleteDocumentChunks(ctx context.Context, tx neo4j.ManagedTransaction, documentID string) error {
	result, err := tx.Run(ctx, `
		MATCH (d:Document {id: $id})
		OPTIONAL MATCH (d)-[:HAS_CHUNK]->(c:Chunk)
		WITH collect(c) AS chunks
		FOREACH (c IN chunks | DETACH DELETE c)
	`, map[string]any{"id": documentID})
	if err != nil {
		return err
	}
	_, err = result.Consume(ctx)
	return err
}

func documentParams(doc *document.Document) map[string]any {
	tags := doc.Tags()
	if tags == nil {
		tags = []string{}
	}
	indexState := doc.IndexState()
	return map[string]any{
		"id":                doc.ID(),
		"title":             doc.Title(),
		"content":           doc.Content(),
		"tags":              tags,
		"indexStatus":       int64(indexState.Status()),
		"createdAt":         doc.CreatedAt(),
		"updatedAt":         doc.UpdatedAt(),
		"indexStartedAt":    optionalTime(indexState.StartedAt()),
		"indexFinishedAt":   optionalTime(indexState.FinishedAt()),
		"indexErrorMessage": optionalString(indexState.ErrorMessage()),
	}
}

func optionalTime(value time.Time) any {
	if value.IsZero() {
		return nil
	}
	return value
}

func optionalString(value string) any {
	if value == "" {
		return nil
	}
	return value
}
