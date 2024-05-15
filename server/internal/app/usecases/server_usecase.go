package usecases

import (
	"context"
	"github.com/Nixonxp/discord/server/internal/app/models"
	"time"
)

type Deps struct {
	ServerRepo ServerStorage
}

type ServerUsecase struct {
	Deps
}

var _ UsecaseInterface = (*ServerUsecase)(nil)

func NewServerUsecase(d Deps) UsecaseInterface {
	return &ServerUsecase{
		Deps: d,
	}
}

func (u *ServerUsecase) CreateServer(ctx context.Context, req CreateServerRequest) (*models.ServerInfo, error) {
	server, err := u.ServerRepo.CreateServer(ctx, models.ServerInfo{
		Id:   1,
		Name: req.Name,
	})

	if err != nil {
		return &models.ServerInfo{}, err
	}

	return server, nil
}

func (u *ServerUsecase) SearchServer(_ context.Context, _ SearchServerRequest) (*models.ServerInfo, error) {
	// todo add repo
	return &models.ServerInfo{
		Id:   1,
		Name: "name",
	}, nil
}

func (u *ServerUsecase) SubscribeServer(_ context.Context, _ SubscribeServerRequest) (*models.ActionInfo, error) {
	// todo add repo
	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ServerUsecase) UnsubscribeServer(_ context.Context, _ UnsubscribeServerRequest) (*models.ActionInfo, error) {
	// todo add repo
	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ServerUsecase) SearchServerByUserId(_ context.Context, _ SearchServerByUserIdRequest) (*models.ServerInfo, error) {
	// todo add repo
	return &models.ServerInfo{
		Id:   1,
		Name: "name",
	}, nil
}

func (u *ServerUsecase) InviteUserToServer(_ context.Context, _ InviteUserToServerRequest) (*models.ActionInfo, error) {
	// todo add repo
	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ServerUsecase) PublishMessageOnServer(_ context.Context, _ PublishMessageOnServerRequest) (*models.ActionInfo, error) {
	// todo add repo
	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ServerUsecase) GetMessagesFromServer(_ context.Context, _ GetMessagesFromServerRequest) (*models.GetMessagesInfo, error) {
	// todo add repo
	return &models.GetMessagesInfo{
		Messages: []*models.Message{
			{
				Id:        1,
				Text:      "text",
				Timestamp: time.Now(),
			},
		},
	}, nil
}
