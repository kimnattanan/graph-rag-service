package adapters

import (
	"context"
	"errors"
	"fmt"
	"time"

	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/query"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/domain/document"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type DocumentMemgraphReadModel struct {
	memgraphDriver neo4j.DriverWithContext
}

func NewDocumentMemgraphReadModel(memgraphDriver neo4j.DriverWithContext) *DocumentMemgraphReadModel {
	return &DocumentMemgraphReadModel{
		memgraphDriver: memgraphDriver,
	}
}

func (r *DocumentMemgraphReadModel) ListDocuments(ctx context.Context, limit int, offset int, tag *string, indexStatus *document.IndexStatus) (*query.DocumentSummaryList, error) {
	params := map[string]any{
		"limit":  limit,
		"offset": offset,
	}
	indexStatusClause := ""
	if indexStatus != nil {
		params["indexStatus"] = int64(*indexStatus)
		indexStatusClause = "WHERE d.index_status = $indexStatus"
	}

	var cypher string
	if tag == nil {
		cypher = fmt.Sprintf(`
			MATCH (d:Document)
			%s
			WITH d
			ORDER BY d.created_at DESC, d.id ASC
			WITH collect(d) AS docs
			WITH size(docs) AS total, docs[$offset..$offset + $limit] AS page
			UNWIND CASE WHEN size(page) = 0 THEN [null] ELSE page END AS d
			OPTIONAL MATCH (d)-[:TAGGED_AS]->(t:Tag)
			RETURN total, d, collect(t.value) AS tags
		`, indexStatusClause)
	} else {
		params["tag"] = *tag
		cypher = fmt.Sprintf(`
			MATCH (t:Tag {value: $tag})
			MATCH (d:Document)-[:TAGGED_AS]->(t)
			%s
			WITH d
			ORDER BY d.created_at DESC, d.id ASC
			WITH collect(d) AS docs
			WITH size(docs) AS total, docs[$offset..$offset + $limit] AS page
			UNWIND CASE WHEN size(page) = 0 THEN [null] ELSE page END AS d
			OPTIONAL MATCH (d)-[:TAGGED_AS]->(t:Tag)
			RETURN total, d, collect(t.value) AS tags
		`, indexStatusClause)
	}

	result, err := neo4j.ExecuteQuery(ctx, r.memgraphDriver, cypher, params, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}

	docs := make([]query.DocumentSummary, 0, len(result.Records))
	total := 0
	if len(result.Records) > 0 {
		rawTotal, _, err := neo4j.GetRecordValue[int64](result.Records[0], "total")
		if err != nil {
			return nil, err
		}
		total = int(rawTotal)
	}
	for _, record := range result.Records {
		node, isNil, err := neo4j.GetRecordValue[neo4j.Node](record, "d")
		if err != nil {
			return nil, err
		}
		if isNil {
			continue
		}
		rawTags, isNil, err := neo4j.GetRecordValue[[]any](record, "tags")
		if err != nil {
			return nil, err
		}
		tags := []string{}
		if !isNil {
			tags, err = tagValues(rawTags)
			if err != nil {
				return nil, err
			}
		}
		doc, err := documentSummaryFromNode(&node, tags)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return &query.DocumentSummaryList{
		Items: docs,
		Total: total,
	}, nil
}

func (r *DocumentMemgraphReadModel) GetDocument(ctx context.Context, documentID string) (*document.Document, error) {
	result, err := neo4j.ExecuteQuery(ctx, r.memgraphDriver, matchDocumentWithTagsCypher, map[string]any{
		"id": documentID,
	}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return nil, err
	}
	if len(result.Records) == 0 {
		return nil, commonerrors.NewNotFoundError("document not found", "document-not-found")
	}
	return domainDocumentFromRecord(result.Records[0])
}

func (r *DocumentMemgraphReadModel) Retrieve(ctx context.Context, embedding []float64, topK int, tags *[]string) (query.RetrieveResult, error) {
	tagFilter := []string{}
	if tags != nil {
		tagFilter = *tags
	}

	cypher := `
		CALL vector_search.search("chunk_embedding", $candidates, $embedding)
		YIELD node AS c, similarity
		MATCH (d:Document)-[:HAS_CHUNK]->(c)
		WHERE d.index_status = $completedStatus
		OPTIONAL MATCH (d)-[:TAGGED_AS]->(t:Tag)
		WITH d, c, similarity, collect(DISTINCT t.value) AS docTags
		WHERE size($tags) = 0 OR any(tag IN $tags WHERE tag IN docTags)
		ORDER BY similarity DESC
		LIMIT $topK
		OPTIONAL MATCH (c)-[:MENTIONS]->(e:Entity)
		RETURN d.id AS documentID,
			d.title AS documentTitle,
			c.id AS chunkID,
			c.text AS text,
			similarity AS score,
			collect(DISTINCT e.name) AS graphPath
	`
	result, err := neo4j.ExecuteQuery(ctx, r.memgraphDriver, cypher, map[string]any{
		"embedding":       embedding,
		"topK":            topK,
		"candidates":      min(topK*10, 100),
		"completedStatus": int64(document.IndexStatusCompleted),
		"tags":            tagFilter,
	}, neo4j.EagerResultTransformer, neo4j.ExecuteQueryWithDatabase(""))
	if err != nil {
		return query.RetrieveResult{}, err
	}

	chunks := make([]query.RetrievedChunk, 0, len(result.Records))
	for _, record := range result.Records {
		chunk, err := retrievedChunkFromRecord(record)
		if err != nil {
			return query.RetrieveResult{}, err
		}
		chunks = append(chunks, chunk)
	}
	return query.RetrieveResult{Chunks: chunks}, nil
}

const matchDocumentWithTagsCypher = `
	MATCH (d:Document {id: $id})
	OPTIONAL MATCH (d)-[:TAGGED_AS]->(t:Tag)
	RETURN d, collect(t.value) AS tags
`

func domainDocumentFromRecord(record *neo4j.Record) (*document.Document, error) {
	node, _, err := neo4j.GetRecordValue[neo4j.Node](record, "d")
	if err != nil {
		return nil, err
	}
	rawTags, isNil, err := neo4j.GetRecordValue[[]any](record, "tags")
	if err != nil {
		return nil, err
	}
	tags := []string{}
	if !isNil {
		tags, err = tagValues(rawTags)
		if err != nil {
			return nil, err
		}
	}
	return domainDocumentFromNode(&node, tags)
}

func domainDocumentFromNode(node *neo4j.Node, tags []string) (*document.Document, error) {
	summary, err := documentSummaryFromNode(node, tags)
	if err != nil {
		return nil, err
	}
	content, ok := node.Props["content"].(string)
	if !ok {
		return nil, errors.New("content is not string")
	}
	startedAt, err := optionalProp[time.Time](node.Props, "index_started_at")
	if err != nil {
		return nil, err
	}
	finishedAt, err := optionalProp[time.Time](node.Props, "index_finished_at")
	if err != nil {
		return nil, err
	}
	errorMessage, err := optionalProp[string](node.Props, "index_error_message")
	if err != nil {
		return nil, err
	}
	return document.UnmarshalDocumentFromDatabase(
		summary.ID,
		summary.Title,
		content,
		summary.Tags,
		summary.CreatedAt,
		summary.UpdatedAt,
		document.UnmarshalIndexStateFromDatabase(
			summary.IndexStatus,
			startedAt,
			finishedAt,
			errorMessage,
		),
	)
}

func documentSummaryFromNode(node *neo4j.Node, tags []string) (query.DocumentSummary, error) {
	props := node.Props
	id, ok := props["id"].(string)
	if !ok {
		return query.DocumentSummary{}, errors.New("id is not string")
	}
	title, ok := props["title"].(string)
	if !ok {
		return query.DocumentSummary{}, errors.New("title is not string")
	}
	if tags == nil {
		tags = []string{}
	}
	indexStatus, ok := props["index_status"].(int64)
	if !ok {
		return query.DocumentSummary{}, errors.New("index_status is not int64")
	}
	createdAt, ok := props["created_at"].(time.Time)
	if !ok {
		return query.DocumentSummary{}, errors.New("created_at is not time.Time")
	}
	updatedAt, ok := props["updated_at"].(time.Time)
	if !ok {
		return query.DocumentSummary{}, errors.New("updated_at is not time.Time")
	}
	return query.DocumentSummary{
		ID:          id,
		Title:       title,
		Tags:        tags,
		IndexStatus: document.IndexStatus(indexStatus),
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
	}, nil
}

func retrievedChunkFromRecord(record *neo4j.Record) (query.RetrievedChunk, error) {
	documentID, isNil, err := neo4j.GetRecordValue[string](record, "documentID")
	if err != nil || isNil {
		return query.RetrievedChunk{}, errors.New("documentID is not string")
	}
	documentTitle, isNil, err := neo4j.GetRecordValue[string](record, "documentTitle")
	if err != nil || isNil {
		return query.RetrievedChunk{}, errors.New("documentTitle is not string")
	}
	chunkID, isNil, err := neo4j.GetRecordValue[string](record, "chunkID")
	if err != nil || isNil {
		return query.RetrievedChunk{}, errors.New("chunkID is not string")
	}
	text, isNil, err := neo4j.GetRecordValue[string](record, "text")
	if err != nil || isNil {
		return query.RetrievedChunk{}, errors.New("text is not string")
	}
	score, isNil, err := neo4j.GetRecordValue[float64](record, "score")
	if err != nil || isNil {
		return query.RetrievedChunk{}, errors.New("score is not float64")
	}
	rawPath, pathIsNil, err := neo4j.GetRecordValue[[]any](record, "graphPath")
	if err != nil {
		return query.RetrievedChunk{}, err
	}
	var graphPath *[]string
	if !pathIsNil {
		values, err := stringValues(rawPath)
		if err != nil {
			return query.RetrievedChunk{}, err
		}
		if len(values) > 0 {
			graphPath = &values
		}
	}
	return query.RetrievedChunk{
		DocumentID:    documentID,
		DocumentTitle: documentTitle,
		ChunkID:       &chunkID,
		Text:          text,
		Score:         score,
		GraphPath:     graphPath,
	}, nil
}

func optionalProp[T any](props map[string]any, key string) (T, error) {
	raw, ok := props[key]
	if !ok || raw == nil {
		var zero T
		return zero, nil
	}
	value, ok := raw.(T)
	if !ok {
		var zero T
		return zero, fmt.Errorf("%s is not %T", key, zero)
	}
	return value, nil
}

func stringValues(raw []any) ([]string, error) {
	values := make([]string, 0, len(raw))
	for _, rawValue := range raw {
		if rawValue == nil {
			continue
		}
		value, ok := rawValue.(string)
		if !ok {
			return nil, errors.New("value is not string")
		}
		values = append(values, value)
	}
	return values, nil
}

func tagValues(rawTags []any) ([]string, error) {
	tags := make([]string, 0, len(rawTags))
	for _, rawTag := range rawTags {
		if rawTag == nil {
			continue
		}
		tag, ok := rawTag.(string)
		if !ok {
			return nil, errors.New("tag is not string")
		}
		tags = append(tags, tag)
	}
	return tags, nil
}
