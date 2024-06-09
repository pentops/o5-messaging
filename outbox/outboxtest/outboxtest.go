package outboxtest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/elgris/sqrl"
	"github.com/pentops/o5-go/messaging/v1/messaging_pb"
	"github.com/pentops/o5-messaging.go/outbox"
	"github.com/pentops/sqrlx.go/sqrlx"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type OutboxAsserter struct {
	db *sqrlx.Wrapper

	outbox.Config
}

type TB interface {
	Fatal(args ...any)
	Fatalf(format string, args ...any)
	Helper()
}

func NewOutboxAsserter(t TB, conn sqrlx.Connection) *OutboxAsserter {
	db, err := sqrlx.New(conn, sq.Dollar)
	if err != nil {
		t.Fatal(err.Error())
	}
	return &OutboxAsserter{
		db:     db,
		Config: outbox.DefaultConfig,
	}
}

func (oa *OutboxAsserter) PopMessage(ctx context.Context, tb TB, msg proto.Message) {
	tb.Helper()
	typeURL := fmt.Sprintf("type.googleapis.com/%s", msg.ProtoReflect().Descriptor().FullName())
	wrapper, err := oa.popWrapper(ctx, tb, sq.Eq{
		fmt.Sprintf("%s->'body'->>'typeUrl'", oa.DataColumn): typeURL,
	})
	if errors.Is(err, sql.ErrNoRows) {
		tb.Fatalf("no message found for type %s", typeURL)
	} else if err != nil {
		tb.Fatal(err)
	}
	if err := protojson.Unmarshal(wrapper.Body.Value, msg); err != nil {
		tb.Fatal(err)
	}
}

func (oa *OutboxAsserter) popWrapper(ctx context.Context, tb TB, condition sq.Eq) (*messaging_pb.Message, error) {

	var msgBody []byte
	if err := oa.db.Transact(ctx, nil, func(ctx context.Context, tx sqrlx.Transaction) error {
		tb.Helper()

		var msgID string

		if err := tx.SelectRow(ctx, sq.Select(oa.IDColumn, oa.DataColumn).
			From(oa.TableName).
			Where(condition).
			OrderBy(oa.IDColumn).
			Limit(1)).Scan(&msgID, &msgBody); err != nil {
			return err
		}

		if _, err := tx.Delete(ctx, sq.Delete(oa.TableName).
			Where(sq.Eq{oa.IDColumn: msgID})); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	wrapper := &messaging_pb.Message{}
	if err := protojson.Unmarshal(msgBody, wrapper); err != nil {
		tb.Fatal(err)
	}

	return wrapper, nil
}
