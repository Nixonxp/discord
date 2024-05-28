package usecases

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/user/internal/app/models"
	pkgerrors "github.com/Nixonxp/discord/user/pkg/errors"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Deps struct {
	UserRepo          UsersStorage
	FriendInvitesRepo FriendInvitesStorage
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
	password := []byte(req.Password)

	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.MinCost)
	if err != nil {
		return &models.User{}, pkgerrors.Wrap("password hashing error", err)
	}

	user := &models.User{
		Id:       userID,
		Login:    req.Login,
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	err = u.UserRepo.CreateUser(ctx, user)
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, req UpdateUserRequest) (*models.User, error) {
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

func (u *UserUsecase) GetUserByLoginAndPassword(ctx context.Context, req GetUserByLoginAndPasswordRequest) (*models.User, error) {
	user, err := u.UserRepo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, models.ErrCredInvalid
		}

		return &models.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return &models.User{}, models.ErrCredInvalid
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

func (u *UserUsecase) GetUserFriends(_ context.Context, req GetUserFriendsRequest) (*models.UserFriendsInfo, error) {
	// todo add repo
	id, _ := uuid.NewUUID()
	return &models.UserFriendsInfo{
		Friends: []*models.User{
			{
				Id:    models.UserID(id),
				Login: "login",
				Name:  "name",
				Email: "test@test.ru",
			},
		},
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
		OwnerId:  models.UserID(uuid.MustParse("e72cebc9-957a-4093-b6b4-5cd820efa9e1")), // todo from auth data
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

func (u *UserUsecase) AcceptFriendInvite(_ context.Context, _ AcceptFriendInviteRequest) (*models.ActionInfo, error) {
	// todo add repo
	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *UserUsecase) DeclineFriendInvite(ctx context.Context, req DeclineFriendInviteRequest) (*models.ActionInfo, error) {
	inviteID := models.InviteId(uuid.MustParse(req.InviteId))

	err := u.FriendInvitesRepo.DeclineInvite(ctx, inviteID)
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}
