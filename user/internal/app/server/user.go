package server

import (
	"context"
	"github.com/Nixonxp/discord/user/internal/app/usecases"
	pb "github.com/Nixonxp/discord/user/pkg/api/v1"
	"log"
	"math/rand/v2"
)

// isServiceOk в зависимости от входящего значения вернет false, например
// передано 5, тогда (100 / 5 = 20) 20% вероятностью вернется false, для теста сервиса
func isServiceOk(probability int) bool {
	randNumber := rand.IntN(probability-1) + 1

	if randNumber == 1 {
		return false
	}

	return true
}

func (s *UserServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserDataResponse, error) {
	log.Printf("update user: received: %s", req.Login)

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	result, err := s.UserUsecase.UpdateUser(ctx, usecases.UpdateUserRequest{
		Id:    req.Id,
		Login: req.Login,
		Name:  req.Name,
		Email: req.Email,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UserDataResponse{
		Id:             result.UserID,
		Login:          result.Login,
		Name:           result.Name,
		Email:          result.Email,
		AvatarPhotoUrl: "url",
	}, nil
}

func (s *UserServer) GetUserByLogin(ctx context.Context, req *pb.GetUserByLoginRequest) (*pb.UserDataResponse, error) {
	log.Printf("get user by login: received: %s", req.Login)

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	result, err := s.UserUsecase.GetUserByLogin(ctx, usecases.GetUserByLoginRequest{
		Login: req.Login,
	})
	if err != nil {
		return nil, err
	}

	return &pb.UserDataResponse{
		Id:             result.UserID,
		Login:          result.Login,
		Name:           result.Name,
		Email:          result.Email,
		AvatarPhotoUrl: "url",
	}, nil
}

func (s *UserServer) GetUserFriends(ctx context.Context, req *pb.GetUserFriendsRequest) (*pb.GetUserFriendsResponse, error) {
	log.Printf("get user friends received")

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	result, err := s.UserUsecase.GetUserFriends(ctx, usecases.GetUserFriendsRequest{})
	if err != nil {
		return nil, err
	}

	return &pb.GetUserFriendsResponse{
		Friends: []*pb.Friend{
			{
				UserId: result.Friends[0].UserID,
				Login:  result.Friends[0].Login,
				Name:   result.Friends[0].Name,
				Email:  result.Friends[0].Email,
			},
		},
	}, nil
}

func (s *UserServer) AddToFriendByUserId(ctx context.Context, req *pb.AddToFriendByUserIdRequest) (*pb.ActionResponse, error) {
	log.Printf("add user to friends received: %d", req.UserId)

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
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
	log.Printf("accept user to friends received: %d", req.GetInviteId())

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
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
	log.Printf("accepdecline user to friends received: %d", req.GetInviteId())

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
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
