// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: internal/app/api/channel/channel.proto

package channel

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

type AddChannelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *AddChannelRequest) Reset() {
	*x = AddChannelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_channel_channel_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AddChannelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AddChannelRequest) ProtoMessage() {}

func (x *AddChannelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_channel_channel_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AddChannelRequest.ProtoReflect.Descriptor instead.
func (*AddChannelRequest) Descriptor() ([]byte, []int) {
	return file_internal_app_api_channel_channel_proto_rawDescGZIP(), []int{0}
}

func (x *AddChannelRequest) GetName() string {
	if x != nil {
		return x.Name
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
		mi := &file_internal_app_api_channel_channel_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ErrorMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ErrorMessage) ProtoMessage() {}

func (x *ErrorMessage) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_channel_channel_proto_msgTypes[1]
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
	return file_internal_app_api_channel_channel_proto_rawDescGZIP(), []int{1}
}

func (x *ErrorMessage) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
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
		mi := &file_internal_app_api_channel_channel_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActionResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActionResponse) ProtoMessage() {}

func (x *ActionResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_channel_channel_proto_msgTypes[2]
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
	return file_internal_app_api_channel_channel_proto_rawDescGZIP(), []int{2}
}

func (x *ActionResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

type DeleteChannelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId string `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
}

func (x *DeleteChannelRequest) Reset() {
	*x = DeleteChannelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_channel_channel_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteChannelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteChannelRequest) ProtoMessage() {}

func (x *DeleteChannelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_channel_channel_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteChannelRequest.ProtoReflect.Descriptor instead.
func (*DeleteChannelRequest) Descriptor() ([]byte, []int) {
	return file_internal_app_api_channel_channel_proto_rawDescGZIP(), []int{3}
}

func (x *DeleteChannelRequest) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

type JoinChannelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId string `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
}

func (x *JoinChannelRequest) Reset() {
	*x = JoinChannelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_channel_channel_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *JoinChannelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*JoinChannelRequest) ProtoMessage() {}

