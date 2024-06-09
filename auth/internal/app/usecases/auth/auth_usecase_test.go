package usecases

import (
	"context"
	"errors"
	config "github.com/Nixonxp/discord/auth/configs"
	"github.com/Nixonxp/discord/auth/internal/app/models"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	"github.com/Nixonxp/discord/auth/internal/app/usecases/mocks"
	log "github.com/Nixonxp/discord/auth/pkg/logger"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_usecase_AuthUsecase_Register(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		UsersService *mocks.UserServiceInterface
		Log          *log.Logger
		OauthSvc     usecases.OAuthServiceInterface
		Cfg          *config.Config
	}

	type args struct {
		ctx  context.Context
		info usecases.RegisterUserInfo
	}
	tests := []struct {
		name    string
		args    args
		want    *models.User
		wantErr bool

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				info: usecases.RegisterUserInfo{
					Login:    "login",
					Name:     "name",
					Email:    "test@test.ru",
					Password: "pass",
				},
			},
			want: &models.User{
				Login:    "login",
				Name:     "name",
				Email:    "test@test.ru",
				Password: "pass",
			},
			wantErr: false,

			on: func(f *fields) {
				f.UsersService.On("Register", ctx, mock.MatchedBy(func(user usecases.RegisterUserInfo) bool {
					return user.Name == "name" &&
						user.Login == "login" &&
						user.Email == "test@test.ru" &&
						user.Password != "pass"
				})).
					Return(&models.User{
						UserID:   models.UserID{},
						Login:    "login",
						Name:     "name",
						Email:    "test@test.ru",
						Password: "pass",
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.UsersService.AssertNumberOfCalls(t, "Register", 1)
			},
		},
		{
			name: "Test 2. Negative. Register return error.",
			args: args{
				ctx: ctx, // dumm
				info: usecases.RegisterUserInfo{
					Login:    "login",
					Name:     "name",
					Email:    "test@test.ru",
					Password: "pass",
				},
			},
			want:    nil,
			wantErr: true,

			on: func(f *fields) {
				f.UsersService.On("Register", ctx, mock.MatchedBy(func(user usecases.RegisterUserInfo) bool {
					return user.Name == "name" &&
						user.Login == "login" &&
						user.Email == "test@test.ru" &&
						user.Password != "pass"
				})).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UsersService.AssertNumberOfCalls(t, "Register", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				UsersService: mocks.NewUserServiceInterface(t),
				Log:          &log.Logger{},
				OauthSvc:     mocks.NewOAuthServiceInterface(t),
				Cfg:          config.NewConfig(),
			}
			au := NewAuthUsecase(Deps{
				UserService: f.UsersService,
				Log:         f.Log,
				OauthSvc:    f.OauthSvc,
				Cfg:         f.Cfg,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.Register(tt.args.ctx, tt.args.info)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != nil { // зануляем так как не можем проверить
				got.UserID = models.UserID{}
			}
			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}

func Test_usecase_AuthUsecase_Login(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		UsersService *mocks.UserServiceInterface
		Log          *log.Logger
		OauthSvc     usecases.OAuthServiceInterface
		Cfg          *config.Config
	}

	type args struct {
		ctx  context.Context
		info usecases.LoginUserInfo
	}
	tests := []struct {
		name        string
		args        args
		want        *models.LoginResult
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				info: usecases.LoginUserInfo{
					Login:    "login",
					Password: "password123",
				},
			},
			want: &models.LoginResult{
				Token:        "token",
				RefreshToken: "refresh_token",
			},
			wantErr: false,

			on: func(f *fields) {
				f.UsersService.On("GetUserForLogin", ctx, usecases.LoginUserInfo{
					Login:    "login",
					Password: "password123",
				}).
					Return(&models.User{
						UserID:   models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Login:    "login",
						Name:     "name",
						Email:    "test@test.ru",
						Password: "$2a$04$cO03pCxW.xurFRFGAobYbOeeU2WqEbaqAuUy2I5o8w4XTCq3OelCq",
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.UsersService.AssertNumberOfCalls(t, "GetUserForLogin", 1)
			},
		},
		{
			name: "Test 2. Negative. GetUserForLogin return error.",
			args: args{
				ctx: ctx, // dumm
				info: usecases.LoginUserInfo{
					Login:    "login",
					Password: "pass",
				},
			},
			want:        nil,
			wantErr:     true,
			errorString: "some error",

			on: func(f *fields) {
				f.UsersService.On("GetUserForLogin", ctx, usecases.LoginUserInfo{
					Login:    "login",
					Password: "pass",
				}).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UsersService.AssertNumberOfCalls(t, "GetUserForLogin", 1)
			},
		},
		{
			name: "Test 3. Negative wrong credentials.",
			args: args{
				ctx: ctx, // dumm
				info: usecases.LoginUserInfo{
					Login:    "login",
					Password: "password123",
				},
			},
			want: &models.LoginResult{
				Token:        "token",
				RefreshToken: "refresh_token",
			},
			wantErr:     true,
			errorString: "wrong password or login: credentials invalid",

			on: func(f *fields) {
				f.UsersService.On("GetUserForLogin", ctx, usecases.LoginUserInfo{
					Login:    "login",
					Password: "password123",
				}).
					Return(&models.User{
						UserID:   models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Login:    "login",
						Name:     "name",
						Email:    "test@test.ru",
						Password: "$2a$04$cO03pCxW.xurFRFGAobYbOeeU2WqEbaqAuUy2I5o8w4XTCq3OelCs",
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.UsersService.AssertNumberOfCalls(t, "GetUserForLogin", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			cfg := config.NewConfig()
			cfg.Application.AuthSecretKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
			cfg.Application.AuthRefreshSecretKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9Refresh"

			f := &fields{
				UsersService: mocks.NewUserServiceInterface(t),
				Log:          &log.Logger{},
				OauthSvc:     mocks.NewOAuthServiceInterface(t),
				Cfg:          cfg,
			}
			au := NewAuthUsecase(Deps{
				UserService: f.UsersService,
				Log:         f.Log,
				OauthSvc:    f.OauthSvc,
				Cfg:         f.Cfg,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.Login(tt.args.ctx, tt.args.info)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			if got != nil { // мокаем так как не можем проверить
				got.Token = "token"
				got.RefreshToken = "refresh_token"
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}

func Test_usecase_AuthUsecase_Refresh_Positive(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)

	type fields struct {
		UsersService *mocks.UserServiceInterface
		Log          *log.Logger
		OauthSvc     usecases.OAuthServiceInterface
		Cfg          *config.Config
	}

	cfg := config.NewConfig()
	cfg.Application.AuthSecretKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"
	cfg.Application.AuthRefreshSecretKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9Refresh"

	refreshToken, err := generateRefreshJWTToken(cfg.Application.AuthRefreshSecretKey, "284fef68-7e3e-4d1d-96a0-8c96f7b3b795")
	assert.NoError(t, err)

	f := &fields{
		UsersService: mocks.NewUserServiceInterface(t),
		Log:          &log.Logger{},
		OauthSvc:     mocks.NewOAuthServiceInterface(t),
		Cfg:          cfg,
	}

	au := NewAuthUsecase(Deps{
		UserService: f.UsersService,
		Log:         f.Log,
		OauthSvc:    f.OauthSvc,
		Cfg:         f.Cfg,
	})

	// act
	got, err := au.Refresh(ctx, refreshToken)
	assert.NoError(t, err)
	assert.NotEqual(t, "", got)
}

func Test_usecase_AuthUsecase_Refresh_Negative(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)

	type fields struct {
		UsersService *mocks.UserServiceInterface
		Log          *log.Logger
		OauthSvc     usecases.OAuthServiceInterface
		Cfg          *config.Config
	}

	cfg := config.NewConfig()
	cfg.Application.AuthRefreshSecretKey = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9Refresh"

	// pastime add
	expirationTime := time.Now().Add(-24 * time.Hour)
	claims := jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Id:        "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(cfg.Application.AuthRefreshSecretKey))
	assert.NoError(t, err)

	f := &fields{
		UsersService: mocks.NewUserServiceInterface(t),
		Log:          &log.Logger{},
		OauthSvc:     mocks.NewOAuthServiceInterface(t),
		Cfg:          cfg,
	}

	au := NewAuthUsecase(Deps{
		UserService: f.UsersService,
		Log:         f.Log,
		OauthSvc:    f.OauthSvc,
		Cfg:         f.Cfg,
	})

	// act
	got, err := au.Refresh(ctx, refreshToken)
	assert.Error(t, err)
	assert.Equal(t, "", got)
}

func Test_usecase_AuthUsecase_OauthLogin(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		UsersService *mocks.UserServiceInterface
		Log          *log.Logger
		OauthSvc     *mocks.OAuthServiceInterface
		Cfg          *config.Config
	}

	type args struct {
		ctx context.Context
		req usecases.OauthLoginRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *models.OauthLoginResult
		wantErr bool

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.OauthLoginRequest{},
			},
			want: &models.OauthLoginResult{
				Code: "url",
			},
			wantErr: false,

			on: func(f *fields) {
				f.OauthSvc.On("AuthCodeURL", "state").
					Return("url")
			},
			assert: func(t *testing.T, f *fields) {
				f.OauthSvc.AssertNumberOfCalls(t, "AuthCodeURL", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				UsersService: mocks.NewUserServiceInterface(t),
				Log:          &log.Logger{},
				OauthSvc:     mocks.NewOAuthServiceInterface(t),
				Cfg:          config.NewConfig(),
			}
			au := NewAuthUsecase(Deps{
				UserService: f.UsersService,
				Log:         f.Log,
				OauthSvc:    f.OauthSvc,
				Cfg:         f.Cfg,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.OauthLogin(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}

func Test_usecase_AuthUsecase_OauthLoginCallback(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		UsersService *mocks.UserServiceInterface
		Log          *log.Logger
		OauthSvc     *mocks.OAuthServiceInterface
		Cfg          *config.Config
	}

	type args struct {
		ctx context.Context
		req usecases.OauthLoginCallbackRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.LoginResult
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.OauthLoginCallbackRequest{
					Code:  "code",
					State: "state",
				},
			},
			want: &models.LoginResult{
				Token:        "token",
				RefreshToken: "refresh_token",
			},
			wantErr: false,

			on: func(f *fields) {
				f.OauthSvc.On("ExchangeClient", ctx, "code").
					Return(&usecases.UserInfo{
						OauthId:        "123",
						Name:           "name",
						Email:          "test@test.com",
						AvatarPhotoUrl: "url",
					}, nil)

				f.UsersService.On("CreateOrCreateUser", ctx, mock.MatchedBy(func(user usecases.GetOrCreateUserRequest) bool {
					return user.Name == "name" &&
						user.Login == "test@test.com" &&
						user.Email == "test@test.com" &&
						user.AvatarPhotoUrl == "url" &&
						user.OauthId == "123" &&
						user.Password != "pass"
				})).
					Return(&models.User{
						UserID:   models.UserID{},
						Login:    "login",
						Name:     "name",
						Email:    "test@test.ru",
						Password: "pass",
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.OauthSvc.AssertNumberOfCalls(t, "ExchangeClient", 1)
			},
		},
		{
			name: "Test 2. Negative. ExchangeClient return error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.OauthLoginCallbackRequest{
					Code:  "code",
					State: "state",
				},
			},
			want: &models.LoginResult{
				Token:        "token",
				RefreshToken: "refresh_token",
			},
			wantErr:     true,
			errorString: "some error",

			on: func(f *fields) {
				f.OauthSvc.On("ExchangeClient", ctx, "code").
					Return(&usecases.UserInfo{
						OauthId:        "123",
						Name:           "name",
						Email:          "test@test.com",
						AvatarPhotoUrl: "url",
					}, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.OauthSvc.AssertNumberOfCalls(t, "ExchangeClient", 1)
			},
		},
		{
			name: "Test 3. Negative. CreateOrCreateUser return error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.OauthLoginCallbackRequest{
					Code:  "code",
					State: "state",
				},
			},
			want: &models.LoginResult{
				Token:        "token",
				RefreshToken: "refresh_token",
			},
			wantErr:     true,
			errorString: "some error",

			on: func(f *fields) {
				f.OauthSvc.On("ExchangeClient", ctx, "code").
					Return(&usecases.UserInfo{
						OauthId:        "123",
						Name:           "name",
						Email:          "test@test.com",
						AvatarPhotoUrl: "url",
					}, nil)

				f.UsersService.On("CreateOrCreateUser", ctx, mock.MatchedBy(func(user usecases.GetOrCreateUserRequest) bool {
					return user.Name == "name" &&
						user.Login == "test@test.com" &&
						user.Email == "test@test.com" &&
						user.AvatarPhotoUrl == "url" &&
						user.OauthId == "123" &&
						user.Password != "pass"
				})).
					Return(&models.User{}, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.OauthSvc.AssertNumberOfCalls(t, "ExchangeClient", 1)
				f.OauthSvc.AssertNumberOfCalls(t, "CreateOrCreateUser", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				UsersService: mocks.NewUserServiceInterface(t),
				Log:          &log.Logger{},
				OauthSvc:     mocks.NewOAuthServiceInterface(t),
				Cfg:          config.NewConfig(),
			}
			au := NewAuthUsecase(Deps{
				UserService: f.UsersService,
				Log:         f.Log,
				OauthSvc:    f.OauthSvc,
				Cfg:         f.Cfg,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.OauthLoginCallback(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			if got != nil { // мокаем так как не можем проверить
				got.Token = "token"
				got.RefreshToken = "refresh_token"
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}
