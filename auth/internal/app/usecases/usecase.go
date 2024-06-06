package usecases

import (
	"context"
	"github.com/Nixonxp/discord/auth/internal/app/models"
)

type UsecaseInterface interface {
	Register(ctx context.Context, registerInfo RegisterUserInfo) (*models.User, error)
	Login(ctx context.Context, loginInfo LoginUserInfo) (*models.LoginResult, error)
	Refresh(ctx context.Context, refreshToken string) (string, error)
	OauthLogin(ctx context.Context, req OauthLoginRequest) (*models.OauthLoginResult, error)
	OauthLoginCallback(ctx context.Context, req OauthLoginCallbackRequest) (*models.LoginResult, error)
}

type UsecaseServiceInterface interface {
	Register(ctx context.Context, registerInfo RegisterUserInfo) (*models.User, error)
	GetUserForLogin(ctx context.Context, loginInfo LoginUserInfo) (*models.User, error)
	CreateOrCreateUser(ctx context.Context, userInfo GetOrCreateUserRequest) (*models.User, error)
}

//go:generate mockery --name=UsersStorage --filename=users_storage_mock.go --disable-version-string
type UsersStorage interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	LoginUser(ctx context.Context, loginInfo *models.Login) (*models.User, error)
}
