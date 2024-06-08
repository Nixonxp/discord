package services

import (
	"context"
	"fmt"
	config "github.com/Nixonxp/discord/chat/configs"
	"github.com/segmentio/kafka-go"
)

type Kafka struct {
	conn *kafka.Conn
}

func (k *Kafka) Init(ctx context.Context, cfg *config.Config) error {
	var err error
	k.conn, err = kafka.DialLeader(ctx, "tcp", cfg.Application.KafkaAddress, cfg.Application.KafkaMessagesTopic, 0)
	if err != nil {
		return fmt.Errorf("failed to dial leader:", err)
	}

	return nil
}

func (k *Kafka) GetInstance() *kafka.Conn {
	return k.conn
}

func (k *Kafka) Ident() string {
	return "kafka"
}

func (k *Kafka) Close(_ context.Context) error {
	err := k.conn.Close()
	if err != nil {
		return err
	}
	return nil
}
