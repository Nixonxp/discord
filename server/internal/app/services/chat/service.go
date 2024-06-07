package chat

import (
	"context"
	"github.com/Nixonxp/discord/server/internal/app/models"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
	"github.com/Nixonxp/discord/server/pkg/api/chat"
	"github.com/opentracing/opentracing-go"
	"time"
)

func (s *ChatClient) SendServerMessage(ctx context.Context, msg usecases.PublishMessageOnServerRequest) (*models.ActionInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "chat_service.SendServerMessage")
	defer span.Finish()
	response, err := s.client.SendServerMessage(ctx, &chat.SendServerMessageRequest{
		ServerId: msg.ServerId,
		Text:     msg.Text,
	})
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{
		Success: response.GetSuccess(),
	}, nil
}

func (s *ChatClient) GetServerMessages(ctx context.Context, msg usecases.GetMessagesFromServerRequest) (*models.GetMessagesInfo, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "chat_service.GetServerMessages")
	defer span.Finish()
	response, err := s.client.GetServerMessages(ctx, &chat.GetServerMessagesRequest{
		ServerId: msg.ServerId,
	})
	if err != nil {
		return nil, err
	}

	messages := make([]*models.Message, len(response.GetMessages()))
	for k, v := range response.GetMessages() {
		messages[k] = &models.Message{
			Id:        v.Id,
			Text:      v.Text,
			Timestamp: time.Unix(v.GetTimestamp().Seconds, int64(v.GetTimestamp().GetNanos())),
		}
	}

	return &models.GetMessagesInfo{
		Messages: messages,
	}, nil
}
