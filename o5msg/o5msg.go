package o5msg

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/pentops/j5/lib/j5codec"
	"github.com/pentops/o5-messaging/gen/o5/messaging/v1/messaging_pb"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Registry allows the sender/collector to handle reflection, pre-loading of validation etc.
// Implementation should store the descriptor only, and process later if
// context and errors are required.
type Registry interface {
	// Register is called when a new Topic is created using ths Sender
	Register(TopicDescriptor)
}

// TxSender immediately marshalls and places the message into the next layer, e.g.
// a database transaction or a network call. Generated code requires this
// interface to build a NewFooSender
type TxSender[T any] interface {
	Registry

	// Send is called with an un-marshalled message and routing information
	// populated by the generated code.
	// The sender is responsible for building and sending an o5.messaging.v1.Message
	Send(ctx context.Context, sendContext T, msg Message) error
}

// Collector is the counterpart to Sender, and is used to queue messages for
// later marshalling and sending. This allows the caller to push error handling
// to a later call but is otherwise identical in purpose and functionality
type Collector[T any] interface {
	Registry

	// Collect queues the message for later marshalling and sending, deferring
	// marshalling until the end of an operation, e.g. a psm state transition.
	Collect(sendContext T, msg Message)
}

// Publisher sends messages directly, no transactions or queues. Useful for
// workers which aren't using an outbox pattern
type Publisher interface {
	Registry

	// Publish sends the message directly, without any transaction context
	Publish(ctx context.Context, msg Message) error
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

// Message is implemented by all method bodies in the generated code. The
// sender/collector should use the O5Header method to get the routing
// information, and 'blindly' marshal the message body (e.g. using protojson)
type Message interface {
	O5MessageHeader() Header
	proto.Message
}

// Header is produced by generated code, responsibility for these fields is on
// the generated code. The sender/collector is responsible for mapping this to
// the fields of the wrapper message.
type Header struct {
	GrpcService      string
	GrpcMethod       string
	Headers          map[string]string
	DestinationTopic string
	Extension        messaging_pb.IsMessage_Extension
}

// WrapMessage returns the *messaging_pb.Message for the given Message, using
// j5-json encoding. It does not set SourceApp and SourceEnv, as they should
// be set by the infrastructure layer when placing the message into a
// broker/bus.
func WrapMessage(msg Message) (*messaging_pb.Message, error) {
	bodyData, err := j5codec.Global.ProtoToJSON(msg.ProtoReflect())
	if err != nil {
		return nil, err
	}

	wrapper := MessageWrapper(msg)

	wrapper.Body = &messaging_pb.Any{
		TypeUrl:  fmt.Sprintf("type.googleapis.com/%s", msg.ProtoReflect().Descriptor().FullName()),
		Value:    bodyData,
		Encoding: messaging_pb.WireEncoding_J5_JSON,
	}

	return wrapper, nil
}

// MessageWrapper returns the *messaging_pb.Message for the given Message, but
// does not set the Body field, allowing for other encodings. It does not set
// SourceApp and SourceEnv, as they should be set by the infrastructure layer
// when placing the message into a broker/bus.
func MessageWrapper(msg Message) *messaging_pb.Message {
	header := msg.O5MessageHeader()
	wrapper := &messaging_pb.Message{
		MessageId:        uuid.New().String(),
		Timestamp:        timestamppb.Now(),
		GrpcService:      header.GrpcService,
		GrpcMethod:       header.GrpcMethod,
		DestinationTopic: header.DestinationTopic,
		Headers:          header.Headers,
		Extension:        header.Extension,
	}
	return wrapper
}
