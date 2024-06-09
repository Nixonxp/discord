package usecases

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/user/internal/app/models"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	"github.com/Nixonxp/discord/user/internal/app/usecases/mocks"
	log "github.com/Nixonxp/discord/user/pkg/logger"
	"github.com/Nixonxp/discord/user/pkg/transaction_manager"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_usecase_FriendUsecase_GetUserFriends(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		TransactionManager *mocks.TransactionManager
		Log                *log.Logger
		UserRepo           *mocks.UsersStorage
		FriendInvitesRepo  *mocks.FriendInvitesStorage
		UserFriendsRepo    *mocks.UserFriendsStorage
	}

	type args struct {
		ctx context.Context
		req usecases.GetUserFriendsRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.UserFriendsInfo
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetUserFriendsRequest{
					UserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.UserFriendsInfo{
				Friends: []*models.Friend{
					{
						UserId:         models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Login:          "login",
						Name:           "name",
						Email:          "test@test.com",
						AvatarPhotoUrl: "url",
					},
				},
			},
			wantErr: false,

			on: func(f *fields) {
				f.UserFriendsRepo.On("GetUserFriendsByUserId",
					ctx,
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795"))).
					Return([]*models.Friend{
						{
							UserId:         models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
							Login:          "login",
							Name:           "name",
							Email:          "test@test.com",
							AvatarPhotoUrl: "url",
						},
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.UserFriendsRepo.AssertNumberOfCalls(t, "GetUserFriendsByUserId", 1)
			},
		},
		{
			name: "Test 1. Negative. GetUserFriendsByUserId returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetUserFriendsRequest{
					UserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.UserFriendsInfo{
				Friends: []*models.Friend{
					{
						UserId:         models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Login:          "login",
						Name:           "name",
						Email:          "test@test.com",
						AvatarPhotoUrl: "url",
					},
				},
			},
			wantErr:     true,
			errorString: "some error",

			on: func(f *fields) {
				f.UserFriendsRepo.On("GetUserFriendsByUserId",
					ctx,
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795"))).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UserFriendsRepo.AssertNumberOfCalls(t, "GetUserFriendsByUserId", 1)
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
				FriendInvitesRepo:  mocks.NewFriendInvitesStorage(t),
				UserFriendsRepo:    mocks.NewUserFriendsStorage(t),
			}
			au := NewFriendUsecase(Deps{
				TransactionManager: f.TransactionManager,
				Log:                f.Log,
				UserRepo:           f.UserRepo,
				FriendInvitesRepo:  f.FriendInvitesRepo,
				UserFriendsRepo:    f.UserFriendsRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.GetUserFriends(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}

func Test_usecase_FriendUsecase_AddToFriendByUserId(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		TransactionManager *mocks.TransactionManager
		Log                *log.Logger
		UserRepo           *mocks.UsersStorage
		FriendInvitesRepo  *mocks.FriendInvitesStorage
		UserFriendsRepo    *mocks.UserFriendsStorage
	}

	type args struct {
		ctx context.Context
		req usecases.AddToFriendByUserIdRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.ActionInfo
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.AddToFriendByUserIdRequest{
					UserId:  "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					OwnerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b999",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr: false,

			on: func(f *fields) {
				f.UserRepo.On("GetUserById",
					ctx,
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795"))).
					Return(&models.User{
						Id:             models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Login:          "login",
						Name:           "name",
						Email:          "test@test.com",
						AvatarPhotoUrl: "url",
					}, nil)

				f.FriendInvitesRepo.On("CreateInvite",
					ctx,
					mock.MatchedBy(func(invite *models.FriendInvite) bool {
						return invite.InviteId.String() != "" &&
							invite.OwnerId == models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b999")) &&
							invite.UserId == models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) &&
							invite.Status == models.PendingStatus
					})).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "GetUserById", 1)
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "CreateInvite", 1)
			},
		},
		{
			name: "Test 2. Negative. CreateInvite returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.AddToFriendByUserIdRequest{
					UserId:  "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					OwnerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b999",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "create user invite error: some error",

			on: func(f *fields) {
				f.UserRepo.On("GetUserById",
					ctx,
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795"))).
					Return(&models.User{
						Id:             models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Login:          "login",
						Name:           "name",
						Email:          "test@test.com",
						AvatarPhotoUrl: "url",
					}, nil)

				f.FriendInvitesRepo.On("CreateInvite",
					ctx,
					mock.MatchedBy(func(invite *models.FriendInvite) bool {
						return invite.InviteId.String() != "" &&
							invite.OwnerId == models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b999")) &&
							invite.UserId == models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) &&
							invite.Status == models.PendingStatus
					})).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "GetUserById", 1)
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "CreateInvite", 1)
			},
		},
		{
			name: "Test 3. Negative. GetUserById returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.AddToFriendByUserIdRequest{
					UserId:  "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					OwnerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b999",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "user not found: not found",

			on: func(f *fields) {
				f.UserRepo.On("GetUserById",
					ctx,
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795"))).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UserRepo.AssertNumberOfCalls(t, "GetUserById", 1)
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
				FriendInvitesRepo:  mocks.NewFriendInvitesStorage(t),
				UserFriendsRepo:    mocks.NewUserFriendsStorage(t),
			}
			au := NewFriendUsecase(Deps{
				TransactionManager: f.TransactionManager,
				Log:                f.Log,
				UserRepo:           f.UserRepo,
				FriendInvitesRepo:  f.FriendInvitesRepo,
				UserFriendsRepo:    f.UserFriendsRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.AddToFriendByUserId(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}

func Test_usecase_FriendUsecase_GetUserInvites(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		TransactionManager *mocks.TransactionManager
		Log                *log.Logger
		UserRepo           *mocks.UsersStorage
		FriendInvitesRepo  *mocks.FriendInvitesStorage
		UserFriendsRepo    *mocks.UserFriendsStorage
	}

	type args struct {
		ctx context.Context
		req usecases.GetUserInvitesRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.UserInvitesInfo
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetUserInvitesRequest{
					UserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.UserInvitesInfo{
				Invites: []*models.FriendInvite{
					{
						InviteId: models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
						OwnerId:  models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b200")),
						UserId:   models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Status:   models.PendingStatus,
					},
				},
			},
			wantErr: false,

			on: func(f *fields) {
				f.FriendInvitesRepo.On("GetInvitesByUserId",
					ctx,
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795"))).
					Return(&models.UserInvitesInfo{
						Invites: []*models.FriendInvite{
							{
								InviteId: models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
								OwnerId:  models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b200")),
								UserId:   models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
								Status:   models.PendingStatus,
							},
						},
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "GetInvitesByUserId", 1)
			},
		},
		{
			name: "Test 2. Negative. GetInvitesByUserId returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetUserInvitesRequest{
					UserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.UserInvitesInfo{
				Invites: []*models.FriendInvite{
					{
						InviteId: models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
						OwnerId:  models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b200")),
						UserId:   models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Status:   models.PendingStatus,
					},
				},
			},
			wantErr:     true,
			errorString: "get user invites error: some error",

			on: func(f *fields) {
				f.FriendInvitesRepo.On("GetInvitesByUserId",
					ctx,
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795"))).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "GetInvitesByUserId", 1)
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
				FriendInvitesRepo:  mocks.NewFriendInvitesStorage(t),
				UserFriendsRepo:    mocks.NewUserFriendsStorage(t),
			}
			au := NewFriendUsecase(Deps{
				TransactionManager: f.TransactionManager,
				Log:                f.Log,
				UserRepo:           f.UserRepo,
				FriendInvitesRepo:  f.FriendInvitesRepo,
				UserFriendsRepo:    f.UserFriendsRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.GetUserInvites(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}

func Test_usecase_FriendUsecase_AcceptFriendInvite(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		TransactionManager *mocks.TransactionManager
		Log                *log.Logger
		UserRepo           *mocks.UsersStorage
		FriendInvitesRepo  *mocks.FriendInvitesStorage
		UserFriendsRepo    *mocks.UserFriendsStorage
	}

	type args struct {
		ctx context.Context
		req usecases.AcceptFriendInviteRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.ActionInfo
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.AcceptFriendInviteRequest{
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					InviteId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr: false,

			on: func(f *fields) {
				f.FriendInvitesRepo.On("GetInviteById",
					ctx,
					models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.FriendInvite{
						InviteId: models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						OwnerId:  models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b200")),
						UserId:   models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Status:   models.PendingStatus,
					}, nil)

				f.TransactionManager.On("RunReadCommitted",
					ctx,
					transaction_manager.ReadWrite,
					mock.Anything).
					Return(nil)

			},
			assert: func(t *testing.T, f *fields) {
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "GetInviteById", 1)
				f.TransactionManager.AssertNumberOfCalls(t, "RunReadCommitted", 1)
			},
		},
		{
			name: "Test 2. Negative. RunReadCommitted returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.AcceptFriendInviteRequest{
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					InviteId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "accept invite error: some error",

			on: func(f *fields) {
				f.FriendInvitesRepo.On("GetInviteById",
					ctx,
					models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.FriendInvite{
						InviteId: models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						OwnerId:  models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b200")),
						UserId:   models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Status:   models.PendingStatus,
					}, nil)

				f.TransactionManager.On("RunReadCommitted",
					ctx,
					transaction_manager.ReadWrite,
					mock.Anything).
					Return(errors.New("some error"))

			},
			assert: func(t *testing.T, f *fields) {
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "GetInviteById", 1)
				f.TransactionManager.AssertNumberOfCalls(t, "RunReadCommitted", 1)
			},
		},
		{
			name: "Test 3. Negative. GetInviteById returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.AcceptFriendInviteRequest{
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					InviteId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "some error",

			on: func(f *fields) {
				f.FriendInvitesRepo.On("GetInviteById",
					ctx,
					models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "GetInviteById", 1)
			},
		},
		{
			name: "Test 4. Negative. GetInviteById returns not found error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.AcceptFriendInviteRequest{
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					InviteId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "invite : not found",

			on: func(f *fields) {
				f.FriendInvitesRepo.On("GetInviteById",
					ctx,
					models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(nil, models.ErrNotFound)
			},
			assert: func(t *testing.T, f *fields) {
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "GetInviteById", 1)
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
				FriendInvitesRepo:  mocks.NewFriendInvitesStorage(t),
				UserFriendsRepo:    mocks.NewUserFriendsStorage(t),
			}
			au := NewFriendUsecase(Deps{
				TransactionManager: f.TransactionManager,
				Log:                f.Log,
				UserRepo:           f.UserRepo,
				FriendInvitesRepo:  f.FriendInvitesRepo,
				UserFriendsRepo:    f.UserFriendsRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.AcceptFriendInvite(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}

func Test_usecase_FriendUsecase_DeclineFriendInvite(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		TransactionManager *mocks.TransactionManager
		Log                *log.Logger
		UserRepo           *mocks.UsersStorage
		FriendInvitesRepo  *mocks.FriendInvitesStorage
		UserFriendsRepo    *mocks.UserFriendsStorage
	}

	type args struct {
		ctx context.Context
		req usecases.DeclineFriendInviteRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.ActionInfo
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.DeclineFriendInviteRequest{
					InviteId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr: false,

			on: func(f *fields) {
				f.FriendInvitesRepo.On("GetInviteById",
					ctx,
					models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.FriendInvite{
						InviteId: models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						OwnerId:  models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b200")),
						UserId:   models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Status:   models.PendingStatus,
					}, nil)

				f.FriendInvitesRepo.On("DeclineInvite",
					ctx,
					models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "GetInviteById", 1)
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "DeclineInvite", 1)
			},
		},
		{
			name: "Test 2. Negative. DeclineInvite returns an error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.DeclineFriendInviteRequest{
					InviteId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "invite decline error: some error",

			on: func(f *fields) {
				f.FriendInvitesRepo.On("GetInviteById",
					ctx,
					models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.FriendInvite{
						InviteId: models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						OwnerId:  models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b200")),
						UserId:   models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						Status:   models.PendingStatus,
					}, nil)

				f.FriendInvitesRepo.On("DeclineInvite",
					ctx,
					models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "GetInviteById", 1)
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "DeclineInvite", 1)
			},
		},
		{
			name: "Test 3. Negative. permission error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.DeclineFriendInviteRequest{
					InviteId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "invite permission denied: permission denied",

			on: func(f *fields) {
				f.FriendInvitesRepo.On("GetInviteById",
					ctx,
					models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.FriendInvite{
						InviteId: models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						OwnerId:  models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b200")),
						UserId:   models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b791")),
						Status:   models.PendingStatus,
					}, nil)

			},
			assert: func(t *testing.T, f *fields) {
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "GetInviteById", 1)
			},
		},
		{
			name: "Test 4. Negative. GetInviteById returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.DeclineFriendInviteRequest{
					InviteId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "some error",

			on: func(f *fields) {
				f.FriendInvitesRepo.On("GetInviteById",
					ctx,
					models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(nil, errors.New("some error"))

			},
			assert: func(t *testing.T, f *fields) {
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "GetInviteById", 1)
			},
		},
		{
			name: "Test 5. Negative. GetInviteById returns not found error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.DeclineFriendInviteRequest{
					InviteId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "invite: not found",

			on: func(f *fields) {
				f.FriendInvitesRepo.On("GetInviteById",
					ctx,
					models.InviteId(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(nil, models.ErrNotFound)

			},
			assert: func(t *testing.T, f *fields) {
				f.FriendInvitesRepo.AssertNumberOfCalls(t, "GetInviteById", 1)
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
				FriendInvitesRepo:  mocks.NewFriendInvitesStorage(t),
				UserFriendsRepo:    mocks.NewUserFriendsStorage(t),
			}
			au := NewFriendUsecase(Deps{
				TransactionManager: f.TransactionManager,
				Log:                f.Log,
				UserRepo:           f.UserRepo,
				FriendInvitesRepo:  f.FriendInvitesRepo,
				UserFriendsRepo:    f.UserFriendsRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.DeclineFriendInvite(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}

func Test_usecase_FriendUsecase_DeleteFromFriend(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		TransactionManager *mocks.TransactionManager
		Log                *log.Logger
		UserRepo           *mocks.UsersStorage
		FriendInvitesRepo  *mocks.FriendInvitesStorage
		UserFriendsRepo    *mocks.UserFriendsStorage
	}

	type args struct {
		ctx context.Context
		req usecases.DeleteFromFriendRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.ActionInfo
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.DeleteFromFriendRequest{
					FriendId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr: false,

			on: func(f *fields) {
				f.UserFriendsRepo.On("DeleteFriend",
					ctx,
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
				).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.UserFriendsRepo.AssertNumberOfCalls(t, "DeleteFriend", 1)
			},
		},
		{
			name: "Test 2. Negative. DeleteFriend returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.DeleteFromFriendRequest{
					FriendId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "delete froms friend error: some  errors",

			on: func(f *fields) {
				f.UserFriendsRepo.On("DeleteFriend",
					ctx,
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
				).
					Return(errors.New("some  errors"))
			},
			assert: func(t *testing.T, f *fields) {
				f.UserFriendsRepo.AssertNumberOfCalls(t, "DeleteFriend", 1)
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
				FriendInvitesRepo:  mocks.NewFriendInvitesStorage(t),
				UserFriendsRepo:    mocks.NewUserFriendsStorage(t),
			}
			au := NewFriendUsecase(Deps{
				TransactionManager: f.TransactionManager,
				Log:                f.Log,
				UserRepo:           f.UserRepo,
				FriendInvitesRepo:  f.FriendInvitesRepo,
				UserFriendsRepo:    f.UserFriendsRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.DeleteFromFriend(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}
