package services

import (
	"context"
	pb_auth "github.com/Nixonxp/discord/gateway/pkg/api/auth"
	pb_channel "github.com/Nixonxp/discord/gateway/pkg/api/channel"
	pb_chat "github.com/Nixonxp/discord/gateway/pkg/api/chat"
	pb_server "github.com/Nixonxp/discord/gateway/pkg/api/server"
	pb_user "github.com/Nixonxp/discord/gateway/pkg/api/user"
	pb "github.com/Nixonxp/discord/gateway/pkg/api/v1"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"log"
)

type Deps struct {
	AuthConn    *grpc.ClientConn
	UserConn    *grpc.ClientConn
	ServerConn  *grpc.ClientConn
	ChannelConn *grpc.ClientConn
	ChatConn    *grpc.ClientConn
}

type DiscordGatewayService struct {
	Deps
}

func NewDiscordGatewayService(d Deps) *DiscordGatewayService {
	return &DiscordGatewayService{
		d,
	}
}

func (s *DiscordGatewayService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	authClient := pb_auth.NewAuthServiceClient(s.AuthConn)
	registerReq := pb_auth.RegisterRequest{
		Body: &pb_auth.User{
			Login:    req.GetLogin(),
			Name:     req.GetName(),
			Email:    req.GetEmail(),
			Password: req.GetPassword(),
		},
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "auth_service.Register")
	defer span.Finish()

	registerResponse, err := authClient.Register(ctx, &registerReq)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		UserId: registerResponse.GetUserId(),
		Login:  registerResponse.GetLogin(),
		Name:   registerResponse.GetName(),
		Email:  registerResponse.GetEmail(),
	}, nil
}

