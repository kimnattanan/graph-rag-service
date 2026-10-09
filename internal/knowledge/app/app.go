package app

import (
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/knowledge/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateDocument    command.CreateDocumentHandler
	UpdateDocument    command.UpdateDocumentHandler
	DeleteDocument    command.DeleteDocumentHandler
	ReindexDocument   command.ReindexDocumentHandler
	IndexNextDocument command.IndexNextDocumentHandler
	EmbedNextEntity   command.EmbedNextEntityHandler
	SweepOrphans      command.SweepOrphansHandler
}

type Queries struct {
	ListDocuments          query.ListDocumentsHandler
	GetDocument            query.GetDocumentHandler
	GetDocumentIndexStatus query.GetDocumentIndexStatusHandler
	Retrieve               query.RetrieveHandler
}
