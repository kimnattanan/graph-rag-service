package document

import (
	"errors"
	"time"

	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
)

var (
	ErrCannotUpdateContent  = commonerrors.NewIncorrectInputError("cannot update content if not pending", "cannot-update-content")
	ErrEmptyDocumentID      = commonerrors.NewIncorrectInputError("empty document id", "empty-document-id")
	ErrEmptyDocumentTitle   = commonerrors.NewIncorrectInputError("empty document title", "empty-document-title")
	ErrEmptyDocumentContent = commonerrors.NewIncorrectInputError("empty document content", "empty-document-content")
)

type Document struct {
	id         string
	title      string
	content    string
	tags       []string
	createdAt  time.Time
	updatedAt  time.Time
	indexState IndexState
}

func NewDocument(id string, title string, content string, tags []string) (*Document, error) {
	if id == "" {
		return nil, ErrEmptyDocumentID
	}
	if title == "" {
		return nil, ErrEmptyDocumentTitle
	}
	if content == "" {
		return nil, ErrEmptyDocumentContent
	}
	now := time.Now()
	return &Document{
		id:         id,
		title:      title,
		content:    content,
		tags:       append([]string(nil), tags...),
		createdAt:  now,
		updatedAt:  now,
		indexState: NewIndexState(),
	}, nil
}

// UnmarshalDocumentFromDatabase unmarshals Training from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use UnmarshalDocumentFromDatabase as constructor - It may put domain into the invalid state!
func UnmarshalDocumentFromDatabase(id string, title string, content string, tags []string, createdAt time.Time, updatedAt time.Time, indexState IndexState) (*Document, error) {
	doc, err := NewDocument(id, title, content, tags)
	if err != nil {
		return nil, err
	}
	doc.createdAt = createdAt
	doc.updatedAt = updatedAt
	doc.indexState = indexState
	return doc, nil
}

func (d *Document) ID() string {
	return d.id
}

func (d *Document) Title() string {
	return d.title
}

func (d *Document) Content() string {
	return d.content
}

func (d *Document) Tags() []string {
	return append([]string(nil), d.tags...)
}

func (d *Document) CreatedAt() time.Time {
	return d.createdAt
}

func (d *Document) UpdatedAt() time.Time {
	return d.updatedAt
}

func (d *Document) IndexState() IndexState {
	return d.indexState
}

func (d *Document) MarkAsPending() error {
	if d.indexState.IsPending() {
		return nil
	}
	if err := d.indexState.MarkAsPending(); err != nil {
		return err
	}
	d.updatedAt = time.Now()
	return nil
}

func (d *Document) MarkAsIndexing() error {
	if d.indexState.IsIndexing() {
		return nil
	}
	if err := d.indexState.MarkAsIndexing(); err != nil {
		return err
	}
	d.updatedAt = time.Now()
	return nil
}

func (d *Document) MarkAsCompleted() error {
	if d.indexState.IsCompleted() {
		return nil
	}
	if err := d.indexState.MarkAsCompleted(); err != nil {
		return err
	}
	d.updatedAt = time.Now()
	return nil
}

func (d *Document) MarkAsFailed(errorMessage string) error {
	if d.indexState.IsFailed() {
		return nil
	}
	if err := d.indexState.MarkAsFailed(errorMessage); err != nil {
		return err
	}
	d.updatedAt = time.Now()
	return nil
}

func (d *Document) UpdateTitle(title string) error {
	if title == "" {
		return errors.New("empty document title")
	}
	d.title = title
	d.updatedAt = time.Now()
	return nil
}

func (d *Document) UpdateContent(content string) error {
	if !d.indexState.IsPending() {
		return ErrCannotUpdateContent
	}
	if content == "" {
		return errors.New("empty document content")
	}
	d.content = content
	d.updatedAt = time.Now()
	return nil
}

func (d *Document) UpdateTags(tags []string) {
	d.tags = append([]string(nil), tags...)
	d.updatedAt = time.Now()
}
