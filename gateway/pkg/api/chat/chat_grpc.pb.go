// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: internal/app/api/chat/chat.proto

package chat

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	ChatService_CreatePrivateChat_FullMethodName      = "/github.com.Nixonxp.discord.chat.api.v1.ChatService/CreatePrivateChat"
	ChatService_SendUserPrivateMessage_FullMethodName = "/github.com.Nixonxp.discord.chat.api.v1.ChatService/SendUserPrivateMessage"
	ChatService_GetUserPrivateMessages_FullMethodName = "/github.com.Nixonxp.discord.chat.api.v1.ChatService/GetUserPrivateMessages"
)

// ChatServiceClient is the client API for ChatService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatServiceClient interface {
	CreatePrivateChat(ctx context.Context, in *CreatePrivateChatRequest, opts ...grpc.CallOption) (*CreatePrivateChatResponse, error)
	SendUserPrivateMessage(ctx context.Context, in *SendUserPrivateMessageRequest, opts ...grpc.CallOption) (*ActionResponse, error)
	GetUserPrivateMessages(ctx context.Context, in *GetUserPrivateMessagesRequest, opts ...grpc.CallOption) (*GetMessagesResponse, error)
}

type chatServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChatServiceClient(cc grpc.ClientConnInterface) ChatServiceClient {
	return &chatServiceClient{cc}
}

func (c *chatServiceClient) CreatePrivateChat(ctx context.Context, in *CreatePrivateChatRequest, opts ...grpc.CallOption) (*CreatePrivateChatResponse, error) {
	out := new(CreatePrivateChatResponse)
	err := c.cc.Invoke(ctx, ChatService_CreatePrivateChat_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) SendUserPrivateMessage(ctx context.Context, in *SendUserPrivateMessageRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, ChatService_SendUserPrivateMessage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *chatServiceClient) GetUserPrivateMessages(ctx context.Context, in *GetUserPrivateMessagesRequest, opts ...grpc.CallOption) (*GetMessagesResponse, error) {
	out := new(GetMessagesResponse)
	err := c.cc.Invoke(ctx, ChatService_GetUserPrivateMessages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServiceServer is the server API for ChatService service.
// All implementations must embed UnimplementedChatServiceServer
// for forward compatibility
type ChatServiceServer interface {
	CreatePrivateChat(context.Context, *CreatePrivateChatRequest) (*CreatePrivateChatResponse, error)
	SendUserPrivateMessage(context.Context, *SendUserPrivateMessageRequest) (*ActionResponse, error)
	GetUserPrivateMessages(context.Context, *GetUserPrivateMessagesRequest) (*GetMessagesResponse, error)
	mustEmbedUnimplementedChatServiceServer()
}

// UnimplementedChatServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChatServiceServer struct {
}

func (UnimplementedChatServiceServer) CreatePrivateChat(context.Context, *CreatePrivateChatRequest) (*CreatePrivateChatResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePrivateChat not implemented")
}
func (UnimplementedChatServiceServer) SendUserPrivateMessage(context.Context, *SendUserPrivateMessageRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendUserPrivateMessage not implemented")
}
func (UnimplementedChatServiceServer) GetUserPrivateMessages(context.Context, *GetUserPrivateMessagesRequest) (*GetMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserPrivateMessages not implemented")
}
func (UnimplementedChatServiceServer) mustEmbedUnimplementedChatServiceServer() {}

// UnsafeChatServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServiceServer will
// result in compilation errors.
type UnsafeChatServiceServer interface {
	mustEmbedUnimplementedChatServiceServer()
}

func RegisterChatServiceServer(s grpc.ServiceRegistrar, srv ChatServiceServer) {
	s.RegisterService(&ChatService_ServiceDesc, srv)
}

func _ChatService_CreatePrivateChat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePrivateChatRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).CreatePrivateChat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_CreatePrivateChat_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).CreatePrivateChat(ctx, req.(*CreatePrivateChatRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_SendUserPrivateMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendUserPrivateMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).SendUserPrivateMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_SendUserPrivateMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).SendUserPrivateMessage(ctx, req.(*SendUserPrivateMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChatService_GetUserPrivateMessages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserPrivateMessagesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServiceServer).GetUserPrivateMessages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChatService_GetUserPrivateMessages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServiceServer).GetUserPrivateMessages(ctx, req.(*GetUserPrivateMessagesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChatService_ServiceDesc is the grpc.ServiceDesc for ChatService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChatService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "github.com.Nixonxp.discord.chat.api.v1.ChatService",
	HandlerType: (*ChatServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePrivateChat",
			Handler:    _ChatService_CreatePrivateChat_Handler,
		},
		{
			MethodName: "SendUserPrivateMessage",
			Handler:    _ChatService_SendUserPrivateMessage_Handler,
		},
		{
			MethodName: "GetUserPrivateMessages",
			Handler:    _ChatService_GetUserPrivateMessages_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/app/api/chat/chat.proto",
}
