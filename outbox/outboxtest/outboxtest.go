package outboxtest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	sq "github.com/elgris/sqrl"
	"github.com/pentops/o5-messaging/gen/o5/messaging/v1/messaging_pb"
	"github.com/pentops/o5-messaging/outbox"
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

var txOptions = &sqrlx.TxOptions{
	Retryable: true,
	ReadOnly:  false,
	Isolation: sql.LevelReadCommitted,
}

func getContext(t TB) context.Context {
	t.Helper()

	ctx, ok := t.(interface{ Context() context.Context })
	if ok {
		return ctx.Context()
	}

	return context.Background()
}

func MessageBodyMatches[T proto.Message](filter func(T) bool) condition {
	return func(conditions *queryConditions) {
		conditions.checkers = append(conditions.checkers, func(wrapper *messaging_pb.Message) (bool, error) {
			body := (*new(T)).ProtoReflect().New().Interface().(T)

			err := protojson.Unmarshal(wrapper.Body.Value, body)
			if err != nil {
				return false, fmt.Errorf("error unmarshalling body: %w", err)
			}

			return filter(body), nil
		})
	}
}

type condition func(*queryConditions)

type queryConditions struct {
	filters  []filter
	checkers []checker
}

func (oa *OutboxAsserter) PopMessage(tb TB, msg proto.Message, conditions ...condition) {
	tb.Helper()
	typeURL := fmt.Sprintf("type.googleapis.com/%s", msg.ProtoReflect().Descriptor().FullName())

	qc := &queryConditions{}
	for _, condition := range conditions {
		condition(qc)
	}

	messageTypeCondition(typeURL)(qc)

	wrapper, err := oa.popWrapper(getContext(tb), tb, *qc)
	if errors.Is(err, sql.ErrNoRows) {
		tb.Fatalf("no message found for type %s", typeURL)
	} else if errors.Is(err, multiMatchError) {
		tb.Fatalf("found multiple messages for type %s", typeURL)
	} else if err != nil {
		tb.Fatal(err)
	}

	err = protojson.Unmarshal(wrapper.Body.Value, msg)
	if err != nil {
		tb.Fatal(err)
	}
}

func (oa *OutboxAsserter) AssertEmpty(tb TB) {
	tb.Helper()

	q := sq.Select("COUNT(*)").
		From(oa.TableName)

	var count int
	err := oa.db.Transact(getContext(tb), txOptions, func(ctx context.Context, tx sqrlx.Transaction) error {
		err := tx.SelectRow(ctx, q).Scan(&count)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		tb.Fatal(err)
	}

	if count > 0 {
		tb.Fatalf("expected outbox to be empty, but found %d messages", count)
	}
}

func (oa *OutboxAsserter) PurgeAll(tb TB) {
	tb.Helper()

	err := oa.db.Transact(getContext(tb), txOptions, func(ctx context.Context, tx sqrlx.Transaction) error {
		_, err := tx.Exec(ctx, sq.Delete(oa.TableName))
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		tb.Fatal(err)
	}
}

func (oa *OutboxAsserter) ForEachMessage(tb TB, fn func(*messaging_pb.Message)) {
	tb.Helper()

	q := sq.Select(oa.DataColumn).
		From(oa.TableName)

	err := oa.db.Transact(getContext(tb), txOptions, func(ctx context.Context, tx sqrlx.Transaction) error {
		rows, err := tx.Query(ctx, q)
		if err != nil {
			return err
		}
		defer rows.Close()

		for rows.Next() {
			var msgBody []byte
			if err := rows.Scan(&msgBody); err != nil {
				return err
			}

			wrapper := &messaging_pb.Message{}
			err := protojson.Unmarshal(msgBody, wrapper)
			if err != nil {
				return err
			}

			fn(wrapper)
		}

		err = rows.Err()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		tb.Fatal(err)
	}
}

type checker func(*messaging_pb.Message) (bool, error)

type filter func(*sq.SelectBuilder, outbox.Config)

func messageTypeCondition(typeURL string) condition {
	return func(conditions *queryConditions) {
		conditions.filters = append(conditions.filters, messageTypeFilter(typeURL))
	}
}

func messageTypeFilter(typeURL string) filter {
	return func(sb *sq.SelectBuilder, cfg outbox.Config) {
		sb.Where(sq.Eq{
			fmt.Sprintf("%s->'body'->>'typeUrl'", cfg.DataColumn): typeURL,
		})
	}
}

var multiMatchError = errors.New("multiple messages matched")

func (oa *OutboxAsserter) popWrapper(ctx context.Context, tb TB, conditions queryConditions) (*messaging_pb.Message, error) {
	query := sq.Select(oa.DataColumn).
		From(oa.TableName).
		OrderBy(oa.IDColumn)

	for _, filter := range conditions.filters {
		filter(query, oa.Config)
	}

	bodies := make([][]byte, 0, 1)

	err := oa.db.Transact(ctx, txOptions, func(ctx context.Context, tx sqrlx.Transaction) error {
		tb.Helper()

		var msgBody []byte

		rows, err := tx.Select(ctx, query)
		if err != nil {
			return err
		}

		defer rows.Close()

		for rows.Next() {
			err := rows.Scan(&msgBody)
			if err != nil {
				return err
			}

			bodies = append(bodies, msgBody)
		}

		err = rows.Err()
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	matchedMessages := make([]*messaging_pb.Message, 0, 1)

bodies:
	for _, body := range bodies {
		wrapper := &messaging_pb.Message{}
		if err := protojson.Unmarshal(body, wrapper); err != nil {
			return nil, err
		}

		for _, check := range conditions.checkers {
			matches, err := check(wrapper)
			if err != nil {
				return nil, err
			} else if !matches {
				continue bodies
			}
		}

		matchedMessages = append(matchedMessages, wrapper)
	}

	if len(matchedMessages) > 1 {
		return nil, multiMatchError
	}

	if len(matchedMessages) < 1 {
		return nil, sql.ErrNoRows
	}

	msgID := matchedMessages[0].MessageId

	q := sq.Delete(oa.TableName).
		Where(sq.Eq{oa.IDColumn: msgID})

	err = oa.db.Transact(ctx, txOptions, func(ctx context.Context, tx sqrlx.Transaction) error {
		tb.Helper()

		_, err := tx.Delete(ctx, q)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return matchedMessages[0], nil
}
