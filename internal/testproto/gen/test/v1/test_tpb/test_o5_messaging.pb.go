// Code generated by protoc-gen-go-o5-messaging. DO NOT EDIT.
// versions:
// - protoc-gen-go-o5-messaging 0.0.0
// source: test/v1/topic/test.proto

package test_tpb

import (
	context "context"
	messaging_j5pb "github.com/pentops/j5/gen/j5/messaging/v1/messaging_j5pb"
	messaging_pb "github.com/pentops/o5-messaging/gen/o5/messaging/v1/messaging_pb"
	o5msg "github.com/pentops/o5-messaging/o5msg"
)

// Service: TestTopic
// Method: Test

func (msg *TestMessage) O5MessageHeader() o5msg.Header {
	header := o5msg.Header{
		GrpcService: "test.v1.topic.TestTopic",
		GrpcMethod:  "Test",
		Headers:     map[string]string{},
	}
	return header
}

type TestTopicTxSender[C any] struct {
	sender o5msg.TxSender[C]
}

func NewTestTopicTxSender[C any](sender o5msg.TxSender[C]) *TestTopicTxSender[C] {
	sender.Register(o5msg.TopicDescriptor{
		Service: "test.v1.topic.TestTopic",
		Methods: []o5msg.MethodDescriptor{
			{
				Name:    "Test",
				Message: (*TestMessage).ProtoReflect(nil).Descriptor(),
			},
		},
	})
	return &TestTopicTxSender[C]{sender: sender}
}

type TestTopicCollector[C any] struct {
	collector o5msg.Collector[C]
}

func NewTestTopicCollector[C any](collector o5msg.Collector[C]) *TestTopicCollector[C] {
	collector.Register(o5msg.TopicDescriptor{
		Service: "test.v1.topic.TestTopic",
		Methods: []o5msg.MethodDescriptor{
			{
				Name:    "Test",
				Message: (*TestMessage).ProtoReflect(nil).Descriptor(),
			},
		},
	})
	return &TestTopicCollector[C]{collector: collector}
}

type TestTopicPublisher struct {
	publisher o5msg.Publisher
}

func NewTestTopicPublisher(publisher o5msg.Publisher) *TestTopicPublisher {
	publisher.Register(o5msg.TopicDescriptor{
		Service: "test.v1.topic.TestTopic",
		Methods: []o5msg.MethodDescriptor{
			{
				Name:    "Test",
				Message: (*TestMessage).ProtoReflect(nil).Descriptor(),
			},
		},
	})
	return &TestTopicPublisher{publisher: publisher}
}

// Method: Test

func (send TestTopicTxSender[C]) Test(ctx context.Context, sendContext C, msg *TestMessage) error {
	return send.sender.Send(ctx, sendContext, msg)
}

func (collect TestTopicCollector[C]) Test(sendContext C, msg *TestMessage) {
	collect.collector.Collect(sendContext, msg)
}

func (publish TestTopicPublisher) Test(ctx context.Context, msg *TestMessage) error {
	return publish.publisher.Publish(ctx, msg)
}

// Service: GreetingRequestTopic
// Expose Request Metadata
func (msg *GreetingMessage) SetJ5RequestMetadata(md *messaging_j5pb.RequestMetadata) {
	msg.Request = md
}
func (msg *GreetingMessage) GetJ5RequestMetadata() *messaging_j5pb.RequestMetadata {
	return msg.Request
}

// Method: Greeting

func (msg *GreetingMessage) O5MessageHeader() o5msg.Header {
	header := o5msg.Header{
		GrpcService: "test.v1.topic.GreetingRequestTopic",
		GrpcMethod:  "Greeting",
		Headers:     map[string]string{},
	}
	if msg.Request != nil {
		header.Extension = &messaging_pb.Message_Request_{
			Request: &messaging_pb.Message_Request{
				ReplyTo: msg.Request.ReplyTo,
			},
		}
	} else {
		header.Extension = &messaging_pb.Message_Request_{
			Request: &messaging_pb.Message_Request{
				ReplyTo: "",
			},
		}
	}
	return header
}

type GreetingRequestTopicTxSender[C any] struct {
	sender o5msg.TxSender[C]
}

func NewGreetingRequestTopicTxSender[C any](sender o5msg.TxSender[C]) *GreetingRequestTopicTxSender[C] {
	sender.Register(o5msg.TopicDescriptor{
		Service: "test.v1.topic.GreetingRequestTopic",
		Methods: []o5msg.MethodDescriptor{
			{
				Name:    "Greeting",
				Message: (*GreetingMessage).ProtoReflect(nil).Descriptor(),
			},
		},
	})
	return &GreetingRequestTopicTxSender[C]{sender: sender}
}

type GreetingRequestTopicCollector[C any] struct {
	collector o5msg.Collector[C]
}

