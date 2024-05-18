package usecases

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/auth/internal/app/models"
	"github.com/Nixonxp/discord/auth/internal/app/usecases/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_usecase_AuthUsecase_Register(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		UsersStorage *mocks.UsersStorage
	}

	type args struct {
		ctx  context.Context
		info RegisterUserInfo
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
				info: RegisterUserInfo{
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
				f.UsersStorage.On("CreateUser", ctx, mock.MatchedBy(func(user *models.User) bool {
					return user != nil &&
						user.Name == "name" &&
						user.Login == "login" &&
						user.Email == "test@test.ru" &&
						user.Password == "pass" &&
						user.UserID != models.UserID{} // not empty
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
				f.UsersStorage.AssertNumberOfCalls(t, "CreateUser", 1)
			},
		},
		{
			name: "Test 2. Negative. CreateUser return error.",
			args: args{
				ctx: ctx, // dumm
				info: RegisterUserInfo{
					Login:    "login",
					Name:     "name",
					Email:    "test@test.ru",
					Password: "pass",
				},
			},
			want:    nil,
			wantErr: true,

			on: func(f *fields) {
				f.UsersStorage.On("CreateUser", ctx, mock.MatchedBy(func(user *models.User) bool {
					return user != nil &&
						user.Name == "name" &&
						user.Login == "login" &&
						user.Email == "test@test.ru" &&
						user.Password == "pass" &&
						user.UserID != models.UserID{} // not empty
				})).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UsersStorage.AssertNumberOfCalls(t, "CreateUser", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				UsersStorage: mocks.NewUsersStorage(t),
			}
			au := NewAuthUsecase(Deps{
				UserRepo: f.UsersStorage,
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
		UsersStorage *mocks.UsersStorage
	}

	type args struct {
		ctx  context.Context
		info LoginUserInfo
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
				info: LoginUserInfo{
					Login:    "login",
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
				f.UsersStorage.On("LoginUser", ctx, &models.Login{
					Login:    "login",
					Password: "pass",
				}).
					Return(&models.User{
						UserID:   models.UserID{},
						Login:    "login",
						Name:     "name",
						Email:    "test@test.ru",
						Password: "pass",
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.UsersStorage.AssertNumberOfCalls(t, "LoginUser", 1)
			},
		},
		{
			name: "Test 2. Negative. LoginUser return error.",
			args: args{
				ctx: ctx, // dumm
				info: LoginUserInfo{
					Login:    "login",
					Password: "pass",
				},
			},
			want:    nil,
			wantErr: true,

			on: func(f *fields) {
				f.UsersStorage.On("LoginUser", ctx, &models.Login{
					Login:    "login",
					Password: "pass",
				}).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UsersStorage.AssertNumberOfCalls(t, "LoginUser", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				UsersStorage: mocks.NewUsersStorage(t),
			}
			au := NewAuthUsecase(Deps{
				UserRepo: f.UsersStorage,
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
