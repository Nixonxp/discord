package chat

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/chat/internal/app/enum"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	"github.com/Nixonxp/discord/chat/internal/app/usecases/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func Test_usecase_ChatUsecase_SendUserPrivateMessage(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		MessagesRepo *mocks.MessagesStorage
		ChatRepo     *mocks.ChatStorage
		KafkaConn    *mocks.KafkaServiceInterface
	}

	type args struct {
		ctx context.Context
		req usecases.SendUserPrivateMessageRequest
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
				req: usecases.SendUserPrivateMessageRequest{
					UserId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					CurrentUser: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					Text:        "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     false,
			errorString: "",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b800_284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					enum.PrivateChatType,
				).
					Return(&models.Chat{
						Id:       models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						Type:     enum.PrivateChatType,
						OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795_284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					}, nil)

				f.KafkaConn.On("SendMessage",
					usecases.MessageDto{
						ChatId:  "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
						OwnerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
						Text:    "text",
					},
				).
					Return(nil, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.KafkaConn.AssertNumberOfCalls(t, "SendMessage", 1)
			},
		},
		{
			name: "Test 2. Negative. SendMessage returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SendUserPrivateMessageRequest{
					UserId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					CurrentUser: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					Text:        "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "chat message send error: some error",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b800_284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					enum.PrivateChatType,
				).
					Return(&models.Chat{
						Id:       models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						Type:     enum.PrivateChatType,
						OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795_284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					}, nil)

				f.KafkaConn.On("SendMessage",
					usecases.MessageDto{
						ChatId:  "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
						OwnerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
						Text:    "text",
					},
				).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.KafkaConn.AssertNumberOfCalls(t, "SendMessage", 1)
			},
		},
		{
			name: "Test 3. Negative. GetChatByMetadataAndType returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SendUserPrivateMessageRequest{
					UserId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					CurrentUser: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					Text:        "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "chat search error: some error",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b800_284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
					enum.PrivateChatType,
				).
					Return(&models.Chat{
						Id:       models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b000")),
						Type:     enum.PrivateChatType,
						OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")),
						MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b795_284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					}, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				MessagesRepo: mocks.NewMessagesStorage(t),
				ChatRepo:     mocks.NewChatStorage(t),
				KafkaConn:    mocks.NewKafkaServiceInterface(t),
			}
			au := NewChatUsecase(Deps{
				MessagesRepo: f.MessagesRepo,
				ChatRepo:     f.ChatRepo,
				KafkaConn:    f.KafkaConn,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.SendUserPrivateMessage(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.SendUserPrivateMessage() error = %v, wantErr %v", err, tt.wantErr)
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

func Test_usecase_ChatUsecase_SendServerMessage(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		MessagesRepo *mocks.MessagesStorage
		ChatRepo     *mocks.ChatStorage
		KafkaConn    *mocks.KafkaServiceInterface
	}

	type args struct {
		ctx context.Context
		req usecases.SendServerMessageRequest
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
				req: usecases.SendServerMessageRequest{
					ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					Text:     "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     false,
			errorString: "",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					enum.ServerChatType,
				).
					Return(&models.Chat{
						Id:       models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						Type:     enum.ServerChatType,
						OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					}, nil)

				f.KafkaConn.On("SendMessage",
					usecases.MessageDto{
						ChatId:  "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
						OwnerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
						Text:    "text",
					},
				).
					Return(nil, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.KafkaConn.AssertNumberOfCalls(t, "SendMessage", 1)
			},
		},
		{
			name: "Test 2. Positive. Create new chat",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SendServerMessageRequest{
					ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					Text:     "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     false,
			errorString: "",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					enum.ServerChatType,
				).
					Return(nil, models.ErrNotFound)

				f.ChatRepo.On("CreateChat",
					ctx,
					mock.MatchedBy(func(chat *models.Chat) bool {
						return chat.Id.String() != "" &&
							chat.Type == enum.ServerChatType &&
							chat.OwnerId == models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")) &&
							chat.MetaData == "284fef68-7e3e-4d1d-96a0-8c96f7b3b800"
					}),
				).
					Return(nil)

				f.KafkaConn.On("SendMessage",
					mock.MatchedBy(func(dto usecases.MessageDto) bool {
						return dto.ChatId != "" &&
							dto.Text == "text" &&
							dto.OwnerId == "284fef68-7e3e-4d1d-96a0-8c96f7b3b800"
					}),
				).
					Return(nil, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.ChatRepo.AssertNumberOfCalls(t, "CreateChat", 1)
				f.KafkaConn.AssertNumberOfCalls(t, "SendMessage", 1)
			},
		},
		{
			name: "Test 3. Negative. SendMessage return error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SendServerMessageRequest{
					ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					Text:     "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "send message to server server: some  error",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					enum.ServerChatType,
				).
					Return(nil, models.ErrNotFound)

				f.ChatRepo.On("CreateChat",
					ctx,
					mock.MatchedBy(func(chat *models.Chat) bool {
						return chat.Id.String() != "" &&
							chat.Type == enum.ServerChatType &&
							chat.OwnerId == models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")) &&
							chat.MetaData == "284fef68-7e3e-4d1d-96a0-8c96f7b3b800"
					}),
				).
					Return(nil)

				f.KafkaConn.On("SendMessage",
					mock.MatchedBy(func(dto usecases.MessageDto) bool {
						return dto.ChatId != "" &&
							dto.Text == "text" &&
							dto.OwnerId == "284fef68-7e3e-4d1d-96a0-8c96f7b3b800"
					}),
				).
					Return(errors.New("some  error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.ChatRepo.AssertNumberOfCalls(t, "CreateChat", 1)
				f.KafkaConn.AssertNumberOfCalls(t, "SendMessage", 1)
			},
		},
		{
			name: "Test 4. Negative. CreateChat return error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SendServerMessageRequest{
					ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					Text:     "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "create new chat for server: some error",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					enum.ServerChatType,
				).
					Return(nil, models.ErrNotFound)

				f.ChatRepo.On("CreateChat",
					ctx,
					mock.MatchedBy(func(chat *models.Chat) bool {
						return chat.Id.String() != "" &&
							chat.Type == enum.ServerChatType &&
							chat.OwnerId == models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")) &&
							chat.MetaData == "284fef68-7e3e-4d1d-96a0-8c96f7b3b800"
					}),
				).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.ChatRepo.AssertNumberOfCalls(t, "CreateChat", 1)
			},
		},
		{
			name: "Test 5. Negative. GetChatByMetadataAndType return error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.SendServerMessageRequest{
					ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					Text:     "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr:     true,
			errorString: "chat search error: some error",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					enum.ServerChatType,
				).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				MessagesRepo: mocks.NewMessagesStorage(t),
				ChatRepo:     mocks.NewChatStorage(t),
				KafkaConn:    mocks.NewKafkaServiceInterface(t),
			}
			au := NewChatUsecase(Deps{
				MessagesRepo: f.MessagesRepo,
				ChatRepo:     f.ChatRepo,
				KafkaConn:    f.KafkaConn,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.SendServerMessage(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.SendUserPrivateMessage() error = %v, wantErr %v", err, tt.wantErr)
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

func Test_usecase_ChatUsecase_GetUserPrivateMessages(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		MessagesRepo *mocks.MessagesStorage
		ChatRepo     *mocks.ChatStorage
		KafkaConn    *mocks.KafkaServiceInterface
	}

	type args struct {
		ctx context.Context
		req usecases.GetUserPrivateMessagesRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.Messages
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetUserPrivateMessagesRequest{
					UserId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					CurrentUser: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600",
				},
			},
			want: &models.Messages{
				Data: []*models.Message{
					{
						Id:        models.MessageID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b111")),
						Text:      "text",
						ChatId:    models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						OwnerId:   models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						Timestamp: time.Time{},
					},
				},
			},
			wantErr:     false,
			errorString: "",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					enum.PrivateChatType,
				).
					Return(&models.Chat{
						Id:       models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						Type:     enum.PrivateChatType,
						OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					}, nil)

				f.MessagesRepo.On("GetMessages",
					ctx,
					models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
				).
					Return([]*models.Message{
						{
							Id:        models.MessageID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b111")),
							Text:      "text",
							ChatId:    models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
							OwnerId:   models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
							Timestamp: time.Time{},
						},
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.MessagesRepo.AssertNumberOfCalls(t, "GetMessages", 1)
			},
		},
		{
			name: "Test 2. Positive. GetMessages is empty",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetUserPrivateMessagesRequest{
					UserId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					CurrentUser: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600",
				},
			},
			want:        nil,
			wantErr:     true,
			errorString: "get messages error: error empty",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					enum.PrivateChatType,
				).
					Return(&models.Chat{
						Id:       models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						Type:     enum.ServerChatType,
						OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					}, nil)

				f.MessagesRepo.On("GetMessages",
					ctx,
					models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
				).
					Return([]*models.Message{}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.MessagesRepo.AssertNumberOfCalls(t, "GetMessages", 1)
			},
		},
		{
			name: "Test 3. Negative. GetMessages is returned error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetUserPrivateMessagesRequest{
					UserId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					CurrentUser: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600",
				},
			},
			want:        nil,
			wantErr:     true,
			errorString: "get messages error: some error",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					enum.PrivateChatType,
				).
					Return(&models.Chat{
						Id:       models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						Type:     enum.ServerChatType,
						OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					}, nil)

				f.MessagesRepo.On("GetMessages",
					ctx,
					models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
				).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.MessagesRepo.AssertNumberOfCalls(t, "GetMessages", 1)
			},
		},
		{
			name: "Test 4. Negative. GetChatByMetadataAndType is returned error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetUserPrivateMessagesRequest{
					UserId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					CurrentUser: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600",
				},
			},
			want:        nil,
			wantErr:     true,
			errorString: "get chat error: some error",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b800",
					enum.PrivateChatType,
				).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				MessagesRepo: mocks.NewMessagesStorage(t),
				ChatRepo:     mocks.NewChatStorage(t),
				KafkaConn:    mocks.NewKafkaServiceInterface(t),
			}
			au := NewChatUsecase(Deps{
				MessagesRepo: f.MessagesRepo,
				ChatRepo:     f.ChatRepo,
				KafkaConn:    f.KafkaConn,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.GetUserPrivateMessages(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.SendUserPrivateMessage() error = %v, wantErr %v", err, tt.wantErr)
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

func Test_usecase_ChatUsecase_GetServerMessagesRequest(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		MessagesRepo *mocks.MessagesStorage
		ChatRepo     *mocks.ChatStorage
		KafkaConn    *mocks.KafkaServiceInterface
	}

	type args struct {
		ctx context.Context
		req usecases.GetServerMessageRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.Messages
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetServerMessageRequest{
					ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
				},
			},
			want: &models.Messages{
				Data: []*models.Message{
					{
						Id:        models.MessageID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b111")),
						Text:      "text",
						ChatId:    models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						OwnerId:   models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						Timestamp: time.Time{},
					},
				},
			},
			wantErr:     false,
			errorString: "",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
					enum.ServerChatType,
				).
					Return(&models.Chat{
						Id:       models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
						Type:     enum.ServerChatType,
						OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
						MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
					}, nil)

				f.MessagesRepo.On("GetMessages",
					ctx,
					models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
				).
					Return([]*models.Message{
						{
							Id:        models.MessageID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b111")),
							Text:      "text",
							ChatId:    models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
							OwnerId:   models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
							Timestamp: time.Time{},
						},
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.MessagesRepo.AssertNumberOfCalls(t, "GetMessages", 1)
			},
		},
		{
			name: "Test 2. Positive. GetMessages return empty messages",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetServerMessageRequest{
					ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
				},
			},
			want: &models.Messages{
				Data: []*models.Message{
					{
						Id:        models.MessageID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b111")),
						Text:      "text",
						ChatId:    models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						OwnerId:   models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						Timestamp: time.Time{},
					},
				},
			},
			wantErr:     true,
			errorString: "get server messages error: error empty",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
					enum.ServerChatType,
				).
					Return(&models.Chat{
						Id:       models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
						Type:     enum.ServerChatType,
						OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
						MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
					}, nil)

				f.MessagesRepo.On("GetMessages",
					ctx,
					models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
				).
					Return([]*models.Message{}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.MessagesRepo.AssertNumberOfCalls(t, "GetMessages", 1)
			},
		},
		{
			name: "Test 3. Negative. GetMessages return errors",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetServerMessageRequest{
					ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
				},
			},
			want: &models.Messages{
				Data: []*models.Message{
					{
						Id:        models.MessageID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b111")),
						Text:      "text",
						ChatId:    models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						OwnerId:   models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						Timestamp: time.Time{},
					},
				},
			},
			wantErr:     true,
			errorString: "get server messages error: some errors",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
					enum.ServerChatType,
				).
					Return(&models.Chat{
						Id:       models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
						Type:     enum.ServerChatType,
						OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
						MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
					}, nil)

				f.MessagesRepo.On("GetMessages",
					ctx,
					models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
				).
					Return(nil, errors.New("some errors"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.MessagesRepo.AssertNumberOfCalls(t, "GetMessages", 1)
			},
		},
		{
			name: "Test 4. Negative. GetChatByMetadataAndType return errors",
			args: args{
				ctx: ctx, // dumm
				req: usecases.GetServerMessageRequest{
					ServerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
				},
			},
			want: &models.Messages{
				Data: []*models.Message{
					{
						Id:        models.MessageID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b111")),
						Text:      "text",
						ChatId:    models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						OwnerId:   models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b800")),
						Timestamp: time.Time{},
					},
				},
			},
			wantErr:     true,
			errorString: "get server chat error: some errors",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
					enum.ServerChatType,
				).
					Return(nil, errors.New("some errors"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				MessagesRepo: mocks.NewMessagesStorage(t),
				ChatRepo:     mocks.NewChatStorage(t),
				KafkaConn:    mocks.NewKafkaServiceInterface(t),
			}
			au := NewChatUsecase(Deps{
				MessagesRepo: f.MessagesRepo,
				ChatRepo:     f.ChatRepo,
				KafkaConn:    f.KafkaConn,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.GetServerMessagesRequest(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.SendUserPrivateMessage() error = %v, wantErr %v", err, tt.wantErr)
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

func Test_usecase_ChatUsecase_CreatePrivateChat(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		MessagesRepo *mocks.MessagesStorage
		ChatRepo     *mocks.ChatStorage
		KafkaConn    *mocks.KafkaServiceInterface
	}

	type args struct {
		ctx context.Context
		req usecases.CreatePrivateChatRequest
	}
	tests := []struct {
		name        string
		args        args
		want        *models.Chat
		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: usecases.CreatePrivateChatRequest{
					UserId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
					CurrentUser: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600",
				},
			},
			want: &models.Chat{
				Id:       models.ChatID{},
				Type:     enum.PrivateChatType,
				OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b300")),
				MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
			},
			wantErr:     false,
			errorString: "",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
					enum.PrivateChatType,
				).
					Return(&models.Chat{
						Id:       models.ChatID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b100")),
						Type:     enum.PrivateChatType,
						OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b300")),
						MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
			},
		},
		{
			name: "Test 2. Positive. Create new chat",
			args: args{
				ctx: ctx, // dumm
				req: usecases.CreatePrivateChatRequest{
					UserId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
					CurrentUser: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600",
				},
			},
			want: &models.Chat{
				Id:       models.ChatID{},
				Type:     enum.PrivateChatType,
				OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b600")),
				MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
			},
			wantErr:     false,
			errorString: "",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
					enum.PrivateChatType,
				).
					Return(nil, models.ErrNotFound)

				f.ChatRepo.On("CreateChat",
					ctx,
					mock.MatchedBy(func(chat *models.Chat) bool {
						return chat.Id.String() != "" &&
							chat.Type == enum.PrivateChatType &&
							chat.OwnerId == models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b600")) &&
							chat.MetaData == "284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b500"
					}),
				).
					Return(nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.ChatRepo.AssertNumberOfCalls(t, "CreateChat", 1)
			},
		},
		{
			name: "Test 4. Negative. CreateChat returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.CreatePrivateChatRequest{
					UserId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
					CurrentUser: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600",
				},
			},
			want: &models.Chat{
				Id:       models.ChatID{},
				Type:     enum.PrivateChatType,
				OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b600")),
				MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
			},
			wantErr:     true,
			errorString: "crate chat error: some error",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
					enum.PrivateChatType,
				).
					Return(nil, models.ErrNotFound)

				f.ChatRepo.On("CreateChat",
					ctx,
					mock.MatchedBy(func(chat *models.Chat) bool {
						return chat.Id.String() != "" &&
							chat.Type == enum.PrivateChatType &&
							chat.OwnerId == models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b600")) &&
							chat.MetaData == "284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b500"
					}),
				).
					Return(errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
				f.ChatRepo.AssertNumberOfCalls(t, "CreateChat", 1)
			},
		},
		{
			name: "Test 5. Negative. GetChatByMetadataAndType returns error",
			args: args{
				ctx: ctx, // dumm
				req: usecases.CreatePrivateChatRequest{
					UserId:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
					CurrentUser: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600",
				},
			},
			want: &models.Chat{
				Id:       models.ChatID{},
				Type:     enum.PrivateChatType,
				OwnerId:  models.OwnerID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b600")),
				MetaData: "284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
			},
			wantErr:     true,
			errorString: "get chat error: not found",

			on: func(f *fields) {
				f.ChatRepo.On("GetChatByMetadataAndType",
					ctx,
					"284fef68-7e3e-4d1d-96a0-8c96f7b3b600_284fef68-7e3e-4d1d-96a0-8c96f7b3b500",
					enum.PrivateChatType,
				).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatRepo.AssertNumberOfCalls(t, "GetChatByMetadataAndType", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				MessagesRepo: mocks.NewMessagesStorage(t),
				ChatRepo:     mocks.NewChatStorage(t),
				KafkaConn:    mocks.NewKafkaServiceInterface(t),
			}
			au := NewChatUsecase(Deps{
				MessagesRepo: f.MessagesRepo,
				ChatRepo:     f.ChatRepo,
				KafkaConn:    f.KafkaConn,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.CreatePrivateChat(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.SendUserPrivateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			if got != nil {
				got.Id = models.ChatID{}
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}
