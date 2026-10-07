package document

import "context"

type Repository interface {
	AddDocument(ctx context.Context, doc *Document) error
	UpdateDocument(ctx context.Context, docID string, updateFn func(doc *Document) error) error
	DeleteDocument(ctx context.Context, docID string) error
	ClaimNextPending(ctx context.Context) (*Document, error)
}
