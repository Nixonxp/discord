package usecases

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/user/internal/app/models"
	pkgerrors "github.com/Nixonxp/discord/user/pkg/errors"
	"github.com/Nixonxp/discord/user/pkg/transaction_manager"
	"github.com/google/uuid"
)

type Deps struct {
	TransactionManager
	UserRepo          UsersStorage
	FriendInvitesRepo FriendInvitesStorage
	UserFriendsRepo   UserFriendsStorage
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

func (u *UserUsecase) CreateUser(ctx context.Context, req CreateUserRequest) (*models.User, error) {
	userID := models.UserID(uuid.New())

	user := &models.User{
		Id:       userID,
		Login:    req.Login,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := u.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (u *UserUsecase) CreateOrGetUser(ctx context.Context, req CreateOrGetUserRequest) (*models.User, error) {
	user, err := u.UserRepo.GetUserByOauthId(ctx, req.OauthId)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			return nil, err
		}

		userID := models.UserID(uuid.New())
		user = &models.User{
			Id:             userID,
			Login:          req.Login,
			Name:           req.Name,
			Email:          req.Email,
			Password:       req.Password,
			OauthId:        req.OauthId,
			AvatarPhotoUrl: req.AvatarPhotoUrl,
		}

		err = u.UserRepo.CreateUser(ctx, user)
		if err != nil {
			return &models.User{}, err
		}
	}

	return user, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, req UpdateUserRequest) (*models.User, error) {
	if req.CurrentUserId != req.Id {
		return nil, models.PermissionDenied
	}
	userID := models.UserID(uuid.MustParse(req.Id))
	err := u.UserRepo.UpdateUser(ctx, &models.User{
		Id:             userID,
		Login:          req.Login,
		Name:           req.Name,
		Email:          req.Email,
		AvatarPhotoUrl: req.AvatarPhotoUrl,
	})
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgerrors.Wrap("user for update not found", err)
		}
		return nil, err
	}

	user := &models.User{
		Id:             userID,
		Login:          req.Login,
		Name:           req.Name,
		Email:          req.Email,
		AvatarPhotoUrl: req.AvatarPhotoUrl,
	}

	return user, nil
}

func (u *UserUsecase) GetUserForLogin(ctx context.Context, req GetUserByLoginRequest) (*models.User, error) {
	user, err := u.UserRepo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrCredInvalid
		}

		return &models.User{}, err
	}

	return user, nil
}

func (u *UserUsecase) GetUserByLogin(ctx context.Context, req GetUserByLoginRequest) (*models.User, error) {
	user, err := u.UserRepo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUsecase) GetUserFriends(ctx context.Context, req GetUserFriendsRequest) (*models.UserFriendsInfo, error) {
	userID := models.UserID(uuid.MustParse(req.UserId))

	friends, err := u.UserFriendsRepo.GetUserFriendsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.UserFriendsInfo{
		Friends: friends,
	}, nil
}

func (u *UserUsecase) AddToFriendByUserId(ctx context.Context, req AddToFriendByUserIdRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse(req.UserId))
	_, err := u.UserRepo.GetUserById(ctx, userID)
	if err != nil {
		return nil, err
	}

	inviteId, _ := uuid.NewUUID()
	err = u.FriendInvitesRepo.CreateInvite(ctx, &models.FriendInvite{
		InviteId: models.InviteId(inviteId),
		OwnerId:  models.UserID(uuid.MustParse(req.OwnerId)),
		UserId:   userID,
		Status:   models.PendingStatus,
	})
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *UserUsecase) GetUserInvites(ctx context.Context, req GetUserInvitesRequest) (*models.UserInvitesInfo, error) {
	userID := models.UserID(uuid.MustParse(req.UserId))

	invites, err := u.FriendInvitesRepo.GetInvitesByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	return invites, nil
}

func (u *UserUsecase) AcceptFriendInvite(ctx context.Context, req AcceptFriendInviteRequest) (*models.ActionInfo, error) {
	inviteID := models.InviteId(uuid.MustParse(req.InviteId))

	invite, err := u.FriendInvitesRepo.GetInviteById(ctx, inviteID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.New("invite not found")
		}
		return nil, err
	}

	if invite.UserId.String() != req.CurrentUserId {
		return nil, models.PermissionDenied
	}

	if invite.Status != models.PendingStatus {
		return nil, errors.New("invite not in pending status")
	}

	err = u.TransactionManager.RunReadCommitted(ctx, transaction_manager.ReadWrite,
		func(txCtx context.Context) error {
			// подтверждаем заявку
			err := u.FriendInvitesRepo.AcceptInvite(ctx, inviteID)
			if err != nil {
				return err
			}

			// добавляем записи для связки друзей
			friendsInfo := []*models.UserFriends{
				{
					UserId:   invite.UserId,
					FriendId: invite.OwnerId,
				},
				{
					UserId:   invite.OwnerId,
					FriendId: invite.UserId,
				},
			}

			err = u.UserFriendsRepo.AddToFriend(ctx, friendsInfo)
			if err != nil {
				return err
			}

			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *UserUsecase) DeclineFriendInvite(ctx context.Context, req DeclineFriendInviteRequest) (*models.ActionInfo, error) {
	inviteID := models.InviteId(uuid.MustParse(req.InviteId))

	invite, err := u.FriendInvitesRepo.GetInviteById(ctx, inviteID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.New("invite not found")
		}
		return nil, err
	}

	if invite.UserId.String() != req.CurrentUserId {
		return nil, models.PermissionDenied
	}

	err = u.FriendInvitesRepo.DeclineInvite(ctx, inviteID)
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *UserUsecase) DeleteFromFriend(ctx context.Context, req DeleteFromFriendRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse(req.CurrentUserId))
	friendID := models.UserID(uuid.MustParse(req.FriendId))

	err := u.UserFriendsRepo.DeleteFriend(ctx, userID, friendID)
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}
