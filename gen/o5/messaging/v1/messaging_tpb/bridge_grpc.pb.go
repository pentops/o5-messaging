// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: o5/messaging/v1/topic/bridge.proto

package messaging_tpb

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	MessageBridgeTopic_Send_FullMethodName = "/o5.messaging.v1.topic.MessageBridgeTopic/Send"
)

// MessageBridgeTopicClient is the client API for MessageBridgeTopic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MessageBridgeTopicClient interface {
	Send(ctx context.Context, in *SendMessage, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type messageBridgeTopicClient struct {
	cc grpc.ClientConnInterface
}

func NewMessageBridgeTopicClient(cc grpc.ClientConnInterface) MessageBridgeTopicClient {
	return &messageBridgeTopicClient{cc}
}

func (c *messageBridgeTopicClient) Send(ctx context.Context, in *SendMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, MessageBridgeTopic_Send_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MessageBridgeTopicServer is the server API for MessageBridgeTopic service.
// All implementations must embed UnimplementedMessageBridgeTopicServer
// for forward compatibility
type MessageBridgeTopicServer interface {
	Send(context.Context, *SendMessage) (*emptypb.Empty, error)
	mustEmbedUnimplementedMessageBridgeTopicServer()
}

// UnimplementedMessageBridgeTopicServer must be embedded to have forward compatible implementations.
type UnimplementedMessageBridgeTopicServer struct {
}

func (UnimplementedMessageBridgeTopicServer) Send(context.Context, *SendMessage) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedMessageBridgeTopicServer) mustEmbedUnimplementedMessageBridgeTopicServer() {}

// UnsafeMessageBridgeTopicServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MessageBridgeTopicServer will
// result in compilation errors.
type UnsafeMessageBridgeTopicServer interface {
	mustEmbedUnimplementedMessageBridgeTopicServer()
}

func RegisterMessageBridgeTopicServer(s grpc.ServiceRegistrar, srv MessageBridgeTopicServer) {
	s.RegisterService(&MessageBridgeTopic_ServiceDesc, srv)
}

func _MessageBridgeTopic_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MessageBridgeTopicServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: MessageBridgeTopic_Send_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MessageBridgeTopicServer).Send(ctx, req.(*SendMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// MessageBridgeTopic_ServiceDesc is the grpc.ServiceDesc for MessageBridgeTopic service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MessageBridgeTopic_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "o5.messaging.v1.topic.MessageBridgeTopic",
	HandlerType: (*MessageBridgeTopicServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _MessageBridgeTopic_Send_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "o5/messaging/v1/topic/bridge.proto",
}
