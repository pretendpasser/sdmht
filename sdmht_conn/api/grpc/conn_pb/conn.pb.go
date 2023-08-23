// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v4.22.3
// source: conn.proto

package conn_pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	error1 "sdmht/lib/protobuf/types/error"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type KickClientReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
}

func (x *KickClientReq) Reset() {
	*x = KickClientReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *KickClientReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*KickClientReq) ProtoMessage() {}

func (x *KickClientReq) ProtoReflect() protoreflect.Message {
	mi := &file_conn_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use KickClientReq.ProtoReflect.Descriptor instead.
func (*KickClientReq) Descriptor() ([]byte, []int) {
	return file_conn_proto_rawDescGZIP(), []int{0}
}

func (x *KickClientReq) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

type ClientEventReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Event  *Event `protobuf:"bytes,2,opt,name=event,proto3" json:"event,omitempty"`
}

func (x *ClientEventReq) Reset() {
	*x = ClientEventReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientEventReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientEventReq) ProtoMessage() {}

func (x *ClientEventReq) ProtoReflect() protoreflect.Message {
	mi := &file_conn_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientEventReq.ProtoReflect.Descriptor instead.
func (*ClientEventReq) Descriptor() ([]byte, []int) {
	return file_conn_proto_rawDescGZIP(), []int{1}
}

func (x *ClientEventReq) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ClientEventReq) GetEvent() *Event {
	if x != nil {
		return x.Event
	}
	return nil
}

type Event struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type    string                 `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`       // event_type
	Content string                 `protobuf:"bytes,2,opt,name=content,proto3" json:"content,omitempty"` // json
	AtTime  *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=at_time,json=atTime,proto3" json:"at_time,omitempty"`
}

func (x *Event) Reset() {
	*x = Event{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Event) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Event) ProtoMessage() {}

func (x *Event) ProtoReflect() protoreflect.Message {
	mi := &file_conn_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Event.ProtoReflect.Descriptor instead.
func (*Event) Descriptor() ([]byte, []int) {
	return file_conn_proto_rawDescGZIP(), []int{2}
}

func (x *Event) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Event) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *Event) GetAtTime() *timestamppb.Timestamp {
	if x != nil {
		return x.AtTime
	}
	return nil
}

type DispatchEventToClientReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ClientReply *ClientReply  `protobuf:"bytes,1,opt,name=client_reply,json=clientReply,proto3" json:"client_reply,omitempty"`
	Err         *error1.Error `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *DispatchEventToClientReply) Reset() {
	*x = DispatchEventToClientReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DispatchEventToClientReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DispatchEventToClientReply) ProtoMessage() {}

func (x *DispatchEventToClientReply) ProtoReflect() protoreflect.Message {
	mi := &file_conn_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DispatchEventToClientReply.ProtoReflect.Descriptor instead.
func (*DispatchEventToClientReply) Descriptor() ([]byte, []int) {
	return file_conn_proto_rawDescGZIP(), []int{3}
}

func (x *DispatchEventToClientReply) GetClientReply() *ClientReply {
	if x != nil {
		return x.ClientReply
	}
	return nil
}

func (x *DispatchEventToClientReply) GetErr() *error1.Error {
	if x != nil {
		return x.Err
	}
	return nil
}

type ClientReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId uint64 `protobuf:"varint,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	Ok     bool   `protobuf:"varint,2,opt,name=ok,proto3" json:"ok,omitempty"`
	Err    string `protobuf:"bytes,3,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *ClientReply) Reset() {
	*x = ClientReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClientReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClientReply) ProtoMessage() {}

func (x *ClientReply) ProtoReflect() protoreflect.Message {
	mi := &file_conn_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClientReply.ProtoReflect.Descriptor instead.
func (*ClientReply) Descriptor() ([]byte, []int) {
	return file_conn_proto_rawDescGZIP(), []int{4}
}

func (x *ClientReply) GetUserId() uint64 {
	if x != nil {
		return x.UserId
	}
	return 0
}

func (x *ClientReply) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *ClientReply) GetErr() string {
	if x != nil {
		return x.Err
	}
	return ""
}

type CommonReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Err *error1.Error `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
}

func (x *CommonReply) Reset() {
	*x = CommonReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conn_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommonReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommonReply) ProtoMessage() {}

