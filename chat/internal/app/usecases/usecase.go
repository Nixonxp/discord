package usecases

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
)

type UsecaseInterface interface {
	SendUserPrivateMessage(ctx context.Context, req SendUserPrivateMessageRequest) (*models.ActionInfo, error)
	GetUserPrivateMessages(ctx context.Context, req GetUserPrivateMessagesRequest) (*models.Messages, error)
	CreatePrivateChat(ctx context.Context, req CreatePrivateChatRequest) (*models.Chat, error)
	GetServerMessagesRequest(ctx context.Context, req GetServerMessageRequest) (*models.Messages, error)
	SendServerMessage(ctx context.Context, req SendServerMessageRequest) (*models.ActionInfo, error)
}

type QueueInterface interface {
	CreateMessage(ctx context.Context, message MessageDto) (*models.ActionInfo, error)
}

//go:generate mockery --name=ChatsStorage --filename=chat_storage_mock.go --disable-version-string
type MessagesStorage interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	GetMessages(ctx context.Context, chatId models.ChatID) ([]*models.Message, error)
}

//go:generate mockery --name=ChatsStorage --filename=chat_storage_mock.go --disable-version-string
type ChatStorage interface {
	CreateChat(ctx context.Context, chat *models.Chat) error
	GetChatById(ctx context.Context, chatId models.ChatID) (*models.Chat, error)
	GetChatByMetadataAndType(ctx context.Context, metadata string, chatType string) (*models.Chat, error)
}
