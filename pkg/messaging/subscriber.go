package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/nats-io/nats.go/jetstream"
)

// Handler processes a received message.
type Handler func(ctx context.Context, data []byte) error

// Subscriber consumes messages from JetStream.
type Subscriber struct {
	conn   *Connection
	logger *slog.Logger
}

// NewSubscriber creates a new event subscriber.
func NewSubscriber(conn *Connection, logger *slog.Logger) *Subscriber {
	return &Subscriber{conn: conn, logger: logger}
}

// Subscribe starts consuming messages for a given stream and subject.
func (s *Subscriber) Subscribe(ctx context.Context, streamName, consumerName, subject string, handler Handler) error {
	stream, err := s.conn.JS.Stream(ctx, streamName)
	if err != nil {
		return fmt.Errorf("get stream %s: %w", streamName, err)
	}

	consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
		Durable:       consumerName,
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return fmt.Errorf("create consumer %s: %w", consumerName, err)
	}

	_, err = consumer.Consume(func(msg jetstream.Msg) {
		if err := handler(ctx, msg.Data()); err != nil {
			s.logger.Error("message handler failed",
				slog.String("subject", subject),
				slog.String("error", err.Error()),
			)
			msg.Nak()
			return
		}
		msg.Ack()
	})
	if err != nil {
		return fmt.Errorf("consume %s: %w", subject, err)
	}

	s.logger.Info("subscribed", slog.String("stream", streamName), slog.String("subject", subject))
	return nil
}

// DecodeEvent is a helper to decode JSON event data into a struct.
func DecodeEvent[T any](data []byte) (T, error) {
	var event T
	if err := json.Unmarshal(data, &event); err != nil {
		return event, fmt.Errorf("decode event: %w", err)
	}
	return event, nil
}
