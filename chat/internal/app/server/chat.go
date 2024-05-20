package server

import (
	"context"
	"github.com/Nixonxp/discord/chat/internal/app/usecases"
	pb "github.com/Nixonxp/discord/chat/pkg/api/v1"
	grpcutils "github.com/Nixonxp/discord/chat/pkg/grpc_utils"
	"google.golang.org/protobuf/types/known/timestamppb"
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

func (s *ChatServer) SendUserPrivateMessage(ctx context.Context, req *pb.SendUserPrivateMessageRequest) (*pb.ActionResponse, error) {
	log.Printf("Send private message: received: %s", req.GetUserId())
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	result, err := s.ChatUsecase.SendUserPrivateMessage(ctx, usecases.SendUserPrivateMessageRequest{
		UserId: req.GetUserId(),
		Text:   req.GetText(),
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

	result, err := s.ChatUsecase.GetUserPrivateMessages(ctx, usecases.GetUserPrivateMessagesRequest{
		UserId: req.GetUserId(),
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
