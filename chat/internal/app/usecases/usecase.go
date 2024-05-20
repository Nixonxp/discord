package usecases

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
)

type UsecaseInterface interface {
	SendUserPrivateMessage(ctx context.Context, req SendUserPrivateMessageRequest) (*models.ActionInfo, error)
	GetUserPrivateMessages(ctx context.Context, req GetUserPrivateMessagesRequest) (*models.Messages, error)
}

//go:generate mockery --name=ChatsStorage --filename=chat_storage_mock.go --disable-version-string
type ChatsStorage interface {
	CreateMessage(ctx context.Context, message *models.Message) error
	GetMessages(ctx context.Context, chatId models.ChatID) ([]*models.Message, error)
}
