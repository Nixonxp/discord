package usecases

import (
	"context"
	"github.com/Nixonxp/discord/auth/internal/app/models"
	"github.com/google/uuid"
)

type Deps struct {
	UserRepo UsersStorage
}

type AuthUsecase struct {
	Deps
}

var _ UsecaseInterface = (*AuthUsecase)(nil)

func NewAuthUsecase(d Deps) UsecaseInterface {
	return &AuthUsecase{
		Deps: d,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, registerInfo RegisterUserInfo) (*models.User, error) {
	userID := models.UserID(uuid.New())
	user, err := u.UserRepo.CreateUser(ctx, models.User{
		UserID:   userID,
		Login:    registerInfo.Login,
		Name:     registerInfo.Name,
		Email:    registerInfo.Email,
		Password: registerInfo.Password,
	})
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, loginInfo LoginUserInfo) (*models.User, error) {
	user, err := u.UserRepo.LoginUser(ctx, models.Login{
		Login:    loginInfo.Login,
		Password: loginInfo.Password,
	})
	if err != nil {
		return &models.User{}, err
	}

	return user, nil
}

func (u *AuthUsecase) OauthLogin(_ context.Context, _ OauthLoginRequest) (*models.OauthLoginResult, error) {
	// todo implementation
	return &models.OauthLoginResult{
		Code: "{code}",
	}, nil
}

func (u *AuthUsecase) OauthLoginCallback(_ context.Context, _ OauthLoginCallbackRequest) (*models.OauthLoginCallbackResult, error) {
	return &models.OauthLoginCallbackResult{
		User: models.User{
			UserID:   models.UserID(uuid.New()),
			Login:    "login user",
			Name:     "user name",
			Email:    "test@test.ru",
			Password: "",
		},
	}, nil
}
