package app

import (
	"github.com/kimnattanan/graph-rag-service/internal/user/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/user/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	Register      command.RegisterHandler
	Logout        command.LogoutHandler
	DeleteAccount command.DeleteAccountHandler
}

type Queries struct {
	Login          query.LoginHandler
	Authenticate   query.AuthenticateHandler
	GetCurrentUser query.GetCurrentUserHandler
}
