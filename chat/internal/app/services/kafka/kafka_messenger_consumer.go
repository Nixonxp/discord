package kafka

import (
	"context"
	config "github.com/Nixonxp/discord/chat/configs"
	"github.com/Nixonxp/discord/chat/internal/app/queue"
)

type KafkaMessengerConsumer struct {
	s *queue.KafkaMessengerConsumer
}

func (k *KafkaMessengerConsumer) Init(_ context.Context, cfg *config.Config) error {
	k.s = queue.NewKafkaMessengerConsumer(cfg)

	return nil
}

func (k *KafkaMessengerConsumer) GetInstance() *queue.KafkaMessengerConsumer {
	return k.s
}

func (k *KafkaMessengerConsumer) Ident() string {
	return "kafka messenger consumer"
}

func (k *KafkaMessengerConsumer) Close(_ context.Context) error {
	err := k.s.Close()
	if err != nil {
		return err
	}
	return nil
}
