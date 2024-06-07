package server

import (
	"context"
	"fmt"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
	pb "github.com/Nixonxp/discord/server/pkg/api/v1"
	grpcutils "github.com/Nixonxp/discord/server/pkg/grpc_utils"
	log "github.com/Nixonxp/discord/server/pkg/logger"
	"github.com/bufbuild/protovalidate-go"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net"
)

// Config - server config
type Config struct {
	ChainUnaryInterceptors  []grpc.UnaryServerInterceptor
	UnaryInterceptors       []grpc.UnaryServerInterceptor
	UnaryClientInterceptors []grpc.UnaryClientInterceptor

	ChatServiceUrl string
}

// Deps - server deps
type Deps struct {
	ServerUsecase usecases.UsecaseInterface
	Log           *log.Logger
}

type ServerServer struct {
	pb.UnimplementedServerServiceServer
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

func NewServerServer(ctx context.Context, d Deps) (*ServerServer, error) {
	srv := &ServerServer{
		Deps: d,
	}

	// validator
	{
		validator, err := protovalidate.New(
			protovalidate.WithDisableLazy(true),
			protovalidate.WithMessages(
				&pb.CreateServerRequest{},
				&pb.SearchServerRequest{},
				&pb.SubscribeServerRequest{},
				&pb.UnsubscribeServerRequest{},
				&pb.SearchServerByUserIdRequest{},
				&pb.InviteUserToServerRequest{},
				&pb.PublishMessageOnServerRequest{},
				&pb.GetMessagesFromServerRequest{},
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

func (s *ServerServer) CreateServer(ctx context.Context, req *pb.CreateServerRequest) (*pb.CreateServerResponse, error) {
	s.Log.WithContext(ctx).WithField("name", req.GetName()).Info("Create server: received")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	created, err := s.ServerUsecase.CreateServer(ctx, usecases.CreateServerRequest{
		Name: req.GetName(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreateServerResponse{
		Id:      created.Id.String(),
		Name:    created.Name,
		OwnerId: created.OwnerId.String(),
	}, nil
}

func (s *ServerServer) SearchServer(ctx context.Context, req *pb.SearchServerRequest) (*pb.SearchServerResponse, error) {
	s.Log.WithContext(ctx).WithField("name", req.GetName()).Info("Search server: received")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ServerUsecase.SearchServer(ctx, usecases.SearchServerRequest{
		Name: req.GetName(),
	})
	if err != nil {
		return nil, err
	}

	servers := make([]*pb.ServerInfo, len(result))
	for i, srv := range result {
		servers[i] = &pb.ServerInfo{
			Id:   srv.Id.String(),
			Name: srv.Name,
		}
	}

	return &pb.SearchServerResponse{
		Servers: servers,
	}, nil
}

func (s *ServerServer) SubscribeServer(ctx context.Context, req *pb.SubscribeServerRequest) (*pb.ActionResponse, error) {
	s.Log.WithContext(ctx).WithField("server id", req.GetServerId()).Info("Subscribe server: received")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ServerUsecase.SubscribeServer(ctx, usecases.SubscribeServerRequest{
		ServerId: req.ServerId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *ServerServer) UnsubscribeServer(ctx context.Context, req *pb.UnsubscribeServerRequest) (*pb.ActionResponse, error) {
	s.Log.WithContext(ctx).WithField("server id", req.GetServerId()).Info("Unsubscribe server: received")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ServerUsecase.UnsubscribeServer(ctx, usecases.UnsubscribeServerRequest{
		ServerId: req.ServerId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *ServerServer) SearchServerByUserId(ctx context.Context, req *pb.SearchServerByUserIdRequest) (*pb.SearchServerByUserIdResponse, error) {
	s.Log.WithContext(ctx).WithField("user id", req.GetUserId()).Info("Search server by user id")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ServerUsecase.SearchServerByUserId(ctx, usecases.SearchServerByUserIdRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.SearchServerByUserIdResponse{
		Id: result,
	}, nil
}

func (s *ServerServer) InviteUserToServer(ctx context.Context, req *pb.InviteUserToServerRequest) (*pb.ActionResponse, error) {
	s.Log.WithContext(ctx).WithField("server id", req.GetServerId()).Info("Invite user to server")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ServerUsecase.InviteUserToServer(ctx, usecases.InviteUserToServerRequest{
		UserId:   req.UserId,
		ServerId: req.ServerId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *ServerServer) PublishMessageOnServer(ctx context.Context, req *pb.PublishMessageOnServerRequest) (*pb.ActionResponse, error) {
	s.Log.WithContext(ctx).WithField("text", req.GetText()).Info("Publish message to server")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ServerUsecase.PublishMessageOnServer(ctx, usecases.PublishMessageOnServerRequest{
		ServerId: req.ServerId,
		Text:     req.Text,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *ServerServer) GetMessagesFromServer(ctx context.Context, req *pb.GetMessagesFromServerRequest) (*pb.GetMessagesResponse, error) {
	s.Log.WithContext(ctx).WithField("server id", req.GetServerId()).Info("Get messages from server")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ServerUsecase.GetMessagesFromServer(ctx, usecases.GetMessagesFromServerRequest{
		ServerId: req.ServerId,
	})
	if err != nil {
		return nil, err
	}

	messages := make([]*pb.Message, len(result.Messages))
	for k, v := range result.Messages {
		messages[k] = &pb.Message{
			Id:   v.Id,
			Text: v.Text,
			Timestamp: &timestamppb.Timestamp{
				Seconds: v.Timestamp.Unix(),
			},
		}
	}

	return &pb.GetMessagesResponse{
		Messages: messages,
	}, nil
}
