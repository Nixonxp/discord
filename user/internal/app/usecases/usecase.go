package usecases

import (
	"context"
	"github.com/Nixonxp/discord/user/internal/app/models"
	"github.com/jackc/pgx/v5"
)

type UserUsecaseInterface interface {
	CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error)
	CreateOrGetUser(ctx context.Context, req CreateOrGetUserRequest) (*models.User, error)
	UpdateUser(ctx context.Context, req UpdateUserRequest) (*models.User, error)
	GetUserForLogin(ctx context.Context, req GetUserByLoginRequest) (*models.User, error)
}

type FriendUsecaseInterface interface {
	GetUserFriends(ctx context.Context, req GetUserFriendsRequest) (*models.UserFriendsInfo, error)
	GetUserInvites(ctx context.Context, req GetUserInvitesRequest) (*models.UserInvitesInfo, error)
	AddToFriendByUserId(ctx context.Context, req AddToFriendByUserIdRequest) (*models.ActionInfo, error)
	AcceptFriendInvite(ctx context.Context, req AcceptFriendInviteRequest) (*models.ActionInfo, error)
	DeclineFriendInvite(ctx context.Context, req DeclineFriendInviteRequest) (*models.ActionInfo, error)
	DeleteFromFriend(ctx context.Context, req DeleteFromFriendRequest) (*models.ActionInfo, error)
}

//go:generate mockery --name=UsersStorage --filename=users_storage_mock.go --disable-version-string
type UsersStorage interface {
	CreateUser(ctx context.Context, user *models.User) error
	UpdateUser(ctx context.Context, user *models.User) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetUserByOauthId(ctx context.Context, oauthId string) (*models.User, error)
	GetUserById(ctx context.Context, userId models.UserID) (*models.User, error)
}

//go:generate mockery --name=TransactionManager --filename=transaction_manager_mock.go --disable-version-string
type TransactionManager interface {
	RunReadCommitted(ctx context.Context, accessMode pgx.TxAccessMode, f func(ctx context.Context) error) error
}

//go:generate mockery --name=FriendInvitesStorage --filename=friend_invites_storage_mock.go --disable-version-string
type FriendInvitesStorage interface {
	CreateInvite(ctx context.Context, invite *models.FriendInvite) error
	GetInvitesByUserId(ctx context.Context, userId models.UserID) (*models.UserInvitesInfo, error)
	GetInviteById(ctx context.Context, inviteId models.InviteId) (*models.FriendInvite, error)
	DeclineInvite(ctx context.Context, inviteId models.InviteId) error
	AcceptInvite(ctx context.Context, inviteId models.InviteId) error
}

//go:generate mockery --name=UserFriendsStorage --filename=user_friends_storage_mock.go --disable-version-string
type UserFriendsStorage interface {
	AddToFriend(ctx context.Context, friendInfo []*models.UserFriends) error
	GetUserFriendsByUserId(ctx context.Context, userId models.UserID) ([]*models.Friend, error)
	DeleteFriend(ctx context.Context, userId models.UserID, friendId models.UserID) error
}
