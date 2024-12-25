package o5msg

// Run buf generate prior to running these tests to ensure the code matches the
// plugin.

import (
	"context"
	"fmt"
	"testing"

	"github.com/pentops/j5/gen/j5/messaging/v1/messaging_j5pb"
	"github.com/pentops/o5-messaging/gen/o5/messaging/v1/messaging_pb"
	"github.com/pentops/o5-messaging/internal/testproto/gen/test/v1/test_tpb"
	"github.com/pentops/o5-messaging/o5msg"
	"github.com/stretchr/testify/assert"
)

func Send(tx *transaction, msg *messaging_pb.Message) error {
	return tx.send(msg)
}

type TestSender struct {
	o5msg.TopicSet
}

func (s *TestSender) Send(ctx context.Context, tx *transaction, msg o5msg.Message) error {
	wrapper, err := o5msg.WrapMessage(msg)
	if err != nil {
		return err
	}
	return Send(tx, wrapper)
}

type TestContext struct {
	Name string
}

type transaction struct {
	id     int
	runner *runner
	t      testing.TB
}

type gotMessage struct {
	txIdx   int
	message *messaging_pb.Message
}

func (tx *transaction) send(msg *messaging_pb.Message) error {
	tx.runner.got = append(tx.runner.got, gotMessage{
		txIdx:   tx.id,
		message: msg,
	})
	tx.t.Logf("send %d %s %v", tx.id, msg.MessageId, msg.Body)
	return nil
}

type runner struct {
	last int
	got  []gotMessage
}

func (r *runner) Transact(t testing.TB, cb func(testing.TB, *transaction) error) error {
	tx := &transaction{
		id:     r.last,
		runner: r,
		t:      t,
	}
	r.last++
	return cb(t, tx)
}

func TestCallback(t *testing.T) {
	ts := &TestSender{}
	ctx := context.Background()
	testTopic := test_tpb.NewTestTopicTxSender(ts)
	greetingRequestTopic := test_tpb.NewGreetingRequestTopicTxSender(ts)
	greetingResponseTopic := test_tpb.NewGreetingResponseTopicTxSender(ts)

	t.Run("Topics should be registered", func(t *testing.T) {
		allTopics := ts.TopicSet
		assert.Len(t, allTopics, 3)
		byFull := make(map[string]o5msg.MethodDescriptor)
		for _, td := range allTopics {
			for _, m := range td.Methods {
				byFull[fmt.Sprintf("/%s/%s", td.Service, m.Name)] = m
			}
		}

		assert.Contains(t, byFull, "/test.v1.topic.TestTopic/Test")
		assert.Contains(t, byFull, "/test.v1.topic.GreetingRequestTopic/Greeting")
		assert.Contains(t, byFull, "/test.v1.topic.GreetingResponseTopic/Response")

		testMethod := byFull["/test.v1.topic.TestTopic/Test"]
		testDest := testMethod.Message
		assert.Equal(t, "test.v1.topic.TestMessage", string(testDest.FullName()))
	})

	db := &runner{}

	requestMeta := &messaging_j5pb.RequestMetadata{
		ReplyTo: "reply-to",
		Context: []byte("request-metadata"),
	}

	err := db.Transact(t, func(t testing.TB, tx *transaction) error {
		if err := testTopic.Test(ctx, tx, &test_tpb.TestMessage{}); err != nil {
			return err
		}
		if err := greetingRequestTopic.Greeting(ctx, tx, &test_tpb.GreetingMessage{
			Request: requestMeta,
		}); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	err = db.Transact(t, func(t testing.TB, tx *transaction) error {
		if err := testTopic.Test(ctx, tx, &test_tpb.TestMessage{}); err != nil {
			return err
		}
		if err := greetingResponseTopic.Response(ctx, tx, &test_tpb.ResponseMessage{
			Request: requestMeta,
		}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Four messages should be sent", func(t *testing.T) {
		if len(db.got) != 4 {
			t.Errorf("expected 1 got %d", len(db.got))
		}
		msgA, msgB, msgC, msgD := db.got[0], db.got[1], db.got[2], db.got[3]

		assert.Equal(t, 0, msgA.txIdx)
		assert.Equal(t, "test.v1.topic.TestTopic", msgA.message.GrpcService)
		assert.Equal(t, "type.googleapis.com/test.v1.topic.TestMessage", msgA.message.Body.TypeUrl)

		assert.Equal(t, 0, msgB.txIdx)
		assert.Equal(t, "test.v1.topic.GreetingRequestTopic", msgB.message.GrpcService)
		assert.Equal(t, "reply-to", msgB.message.Extension.(*messaging_pb.Message_Request_).Request.ReplyTo)

		assert.Equal(t, 1, msgC.txIdx)
		assert.Equal(t, "test.v1.topic.TestTopic", msgC.message.GrpcService)

		assert.Equal(t, 1, msgD.txIdx)
		assert.Equal(t, "test.v1.topic.GreetingResponseTopic", msgD.message.GrpcService)
		assert.Equal(t, "reply-to", msgD.message.Extension.(*messaging_pb.Message_Reply_).Reply.ReplyTo)
	})

}
