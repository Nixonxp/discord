package server

import (
	"context"
	"github.com/Nixonxp/discord/channel/internal/app/usecases"
	pb "github.com/Nixonxp/discord/channel/pkg/api/v1"
	"log"
)

func (s *ChannelServer) AddChannel(ctx context.Context, req *pb.AddChannelRequest) (*pb.ActionResponse, error) {
	log.Printf("add channel: received: %s", req.GetName())

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
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
	log.Printf("Delete channel: received: %d", req.GetChannelId())

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
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
	log.Printf("Join channel: received: %d", req.GetChannelId())

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
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
	log.Printf("Leave channel: received: %d", req.GetChannelId())

	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
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
