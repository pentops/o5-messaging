package outbox

import (
	"context"
	"database/sql"
	"errors"
	"net/url"
	"time"

	"github.com/elgris/sqrl"
	"github.com/pentops/j5/lib/j5codec"
	"github.com/pentops/o5-messaging/o5msg"
	"github.com/pentops/sqrlx.go/sqrlx"
)

type Config struct {
	TableName       string
	IDColumn        string
	HeadersColumn   string
	DataColumn      string
	SendAfterColumn string
}

var DefaultConfig = Config{
	TableName:       "outbox",
	IDColumn:        "id",
	HeadersColumn:   "headers",
	DataColumn:      "data",
	SendAfterColumn: "send_after",
}

var ErrMaxDelayExceeded = errors.New("maximum delay exceeded")

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
	return ss.SendDelayed(ctx, tx, 0, msg)
}

// SendDelayed places the message in the outbox table to be sent at approximately
// the current time plus the delay. The delay must be less than 15 minutes.
func (ss *Sender) SendDelayed(ctx context.Context, tx sqrlx.Transaction, approximateDelay time.Duration, msg o5msg.Message) error {
	if approximateDelay > 15*time.Minute {
		return ErrMaxDelayExceeded
	}

	wrapper, err := o5msg.WrapMessage(msg)
	if err != nil {
		return err
	}

	msgBytes, err := j5codec.Global.ProtoToJSON(wrapper.ProtoReflect())
	if err != nil {
		return err
	}

	headers := &url.Values{
		"Content-Type": []string{"application/json"},
	}

	sets := map[string]any{
		ss.IDColumn:      wrapper.MessageId,
		ss.HeadersColumn: headers.Encode(),
		ss.DataColumn:    msgBytes,
	}

	if approximateDelay > 0 {
		sets[ss.SendAfterColumn] = time.Now().Add(approximateDelay)
	}

	_, err = tx.Insert(ctx, sqrl.Insert(ss.TableName).SetMap(sets))

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
