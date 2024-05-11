// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: internal/app/api/user/user.proto

package user

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UpdateUserRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Login          string `protobuf:"bytes,2,opt,name=login,proto3" json:"login,omitempty"`
	Name           string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Email          string `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	Password       string `protobuf:"bytes,5,opt,name=password,proto3" json:"password,omitempty"`
	AvatarPhotoUrl string `protobuf:"bytes,6,opt,name=avatar_photo_url,json=avatarPhotoUrl,proto3" json:"avatar_photo_url,omitempty"`
}

func (x *UpdateUserRequest) Reset() {
	*x = UpdateUserRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_user_user_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateUserRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateUserRequest) ProtoMessage() {}

func (x *UpdateUserRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_user_user_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateUserRequest.ProtoReflect.Descriptor instead.
func (*UpdateUserRequest) Descriptor() ([]byte, []int) {
	return file_internal_app_api_user_user_proto_rawDescGZIP(), []int{0}
}

func (x *UpdateUserRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UpdateUserRequest) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *UpdateUserRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UpdateUserRequest) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UpdateUserRequest) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *UpdateUserRequest) GetAvatarPhotoUrl() string {
	if x != nil {
		return x.AvatarPhotoUrl
	}
	return ""
}

type UserDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Login          string `protobuf:"bytes,2,opt,name=login,proto3" json:"login,omitempty"`
	Name           string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Email          string `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
	AvatarPhotoUrl string `protobuf:"bytes,6,opt,name=avatar_photo_url,json=avatarPhotoUrl,proto3" json:"avatar_photo_url,omitempty"`
}

func (x *UserDataResponse) Reset() {
	*x = UserDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_user_user_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UserDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UserDataResponse) ProtoMessage() {}

func (x *UserDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_user_user_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UserDataResponse.ProtoReflect.Descriptor instead.
func (*UserDataResponse) Descriptor() ([]byte, []int) {
	return file_internal_app_api_user_user_proto_rawDescGZIP(), []int{1}
}

func (x *UserDataResponse) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *UserDataResponse) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *UserDataResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *UserDataResponse) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *UserDataResponse) GetAvatarPhotoUrl() string {
	if x != nil {
		return x.AvatarPhotoUrl
	}
	return ""
}

type ErrorMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *ErrorMessage) Reset() {
	*x = ErrorMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_user_user_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrorMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorMessage) ProtoMessage() {}

func (x *ErrorMessage) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_user_user_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ErrorMessage.ProtoReflect.Descriptor instead.
func (*ErrorMessage) Descriptor() ([]byte, []int) {
	return file_internal_app_api_user_user_proto_rawDescGZIP(), []int{2}
}

func (x *ErrorMessage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

type GetUserByLoginRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Login string `protobuf:"bytes,1,opt,name=login,proto3" json:"login,omitempty"`
}

func (x *GetUserByLoginRequest) Reset() {
	*x = GetUserByLoginRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_user_user_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserByLoginRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserByLoginRequest) ProtoMessage() {}

func (x *GetUserByLoginRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_user_user_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserByLoginRequest.ProtoReflect.Descriptor instead.
func (*GetUserByLoginRequest) Descriptor() ([]byte, []int) {
	return file_internal_app_api_user_user_proto_rawDescGZIP(), []int{3}
}

func (x *GetUserByLoginRequest) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

type GetUserFriendsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *GetUserFriendsRequest) Reset() {
	*x = GetUserFriendsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_user_user_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserFriendsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserFriendsRequest) ProtoMessage() {}

func (x *GetUserFriendsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_user_user_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserFriendsRequest.ProtoReflect.Descriptor instead.
func (*GetUserFriendsRequest) Descriptor() ([]byte, []int) {
	return file_internal_app_api_user_user_proto_rawDescGZIP(), []int{4}
}

type GetUserFriendsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Friends []*Friend `protobuf:"bytes,1,rep,name=friends,proto3" json:"friends,omitempty"`
}

func (x *GetUserFriendsResponse) Reset() {
	*x = GetUserFriendsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_user_user_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetUserFriendsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetUserFriendsResponse) ProtoMessage() {}

func (x *GetUserFriendsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_user_user_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetUserFriendsResponse.ProtoReflect.Descriptor instead.
func (*GetUserFriendsResponse) Descriptor() ([]byte, []int) {
	return file_internal_app_api_user_user_proto_rawDescGZIP(), []int{5}
}

func (x *GetUserFriendsResponse) GetFriends() []*Friend {
	if x != nil {
		return x.Friends
	}
	return nil
}

type Friend struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint64 `protobuf:"varint,1,opt,name=userId,proto3" json:"userId,omitempty"`
	Login  string `protobuf:"bytes,2,opt,name=login,proto3" json:"login,omitempty"`
	Name   string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Email  string `protobuf:"bytes,4,opt,name=email,proto3" json:"email,omitempty"`
}

