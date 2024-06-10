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

//go:generate mockery --name=ServerStorage --filename=server_storage_mock.go --disable-version-string
type ServerStorage interface {
	CreateServer(ctx context.Context, server models.ServerInfo) error
	SearchServers(ctx context.Context, serverName string) ([]*models.ServerInfo, error)
	GetServerById(ctx context.Context, id string) (*models.ServerInfo, error)
}

//go:generate mockery --name=SubscribeStorage --filename=subscribe_storage_mock.go --disable-version-string
type SubscribeStorage interface {
	CreateSubscribe(ctx context.Context, server models.SubscribeInfo) error
	DeleteSubscribe(ctx context.Context, serverId models.ServerID, userId models.UserID) error
	GetByUserId(ctx context.Context, userId models.UserID) ([]*models.SubscribeInfo, error)
}

//go:generate mockery --name=ServiceChatInterface --filename=service_chat_mock.go --disable-version-string
type ServiceChatInterface interface {
	SendServerMessage(ctx context.Context, msg PublishMessageOnServerRequest) (*models.ActionInfo, error)
	GetServerMessages(ctx context.Context, msg GetMessagesFromServerRequest) (*models.GetMessagesInfo, error)
}
