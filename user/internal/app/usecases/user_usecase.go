package usecases

import (
	"context"
	"github.com/Nixonxp/discord/user/internal/app/models"
)

type Deps struct {
	UserRepo UsersStorage
}

type UserUsecase struct {
	Deps
}

var _ UsecaseInterface = (*UserUsecase)(nil)

func NewUserUsecase(d Deps) UsecaseInterface {
	return &UserUsecase{
		Deps: d,
	}
}

func (u *UserUsecase) UpdateUser(ctx context.Context, req UpdateUserRequest) (*models.User, error) {
	user, err := u.UserRepo.UpdateUser(ctx, &models.User{
		UserID: req.Id,
		Login:  req.Login,
		Name:   req.Name,
		Email:  req.Email,
	})
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (u *UserUsecase) GetUserByLogin(_ context.Context, req GetUserByLoginRequest) (*models.User, error) {
	// todo add repo
	return &models.User{
		UserID: 1,
		Login:  req.Login,
		Name:   "name",
		Email:  "test@test.ru",
	}, nil
}

func (u *UserUsecase) GetUserFriends(_ context.Context, req GetUserFriendsRequest) (*models.UserFriendsInfo, error) {
	// todo add repo
	return &models.UserFriendsInfo{
		Friends: []*models.User{
			{
				UserID: 1,
				Login:  "login",
				Name:   "name",
				Email:  "test@test.ru",
			},
		},
	}, nil
}

func (u *UserUsecase) AddToFriendByUserId(_ context.Context, _ AddToFriendByUserIdRequest) (*models.ActionInfo, error) {
	// todo add repo
	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *UserUsecase) AcceptFriendInvite(_ context.Context, _ AcceptFriendInviteRequest) (*models.ActionInfo, error) {
	// todo add repo
	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *UserUsecase) DeclineFriendInvite(_ context.Context, _ DeclineFriendInviteRequest) (*models.ActionInfo, error) {
	// todo add repo
	return &models.ActionInfo{
		Success: true,
	}, nil
}