func (x *Friend) Reset() {
	*x = Friend{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_user_user_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Friend) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Friend) ProtoMessage() {}

func (x *Friend) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_user_user_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Friend.ProtoReflect.Descriptor instead.
func (*Friend) Descriptor() ([]byte, []int) {
	return file_internal_app_api_user_user_proto_rawDescGZIP(), []int{6}
}

func (x *Friend) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *Friend) GetLogin() string {
	if x != nil {
		return x.Login
	}
	return ""
}

func (x *Friend) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Friend) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

type AddToFriendByUserIdRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *AddToFriendByUserIdRequest) Reset() {
	*x = AddToFriendByUserIdRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_user_user_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddToFriendByUserIdRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddToFriendByUserIdRequest) ProtoMessage() {}

func (x *AddToFriendByUserIdRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_user_user_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddToFriendByUserIdRequest.ProtoReflect.Descriptor instead.
func (*AddToFriendByUserIdRequest) Descriptor() ([]byte, []int) {
	return file_internal_app_api_user_user_proto_rawDescGZIP(), []int{7}
}

func (x *AddToFriendByUserIdRequest) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type AcceptFriendInviteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InviteId uint64 `protobuf:"varint,1,opt,name=invite_id,json=inviteId,proto3" json:"invite_id,omitempty"`
}

func (x *AcceptFriendInviteRequest) Reset() {
	*x = AcceptFriendInviteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_user_user_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AcceptFriendInviteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AcceptFriendInviteRequest) ProtoMessage() {}

func (x *AcceptFriendInviteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_user_user_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AcceptFriendInviteRequest.ProtoReflect.Descriptor instead.
func (*AcceptFriendInviteRequest) Descriptor() ([]byte, []int) {
	return file_internal_app_api_user_user_proto_rawDescGZIP(), []int{8}
}

func (x *AcceptFriendInviteRequest) GetInviteId() uint64 {
	if x != nil {
		return x.InviteId
	}
	return 0
}

type DeclineFriendInviteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	InviteId uint64 `protobuf:"varint,1,opt,name=invite_id,json=inviteId,proto3" json:"invite_id,omitempty"`
}

func (x *DeclineFriendInviteRequest) Reset() {
	*x = DeclineFriendInviteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_user_user_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeclineFriendInviteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeclineFriendInviteRequest) ProtoMessage() {}

func (x *DeclineFriendInviteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_user_user_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeclineFriendInviteRequest.ProtoReflect.Descriptor instead.
func (*DeclineFriendInviteRequest) Descriptor() ([]byte, []int) {
	return file_internal_app_api_user_user_proto_rawDescGZIP(), []int{9}
}

func (x *DeclineFriendInviteRequest) GetInviteId() uint64 {
	if x != nil {
		return x.InviteId
	}
	return 0
}

type ActionResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *ActionResponse) Reset() {
	*x = ActionResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_user_user_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActionResponse) ProtoMessage() {}

