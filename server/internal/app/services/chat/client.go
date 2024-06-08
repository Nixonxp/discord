package chat

import (
	"context"
	config "github.com/Nixonxp/discord/server/configs"
	"github.com/Nixonxp/discord/server/internal/app/services"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
	"github.com/Nixonxp/discord/server/pkg/api/chat"
	log "github.com/Nixonxp/discord/server/pkg/logger"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	ratelimitCustom "github.com/tommy-sho/rate-limiter-grpc-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type ChatClient struct {
	client chat.ChatServiceClient
	log    *log.Logger
}

var _ usecases.UsecaseChatInterface = (*ChatClient)(nil)

func (s *ChatClient) Init(ctx context.Context, cfg *config.Config) error {
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Canceled, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(3)),
		grpcretry.WithPerRetryTimeout(time.Second * 15),
	}

	chatConn, err := grpc.DialContext(ctx,
		cfg.Application.ChatServiceHost,
		grpc.WithIdleTimeout(10*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
			ratelimitCustom.UnaryClientInterceptor(ratelimitCustom.NewLimiter(10000)),
			grpc_opentracing.OpenTracingClientInterceptor(opentracing.GlobalTracer(), grpc_opentracing.LogPayloads()),
		),
		grpc.WithInsecure())
	if err != nil {
		return err
	}

	logger := services.Logger{}
	err = logger.Init(ctx, cfg)
	if err != nil {
		return err
	}

	s.client = chat.NewChatServiceClient(chatConn)
	s.log = logger.GetInstance()

	return nil
}

func (s *ChatClient) Ident() string {
	return "chat service"
}

func (s *ChatClient) GetInstance() *ChatClient {
	return s
}

func (s *ChatClient) Close(_ context.Context) error {
	return nil
}
