package usecases

import (
	"context"
	"github.com/Nixonxp/discord/server/internal/app/models"
)

type UsecaseInterface interface {
	CreateServer(ctx context.Context, req CreateServerRequest) (*models.ServerInfo, error)
	SearchServer(ctx context.Context, req SearchServerRequest) ([]*models.ServerInfo, error)
	SubscribeServer(ctx context.Context, req SubscribeServerRequest) (*models.ActionInfo, error)
	UnsubscribeServer(ctx context.Context, req UnsubscribeServerRequest) (*models.ActionInfo, error)
	SearchServerByUserId(ctx context.Context, req SearchServerByUserIdRequest) ([]string, error)
	InviteUserToServer(ctx context.Context, req InviteUserToServerRequest) (*models.ActionInfo, error)
	PublishMessageOnServer(ctx context.Context, req PublishMessageOnServerRequest) (*models.ActionInfo, error)
	GetMessagesFromServer(ctx context.Context, req GetMessagesFromServerRequest) (*models.GetMessagesInfo, error)
}

type ServerStorage interface {
	CreateServer(ctx context.Context, server models.ServerInfo) error
	SearchServers(ctx context.Context, serverName string) ([]*models.ServerInfo, error)
	GetServerById(ctx context.Context, id string) (*models.ServerInfo, error)
}

type SubscribeStorage interface {
	CreateSubscribe(ctx context.Context, server models.SubscribeInfo) error
	DeleteSubscribe(ctx context.Context, serverId models.ServerID, userId models.UserID) error
	GetByUserId(ctx context.Context, userId models.UserID) ([]*models.SubscribeInfo, error)
}
