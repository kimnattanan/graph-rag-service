package conversation

import (
	"time"

	commonerrors "github.com/kimnattanan/graph-rag-service/internal/common/errors"
)

var (
	ErrEmptyConversationID         = commonerrors.NewIncorrectInputError("empty conversation id", "empty-conversation-id")
	ErrEmptyUserID                 = commonerrors.NewIncorrectInputError("empty user id", "empty-user-id")
	ErrMessageConversationMismatch = commonerrors.NewIncorrectInputError("message does not belong to this conversation", "message-conversation-mismatch")
	ErrConversationNotOwnedByUser  = commonerrors.NewIncorrectInputError("conversation not owned by user", "conversation-not-owned-by-user")
)

type Conversation struct {
	id        string
	userID    string
	title     string
	createdAt time.Time
	updatedAt time.Time
	messages  []Message
}

func NewConversation(id string, userID string, title string) (*Conversation, error) {
	if id == "" {
		return nil, ErrEmptyConversationID
	}
	if userID == "" {
		return nil, ErrEmptyUserID
	}
	now := time.Now()
	return &Conversation{
		id:        id,
		userID:    userID,
		title:     title,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// UnmarshalConversationFromDatabase unmarshals Conversation from the database.
//
// It should be used only for unmarshalling from the database!
// You can't use UnmarshalConversationFromDatabase as constructor - It may put domain into the invalid state!
func UnmarshalConversationFromDatabase(
	id string,
	userID string,
	title string,
	createdAt time.Time,
	updatedAt time.Time,
	messages []Message,
) (*Conversation, error) {
	conversation, err := NewConversation(id, userID, title)
	if err != nil {
		return nil, err
	}
	conversation.createdAt = createdAt
	conversation.updatedAt = updatedAt
	conversation.messages = append([]Message(nil), messages...)
	return conversation, nil
}

func (c *Conversation) ID() string {
	return c.id
}

func (c *Conversation) UserID() string {
	return c.userID
}

func (c *Conversation) Title() string {
	return c.title
}

func (c *Conversation) CreatedAt() time.Time {
	return c.createdAt
}

func (c *Conversation) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *Conversation) Messages() []Message {
	return append([]Message(nil), c.messages...)
}

func (c *Conversation) UpdateTitle(title string) {
	c.title = title
	c.updatedAt = time.Now()
}

func (c *Conversation) AddMessage(message Message) error {
	if message.ConversationID() != c.id {
		return ErrMessageConversationMismatch
	}
	c.messages = append(c.messages, message)
	c.updatedAt = time.Now()
	return nil
}
