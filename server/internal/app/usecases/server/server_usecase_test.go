package server

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/server/internal/app/models"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
	"github.com/Nixonxp/discord/server/internal/app/usecases/mocks"
	log "github.com/Nixonxp/discord/server/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_usecase_ServerUsecase_CreateServer(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		ServerRepo    *mocks.ServerStorage
		Log           *log.Logger
		SubscribeRepo *mocks.SubscribeStorage
		ChatService   *mocks.ServiceChatInterface
	}

	type args struct {
		ctx context.Context
		req usecases.CreateServerRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.ServerInfo
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.CreateServerRequest{
					Name:          "name",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ServerInfo{
				Id:      models.ServerID{},
				Name:    "name",
				OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
			},
			wantErr: false,

			on: func(f *fields) {
				f.ServerRepo.On("CreateServer",
					ctx,
					mock.MatchedBy(func(server models.ServerInfo) bool {
						return server.Id.String() != "" &&
							server.OwnerId == models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) &&
							server.Name == "name"
					})).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "CreateServer", 1)
			},
		},
		{
			name: "Test 2. Negative. CreateServer returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.CreateServerRequest{
					Name:          "name",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ServerInfo{
				Id:      models.ServerID{},
				Name:    "name",
				OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
			},
			wantErr:     true,
			errorString: "create server error: some error",

			on: func(f *fields) {
				f.ServerRepo.On("CreateServer",
					ctx,
					mock.MatchedBy(func(server models.ServerInfo) bool {
						return server.Id.String() != "" &&
							server.OwnerId == models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) &&
							server.Name == "name"
					})).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "CreateServer", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ServerRepo:    mocks.NewServerStorage(t),
				Log:           &log.Logger{},
				SubscribeRepo: mocks.NewSubscribeStorage(t),
				ChatService:   mocks.NewServiceChatInterface(t),
			}
			au := NewServerUsecase(Deps{
				ServerRepo:    f.ServerRepo,
				Log:           f.Log,
				SubscribeRepo: f.SubscribeRepo,
				ChatService:   f.ChatService,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.CreateServer(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			if got != nil {
				got.Id = models.ServerID{}
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}

func Test_usecase_ServerUsecase_SearchServer(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		ServerRepo    *mocks.ServerStorage
		Log           *log.Logger
		SubscribeRepo *mocks.SubscribeStorage
		ChatService   *mocks.ServiceChatInterface
	}

	type args struct {
		ctx context.Context
		req usecases.SearchServerRequest
	}
	tests := []struct {
		name        string
		args        args
		want        []*models.ServerInfo
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SearchServerRequest{
					Name: "name",
				},
			},
			want: []*models.ServerInfo{
				{
					Id:      models.ServerID{},
					Name:    "name",
					OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
				},
			},
			wantErr: false,

			on: func(f *fields) {
				f.ServerRepo.On("SearchServers",
					ctx,
					"name").
					Return([]*models.ServerInfo{
						{
							Id:      models.ServerID{},
							Name:    "name",
							OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						},
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "SearchServers", 1)
			},
		},
		{
			name: "Test 2. Negative. SearchServers return empty result",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SearchServerRequest{
					Name: "name",
				},
			},
			want: []*models.ServerInfo{
				{
					Id:      models.ServerID{},
					Name:    "name",
					OwnerId: models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
				},
			},
			wantErr:     true,
			errorString: "servers: not found",

			on: func(f *fields) {
				f.ServerRepo.On("SearchServers",
					ctx,
					"name").
					Return([]*models.ServerInfo{}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "SearchServers", 1)
			},
		},
		{
			name: "Test 3. Negative. SearchServers return error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SearchServerRequest{
					Name: "name",
				},
			},
			want:        nil,
			wantErr:     true,
			errorString: "search server error: some  error",

			on: func(f *fields) {
				f.ServerRepo.On("SearchServers",
					ctx,
					"name").
					Return(nil, errors.New("some  error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "SearchServers", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ServerRepo:    mocks.NewServerStorage(t),
				Log:           &log.Logger{},
				SubscribeRepo: mocks.NewSubscribeStorage(t),
				ChatService:   mocks.NewServiceChatInterface(t),
			}
			au := NewServerUsecase(Deps{
				ServerRepo:    f.ServerRepo,
				Log:           f.Log,
				SubscribeRepo: f.SubscribeRepo,
				ChatService:   f.ChatService,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.SearchServer(tt.args.ctx, tt.args.req)

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

func Test_usecase_ServerUsecase_SubscribeServer(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		ServerRepo    *mocks.ServerStorage
		Log           *log.Logger
		SubscribeRepo *mocks.SubscribeStorage
		ChatService   *mocks.ServiceChatInterface
	}

	type args struct {
		ctx context.Context
		req usecases.SubscribeServerRequest
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
				req: usecases.SubscribeServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr: false,

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000").
					Return(nil, nil)

				f.SubscribeRepo.On("CreateSubscribe",
					ctx,
					mock.MatchedBy(func(subscribe models.SubscribeInfo) bool {
						return subscribe.Id.String() != "" &&
							subscribe.UserId == models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) &&
							subscribe.ServerId.String() == "284fef68-7e3e-4d1d-96a0-8c96f7b3b000"
					})).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
				f.SubscribeRepo.AssertNumberOfCalls(t, "CreateSubscribe", 1)
			},
		},
		{
			name: "Test 2. Negative. CreateSubscribe returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SubscribeServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "subscribe create error: some error",

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000").
					Return(nil, nil)

				f.SubscribeRepo.On("CreateSubscribe",
					ctx,
					mock.MatchedBy(func(subscribe models.SubscribeInfo) bool {
						return subscribe.Id.String() != "" &&
							subscribe.UserId == models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) &&
							subscribe.ServerId.String() == "284fef68-7e3e-4d1d-96a0-8c96f7b3b000"
					})).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
				f.SubscribeRepo.AssertNumberOfCalls(t, "CreateSubscribe", 1)
			},
		},
		{
			name: "Test 3. Negative. GetServerById returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SubscribeServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "subscribe server error: some error",

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000").
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
			},
		},
		{
			name: "Test 4. Negative. GetServerById returns not found error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SubscribeServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "subscribe server error: not found",

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000").
					Return(nil, models.ErrNotFound)
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ServerRepo:    mocks.NewServerStorage(t),
				Log:           &log.Logger{},
				SubscribeRepo: mocks.NewSubscribeStorage(t),
				ChatService:   mocks.NewServiceChatInterface(t),
			}
			au := NewServerUsecase(Deps{
				ServerRepo:    f.ServerRepo,
				Log:           f.Log,
				SubscribeRepo: f.SubscribeRepo,
				ChatService:   f.ChatService,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.SubscribeServer(tt.args.ctx, tt.args.req)

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

func Test_usecase_ServerUsecase_UnsubscribeServer(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		ServerRepo    *mocks.ServerStorage
		Log           *log.Logger
		SubscribeRepo *mocks.SubscribeStorage
		ChatService   *mocks.ServiceChatInterface
	}

	type args struct {
		ctx context.Context
		req usecases.UnsubscribeServerRequest
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
				req: usecases.UnsubscribeServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
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
					models.ServerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
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
				req: usecases.UnsubscribeServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "unsubscribe error: some  error",

			on: func(f *fields) {
				f.SubscribeRepo.On("DeleteSubscribe",
					ctx,
					models.ServerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
				).
					Return(errors.New("some  error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.SubscribeRepo.AssertNumberOfCalls(t, "DeleteSubscribe", 1)
			},
		},
		{
			name: "Test 3. Negative. DeleteSubscribe returns not found error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.UnsubscribeServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "unsubscribe error: not found",

			on: func(f *fields) {
				f.SubscribeRepo.On("DeleteSubscribe",
					ctx,
					models.ServerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
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
				ServerRepo:    mocks.NewServerStorage(t),
				Log:           &log.Logger{},
				SubscribeRepo: mocks.NewSubscribeStorage(t),
				ChatService:   mocks.NewServiceChatInterface(t),
			}
			au := NewServerUsecase(Deps{
				ServerRepo:    f.ServerRepo,
				Log:           f.Log,
				SubscribeRepo: f.SubscribeRepo,
				ChatService:   f.ChatService,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.UnsubscribeServer(tt.args.ctx, tt.args.req)

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

func Test_usecase_ServerUsecase_SearchServerByUserId(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		ServerRepo    *mocks.ServerStorage
		Log           *log.Logger
		SubscribeRepo *mocks.SubscribeStorage
		ChatService   *mocks.ServiceChatInterface
	}

	type args struct {
		ctx context.Context
		req usecases.SearchServerByUserIdRequest
	}
	tests := []struct {
		name        string
		args        args
		want        []string
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SearchServerByUserIdRequest{
					UserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: []string{
				"284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
			},
			wantErr: false,

			on: func(f *fields) {
				f.SubscribeRepo.On("GetByUserId",
					ctx,
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
				).
					Return([]*models.SubscribeInfo{
						{
							Id:       models.SubscribeID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b300")),
							ServerId: models.ServerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
							UserId:   models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						},
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.SubscribeRepo.AssertNumberOfCalls(t, "GetByUserId", 1)
			},
		},
		{
			name: "Test 2. Negative. GetByUserId returns empty",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SearchServerByUserIdRequest{
					UserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: []string{
				"284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
			},
			wantErr:     true,
			errorString: "servers: not found",

			on: func(f *fields) {
				f.SubscribeRepo.On("GetByUserId",
					ctx,
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
				).
					Return([]*models.SubscribeInfo{}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.SubscribeRepo.AssertNumberOfCalls(t, "GetByUserId", 1)
			},
		},
		{
			name: "Test 3. Negative. GetByUserId returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SearchServerByUserIdRequest{
					UserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: []string{
				"284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
			},
			wantErr:     true,
			errorString: "search server error: some error",

			on: func(f *fields) {
				f.SubscribeRepo.On("GetByUserId",
					ctx,
					models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
				).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.SubscribeRepo.AssertNumberOfCalls(t, "GetByUserId", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ServerRepo:    mocks.NewServerStorage(t),
				Log:           &log.Logger{},
				SubscribeRepo: mocks.NewSubscribeStorage(t),
				ChatService:   mocks.NewServiceChatInterface(t),
			}
			au := NewServerUsecase(Deps{
				ServerRepo:    f.ServerRepo,
				Log:           f.Log,
				SubscribeRepo: f.SubscribeRepo,
				ChatService:   f.ChatService,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.SearchServerByUserId(tt.args.ctx, tt.args.req)

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

func Test_usecase_ServerUsecase_InviteUserToServer(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		ServerRepo    *mocks.ServerStorage
		Log           *log.Logger
		SubscribeRepo *mocks.SubscribeStorage
		ChatService   *mocks.ServiceChatInterface
	}

	type args struct {
		ctx context.Context
		req usecases.InviteUserToServerRequest
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
				req: usecases.InviteUserToServerRequest{
					ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					UserId:   "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr: false,

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				).
					Return(nil, nil)

				f.SubscribeRepo.On("CreateSubscribe",
					ctx,
					mock.MatchedBy(func(subscribe models.SubscribeInfo) bool {
						return subscribe.Id.String() != "" &&
							subscribe.UserId == models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) &&
							subscribe.ServerId.String() == "284fef68-7e3e-4d1d-96a0-8c96f7b3b000"
					}),
				).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
				f.SubscribeRepo.AssertNumberOfCalls(t, "CreateSubscribe", 1)
			},
		},
		{
			name: "Test 2. Negative. InviteUserToServer returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.InviteUserToServerRequest{
					ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					UserId:   "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "create subscribe: some error",

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				).
					Return(nil, nil)

				f.SubscribeRepo.On("CreateSubscribe",
					ctx,
					mock.MatchedBy(func(subscribe models.SubscribeInfo) bool {
						return subscribe.Id.String() != "" &&
							subscribe.UserId == models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) &&
							subscribe.ServerId.String() == "284fef68-7e3e-4d1d-96a0-8c96f7b3b000"
					}),
				).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
				f.SubscribeRepo.AssertNumberOfCalls(t, "CreateSubscribe", 1)
			},
		},
		{
			name: "Test 3. Negative. GetServerById returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.InviteUserToServerRequest{
					ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					UserId:   "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "server search: some error",

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ServerRepo:    mocks.NewServerStorage(t),
				Log:           &log.Logger{},
				SubscribeRepo: mocks.NewSubscribeStorage(t),
				ChatService:   mocks.NewServiceChatInterface(t),
			}
			au := NewServerUsecase(Deps{
				ServerRepo:    f.ServerRepo,
				Log:           f.Log,
				SubscribeRepo: f.SubscribeRepo,
				ChatService:   f.ChatService,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.InviteUserToServer(tt.args.ctx, tt.args.req)

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

func Test_usecase_ServerUsecase_PublishMessageOnServer(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		ServerRepo    *mocks.ServerStorage
		Log           *log.Logger
		SubscribeRepo *mocks.SubscribeStorage
		ChatService   *mocks.ServiceChatInterface
	}

	type args struct {
		ctx context.Context
		req usecases.PublishMessageOnServerRequest
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
				req: usecases.PublishMessageOnServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					Text:          "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr: false,

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				).
					Return(nil, nil)

				f.ChatService.On("SendServerMessage",
					ctx,
					usecases.PublishMessageOnServerRequest{
						ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
						Text:     "text",
					},
				).
					Return(nil, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
				f.ChatService.AssertNumberOfCalls(t, "SendServerMessage", 1)
			},
		},
		{
			name: "Test 2. Negative. SendServerMessage returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.PublishMessageOnServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					Text:          "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "publish message on server: some error",

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				).
					Return(nil, nil)

				f.ChatService.On("SendServerMessage",
					ctx,
					usecases.PublishMessageOnServerRequest{
						ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
						Text:     "text",
					},
				).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
				f.ChatService.AssertNumberOfCalls(t, "SendServerMessage", 1)
			},
		},
		{
			name: "Test 3. Negative. SendServerMessage returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.PublishMessageOnServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					Text:          "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "publish message on server: some error",

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
			},
		},
		{
			name: "Test 4. Negative. SendServerMessage returns not found error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.PublishMessageOnServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					Text:          "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "publish message on server: not found",

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				).
					Return(nil, models.ErrNotFound)
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ServerRepo:    mocks.NewServerStorage(t),
				Log:           &log.Logger{},
				SubscribeRepo: mocks.NewSubscribeStorage(t),
				ChatService:   mocks.NewServiceChatInterface(t),
			}
			au := NewServerUsecase(Deps{
				ServerRepo:    f.ServerRepo,
				Log:           f.Log,
				SubscribeRepo: f.SubscribeRepo,
				ChatService:   f.ChatService,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.PublishMessageOnServer(tt.args.ctx, tt.args.req)

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

func Test_usecase_ServerUsecase_GetMessagesFromServer(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		ServerRepo    *mocks.ServerStorage
		Log           *log.Logger
		SubscribeRepo *mocks.SubscribeStorage
		ChatService   *mocks.ServiceChatInterface
	}

	timeMock, _ := time.Parse("2006 Jan 02 15:04:05", "2012 Dec 07 00:00:00")

	type args struct {
		ctx context.Context
		req usecases.GetMessagesFromServerRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.GetMessagesInfo
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetMessagesFromServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.GetMessagesInfo{
				Messages: []*models.Message{
					{
						Id:        "284fef68-7e3e-4d1d-96a0-8c96f7b3b200",
						Text:      "text",
						Timestamp: timeMock,
					},
				},
			},
			wantErr: false,

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				).
					Return(nil, nil)

				f.ChatService.On("GetServerMessages",
					ctx,
					usecases.GetMessagesFromServerRequest{
						ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					},
				).
					Return(&models.GetMessagesInfo{
						Messages: []*models.Message{
							{
								Id:        "284fef68-7e3e-4d1d-96a0-8c96f7b3b200",
								Text:      "text",
								Timestamp: timeMock,
							},
						},
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
				f.ChatService.AssertNumberOfCalls(t, "GetServerMessages", 1)
			},
		},
		{
			name: "Test 2. Negative. GetServerMessages returns errors",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetMessagesFromServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.GetMessagesInfo{
				Messages: []*models.Message{
					{
						Id:        "284fef68-7e3e-4d1d-96a0-8c96f7b3b200",
						Text:      "text",
						Timestamp: timeMock,
					},
				},
			},
			wantErr:     true,
			errorString: "get messages on server: some error",

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				).
					Return(nil, nil)

				f.ChatService.On("GetServerMessages",
					ctx,
					usecases.GetMessagesFromServerRequest{
						ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					},
				).
					Return(&models.GetMessagesInfo{
						Messages: []*models.Message{
							{
								Id:        "284fef68-7e3e-4d1d-96a0-8c96f7b3b200",
								Text:      "text",
								Timestamp: timeMock,
							},
						},
					}, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
				f.ChatService.AssertNumberOfCalls(t, "GetServerMessages", 1)
			},
		},
		{
			name: "Test 3. Negative. GetServerById returns errors",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetMessagesFromServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.GetMessagesInfo{
				Messages: []*models.Message{
					{
						Id:        "284fef68-7e3e-4d1d-96a0-8c96f7b3b200",
						Text:      "text",
						Timestamp: timeMock,
					},
				},
			},
			wantErr:     true,
			errorString: "get messages on server: some error",

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
			},
		},
		{
			name: "Test 4. Negative. GetServerById returns not found errors",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetMessagesFromServerRequest{
					ServerId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
					CurrentUserId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
				},
			},
			want: &models.GetMessagesInfo{
				Messages: []*models.Message{
					{
						Id:        "284fef68-7e3e-4d1d-96a0-8c96f7b3b200",
						Text:      "text",
						Timestamp: timeMock,
					},
				},
			},
			wantErr:     true,
			errorString: "get messages on server: not found",

			on: func(f *fields) {
				f.ServerRepo.On("GetServerById",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
				).
					Return(nil, models.ErrNotFound)
			},
			assert: func(t *testing.T, f *fields) {
				f.ServerRepo.AssertNumberOfCalls(t, "GetServerById", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ServerRepo:    mocks.NewServerStorage(t),
				Log:           &log.Logger{},
				SubscribeRepo: mocks.NewSubscribeStorage(t),
				ChatService:   mocks.NewServiceChatInterface(t),
			}
			au := NewServerUsecase(Deps{
				ServerRepo:    f.ServerRepo,
				Log:           f.Log,
				SubscribeRepo: f.SubscribeRepo,
				ChatService:   f.ChatService,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.GetMessagesFromServer(tt.args.ctx, tt.args.req)

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