func (x *JoinChannelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_channel_channel_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use JoinChannelRequest.ProtoReflect.Descriptor instead.
func (*JoinChannelRequest) Descriptor() ([]byte, []int) {
	return file_internal_app_api_channel_channel_proto_rawDescGZIP(), []int{4}
}

func (x *JoinChannelRequest) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

type LeaveChannelRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ChannelId string `protobuf:"bytes,1,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
}

func (x *LeaveChannelRequest) Reset() {
	*x = LeaveChannelRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_app_api_channel_channel_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LeaveChannelRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LeaveChannelRequest) ProtoMessage() {}

func (x *LeaveChannelRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_app_api_channel_channel_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LeaveChannelRequest.ProtoReflect.Descriptor instead.
func (*LeaveChannelRequest) Descriptor() ([]byte, []int) {
	return file_internal_app_api_channel_channel_proto_rawDescGZIP(), []int{5}
}

func (x *LeaveChannelRequest) GetChannelId() string {
	if x != nil {
		return x.ChannelId
	}
	return ""
}

var File_internal_app_api_channel_channel_proto protoreflect.FileDescriptor

var file_internal_app_api_channel_channel_proto_rawDesc = []byte{
	0x0a, 0x26, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x70, 0x2f, 0x61,
	0x70, 0x69, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x29, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73,
	0x63, 0x6f, 0x72, 0x64, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x22, 0x27, 0x0a, 0x11, 0x41, 0x64, 0x64, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x28, 0x0a, 0x0c,
	0x45, 0x72, 0x72, 0x6f, 0x72, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x2a, 0x0a, 0x0e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65,
	0x73, 0x73, 0x22, 0x35, 0x0a, 0x14, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64, 0x22, 0x33, 0x0a, 0x12, 0x4a, 0x6f, 0x69,
	0x6e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x49, 0x64, 0x22, 0x34,
	0x0a, 0x13, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x68, 0x61, 0x6e, 0x6e,
	0x65, 0x6c, 0x49, 0x64, 0x32, 0xc4, 0x04, 0x0a, 0x0e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x87, 0x01, 0x0a, 0x0a, 0x41, 0x64, 0x64, 0x43,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x3c, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63,
	0x6f, 0x72, 0x64, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x41, 0x64, 0x64, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x39, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72,
	0x64, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x8d, 0x01, 0x0a, 0x0d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x12, 0x3f, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64,
	0x2e, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x39, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72,
	0x64, 0x2e, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31,
	0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x12, 0x89, 0x01, 0x0a, 0x0b, 0x4a, 0x6f, 0x69, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65,
	0x6c, 0x12, 0x3d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e,
	0x69, 0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x63,
	0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x4a, 0x6f,
	0x69, 0x6e, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x39, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69,
	0x78, 0x6f, 0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x63, 0x68,
	0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x12, 0x8b, 0x01,
	0x0a, 0x0c, 0x4c, 0x65, 0x61, 0x76, 0x65, 0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x12, 0x3e,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f,
	0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x63, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x65, 0x61, 0x76, 0x65,
	0x43, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x39,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x4e, 0x69, 0x78, 0x6f,
	0x6e, 0x78, 0x70, 0x2e, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x72, 0x64, 0x2e, 0x63, 0x68, 0x61, 0x6e,
	0x6e, 0x65, 0x6c, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0d, 0x5a, 0x0b, 0x61,
	0x70, 0x69, 0x2f, 0x63, 0x68, 0x61, 0x6e, 0x6e, 0x65, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_internal_app_api_channel_channel_proto_rawDescOnce sync.Once
	file_internal_app_api_channel_channel_proto_rawDescData = file_internal_app_api_channel_channel_proto_rawDesc
)

func file_internal_app_api_channel_channel_proto_rawDescGZIP() []byte {
	file_internal_app_api_channel_channel_proto_rawDescOnce.Do(func() {
		file_internal_app_api_channel_channel_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_app_api_channel_channel_proto_rawDescData)
	})
	return file_internal_app_api_channel_channel_proto_rawDescData
}

var file_internal_app_api_channel_channel_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_internal_app_api_channel_channel_proto_goTypes = []interface{}{
	(*AddChannelRequest)(nil),    // 0: github.com.Nixonxp.discord.channel.api.v1.AddChannelRequest
	(*ErrorMessage)(nil),         // 1: github.com.Nixonxp.discord.channel.api.v1.ErrorMessage
	(*ActionResponse)(nil),       // 2: github.com.Nixonxp.discord.channel.api.v1.ActionResponse
	(*DeleteChannelRequest)(nil), // 3: github.com.Nixonxp.discord.channel.api.v1.DeleteChannelRequest
	(*JoinChannelRequest)(nil),   // 4: github.com.Nixonxp.discord.channel.api.v1.JoinChannelRequest
	(*LeaveChannelRequest)(nil),  // 5: github.com.Nixonxp.discord.channel.api.v1.LeaveChannelRequest
}
var file_internal_app_api_channel_channel_proto_depIdxs = []int32{
	0, // 0: github.com.Nixonxp.discord.channel.api.v1.ChannelService.AddChannel:input_type -> github.com.Nixonxp.discord.channel.api.v1.AddChannelRequest
	3, // 1: github.com.Nixonxp.discord.channel.api.v1.ChannelService.DeleteChannel:input_type -> github.com.Nixonxp.discord.channel.api.v1.DeleteChannelRequest
	4, // 2: github.com.Nixonxp.discord.channel.api.v1.ChannelService.JoinChannel:input_type -> github.com.Nixonxp.discord.channel.api.v1.JoinChannelRequest
	5, // 3: github.com.Nixonxp.discord.channel.api.v1.ChannelService.LeaveChannel:input_type -> github.com.Nixonxp.discord.channel.api.v1.LeaveChannelRequest
	2, // 4: github.com.Nixonxp.discord.channel.api.v1.ChannelService.AddChannel:output_type -> github.com.Nixonxp.discord.channel.api.v1.ActionResponse
	2, // 5: github.com.Nixonxp.discord.channel.api.v1.ChannelService.DeleteChannel:output_type -> github.com.Nixonxp.discord.channel.api.v1.ActionResponse
	2, // 6: github.com.Nixonxp.discord.channel.api.v1.ChannelService.JoinChannel:output_type -> github.com.Nixonxp.discord.channel.api.v1.ActionResponse
	2, // 7: github.com.Nixonxp.discord.channel.api.v1.ChannelService.LeaveChannel:output_type -> github.com.Nixonxp.discord.channel.api.v1.ActionResponse
	4, // [4:8] is the sub-list for method output_type
	0, // [0:4] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_app_api_channel_channel_proto_init() }
func file_internal_app_api_channel_channel_proto_init() {
	if File_internal_app_api_channel_channel_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_app_api_channel_channel_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AddChannelRequest); i {
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
		file_internal_app_api_channel_channel_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
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
		file_internal_app_api_channel_channel_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
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
		file_internal_app_api_channel_channel_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteChannelRequest); i {
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
		file_internal_app_api_channel_channel_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*JoinChannelRequest); i {
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
		file_internal_app_api_channel_channel_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LeaveChannelRequest); i {
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
			RawDescriptor: file_internal_app_api_channel_channel_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_app_api_channel_channel_proto_goTypes,
		DependencyIndexes: file_internal_app_api_channel_channel_proto_depIdxs,
		MessageInfos:      file_internal_app_api_channel_channel_proto_msgTypes,
	}.Build()
	File_internal_app_api_channel_channel_proto = out.File
	file_internal_app_api_channel_channel_proto_rawDesc = nil
	file_internal_app_api_channel_channel_proto_goTypes = nil
	file_internal_app_api_channel_channel_proto_depIdxs = nil
}
