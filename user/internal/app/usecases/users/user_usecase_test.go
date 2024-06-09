package usecases

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/user/internal/app/models"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	"github.com/Nixonxp/discord/user/internal/app/usecases/mocks"
	log "github.com/Nixonxp/discord/user/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_usecase_UserUsecase_CreateUser(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		TransactionManager *mocks.TransactionManager
		Log                *log.Logger
		UserRepo           *mocks.UsersStorage
	}

	type args struct {
		ctx context.Context
		req usecases.CreateUserRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.User
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.CreateUserRequest{
					Login:    "login",
					Name:     "name",
					Email:    "test@test.com",
					Password: "pass",
				},
			},
			want: &models.User{
				Id:             models.UserID{},
				Login:          "login",
				Name:           "name",
				Email:          "test@test.com",
				AvatarPhotoUrl: "",
				Password:       "pass",
			},
			wantErr: false,

			on: func(f *fields) {
				f.UserRepo.On("CreateUser",
					ctx,
					mock.MatchedBy(func(user *models.User) bool {
						return user.Id.String() != "" &&
							user.Login == "login" &&
							user.Name == "name" &&
							user.AvatarPhotoUrl == "" &&
							user.Password == "pass" &&
							user.Email == "test@test.com"
					})).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "CreateUser", 1)
			},
		},
		{
			name: "Test 2. Negative. CreateUser returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.CreateUserRequest{
					Login:    "login",
					Name:     "name",
					Email:    "test@test.com",
					Password: "pass",
				},
			},
			want: &models.User{
				Id:             models.UserID{},
				Login:          "login",
				Name:           "name",
				Email:          "test@test.com",
				AvatarPhotoUrl: "",
				Password:       "pass",
			},
			wantErr:     true,
			errorString: "create user error: some error",

			on: func(f *fields) {
				f.UserRepo.On("CreateUser",
					ctx,
					mock.MatchedBy(func(user *models.User) bool {
						return user.Id.String() != "" &&
							user.Login == "login" &&
							user.Name == "name" &&
							user.AvatarPhotoUrl == "" &&
							user.Password == "pass" &&
							user.Email == "test@test.com"
					})).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "CreateUser", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				TransactionManager: mocks.NewTransactionManager(t),
				Log:                &log.Logger{},
				UserRepo:           mocks.NewUsersStorage(t),
			}
			au := NewUserUsecase(Deps{
				TransactionManager: f.TransactionManager,
				Log:                f.Log,
				UserRepo:           f.UserRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.CreateUser(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			if got != nil { // зануляем так как не можем проверить
				got.Id = models.UserID{}
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}

func Test_usecase_UserUsecase_CreateOrGetUser(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		TransactionManager *mocks.TransactionManager
		Log                *log.Logger
		UserRepo           *mocks.UsersStorage
	}

	type args struct {
		ctx context.Context
		req usecases.CreateOrGetUserRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.User
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.CreateOrGetUserRequest{
					Login:    "login",
					Name:     "name",
					Email:    "test@test.com",
					Password: "pass",
					OauthId:  "123",
				},
			},
			want: &models.User{
				Id:             models.UserID{},
				Login:          "login",
				Name:           "name",
				Email:          "test@test.com",
				AvatarPhotoUrl: "url",
				Password:       "pass",
			},
			wantErr: false,

			on: func(f *fields) {
				f.UserRepo.On("GetUserByOauthId",
					ctx,
					"123").
					Return(&models.User{
						Id:             models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Login:          "login",
						Name:           "name",
						Email:          "test@test.com",
						AvatarPhotoUrl: "url",
						Password:       "pass",
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "GetUserByOauthId", 1)
			},
		},
		{
			name: "Test 2. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.CreateOrGetUserRequest{
					Login:    "login",
					Name:     "name",
					Email:    "test@test.com",
					Password: "pass",
					OauthId:  "123",
				},
			},
			want: &models.User{
				Id:             models.UserID{},
				Login:          "login",
				Name:           "name",
				Email:          "test@test.com",
				AvatarPhotoUrl: "",
				Password:       "pass",
				OauthId:        "123",
			},
			wantErr: false,

			on: func(f *fields) {
				f.UserRepo.On("GetUserByOauthId",
					ctx,
					"123").
					Return(nil, models.ErrNotFound)

				f.UserRepo.On("CreateUser",
					ctx,
					mock.MatchedBy(func(user *models.User) bool {
						return user.Id.String() != "" &&
							user.Login == "login" &&
							user.Name == "name" &&
							user.AvatarPhotoUrl == "" &&
							user.Password == "pass" &&
							user.Email == "test@test.com"
					})).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "GetUserByOauthId", 1)
				f.UserRepo.AssertNumberOfCalls(t, "CreateUser", 1)
			},
		},
		{
			name: "Test 3. Negative. CreateUser returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.CreateOrGetUserRequest{
					Login:    "login",
					Name:     "name",
					Email:    "test@test.com",
					Password: "pass",
					OauthId:  "123",
				},
			},
			want: &models.User{
				Id:             models.UserID{},
				Login:          "login",
				Name:           "name",
				Email:          "test@test.com",
				AvatarPhotoUrl: "",
				Password:       "pass",
				OauthId:        "123",
			},
			wantErr:     true,
			errorString: "user create error: some error",

			on: func(f *fields) {
				f.UserRepo.On("GetUserByOauthId",
					ctx,
					"123").
					Return(nil, models.ErrNotFound)

				f.UserRepo.On("CreateUser",
					ctx,
					mock.MatchedBy(func(user *models.User) bool {
						return user.Id.String() != "" &&
							user.Login == "login" &&
							user.Name == "name" &&
							user.AvatarPhotoUrl == "" &&
							user.Password == "pass" &&
							user.Email == "test@test.com"
					})).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "GetUserByOauthId", 1)
				f.UserRepo.AssertNumberOfCalls(t, "CreateUser", 1)
			},
		},
		{
			name: "Test 4. Negative. GetUserByOauthId returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.CreateOrGetUserRequest{
					Login:    "login",
					Name:     "name",
					Email:    "test@test.com",
					Password: "pass",
					OauthId:  "123",
				},
			},
			want: &models.User{
				Id:             models.UserID{},
				Login:          "login",
				Name:           "name",
				Email:          "test@test.com",
				AvatarPhotoUrl: "",
				Password:       "pass",
				OauthId:        "123",
			},
			wantErr:     true,
			errorString: "user not found: not found",

			on: func(f *fields) {
				f.UserRepo.On("GetUserByOauthId",
					ctx,
					"123").
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "GetUserByOauthId", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				TransactionManager: mocks.NewTransactionManager(t),
				Log:                &log.Logger{},
				UserRepo:           mocks.NewUsersStorage(t),
			}
			au := NewUserUsecase(Deps{
				TransactionManager: f.TransactionManager,
				Log:                f.Log,
				UserRepo:           f.UserRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.CreateOrGetUser(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			if got != nil { // зануляем так как не можем проверить
				got.Id = models.UserID{}
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}

func Test_usecase_UserUsecase_UpdateUser(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		TransactionManager *mocks.TransactionManager
		Log                *log.Logger
		UserRepo           *mocks.UsersStorage
	}

	type args struct {
		ctx context.Context
		req usecases.UpdateUserRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.User
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.UpdateUserRequest{
					Id:             "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					Login:          "login",
					Name:           "name",
					Email:          "test@test.com",
					AvatarPhotoUrl: "url",
					CurrentUserId:  "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.User{
				Id:             models.UserID{},
				Login:          "login",
				Name:           "name",
				Email:          "test@test.com",
				AvatarPhotoUrl: "url",
			},
			wantErr: false,

			on: func(f *fields) {
				f.UserRepo.On("UpdateUser",
					ctx,
					&models.User{
						Id:             models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Login:          "login",
						Name:           "name",
						Email:          "test@test.com",
						AvatarPhotoUrl: "url",
					}).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "UpdateUser", 1)
			},
		},
		{
			name: "Test 2. Negative. UpdateUser returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.UpdateUserRequest{
					Id:             "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					Login:          "login",
					Name:           "name",
					Email:          "test@test.com",
					AvatarPhotoUrl: "url",
					CurrentUserId:  "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.User{
				Id:             models.UserID{},
				Login:          "login",
				Name:           "name",
				Email:          "test@test.com",
				AvatarPhotoUrl: "url",
			},
			wantErr:     true,
			errorString: "some error",

			on: func(f *fields) {
				f.UserRepo.On("UpdateUser",
					ctx,
					&models.User{
						Id:             models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Login:          "login",
						Name:           "name",
						Email:          "test@test.com",
						AvatarPhotoUrl: "url",
					}).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "UpdateUser", 1)
			},
		},
		{
			name: "Test 3. Negative. UpdateUser returns not found error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.UpdateUserRequest{
					Id:             "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					Login:          "login",
					Name:           "name",
					Email:          "test@test.com",
					AvatarPhotoUrl: "url",
					CurrentUserId:  "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.User{
				Id:             models.UserID{},
				Login:          "login",
				Name:           "name",
				Email:          "test@test.com",
				AvatarPhotoUrl: "url",
			},
			wantErr:     true,
			errorString: "user for update not found: not found",

			on: func(f *fields) {
				f.UserRepo.On("UpdateUser",
					ctx,
					&models.User{
						Id:             models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Login:          "login",
						Name:           "name",
						Email:          "test@test.com",
						AvatarPhotoUrl: "url",
					}).
					Return(models.ErrNotFound)
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "UpdateUser", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				TransactionManager: mocks.NewTransactionManager(t),
				Log:                &log.Logger{},
				UserRepo:           mocks.NewUsersStorage(t),
			}
			au := NewUserUsecase(Deps{
				TransactionManager: f.TransactionManager,
				Log:                f.Log,
				UserRepo:           f.UserRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.UpdateUser(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			if got != nil { // зануляем так как не можем проверить
				got.Id = models.UserID{}
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}

func Test_usecase_UserUsecase_GetUserForLogin(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		TransactionManager *mocks.TransactionManager
		Log                *log.Logger
		UserRepo           *mocks.UsersStorage
	}

	type args struct {
		ctx context.Context
		req usecases.GetUserByLoginRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.User
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetUserByLoginRequest{
					Login: "login",
				},
			},
			want: &models.User{
				Id:             models.UserID{},
				Login:          "login",
				Name:           "name",
				Email:          "test@test.com",
				AvatarPhotoUrl: "url",
			},
			wantErr: false,

			on: func(f *fields) {
				f.UserRepo.On("GetUserByLogin",
					ctx,
					"login").
					Return(&models.User{
						Id:             models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Login:          "login",
						Name:           "name",
						Email:          "test@test.com",
						AvatarPhotoUrl: "url",
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "GetUserByLogin", 1)
			},
		},
		{
			name: "Test 2. Negative. GetUserByLogin returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetUserByLoginRequest{
					Login: "login",
				},
			},
			want: &models.User{
				Id:             models.UserID{},
				Login:          "login",
				Name:           "name",
				Email:          "test@test.com",
				AvatarPhotoUrl: "url",
			},
			wantErr:     true,
			errorString: "some error",

			on: func(f *fields) {
				f.UserRepo.On("GetUserByLogin",
					ctx,
					"login").
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "GetUserByLogin", 1)
			},
		},
		{
			name: "Test 3. Negative. GetUserByLogin returns not found error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetUserByLoginRequest{
					Login: "login",
				},
			},
			want: &models.User{
				Id:             models.UserID{},
				Login:          "login",
				Name:           "name",
				Email:          "test@test.com",
				AvatarPhotoUrl: "url",
			},
			wantErr:     true,
			errorString: "user by login not found: not found",

			on: func(f *fields) {
				f.UserRepo.On("GetUserByLogin",
					ctx,
					"login").
					Return(nil, models.ErrNotFound)
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "GetUserByLogin", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				TransactionManager: mocks.NewTransactionManager(t),
				Log:                &log.Logger{},
				UserRepo:           mocks.NewUsersStorage(t),
			}
			au := NewUserUsecase(Deps{
				TransactionManager: f.TransactionManager,
				Log:                f.Log,
				UserRepo:           f.UserRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.GetUserForLogin(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			if got != nil { // зануляем так как не можем проверить
				got.Id = models.UserID{}
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}
