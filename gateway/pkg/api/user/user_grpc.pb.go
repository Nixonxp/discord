// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.25.1
// source: internal/app/api/user/user.proto

package user

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
	UserService_UpdateUser_FullMethodName          = "/github.com.Nixonxp.discord.user.api.v1.UserService/UpdateUser"
	UserService_GetUserByLogin_FullMethodName      = "/github.com.Nixonxp.discord.user.api.v1.UserService/GetUserByLogin"
	UserService_GetUserFriends_FullMethodName      = "/github.com.Nixonxp.discord.user.api.v1.UserService/GetUserFriends"
	UserService_AddToFriendByUserId_FullMethodName = "/github.com.Nixonxp.discord.user.api.v1.UserService/AddToFriendByUserId"
	UserService_AcceptFriendInvite_FullMethodName  = "/github.com.Nixonxp.discord.user.api.v1.UserService/AcceptFriendInvite"
	UserService_DeclineFriendInvite_FullMethodName = "/github.com.Nixonxp.discord.user.api.v1.UserService/DeclineFriendInvite"
)

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UserDataResponse, error)
	GetUserByLogin(ctx context.Context, in *GetUserByLoginRequest, opts ...grpc.CallOption) (*UserDataResponse, error)
	GetUserFriends(ctx context.Context, in *GetUserFriendsRequest, opts ...grpc.CallOption) (*GetUserFriendsResponse, error)
	AddToFriendByUserId(ctx context.Context, in *AddToFriendByUserIdRequest, opts ...grpc.CallOption) (*ActionResponse, error)
	AcceptFriendInvite(ctx context.Context, in *AcceptFriendInviteRequest, opts ...grpc.CallOption) (*ActionResponse, error)
	DeclineFriendInvite(ctx context.Context, in *DeclineFriendInviteRequest, opts ...grpc.CallOption) (*ActionResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*UserDataResponse, error) {
	out := new(UserDataResponse)
	err := c.cc.Invoke(ctx, UserService_UpdateUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserByLogin(ctx context.Context, in *GetUserByLoginRequest, opts ...grpc.CallOption) (*UserDataResponse, error) {
	out := new(UserDataResponse)
	err := c.cc.Invoke(ctx, UserService_GetUserByLogin_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) GetUserFriends(ctx context.Context, in *GetUserFriendsRequest, opts ...grpc.CallOption) (*GetUserFriendsResponse, error) {
	out := new(GetUserFriendsResponse)
	err := c.cc.Invoke(ctx, UserService_GetUserFriends_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) AddToFriendByUserId(ctx context.Context, in *AddToFriendByUserIdRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, UserService_AddToFriendByUserId_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) AcceptFriendInvite(ctx context.Context, in *AcceptFriendInviteRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, UserService_AcceptFriendInvite_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) DeclineFriendInvite(ctx context.Context, in *DeclineFriendInviteRequest, opts ...grpc.CallOption) (*ActionResponse, error) {
	out := new(ActionResponse)
	err := c.cc.Invoke(ctx, UserService_DeclineFriendInvite_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	UpdateUser(context.Context, *UpdateUserRequest) (*UserDataResponse, error)
	GetUserByLogin(context.Context, *GetUserByLoginRequest) (*UserDataResponse, error)
	GetUserFriends(context.Context, *GetUserFriendsRequest) (*GetUserFriendsResponse, error)
	AddToFriendByUserId(context.Context, *AddToFriendByUserIdRequest) (*ActionResponse, error)
	AcceptFriendInvite(context.Context, *AcceptFriendInviteRequest) (*ActionResponse, error)
	DeclineFriendInvite(context.Context, *DeclineFriendInviteRequest) (*ActionResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) UpdateUser(context.Context, *UpdateUserRequest) (*UserDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUser not implemented")
}
func (UnimplementedUserServiceServer) GetUserByLogin(context.Context, *GetUserByLoginRequest) (*UserDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByLogin not implemented")
}
func (UnimplementedUserServiceServer) GetUserFriends(context.Context, *GetUserFriendsRequest) (*GetUserFriendsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserFriends not implemented")
}
func (UnimplementedUserServiceServer) AddToFriendByUserId(context.Context, *AddToFriendByUserIdRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddToFriendByUserId not implemented")
}
func (UnimplementedUserServiceServer) AcceptFriendInvite(context.Context, *AcceptFriendInviteRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AcceptFriendInvite not implemented")
}
func (UnimplementedUserServiceServer) DeclineFriendInvite(context.Context, *DeclineFriendInviteRequest) (*ActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeclineFriendInvite not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_UpdateUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).UpdateUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_UpdateUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).UpdateUser(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserByLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserByLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserByLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetUserByLogin_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserByLogin(ctx, req.(*GetUserByLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_GetUserFriends_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserFriendsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).GetUserFriends(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_GetUserFriends_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).GetUserFriends(ctx, req.(*GetUserFriendsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_AddToFriendByUserId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddToFriendByUserIdRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).AddToFriendByUserId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_AddToFriendByUserId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).AddToFriendByUserId(ctx, req.(*AddToFriendByUserIdRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_AcceptFriendInvite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AcceptFriendInviteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).AcceptFriendInvite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_AcceptFriendInvite_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).AcceptFriendInvite(ctx, req.(*AcceptFriendInviteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_DeclineFriendInvite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeclineFriendInviteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).DeclineFriendInvite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserService_DeclineFriendInvite_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).DeclineFriendInvite(ctx, req.(*DeclineFriendInviteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "github.com.Nixonxp.discord.user.api.v1.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateUser",
			Handler:    _UserService_UpdateUser_Handler,
		},
		{
			MethodName: "GetUserByLogin",
			Handler:    _UserService_GetUserByLogin_Handler,
		},
		{
			MethodName: "GetUserFriends",
			Handler:    _UserService_GetUserFriends_Handler,
		},
		{
			MethodName: "AddToFriendByUserId",
			Handler:    _UserService_AddToFriendByUserId_Handler,
		},
		{
			MethodName: "AcceptFriendInvite",
			Handler:    _UserService_AcceptFriendInvite_Handler,
		},
		{
			MethodName: "DeclineFriendInvite",
			Handler:    _UserService_DeclineFriendInvite_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/app/api/user/user.proto",
}
