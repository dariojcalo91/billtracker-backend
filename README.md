# BillTracker — Backend

> Never miss a bill again. BillTracker is a personal finance tool for managing monthly bills, tracking payment status, and staying ahead of due dates — with proof-of-payment uploads and a clear monthly spending overview.

This repository contains the **Go REST API** that powers the BillTracker ecosystem. The mobile and web clients (Flutter) live in [`billtracker-mobile`](https://github.com/dariojcalo91/billtracker-mobile) _(coming soon)_.

---

## Tech Stack

| Layer | Choice | Why |
|---|---|---|
| Language | Go 1.26 | High performance, strong concurrency, familiar to the team |
| HTTP framework | Gin | Minimal, fast, production-proven |
| Database | PostgreSQL 16 | Relational data, strong consistency, ACID guarantees |
| DB access | sqlc + pgx/v5 | Type-safe SQL — no ORM magic, full control over queries |
| Migrations | golang-migrate | Version-controlled, reversible schema changes |
| Auth | JWT (golang-jwt) + bcrypt | Self-implemented to enforce security best practices |
| File storage | MinIO (dev) / AWS S3 or Cloudflare R2 (prod) | Binary files stay out of Postgres |
| Notifications | flutter_local_notifications (client-side) | Due dates are known at creation time — no server push needed |
| Testing | testify + testcontainers-go | Unit tests (fake repos) + integration tests (real ephemeral Postgres) |
| Containers | Docker + Docker Compose | Reproducible local environment |

## Architecture

This project follows **Hexagonal Architecture** (Ports & Adapters):

```
cmd/api/              ← composition root (wires everything together)
internal/
  domain/             ← pure business entities and validation rules
  ports/              ← interfaces (contracts) the domain depends on
  usecase/            ← application logic (orchestrates domain + ports)
  adapter/
    http/             ← Gin handlers (inbound adapter)
    postgres/         ← PostgreSQL repository implementations (outbound adapter)
    storage/          ← S3/MinIO file storage (outbound adapter)
migrations/           ← SQL migration files (golang-migrate)
```

The domain and use case layers have **zero knowledge** of Gin, Postgres, or any infrastructure. This means:
- Business logic is unit-testable with no database or HTTP server running.
- Infrastructure can be swapped (e.g. Postgres → another DB) without touching business logic.
- Each layer has a single, clear responsibility.

---

## Getting Started

### Prerequisites

- [Go 1.26+](https://go.dev/dl/)
- [Docker + Docker Compose](https://docs.docker.com/get-docker/)
- [golang-migrate CLI](https://github.com/golang-migrate/migrate)
- [sqlc](https://sqlc.dev/)

### 1. Clone the repo

```bash
git clone https://github.com/dariojcalo91/billtracker-backend.git
cd billtracker-backend
```

### 2. Start local infrastructure

```bash
docker compose up -d
```

This starts:
- PostgreSQL on `localhost:5432`
- MinIO on `localhost:9000` (console at `localhost:9001`)

### 3. Run migrations

```bash
migrate -path migrations \
  -database "postgres://billtracker:billtracker@localhost:5432/billtracker?sslmode=disable" \
  up
```

### 4. Set environment variables

```bash
cp .env.example .env
# edit .env with your values
```

### 5. Run the server

```bash
go run ./cmd/api
```

Server starts on `http://localhost:8080`.

---

## API Endpoints

### Auth
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/auth/register` | ❌ | Register a new user |
| POST | `/auth/login` | ❌ | Login, returns JWT |

### Bills
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| POST | `/bills` | ✅ | Create a bill |
| GET | `/bills` | ✅ | List all bills for the authenticated user |
| GET | `/bills/:id` | ✅ | Get a single bill |
| PUT | `/bills/:id` | ✅ | Update a bill |
| DELETE | `/bills/:id` | ✅ | Delete a bill |

### Dashboard _(coming soon)_
| Method | Path | Auth | Description |
|--------|------|------|-------------|
| GET | `/dashboard` | ✅ | Monthly summary: done / upcoming / overdue |

All authenticated endpoints require:
```
Authorization: Bearer <token>
```

---

## Running Tests

### Unit tests (no Docker needed, fast)
```bash
go test ./...
```

### Integration tests (requires Docker)
```bash
go test -tags=integration ./internal/adapter/postgres/... -v
```

Integration tests use `testcontainers-go` to spin up an ephemeral Postgres instance per test, fully isolated and cleaned up automatically.

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `DATABASE_URL` | `postgres://billtracker:billtracker@localhost:5432/billtracker?sslmode=disable` | Postgres connection string |
| `JWT_SECRET` | `dev-secret-change-me` | Secret key for signing JWTs — **change this in production** |
| `PORT` | `8080` | HTTP server port |

> ⚠️ Never commit real secrets. Always use `.env` locally and a secrets manager in production.

---

## Project Status

| Phase | Status | Description |
|-------|--------|-------------|
| Phase 1 — MVP | 🟡 In progress | Auth + Bill CRUD + Dashboard |
| Phase 2 — Notifications & Uploads | ⬜ Planned | Local push notifications + proof-of-payment file uploads |
| Phase 3 — Web & Spending Insights | ⬜ Planned | Flutter web build + monthly spending dashboard |

See [`bill-tracker-project-plan.md`](./bill-tracker-project-plan.md) for the full roadmap, architecture decisions, and changelog.

---

## Contributing

This is a personal learning project. If you're reading this and want to suggest improvements, feel free to open an issue.

---

## License

MIT
