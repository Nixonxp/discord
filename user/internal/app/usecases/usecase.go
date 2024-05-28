package usecases

import (
	"context"
	"github.com/Nixonxp/discord/user/internal/app/models"
)

type UsecaseInterface interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error)
	UpdateUser(ctx context.Context, req UpdateUserRequest) (*models.User, error)
	GetUserByLoginAndPassword(ctx context.Context, req GetUserByLoginAndPasswordRequest) (*models.User, error)
	GetUserByLogin(ctx context.Context, req GetUserByLoginRequest) (*models.User, error)
	GetUserFriends(ctx context.Context, req GetUserFriendsRequest) (*models.UserFriendsInfo, error)
	GetUserInvites(ctx context.Context, req GetUserInvitesRequest) (*models.UserInvitesInfo, error)
	AddToFriendByUserId(ctx context.Context, req AddToFriendByUserIdRequest) (*models.ActionInfo, error)
	AcceptFriendInvite(ctx context.Context, req AcceptFriendInviteRequest) (*models.ActionInfo, error)
	DeclineFriendInvite(ctx context.Context, req DeclineFriendInviteRequest) (*models.ActionInfo, error)
}

type UsersStorage interface {
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserById(ctx context.Context, userId models.UserID) (*models.User, error)
}

type FriendInvitesStorage interface {
	CreateInvite(ctx context.Context, invite *models.FriendInvite) error
	GetInvitesByUserId(ctx context.Context, userId models.UserID) (*models.UserInvitesInfo, error)
	DeclineInvite(ctx context.Context, inviteId models.InviteId) error
}
