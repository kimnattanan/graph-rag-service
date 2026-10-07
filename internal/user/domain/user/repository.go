package user

import "context"

type Repository interface {
	AddUser(ctx context.Context, account *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, id string) (*User, error)
	DeleteUser(ctx context.Context, id string) error
}

type SessionRepository interface {
	AddSession(ctx context.Context, session Session) error
	GetSession(ctx context.Context, sessionID string) (Session, error)
	DeleteSession(ctx context.Context, sessionID string) error
}
