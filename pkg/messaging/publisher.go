package messaging

import (
	"context"
	"encoding/json"
	"fmt"

)

// Event represents a domain event to publish.
type Event struct {
	Subject string
	Data    any
}

// Publisher publishes events to NATS JetStream.
type Publisher struct {
	conn *Connection
}

// NewPublisher creates a new event publisher.
func NewPublisher(conn *Connection) *Publisher {
	return &Publisher{conn: conn}
}

// Publish publishes an event to JetStream.
func (p *Publisher) Publish(ctx context.Context, event Event) error {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	_, err = p.conn.JS.Publish(ctx, event.Subject, data)
	if err != nil {
		return fmt.Errorf("publish event %s: %w", event.Subject, err)
	}

	p.conn.logger.Info("event published",
		"subject", event.Subject,
	)
	return nil
}
