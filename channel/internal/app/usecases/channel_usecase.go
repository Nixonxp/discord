package usecases

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/channel/internal/app/models"
	pkgErrors "github.com/Nixonxp/discord/channel/pkg/errors"
	log "github.com/Nixonxp/discord/channel/pkg/logger"
	"github.com/google/uuid"
)

type Deps struct {
	ChannelRepo   ChannelStorage
	SubscribeRepo SubscribeStorage
	Log           *log.Logger
}

type ChannelUsecase struct {
	Deps
}

var _ UsecaseInterface = (*ChannelUsecase)(nil)

func NewChannelUsecase(d Deps) UsecaseInterface {
	return &ChannelUsecase{
		Deps: d,
	}
}

func (u *ChannelUsecase) AddChannel(ctx context.Context, req AddChannelRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) // todo from auth
	channelID := models.ChannelID(uuid.New())

	err := u.ChannelRepo.CreateChannel(ctx, models.Channel{
		Id:      channelID,
		Name:    req.Name,
		OwnerId: userID,
	})
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{Success: true}, nil
}

func (u *ChannelUsecase) DeleteChannel(ctx context.Context, req DeleteChannelRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) // todo from auth
	channelID := models.ChannelID(uuid.MustParse(req.ChannelId))

	channel, err := u.ChannelRepo.GetChannelById(ctx, channelID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgErrors.Wrap("channel not found", err)
		}
		return nil, err
	}

	if channel.OwnerId.String() != userID.String() {
		return nil, models.ErrPermDenied
	}

	err = u.ChannelRepo.DeleteChannel(ctx, channelID)
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{Success: true}, nil
}

func (u *ChannelUsecase) JoinChannel(ctx context.Context, req JoinChannelRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) // todo from auth
	channelID := models.ChannelID(uuid.MustParse(req.ChannelId))

	_, err := u.ChannelRepo.GetChannelById(ctx, channelID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.New("channel not found")
		}
		return nil, err
	}

	newSubscribe := models.SubscribeInfo{
		Id:        models.SubscribeID(uuid.New()),
		ChannelId: channelID,
		UserId:    userID,
	}
	// todo create index to unique subscribes
	err = u.SubscribeRepo.CreateSubscribe(ctx, newSubscribe)
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ChannelUsecase) LeaveChannel(ctx context.Context, req LeaveChannelRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) // todo from auth
	channelID := models.ChannelID(uuid.MustParse(req.ChannelId))

	err := u.SubscribeRepo.DeleteSubscribe(ctx, channelID, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.New("subscribe not found")
		}
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}
