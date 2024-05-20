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
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
)

// Config - server config
type Config struct {
	GRPCPort string
	HTTPPort string

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

func NewServerServer(ctx context.Context, cfg Config, d Deps) (*ServerServer, error) {
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

	// grpc
	{
		grpcServerOptions := unaryInterceptorsToGrpcServerOptions(cfg.UnaryInterceptors...)
		grpcServerOptions = append(grpcServerOptions,
			grpc.ChainUnaryInterceptor(cfg.ChainUnaryInterceptors...),
		)

		grpcServer := grpc.NewServer(grpcServerOptions...)
		pb.RegisterServerServiceServer(grpcServer, srv)

		reflection.Register(grpcServer)

		lis, err := net.Listen("tcp", cfg.GRPCPort)
		if err != nil {
			return nil, fmt.Errorf("server: failed to listen: %v", err)
		}

		srv.grpc.server = grpcServer
		srv.grpc.lis = lis
	}

	{
		e := echo.New()

		e.Use(middleware.Logger())
		e.Use(middleware.Recover())

		e.GET("/health", func(c echo.Context) error {
			status := http.StatusOK
			statusMessage := "OK"

			if !isServiceOk(10) {
				status = http.StatusInternalServerError
				statusMessage = "Error"
			}

			return c.JSON(status, struct{ Status string }{Status: statusMessage})
		})

		e.GET("/ready", func(c echo.Context) error {
			status := http.StatusOK
			statusMessage := "OK"

			if !isServiceOk(5) {
				status = http.StatusInternalServerError
				statusMessage = "Error"
			}

			return c.JSON(status, struct{ Status string }{Status: statusMessage})
		})

		srv.http.lis = e
		srv.http.port = cfg.HTTPPort
	}

	return srv, nil
}

// Run - serve
func (s *ServerServer) Run(ctx context.Context) error {
	group := errgroup.Group{}

	group.Go(func() error {
		log.Println("start serve", s.grpc.lis.Addr())
		if err := s.grpc.server.Serve(s.grpc.lis); err != nil {
			return fmt.Errorf("server: serve grpc: %v", err)
		}
		return nil
	})

	group.Go(func() error {
		log.Println("start http server", s.http.port)
		err := s.http.lis.Start(s.http.port)
		if err != nil {
			return fmt.Errorf("server: serve http: %v", err)
		}

		return nil
	})

	return group.Wait()
}

func unaryInterceptorsToGrpcServerOptions(interceptors ...grpc.UnaryServerInterceptor) []grpc.ServerOption {
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

// isServiceOk в зависимости от входящего значения вернет false, например
// передано 5, тогда (100 / 5 = 20) 20% вероятностью вернется false, для теста сервиса
func isServiceOk(probability int) bool {
	randNumber := rand.IntN(probability-1) + 1

	if randNumber == 1 {
		return false
	}

	return true
}
