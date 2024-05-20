package server

import (
	"context"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	pb "github.com/Nixonxp/discord/auth/pkg/api/v1"
	grpcutils "github.com/Nixonxp/discord/auth/pkg/grpc_utils"
	"log"
	"math/rand/v2"
)

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	loginInfo := req.Body.Login
	log.Printf("Register: received: %s", loginInfo)

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	user, err := s.AuthUsecase.Register(ctx, usecases.RegisterUserInfo{
		Login:    req.GetBody().Login,
		Name:     req.GetBody().Name,
		Email:    req.GetBody().Email,
		Password: req.GetBody().Password,
	})
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		UserId: user.UserID.String(),
		Login:  user.Login,
		Name:   user.Name,
		Email:  user.Email,
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

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginInfo := req.GetLogin()
	log.Printf("Login: received: %s", loginInfo)

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	userInfo, err := s.AuthUsecase.Login(ctx, usecases.LoginUserInfo{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		UserId: userInfo.UserID.String(),
		Login:  userInfo.Login,
		Name:   userInfo.Name,
		Token:  "{token}",
	}, nil
}

func (s *AuthServer) OauthLogin(ctx context.Context, req *pb.OauthLoginRequest) (*pb.OauthLoginResponse, error) {
	log.Printf("Oauth Login: received:")
	loginResult, err := s.AuthUsecase.OauthLogin(ctx, usecases.OauthLoginRequest{})
	if err != nil {
		return nil, err
	}

	return &pb.OauthLoginResponse{
		Code: loginResult.Code,
	}, nil
}

func (s *AuthServer) OauthLoginCallback(ctx context.Context, req *pb.OauthLoginCallbackRequest) (*pb.OauthLoginCallbackResponse, error) {
	log.Printf("Oauth login callback: received: %s", req.GetCode())
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	loginResult, err := s.AuthUsecase.OauthLoginCallback(ctx, usecases.OauthLoginCallbackRequest{
		Code: req.GetCode(),
	})

	if err != nil {
		return nil, err
	}

	return &pb.OauthLoginCallbackResponse{
		UserId: loginResult.UserID.String(),
		Login:  loginResult.Login,
		Name:   loginResult.Name,
		Token:  "{token}",
	}, nil
}
