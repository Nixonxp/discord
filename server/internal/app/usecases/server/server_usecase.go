package server

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/server/internal/app/models"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
	pkgerrors "github.com/Nixonxp/discord/server/pkg/errors"
	log "github.com/Nixonxp/discord/server/pkg/logger"
	"github.com/google/uuid"
)

type Deps struct {
	ServerRepo    usecases.ServerStorage
	SubscribeRepo usecases.SubscribeStorage
	ChatService   usecases.ServiceChatInterface
	Log           *log.Logger
}

type ServerUsecase struct {
	Deps
}

var _ usecases.UsecaseInterface = (*ServerUsecase)(nil)

func NewServerUsecase(d Deps) usecases.UsecaseInterface {
	return &ServerUsecase{
		Deps: d,
	}
}

func (u *ServerUsecase) CreateServer(ctx context.Context, req usecases.CreateServerRequest) (*models.ServerInfo, error) {
	chatID := models.ServerID(uuid.New())
	userID := models.UserID(uuid.MustParse(req.CurrentUserId))

	newServer := models.ServerInfo{
		Id:      chatID,
		Name:    req.Name,
		OwnerId: userID,
	}
	err := u.ServerRepo.CreateServer(ctx, newServer)

	if err != nil {
		return nil, pkgerrors.Wrap("create server error", err)
	}

	return &newServer, nil
}

func (u *ServerUsecase) SearchServer(ctx context.Context, req usecases.SearchServerRequest) ([]*models.ServerInfo, error) {
	servers, err := u.ServerRepo.SearchServers(ctx, req.Name)
	if err != nil {
		return nil, pkgerrors.Wrap("search server error", err)
	}

	if len(servers) == 0 {
		return nil, pkgerrors.Wrap("servers", models.ErrNotFound)
	}

	return servers, nil
}

func (u *ServerUsecase) SubscribeServer(ctx context.Context, req usecases.SubscribeServerRequest) (*models.ActionInfo, error) {
	_, err := u.ServerRepo.GetServerById(ctx, req.ServerId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgerrors.Wrap("subscribe server error", models.ErrNotFound)
		}
		return nil, pkgerrors.Wrap("subscribe server error", err)
	}

	userID := models.UserID(uuid.MustParse(req.CurrentUserId))
	serverID := models.ServerID(uuid.MustParse(req.ServerId))

	newSubscribe := models.SubscribeInfo{
		Id:       models.SubscribeID(uuid.New()),
		ServerId: serverID,
		UserId:   userID,
	}
	// todo create index to unique subscribes
	err = u.SubscribeRepo.CreateSubscribe(ctx, newSubscribe)
	if err != nil {
		return nil, pkgerrors.Wrap("subscribe create error", err)
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ServerUsecase) UnsubscribeServer(ctx context.Context, req usecases.UnsubscribeServerRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse(req.CurrentUserId))
	serverID := models.ServerID(uuid.MustParse(req.ServerId))

	err := u.SubscribeRepo.DeleteSubscribe(ctx, serverID, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgerrors.Wrap("unsubscribe error", models.ErrNotFound)
		}
		return nil, pkgerrors.Wrap("unsubscribe error", err)
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ServerUsecase) SearchServerByUserId(ctx context.Context, req usecases.SearchServerByUserIdRequest) ([]string, error) {
	userID := models.UserID(uuid.MustParse(req.UserId))
	subscribes, err := u.SubscribeRepo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, pkgerrors.Wrap("search server error", err)
	}

	if len(subscribes) == 0 {
		return nil, pkgerrors.Wrap("servers", models.ErrNotFound)
	}

	userServers := make([]string, len(subscribes))
	for i, s := range subscribes {
		userServers[i] = s.ServerId.String()
	}

	return userServers, nil
}

func (u *ServerUsecase) InviteUserToServer(ctx context.Context, req usecases.InviteUserToServerRequest) (*models.ActionInfo, error) {
	_, err := u.ServerRepo.GetServerById(ctx, req.ServerId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgerrors.Wrap("servers", models.ErrNotFound)
		}
		return nil, pkgerrors.Wrap("server search", err)
	}

	userID := models.UserID(uuid.MustParse(req.UserId))
	serverID := models.ServerID(uuid.MustParse(req.ServerId))

	newSubscribe := models.SubscribeInfo{
		Id:       models.SubscribeID(uuid.New()),
		ServerId: serverID,
		UserId:   userID,
	}
	// todo create index to unique subscribes
	err = u.SubscribeRepo.CreateSubscribe(ctx, newSubscribe)
	if err != nil {
		return nil, pkgerrors.Wrap("create subscribe", err)
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ServerUsecase) PublishMessageOnServer(ctx context.Context, req usecases.PublishMessageOnServerRequest) (*models.ActionInfo, error) {
	_, err := u.ServerRepo.GetServerById(ctx, req.ServerId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgerrors.Wrap("publish message on server", models.ErrNotFound)
		}
		return nil, pkgerrors.Wrap("publish message on server", err)
	}

	_, err = u.ChatService.SendServerMessage(ctx, usecases.PublishMessageOnServerRequest{
		ServerId: req.ServerId,
		Text:     req.Text,
	})
	if err != nil {
		return nil, pkgerrors.Wrap("publish message on server", err)
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ServerUsecase) GetMessagesFromServer(ctx context.Context, req usecases.GetMessagesFromServerRequest) (*models.GetMessagesInfo, error) {
	_, err := u.ServerRepo.GetServerById(ctx, req.ServerId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, pkgerrors.Wrap("get messages on server", models.ErrNotFound)
		}
		return nil, pkgerrors.Wrap("get messages on server", err)
	}

	messages, err := u.ChatService.GetServerMessages(ctx, usecases.GetMessagesFromServerRequest{
		ServerId: req.ServerId,
	})
	if err != nil {
		return nil, pkgerrors.Wrap("get messages on server", err)
	}

	return &models.GetMessagesInfo{
		Messages: messages.Messages,
	}, nil
}
