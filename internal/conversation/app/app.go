package app

import (
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app/command"
	"github.com/kimnattanan/graph-rag-service/internal/conversation/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	CreateConversation command.CreateConversationHandler
	DeleteConversation command.DeleteConversationHandler
	SendMessage        command.SendMessageHandler
}

type Queries struct {
	ListConversations query.ListConversationsHandler
	GetConversation   query.GetConversationHandler
	ListMessages      query.ListMessagesHandler
	GetMessage        query.GetMessageHandler
}
