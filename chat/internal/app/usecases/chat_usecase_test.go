package usecases

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"github.com/Nixonxp/discord/chat/internal/app/usecases/mocks"
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
		ChatsStorage *mocks.ChatsStorage
	}

	type args struct {
		ctx context.Context
		req SendUserPrivateMessageRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *models.ActionInfo
		wantErr bool

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: SendUserPrivateMessageRequest{
					UserId: 1,
					Text:   "text",
				},
			},
			want: &models.ActionInfo{
				Success: true,
			},
			wantErr: false,

			on: func(f *fields) {
				f.ChatsStorage.On("CreateMessage", ctx, mock.MatchedBy(func(message *models.Message) bool {
					return message != nil &&
						message.Id != 0 &&
						message.UserId != 0 &&
						message.Text == "text" &&
						message.Timestamp != time.Time{}
				})).
					Return(&models.Message{
						Id:        1,
						ChatId:    2,
						UserId:    3,
						Text:      "text",
						Timestamp: time.Time{},
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatsStorage.AssertNumberOfCalls(t, "CreateMessage", 1)
			},
		},
		{
			name: "Test 2. Negative. CreateMessage return error",
			args: args{
				ctx: ctx, // dumm
				req: SendUserPrivateMessageRequest{
					UserId: 1,
					Text:   "text",
				},
			},
			want: &models.ActionInfo{
				Success: false,
			},
			wantErr: true,

			on: func(f *fields) {
				f.ChatsStorage.On("CreateMessage", ctx, mock.MatchedBy(func(message *models.Message) bool {
					return message != nil &&
						message.Id != 0 &&
						message.UserId != 0 &&
						message.Text == "text" &&
						message.Timestamp != time.Time{}
				})).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatsStorage.AssertNumberOfCalls(t, "CreateMessage", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ChatsStorage: mocks.NewChatsStorage(t),
			}
			au := NewChatUsecase(Deps{
				ChatRepo: f.ChatsStorage,
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
		ChatsStorage *mocks.ChatsStorage
	}

	type args struct {
		ctx context.Context
		req GetUserPrivateMessagesRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *models.Messages
		wantErr bool

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: GetUserPrivateMessagesRequest{
					UserId: 1,
				},
			},
			want: &models.Messages{
				Data: []*models.Message{
					{
						Id:        1,
						ChatId:    2,
						UserId:    3,
						Text:      "text",
						Timestamp: time.Time{},
					},
				},
			},
			wantErr: false,

			on: func(f *fields) {
				f.ChatsStorage.On("GetMessages", ctx, uint64(1), mock.AnythingOfType("uint64")).
					Return([]*models.Message{
						{
							Id:        1,
							ChatId:    2,
							UserId:    3,
							Text:      "text",
							Timestamp: time.Time{},
						},
					}, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatsStorage.AssertNumberOfCalls(t, "GetMessages", 1)
			},
		},
		{
			name: "Test 2. Negative. GetMessages return error",
			args: args{
				ctx: ctx, // dumm
				req: GetUserPrivateMessagesRequest{
					UserId: 1,
				},
			},
			want:    nil,
			wantErr: true,

			on: func(f *fields) {
				f.ChatsStorage.On("GetMessages", ctx, uint64(1), mock.AnythingOfType("uint64")).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatsStorage.AssertNumberOfCalls(t, "GetMessages", 1)
			},
		},
		{
			name: "Test 3. Negative. GetMessages return empty data",
			args: args{
				ctx: ctx, // dumm
				req: GetUserPrivateMessagesRequest{
					UserId: 1,
				},
			},
			want:    nil,
			wantErr: true,

			on: func(f *fields) {
				f.ChatsStorage.On("GetMessages", ctx, uint64(1), mock.AnythingOfType("uint64")).
					Return(nil, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.ChatsStorage.AssertNumberOfCalls(t, "GetMessages", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				ChatsStorage: mocks.NewChatsStorage(t),
			}
			au := NewChatUsecase(Deps{
				ChatRepo: f.ChatsStorage,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			got, err := au.GetUserPrivateMessages(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.GetUserPrivateMessages() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assert.Equal(t, tt.want, got)

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}
