package server

import (
	"context"
	"github.com/Nixonxp/discord/user/internal/app/models"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	pb "github.com/Nixonxp/discord/user/pkg/api/v1"
	auth "github.com/Nixonxp/discord/user/pkg/auth"
	grpcutils "github.com/Nixonxp/discord/user/pkg/grpc_utils"
	log "github.com/Nixonxp/discord/user/pkg/logger"
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

	result, err := s.UserUsecase.GetUserByLogin(ctx, usecases.GetUserByLoginRequest{
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

func (s *UserServer) GetUserFriends(ctx context.Context, req *pb.GetUserFriendsRequest) (*pb.GetUserFriendsResponse, error) {
	s.ctxLog(ctx).Infof("get user friends received")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	userId, err := auth.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, models.Unauthenticated
	}

	result, err := s.UserUsecase.GetUserFriends(ctx, usecases.GetUserFriendsRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	friends := make([]*pb.Friend, len(result.Friends))
	for i, friend := range result.Friends {
		friends[i] = &pb.Friend{
			UserId: friend.UserId.String(),
			Login:  friend.Login,
			Name:   friend.Name,
			Email:  friend.Email,
		}
	}

	return &pb.GetUserFriendsResponse{
		Friends: friends,
	}, nil
}

func (s *UserServer) GetUserInvites(ctx context.Context, req *pb.GetUserInvitesRequest) (*pb.GetUserInvitesResponse, error) {
	s.ctxLog(ctx).Infof("get user friends received")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	userId, err := auth.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, models.Unauthenticated
	}

	result, err := s.UserUsecase.GetUserInvites(ctx, usecases.GetUserInvitesRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	pbInvites := make([]*pb.FriendInvite, len(result.Invites))
	for i, invite := range result.Invites {
		pbInvites[i] = &pb.FriendInvite{
			InviteId: invite.InviteId.String(),
			OwnerId:  invite.OwnerId.String(),
			UserId:   invite.UserId.String(),
			Status:   invite.Status,
		}
	}

	return &pb.GetUserInvitesResponse{
		Invites: pbInvites,
	}, nil
}

func (s *UserServer) AddToFriendByUserId(ctx context.Context, req *pb.AddToFriendByUserIdRequest) (*pb.ActionResponse, error) {
	s.ctxLog(ctx).Infof("add user to friends received: %d", req.UserId)
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	userId, err := auth.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, models.Unauthenticated
	}

	result, err := s.UserUsecase.AddToFriendByUserId(ctx, usecases.AddToFriendByUserIdRequest{
		UserId:  req.UserId,
		OwnerId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *UserServer) CreateOrGetUser(ctx context.Context, req *pb.CreateOrGetUserRequest) (*pb.UserDataResponse, error) {
	s.ctxLog(ctx).Infof("create or get user: received: %s", req.Login)

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.UserUsecase.CreateOrGetUser(ctx, usecases.CreateOrGetUserRequest{
		Login:          req.Login,
		Name:           req.Name,
		Email:          req.Email,
		Password:       req.Password,
		OauthId:        req.OauthId,
		AvatarPhotoUrl: req.AvatarPhotoUrl,
	})
	if err != nil {
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

func (s *UserServer) AcceptFriendInvite(ctx context.Context, req *pb.AcceptFriendInviteRequest) (*pb.ActionResponse, error) {
	s.ctxLog(ctx).Infof("accept user to friends received: %d", req.GetInviteId())

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	userId, err := auth.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, models.Unauthenticated
	}

	result, err := s.UserUsecase.AcceptFriendInvite(ctx, usecases.AcceptFriendInviteRequest{
		InviteId:      req.GetInviteId(),
		CurrentUserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *UserServer) DeclineFriendInvite(ctx context.Context, req *pb.DeclineFriendInviteRequest) (*pb.ActionResponse, error) {
	s.ctxLog(ctx).Infof("decline user to friends received: %s", req.GetInviteId())

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	userId, err := auth.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, models.Unauthenticated
	}

	result, err := s.UserUsecase.DeclineFriendInvite(ctx, usecases.DeclineFriendInviteRequest{
		InviteId:      req.GetInviteId(),
		CurrentUserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *UserServer) DeleteFromFriend(ctx context.Context, req *pb.DeleteFromFriendRequest) (*pb.ActionResponse, error) {
	s.ctxLog(ctx).Infof("delete friends received: %s", req.GetFriendId())

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	userId, err := auth.GetUserIdFromContext(ctx)
	if err != nil {
		return nil, models.Unauthenticated
	}

	result, err := s.UserUsecase.DeleteFromFriend(ctx, usecases.DeleteFromFriendRequest{
		FriendId:      req.GetFriendId(),
		CurrentUserId: userId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *UserServer) ctxLog(ctx context.Context) *log.Logger {
	return s.Log.WithContext(ctx)
}
