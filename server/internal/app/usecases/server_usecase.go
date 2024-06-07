package usecases

import (
	"context"
	"errors"
	"github.com/Nixonxp/discord/server/internal/app/models"
	log "github.com/Nixonxp/discord/server/pkg/logger"
	"github.com/google/uuid"
)

type Deps struct {
	ServerRepo    ServerStorage
	SubscribeRepo SubscribeStorage
	ChatService   UsecaseChatInterface
	Log           *log.Logger
}

type ServerUsecase struct {
	Deps
}

var _ UsecaseInterface = (*ServerUsecase)(nil)

func NewServerUsecase(d Deps) UsecaseInterface {
	return &ServerUsecase{
		Deps: d,
	}
}

func (u *ServerUsecase) CreateServer(ctx context.Context, req CreateServerRequest) (*models.ServerInfo, error) {
	chatID := models.ServerID(uuid.New())
	userID := models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) // todo from auth

	newServer := models.ServerInfo{
		Id:      chatID,
		Name:    req.Name,
		OwnerId: userID,
	}
	err := u.ServerRepo.CreateServer(ctx, newServer)

	if err != nil {
		return &models.ServerInfo{}, err
	}

	return &newServer, nil
}

func (u *ServerUsecase) SearchServer(ctx context.Context, req SearchServerRequest) ([]*models.ServerInfo, error) {
	servers, err := u.ServerRepo.SearchServers(ctx, req.Name)
	if err != nil {
		return nil, err
	}

	if len(servers) == 0 {
		return nil, models.ErrNotFound
	}

	return servers, nil
}

func (u *ServerUsecase) SubscribeServer(ctx context.Context, req SubscribeServerRequest) (*models.ActionInfo, error) {
	_, err := u.ServerRepo.GetServerById(ctx, req.ServerId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.New("server not found")
		}
		return nil, err
	}

	userID := models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) // todo from auth
	serverID := models.ServerID(uuid.MustParse(req.ServerId))

	newSubscribe := models.SubscribeInfo{
		Id:       models.SubscribeID(uuid.New()),
		ServerId: serverID,
		UserId:   userID,
	}
	// todo create index to unique subscribes
	err = u.SubscribeRepo.CreateSubscribe(ctx, newSubscribe)
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ServerUsecase) UnsubscribeServer(ctx context.Context, req UnsubscribeServerRequest) (*models.ActionInfo, error) {
	userID := models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) // todo from auth
	serverID := models.ServerID(uuid.MustParse(req.ServerId))

	err := u.SubscribeRepo.DeleteSubscribe(ctx, serverID, userID)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.New("subscribe not found")
		}
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ServerUsecase) SearchServerByUserId(ctx context.Context, req SearchServerByUserIdRequest) ([]string, error) {
	userID := models.UserID(uuid.MustParse(req.UserId))
	subscribes, err := u.SubscribeRepo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, err
	}

	if len(subscribes) == 0 {
		return nil, models.ErrNotFound
	}

	userServers := make([]string, len(subscribes))
	for i, s := range subscribes {
		userServers[i] = s.ServerId.String()
	}

	return userServers, nil
}

func (u *ServerUsecase) InviteUserToServer(ctx context.Context, req InviteUserToServerRequest) (*models.ActionInfo, error) {
	_, err := u.ServerRepo.GetServerById(ctx, req.ServerId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.New("server not found")
		}
		return nil, err
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
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ServerUsecase) PublishMessageOnServer(ctx context.Context, req PublishMessageOnServerRequest) (*models.ActionInfo, error) {
	_, err := u.ServerRepo.GetServerById(ctx, req.ServerId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.New("server not found")
		}
		return nil, err
	}

	_, err = u.ChatService.SendServerMessage(ctx, PublishMessageOnServerRequest{
		ServerId: req.ServerId,
		Text:     req.Text,
	})
	if err != nil {
		return nil, err
	}

	return &models.ActionInfo{
		Success: true,
	}, nil
}

func (u *ServerUsecase) GetMessagesFromServer(ctx context.Context, req GetMessagesFromServerRequest) (*models.GetMessagesInfo, error) {
	_, err := u.ServerRepo.GetServerById(ctx, req.ServerId)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, errors.New("server not found")
		}
		return nil, err
	}

	messages, err := u.ChatService.GetServerMessages(ctx, GetMessagesFromServerRequest{
		ServerId: req.ServerId,
	})
	if err != nil {
		return nil, err
	}

	return &models.GetMessagesInfo{
		Messages: messages.Messages,
	}, nil
}
