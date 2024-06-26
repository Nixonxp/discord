package server

import (
	"context"
	"fmt"
	"github.com/Nixonxp/discord/gateway/internal/app/services/gateway"
	pb "github.com/Nixonxp/discord/gateway/pkg/api/v1"
	grpcutils "github.com/Nixonxp/discord/gateway/pkg/grpc_utils"
	"github.com/bufbuild/protovalidate-go"
	"log"
)

type Deps struct {
	DiscordGatewayService *services.DiscordGatewayService
}

type DiscordGatewayServiceServer struct {
	pb.UnimplementedGatewayServiceServer
	Deps

	validator *protovalidate.Validator
}

func NewDiscordGatewayServiceServer(ctx context.Context, d Deps) (*DiscordGatewayServiceServer, error) {
	srv := &DiscordGatewayServiceServer{
		Deps: d,
	}

	// validator
	{
		validator, err := protovalidate.New(
			protovalidate.WithDisableLazy(true),
			protovalidate.WithMessages(
				&pb.RegisterRequest{},
				&pb.LoginRequest{},
				&pb.OauthLoginRequest{},
				&pb.OauthLoginCallbackRequest{},
				&pb.UpdateUserRequest{},
				&pb.GetUserByLoginRequest{},
				&pb.GetUserFriendsRequest{},
				&pb.AddToFriendByUserIdRequest{},
				&pb.AcceptFriendInviteRequest{},
				&pb.DeclineFriendInviteRequest{},
				&pb.CreateServerRequest{},
				&pb.SearchServerRequest{},
				&pb.SubscribeServerRequest{},
				&pb.UnsubscribeServerRequest{},
				&pb.SearchServerByUserIdRequest{},
				&pb.InviteUserToServerRequest{},
				&pb.PublishMessageOnServerRequest{},
				&pb.GetMessagesFromServerRequest{},
				&pb.AddChannelRequest{},
				&pb.DeleteChannelRequest{},
				&pb.JoinChannelRequest{},
				&pb.LeaveChannelRequest{},
				&pb.SendUserPrivateMessageRequest{},
				&pb.GetUserPrivateMessagesRequest{},
				&pb.DeleteFromFriendRequest{},
				&pb.GetUserInvitesRequest{},
				&pb.RefreshRequest{},
				&pb.CreatePrivateChatRequest{},
			),
		)
		if err != nil {
			return nil, fmt.Errorf("server: failed to initialize validator: %w", err)
		}
		srv.validator = validator
	}

	return srv, nil
}

func (s *DiscordGatewayServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.Register(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.Login(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.Refresh(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) OauthLogin(ctx context.Context, req *pb.OauthLoginRequest) (*pb.OauthLoginResponse, error) {
	resp, err := s.DiscordGatewayService.OauthLogin(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) OauthLoginCallback(ctx context.Context, req *pb.OauthLoginCallbackRequest) (*pb.OauthLoginCallbackResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.OauthLoginCallback(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserDataResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) GetUserByLogin(ctx context.Context, req *pb.GetUserByLoginRequest) (*pb.UserDataResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.GetUserByLogin(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) GetUserFriends(ctx context.Context, req *pb.GetUserFriendsRequest) (*pb.GetUserFriendsResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.GetUserFriends(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) AddToFriendByUserId(ctx context.Context, req *pb.AddToFriendByUserIdRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.AddToFriendByUserId(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) AcceptFriendInvite(ctx context.Context, req *pb.AcceptFriendInviteRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.AcceptFriendInvite(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) DeclineFriendInvite(ctx context.Context, req *pb.DeclineFriendInviteRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.DeclineFriendInvite(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) GetUserInvites(ctx context.Context, req *pb.GetUserInvitesRequest) (*pb.GetUserInvitesResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.GetUserInvites(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) DeleteFromFriend(ctx context.Context, req *pb.DeleteFromFriendRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.DeleteFromFriend(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) CreateServer(ctx context.Context, req *pb.CreateServerRequest) (*pb.CreateServerResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.CreateServer(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) SearchServer(ctx context.Context, req *pb.SearchServerRequest) (*pb.SearchServerResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.SearchServer(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) SubscribeServer(ctx context.Context, req *pb.SubscribeServerRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}
	resp, err := s.DiscordGatewayService.SubscribeServer(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) UnsubscribeServer(ctx context.Context, req *pb.UnsubscribeServerRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.UnsubscribeServer(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) SearchServerByUserId(ctx context.Context, req *pb.SearchServerByUserIdRequest) (*pb.SearchServerByUserIdResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.SearchServerByUserId(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) InviteUserToServer(ctx context.Context, req *pb.InviteUserToServerRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.InviteUserToServer(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) PublishMessageOnServer(ctx context.Context, req *pb.PublishMessageOnServerRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.PublishMessageOnServer(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) GetMessagesFromServer(ctx context.Context, req *pb.GetMessagesFromServerRequest) (*pb.GetMessagesResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.GetMessagesFromServer(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) AddChannel(ctx context.Context, req *pb.AddChannelRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.AddChannel(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) DeleteChannel(ctx context.Context, req *pb.DeleteChannelRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.DeleteChannel(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) JoinChannel(ctx context.Context, req *pb.JoinChannelRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.JoinChannel(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) LeaveChannel(ctx context.Context, req *pb.LeaveChannelRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.LeaveChannel(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) SendUserPrivateMessage(ctx context.Context, req *pb.SendUserPrivateMessageRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.SendUserPrivateMessage(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) CreatePrivateChat(ctx context.Context, req *pb.CreatePrivateChatRequest) (*pb.CreatePrivateChatResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.CreatePrivateChat(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *DiscordGatewayServiceServer) GetUserPrivateMessages(ctx context.Context, req *pb.GetUserPrivateMessagesRequest) (*pb.GetMessagesResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, grpcutils.RPCValidationError(err)
	}

	resp, err := s.DiscordGatewayService.GetUserPrivateMessages(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
