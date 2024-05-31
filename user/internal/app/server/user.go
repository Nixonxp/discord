package server

import (
	"context"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	pb "github.com/Nixonxp/discord/user/pkg/api/v1"
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

func (s *UserServer) GetUserByLoginAndPassword(ctx context.Context, req *pb.GetUserByLoginAndPasswordRequest) (*pb.UserDataResponse, error) {
	s.ctxLog(ctx).Infof("get user by login: received: %s", req.Login)

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.UserUsecase.GetUserByLoginAndPassword(ctx, usecases.GetUserByLoginAndPasswordRequest{
		Login:    req.Login,
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
		AvatarPhotoUrl: "url", // todo add
	}, nil
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserDataResponse, error) {
	s.ctxLog(ctx).Infof("update user: received: %s", req.Login)

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.UserUsecase.UpdateUser(ctx, usecases.UpdateUserRequest{
		Id:             req.Id,
		Login:          req.Login,
		Name:           req.Name,
		Email:          req.Email,
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
		AvatarPhotoUrl: req.AvatarPhotoUrl,
	}, nil
}

func (s *UserServer) GetUserByLogin(ctx context.Context, req *pb.GetUserByLoginRequest) (*pb.UserDataResponse, error) {
	s.ctxLog(ctx).Infof("get user by login: received: %s", req.Login)

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
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

	result, err := s.UserUsecase.GetUserFriends(ctx, usecases.GetUserFriendsRequest{})
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

	result, err := s.UserUsecase.GetUserInvites(ctx, usecases.GetUserInvitesRequest{
		UserId: "4aee4258-1cdd-11ef-b2b5-4612de44ab9f", // todo from auth data,
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

	result, err := s.UserUsecase.AddToFriendByUserId(ctx, usecases.AddToFriendByUserIdRequest{
		UserId: req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *UserServer) AcceptFriendInvite(ctx context.Context, req *pb.AcceptFriendInviteRequest) (*pb.ActionResponse, error) {
	s.ctxLog(ctx).Infof("accept user to friends received: %d", req.GetInviteId())

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.UserUsecase.AcceptFriendInvite(ctx, usecases.AcceptFriendInviteRequest{
		InviteId: req.GetInviteId(),
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

	result, err := s.UserUsecase.DeclineFriendInvite(ctx, usecases.DeclineFriendInviteRequest{
		InviteId: req.GetInviteId(),
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

	result, err := s.UserUsecase.DeleteFromFriend(ctx, usecases.DeleteFromFriendRequest{
		FriendId: req.GetFriendId(),
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
