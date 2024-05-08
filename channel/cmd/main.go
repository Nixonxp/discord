package main

import (
	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"context"
	"errors"
	"fmt"
	pb "github.com/Nixonxp/discord/channel/pkg/api/v1"
	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"math/rand/v2"
	"net"
	"net/http"
	"os"
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type server struct {
	pb.UnimplementedChannelServiceServer
	validator *protovalidate.Validator
}

func NewServer() (*server, error) {
	srv := &server{}

	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(
			&pb.AddChannelRequest{},
			&pb.DeleteChannelRequest{},
			&pb.JoinChannelRequest{},
			&pb.LeaveChannelRequest{},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	srv.validator = validator
	return srv, nil
}

func rpcValidationError(err error) error {
	if err == nil {
		return nil
	}

	var valErr *protovalidate.ValidationError
	if ok := errors.As(err, &valErr); ok {
		st, err := status.New(codes.InvalidArgument, codes.InvalidArgument.String()).
			WithDetails(convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr))
		if err == nil {
			return st.Err()
		}
	}

	return status.Error(codes.Internal, err.Error())
}

func convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr *protovalidate.ValidationError) *errdetails.BadRequest {
	return &errdetails.BadRequest{
		FieldViolations: protovalidateVialationsToGoogleViolations(valErr.Violations),
	}
}

func protovalidateVialationsToGoogleViolations(vs []*validate.Violation) []*errdetails.BadRequest_FieldViolation {
	res := make([]*errdetails.BadRequest_FieldViolation, len(vs))
	for i, v := range vs {
		res[i] = &errdetails.BadRequest_FieldViolation{
			Field:       v.FieldPath,
			Description: v.Message,
		}
	}
	return res
}

func (s *server) AddChannel(_ context.Context, req *pb.AddChannelRequest) (*pb.ActionResponse, error) {
	log.Printf("Add channel: received: %s", req.GetName())

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	return &pb.ActionResponse{
		Success: true,
	}, nil
}

func (s *server) DeleteChannel(_ context.Context, req *pb.DeleteChannelRequest) (*pb.ActionResponse, error) {
	log.Printf("Delete channel: received: %d", req.GetChannelId())

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	return &pb.ActionResponse{
		Success: true,
	}, nil
}

func (s *server) JoinChannel(_ context.Context, req *pb.JoinChannelRequest) (*pb.ActionResponse, error) {
	log.Printf("Join channel: received: %d", req.GetChannelId())

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	return &pb.ActionResponse{
		Success: true,
	}, nil
}

func (s *server) LeaveChannel(_ context.Context, req *pb.LeaveChannelRequest) (*pb.ActionResponse, error) {
	log.Printf("Leave channel: received: %d", req.GetChannelId())

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	return &pb.ActionResponse{
		Success: true,
	}, nil
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	server, err := NewServer()
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		grpcServer := grpc.NewServer()
		pb.RegisterChannelServiceServer(grpcServer, server)

		reflection.Register(grpcServer)

		lis, err := net.Listen("tcp", ":8804")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		log.Printf("server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

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

		httpPort := os.Getenv("PORT")
		if httpPort == "" {
			httpPort = "8084"
		}

		e.Logger.Fatal(e.Start(":" + httpPort))
	}()

	wg.Wait()
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
