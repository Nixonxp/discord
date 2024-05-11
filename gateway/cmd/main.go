package main

import (
	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"context"
	"errors"
	"fmt"
	pb_auth "github.com/Nixonxp/discord/gateway/pkg/api/auth"
	pb_channel "github.com/Nixonxp/discord/gateway/pkg/api/channel"
	pb_chat "github.com/Nixonxp/discord/gateway/pkg/api/chat"
	pb_server "github.com/Nixonxp/discord/gateway/pkg/api/server"
	pb_user "github.com/Nixonxp/discord/gateway/pkg/api/user"
	pb "github.com/Nixonxp/discord/gateway/pkg/api/v1"
	"github.com/bufbuild/protovalidate-go"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

type server struct {
	pb.UnimplementedGatewayServiceServer
	validator *protovalidate.Validator

	authConn    *grpc.ClientConn
	userConn    *grpc.ClientConn
	serverConn  *grpc.ClientConn
	channelConn *grpc.ClientConn
	chatConn    *grpc.ClientConn
}

func NewServer(authConn *grpc.ClientConn, userConn *grpc.ClientConn, serverConn *grpc.ClientConn, channelConn *grpc.ClientConn, chatConn *grpc.ClientConn) (*server, error) {
	srv := &server{
		authConn:    authConn,
		userConn:    userConn,
		serverConn:  serverConn,
		channelConn: channelConn,
		chatConn:    chatConn,
	}

	validator, err := protovalidate.New(
		protovalidate.WithDisableLazy(true),
		protovalidate.WithMessages(
			// Добавляем сюда все запросы наши
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
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	srv.validator = validator
	return srv, nil
}

func convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr *protovalidate.ValidationError) *errdetails.BadRequest {
	return &errdetails.BadRequest{
		FieldViolations: protovalidateVialationsToGoogleViolations(valErr.Violations),
	}
}

func protovalidateVialationsToGoogleViolations(vs []*validate.Violation) []*errdetails.BadRequest_FieldViolation {
	res := make([]*errdetails.BadRequest_FieldViolation, len(vs))
	for i, v := range vs {
		res[i] = &errdetails.BadRequest_FieldViolation{
			Field:       v.FieldPath,
			Description: v.Message,
		}
	}
	return res
}

func (s *server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	authClient := pb_auth.NewAuthServiceClient(s.authConn)
	registerReq := pb_auth.RegisterRequest{
		Body: &pb_auth.User{
			Login:    req.GetLogin(),
			Name:     req.GetName(),
			Email:    req.GetEmail(),
			Password: req.GetPassword(),
		},
	}

	registerResponse, err := authClient.Register(ctx, &registerReq)
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{
		Id:    registerResponse.GetId(),
		Login: registerResponse.GetLogin(),
		Name:  registerResponse.GetName(),
		Email: registerResponse.GetEmail(),
	}, nil
}

func (s *server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, rpcValidationError(err)
	}

	authClient := pb_auth.NewAuthServiceClient(s.authConn)
	loginReq := pb_auth.LoginRequest{
		Login:    req.GetLogin(),
		Password: req.GetPassword(),
	}

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

func (s *server) OauthLogin(ctx context.Context, _ *pb.OauthLoginRequest) (*pb.OauthLoginResponse, error) {
	authClient := pb_auth.NewAuthServiceClient(s.authConn)
	loginReq := pb_auth.OauthLoginRequest{}

	loginResponse, err := authClient.OauthLogin(ctx, &loginReq)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.OauthLoginResponse{
		Code: loginResponse.GetCode(),
	}, nil
}

func rpcValidationError(err error) error {
	if err == nil {
		return nil
	}

	var valErr *protovalidate.ValidationError
	if ok := errors.As(err, &valErr); ok {
		st, err := status.New(codes.InvalidArgument, codes.InvalidArgument.String()).
			WithDetails(convertProtovalidateValidationErrorToErrdetailsBadRequest(valErr))
		if err == nil {
			return st.Err()
		}
	}

	return status.Error(codes.Internal, err.Error())
}

func (s *server) OauthLoginCallback(ctx context.Context, req *pb.OauthLoginCallbackRequest) (*pb.OauthLoginCallbackResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, rpcValidationError(err)
	}

	authClient := pb_auth.NewAuthServiceClient(s.authConn)
	loginReq := pb_auth.OauthLoginCallbackRequest{
		Code: "code ...",
	}

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

func (s *server) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UserDataResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, rpcValidationError(err)
	}

	userClient := pb_user.NewUserServiceClient(s.userConn)
	request := pb_user.UpdateUserRequest{
		Id:             req.GetId(),
		Login:          req.GetBody().GetLogin(),
		Name:           req.GetBody().GetName(),
		Email:          req.GetBody().GetEmail(),
		Password:       req.GetBody().GetPassword(),
		AvatarPhotoUrl: req.GetBody().GetAvatarPhotoUrl(),
	}

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

func (s *server) GetUserByLogin(ctx context.Context, req *pb.GetUserByLoginRequest) (*pb.UserDataResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	userClient := pb_user.NewUserServiceClient(s.userConn)
	request := pb_user.GetUserByLoginRequest{
		Login: req.GetLogin(),
	}

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

func (s *server) GetUserFriends(ctx context.Context, req *pb.GetUserFriendsRequest) (*pb.GetUserFriendsResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	userClient := pb_user.NewUserServiceClient(s.userConn)
	request := pb_user.GetUserFriendsRequest{}

	response, err := userClient.GetUserFriends(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

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

func (s *server) AddToFriendByUserId(ctx context.Context, req *pb.AddToFriendByUserIdRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	userClient := pb_user.NewUserServiceClient(s.userConn)
	request := pb_user.AddToFriendByUserIdRequest{
		UserId: req.GetUserId(),
	}

	response, err := userClient.AddToFriendByUserId(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *server) AcceptFriendInvite(ctx context.Context, req *pb.AcceptFriendInviteRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	userClient := pb_user.NewUserServiceClient(s.userConn)
	request := pb_user.AcceptFriendInviteRequest{
		InviteId: req.GetInviteId(),
	}

	response, err := userClient.AcceptFriendInvite(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *server) DeclineFriendInvite(ctx context.Context, req *pb.DeclineFriendInviteRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	userClient := pb_user.NewUserServiceClient(s.userConn)
	request := pb_user.DeclineFriendInviteRequest{
		InviteId: req.GetInviteId(),
	}

	response, err := userClient.DeclineFriendInvite(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *server) CreateServer(ctx context.Context, req *pb.CreateServerRequest) (*pb.CreateServerResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, rpcValidationError(err)
	}

	serverClient := pb_server.NewServerServiceClient(s.serverConn)
	request := pb_server.CreateServerRequest{
		Name: req.GetName(),
	}

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

func (s *server) SearchServer(ctx context.Context, req *pb.SearchServerRequest) (*pb.SearchServerResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	serverClient := pb_server.NewServerServiceClient(s.serverConn)
	request := pb_server.SearchServerRequest{
		Name: req.GetName(),
	}

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

func (s *server) SubscribeServer(ctx context.Context, req *pb.SubscribeServerRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	serverClient := pb_server.NewServerServiceClient(s.serverConn)
	request := pb_server.SubscribeServerRequest{
		ServerId: req.GetServerId(),
	}

	response, err := serverClient.SubscribeServer(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *server) UnsubscribeServer(ctx context.Context, req *pb.UnsubscribeServerRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	serverClient := pb_server.NewServerServiceClient(s.serverConn)
	request := pb_server.UnsubscribeServerRequest{
		ServerId: req.GetServerId(),
	}

	response, err := serverClient.UnsubscribeServer(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *server) SearchServerByUserId(ctx context.Context, req *pb.SearchServerByUserIdRequest) (*pb.SearchServerByUserIdResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	serverClient := pb_server.NewServerServiceClient(s.serverConn)
	request := pb_server.SearchServerByUserIdRequest{
		UserId: req.GetUserId(),
	}

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

func (s *server) InviteUserToServer(ctx context.Context, req *pb.InviteUserToServerRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	serverClient := pb_server.NewServerServiceClient(s.serverConn)
	request := pb_server.InviteUserToServerRequest{
		UserId:   req.GetUserId(),
		ServerId: req.GetServerId(),
	}

	response, err := serverClient.InviteUserToServer(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *server) PublishMessageOnServer(ctx context.Context, req *pb.PublishMessageOnServerRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	serverClient := pb_server.NewServerServiceClient(s.serverConn)
	request := pb_server.PublishMessageOnServerRequest{
		ServerId: req.GetServerId(),
		Text:     req.GetText(),
	}

	response, err := serverClient.PublishMessageOnServer(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *server) GetMessagesFromServer(ctx context.Context, req *pb.GetMessagesFromServerRequest) (*pb.GetMessagesResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, rpcValidationError(err)
	}

	serverClient := pb_server.NewServerServiceClient(s.serverConn)
	request := pb_server.GetMessagesFromServerRequest{
		ServerId: req.GetServerId(),
	}

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

func (s *server) AddChannel(ctx context.Context, req *pb.AddChannelRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	channelClient := pb_channel.NewChannelServiceClient(s.channelConn)
	request := pb_channel.AddChannelRequest{
		Name: "channel name",
	}

	response, err := channelClient.AddChannel(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *server) DeleteChannel(ctx context.Context, req *pb.DeleteChannelRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	channelClient := pb_channel.NewChannelServiceClient(s.channelConn)
	request := pb_channel.DeleteChannelRequest{
		ChannelId: req.GetChannelId(),
	}

	response, err := channelClient.DeleteChannel(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *server) JoinChannel(ctx context.Context, req *pb.JoinChannelRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	channelClient := pb_channel.NewChannelServiceClient(s.channelConn)
	request := pb_channel.JoinChannelRequest{
		ChannelId: req.GetChannelId(),
	}

	response, err := channelClient.JoinChannel(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *server) LeaveChannel(ctx context.Context, req *pb.LeaveChannelRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	channelClient := pb_channel.NewChannelServiceClient(s.channelConn)
	request := pb_channel.LeaveChannelRequest{
		ChannelId: req.GetChannelId(),
	}

	response, err := channelClient.LeaveChannel(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *server) SendUserPrivateMessage(ctx context.Context, req *pb.SendUserPrivateMessageRequest) (*pb.ActionResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, rpcValidationError(err)
	}

	chatClient := pb_chat.NewChatServiceClient(s.chatConn)
	request := pb_chat.SendUserPrivateMessageRequest{
		UserId: req.GetUserId(),
		Text:   req.GetText(),
	}

	response, err := chatClient.SendUserPrivateMessage(ctx, &request)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.ActionResponse{
		Success: response.GetSuccess(),
	}, nil
}

func (s *server) GetUserPrivateMessages(ctx context.Context, req *pb.GetUserPrivateMessagesRequest) (*pb.GetMessagesResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		log.Println(err)
		return nil, rpcValidationError(err)
	}

	chatClient := pb_chat.NewChatServiceClient(s.chatConn)
	request := pb_chat.GetUserPrivateMessagesRequest{
		UserId: req.GetUserId(),
	}

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
		}
	}

	return &pb.GetMessagesResponse{
		Messages: messages,
	}, nil
}

// allowCORS allows Cross Origin Resource Sharing from any origin.
// Don't do this without consideration in production systems.
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept", "Authorization"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	grpclog.Infof("Preflight request for %s", r.URL.Path)
}

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		ctx := context.Background()
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		// grpc connection to auth service
		authConn, err := grpc.DialContext(context.Background(), "auth:8080", grpc.WithIdleTimeout(5*time.Second), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}

		// grpc connection to user service
		userConn, err := grpc.DialContext(context.Background(), "user:8080", grpc.WithIdleTimeout(5*time.Second), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}

		// grpc connection to server service
		serverConn, err := grpc.DialContext(context.Background(), "server:8080", grpc.WithIdleTimeout(5*time.Second), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}

		// grpc connection to channel service
		channelConn, err := grpc.DialContext(context.Background(), "channel:8080", grpc.WithIdleTimeout(5*time.Second), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}

		// grpc connection to chat service
		chatConn, err := grpc.DialContext(context.Background(), "chat:8080", grpc.WithIdleTimeout(5*time.Second), grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithInsecure())
		if err != nil {
			log.Fatal(err)
		}

		server, err := NewServer(authConn, userConn, serverConn, channelConn, chatConn)
		if err != nil {
			log.Fatalf("failed to create server: %v", err)
		}

		// Register gRPC server endpoint
		// Note: Make sure the gRPC server is running properly and accessible
		mux := runtime.NewServeMux()
		if err = pb.RegisterGatewayServiceHandlerServer(ctx, mux, server); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}

		httpServer := &http.Server{Handler: allowCORS(mux)}

		lis, err := net.Listen("tcp", ":8080")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		// Start HTTP server (and proxy calls to gRPC server endpoint)
		log.Printf("server listening at %v", lis.Addr())
		if err := httpServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		router := mux.NewRouter().StrictSlash(true)
		sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
		router.PathPrefix("/swaggerui/").Handler(sh)

		log.Fatal(http.ListenAndServe(":8081", router))
	}()

	wg.Wait()

}
