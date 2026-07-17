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

func setupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()
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
	t.Cleanup(func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("failed to terminate container: %v", err)
		}
	})

	connStr, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("failed to get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("failed to create pool: %v", err)
	}
	t.Cleanup(pool.Close)

	// retry ping — postgres image restarts once after initdb
	var pingErr error
	for i := 0; i < 10; i++ {
		if pingErr = pool.Ping(ctx); pingErr == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	if pingErr != nil {
		t.Fatalf("failed to ping db after retries: %v", pingErr)
	}

	_, err = pool.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS pgcrypto;

		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email TEXT NOT NULL UNIQUE,
			password_hash TEXT NOT NULL,
			created_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);

		CREATE TABLE IF NOT EXISTS bills (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name TEXT NOT NULL,
			category TEXT NOT NULL,
			service_provider TEXT NOT NULL,
			expected_amount NUMERIC(12,2) NOT NULL,
			due_day SMALLINT NOT NULL CHECK (due_day BETWEEN 1 AND 31),
			status TEXT NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
			created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
			updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
		);
	`)
	if err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return pool
}

func createTestUser(t *testing.T, pool *pgxpool.Pool) string {
	t.Helper()
	var userID string
	err := pool.QueryRow(context.Background(),
		`INSERT INTO users (email, password_hash) VALUES ($1, $2) RETURNING id`,
		"billtest@example.com", "hash",
	).Scan(&userID)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return userID
}

func TestBillRepository_CreateAndGet(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewBillRepository(pool)
	userID := createTestUser(t, pool)
	ctx := context.Background()

	bill, err := domain.NewBill(userID, "Electricity", "utilities", "ACME Power", 120.50, 15)
	if err != nil {
		t.Fatalf("failed to build bill: %v", err)
	}

	created, err := repo.Create(ctx, bill)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if created.ID == "" {
		t.Error("expected generated ID")
	}

	fetched, err := repo.GetByID(ctx, created.ID, userID)
	if err != nil {
		t.Fatalf("GetByID failed: %v", err)
	}
	if fetched == nil {
		t.Fatal("expected bill, got nil")
	}
	if fetched.Name != "Electricity" {
		t.Errorf("expected name Electricity, got %s", fetched.Name)
	}
}

func TestBillRepository_ListByUser(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewBillRepository(pool)
	userID := createTestUser(t, pool)
	ctx := context.Background()

	for _, name := range []string{"Electricity", "Water", "Gas"} {
		b, _ := domain.NewBill(userID, name, "utilities", "Provider", 50.00, 10)
		_, err := repo.Create(ctx, b)
		if err != nil {
			t.Fatalf("Create failed for %s: %v", name, err)
		}
	}

	bills, err := repo.ListByUser(ctx, userID)
	if err != nil {
		t.Fatalf("ListByUser failed: %v", err)
	}
	if len(bills) != 3 {
		t.Errorf("expected 3 bills, got %d", len(bills))
	}
}

func TestBillRepository_Update(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewBillRepository(pool)
	userID := createTestUser(t, pool)
	ctx := context.Background()

	bill, _ := domain.NewBill(userID, "Electricity", "utilities", "ACME Power", 120.50, 15)
	created, err := repo.Create(ctx, bill)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	created.Name = "Electricity Updated"
	created.ServiceProvider = "New Provider"
	created.ExpectedAmount = 150.00

	updated, err := repo.Update(ctx, created)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if updated.Name != "Electricity Updated" {
		t.Errorf("expected updated name, got %s", updated.Name)
	}
	if updated.ServiceProvider != "New Provider" {
		t.Errorf("expected updated provider, got %s", updated.ServiceProvider)
	}
}

func TestBillRepository_Delete(t *testing.T) {
	pool := setupTestDB(t)
	repo := postgres.NewBillRepository(pool)
	userID := createTestUser(t, pool)
	ctx := context.Background()

	bill, _ := domain.NewBill(userID, "Electricity", "utilities", "ACME Power", 120.50, 15)
	created, err := repo.Create(ctx, bill)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if err := repo.Delete(ctx, created.ID, userID); err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	fetched, err := repo.GetByID(ctx, created.ID, userID)
	if err != nil {
		t.Fatalf("GetByID after delete failed: %v", err)
	}
	if fetched != nil {
		t.Error("expected nil after delete, got a bill")
	}
}
