package outbox

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/elgris/sqrl"
	"github.com/pentops/o5-messaging/o5msg"
	"github.com/pentops/sqrlx.go/sqrlx"
	"google.golang.org/protobuf/encoding/protojson"
)

type Config struct {
	TableName     string
	IDColumn      string
	HeadersColumn string
	DataColumn    string
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

type Sender struct {
	Config
	o5msg.TopicSet
}

func NewSender(config Config) *Sender {
	return &Sender{
		Config:   config,
		TopicSet: o5msg.TopicSet{},
	}
}

// Send places the message in the outbox table.
func (ss *Sender) Send(ctx context.Context, tx sqrlx.Transaction, msg o5msg.Message) error {
	wrapper, err := o5msg.WrapMessage(msg)
	if err != nil {
		return err
	}

	msgBytes, err := protojson.Marshal(wrapper)
	if err != nil {
		return err
	}

	headers := &url.Values{
		"Content-Type": []string{"application/json"},
	}

	_, err = tx.Insert(ctx, sqrl.Insert(ss.TableName).
		Columns(ss.IDColumn, ss.HeadersColumn, ss.DataColumn).
		Values(wrapper.MessageId, headers.Encode(), msgBytes))

	return err
}

// DEPRECATED: Send bypasses registration and is included for an easier
// transition to typed senders / collectors
func Send(ctx context.Context, tx sqrlx.Transaction, msg o5msg.Message) error {
	return DefaultSender.Send(ctx, tx, msg)
}

type DirectPublisher struct {
	sender *Sender
	db     sqrlx.Transactor
}

func NewDirectPublisher(db sqrlx.Transactor, sender *Sender) (*DirectPublisher, error) {
	return &DirectPublisher{
		sender: sender,
		db:     db,
	}, nil
}

func (dp *DirectPublisher) Register(md o5msg.TopicDescriptor) {
	dp.sender.Register(md)
}

func (dp *DirectPublisher) Publish(ctx context.Context, msg o5msg.Message) error {
	return dp.db.Transact(ctx, &sqrlx.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
		Retryable: true,
	}, func(ctx context.Context, tx sqrlx.Transaction) error {
		return dp.sender.Send(ctx, tx, msg)
	})
}
