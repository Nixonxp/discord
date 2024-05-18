package usecases

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"math/rand/v2"
	"time"
)

type Deps struct {
	ChatRepo ChatsStorage
}

type ChatUsecase struct {
	Deps
}

var _ UsecaseInterface = (*ChatUsecase)(nil)

func NewChatUsecase(d Deps) UsecaseInterface {
	return &ChatUsecase{
		Deps: d,
	}
}

func (u *ChatUsecase) SendUserPrivateMessage(ctx context.Context, req SendUserPrivateMessageRequest) (*models.ActionInfo, error) {
	// временное решение до настоящей реализации
	messageId := rand.IntN(999-1) + 1
	chatId := rand.IntN(3-1) + 1

	_, err := u.ChatRepo.CreateMessage(ctx, &models.Message{
		Id:        uint64(messageId),
		ChatId:    uint64(chatId),
		UserId:    req.UserId,
		Text:      req.Text,
		Timestamp: time.Now(),
	})
	if err != nil {
		return &models.ActionInfo{}, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ChatUsecase) GetUserPrivateMessages(ctx context.Context, req GetUserPrivateMessagesRequest) (*models.Messages, error) {
	chatId := rand.IntN(3-1) + 1
	messages, err := u.ChatRepo.GetMessages(ctx, req.UserId, uint64(chatId))
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, models.ErrEmpty
	}

	return &models.Messages{
		Data: messages,
	}, nil
}
