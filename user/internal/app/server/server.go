package server

import (
	"context"
	"fmt"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	pb "github.com/Nixonxp/discord/user/pkg/api/v1"
	log "github.com/Nixonxp/discord/user/pkg/logger"
	"github.com/bufbuild/protovalidate-go"
	"github.com/labstack/echo/v4"
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
	UserUsecase usecases.UsecaseInterface
	Log         *log.Logger
}

type UserServer struct {
	pb.UnimplementedUserServiceServer
	Deps

	validator *protovalidate.Validator

	grpc struct {
		lis    net.Listener
		server *grpc.Server
	}

	http struct {
		lis  *echo.Echo
		port string
	}
}

func NewUserServer(ctx context.Context, d Deps) (*UserServer, error) {
	srv := &UserServer{
		Deps: d,
	}

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

	return srv, nil
}

func UnaryInterceptorsToGrpcServerOptions(interceptors ...grpc.UnaryServerInterceptor) []grpc.ServerOption {
	opts := make([]grpc.ServerOption, 0, len(interceptors))
	for _, interceptor := range interceptors {
		opts = append(opts, grpc.UnaryInterceptor(interceptor))
	}
	return opts
}
