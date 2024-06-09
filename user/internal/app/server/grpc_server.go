package server

import (
	"context"
	"fmt"
	repositoryFriendInvte "github.com/Nixonxp/discord/user/internal/app/repository/friend_invites_storage"
	repositoryUserFriends "github.com/Nixonxp/discord/user/internal/app/repository/user_friends_storage"
	repository "github.com/Nixonxp/discord/user/internal/app/repository/user_storage"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	friend_usc "github.com/Nixonxp/discord/user/internal/app/usecases/friends"
	user_usc "github.com/Nixonxp/discord/user/internal/app/usecases/users"
	middleware "github.com/Nixonxp/discord/user/internal/middleware/errors"
	middleware_metrics "github.com/Nixonxp/discord/user/internal/middleware/metrics"
	middleware_tracing "github.com/Nixonxp/discord/user/internal/middleware/tracing"
	pb "github.com/Nixonxp/discord/user/pkg/api/v1"
	log "github.com/Nixonxp/discord/user/pkg/logger"
	"github.com/Nixonxp/discord/user/pkg/rate_limiter"
	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/go-grpc-middleware/ratelimit"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_opentracing "github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"net"
)

// Config - server config
type Config struct {
	ChainUnaryInterceptors []grpc.UnaryServerInterceptor
	UnaryInterceptors      []grpc.UnaryServerInterceptor
}

// Deps - server deps
type Deps struct {
	UserUsecase    usecases.UserUsecaseInterface
	FriendsUsecase usecases.FriendUsecaseInterface
	Log            *log.Logger
}

type UserServer struct {
	pb.UnimplementedUserServiceServer
	Deps

	validator  *protovalidate.Validator
	grpcServer *grpc.Server

	grpc struct {
		lis    net.Listener
		server *grpc.Server
	}

	http struct {
		lis  *echo.Echo
		port string
	}
}

func NewUserServer(_ context.Context, s *MainServer) (*UserServer, error) {
	srv := &UserServer{}

	// validator
	{
		validator, err := protovalidate.New(
			protovalidate.WithDisableLazy(true),
			protovalidate.WithMessages(
				&pb.CreateUserRequest{},
				&pb.GetUserForLoginRequest{},
				&pb.UpdateUserRequest{},
				&pb.GetUserByLoginRequest{},
				&pb.GetUserFriendsRequest{},
				&pb.AddToFriendByUserIdRequest{},
				&pb.AcceptFriendInviteRequest{},
				&pb.DeclineFriendInviteRequest{},
				&pb.GetUserInvitesRequest{},
				&pb.DeleteFromFriendRequest{},
				&pb.CreateOrGetUserRequest{},
			),
		)
		if err != nil {
			return nil, fmt.Errorf("server: failed to initialize validator: %w", err)
		}
		srv.validator = validator
	}

	userRepo := repository.NewUserPostgresqlRepository(s.postgres.GetInstance(), s.logger.GetInstance())
	friendInvitesRepo := repositoryFriendInvte.NewFriendInvitesPostgresqlRepository(s.postgres.GetInstance(), s.logger.GetInstance())
	friendsRepo := repositoryUserFriends.NewUserFriendsPostgresqlRepository(s.postgres.GetInstance(), s.logger.GetInstance())

	userUsecase := user_usc.NewUserUsecase(user_usc.Deps{
		UserRepo:           userRepo,
		TransactionManager: s.postgres.GetInstance(),
		Log:                s.logger.GetInstance(),
	})

	friendUsecase := friend_usc.NewFriendUsecase(friend_usc.Deps{
		UserRepo:           userRepo,
		FriendInvitesRepo:  friendInvitesRepo,
		UserFriendsRepo:    friendsRepo,
		TransactionManager: s.postgres.GetInstance(),
		Log:                s.logger.GetInstance(),
	})

	globalLimiter := rate_limiter.NewRateLimiter(10000)
	grpcConfig := Config{
		ChainUnaryInterceptors: []grpc.UnaryServerInterceptor{
			ratelimit.UnaryServerInterceptor(globalLimiter),
			grpc_recovery.UnaryServerInterceptor(),
			grpc_opentracing.OpenTracingServerInterceptor(opentracing.GlobalTracer(), grpc_opentracing.LogPayloads()),
			middleware_tracing.DebugOpenTracingUnaryServerInterceptor(true, true),
			middleware_metrics.MetricsUnaryInterceptor(),
		},
		UnaryInterceptors: []grpc.UnaryServerInterceptor{
			middleware.ErrorsUnaryInterceptor(s.logger.GetInstance()),
		},
	}

	srv.UserUsecase = userUsecase
	srv.FriendsUsecase = friendUsecase
	srv.Log = s.logger.GetInstance()

	grpcServerOptions := UnaryInterceptorsToGrpcServerOptions(grpcConfig.UnaryInterceptors...)
	grpcServerOptions = append(grpcServerOptions,
		grpc.ChainUnaryInterceptor(grpcConfig.ChainUnaryInterceptors...),
	)

	grpcServer := grpc.NewServer(grpcServerOptions...)
	pb.RegisterUserServiceServer(grpcServer, srv)

	srv.grpcServer = grpcServer

	return srv, nil
}

func UnaryInterceptorsToGrpcServerOptions(interceptors ...grpc.UnaryServerInterceptor) []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0, len(interceptors))
	for _, interceptor := range interceptors {
		opts = append(opts, grpc.UnaryInterceptor(interceptor))
	}
	return opts
}
