package usecases

import (
	"context"
	"github.com/Nixonxp/discord/user/internal/app/models"
)

type UsecaseInterface interface {
	UpdateUser(ctx context.Context, req UpdateUserRequest) (*models.User, error)
	GetUserByLogin(ctx context.Context, req GetUserByLoginRequest) (*models.User, error)
	GetUserFriends(ctx context.Context, req GetUserFriendsRequest) (*models.UserFriendsInfo, error)
	AddToFriendByUserId(ctx context.Context, req AddToFriendByUserIdRequest) (*models.ActionInfo, error)
	AcceptFriendInvite(ctx context.Context, req AcceptFriendInviteRequest) (*models.ActionInfo, error)
	DeclineFriendInvite(ctx context.Context, req DeclineFriendInviteRequest) (*models.ActionInfo, error)
}

type UsersStorage interface {
	UpdateUser(_ context.Context, user *models.User) (*models.User, error)
}
