package queue

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	"github.com/google/uuid"
	"time"
)

type QueueUsecase struct {
	chatRepo usecases.MessagesStorage
}

func NewQueueUsecase(chatRepo usecases.MessagesStorage) *QueueUsecase {
	return &QueueUsecase{
		chatRepo: chatRepo,
	}
}

func (u *QueueUsecase) CreateMessage(ctx context.Context, message usecases.MessageDto) (*models.ActionInfo, error) {
	if message.Id == "" {
		message.Id = uuid.New().String()
	}

	err := u.chatRepo.CreateMessage(ctx, &models.Message{
		Id:        models.MessageID(uuid.MustParse(message.Id)),
		ChatId:    models.ChatID(uuid.MustParse(message.ChatId)),
		OwnerId:   models.OwnerID(uuid.MustParse(message.OwnerId)),
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
