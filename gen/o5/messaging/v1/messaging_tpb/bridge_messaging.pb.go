// Code generated by protoc-gen-go-messaging. DO NOT EDIT.

package messaging_tpb

// Service: MessageBridgeTopic
// Method: Send

func (msg *SendMessage) MessagingHeaders() map[string]string {
	headers := map[string]string{
		"grpc-service": "/o5.messaging.v1.topic.MessageBridgeTopic/Send",
		"grpc-message": "o5.messaging.v1.topic.SendMessage",
	}
	return headers
}
