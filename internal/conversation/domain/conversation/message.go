package conversation

import (
	"time"

	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
)

var (
	ErrEmptyMessageID      = commonerrors.NewIncorrectInputError("empty message id", "empty-message-id")
	ErrEmptyMessageContent = commonerrors.NewIncorrectInputError("empty message content", "empty-message-content")
	ErrInvalidMessageRole  = commonerrors.NewIncorrectInputError("invalid message role", "invalid-message-role")
)

type Role string

const (
	RoleUser      Role = "User"
	RoleAssistant Role = "Assistant"
)

func (r Role) IsValid() bool {
	return r == RoleUser || r == RoleAssistant
}

type Source struct {
	documentID    string
	documentTitle string
	chunkID       string
	text          string
	score         float64
}

func NewSource(documentID string, documentTitle string, chunkID string, text string, score float64) Source {
	return Source{
		documentID:    documentID,
		documentTitle: documentTitle,
		chunkID:       chunkID,
		text:          text,
		score:         score,
	}
}

func (s Source) DocumentID() string {
	return s.documentID
}

func (s Source) DocumentTitle() string {
	return s.documentTitle
}

func (s Source) ChunkID() string {
	return s.chunkID
}

func (s Source) Text() string {
	return s.text
}

func (s Source) Score() float64 {
	return s.score
}

type Message struct {
	id             string
	conversationID string
	role           Role
	content        string
	sources        []Source
	createdAt      time.Time
}

func NewMessage(id string, conversationID string, role Role, content string, sources []Source) (*Message, error) {
	if id == "" {
		return nil, ErrEmptyMessageID
	}
	if conversationID == "" {
		return nil, ErrEmptyConversationID
	}
	if !role.IsValid() {
		return nil, ErrInvalidMessageRole
	}
	if content == "" {
		return nil, ErrEmptyMessageContent
	}
	return &Message{
		id:             id,
		conversationID: conversationID,
		role:           role,
		content:        content,
		sources:        append([]Source(nil), sources...),
		createdAt:      time.Now(),
	}, nil
}

// UnmarshalMessageFromDatabase unmarshals Message from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use UnmarshalMessageFromDatabase as constructor - It may put domain into the invalid state!
func UnmarshalMessageFromDatabase(
	id string,
	conversationID string,
	role Role,
	content string,
	sources []Source,
	createdAt time.Time,
) (*Message, error) {
	message, err := NewMessage(id, conversationID, role, content, sources)
	if err != nil {
		return nil, err
	}
	message.createdAt = createdAt
	return message, nil
}

func (m Message) ID() string {
	return m.id
}

func (m Message) ConversationID() string {
	return m.conversationID
}

func (m Message) Role() Role {
	return m.role
}

func (m Message) Content() string {
	return m.content
}

func (m Message) Sources() []Source {
	return append([]Source(nil), m.sources...)
}

func (m Message) CreatedAt() time.Time {
	return m.createdAt
}
