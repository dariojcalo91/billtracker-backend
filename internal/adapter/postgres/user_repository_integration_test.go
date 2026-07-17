//go:build integration

package postgres_test

import (
	"context"
	"testing"
	"time"

	"github.com/dariojcalo91/billtracker/internal/adapter/postgres"
	"github.com/dariojcalo91/billtracker/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestUserRepository_CreateAndGetByEmail(t *testing.T) {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx,
		"postgres:16",
		tcpostgres.WithDatabase("billtracker"),
		tcpostgres.WithUsername("billtracker"),
		tcpostgres.WithPassword("billtracker"),
	)
	if err != nil {
		t.Fatalf("failed to start container: %v", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	}()

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to connect: %v", err)
	}
	defer pool.Close()

	// Retry the first ping — the postgres image restarts once after
	// initdb, so the port can briefly refuse connections even after
	// testcontainers reports the container as ready.
	for i := 0; i < 10; i++ {
		if err = pool.Ping(ctx); err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		t.Fatalf("failed to ping db after retries: %v", err)
	}

	// Run our migration manually since this is a fresh container
	_, err = pool.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS pgcrypto;
		CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	repo := postgres.NewUserRepository(pool)

	user, err := domain.NewUser("integration@example.com", "some-hash")
	if err != nil {
		t.Fatalf("failed to build domain user: %v", err)
	}

	created, err := repo.Create(ctx, user)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.ID == "" {
		t.Error("expected generated ID, got empty string")
	}

	fetched, err := repo.GetByEmail(ctx, "integration@example.com")
	if err != nil {
		t.Fatalf("GetByEmail failed: %v", err)
	}
	if fetched == nil {
		t.Fatal("expected to find user, got nil")
	}
	if fetched.Email != "integration@example.com" {
		t.Errorf("expected email integration@example.com, got %s", fetched.Email)
	}
}
