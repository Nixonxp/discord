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

//go:generate mockery --name=UserServiceInterface --filename=users_service_mock.go --disable-version-string
type UserServiceInterface interface {
	Register(ctx context.Context, registerInfo RegisterUserInfo) (*models.User, error)
	GetUserForLogin(ctx context.Context, loginInfo LoginUserInfo) (*models.User, error)
	CreateOrCreateUser(ctx context.Context, userInfo GetOrCreateUserRequest) (*models.User, error)
}

//go:generate mockery --name=OAuthServiceInterface --filename=oauth_service_mock.go --disable-version-string
type OAuthServiceInterface interface {
	AuthCodeURL(state string) string
	ExchangeClient(ctx context.Context, code string) (*UserInfo, error)
}
