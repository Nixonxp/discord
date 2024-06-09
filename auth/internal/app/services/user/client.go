package services

import (
	"context"
	config "github.com/Nixonxp/discord/auth/configs"
	"github.com/Nixonxp/discord/auth/internal/app/services"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	"github.com/Nixonxp/discord/auth/pkg/api/user"
	log "github.com/Nixonxp/discord/auth/pkg/logger"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	ratelimitCustom "github.com/tommy-sho/rate-limiter-grpc-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type UserClient struct {
	client user.UserServiceClient
	log    *log.Logger
}

var _ usecases.UserServiceInterface = (*UserClient)(nil)

func (s *UserClient) Init(ctx context.Context, cfg *config.Config) error {
	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Canceled, codes.Aborted, codes.DeadlineExceeded),
		grpcretry.WithMax(uint(3)),
		grpcretry.WithPerRetryTimeout(time.Second * 15),
	}

	userConn, err := grpc.DialContext(ctx,
		cfg.Application.UserServiceHost,
		grpc.WithIdleTimeout(10*time.Second),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpcretry.UnaryClientInterceptor(retryOpts...),
			ratelimitCustom.UnaryClientInterceptor(ratelimitCustom.NewLimiter(1000)),
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

	s.client = user.NewUserServiceClient(userConn)
	s.log = logger.GetInstance()

	return nil
}

func (s *UserClient) Ident() string {
	return "user service"
}

func (s *UserClient) GetInstance() *UserClient {
	return s
}

func (s *UserClient) Close(_ context.Context) error {
	return nil
}
