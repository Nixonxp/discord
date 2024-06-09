package usecases

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/user/internal/app/models"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	pkgerrors "github.com/Nixonxp/discord/user/pkg/errors"
	log "github.com/Nixonxp/discord/user/pkg/logger"
	"github.com/google/uuid"
)

type Deps struct {
	usecases.TransactionManager
	UserRepo usecases.UsersStorage
	Log      *log.Logger
}

type UserUsecase struct {
	Deps
}

var _ usecases.UserUsecaseInterface = (*UserUsecase)(nil)

func NewUserUsecase(d Deps) *UserUsecase {
	return &UserUsecase{
		Deps: d,
	}
}

func (u *UserUsecase) CreateUser(ctx context.Context, req usecases.CreateUserRequest) (*models.User, error) {
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
		return nil, pkgerrors.Wrap("create user error", err)
	}

	return user, nil
}

func (u *UserUsecase) CreateOrGetUser(ctx context.Context, req usecases.CreateOrGetUserRequest) (*models.User, error) {
	user, err := u.UserRepo.GetUserByOauthId(ctx, req.OauthId)
	if err != nil {
		if !errors.Is(err, models.ErrNotFound) {
			return nil, pkgerrors.Wrap("user not found", models.ErrNotFound)
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
			return nil, pkgerrors.Wrap("user create error", err)
		}
	}

	return user, nil
}

func (u *UserUsecase) UpdateUser(ctx context.Context, req usecases.UpdateUserRequest) (*models.User, error) {
	if req.CurrentUserId != req.Id {
		return nil, pkgerrors.Wrap("user update permission error", models.PermissionDenied)
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

func (u *UserUsecase) GetUserForLogin(ctx context.Context, req usecases.GetUserByLoginRequest) (*models.User, error) {
	user, err := u.UserRepo.GetUserByLogin(ctx, req.Login)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgerrors.Wrap("user by login not found", models.ErrNotFound)
		}

		return &models.User{}, err
	}

	return user, nil
}
