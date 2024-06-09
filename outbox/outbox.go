package outbox

import (
	"context"
	"fmt"
	"net/url"

	"github.com/elgris/sqrl"
	"github.com/google/uuid"
	"github.com/pentops/o5-go/messaging/v1/messaging_pb"
	"github.com/pentops/o5-messaging.go/o5msg"
	"github.com/pentops/sqrlx.go/sqrlx"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type Config struct {
	TableName     string
	IDColumn      string
	HeadersColumn string
	DataColumn    string
}

type Sender struct {
	Config
	o5msg.TopicSet
}

var DefaultConfig = Config{
	TableName:     "outbox",
	IDColumn:      "id",
	HeadersColumn: "headers",
	DataColumn:    "data",
}

var DefaultSender *Sender = &Sender{
	Config:   DefaultConfig,
	TopicSet: o5msg.TopicSet{},
}

func Send(ctx context.Context, tx sqrlx.Transaction, msg *messaging_pb.Message, payload proto.Message) error {
	return DefaultSender.Send(ctx, tx, msg, payload)
}

func (ss *Sender) Send(ctx context.Context, tx sqrlx.Transaction, msg *messaging_pb.Message, payload proto.Message) error {
	bodyBytes, err := protojson.Marshal(payload)
	if err != nil {
		return err
	}

	msg.Body = &messaging_pb.Any{
		TypeUrl:  fmt.Sprintf("type.googleapis.com/%s", payload.ProtoReflect().Descriptor().FullName()),
		Value:    bodyBytes,
		Encoding: messaging_pb.WireEncoding_PROTOJSON,
	}

	msgBytes, err := protojson.Marshal(msg)
	if err != nil {
		return err
	}

	headers := &url.Values{
		"Content-Type": []string{"application/json"},
	}

	msg.MessageId = uuid.New().String()

	_, err = tx.Insert(ctx, sqrl.Insert(ss.TableName).
		Columns(ss.IDColumn, ss.HeadersColumn, ss.DataColumn).
		Values(msg.MessageId, headers.Encode(), msgBytes))

	return err
}