func (x *ActionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_user_user_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActionResponse.ProtoReflect.Descriptor instead.
func (*ActionResponse) Descriptor() ([]byte, []int) {
	return file_internal_app_api_user_user_proto_rawDescGZIP(), []int{10}
}

func (x *ActionResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_internal_app_api_user_user_proto protoreflect.FileDescriptor

var file_internal_app_api_user_user_proto_rawDesc = []byte{
	0x0a, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x26, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e,
	0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x75,
	0x73, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x22, 0xa9, 0x01, 0x0a, 0x11, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d,
	0x61, 0x69, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x28, 0x0a, 0x10,
	0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x5f, 0x75, 0x72, 0x6c,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x50, 0x68,
	0x6f, 0x74, 0x6f, 0x55, 0x72, 0x6c, 0x22, 0x8c, 0x01, 0x0a, 0x10, 0x55, 0x73, 0x65, 0x72, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x6c,
	0x6f, 0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69,
	0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x28, 0x0a, 0x10, 0x61,
	0x76, 0x61, 0x74, 0x61, 0x72, 0x5f, 0x70, 0x68, 0x6f, 0x74, 0x6f, 0x5f, 0x75, 0x72, 0x6c, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x61, 0x76, 0x61, 0x74, 0x61, 0x72, 0x50, 0x68, 0x6f,
	0x74, 0x6f, 0x55, 0x72, 0x6c, 0x22, 0x28, 0x0a, 0x0c, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22,
	0x2d, 0x0a, 0x15, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x42, 0x79, 0x4c, 0x6f, 0x67, 0x69,
	0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x6f, 0x67, 0x69,
	0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x22, 0x17,
	0x0a, 0x15, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x73,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x62, 0x0a, 0x16, 0x47, 0x65, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x48, 0x0a, 0x07, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x2e, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e,
	0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e,
	0x75, 0x73, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x72, 0x69, 0x65,
	0x6e, 0x64, 0x52, 0x07, 0x66, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x73, 0x22, 0x60, 0x0a, 0x06, 0x46,
	0x72, 0x69, 0x65, 0x6e, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x14, 0x0a,
	0x05, 0x6c, 0x6f, 0x67, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6c, 0x6f,
	0x67, 0x69, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x22, 0x35, 0x0a,
	0x1a, 0x41, 0x64, 0x64, 0x54, 0x6f, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x42, 0x79, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x22, 0x38, 0x0a, 0x19, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x46, 0x72,
	0x69, 0x65, 0x6e, 0x64, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x49, 0x64, 0x22, 0x39,
	0x0a, 0x1a, 0x44, 0x65, 0x63, 0x6c, 0x69, 0x6e, 0x65, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x49,
	0x6e, 0x76, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09,
	0x69, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x08, 0x69, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x49, 0x64, 0x22, 0x2a, 0x0a, 0x0e, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73,
	0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75,
	0x63, 0x63, 0x65, 0x73, 0x73, 0x32, 0xf5, 0x06, 0x0a, 0x0b, 0x55, 0x73, 0x65, 0x72, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x83, 0x01, 0x0a, 0x0a, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x55, 0x73, 0x65, 0x72, 0x12, 0x39, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72,
	0x64, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x55, 0x73, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x38, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78,
	0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x8b, 0x01, 0x0a, 0x0e,
	0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x42, 0x79, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x12, 0x3d,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f,
	0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65, 0x72,
	0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x42,
	0x79, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x38, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e,
	0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x91, 0x01, 0x0a, 0x0e, 0x47, 0x65,
	0x74, 0x55, 0x73, 0x65, 0x72, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x73, 0x12, 0x3d, 0x2e, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78,
	0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x46, 0x72, 0x69,
	0x65, 0x6e, 0x64, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3e, 0x2e, 0x67, 0x69,
	0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70,
	0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x47, 0x65, 0x74, 0x55, 0x73, 0x65, 0x72, 0x46, 0x72, 0x69, 0x65,
	0x6e, 0x64, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x93, 0x01,
	0x0a, 0x13, 0x41, 0x64, 0x64, 0x54, 0x6f, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x42, 0x79, 0x55,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x42, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f,
	0x72, 0x64, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41,
	0x64, 0x64, 0x54, 0x6f, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72,
	0x49, 0x64, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x36, 0x2e, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64,
	0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x91, 0x01, 0x0a, 0x12, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x46, 0x72,
	0x69, 0x65, 0x6e, 0x64, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x12, 0x41, 0x2e, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e,
	0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x63, 0x65, 0x70, 0x74, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64,
	0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x36, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e,
	0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e,
	0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x93, 0x01, 0x0a, 0x13, 0x44, 0x65, 0x63, 0x6c,
	0x69, 0x6e, 0x65, 0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x12,
	0x42, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78,
	0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x75, 0x73, 0x65,
	0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x44, 0x65, 0x63, 0x6c, 0x69, 0x6e, 0x65,
	0x46, 0x72, 0x69, 0x65, 0x6e, 0x64, 0x49, 0x6e, 0x76, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x36, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64,
	0x2e, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0a, 0x5a,
	0x08, 0x61, 0x70, 0x69, 0x2f, 0x75, 0x73, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_internal_app_api_user_user_proto_rawDescOnce sync.Once
	file_internal_app_api_user_user_proto_rawDescData = file_internal_app_api_user_user_proto_rawDesc
)

func file_internal_app_api_user_user_proto_rawDescGZIP() []byte {
	file_internal_app_api_user_user_proto_rawDescOnce.Do(func() {
		file_internal_app_api_user_user_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_app_api_user_user_proto_rawDescData)
	})
	return file_internal_app_api_user_user_proto_rawDescData
}

var file_internal_app_api_user_user_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_internal_app_api_user_user_proto_goTypes = []interface{}{
	(*UpdateUserRequest)(nil),          // 0: github.com.Nixonxp.discord.user.api.v1.UpdateUserRequest
	(*UserDataResponse)(nil),           // 1: github.com.Nixonxp.discord.user.api.v1.UserDataResponse
	(*ErrorMessage)(nil),               // 2: github.com.Nixonxp.discord.user.api.v1.ErrorMessage
	(*GetUserByLoginRequest)(nil),      // 3: github.com.Nixonxp.discord.user.api.v1.GetUserByLoginRequest
	(*GetUserFriendsRequest)(nil),      // 4: github.com.Nixonxp.discord.user.api.v1.GetUserFriendsRequest
	(*GetUserFriendsResponse)(nil),     // 5: github.com.Nixonxp.discord.user.api.v1.GetUserFriendsResponse
	(*Friend)(nil),                     // 6: github.com.Nixonxp.discord.user.api.v1.Friend
	(*AddToFriendByUserIdRequest)(nil), // 7: github.com.Nixonxp.discord.user.api.v1.AddToFriendByUserIdRequest
	(*AcceptFriendInviteRequest)(nil),  // 8: github.com.Nixonxp.discord.user.api.v1.AcceptFriendInviteRequest
	(*DeclineFriendInviteRequest)(nil), // 9: github.com.Nixonxp.discord.user.api.v1.DeclineFriendInviteRequest
	(*ActionResponse)(nil),             // 10: github.com.Nixonxp.discord.user.api.v1.ActionResponse
}
var file_internal_app_api_user_user_proto_depIdxs = []int32{
	6,  // 0: github.com.Nixonxp.discord.user.api.v1.GetUserFriendsResponse.friends:type_name -> github.com.Nixonxp.discord.user.api.v1.Friend
	0,  // 1: github.com.Nixonxp.discord.user.api.v1.UserService.UpdateUser:input_type -> github.com.Nixonxp.discord.user.api.v1.UpdateUserRequest
	3,  // 2: github.com.Nixonxp.discord.user.api.v1.UserService.GetUserByLogin:input_type -> github.com.Nixonxp.discord.user.api.v1.GetUserByLoginRequest
	4,  // 3: github.com.Nixonxp.discord.user.api.v1.UserService.GetUserFriends:input_type -> github.com.Nixonxp.discord.user.api.v1.GetUserFriendsRequest
	7,  // 4: github.com.Nixonxp.discord.user.api.v1.UserService.AddToFriendByUserId:input_type -> github.com.Nixonxp.discord.user.api.v1.AddToFriendByUserIdRequest
	8,  // 5: github.com.Nixonxp.discord.user.api.v1.UserService.AcceptFriendInvite:input_type -> github.com.Nixonxp.discord.user.api.v1.AcceptFriendInviteRequest
	9,  // 6: github.com.Nixonxp.discord.user.api.v1.UserService.DeclineFriendInvite:input_type -> github.com.Nixonxp.discord.user.api.v1.DeclineFriendInviteRequest
	1,  // 7: github.com.Nixonxp.discord.user.api.v1.UserService.UpdateUser:output_type -> github.com.Nixonxp.discord.user.api.v1.UserDataResponse
	1,  // 8: github.com.Nixonxp.discord.user.api.v1.UserService.GetUserByLogin:output_type -> github.com.Nixonxp.discord.user.api.v1.UserDataResponse
	5,  // 9: github.com.Nixonxp.discord.user.api.v1.UserService.GetUserFriends:output_type -> github.com.Nixonxp.discord.user.api.v1.GetUserFriendsResponse
	10, // 10: github.com.Nixonxp.discord.user.api.v1.UserService.AddToFriendByUserId:output_type -> github.com.Nixonxp.discord.user.api.v1.ActionResponse
	10, // 11: github.com.Nixonxp.discord.user.api.v1.UserService.AcceptFriendInvite:output_type -> github.com.Nixonxp.discord.user.api.v1.ActionResponse
	10, // 12: github.com.Nixonxp.discord.user.api.v1.UserService.DeclineFriendInvite:output_type -> github.com.Nixonxp.discord.user.api.v1.ActionResponse
	7,  // [7:13] is the sub-list for method output_type
	1,  // [1:7] is the sub-list for method input_type
	1,  // [1:1] is the sub-list for extension type_name
	1,  // [1:1] is the sub-list for extension extendee
	0,  // [0:1] is the sub-list for field type_name
}

func init() { file_internal_app_api_user_user_proto_init() }
func file_internal_app_api_user_user_proto_init() {
	if File_internal_app_api_user_user_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_app_api_user_user_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateUserRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_app_api_user_user_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UserDataResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_app_api_user_user_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ErrorMessage); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_app_api_user_user_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserByLoginRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_app_api_user_user_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserFriendsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_app_api_user_user_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetUserFriendsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_app_api_user_user_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Friend); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_app_api_user_user_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddToFriendByUserIdRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_app_api_user_user_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AcceptFriendInviteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_app_api_user_user_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeclineFriendInviteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_app_api_user_user_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActionResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_app_api_user_user_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_app_api_user_user_proto_goTypes,
		DependencyIndexes: file_internal_app_api_user_user_proto_depIdxs,
		MessageInfos:      file_internal_app_api_user_user_proto_msgTypes,
	}.Build()
	File_internal_app_api_user_user_proto = out.File
	file_internal_app_api_user_user_proto_rawDesc = nil
	file_internal_app_api_user_user_proto_goTypes = nil
	file_internal_app_api_user_user_proto_depIdxs = nil
}
