// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: api/v1/server.proto

package server

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
	ServerService_CreateServer_FullMethodName           = "/github.com.Nixonxp.discord.server.api.v1.ServerService/CreateServer"
	ServerService_SearchServer_FullMethodName           = "/github.com.Nixonxp.discord.server.api.v1.ServerService/SearchServer"
	ServerService_SubscribeServer_FullMethodName        = "/github.com.Nixonxp.discord.server.api.v1.ServerService/SubscribeServer"
	ServerService_UnsubscribeServer_FullMethodName      = "/github.com.Nixonxp.discord.server.api.v1.ServerService/UnsubscribeServer"
	ServerService_SearchServerByUserId_FullMethodName   = "/github.com.Nixonxp.discord.server.api.v1.ServerService/SearchServerByUserId"
	ServerService_InviteUserToServer_FullMethodName     = "/github.com.Nixonxp.discord.server.api.v1.ServerService/InviteUserToServer"
	ServerService_PublishMessageOnServer_FullMethodName = "/github.com.Nixonxp.discord.server.api.v1.ServerService/PublishMessageOnServer"
	ServerService_GetMessagesFromServer_FullMethodName  = "/github.com.Nixonxp.discord.server.api.v1.ServerService/GetMessagesFromServer"
)

// ServerServiceClient is the client API for ServerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServerServiceClient interface {
	CreateServer(ctx context.Context, in *CreateServerRequest, opts ...grpc.CallOption) (*CreateServerResponse, error)
	SearchServer(ctx context.Context, in *SearchServerRequest, opts ...grpc.CallOption) (*SearchServerResponse, error)
	SubscribeServer(ctx context.Context, in *SubscribeServerRequest, opts ...grpc.CallOption) (*ActionResponse, error)
	UnsubscribeServer(ctx context.Context, in *UnsubscribeServerRequest, opts ...grpc.CallOption) (*ActionResponse, error)
	SearchServerByUserId(ctx context.Context, in *SearchServerByUserIdRequest, opts ...grpc.CallOption) (*SearchServerByUserIdResponse, error)
	InviteUserToServer(ctx context.Context, in *InviteUserToServerRequest, opts ...grpc.CallOption) (*ActionResponse, error)
	PublishMessageOnServer(ctx context.Context, in *PublishMessageOnServerRequest, opts ...grpc.CallOption) (*ActionResponse, error)
	GetMessagesFromServer(ctx context.Context, in *GetMessagesFromServerRequest, opts ...grpc.CallOption) (*GetMessagesResponse, error)
}

type serverServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewServerServiceClient(cc grpc.ClientConnInterface) ServerServiceClient {
	return &serverServiceClient{cc}
}

