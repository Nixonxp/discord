package usecases

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
)

type UsecaseInterface interface {
	SendUserPrivateMessage(_ context.Context, req SendUserPrivateMessageRequest) (*models.ActionInfo, error)
	GetUserPrivateMessages(_ context.Context, req GetUserPrivateMessagesRequest) (*models.Messages, error)
}

//go:generate mockery --name=ChatsStorage --filename=chat_storage_mock.go --disable-version-string
type ChatsStorage interface {
	CreateMessage(_ context.Context, message *models.Message) (*models.Message, error)
	GetMessages(_ context.Context, userId uint64, chatId uint64) ([]*models.Message, error)
}
