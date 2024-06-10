package channel

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/channel/internal/app/models"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
	"github.com/Nixonxp/discord/channel/internal/app/usecases/mocks"
	log "github.com/Nixonxp/discord/channel/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func Test_usecase_ChannelUsecase_AddChannel(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		ChannelRepo   *mocks.ChannelStorage
		Log           *log.Logger
		SubscribeRepo *mocks.SubscribeStorage
	}

	type args struct {
		ctx context.Context
		req usecases.AddChannelRequest
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
				req: usecases.AddChannelRequest{
					Name:          "channel name",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr: false,

			on: func(f *fields) {
				f.ChannelRepo.On("CreateChannel",
					ctx,
					mock.MatchedBy(func(channel models.Channel) bool {
						return channel.Id.String() != "" &&
							channel.Name == "channel name" &&
							channel.OwnerId.String() == "284fef68-7e3e-4d1d-96a0-8c96f7b3b795"
					})).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChannelRepo.AssertNumberOfCalls(t, "CreateChannel", 1)
			},
		},
		{
			name: "Test 2. Negative. CreateChannel returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.AddChannelRequest{
					Name:          "channel name",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "create channel: some error",

			on: func(f *fields) {
				f.ChannelRepo.On("CreateChannel",
					ctx,
					mock.MatchedBy(func(channel models.Channel) bool {
						return channel.Id.String() != "" &&
							channel.Name == "channel name" &&
							channel.OwnerId.String() == "284fef68-7e3e-4d1d-96a0-8c96f7b3b795"
					})).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChannelRepo.AssertNumberOfCalls(t, "CreateChannel", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ChannelRepo:   mocks.NewChannelStorage(t),
				Log:           &log.Logger{},
				SubscribeRepo: mocks.NewSubscribeStorage(t),
			}
			au := NewChannelUsecase(Deps{
				ChannelRepo:   f.ChannelRepo,
				Log:           f.Log,
				SubscribeRepo: f.SubscribeRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.AddChannel(tt.args.ctx, tt.args.req)

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

func Test_usecase_ChannelUsecase_DeleteChannel(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		ChannelRepo   *mocks.ChannelStorage
		Log           *log.Logger
		SubscribeRepo *mocks.SubscribeStorage
	}

	type args struct {
		ctx context.Context
		req usecases.DeleteChannelRequest
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
				req: usecases.DeleteChannelRequest{
					ChannelId:     "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr: false,

			on: func(f *fields) {
				f.ChannelRepo.On("GetChannelById",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.Channel{
						Id:      models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						Name:    "name",
						OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
					}, nil)

				f.ChannelRepo.On("DeleteChannel",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChannelRepo.AssertNumberOfCalls(t, "GetChannelById", 1)
				f.ChannelRepo.AssertNumberOfCalls(t, "DeleteChannel", 1)
			},
		},
		{
			name: "Test 2. Negative. DeleteChannel returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.DeleteChannelRequest{
					ChannelId:     "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "delete channel error: some error",

			on: func(f *fields) {
				f.ChannelRepo.On("GetChannelById",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.Channel{
						Id:      models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						Name:    "name",
						OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
					}, nil)

				f.ChannelRepo.On("DeleteChannel",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChannelRepo.AssertNumberOfCalls(t, "GetChannelById", 1)
				f.ChannelRepo.AssertNumberOfCalls(t, "DeleteChannel", 1)
			},
		},
		{
			name: "Test 3. Negative. Permission error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.DeleteChannelRequest{
					ChannelId:     "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "delete channel error: permission denied",

			on: func(f *fields) {
				f.ChannelRepo.On("GetChannelById",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.Channel{
						Id:      models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						Name:    "name",
						OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b796")),
					}, nil)

			},
			assert: func(t *testing.T, f *fields) {
				f.ChannelRepo.AssertNumberOfCalls(t, "GetChannelById", 1)
			},
		},
		{
			name: "Test 4. Negative. GetChannelById return error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.DeleteChannelRequest{
					ChannelId:     "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "get channel error: some error",

			on: func(f *fields) {
				f.ChannelRepo.On("GetChannelById",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.Channel{
						Id:      models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						Name:    "name",
						OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b796")),
					}, errors.New("some error"))

			},
			assert: func(t *testing.T, f *fields) {
				f.ChannelRepo.AssertNumberOfCalls(t, "GetChannelById", 1)
			},
		},
		{
			name: "Test 5. Negative. GetChannelById return not found error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.DeleteChannelRequest{
					ChannelId:     "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "channel: not found",

			on: func(f *fields) {
				f.ChannelRepo.On("GetChannelById",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.Channel{
						Id:      models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						Name:    "name",
						OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b796")),
					}, models.ErrNotFound)

			},
			assert: func(t *testing.T, f *fields) {
				f.ChannelRepo.AssertNumberOfCalls(t, "GetChannelById", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ChannelRepo:   mocks.NewChannelStorage(t),
				Log:           &log.Logger{},
				SubscribeRepo: mocks.NewSubscribeStorage(t),
			}
			au := NewChannelUsecase(Deps{
				ChannelRepo:   f.ChannelRepo,
				Log:           f.Log,
				SubscribeRepo: f.SubscribeRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.DeleteChannel(tt.args.ctx, tt.args.req)

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

func Test_usecase_ChannelUsecase_JoinChannel(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		ChannelRepo   *mocks.ChannelStorage
		Log           *log.Logger
		SubscribeRepo *mocks.SubscribeStorage
	}

	type args struct {
		ctx context.Context
		req usecases.JoinChannelRequest
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
				req: usecases.JoinChannelRequest{
					ChannelId:     "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr: false,

			on: func(f *fields) {
				f.ChannelRepo.On("GetChannelById",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.Channel{
						Id:      models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						Name:    "name",
						OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
					}, nil)

				f.SubscribeRepo.On("CreateSubscribe",
					ctx,
					mock.MatchedBy(func(info models.SubscribeInfo) bool {
						return info.Id.String() != "" &&
							info.ChannelId.String() == "284fef68-7e3e-4d1d-96a0-8c96f7b3b000" &&
							info.UserId.String() == "284fef68-7e3e-4d1d-96a0-8c96f7b3b795"
					})).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChannelRepo.AssertNumberOfCalls(t, "GetChannelById", 1)
				f.SubscribeRepo.AssertNumberOfCalls(t, "CreateSubscribe", 1)
			},
		},
		{
			name: "Test 2. Negative. CreateSubscribe returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.JoinChannelRequest{
					ChannelId:     "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "create subscribe channel error: some error",

			on: func(f *fields) {
				f.ChannelRepo.On("GetChannelById",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.Channel{
						Id:      models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						Name:    "name",
						OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
					}, nil)

				f.SubscribeRepo.On("CreateSubscribe",
					ctx,
					mock.MatchedBy(func(info models.SubscribeInfo) bool {
						return info.Id.String() != "" &&
							info.ChannelId.String() == "284fef68-7e3e-4d1d-96a0-8c96f7b3b000" &&
							info.UserId.String() == "284fef68-7e3e-4d1d-96a0-8c96f7b3b795"
					})).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChannelRepo.AssertNumberOfCalls(t, "GetChannelById", 1)
				f.SubscribeRepo.AssertNumberOfCalls(t, "CreateSubscribe", 1)
			},
		},
		{
			name: "Test 3. Negative. GetChannelById return error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.JoinChannelRequest{
					ChannelId:     "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "search channel error: some error",

			on: func(f *fields) {
				f.ChannelRepo.On("GetChannelById",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.Channel{
						Id:      models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						Name:    "name",
						OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b796")),
					}, errors.New("some error"))

			},
			assert: func(t *testing.T, f *fields) {
				f.ChannelRepo.AssertNumberOfCalls(t, "GetChannelById", 1)
			},
		},
		{
			name: "Test 5. Negative. GetChannelById return not found error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.JoinChannelRequest{
					ChannelId:     "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "channel: not found",

			on: func(f *fields) {
				f.ChannelRepo.On("GetChannelById",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000"))).
					Return(&models.Channel{
						Id:      models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						Name:    "name",
						OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b796")),
					}, models.ErrNotFound)

			},
			assert: func(t *testing.T, f *fields) {
				f.ChannelRepo.AssertNumberOfCalls(t, "GetChannelById", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ChannelRepo:   mocks.NewChannelStorage(t),
				Log:           &log.Logger{},
				SubscribeRepo: mocks.NewSubscribeStorage(t),
			}
			au := NewChannelUsecase(Deps{
				ChannelRepo:   f.ChannelRepo,
				Log:           f.Log,
				SubscribeRepo: f.SubscribeRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.JoinChannel(tt.args.ctx, tt.args.req)

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

func Test_usecase_ChannelUsecase_LeaveChannel(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		ChannelRepo   *mocks.ChannelStorage
		Log           *log.Logger
		SubscribeRepo *mocks.SubscribeStorage
	}

	type args struct {
		ctx context.Context
		req usecases.LeaveChannelRequest
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
				req: usecases.LeaveChannelRequest{
					ChannelId:     "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr: false,

			on: func(f *fields) {
				f.SubscribeRepo.On("DeleteSubscribe",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
				).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.SubscribeRepo.AssertNumberOfCalls(t, "DeleteSubscribe", 1)
			},
		},
		{
			name: "Test 2. Negative. DeleteSubscribe returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.LeaveChannelRequest{
					ChannelId:     "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "delete subscribe error: some error",

			on: func(f *fields) {
				f.SubscribeRepo.On("DeleteSubscribe",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
				).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.SubscribeRepo.AssertNumberOfCalls(t, "DeleteSubscribe", 1)
			},
		},
		{
			name: "Test 3. Negative. DeleteSubscribe returns not found error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.LeaveChannelRequest{
					ChannelId:     "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "delete channel error: not found",

			on: func(f *fields) {
				f.SubscribeRepo.On("DeleteSubscribe",
					ctx,
					models.ChannelID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
				).
					Return(models.ErrNotFound)
			},
			assert: func(t *testing.T, f *fields) {
				f.SubscribeRepo.AssertNumberOfCalls(t, "DeleteSubscribe", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ChannelRepo:   mocks.NewChannelStorage(t),
				Log:           &log.Logger{},
				SubscribeRepo: mocks.NewSubscribeStorage(t),
			}
			au := NewChannelUsecase(Deps{
				ChannelRepo:   f.ChannelRepo,
				Log:           f.Log,
				SubscribeRepo: f.SubscribeRepo,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.LeaveChannel(tt.args.ctx, tt.args.req)

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
