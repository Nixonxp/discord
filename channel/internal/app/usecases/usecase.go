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
	CreateChannel(ctx context.Context, channel models.Channel) error
	GetChannelById(ctx context.Context, id models.ChannelID) (*models.Channel, error)
	DeleteChannel(ctx context.Context, channelId models.ChannelID) error
}

type SubscribeStorage interface {
	CreateSubscribe(ctx context.Context, subscribe models.SubscribeInfo) error
	DeleteSubscribe(ctx context.Context, channelId models.ChannelID, userId models.UserID) error
}
