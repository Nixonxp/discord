package queue

import (
	"context"
	config "github.com/Nixonxp/discord/chat/configs"
	"github.com/segmentio/kafka-go"
	"log"
)

type KafkaMessengerConsumer struct {
	r *kafka.Reader
}

func NewKafkaMessengerConsumer(cfg *config.Config) *KafkaMessengerConsumer {
	log.Printf("start listening queue from %s", cfg.Application.KafkaAddress)
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.Application.KafkaAddress},
		Topic:    cfg.Application.KafkaMessagesTopic,
		GroupID:  "group-1",
		MaxBytes: 10e6, // 10MB
	})

	return &KafkaMessengerConsumer{
		r: r,
	}
}

func (m *KafkaMessengerConsumer) ReadMessage(ctx context.Context) (kafka.Message, error) {
	return m.r.ReadMessage(ctx)
}

func (m *KafkaMessengerConsumer) Close() error {
	return m.r.Close()
}
