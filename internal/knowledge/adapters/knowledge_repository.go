package adapters

import (
	"context"
	"errors"
	"time"

	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/query"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type KnowledgeRepository struct {
	memgraphDriver neo4j.DriverWithContext
}

func NewKnowledgeRepository(memgraphDriver neo4j.DriverWithContext) *KnowledgeRepository {
	return &KnowledgeRepository{
		memgraphDriver: memgraphDriver,
	}
}

func (r *KnowledgeRepository) ListDocuments(ctx context.Context) ([]*query.DocumentSummary, error) {
	cypher := "MATCH (d:Document) RETURN d;"
	result, err := neo4j.ExecuteQuery(ctx, r.memgraphDriver, cypher, nil, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	docs := make([]*query.DocumentSummary, 0, len(result.Records))
	for _, record := range result.Records {
		node, _, err := neo4j.GetRecordValue[neo4j.Node](record, "d")
		if err != nil {
			return nil, err
		}
		doc, err := r.documentSummaryFromNode(&node)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, nil
}

// Mapper functions

func (r *KnowledgeRepository) documentFromNode(node *neo4j.Node) (*query.Document, error) {
	summary, err := r.documentSummaryFromNode(node)
	if err != nil {
		return nil, err
	}
	content, ok := node.Props["content"].(string)
	if !ok {
		return nil, errors.New("content is not string")
	}
	return &query.Document{
		ID:          summary.ID,
		Title:       summary.Title,
		Content:     content,
		Tags:        summary.Tags,
		IndexStatus: summary.IndexStatus,
		CreatedAt:   summary.CreatedAt,
		UpdatedAt:   summary.UpdatedAt,
	}, nil
}

func (r *KnowledgeRepository) documentSummaryFromNode(node *neo4j.Node) (*query.DocumentSummary, error) {
	props := node.Props
	id, ok := props["id"].(string)
	if !ok {
		return nil, errors.New("id is not string")
	}
	title, ok := props["title"].(string)
	if !ok {
		return nil, errors.New("title is not string")
	}
	tags, ok := props["tags"].([]string)
	if !ok {
		return nil, errors.New("tags is not []string")
	}
	indexStatus, ok := props["index_status"].(int64)
	if !ok {
		return nil, errors.New("index_status is not int64")
	}
	createdAt, ok := props["created_at"].(time.Time)
	if !ok {
		return nil, errors.New("created_at is not time.Time")
	}
	updatedAt, ok := props["updated_at"].(time.Time)
	if !ok {
		return nil, errors.New("updated_at is not time.Time")
	}
	return &query.DocumentSummary{
		ID:          id,
		Title:       title,
		Tags:        tags,
		IndexStatus: query.IndexStatus(indexStatus),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}
