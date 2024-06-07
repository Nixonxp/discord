package usecases

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Nixonxp/discord/chat/internal/app/enum"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
	"sort"
	"strings"
)

type Deps struct {
	MessagesRepo MessagesStorage
	ChatRepo     ChatStorage
	KafkaConn    *kafka.Conn
}

type ChatUsecase struct {
	Deps
}

var _ UsecaseInterface = (*ChatUsecase)(nil)

func NewChatUsecase(d Deps) UsecaseInterface {
	return &ChatUsecase{
		Deps: d,
	}
}

// мок, в будущем будет браться из таблицы
const chatIdStringMock = "3d0222e1-4b58-4fa7-a38c-171ee345b14e"

func (u *ChatUsecase) SendUserPrivateMessage(ctx context.Context, req SendUserPrivateMessageRequest) (*models.ActionInfo, error) {
	dataSlice := []string{req.CurrentUser, req.UserId}
	sort.StringsAreSorted(dataSlice)

	meta := strings.Join(dataSlice, "_")
	existChat, err := u.ChatRepo.GetChatByMetadataAndType(ctx, meta, enum.PrivateChatType)
	if err != nil {
		return nil, err
	}

	err = u.send(MessageDto{
		Id:      uuid.New().String(),
		ChatId:  existChat.Id.String(),
		OwnerId: req.CurrentUser,
		Text:    req.Text,
	})
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ChatUsecase) send(msgData MessageDto) error {
	msgBytes, _ := json.Marshal(msgData)

	_, err := u.KafkaConn.WriteMessages(
		kafka.Message{Value: msgBytes},
	)
	if err != nil {
		return fmt.Errorf("failed to send messages", err)
	}

	return nil
}

func (u *ChatUsecase) SendServerMessage(ctx context.Context, req SendServerMessageRequest) (*models.ActionInfo, error) {
	currentChat, err := u.ChatRepo.GetChatByMetadataAndType(ctx, req.ServerId, enum.ServerChatType)
	if err != nil {
		if !errors.Is(models.ErrNotFound, err) {
			return nil, err
		}

		chatId := models.ChatID(uuid.New())
		currentChat = &models.Chat{
			Id:       chatId,
			Type:     enum.ServerChatType,
			OwnerId:  models.OwnerID(uuid.MustParse(req.ServerId)),
			MetaData: req.ServerId,
		}
		err := u.ChatRepo.CreateChat(ctx, currentChat)
		if err != nil {
			return nil, err
		}
	}

	chatId := models.MessageID(uuid.New())
	err = u.send(MessageDto{
		Id:      chatId.String(),
		ChatId:  currentChat.Id.String(),
		OwnerId: req.ServerId,
		Text:    req.Text,
	})
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ChatUsecase) GetUserPrivateMessages(ctx context.Context, req GetUserPrivateMessagesRequest) (*models.Messages, error) {
	dataSlice := []string{req.CurrentUser, req.UserId}
	sort.StringsAreSorted(dataSlice)

	meta := strings.Join(dataSlice, "_")
	existChat, err := u.ChatRepo.GetChatByMetadataAndType(ctx, meta, enum.PrivateChatType)
	if err != nil {
		return nil, err
	}

	messages, err := u.MessagesRepo.GetMessages(ctx, existChat.Id)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, models.ErrEmpty
	}

	return &models.Messages{
		Data: messages,
	}, nil
}

func (u *ChatUsecase) GetServerMessagesRequest(ctx context.Context, req GetServerMessageRequest) (*models.Messages, error) {
	existChat, err := u.ChatRepo.GetChatByMetadataAndType(ctx, req.ServerId, enum.ServerChatType)
	if err != nil {
		return nil, err
	}

	messages, err := u.MessagesRepo.GetMessages(ctx, existChat.Id)
	if err != nil {
		return nil, err
	}

	if len(messages) == 0 {
		return nil, models.ErrEmpty
	}

	return &models.Messages{
		Data: messages,
	}, nil
}

func (u *ChatUsecase) CreatePrivateChat(ctx context.Context, req CreatePrivateChatRequest) (*models.Chat, error) {
	dataSlice := []string{req.CurrentUser, req.UserId}
	sort.StringsAreSorted(dataSlice)

	meta := strings.Join(dataSlice, "_")
	existChat, err := u.ChatRepo.GetChatByMetadataAndType(ctx, meta, enum.PrivateChatType)
	if err != nil {
		if !errors.Is(models.ErrNotFound, err) {
			return nil, err
		}

		chatId := models.ChatID(uuid.New())
		chat := &models.Chat{
			Id:       chatId,
			Type:     enum.PrivateChatType,
			OwnerId:  models.OwnerID(uuid.MustParse(req.CurrentUser)),
			MetaData: meta,
		}
		err := u.ChatRepo.CreateChat(ctx, chat)
		if err != nil {
			return nil, err
		}

		return chat, err
	}

	return existChat, nil
}
