// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: internal/app/api/channel/channel.proto

package channel

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
	ChannelService_AddChannel_FullMethodName    = "/github.com.Nixonxp.discord.channel.api.v1.ChannelService/AddChannel"
	ChannelService_DeleteChannel_FullMethodName = "/github.com.Nixonxp.discord.channel.api.v1.ChannelService/DeleteChannel"
	ChannelService_JoinChannel_FullMethodName   = "/github.com.Nixonxp.discord.channel.api.v1.ChannelService/JoinChannel"
	ChannelService_LeaveChannel_FullMethodName  = "/github.com.Nixonxp.discord.channel.api.v1.ChannelService/LeaveChannel"
)

// ChannelServiceClient is the client API for ChannelService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChannelServiceClient interface {
	AddChannel(ctx context.Context, in *AddChannelRequest, opts ...grpc.CallOption) (*ActionResponse, error)
	DeleteChannel(ctx context.Context, in *DeleteChannelRequest, opts ...grpc.CallOption) (*ActionResponse, error)
	JoinChannel(ctx context.Context, in *JoinChannelRequest, opts ...grpc.CallOption) (*ActionResponse, error)
	LeaveChannel(ctx context.Context, in *LeaveChannelRequest, opts ...grpc.CallOption) (*ActionResponse, error)
}

type channelServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewChannelServiceClient(cc grpc.ClientConnInterface) ChannelServiceClient {
	return &channelServiceClient{cc}
}

func (c *channelServiceClient) AddChannel(ctx context.Context, in *AddChannelRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, ChannelService_AddChannel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelServiceClient) DeleteChannel(ctx context.Context, in *DeleteChannelRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, ChannelService_DeleteChannel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelServiceClient) JoinChannel(ctx context.Context, in *JoinChannelRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, ChannelService_JoinChannel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *channelServiceClient) LeaveChannel(ctx context.Context, in *LeaveChannelRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, ChannelService_LeaveChannel_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChannelServiceServer is the server API for ChannelService service.
// All implementations must embed UnimplementedChannelServiceServer
// for forward compatibility
type ChannelServiceServer interface {
	AddChannel(context.Context, *AddChannelRequest) (*ActionResponse, error)
	DeleteChannel(context.Context, *DeleteChannelRequest) (*ActionResponse, error)
	JoinChannel(context.Context, *JoinChannelRequest) (*ActionResponse, error)
	LeaveChannel(context.Context, *LeaveChannelRequest) (*ActionResponse, error)
	mustEmbedUnimplementedChannelServiceServer()
}

// UnimplementedChannelServiceServer must be embedded to have forward compatible implementations.
type UnimplementedChannelServiceServer struct {
}

func (UnimplementedChannelServiceServer) AddChannel(context.Context, *AddChannelRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddChannel not implemented")
}
func (UnimplementedChannelServiceServer) DeleteChannel(context.Context, *DeleteChannelRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteChannel not implemented")
}
func (UnimplementedChannelServiceServer) JoinChannel(context.Context, *JoinChannelRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method JoinChannel not implemented")
}
func (UnimplementedChannelServiceServer) LeaveChannel(context.Context, *LeaveChannelRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LeaveChannel not implemented")
}
func (UnimplementedChannelServiceServer) mustEmbedUnimplementedChannelServiceServer() {}

// UnsafeChannelServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChannelServiceServer will
// result in compilation errors.
type UnsafeChannelServiceServer interface {
	mustEmbedUnimplementedChannelServiceServer()
}

func RegisterChannelServiceServer(s grpc.ServiceRegistrar, srv ChannelServiceServer) {
	s.RegisterService(&ChannelService_ServiceDesc, srv)
}

func _ChannelService_AddChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelServiceServer).AddChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChannelService_AddChannel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelServiceServer).AddChannel(ctx, req.(*AddChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelService_DeleteChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelServiceServer).DeleteChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChannelService_DeleteChannel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelServiceServer).DeleteChannel(ctx, req.(*DeleteChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelService_JoinChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelServiceServer).JoinChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChannelService_JoinChannel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelServiceServer).JoinChannel(ctx, req.(*JoinChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ChannelService_LeaveChannel_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveChannelRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChannelServiceServer).LeaveChannel(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ChannelService_LeaveChannel_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChannelServiceServer).LeaveChannel(ctx, req.(*LeaveChannelRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ChannelService_ServiceDesc is the grpc.ServiceDesc for ChannelService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ChannelService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "github.com.Nixonxp.discord.channel.api.v1.ChannelService",
	HandlerType: (*ChannelServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddChannel",
			Handler:    _ChannelService_AddChannel_Handler,
		},
		{
			MethodName: "DeleteChannel",
			Handler:    _ChannelService_DeleteChannel_Handler,
		},
		{
			MethodName: "JoinChannel",
			Handler:    _ChannelService_JoinChannel_Handler,
		},
		{
			MethodName: "LeaveChannel",
			Handler:    _ChannelService_LeaveChannel_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/app/api/channel/channel.proto",
}
