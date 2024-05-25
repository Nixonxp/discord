package usecases

import (
	"context"
	"github.com/Nixonxp/discord/auth/internal/app/models"
	"github.com/google/uuid"
)

type Deps struct {
	// deprecated
	UserRepo    UsersStorage
	UserService UsecaseServiceInterface
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
	user, err := u.UserService.Register(ctx, registerInfo)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, loginInfo LoginUserInfo) (*models.User, error) {
	user, err := u.UserService.Login(ctx, loginInfo)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *AuthUsecase) OauthLogin(_ context.Context, _ OauthLoginRequest) (*models.OauthLoginResult, error) {
	// todo implementation oauth service
	return &models.OauthLoginResult{
		Code: "{code}",
	}, nil
}

func (u *AuthUsecase) OauthLoginCallback(_ context.Context, _ OauthLoginCallbackRequest) (*models.OauthLoginCallbackResult, error) {
	// todo implementation oauth service
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
