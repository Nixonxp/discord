package usecases

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"github.com/google/uuid"
	"time"
)

type QueueUsecase struct {
	chatRepo MessagesStorage
}

func NewQueueUsecase(chatRepo MessagesStorage) *QueueUsecase {
	return &QueueUsecase{
		chatRepo: chatRepo,
	}
}

func (u *QueueUsecase) CreateMessage(ctx context.Context, message MessageDto) (*models.ActionInfo, error) {
	err := u.chatRepo.CreateMessage(ctx, &models.Message{
		Id:        models.MessageID(uuid.MustParse(message.Id)),
		ChatId:    models.ChatID(uuid.MustParse(message.ChatId)),
		OwnerId:   models.OwnerID(uuid.MustParse(message.ChatId)),
		Text:      message.Text,
		Timestamp: time.Now(),
	})
	if err != nil {
		return &models.ActionInfo{}, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}
