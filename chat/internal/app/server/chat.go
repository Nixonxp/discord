package server

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/models"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	pb "github.com/Nixonxp/discord/chat/pkg/api/v1"
	grpcutils "github.com/Nixonxp/discord/chat/pkg/grpc_utils"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
)

func (s *ChatServer) SendUserPrivateMessage(ctx context.Context, req *pb.SendUserPrivateMessageRequest) (*pb.ActionResponse, error) {
	log.Printf("Send private message: received: %s", req.GetUserId())
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	userID := models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) // todo from auth

	result, err := s.ChatUsecase.SendUserPrivateMessage(ctx, usecases.SendUserPrivateMessageRequest{
		UserId:      req.GetUserId(),
		Text:        req.GetText(),
		CurrentUser: userID.String(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: result.Success,
	}, nil
}

func (s *ChatServer) GetUserPrivateMessages(ctx context.Context, req *pb.GetUserPrivateMessagesRequest) (*pb.GetMessagesResponse, error) {
	log.Printf("Send private message: received: %d", req.GetUserId())

	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	userID := models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) // todo from auth

	result, err := s.ChatUsecase.GetUserPrivateMessages(ctx, usecases.GetUserPrivateMessagesRequest{
		UserId:      req.GetUserId(),
		CurrentUser: userID.String(),
	})
	if err != nil {
		return nil, err
	}

	messages := make([]*pb.Message, len(result.Data))
	for k, v := range result.Data {
		messages[k] = &pb.Message{
			Id:   v.Id.String(),
			Text: v.Text,
			Timestamp: &timestamppb.Timestamp{
				Seconds: v.Timestamp.Unix(),
			},
			OwnerId: v.OwnerId.String(),
			ChatId:  v.ChatId.String(),
		}
	}

	return &pb.GetMessagesResponse{
		Messages: messages,
	}, nil
}

func (s *ChatServer) CreatePrivateChat(ctx context.Context, req *pb.CreatePrivateChatRequest) (*pb.CreatePrivateChatResponse, error) {
	log.Printf("create chat: received: %s", req.GetUserId())
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	userID := models.UserID(uuid.MustParse("284fef68-7e3e-4d1d-96a0-8c96f7b3b795")) // todo from auth

	result, err := s.ChatUsecase.CreatePrivateChat(ctx, usecases.CreatePrivateChatRequest{
		UserId:      req.GetUserId(),
		CurrentUser: userID.String(),
	})
	if err != nil {
		return nil, err
	}

	return &pb.CreatePrivateChatResponse{
		Success: true,
		ChatId:  result.Id.String(),
	}, nil
}

func (s *ChatServer) SendServerMessage(ctx context.Context, req *pb.SendServerMessageRequest) (*pb.ActionResponse, error) {
	log.Printf("send server message: received: %s", req.ServerId)
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	_, err := s.ChatUsecase.SendServerMessage(ctx, usecases.SendServerMessageRequest{
		ServerId: req.ServerId,
		Text:     req.Text,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ActionResponse{
		Success: true,
	}, nil
}

func (s *ChatServer) GetServerMessages(ctx context.Context, req *pb.GetServerMessagesRequest) (*pb.GetMessagesResponse, error) {
	log.Printf("get server messages received: %s", req.GetServerId())
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ChatUsecase.GetServerMessagesRequest(ctx, usecases.GetServerMessageRequest{
		ServerId: req.ServerId,
	})
	if err != nil {
		return nil, err
	}

	messages := make([]*pb.Message, len(result.Data))
	for k, v := range result.Data {
		messages[k] = &pb.Message{
			Id:   v.Id.String(),
			Text: v.Text,
			Timestamp: &timestamppb.Timestamp{
				Seconds: v.Timestamp.Unix(),
			},
			OwnerId: v.OwnerId.String(),
			ChatId:  v.ChatId.String(),
		}
	}

	return &pb.GetMessagesResponse{
		Messages: messages,
	}, nil
}
