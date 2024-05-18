package usecases

import (
	"context"
	"github.com/Nixonxp/discord/auth/internal/app/models"
)

type UsecaseInterface interface {
	Register(ctx context.Context, registerInfo RegisterUserInfo) (*models.User, error)
	Login(ctx context.Context, loginInfo LoginUserInfo) (*models.User, error)
	OauthLogin(ctx context.Context, req OauthLoginRequest) (*models.OauthLoginResult, error)
	OauthLoginCallback(ctx context.Context, req OauthLoginCallbackRequest) (*models.OauthLoginCallbackResult, error)
}

//go:generate mockery --name=UsersStorage --filename=users_storage_mock.go --disable-version-string
type UsersStorage interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	LoginUser(ctx context.Context, loginInfo *models.Login) (*models.User, error)
}
