package main

import (
	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"context"
	"errors"
	"fmt"
	pb "github.com/Nixonxp/discord/auth/pkg/api/v1"
	"github.com/bufbuild/protovalidate-go"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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
)

type server struct {
	pb.UnimplementedAuthServiceServer
	validator *protovalidate.Validator
}

func NewServer() (*server, error) {
	srv := &server{}

	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(
			// Добавляем сюда все запросы наши
			&pb.RegisterRequest{},
			&pb.LoginRequest{},
			&pb.OauthLoginRequest{},
			&pb.OauthLoginCallbackRequest{},
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	srv.validator = validator
	return srv, nil
}

func (s *server) Register(_ context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	loginInfo := req.Body.Login
	log.Printf("Register: received: %s", loginInfo)

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	return &pb.RegisterResponse{
		Id:    1,
		Login: req.Body.Login,
		Name:  req.Body.Name,
		Email: req.Body.Email,
	}, nil
}

func (s *server) Login(_ context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginInfo := req.GetLogin()
	log.Printf("Login: received: %s", loginInfo)

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	return &pb.LoginResponse{
		UserId: 1,
		Login:  loginInfo,
		Name:   "user name",
		Token:  "{token}",
	}, nil
}

func (s *server) OauthLogin(_ context.Context, req *pb.OauthLoginRequest) (*pb.OauthLoginResponse, error) {
	log.Printf("Oauth Login: received:")

	return &pb.OauthLoginResponse{
		Code: "{code}",
	}, nil
}

func (s *server) OauthLoginCallback(_ context.Context, req *pb.OauthLoginCallbackRequest) (*pb.OauthLoginCallbackResponse, error) {
	log.Printf("Oauth login callback: received: %s", req.GetCode())

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	return &pb.OauthLoginCallbackResponse{
		UserId: 1,
		Login:  "login user",
		Name:   "user name",
		Token:  "{token}",
	}, nil
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
		pb.RegisterAuthServiceServer(grpcServer, server)

		reflection.Register(grpcServer)

		lis, err := net.Listen("tcp", ":8080")
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
			httpPort = "8081"
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
