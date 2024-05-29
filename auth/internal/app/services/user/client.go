package user

import (
	"context"
	"github.com/Nixonxp/discord/auth/internal/app/server"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	"github.com/Nixonxp/discord/auth/pkg/api/user"
	log "github.com/Nixonxp/discord/auth/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type UserClient struct {
	client user.UserServiceClient
	log    *log.Logger
}

var _ usecases.UsecaseServiceInterface = (*UserClient)(nil)

func NewClient(cfg server.Config, log *log.Logger) (*UserClient, error) {
	userConn, err := grpc.DialContext(context.Background(),
		cfg.UserServiceUrl,
		grpc.WithIdleTimeout(10*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			cfg.UnaryClientInterceptors...,
		),
		grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &UserClient{
		client: user.NewUserServiceClient(userConn),
		log:    log,
	}, nil
}