func (x *CommonReply) ProtoReflect() protoreflect.Message {
	mi := &file_conn_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommonReply.ProtoReflect.Descriptor instead.
func (*CommonReply) Descriptor() ([]byte, []int) {
	return file_conn_proto_rawDescGZIP(), []int{5}
}

func (x *CommonReply) GetErr() *error1.Error {
	if x != nil {
		return x.Err
	}
	return nil
}

var File_conn_proto protoreflect.FileDescriptor

var file_conn_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x6f,
	0x6e, 0x6e, 0x5f, 0x70, 0x62, 0x1a, 0x18, 0x6c, 0x69, 0x62, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a,
	0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x28, 0x0a, 0x0d, 0x4b, 0x69, 0x63, 0x6b, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x71, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x4f, 0x0a, 0x0e, 0x43, 0x6c,
	0x69, 0x65, 0x6e, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x12, 0x17, 0x0a, 0x07,
	0x75, 0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x12, 0x24, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x5f, 0x70, 0x62, 0x2e, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x22, 0x6a, 0x0a, 0x05, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x12, 0x33, 0x0a, 0x07, 0x61, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x06, 0x61, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x73, 0x0a, 0x1a, 0x44, 0x69, 0x73, 0x70, 0x61,
	0x74, 0x63, 0x68, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x37, 0x0a, 0x0c, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f,
	0x72, 0x65, 0x70, 0x6c, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f,
	0x6e, 0x6e, 0x5f, 0x70, 0x62, 0x2e, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x70, 0x6c,
	0x79, 0x52, 0x0b, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1c,
	0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x6c, 0x69,
	0x62, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x03, 0x65, 0x72, 0x72, 0x22, 0x48, 0x0a, 0x0b,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x17, 0x0a, 0x07, 0x75,
	0x73, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x75, 0x73,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x02, 0x6f, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x65, 0x72, 0x72, 0x22, 0x2b, 0x0a, 0x0b, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x1c, 0x0a, 0x03, 0x65, 0x72, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x6c, 0x69, 0x62, 0x2e, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x52, 0x03,
	0x65, 0x72, 0x72, 0x32, 0x99, 0x01, 0x0a, 0x04, 0x43, 0x6f, 0x6e, 0x6e, 0x12, 0x55, 0x0a, 0x15,
	0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63, 0x68, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x12, 0x17, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x5f, 0x70, 0x62, 0x2e,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x23,
	0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x5f, 0x70, 0x62, 0x2e, 0x44, 0x69, 0x73, 0x70, 0x61, 0x74, 0x63,
	0x68, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x54, 0x6f, 0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65,
	0x70, 0x6c, 0x79, 0x12, 0x3a, 0x0a, 0x0a, 0x4b, 0x69, 0x63, 0x6b, 0x43, 0x6c, 0x69, 0x65, 0x6e,
	0x74, 0x12, 0x16, 0x2e, 0x63, 0x6f, 0x6e, 0x6e, 0x5f, 0x70, 0x62, 0x2e, 0x4b, 0x69, 0x63, 0x6b,
	0x43, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x1a, 0x14, 0x2e, 0x63, 0x6f, 0x6e, 0x6e,
	0x5f, 0x70, 0x62, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42,
	0x23, 0x5a, 0x21, 0x73, 0x64, 0x6d, 0x68, 0x74, 0x2f, 0x73, 0x64, 0x6d, 0x68, 0x74, 0x5f, 0x63,
	0x6f, 0x6e, 0x6e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x63, 0x6f, 0x6e,
	0x6e, 0x5f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_conn_proto_rawDescOnce sync.Once
	file_conn_proto_rawDescData = file_conn_proto_rawDesc
)

func file_conn_proto_rawDescGZIP() []byte {
	file_conn_proto_rawDescOnce.Do(func() {
		file_conn_proto_rawDescData = protoimpl.X.CompressGZIP(file_conn_proto_rawDescData)
	})
	return file_conn_proto_rawDescData
}

var file_conn_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_conn_proto_goTypes = []interface{}{
	(*KickClientReq)(nil),              // 0: conn_pb.KickClientReq
	(*ClientEventReq)(nil),             // 1: conn_pb.ClientEventReq
	(*Event)(nil),                      // 2: conn_pb.Event
	(*DispatchEventToClientReply)(nil), // 3: conn_pb.DispatchEventToClientReply
	(*ClientReply)(nil),                // 4: conn_pb.ClientReply
	(*CommonReply)(nil),                // 5: conn_pb.CommonReply
	(*timestamppb.Timestamp)(nil),      // 6: google.protobuf.Timestamp
	(*error1.Error)(nil),               // 7: lib.Error
}
var file_conn_proto_depIdxs = []int32{
	2, // 0: conn_pb.ClientEventReq.event:type_name -> conn_pb.Event
	6, // 1: conn_pb.Event.at_time:type_name -> google.protobuf.Timestamp
	4, // 2: conn_pb.DispatchEventToClientReply.client_reply:type_name -> conn_pb.ClientReply
	7, // 3: conn_pb.DispatchEventToClientReply.err:type_name -> lib.Error
	7, // 4: conn_pb.CommonReply.err:type_name -> lib.Error
	1, // 5: conn_pb.Conn.DispatchEventToClient:input_type -> conn_pb.ClientEventReq
	0, // 6: conn_pb.Conn.KickClient:input_type -> conn_pb.KickClientReq
	3, // 7: conn_pb.Conn.DispatchEventToClient:output_type -> conn_pb.DispatchEventToClientReply
	5, // 8: conn_pb.Conn.KickClient:output_type -> conn_pb.CommonReply
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_conn_proto_init() }
func file_conn_proto_init() {
	if File_conn_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_conn_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*KickClientReq); i {
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
		file_conn_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientEventReq); i {
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
		file_conn_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Event); i {
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
		file_conn_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DispatchEventToClientReply); i {
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
		file_conn_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClientReply); i {
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
		file_conn_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommonReply); i {
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
			RawDescriptor: file_conn_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_conn_proto_goTypes,
		DependencyIndexes: file_conn_proto_depIdxs,
		MessageInfos:      file_conn_proto_msgTypes,
	}.Build()
	File_conn_proto = out.File
	file_conn_proto_rawDesc = nil
	file_conn_proto_goTypes = nil
	file_conn_proto_depIdxs = nil
}