func NewGreetingRequestTopicCollector[C any](collector o5msg.Collector[C]) *GreetingRequestTopicCollector[C] {
	collector.Register(o5msg.TopicDescriptor{
		Service: "test.v1.topic.GreetingRequestTopic",
		Methods: []o5msg.MethodDescriptor{
			{
				Name:    "Greeting",
				Message: (*GreetingMessage).ProtoReflect(nil).Descriptor(),
			},
		},
	})
	return &GreetingRequestTopicCollector[C]{collector: collector}
}

type GreetingRequestTopicPublisher struct {
	publisher o5msg.Publisher
}

func NewGreetingRequestTopicPublisher(publisher o5msg.Publisher) *GreetingRequestTopicPublisher {
	publisher.Register(o5msg.TopicDescriptor{
		Service: "test.v1.topic.GreetingRequestTopic",
		Methods: []o5msg.MethodDescriptor{
			{
				Name:    "Greeting",
				Message: (*GreetingMessage).ProtoReflect(nil).Descriptor(),
			},
		},
	})
	return &GreetingRequestTopicPublisher{publisher: publisher}
}

// Method: Greeting

func (send GreetingRequestTopicTxSender[C]) Greeting(ctx context.Context, sendContext C, msg *GreetingMessage) error {
	return send.sender.Send(ctx, sendContext, msg)
}

func (collect GreetingRequestTopicCollector[C]) Greeting(sendContext C, msg *GreetingMessage) {
	collect.collector.Collect(sendContext, msg)
}

func (publish GreetingRequestTopicPublisher) Greeting(ctx context.Context, msg *GreetingMessage) error {
	return publish.publisher.Publish(ctx, msg)
}

// Service: GreetingResponseTopic
// Expose Request Metadata
func (msg *ResponseMessage) SetJ5RequestMetadata(md *messaging_j5pb.RequestMetadata) {
	msg.Request = md
}
func (msg *ResponseMessage) GetJ5RequestMetadata() *messaging_j5pb.RequestMetadata {
	return msg.Request
}

// Method: Response

func (msg *ResponseMessage) O5MessageHeader() o5msg.Header {
	header := o5msg.Header{
		GrpcService: "test.v1.topic.GreetingResponseTopic",
		GrpcMethod:  "Response",
		Headers:     map[string]string{},
	}
	if msg.Request != nil {
		header.Extension = &messaging_pb.Message_Reply_{
			Reply: &messaging_pb.Message_Reply{
				ReplyTo: msg.Request.ReplyTo,
			},
		}
	}
	return header
}

type GreetingResponseTopicTxSender[C any] struct {
	sender o5msg.TxSender[C]
}

func NewGreetingResponseTopicTxSender[C any](sender o5msg.TxSender[C]) *GreetingResponseTopicTxSender[C] {
	sender.Register(o5msg.TopicDescriptor{
		Service: "test.v1.topic.GreetingResponseTopic",
		Methods: []o5msg.MethodDescriptor{
			{
				Name:    "Response",
				Message: (*ResponseMessage).ProtoReflect(nil).Descriptor(),
			},
		},
	})
	return &GreetingResponseTopicTxSender[C]{sender: sender}
}

type GreetingResponseTopicCollector[C any] struct {
	collector o5msg.Collector[C]
}

func NewGreetingResponseTopicCollector[C any](collector o5msg.Collector[C]) *GreetingResponseTopicCollector[C] {
	collector.Register(o5msg.TopicDescriptor{
		Service: "test.v1.topic.GreetingResponseTopic",
		Methods: []o5msg.MethodDescriptor{
			{
				Name:    "Response",
				Message: (*ResponseMessage).ProtoReflect(nil).Descriptor(),
			},
		},
	})
	return &GreetingResponseTopicCollector[C]{collector: collector}
}

type GreetingResponseTopicPublisher struct {
	publisher o5msg.Publisher
}

func NewGreetingResponseTopicPublisher(publisher o5msg.Publisher) *GreetingResponseTopicPublisher {
	publisher.Register(o5msg.TopicDescriptor{
		Service: "test.v1.topic.GreetingResponseTopic",
		Methods: []o5msg.MethodDescriptor{
			{
				Name:    "Response",
				Message: (*ResponseMessage).ProtoReflect(nil).Descriptor(),
			},
		},
	})
	return &GreetingResponseTopicPublisher{publisher: publisher}
}

// Method: Response

func (send GreetingResponseTopicTxSender[C]) Response(ctx context.Context, sendContext C, msg *ResponseMessage) error {
	return send.sender.Send(ctx, sendContext, msg)
}

func (collect GreetingResponseTopicCollector[C]) Response(sendContext C, msg *ResponseMessage) {
	collect.collector.Collect(sendContext, msg)
}

func (publish GreetingResponseTopicPublisher) Response(ctx context.Context, msg *ResponseMessage) error {
	return publish.publisher.Publish(ctx, msg)
}
