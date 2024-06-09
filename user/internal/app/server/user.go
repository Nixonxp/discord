package server

import (
	"context"
	"github.com/Nixonxp/discord/user/internal/app/models"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	pb "github.com/Nixonxp/discord/user/pkg/api/v1"
	auth "github.com/Nixonxp/discord/user/pkg/auth"
	grpcutils "github.com/Nixonxp/discord/user/pkg/grpc_utils"
)

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.UserDataResponse, error) {
	s.ctxLog(ctx).Infof("create user: received: %s", req.Login)

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.UserUsecase.CreateUser(ctx, usecases.CreateUserRequest{
		Login:    req.Login,
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UserDataResponse{
		Id:             result.Id.String(),
		Login:          result.Login,
		Name:           result.Name,
		Email:          result.Email,
		AvatarPhotoUrl: "",
	}, nil
}

func (s *UserServer) GetUserForLogin(ctx context.Context, req *pb.GetUserForLoginRequest) (*pb.GetUserForLoginResponse, error) {
	s.ctxLog(ctx).Infof("get user for login: received: %s", req.Login)

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.UserUsecase.GetUserForLogin(ctx, usecases.GetUserByLoginRequest{
		Login: req.Login,
	})
	if err != nil {
		return nil, err
	}

	return &pb.GetUserForLoginResponse{
		Id:       result.Id.String(),
		Login:    result.Login,
		Name:     result.Name,
		Email:    result.Email,
		Password: result.Password,
	}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserDataResponse, error) {
	s.ctxLog(ctx).Infof("update user: received: %s", req.Login)

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	userId, err := auth.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, models.Unauthenticated
	}

	result, err := s.UserUsecase.UpdateUser(ctx, usecases.UpdateUserRequest{
		Id:             req.Id,
		Login:          req.Login,
		Name:           req.Name,
		Email:          req.Email,
		AvatarPhotoUrl: req.AvatarPhotoUrl,
		CurrentUserId:  userId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UserDataResponse{
		Id:             result.Id.String(),
		Login:          result.Login,
		Name:           result.Name,
		Email:          result.Email,
		AvatarPhotoUrl: req.AvatarPhotoUrl,
	}, nil
}

func (s *UserServer) GetUserByLogin(ctx context.Context, req *pb.GetUserByLoginRequest) (*pb.UserDataResponse, error) {
	s.ctxLog(ctx).Infof("get user by login: received: %s", req.Login)

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	_, err := auth.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, models.Unauthenticated
	}

	result, err := s.UserUsecase.GetUserForLogin(ctx, usecases.GetUserByLoginRequest{
		Login: req.Login,
	})
	if err != nil {
		s.Log.WithContext(ctx).Error(err)
		return nil, err
	}

	return &pb.UserDataResponse{
		Id:             result.Id.String(),
		Login:          result.Login,
		Name:           result.Name,
		Email:          result.Email,
		AvatarPhotoUrl: result.AvatarPhotoUrl,
	}, nil
}
