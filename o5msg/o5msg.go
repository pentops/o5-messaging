package o5msg

import (
	"context"

	"github.com/pentops/o5-go/messaging/v1/messaging_pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Sender[C any] interface {

	// Send is called with a partially populated message and a payload.
	// The sender is responsible for marshalling the payload into the message
	// body, setting MessageID and may add headers or other metadata to the message.
	// The generated code will populate GrpcService, GrpcMethod, and
	// where applicable, DestinationTopic and ReplyTo fields.
	Send(ctx context.Context, sendContext C, msg *messaging_pb.Message, payload proto.Message) error

	// Register is called when a new Topic is created using ths Sender, allowing
	// the sender to handle reflection, pre-loading of validation etc.
	// Implementation should store the descriptor only, and process later if
	// context and errors are required.
	Register(TopicDescriptor)
}

type TopicDescriptor struct {
	Service string
	Methods []MethodDescriptor
}

type MethodDescriptor struct {
	Name    string
	Message protoreflect.MessageDescriptor
}

// TopicSet can be anonymously embedded to senders for a simple implementation of Register
type TopicSet []TopicDescriptor

// Register adds a new TopicDescriptor to the TopicSet
func (ts *TopicSet) Register(td TopicDescriptor) {
	*ts = append(*ts, td)
}
