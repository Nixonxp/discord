package queue

import (
	"context"
	"encoding/json"
	"errors"
	config "github.com/Nixonxp/discord/chat/configs"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	"github.com/Nixonxp/discord/chat/internal/app/usecases/mocks"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_usecase_ConsumerHandler_CreateMessage(t *testing.T) {
	// prepare
	var (
		ctx = context.Background() // dummy
	)
	type fields struct {
		QueueUsecase    *mocks.QueueInterface
		Cfg             *config.Config
		ConsumerService *mocks.KafkaConsumerServiceInterface
	}

	type args struct {
		ctx context.Context
		req kafka.Message
	}

	msgDto := MessageKafkaMessage{
		Id:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
		Text:    "text",
		ChatId:  "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
		OwnerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
	}

	msgBytes, _ := json.Marshal(msgDto)

	tests := []struct {
		name string
		args args

		wantErr     bool
		errorString string

		on     func(*fields)
		assert func(*testing.T, *fields)
	}{
		{
			name: "Test 1. Positive.",
			args: args{
				ctx: ctx, // dumm
				req: kafka.Message{
					Topic: "",
					Value: msgBytes,
				},
			},

			wantErr:     false,
			errorString: "",

			on: func(f *fields) {
				f.QueueUsecase.On("CreateMessage",
					ctx,
					usecases.MessageDto{
						Id:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
						Text:    "text",
						ChatId:  "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
						OwnerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
					},
				).
					Return(nil, nil)
			},
			assert: func(t *testing.T, f *fields) {
				f.QueueUsecase.AssertNumberOfCalls(t, "CreateMessage", 1)
			},
		},
		{
			name: "Test 2. Negative. CreateMessage returns error",
			args: args{
				ctx: ctx, // dumm
				req: kafka.Message{
					Topic: "",
					Value: msgBytes,
				},
			},

			wantErr:     true,
			errorString: "fail create message: some error",

			on: func(f *fields) {
				f.QueueUsecase.On("CreateMessage",
					ctx,
					usecases.MessageDto{
						Id:      "284fef68-7e3e-4d1d-96a0-8c96f7b3b795",
						Text:    "text",
						ChatId:  "284fef68-7e3e-4d1d-96a0-8c96f7b3b000",
						OwnerId: "284fef68-7e3e-4d1d-96a0-8c96f7b3b100",
					},
				).
					Return(nil, errors.New("some error"))
			},
			assert: func(t *testing.T, f *fields) {
				f.QueueUsecase.AssertNumberOfCalls(t, "CreateMessage", 1)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			f := &fields{
				QueueUsecase:    mocks.NewQueueInterface(t),
				Cfg:             &config.Config{},
				ConsumerService: mocks.NewKafkaConsumerServiceInterface(t),
			}
			au := NewQueue(Deps{
				QueueUsecase:    f.QueueUsecase,
				Cfg:             f.Cfg,
				ConsumerService: f.ConsumerService,
			})
			if tt.on != nil {
				tt.on(f)
			}

			// act
			err := au.CreateMessage(tt.args.ctx, tt.args.req)

			// assert
			if (err != nil) != tt.wantErr {
				t.Errorf("usecase.SendUserPrivateMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if (err != nil) && tt.wantErr {
				assert.Equal(t, err.Error(), tt.errorString)
				return
			}

			if tt.assert != nil {
				tt.assert(t, f)
			}
		})
	}
}
