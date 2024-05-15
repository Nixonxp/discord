package usecases

import (
	"context"
	"github.com/Nixonxp/discord/server/internal/app/models"
)

type UsecaseInterface interface {
	CreateServer(ctx context.Context, req CreateServerRequest) (*models.ServerInfo, error)
	SearchServer(ctx context.Context, req SearchServerRequest) (*models.ServerInfo, error)
	SubscribeServer(ctx context.Context, req SubscribeServerRequest) (*models.ActionInfo, error)
	UnsubscribeServer(ctx context.Context, req UnsubscribeServerRequest) (*models.ActionInfo, error)
	SearchServerByUserId(ctx context.Context, req SearchServerByUserIdRequest) (*models.ServerInfo, error)
	InviteUserToServer(ctx context.Context, req InviteUserToServerRequest) (*models.ActionInfo, error)
	PublishMessageOnServer(ctx context.Context, req PublishMessageOnServerRequest) (*models.ActionInfo, error)
	GetMessagesFromServer(ctx context.Context, req GetMessagesFromServerRequest) (*models.GetMessagesInfo, error)
}

type ServerStorage interface {
	CreateServer(_ context.Context, server models.ServerInfo) (*models.ServerInfo, error)
}
