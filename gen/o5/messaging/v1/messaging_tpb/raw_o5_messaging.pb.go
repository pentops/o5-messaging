// Code generated by protoc-gen-go-o5-messaging. DO NOT EDIT.
// versions:
// - protoc-gen-go-o5-messaging 0.0.0
// source: o5/messaging/v1/topic/raw.proto

package messaging_tpb

import (
	context "context"
	o5msg "github.com/pentops/o5-messaging/o5msg"
)

// Service: RawMessageTopic
type RawMessageTopicTxSender[C any] struct {
	sender o5msg.TxSender[C]
}

func NewRawMessageTopicTxSender[C any](sender o5msg.TxSender[C]) *RawMessageTopicTxSender[C] {
	sender.Register(o5msg.TopicDescriptor{
		Service: "o5.messaging.v1.topic.RawMessageTopic",
		Methods: []o5msg.MethodDescriptor{
			{
				Name:    "Raw",
				Message: (*RawMessage).ProtoReflect(nil).Descriptor(),
			},
		},
	})
	return &RawMessageTopicTxSender[C]{sender: sender}
}

type RawMessageTopicCollector[C any] struct {
	collector o5msg.Collector[C]
}

func NewRawMessageTopicCollector[C any](collector o5msg.Collector[C]) *RawMessageTopicCollector[C] {
	collector.Register(o5msg.TopicDescriptor{
		Service: "o5.messaging.v1.topic.RawMessageTopic",
		Methods: []o5msg.MethodDescriptor{
			{
				Name:    "Raw",
				Message: (*RawMessage).ProtoReflect(nil).Descriptor(),
			},
		},
	})
	return &RawMessageTopicCollector[C]{collector: collector}
}

type RawMessageTopicPublisher struct {
	publisher o5msg.Publisher
}

func NewRawMessageTopicPublisher(publisher o5msg.Publisher) *RawMessageTopicPublisher {
	publisher.Register(o5msg.TopicDescriptor{
		Service: "o5.messaging.v1.topic.RawMessageTopic",
		Methods: []o5msg.MethodDescriptor{
			{
				Name:    "Raw",
				Message: (*RawMessage).ProtoReflect(nil).Descriptor(),
			},
		},
	})
	return &RawMessageTopicPublisher{publisher: publisher}
}

// Method: Raw

func (msg *RawMessage) O5MessageHeader() o5msg.Header {
	header := o5msg.Header{
		GrpcService: "o5.messaging.v1.topic.RawMessageTopic",
		GrpcMethod:  "Raw",
		Headers:     map[string]string{},
	}
	return header
}

func (send RawMessageTopicTxSender[C]) Raw(ctx context.Context, sendContext C, msg *RawMessage) error {
	return send.sender.Send(ctx, sendContext, msg)
}

func (collect RawMessageTopicCollector[C]) Raw(sendContext C, msg *RawMessage) {
	collect.collector.Collect(sendContext, msg)
}

func (publish RawMessageTopicPublisher) Raw(ctx context.Context, msg *RawMessage) {
	publish.publisher.Publish(ctx, msg)
}