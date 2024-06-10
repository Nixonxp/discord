package kafka

import (
	"context"
	"fmt"
	config "github.com/Nixonxp/discord/chat/configs"
	"github.com/Nixonxp/discord/chat/internal/app/queue"
	"github.com/segmentio/kafka-go"
)

type KafkaMessengerProducer struct {
	conn *kafka.Conn
	mess *queue.KafkaMessenger
}

func (k *KafkaMessengerProducer) Init(ctx context.Context, cfg *config.Config) error {
	var err error
	k.conn, err = kafka.DialLeader(ctx, "tcp", cfg.Application.KafkaAddress, cfg.Application.KafkaMessagesTopic, 0)
	if err != nil {
		return fmt.Errorf("failed to dial leader:", err)
	}

	k.mess = queue.NewKafkaMessenger(k.conn)

	return nil
}

func (k *KafkaMessengerProducer) GetInstance() *queue.KafkaMessenger {
	return k.mess
}

func (k *KafkaMessengerProducer) Ident() string {
	return "kafka messenger producer"
}

func (k *KafkaMessengerProducer) Close(_ context.Context) error {
	err := k.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
