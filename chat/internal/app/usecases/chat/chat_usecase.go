package chat

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/chat/internal/app/enum"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	pkgerrors "github.com/Nixonxp/discord/chat/pkg/errors"
	"github.com/google/uuid"
	"sort"
	"strings"
)

type Deps struct {
	MessagesRepo usecases.MessagesStorage
	ChatRepo     usecases.ChatStorage
	KafkaConn    usecases.KafkaProducerServiceInterface
}

type ChatUsecase struct {
	Deps
}

var _ usecases.UsecaseInterface = (*ChatUsecase)(nil)

func NewChatUsecase(d Deps) usecases.UsecaseInterface {
	return &ChatUsecase{
		Deps: d,
	}
}

func (u *ChatUsecase) SendUserPrivateMessage(ctx context.Context, req usecases.SendUserPrivateMessageRequest) (*models.ActionInfo, error) {
	dataSlice := []string{req.CurrentUser, req.UserId}
	sort.StringsAreSorted(dataSlice)

	meta := strings.Join(dataSlice, "_")
	existChat, err := u.ChatRepo.GetChatByMetadataAndType(ctx, meta, enum.PrivateChatType)
	if err != nil {
		return nil, pkgerrors.Wrap("chat search error", err)
	}

	err = u.KafkaConn.SendMessage(usecases.MessageDto{
		ChatId:  existChat.Id.String(),
		OwnerId: req.CurrentUser,
		Text:    req.Text,
	})
	if err != nil {
		return nil, pkgerrors.Wrap("chat message send error", err)
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ChatUsecase) SendServerMessage(ctx context.Context, req usecases.SendServerMessageRequest) (*models.ActionInfo, error) {
	currentChat, err := u.ChatRepo.GetChatByMetadataAndType(ctx, req.ServerId, enum.ServerChatType)
	if err != nil {
		if !errors.Is(models.ErrNotFound, err) {
			return nil, pkgerrors.Wrap("chat search error", err)
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
			return nil, pkgerrors.Wrap("create new chat for server", err)
		}
	}

	err = u.KafkaConn.SendMessage(usecases.MessageDto{
		ChatId:  currentChat.Id.String(),
		OwnerId: req.ServerId,
		Text:    req.Text,
	})
	if err != nil {
		return nil, pkgerrors.Wrap("send message to server server", err)
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ChatUsecase) GetUserPrivateMessages(ctx context.Context, req usecases.GetUserPrivateMessagesRequest) (*models.Messages, error) {
	dataSlice := []string{req.CurrentUser, req.UserId}
	sort.StringsAreSorted(dataSlice)

	meta := strings.Join(dataSlice, "_")
	existChat, err := u.ChatRepo.GetChatByMetadataAndType(ctx, meta, enum.PrivateChatType)
	if err != nil {
		return nil, pkgerrors.Wrap("get chat error", err)
	}

	messages, err := u.MessagesRepo.GetMessages(ctx, existChat.Id)
	if err != nil {
		return nil, pkgerrors.Wrap("get messages error", err)
	}

	if len(messages) == 0 {
		return nil, pkgerrors.Wrap("get messages error", models.ErrEmpty)
	}

	return &models.Messages{
		Data: messages,
	}, nil
}

func (u *ChatUsecase) GetServerMessagesRequest(ctx context.Context, req usecases.GetServerMessageRequest) (*models.Messages, error) {
	existChat, err := u.ChatRepo.GetChatByMetadataAndType(ctx, req.ServerId, enum.ServerChatType)
	if err != nil {
		return nil, pkgerrors.Wrap("get server chat error", err)
	}

	messages, err := u.MessagesRepo.GetMessages(ctx, existChat.Id)
	if err != nil {
		return nil, pkgerrors.Wrap("get server messages error", err)
	}

	if len(messages) == 0 {
		return nil, pkgerrors.Wrap("get server messages error", models.ErrEmpty)
	}

	return &models.Messages{
		Data: messages,
	}, nil
}

func (u *ChatUsecase) CreatePrivateChat(ctx context.Context, req usecases.CreatePrivateChatRequest) (*models.Chat, error) {
	dataSlice := []string{req.CurrentUser, req.UserId}
	sort.StringsAreSorted(dataSlice)

	meta := strings.Join(dataSlice, "_")
	existChat, err := u.ChatRepo.GetChatByMetadataAndType(ctx, meta, enum.PrivateChatType)
	if err != nil {
		if !errors.Is(models.ErrNotFound, err) {
			return nil, pkgerrors.Wrap("get chat error", models.ErrNotFound)
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
			return nil, pkgerrors.Wrap("crate chat error", err)
		}

		return chat, pkgerrors.Wrap("crate chat error", err)
	}

	return existChat, nil
}
