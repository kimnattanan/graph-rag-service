package document

import (
	"time"

	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
)

type IndexStatus int

const (
	IndexStatusPending IndexStatus = iota
	IndexStatusIndexing
	IndexStatusCompleted
	IndexStatusFailed
)

var (
	ErrCannotMarkAsPending   = commonerrors.NewIncorrectInputError("cannot mark as pending if indexing", "cannot-mark-as-pending")
	ErrCannotMarkAsIndexing  = commonerrors.NewIncorrectInputError("cannot mark as indexing if not pending", "cannot-mark-as-indexing")
	ErrCannotMarkAsCompleted = commonerrors.NewIncorrectInputError("cannot mark as completed if not indexing", "cannot-mark-as-completed")
	ErrCannotMarkAsFailed    = commonerrors.NewIncorrectInputError("cannot mark as failed if not indexing", "cannot-mark-as-failed")
)

type IndexState struct {
	status       IndexStatus
	startedAt    time.Time
	finishedAt   time.Time
	errorMessage string
}

func NewIndexState() IndexState {
	return IndexState{
		status:       IndexStatusPending,
		startedAt:    time.Time{},
		finishedAt:   time.Time{},
		errorMessage: "",
	}
}

// UnmarshalIndexStateFromDatabase unmarshals Training from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use UnmarshalIndexStateFromDatabase as constructor - It may put domain into the invalid state!
func UnmarshalIndexStateFromDatabase(status IndexStatus, startedAt time.Time, finishedAt time.Time, errorMessage string) IndexState {
	return IndexState{
		status:       status,
		startedAt:    startedAt,
		finishedAt:   finishedAt,
		errorMessage: errorMessage,
	}
}

func (i *IndexState) Status() IndexStatus {
	return i.status
}

func (i *IndexState) IsPending() bool {
	return i.status == IndexStatusPending
}

func (i *IndexState) IsIndexing() bool {
	return i.status == IndexStatusIndexing
}

func (i *IndexState) IsCompleted() bool {
	return i.status == IndexStatusCompleted
}

func (i *IndexState) IsFailed() bool {
	return i.status == IndexStatusFailed
}

func (i *IndexState) StartedAt() time.Time {
	return i.startedAt
}

func (i *IndexState) FinishedAt() time.Time {
	return i.finishedAt
}

func (i *IndexState) ErrorMessage() string {
	return i.errorMessage
}

func (i *IndexState) MarkAsPending() error {
	if i.IsPending() {
		return nil
	}
	if i.IsIndexing() {
		return ErrCannotMarkAsPending
	}
	i.status = IndexStatusPending
	i.startedAt = time.Time{}
	i.finishedAt = time.Time{}
	i.errorMessage = ""
	return nil
}

func (i *IndexState) MarkAsIndexing() error {
	if i.IsIndexing() {
		return nil
	}
	if !i.IsPending() {
		return ErrCannotMarkAsIndexing
	}
	i.status = IndexStatusIndexing
	i.startedAt = time.Now()
	i.finishedAt = time.Time{}
	i.errorMessage = ""
	return nil
}

func (i *IndexState) MarkAsCompleted() error {
	if i.IsCompleted() {
		return nil
	}
	if !i.IsIndexing() {
		return ErrCannotMarkAsCompleted
	}
	i.status = IndexStatusCompleted
	i.finishedAt = time.Now()
	i.errorMessage = ""
	return nil
}

func (i *IndexState) MarkAsFailed(errorMessage string) error {
	if i.IsFailed() {
		return nil
	}
	if !i.IsIndexing() {
		return ErrCannotMarkAsFailed
	}
	i.status = IndexStatusFailed
	i.finishedAt = time.Now()
	i.errorMessage = errorMessage
	return nil
}
