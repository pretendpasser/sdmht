// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package signaling

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

// SignalingClient is the client API for Signaling service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SignalingClient interface {
	Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginReply, error)
	NewMatch(ctx context.Context, in *NewMatchReq, opts ...grpc.CallOption) (*NewMatchReply, error)
	KeepAlive(ctx context.Context, in *KeepAliveReq, opts ...grpc.CallOption) (*CommonReply, error)
	Offline(ctx context.Context, in *LogoutReq, opts ...grpc.CallOption) (*CommonReply, error)
}

type signalingClient struct {
	cc grpc.ClientConnInterface
}

func NewSignalingClient(cc grpc.ClientConnInterface) SignalingClient {
	return &signalingClient{cc}
}

func (c *signalingClient) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginReply, error) {
	out := new(LoginReply)
	err := c.cc.Invoke(ctx, "/signaling_pb.Signaling/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *signalingClient) NewMatch(ctx context.Context, in *NewMatchReq, opts ...grpc.CallOption) (*NewMatchReply, error) {
	out := new(NewMatchReply)
	err := c.cc.Invoke(ctx, "/signaling_pb.Signaling/NewMatch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *signalingClient) KeepAlive(ctx context.Context, in *KeepAliveReq, opts ...grpc.CallOption) (*CommonReply, error) {
	out := new(CommonReply)
	err := c.cc.Invoke(ctx, "/signaling_pb.Signaling/KeepAlive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *signalingClient) Offline(ctx context.Context, in *LogoutReq, opts ...grpc.CallOption) (*CommonReply, error) {
	out := new(CommonReply)
	err := c.cc.Invoke(ctx, "/signaling_pb.Signaling/Offline", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SignalingServer is the server API for Signaling service.
// All implementations must embed UnimplementedSignalingServer
// for forward compatibility
type SignalingServer interface {
	Login(context.Context, *LoginReq) (*LoginReply, error)
	NewMatch(context.Context, *NewMatchReq) (*NewMatchReply, error)
	KeepAlive(context.Context, *KeepAliveReq) (*CommonReply, error)
	Offline(context.Context, *LogoutReq) (*CommonReply, error)
	mustEmbedUnimplementedSignalingServer()
}

// UnimplementedSignalingServer must be embedded to have forward compatible implementations.
type UnimplementedSignalingServer struct {
}

func (UnimplementedSignalingServer) Login(context.Context, *LoginReq) (*LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedSignalingServer) NewMatch(context.Context, *NewMatchReq) (*NewMatchReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NewMatch not implemented")
}
func (UnimplementedSignalingServer) KeepAlive(context.Context, *KeepAliveReq) (*CommonReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method KeepAlive not implemented")
}
func (UnimplementedSignalingServer) Offline(context.Context, *LogoutReq) (*CommonReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Offline not implemented")
}
func (UnimplementedSignalingServer) mustEmbedUnimplementedSignalingServer() {}

// UnsafeSignalingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SignalingServer will
// result in compilation errors.
type UnsafeSignalingServer interface {
	mustEmbedUnimplementedSignalingServer()
}

func RegisterSignalingServer(s grpc.ServiceRegistrar, srv SignalingServer) {
	s.RegisterService(&Signaling_ServiceDesc, srv)
}

func _Signaling_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SignalingServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/signaling_pb.Signaling/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SignalingServer).Login(ctx, req.(*LoginReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Signaling_NewMatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NewMatchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SignalingServer).NewMatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/signaling_pb.Signaling/NewMatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SignalingServer).NewMatch(ctx, req.(*NewMatchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Signaling_KeepAlive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeepAliveReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SignalingServer).KeepAlive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/signaling_pb.Signaling/KeepAlive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SignalingServer).KeepAlive(ctx, req.(*KeepAliveReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Signaling_Offline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SignalingServer).Offline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/signaling_pb.Signaling/Offline",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SignalingServer).Offline(ctx, req.(*LogoutReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Signaling_ServiceDesc is the grpc.ServiceDesc for Signaling service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Signaling_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "signaling_pb.Signaling",
	HandlerType: (*SignalingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _Signaling_Login_Handler,
		},
		{
			MethodName: "NewMatch",
			Handler:    _Signaling_NewMatch_Handler,
		},
		{
			MethodName: "KeepAlive",
			Handler:    _Signaling_KeepAlive_Handler,
		},
		{
			MethodName: "Offline",
			Handler:    _Signaling_Offline_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "signaling.proto",
}
