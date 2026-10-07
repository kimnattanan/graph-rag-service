package query

import (
	"time"

	"github.com/kimnattanan/graph-rag-service/internal/user/domain/user"
)

type User struct {
	ID          string
	Email       string
	Username    string
	Role        string
	Permissions []string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func UserFromDomain(account *user.User) User {
	return User{
		ID:          account.ID(),
		Email:       account.Email(),
		Username:    account.Username(),
		Role:        string(account.Role()),
		Permissions: account.Permissions(),
		CreatedAt:   account.CreatedAt(),
		UpdatedAt:   account.UpdatedAt(),
	}
}
