package queue

import (
	"encoding/json"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	pkgerrors "github.com/Nixonxp/discord/chat/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type KafkaMessenger struct {
	conn *kafka.Conn
}

func NewKafkaMessenger(conn *kafka.Conn) *KafkaMessenger {
	return &KafkaMessenger{
		conn: conn,
	}
}

func (m *KafkaMessenger) SendMessage(msgData usecases.MessageDto) error {
	msgBytes, _ := json.Marshal(msgData)

	_, err := m.conn.WriteMessages(
		kafka.Message{Value: msgBytes},
	)
	if err != nil {
		return pkgerrors.Wrap("failed to send messages", err)
	}

	return nil
}
