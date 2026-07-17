# AGENTS.md

Guidelines for AI agents (Claude Code, Copilot, etc.) working on this repository.

---

## Project Overview

BillTracker backend is a Go REST API following **Hexagonal Architecture** (Ports & Adapters). Before making any changes, understand the layer boundaries described below — violating them is the most common mistake.

---

## Architecture Rules

```
domain/     ← NO external imports allowed. Pure Go structs and validation only.
ports/      ← Interfaces only. No implementations here.
usecase/    ← Imports domain + ports. NO imports from adapter/, gin, pgx, etc.
adapter/    ← The only layer allowed to import external libraries (Gin, pgx, sqlc, etc.)
cmd/api/    ← Wiring only. Instantiates and connects all layers.
```

**If you find yourself importing Gin inside `usecase/`, or importing `pgx` inside `domain/`, stop — that's an architecture violation.**

---

## Commands

### Run the server
```bash
go run ./cmd/api
```

### Run all unit tests
```bash
go test ./...
```

### Run integration tests (requires Docker running)
```bash
go test -tags=integration ./internal/adapter/postgres/... -v
```

### Build
```bash
go build ./...
```

### Run migrations
```bash
migrate -path migrations \
  -database "postgres://billtracker:billtracker@localhost:5432/billtracker?sslmode=disable" \
  up
```

### Create a new migration
```bash
migrate create -ext sql -dir migrations -seq <migration_name>
# Always fill in both the .up.sql and .down.sql files
```

### Regenerate sqlc code after changing queries or schema
```bash
sqlc generate
```

---

## TDD Workflow

This project follows strict TDD. Always follow this order:

1. Write a failing test first (unit test in `usecase/`, using a fake repo from `ports/`).
2. Run it and confirm it fails with the expected error.
3. Write the minimum implementation to make it pass.
4. Run tests again and confirm green.
5. Refactor if needed, keeping tests green.

**Never write implementation code without a failing test first.**

---

## Testing Conventions

| Test type | Location | Tag | Dependencies |
|-----------|----------|-----|-------------|
| Unit | `usecase/`, `adapter/http/` | none | Fake repos only, no I/O |
| Integration | `adapter/postgres/` | `//go:build integration` | Docker + testcontainers-go |

- Use `fakeUserRepo` and `fakeBillRepo` (defined in `usecase/*_test.go`) for unit tests.
- Integration tests use `setupTestDB(t)` helper which spins up an ephemeral Postgres container per test via `testcontainers-go`. Always use `t.Cleanup` (not `defer`) for teardown inside helpers.
- Never share state between tests — each integration test gets its own container.

---

## Code Conventions

- All domain structs must have `json` tags using `snake_case`.
- `PasswordHash` and any sensitive fields must have `json:"-"` to prevent accidental exposure.
- Repository methods that query by both `id` and `userID` enforce ownership at the query level — never skip the `userID` filter.
- Return the same error (`ErrInvalidCredentials`) for both wrong password and unknown email — never leak which emails are registered.
- Money values use `NUMERIC(12,2)` in Postgres and `float64` in Go. When scanning into `pgtype.Numeric`, always use `fmt.Sprintf("%.2f", amount)` — do not scan a raw `float64`.
- UUID fields in `pgtype.UUID` convert to string via `uuid.UUID(row.ID.Bytes).String()` (using `github.com/google/uuid`).

---

## Environment

Local infrastructure runs via Docker Compose:
- PostgreSQL on `localhost:5432`
- MinIO on `localhost:9000` (console at `localhost:9001`)

Start with:
```bash
docker compose up -d
```

Default local credentials are in `docker-compose.yml`. Production secrets must come from environment variables — never hardcode them.

---

## What NOT to Do

- Do not add business logic to Gin handlers (`adapter/http/`). Handlers translate HTTP ↔ use case calls only.
- Do not call `pgx` or `sqlc` from `usecase/` or `domain/`.
- Do not store binary files (images, PDFs) in Postgres — use the storage adapter (MinIO/S3).
- Do not commit `.env` files or secrets.
- Do not skip writing the `.down.sql` migration — every migration must be reversible.
- Do not use `defer` for cleanup inside `t.Helper()` functions — use `t.Cleanup()`.
