// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: o5/messaging/v1/topic/dead.proto

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
	DeadMessageTopic_Dead_FullMethodName = "/o5.messaging.v1.topic.DeadMessageTopic/Dead"
)

// DeadMessageTopicClient is the client API for DeadMessageTopic service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeadMessageTopicClient interface {
	Dead(ctx context.Context, in *DeadMessage, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type deadMessageTopicClient struct {
	cc grpc.ClientConnInterface
}

func NewDeadMessageTopicClient(cc grpc.ClientConnInterface) DeadMessageTopicClient {
	return &deadMessageTopicClient{cc}
}

func (c *deadMessageTopicClient) Dead(ctx context.Context, in *DeadMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, DeadMessageTopic_Dead_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeadMessageTopicServer is the server API for DeadMessageTopic service.
// All implementations must embed UnimplementedDeadMessageTopicServer
// for forward compatibility
type DeadMessageTopicServer interface {
	Dead(context.Context, *DeadMessage) (*emptypb.Empty, error)
	mustEmbedUnimplementedDeadMessageTopicServer()
}

// UnimplementedDeadMessageTopicServer must be embedded to have forward compatible implementations.
type UnimplementedDeadMessageTopicServer struct {
}

func (UnimplementedDeadMessageTopicServer) Dead(context.Context, *DeadMessage) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Dead not implemented")
}
func (UnimplementedDeadMessageTopicServer) mustEmbedUnimplementedDeadMessageTopicServer() {}

// UnsafeDeadMessageTopicServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeadMessageTopicServer will
// result in compilation errors.
type UnsafeDeadMessageTopicServer interface {
	mustEmbedUnimplementedDeadMessageTopicServer()
}

func RegisterDeadMessageTopicServer(s grpc.ServiceRegistrar, srv DeadMessageTopicServer) {
	s.RegisterService(&DeadMessageTopic_ServiceDesc, srv)
}

func _DeadMessageTopic_Dead_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeadMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeadMessageTopicServer).Dead(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DeadMessageTopic_Dead_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeadMessageTopicServer).Dead(ctx, req.(*DeadMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// DeadMessageTopic_ServiceDesc is the grpc.ServiceDesc for DeadMessageTopic service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DeadMessageTopic_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "o5.messaging.v1.topic.DeadMessageTopic",
	HandlerType: (*DeadMessageTopicServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Dead",
			Handler:    _DeadMessageTopic_Dead_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "o5/messaging/v1/topic/dead.proto",
}
