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

	userLoginInfo, err := s.AuthUsecase.Login(ctx, usecases.LoginUserInfo{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		Token:        userLoginInfo.Token,
		RefreshToken: userLoginInfo.RefreshToken,
	}, nil
}

func (s *AuthServer) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	token := req.GetRefreshToken()
	s.Log.WithContext(ctx).WithField("token", token).Info("refresh token received")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	newToken, err := s.AuthUsecase.Refresh(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &pb.RefreshResponse{
		Token: newToken,
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
		Code:  req.GetCode(),
		State: req.GetState(),
	})

	if err != nil {
		return nil, err
	}

	return &pb.OauthLoginCallbackResponse{
		Token:        loginResult.Token,
		RefreshToken: loginResult.RefreshToken,
	}, nil
}
