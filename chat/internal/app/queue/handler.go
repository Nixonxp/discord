package queue

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	"github.com/segmentio/kafka-go"
	"log"
)

type Deps struct {
	QueueUsecase usecases.QueueInterface
}

type Queue struct {
	KafkaReader *kafka.Reader
	Deps
}

func NewQueue(d Deps) *Queue {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "messages",
		GroupID:  "group-1",
		MaxBytes: 10e6, // 10MB
	})
	return &Queue{
		KafkaReader: r,
		Deps:        d,
	}
}

func (q *Queue) Run(ctx context.Context) error {
	var err error
	for {
		m, err := q.KafkaReader.ReadMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				break
			}
			log.Printf("failed to read messages: %s", err)
			break
		}
		if m.Topic == "messages" {
			err = q.CreateMessage(ctx, m)
			if err != nil {
				break
			}
		}
		fmt.Printf("message at topic/partition/offset %v/%v/%v: %s = %s\n", m.Topic, m.Partition, m.Offset, string(m.Key), string(m.Value))
		/*err = q.KafkaReader.CommitMessages(ctx, m)
		if err != nil {
			err = fmt.Errorf("failed to commit messages:", err)
			break
		}*/
	}

	if err != nil {
		return err
	}

	defer log.Println("consumer closed")

	defer func() {
		if err := q.KafkaReader.Close(); err != nil {
			log.Fatalf("failed to close reader:", err)
		}
	}()

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
	msgDto := MessageKafkaMessage{}
	err := q.ReadInMessage(message, &msgDto)
	if err != nil {
		return err
	}

	log.Println("get message - " + msgDto.Text)

	_, err = q.QueueUsecase.CreateMessage(ctx, usecases.MessageDto{
		Id:      msgDto.Id,
		Text:    msgDto.Text,
		ChatId:  msgDto.ChatId,
		OwnerId: msgDto.OwnerId,
	})
	if err != nil {
		return err
	}

	return nil
}