func (s *DiscordGatewayService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	authClient := pb_auth.NewAuthServiceClient(s.AuthConn)
	loginReq := pb_auth.LoginRequest{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "auth_service.Login")
	defer span.Finish()

	loginResponse, err := authClient.Login(ctx, &loginReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.LoginResponse{
		UserId: loginResponse.GetUserId(),
		Login:  loginResponse.GetLogin(),
		Name:   loginResponse.GetName(),
		Token:  loginResponse.GetToken(),
	}, nil
}

func (s *DiscordGatewayService) OauthLogin(ctx context.Context, _ *pb.OauthLoginRequest) (*pb.OauthLoginResponse, error) {
	authClient := pb_auth.NewAuthServiceClient(s.AuthConn)
	loginReq := pb_auth.OauthLoginRequest{}

	span, ctx := opentracing.StartSpanFromContext(ctx, "auth_service.OauthLogin")
	defer span.Finish()

	loginResponse, err := authClient.OauthLogin(ctx, &loginReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.OauthLoginResponse{
		Code: loginResponse.GetCode(),
	}, nil
}

func (s *DiscordGatewayService) OauthLoginCallback(ctx context.Context, req *pb.OauthLoginCallbackRequest) (*pb.OauthLoginCallbackResponse, error) {
	authClient := pb_auth.NewAuthServiceClient(s.AuthConn)
	loginReq := pb_auth.OauthLoginCallbackRequest{
		Code: req.GetCode(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "auth_service.OauthLoginCallback")
	defer span.Finish()

	loginResponse, err := authClient.OauthLoginCallback(ctx, &loginReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.OauthLoginCallbackResponse{
		UserId: loginResponse.GetUserId(),
		Login:  loginResponse.GetLogin(),
		Name:   loginResponse.GetName(),
		Token:  loginResponse.GetToken(),
	}, nil
}

func (s *DiscordGatewayService) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserDataResponse, error) {
	userClient := pb_user.NewUserServiceClient(s.UserConn)
	request := pb_user.UpdateUserRequest{
		Id:             req.GetId(),
		Login:          req.GetBody().GetLogin(),
		Name:           req.GetBody().GetName(),
		Email:          req.GetBody().GetEmail(),
		AvatarPhotoUrl: req.GetBody().GetAvatarPhotoUrl(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "user_service.UpdateUser")
	defer span.Finish()

	response, err := userClient.UpdateUser(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.UserDataResponse{
		Id:             response.GetId(),
		Login:          response.GetLogin(),
		Name:           response.GetName(),
		Email:          response.GetEmail(),
		AvatarPhotoUrl: response.GetAvatarPhotoUrl(),
	}, nil
}

func (s *DiscordGatewayService) GetUserByLogin(ctx context.Context, req *pb.GetUserByLoginRequest) (*pb.UserDataResponse, error) {
	userClient := pb_user.NewUserServiceClient(s.UserConn)
	request := pb_user.GetUserByLoginRequest{
		Login: req.GetLogin(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "user_service.GetUserByLogin")
	defer span.Finish()

	response, err := userClient.GetUserByLogin(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.UserDataResponse{
		Id:             response.GetId(),
		Login:          response.GetLogin(),
		Name:           response.GetName(),
		Email:          response.GetEmail(),
		AvatarPhotoUrl: response.GetAvatarPhotoUrl(),
	}, nil
}

func (s *DiscordGatewayService) GetUserFriends(ctx context.Context, req *pb.GetUserFriendsRequest) (*pb.GetUserFriendsResponse, error) {
	userClient := pb_user.NewUserServiceClient(s.UserConn)
	request := pb_user.GetUserFriendsRequest{}

	response, err := userClient.GetUserFriends(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "user_service.GetUserFriends")
	defer span.Finish()

	res := make([]*pb.Friend, len(response.GetFriends()))
	for i, v := range response.GetFriends() {
		res[i] = &pb.Friend{
			UserId: v.GetUserId(),
			Login:  v.GetLogin(),
			Name:   v.GetName(),
			Email:  v.GetEmail(),
		}
	}

	return &pb.GetUserFriendsResponse{
		Friends: res,
	}, nil
}

func (s *DiscordGatewayService) AddToFriendByUserId(ctx context.Context, req *pb.AddToFriendByUserIdRequest) (*pb.ActionResponse, error) {
	userClient := pb_user.NewUserServiceClient(s.UserConn)
	request := pb_user.AddToFriendByUserIdRequest{
		UserId: req.GetUserId(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "user_service.AddToFriendByUserId")
	defer span.Finish()

	response, err := userClient.AddToFriendByUserId(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *DiscordGatewayService) AcceptFriendInvite(ctx context.Context, req *pb.AcceptFriendInviteRequest) (*pb.ActionResponse, error) {
	userClient := pb_user.NewUserServiceClient(s.UserConn)
	request := pb_user.AcceptFriendInviteRequest{
		InviteId: req.GetInviteId(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "user_service.AcceptFriendInvite")
	defer span.Finish()

	response, err := userClient.AcceptFriendInvite(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *DiscordGatewayService) DeclineFriendInvite(ctx context.Context, req *pb.DeclineFriendInviteRequest) (*pb.ActionResponse, error) {
	userClient := pb_user.NewUserServiceClient(s.UserConn)
	request := pb_user.DeclineFriendInviteRequest{
		InviteId: req.GetInviteId(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "user_service.DeclineFriendInvite")
	defer span.Finish()

	response, err := userClient.DeclineFriendInvite(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *DiscordGatewayService) CreateServer(ctx context.Context, req *pb.CreateServerRequest) (*pb.CreateServerResponse, error) {
	serverClient := pb_server.NewServerServiceClient(s.ServerConn)
	request := pb_server.CreateServerRequest{
		Name: req.GetName(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "server_service.CreateServer")
	defer span.Finish()

	response, err := serverClient.CreateServer(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.CreateServerResponse{
		Id:   response.GetId(),
		Name: req.GetName(),
	}, nil
}

func (s *DiscordGatewayService) SearchServer(ctx context.Context, req *pb.SearchServerRequest) (*pb.SearchServerResponse, error) {
	serverClient := pb_server.NewServerServiceClient(s.ServerConn)
	request := pb_server.SearchServerRequest{
		Name: req.GetName(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "server_service.SearchServer")
	defer span.Finish()

	response, err := serverClient.SearchServer(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.SearchServerResponse{
		Id:   response.GetId(),
		Name: response.GetName(),
	}, nil
}

func (s *DiscordGatewayService) SubscribeServer(ctx context.Context, req *pb.SubscribeServerRequest) (*pb.ActionResponse, error) {
	serverClient := pb_server.NewServerServiceClient(s.ServerConn)
	request := pb_server.SubscribeServerRequest{
		ServerId: req.GetServerId(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "server_service.SubscribeServer")
	defer span.Finish()

	response, err := serverClient.SubscribeServer(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *DiscordGatewayService) UnsubscribeServer(ctx context.Context, req *pb.UnsubscribeServerRequest) (*pb.ActionResponse, error) {
	serverClient := pb_server.NewServerServiceClient(s.ServerConn)
	request := pb_server.UnsubscribeServerRequest{
		ServerId: req.GetServerId(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "server_service.UnsubscribeServer")
	defer span.Finish()

	response, err := serverClient.UnsubscribeServer(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *DiscordGatewayService) SearchServerByUserId(ctx context.Context, req *pb.SearchServerByUserIdRequest) (*pb.SearchServerByUserIdResponse, error) {
	serverClient := pb_server.NewServerServiceClient(s.ServerConn)
	request := pb_server.SearchServerByUserIdRequest{
		UserId: req.GetUserId(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "server_service.SearchServerByUserId")
	defer span.Finish()

	response, err := serverClient.SearchServerByUserId(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.SearchServerByUserIdResponse{
		Id:   response.GetId(),
		Name: response.GetName(),
	}, nil
}

func (s *DiscordGatewayService) InviteUserToServer(ctx context.Context, req *pb.InviteUserToServerRequest) (*pb.ActionResponse, error) {
	serverClient := pb_server.NewServerServiceClient(s.ServerConn)
	request := pb_server.InviteUserToServerRequest{
		UserId:   req.GetUserId(),
		ServerId: req.GetServerId(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "server_service.InviteUserToServer")
	defer span.Finish()

	response, err := serverClient.InviteUserToServer(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *DiscordGatewayService) PublishMessageOnServer(ctx context.Context, req *pb.PublishMessageOnServerRequest) (*pb.ActionResponse, error) {
	serverClient := pb_server.NewServerServiceClient(s.ServerConn)
	request := pb_server.PublishMessageOnServerRequest{
		ServerId: req.GetServerId(),
		Text:     req.GetText(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "server_service.PublishMessageOnServer")
	defer span.Finish()

	response, err := serverClient.PublishMessageOnServer(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *DiscordGatewayService) GetMessagesFromServer(ctx context.Context, req *pb.GetMessagesFromServerRequest) (*pb.GetMessagesResponse, error) {
	serverClient := pb_server.NewServerServiceClient(s.ServerConn)
	request := pb_server.GetMessagesFromServerRequest{
		ServerId: req.GetServerId(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "server_service.GetMessagesFromServer")
	defer span.Finish()

	response, err := serverClient.GetMessagesFromServer(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	messages := make([]*pb.Message, len(response.GetMessages()))
	for i, m := range response.GetMessages() {
		messages[i] = &pb.Message{
			Id:        m.GetId(),
			Text:      m.GetText(),
			Timestamp: m.GetTimestamp(),
		}
	}

	return &pb.GetMessagesResponse{
		Messages: messages,
	}, nil
}

func (s *DiscordGatewayService) AddChannel(ctx context.Context, req *pb.AddChannelRequest) (*pb.ActionResponse, error) {
	channelClient := pb_channel.NewChannelServiceClient(s.ChannelConn)
	request := pb_channel.AddChannelRequest{
		Name: req.GetName(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "channel_service.AddChannel")
	defer span.Finish()

	response, err := channelClient.AddChannel(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *DiscordGatewayService) DeleteChannel(ctx context.Context, req *pb.DeleteChannelRequest) (*pb.ActionResponse, error) {
	channelClient := pb_channel.NewChannelServiceClient(s.ChannelConn)
	request := pb_channel.DeleteChannelRequest{
		ChannelId: req.GetChannelId(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "channel_service.DeleteChannel")
	defer span.Finish()

	response, err := channelClient.DeleteChannel(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *DiscordGatewayService) JoinChannel(ctx context.Context, req *pb.JoinChannelRequest) (*pb.ActionResponse, error) {
	channelClient := pb_channel.NewChannelServiceClient(s.ChannelConn)
	request := pb_channel.JoinChannelRequest{
		ChannelId: req.GetChannelId(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "channel_service.JoinChannel")
	defer span.Finish()

	response, err := channelClient.JoinChannel(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *DiscordGatewayService) LeaveChannel(ctx context.Context, req *pb.LeaveChannelRequest) (*pb.ActionResponse, error) {
	channelClient := pb_channel.NewChannelServiceClient(s.ChannelConn)
	request := pb_channel.LeaveChannelRequest{
		ChannelId: req.GetChannelId(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "channel_service.LeaveChannel")
	defer span.Finish()

	response, err := channelClient.LeaveChannel(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *DiscordGatewayService) SendUserPrivateMessage(ctx context.Context, req *pb.SendUserPrivateMessageRequest) (*pb.ActionResponse, error) {
	chatClient := pb_chat.NewChatServiceClient(s.ChatConn)
	request := pb_chat.SendUserPrivateMessageRequest{
		UserId: req.GetUserId(),
		Text:   req.GetText(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "chat_service.SendUserPrivateMessage")
	defer span.Finish()

	response, err := chatClient.SendUserPrivateMessage(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *DiscordGatewayService) GetUserPrivateMessages(ctx context.Context, req *pb.GetUserPrivateMessagesRequest) (*pb.GetMessagesResponse, error) {
	chatClient := pb_chat.NewChatServiceClient(s.ChatConn)
	request := pb_chat.GetUserPrivateMessagesRequest{
		UserId: req.GetUserId(),
	}

	span, ctx := opentracing.StartSpanFromContext(ctx, "chat_service.GetUserPrivateMessages")
	defer span.Finish()

	response, err := chatClient.GetUserPrivateMessages(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	messages := make([]*pb.Message, len(response.GetMessages()))
	for i, m := range response.GetMessages() {
		messages[i] = &pb.Message{
			Id:        m.GetId(),
			Text:      m.GetText(),
			Timestamp: m.GetTimestamp(),
			OwnerId:   m.OwnerId,
			ChatId:    m.ChatId,
		}
	}

	return &pb.GetMessagesResponse{
		Messages: messages,
	}, nil
}
