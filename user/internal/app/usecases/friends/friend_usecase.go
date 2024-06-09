package usecases

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/user/internal/app/models"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	pkgerrors "github.com/Nixonxp/discord/user/pkg/errors"
	log "github.com/Nixonxp/discord/user/pkg/logger"
	"github.com/Nixonxp/discord/user/pkg/transaction_manager"
	"github.com/google/uuid"
)

type Deps struct {
	TransactionManager usecases.TransactionManager
	UserRepo           usecases.UsersStorage
	FriendInvitesRepo  usecases.FriendInvitesStorage
	UserFriendsRepo    usecases.UserFriendsStorage
	Log                *log.Logger
}

type FriendUsecase struct {
	Deps
}

var _ usecases.FriendUsecaseInterface = (*FriendUsecase)(nil)

func NewFriendUsecase(d Deps) *FriendUsecase {
	return &FriendUsecase{
		Deps: d,
	}
}

func (u *FriendUsecase) GetUserFriends(ctx context.Context, req usecases.GetUserFriendsRequest) (*models.UserFriendsInfo, error) {
	userID := models.UserID(uuid.MustParse(req.UserId))

	friends, err := u.UserFriendsRepo.GetUserFriendsByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &models.UserFriendsInfo{
		Friends: friends,
	}, nil
}

func (u *FriendUsecase) AddToFriendByUserId(ctx context.Context, req usecases.AddToFriendByUserIdRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse(req.UserId))
	_, err := u.UserRepo.GetUserById(ctx, userID)
	if err != nil {
		return nil, pkgerrors.Wrap("user not found", models.ErrNotFound)
	}

	inviteId, _ := uuid.NewUUID()
	err = u.FriendInvitesRepo.CreateInvite(ctx, &models.FriendInvite{
		InviteId: models.InviteId(inviteId),
		OwnerId:  models.UserID(uuid.MustParse(req.OwnerId)),
		UserId:   userID,
		Status:   models.PendingStatus,
	})
	if err != nil {
		return nil, pkgerrors.Wrap("create user invite error", err)
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *FriendUsecase) GetUserInvites(ctx context.Context, req usecases.GetUserInvitesRequest) (*models.UserInvitesInfo, error) {
	userID := models.UserID(uuid.MustParse(req.UserId))

	invites, err := u.FriendInvitesRepo.GetInvitesByUserId(ctx, userID)
	if err != nil {
		return nil, pkgerrors.Wrap("get user invites error", err)
	}

	return invites, nil
}

func (u *FriendUsecase) AcceptFriendInvite(ctx context.Context, req usecases.AcceptFriendInviteRequest) (*models.ActionInfo, error) {
	inviteID := models.InviteId(uuid.MustParse(req.InviteId))

	invite, err := u.FriendInvitesRepo.GetInviteById(ctx, inviteID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgerrors.Wrap("invite ", models.ErrNotFound)
		}
		return nil, err
	}

	if invite.UserId.String() != req.CurrentUserId {
		return nil, pkgerrors.Wrap("invite permission denied", models.PermissionDenied)
	}

	if invite.Status != models.PendingStatus {
		return nil, pkgerrors.Wrap("invite not in pending status", models.PermissionDenied)
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
		return nil, pkgerrors.Wrap("accept invite error", err)
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *FriendUsecase) DeclineFriendInvite(ctx context.Context, req usecases.DeclineFriendInviteRequest) (*models.ActionInfo, error) {
	inviteID := models.InviteId(uuid.MustParse(req.InviteId))

	invite, err := u.FriendInvitesRepo.GetInviteById(ctx, inviteID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgerrors.Wrap("invite", models.ErrNotFound)
		}
		return nil, err
	}

	if invite.UserId.String() != req.CurrentUserId {
		return nil, pkgerrors.Wrap("invite permission denied", models.PermissionDenied)
	}

	err = u.FriendInvitesRepo.DeclineInvite(ctx, inviteID)
	if err != nil {
		return nil, pkgerrors.Wrap("invite decline error", err)
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *FriendUsecase) DeleteFromFriend(ctx context.Context, req usecases.DeleteFromFriendRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse(req.CurrentUserId))
	friendID := models.UserID(uuid.MustParse(req.FriendId))

	err := u.UserFriendsRepo.DeleteFriend(ctx, userID, friendID)
	if err != nil {
		return nil, pkgerrors.Wrap("delete froms friend error", err)
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}
