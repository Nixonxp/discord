package channel

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/channel/internal/app/models"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
	pkgErrors "github.com/Nixonxp/discord/channel/pkg/errors"
	log "github.com/Nixonxp/discord/channel/pkg/logger"
	"github.com/google/uuid"
)

type Deps struct {
	ChannelRepo   usecases.ChannelStorage
	SubscribeRepo usecases.SubscribeStorage
	Log           *log.Logger
}

type ChannelUsecase struct {
	Deps
}

var _ usecases.UsecaseInterface = (*ChannelUsecase)(nil)

func NewChannelUsecase(d Deps) usecases.UsecaseInterface {
	return &ChannelUsecase{
		Deps: d,
	}
}

func (u *ChannelUsecase) AddChannel(ctx context.Context, req usecases.AddChannelRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse(req.CurrentUserId))
	channelID := models.ChannelID(uuid.New())

	err := u.ChannelRepo.CreateChannel(ctx, models.Channel{
		Id:      channelID,
		Name:    req.Name,
		OwnerId: userID,
	})
	if err != nil {
		return nil, pkgErrors.Wrap("create channel", err)
	}

	return &models.ActionInfo{Success: true}, nil
}

func (u *ChannelUsecase) DeleteChannel(ctx context.Context, req usecases.DeleteChannelRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse(req.CurrentUserId))
	channelID := models.ChannelID(uuid.MustParse(req.ChannelId))

	channel, err := u.ChannelRepo.GetChannelById(ctx, channelID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgErrors.Wrap("channel", err)
		}
		return nil, pkgErrors.Wrap("get channel error", err)
	}

	if channel.OwnerId.String() != userID.String() {
		return nil, pkgErrors.Wrap("delete channel error", models.ErrPermDenied)
	}

	err = u.ChannelRepo.DeleteChannel(ctx, channelID)
	if err != nil {
		return nil, pkgErrors.Wrap("delete channel error", err)
	}

	return &models.ActionInfo{Success: true}, nil
}

func (u *ChannelUsecase) JoinChannel(ctx context.Context, req usecases.JoinChannelRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse(req.CurrentUserId))
	channelID := models.ChannelID(uuid.MustParse(req.ChannelId))

	_, err := u.ChannelRepo.GetChannelById(ctx, channelID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgErrors.Wrap("channel", err)
		}
		return nil, pkgErrors.Wrap("search channel error", err)
	}

	newSubscribe := models.SubscribeInfo{
		Id:        models.SubscribeID(uuid.New()),
		ChannelId: channelID,
		UserId:    userID,
	}
	// todo create index to unique subscribes
	err = u.SubscribeRepo.CreateSubscribe(ctx, newSubscribe)
	if err != nil {
		return nil, pkgErrors.Wrap("create subscribe channel error", err)
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ChannelUsecase) LeaveChannel(ctx context.Context, req usecases.LeaveChannelRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse(req.CurrentUserId))
	channelID := models.ChannelID(uuid.MustParse(req.ChannelId))

	err := u.SubscribeRepo.DeleteSubscribe(ctx, channelID, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgErrors.Wrap("delete channel error", models.ErrNotFound)
		}
		return nil, pkgErrors.Wrap("delete subscribe error", err)
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}
