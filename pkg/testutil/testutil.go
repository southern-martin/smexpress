package testutil

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// PostgresContainer holds a test PostgreSQL container.
type PostgresContainer struct {
	Container testcontainers.Container
	Pool      *pgxpool.Pool
	DSN       string
}

// NewPostgresContainer starts a PostgreSQL container for integration tests.
func NewPostgresContainer(t *testing.T, schema string) *PostgresContainer {
	t.Helper()
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16-alpine",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_DB":       "test",
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections").
			WithOccurrence(2).
			WithStartupTimeout(30 * time.Second),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("get mapped port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://test:test@%s:%s/test?search_path=%s&sslmode=disable", host, port.Port(), schema)

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}

	// Create the schema
	_, err = pool.Exec(ctx, fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema))
	if err != nil {
		t.Fatalf("create schema: %v", err)
	}

	t.Cleanup(func() {
		pool.Close()
		container.Terminate(ctx)
	})

	return &PostgresContainer{
		Container: container,
		Pool:      pool,
		DSN:       dsn,
	}
}
