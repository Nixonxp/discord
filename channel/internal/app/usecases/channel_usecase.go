package usecases

import (
	"context"
	"github.com/Nixonxp/discord/channel/internal/app/models"
)

type Deps struct {
	ChannelRepo ChannelStorage
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
	_, err := u.ChannelRepo.CreateChannel(ctx, models.Channel{
		Id:   1,
		Name: req.Name,
	})
	if err != nil {
		return &models.ActionInfo{}, err
	}

	return &models.ActionInfo{Success: true}, nil
}

func (u *ChannelUsecase) DeleteChannel(_ context.Context, _ DeleteChannelRequest) (*models.ActionInfo, error) {
	// todo add repo
	return &models.ActionInfo{Success: true}, nil
}

func (u *ChannelUsecase) JoinChannel(_ context.Context, _ JoinChannelRequest) (*models.ActionInfo, error) {
	// todo add repo
	return &models.ActionInfo{Success: true}, nil
}

func (u *ChannelUsecase) LeaveChannel(_ context.Context, _ LeaveChannelRequest) (*models.ActionInfo, error) {
	// todo add repo
	return &models.ActionInfo{Success: true}, nil
}