func (c *serverServiceClient) CreateServer(ctx context.Context, in *CreateServerRequest, opts ...grpc.CallOption) (*CreateServerResponse, error) {
	out := new(CreateServerResponse)
	err := c.cc.Invoke(ctx, ServerService_CreateServer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverServiceClient) SearchServer(ctx context.Context, in *SearchServerRequest, opts ...grpc.CallOption) (*SearchServerResponse, error) {
	out := new(SearchServerResponse)
	err := c.cc.Invoke(ctx, ServerService_SearchServer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverServiceClient) SubscribeServer(ctx context.Context, in *SubscribeServerRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, ServerService_SubscribeServer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverServiceClient) UnsubscribeServer(ctx context.Context, in *UnsubscribeServerRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, ServerService_UnsubscribeServer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverServiceClient) SearchServerByUserId(ctx context.Context, in *SearchServerByUserIdRequest, opts ...grpc.CallOption) (*SearchServerByUserIdResponse, error) {
	out := new(SearchServerByUserIdResponse)
	err := c.cc.Invoke(ctx, ServerService_SearchServerByUserId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverServiceClient) InviteUserToServer(ctx context.Context, in *InviteUserToServerRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, ServerService_InviteUserToServer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverServiceClient) PublishMessageOnServer(ctx context.Context, in *PublishMessageOnServerRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, ServerService_PublishMessageOnServer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serverServiceClient) GetMessagesFromServer(ctx context.Context, in *GetMessagesFromServerRequest, opts ...grpc.CallOption) (*GetMessagesResponse, error) {
	out := new(GetMessagesResponse)
	err := c.cc.Invoke(ctx, ServerService_GetMessagesFromServer_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServerServiceServer is the server API for ServerService service.
// All implementations must embed UnimplementedServerServiceServer
// for forward compatibility
type ServerServiceServer interface {
	CreateServer(context.Context, *CreateServerRequest) (*CreateServerResponse, error)
	SearchServer(context.Context, *SearchServerRequest) (*SearchServerResponse, error)
	SubscribeServer(context.Context, *SubscribeServerRequest) (*ActionResponse, error)
	UnsubscribeServer(context.Context, *UnsubscribeServerRequest) (*ActionResponse, error)
	SearchServerByUserId(context.Context, *SearchServerByUserIdRequest) (*SearchServerByUserIdResponse, error)
	InviteUserToServer(context.Context, *InviteUserToServerRequest) (*ActionResponse, error)
	PublishMessageOnServer(context.Context, *PublishMessageOnServerRequest) (*ActionResponse, error)
	GetMessagesFromServer(context.Context, *GetMessagesFromServerRequest) (*GetMessagesResponse, error)
	mustEmbedUnimplementedServerServiceServer()
}

// UnimplementedServerServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServerServiceServer struct {
}

func (UnimplementedServerServiceServer) CreateServer(context.Context, *CreateServerRequest) (*CreateServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateServer not implemented")
}
func (UnimplementedServerServiceServer) SearchServer(context.Context, *SearchServerRequest) (*SearchServerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchServer not implemented")
}
func (UnimplementedServerServiceServer) SubscribeServer(context.Context, *SubscribeServerRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeServer not implemented")
}
func (UnimplementedServerServiceServer) UnsubscribeServer(context.Context, *UnsubscribeServerRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnsubscribeServer not implemented")
}
func (UnimplementedServerServiceServer) SearchServerByUserId(context.Context, *SearchServerByUserIdRequest) (*SearchServerByUserIdResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SearchServerByUserId not implemented")
}
func (UnimplementedServerServiceServer) InviteUserToServer(context.Context, *InviteUserToServerRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method InviteUserToServer not implemented")
}
func (UnimplementedServerServiceServer) PublishMessageOnServer(context.Context, *PublishMessageOnServerRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PublishMessageOnServer not implemented")
}
func (UnimplementedServerServiceServer) GetMessagesFromServer(context.Context, *GetMessagesFromServerRequest) (*GetMessagesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMessagesFromServer not implemented")
}
func (UnimplementedServerServiceServer) mustEmbedUnimplementedServerServiceServer() {}

// UnsafeServerServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServerServiceServer will
// result in compilation errors.
type UnsafeServerServiceServer interface {
	mustEmbedUnimplementedServerServiceServer()
}

func RegisterServerServiceServer(s grpc.ServiceRegistrar, srv ServerServiceServer) {
	s.RegisterService(&ServerService_ServiceDesc, srv)
}

func _ServerService_CreateServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServiceServer).CreateServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServerService_CreateServer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServiceServer).CreateServer(ctx, req.(*CreateServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerService_SearchServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServiceServer).SearchServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServerService_SearchServer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServiceServer).SearchServer(ctx, req.(*SearchServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerService_SubscribeServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscribeServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServiceServer).SubscribeServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServerService_SubscribeServer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServiceServer).SubscribeServer(ctx, req.(*SubscribeServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerService_UnsubscribeServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnsubscribeServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServiceServer).UnsubscribeServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServerService_UnsubscribeServer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServiceServer).UnsubscribeServer(ctx, req.(*UnsubscribeServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerService_SearchServerByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SearchServerByUserIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServiceServer).SearchServerByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServerService_SearchServerByUserId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServiceServer).SearchServerByUserId(ctx, req.(*SearchServerByUserIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerService_InviteUserToServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InviteUserToServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServiceServer).InviteUserToServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServerService_InviteUserToServer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServiceServer).InviteUserToServer(ctx, req.(*InviteUserToServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerService_PublishMessageOnServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PublishMessageOnServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServiceServer).PublishMessageOnServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServerService_PublishMessageOnServer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServiceServer).PublishMessageOnServer(ctx, req.(*PublishMessageOnServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ServerService_GetMessagesFromServer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMessagesFromServerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServerServiceServer).GetMessagesFromServer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ServerService_GetMessagesFromServer_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServerServiceServer).GetMessagesFromServer(ctx, req.(*GetMessagesFromServerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ServerService_ServiceDesc is the grpc.ServiceDesc for ServerService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ServerService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "github.com.Nixonxp.discord.server.api.v1.ServerService",
	HandlerType: (*ServerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateServer",
			Handler:    _ServerService_CreateServer_Handler,
		},
		{
			MethodName: "SearchServer",
			Handler:    _ServerService_SearchServer_Handler,
		},
		{
			MethodName: "SubscribeServer",
			Handler:    _ServerService_SubscribeServer_Handler,
		},
		{
			MethodName: "UnsubscribeServer",
			Handler:    _ServerService_UnsubscribeServer_Handler,
		},
		{
			MethodName: "SearchServerByUserId",
			Handler:    _ServerService_SearchServerByUserId_Handler,
		},
		{
			MethodName: "InviteUserToServer",
			Handler:    _ServerService_InviteUserToServer_Handler,
		},
		{
			MethodName: "PublishMessageOnServer",
			Handler:    _ServerService_PublishMessageOnServer_Handler,
		},
		{
			MethodName: "GetMessagesFromServer",
			Handler:    _ServerService_GetMessagesFromServer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/v1/server.proto",
}
