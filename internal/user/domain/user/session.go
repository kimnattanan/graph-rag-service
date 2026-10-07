package user

import (
	"time"

	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
)

type Session struct {
	id        string
	userID    string
	expiresAt time.Time
	createdAt time.Time
}

func NewSession(id string, userID string, ttl time.Duration) (Session, error) {
	if id == "" || userID == "" || ttl <= 0 {
		return Session{}, commonerrors.NewIncorrectInputError("invalid session", "invalid-session")
	}
	now := time.Now().UTC()
	return Session{
		id:        id,
		userID:    userID,
		expiresAt: now.Add(ttl),
		createdAt: now,
	}, nil
}

// UnmarshalSessionFromDatabase unmarshals a session from the database.
func UnmarshalSessionFromDatabase(id string, userID string, expiresAt time.Time, createdAt time.Time) (Session, error) {
	if id == "" || userID == "" {
		return Session{}, commonerrors.NewIncorrectInputError("invalid session", "invalid-session")
	}
	return Session{
		id:        id,
		userID:    userID,
		expiresAt: expiresAt,
		createdAt: createdAt,
	}, nil
}

func (s Session) ID() string {
	return s.id
}

func (s Session) UserID() string {
	return s.userID
}

func (s Session) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s Session) CreatedAt() time.Time {
	return s.createdAt
}

func (s Session) Expired(now time.Time) bool {
	return !now.Before(s.expiresAt)
}
