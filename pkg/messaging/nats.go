package messaging

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Config holds NATS connection configuration.
type Config struct {
	URL       string
	Name      string // client name (service name)
	MaxReconn int
}

// Connection wraps NATS connection and JetStream context.
type Connection struct {
	NC     *nats.Conn
	JS     jetstream.JetStream
	logger *slog.Logger
}

// Connect establishes a NATS connection with JetStream.
func Connect(ctx context.Context, cfg Config, logger *slog.Logger) (*Connection, error) {
	opts := []nats.Option{
		nats.Name(cfg.Name),
		nats.MaxReconnects(cfg.MaxReconn),
		nats.ReconnectWait(2 * time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				logger.Warn("nats disconnected", slog.String("error", err.Error()))
			}
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			logger.Info("nats reconnected")
		}),
	}

	nc, err := nats.Connect(cfg.URL, opts...)
	if err != nil {
		return nil, fmt.Errorf("connect to nats: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("create jetstream context: %w", err)
	}

	return &Connection{NC: nc, JS: js, logger: logger}, nil
}

// Close closes the NATS connection.
func (c *Connection) Close() {
	if c.NC != nil {
		c.NC.Drain()
	}
}
