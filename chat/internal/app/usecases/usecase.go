package usecases

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
)

type UsecaseInterface interface {
	SendUserPrivateMessage(_ context.Context, req SendUserPrivateMessageRequest) (*models.ActionInfo, error)
	GetUserPrivateMessages(_ context.Context, req GetUserPrivateMessagesRequest) (*models.Messages, error)
}

type ChatsStorage interface {
	CreateMessage(_ context.Context, message models.Message) (*models.Message, error)
}
