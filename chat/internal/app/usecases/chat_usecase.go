package usecases

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"github.com/google/uuid"
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

// мок, в будущем будет браться из таблицы
const chatIdStringMock = "3d0222e1-4b58-4fa7-a38c-171ee345b14e"

func (u *ChatUsecase) SendUserPrivateMessage(ctx context.Context, req SendUserPrivateMessageRequest) (*models.ActionInfo, error) {

	chatUUID := uuid.MustParse(chatIdStringMock)

	err := u.ChatRepo.CreateMessage(ctx, &models.Message{
		Id:        models.MessageID(uuid.New()),
		ChatId:    models.ChatID(chatUUID),
		OwnerId:   models.OwnerID(uuid.New()),
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
	chatUUID := uuid.MustParse(chatIdStringMock)
	chatId := models.ChatID(chatUUID)

	messages, err := u.ChatRepo.GetMessages(ctx, chatId)
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
