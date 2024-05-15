package usecases

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
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
	_, err := u.ChatRepo.CreateMessage(ctx, models.Message{
		Id:        1,
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

func (u *ChatUsecase) GetUserPrivateMessages(_ context.Context, _ GetUserPrivateMessagesRequest) (*models.Messages, error) {
	// todo add repo
	return &models.Messages{
		Data: []*models.Message{
			{
				Id:        1,
				Text:      "text",
				Timestamp: time.Now(),
			},
		},
	}, nil
}
