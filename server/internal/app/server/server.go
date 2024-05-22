package server

import (
	"context"
	"fmt"
	"github.com/Nixonxp/discord/server/internal/app/usecases"
	pb "github.com/Nixonxp/discord/server/pkg/api/v1"
	grpcutils "github.com/Nixonxp/discord/server/pkg/grpc_utils"
	"github.com/bufbuild/protovalidate-go"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"log"
	"net"
)

// Config - server config
type Config struct {
	ChainUnaryInterceptors []grpc.UnaryServerInterceptor
	UnaryInterceptors      []grpc.UnaryServerInterceptor
}

// Deps - server deps
type Deps struct {
	ServerUsecase usecases.UsecaseInterface
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
	log.Printf("Create server: received: %s", req.GetName())

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
		Id:   created.Id,
		Name: created.Name,
	}, nil
}

func (s *ServerServer) SearchServer(ctx context.Context, req *pb.SearchServerRequest) (*pb.SearchServerResponse, error) {
	log.Printf("Search server: received: %d", req.GetId())

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ServerUsecase.SearchServer(ctx, usecases.SearchServerRequest{
		Id:   req.GetId(),
		Name: req.GetName(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.SearchServerResponse{
		Id:   result.Id,
		Name: result.Name,
	}, nil
}

func (s *ServerServer) SubscribeServer(ctx context.Context, req *pb.SubscribeServerRequest) (*pb.ActionResponse, error) {
	log.Printf("Subscribe server: received: %d", req.GetServerId())

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
	log.Printf("Unsubscribe server: received: %d", req.GetServerId())

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
	log.Printf("Search server by user id: received: %d", req.GetUserId())

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
		Id:   result.Id,
		Name: result.Name,
	}, nil
}

func (s *ServerServer) InviteUserToServer(ctx context.Context, req *pb.InviteUserToServerRequest) (*pb.ActionResponse, error) {
	log.Printf("Invite user to server: received: %d", req.GetServerId())

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
	log.Printf("Publish message to server: received: %s", req.GetText())

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
	log.Printf("Get messages from server: received: %d", req.GetServerId())

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ServerUsecase.GetMessagesFromServer(ctx, usecases.GetMessagesFromServerRequest{
		ServerId: req.ServerId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetMessagesResponse{
		Messages: []*pb.Message{
			{
				Id:   result.Messages[0].Id,
				Text: result.Messages[0].Text,
				Timestamp: &timestamp.Timestamp{
					Seconds: result.Messages[0].Timestamp.Unix(),
				},
			},
		},
	}, nil
}
