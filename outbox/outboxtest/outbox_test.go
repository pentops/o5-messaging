package outboxtest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/pentops/o5-messaging/internal/testproto/gen/test/v1/test_tpb"
	"github.com/pentops/o5-messaging/outbox"
	"github.com/pentops/pgtest.go/pgtest"
	"github.com/pentops/sqrlx.go/sqrlx"
)

func TestOutboxSend(t *testing.T) {
	ctx := context.Background()

	conn := pgtest.GetTestDB(t)
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `
	CREATE TABLE outbox (
		id uuid PRIMARY KEY,
		headers text,
		data jsonb
	)`); err != nil {
		t.Fatal(err)
	}

	db, err := sqrlx.New(conn, sqrlx.Dollar)
	if err != nil {
		t.Fatal(err)
	}

	testSender := test_tpb.NewTestTopicTxSender(outbox.DefaultSender)
	outboxAsserter := NewOutboxAsserter(t, conn)

	if err := db.Transact(context.Background(), &sqrlx.TxOptions{
		Isolation: sql.LevelReadCommitted,
		ReadOnly:  false,
		Retryable: true,
	}, func(ctx context.Context, tx sqrlx.Transaction) error {
		return testSender.Test(ctx, tx, &test_tpb.TestMessage{
			Message: "Hello, World!",
		})
	}); err != nil {
		t.Fatal(err)
	}

	msg := &test_tpb.TestMessage{}
	outboxAsserter.PopMessage(t, msg)
	if msg.Message != "Hello, World!" {
		t.Fatalf("expected message to be 'Hello, World!', got %q", msg.Message)
	}
}
