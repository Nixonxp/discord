package server

import (
	"context"
	"github.com/Nixonxp/discord/auth/internal/app/usecases"
	pb "github.com/Nixonxp/discord/auth/pkg/api/v1"
	grpcutils "github.com/Nixonxp/discord/auth/pkg/grpc_utils"
)

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	loginInfo := req.Body.Login
	s.Log.WithContext(ctx).WithField("loginInfo", loginInfo).Info("register: received")

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

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	loginInfo := req.GetLogin()
	s.Log.WithContext(ctx).WithField("loginInfo", loginInfo).Info("login received")

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
	s.Log.WithContext(ctx).Info("Oauth Login: received")
	loginResult, err := s.AuthUsecase.OauthLogin(ctx, usecases.OauthLoginRequest{})
	if err != nil {
		return nil, err
	}

	return &pb.OauthLoginResponse{
		Code: loginResult.Code,
	}, nil
}

func (s *AuthServer) OauthLoginCallback(ctx context.Context, req *pb.OauthLoginCallbackRequest) (*pb.OauthLoginCallbackResponse, error) {
	s.Log.WithContext(ctx).WithField("code", req.GetCode()).Info("Oauth login callback: received")
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
