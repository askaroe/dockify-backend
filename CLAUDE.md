# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the server
go run main.go

# Build
go build ./...

# Run tests
go test ./...

# Run a single test
go test ./internal/services/user/... -run TestFunctionName

# Regenerate Swagger docs (requires swag CLI)
swag init

# Lint (requires golangci-lint)
golangci-lint run
```

## Configuration

The server reads config from `config/config.json` by default. Override with the `config` env variable:

```bash
config=/path/to/config.json go run main.go
```

Config fields: `host`, `port`, PostgreSQL connection (`db_host`, `db_name`, `db_port`, `db_username`, `db_password`, `db_sslmode`), and `mindspore_model_url` for the ML service.

## Architecture

The project follows a layered architecture with strict dependency direction: **Handler → Service → Repository**.

```
cmd/main.go (entry point)
├── config/           — JSON config loading
├── internal/
│   ├── router/       — Gin route registration
│   ├── handlers/     — HTTP layer (request parsing, response writing)
│   ├── services/     — Business logic
│   ├── repository/   — Database queries (raw SQL via pgx)
│   ├── gateway/      — External HTTP calls (MindSpore ML service)
│   ├── scheduler/    — Background jobs (document cleanup every 3h)
│   ├── models/       — DB row structs (User, HealthMetrics, Location, Document)
│   └── entity/       — Request/response DTOs
└── pkg/
    ├── psql/         — pgxpool connection wrapper
    └── utils/        — Logger (logrus), HTTP client helpers
```

### Aggregator pattern

Each layer has an aggregator struct in its package root that bundles all domain-specific sub-services:
- `repository.Repository` embeds `health.Health`, `user.User`, `location.Location`, `document.Document`
- `services.Service` embeds the same domains plus `Gateway *gateway.Gateway`
- `handlers.Handler` embeds `user.User`, `health.Health`, `location.Location`, `document.Document`

When adding a new domain, follow this pattern: create the interface and implementation in a sub-package, then embed the interface in the aggregator struct.

### Gateway

`internal/gateway/mindspore/` wraps the external MindSpore ML service. It is called directly from `handlers.Handler.GetRecommendation` via `h.s.Gateway.MindSpore`. Errors from the ML service are non-fatal — the handler falls back to a default recommendation string.

### Document storage

Uploaded files are saved to a `documents/` directory at the working directory root. The scheduler (`internal/scheduler/scheduler.go`) deletes all files and truncates the `documents` DB table every 3 hours. Document IDs are UUIDs (custom generation to work around openGauss UUID limitations — see recent commits).

### Database

PostgreSQL via `pgxpool` (no ORM). Migrations live in `db/migrations/` as plain SQL files. The DB is targeted at openGauss compatibility (standard Postgres also works).

### API

All routes are under `/api/v1`. Swagger UI is at `/swagger/index.html`. Swagger annotations live in handler files; regenerate with `swag init` after changing annotations.
