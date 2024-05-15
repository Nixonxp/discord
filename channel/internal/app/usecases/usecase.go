package usecases

import (
	"context"
	"github.com/Nixonxp/discord/channel/internal/app/models"
)

type UsecaseInterface interface {
	AddChannel(ctx context.Context, req AddChannelRequest) (*models.ActionInfo, error)
	DeleteChannel(ctx context.Context, req DeleteChannelRequest) (*models.ActionInfo, error)
	JoinChannel(ctx context.Context, req JoinChannelRequest) (*models.ActionInfo, error)
	LeaveChannel(ctx context.Context, req LeaveChannelRequest) (*models.ActionInfo, error)
}

type ChannelStorage interface {
	CreateChannel(ctx context.Context, channel models.Channel) (*models.Channel, error)
}
