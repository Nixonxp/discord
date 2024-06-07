package chat

import (
	"context"
	"github.com/Nixonxp/discord/server/internal/app/server"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
	"github.com/Nixonxp/discord/server/pkg/api/chat"
	log "github.com/Nixonxp/discord/server/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type ChatClient struct {
	client chat.ChatServiceClient
	log    *log.Logger
}

var _ usecases.UsecaseChatInterface = (*ChatClient)(nil)

func NewClient(cfg server.Config, log *log.Logger) (*ChatClient, error) {
	conn, err := grpc.DialContext(context.Background(),
		cfg.ChatServiceUrl,
		grpc.WithIdleTimeout(10*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			cfg.UnaryClientInterceptors...,
		),
		grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &ChatClient{
		client: chat.NewChatServiceClient(conn),
		log:    log,
	}, nil
}
