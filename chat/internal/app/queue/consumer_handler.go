package queue

import (
	"context"
	"encoding/json"
	"errors"
	config "github.com/Nixonxp/discord/chat/configs"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	pkgerrors "github.com/Nixonxp/discord/chat/pkg/errors"
	log "github.com/Nixonxp/discord/chat/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type Deps struct {
	QueueUsecase    usecases.QueueInterface
	Cfg             *config.Config
	ConsumerService usecases.KafkaConsumerServiceInterface
	Log             *log.Logger
}

type Queue struct {
	KafkaReader *kafka.Reader
	Deps
}

func NewQueue(d Deps) *Queue {
	return &Queue{
		Deps: d,
	}
}

func (q *Queue) Run(ctx context.Context) error {
	for {
		m, err := q.ConsumerService.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			q.Log.WithContext(ctx).WithError(err).Error("failed to read messages")
			break
		}

		q.Log.WithContext(ctx).Info("get  message from kafka")
		if m.Topic == "messages" {
			err = q.CreateMessage(ctx, m)
			if err != nil {
				break
			}
		}
	}

	return nil
}

func (q *Queue) ReadInMessage(message kafka.Message, in any) error {
	err := json.Unmarshal(message.Value, in)
	if err != nil {
		return err
	}

	return nil
}

func (q *Queue) CreateMessage(ctx context.Context, message kafka.Message) error {
	q.Log.WithContext(ctx).Infof("get message from kafka topic - %s", q.Cfg.Application.KafkaMessagesTopic)
	msgDto := MessageKafkaMessage{}
	err := q.ReadInMessage(message, &msgDto)
	if err != nil {
		return pkgerrors.Wrap("unmarshal message", err)
	}

	_, err = q.QueueUsecase.CreateMessage(ctx, usecases.MessageDto{
		Id:      msgDto.Id,
		Text:    msgDto.Text,
		ChatId:  msgDto.ChatId,
		OwnerId: msgDto.OwnerId,
	})
	if err != nil {
		return pkgerrors.Wrap("fail create message", err)
	}

	return nil
}
