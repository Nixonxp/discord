package server

import (
	"context"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
	pb "github.com/Nixonxp/discord/channel/pkg/api/v1"
	grpcutils "github.com/Nixonxp/discord/channel/pkg/grpc_utils"
)

func (s *ChannelServer) AddChannel(ctx context.Context, req *pb.AddChannelRequest) (*pb.ActionResponse, error) {
	s.Log.WithContext(ctx).WithField("name", req.GetName()).Info("add channel: received")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ChannelUsecase.AddChannel(ctx, usecases.AddChannelRequest{
		Name: req.GetName(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *ChannelServer) DeleteChannel(ctx context.Context, req *pb.DeleteChannelRequest) (*pb.ActionResponse, error) {
	s.Log.WithContext(ctx).WithField("id", req.GetChannelId()).Info("Delete channel: received")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ChannelUsecase.DeleteChannel(ctx, usecases.DeleteChannelRequest{
		ChannelId: req.ChannelId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *ChannelServer) JoinChannel(ctx context.Context, req *pb.JoinChannelRequest) (*pb.ActionResponse, error) {
	s.Log.WithContext(ctx).WithField("id", req.GetChannelId()).Info("Join channel: received")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ChannelUsecase.JoinChannel(ctx, usecases.JoinChannelRequest{
		ChannelId: req.ChannelId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *ChannelServer) LeaveChannel(ctx context.Context, req *pb.LeaveChannelRequest) (*pb.ActionResponse, error) {
	s.Log.WithContext(ctx).WithField("id", req.GetChannelId()).Info("Leave channel: received")

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ChannelUsecase.LeaveChannel(ctx, usecases.LeaveChannelRequest{
		ChannelId: req.ChannelId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}
