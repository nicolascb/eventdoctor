# GEMINI.md

This file provides guidance to Gemini/Antigravity when working with code in this repository.

## Commands

### Backend (Go)

All commands run from `backend/` or use the `-C ./backend` flag. `CGO_ENABLED=1` is required for the API (SQLite/mattn/go-sqlite3).

```bash
# Build
CGO_ENABLED=1 go build -C ./backend -o eventdoctor ./cmd/api/main.go
CGO_ENABLED=0 go build -C ./backend -o eventdoctor-cli ./cmd/cli/main.go

# Run API (mock data)
cd backend && WITH_MOCK=1 go run ./cmd/api/main.go

# Tests (all)
CGO_ENABLED=1 go test -C ./backend -v -race ./...

# Single test / package
go test -C ./backend -run TestOrderedMap ./internal/service/...
```

### Frontend (Node/React)

```bash
cd frontend
npm ci          # install
npm run dev     # Vite dev server → http://localhost:5173
npm run build   # tsc -b + Vite bundle (TypeScript errors fail the build)
npm run lint    # ESLint
```

### Environment Variables (API)

| Var           | Default      | Description              |
|---------------|--------------|--------------------------|
| `PORT`        | `:8087`      | HTTP listen address      |
| `SQLITE_PATH` | `./data.db`  | SQLite file path         |
| `WITH_MOCK`   | —            | Set to `1` for demo data |

## Architecture

Monorepo: `backend/` (Go module `github.com/nicolascb/eventdoctor`) + `frontend/` (React 19 + Vite).

### Binaries

| Binary | Entry point | CGO | Description |
|--------|-------------|-----|-------------|
| `eventdoctor` | `cmd/api/main.go` | Required (`CGO_ENABLED=1`) | HTTP API server |
| `eventdoctor-cli` | `cmd/cli/main.go` | Not needed (`CGO_ENABLED=0`) | CLI client |
| `eventdoctor-auditor` | `cmd/auditor/main.go` | — | Auditor process |

### Request flow

```
HTTP → internal/api/api.go (routing, CORS)
      → api.Service interface (internal/api/service.go)
      → service.Service (internal/service/)
      → db.* functions (internal/db/sqlite.go, repository.go)
      → SQLite
```

### API endpoints

| Method | Path | Handler | Description |
|--------|------|---------|-------------|
| POST | `/v1/config` | `handlerApplyConfig` | Upload YAML spec |
| GET | `/v1/producers` | `handlerListProducers` | List all producers |
| GET | `/v1/events` | `handlerListEvents` | List events by topic |
| GET | `/v1/consumers` | `handlerListConsumers` | List all consumers |
| GET | `/v1/overview` | `handlerOverview` | Dashboard overview |
| GET | `/v1/topics/{name}` | `handlerGetTopic` | Topic detail view |
| GET | `/v1/services/{name}` | `handlerGetService` | Service detail view |

### Key abstractions

- **`api.Service` interface** (`internal/api/service.go`) — the only contract between HTTP and business logic. Every new endpoint method must be added here first.
- **`service.NewService(db)`** — the single implementation, injected at startup in `cmd/api/main.go`.
- **`db.SQLExecutor` interface** (`internal/db/transaction.go`) — abstracts `*sql.DB` and `*sql.Tx` so DB functions work in both contexts. All DB query/write functions should accept `SQLExecutor`, not `*sql.DB`.
- **`db.WithTransaction(ctx, db, fn)`** — wraps multi-step writes in a transaction; the callback receives a `SQLExecutor`.
- **`orderedMap[K, V]` + `findOrAppend()`** (`internal/service/aggregate.go`) — preserve insertion order when building aggregated responses. Reuse these for any new aggregation.
- **Data Fetching Patterns** — When retrieving related entities (e.g., Producers for a list of Events), **ALWAYS use batch queries** (e.g., `WHERE foreign_id IN (...)`) or joins. **NEVER** use N+1 queries inside loops or fetch entire tables into memory.
- **`eventdoctor.EventDoctorSpec`** (`internal/eventdoctor/spec.go`) — core domain model for YAML spec parsing and validation. Validation rules in `spec.go`, helper functions in `utils.go`.
- **`internal/logger/`** — structured logging via `slog` + `tint` formatter. Use this instead of `fmt.Println` or `log.Printf`.

### Adding a new endpoint (checklist)

1. Add method to `api.Service` interface (`internal/api/service.go`). **Ensure list endpoints accept pagination parameters (`page`, `pageSize`).**
2. Implement in `internal/service/` (new file if non-trivial). Implement efficient data fetching.
3. Add DB query function(s) to `internal/db/sqlite.go` (reads) or `internal/db/repository.go` (deletes/writes); add row types to `internal/db/models/models.go`. Use `SQLExecutor` as the db parameter type.
4. Add response type to `internal/api/response/` (base types in `response.go`: `ErrorResponse`, `SuccessResponse`)
5. Register route in `internal/api/api.go`
6. Add client method to `internal/client/`

### CLI

Uses `urfave/cli/v3`. Command logic lives in `internal/commands/` (testable); wiring lives in `cmd/cli/`.

- `cmd/cli/config.go` → `config validate / apply`
- `cmd/cli/get.go` → `get topic / get service`
- Table output via `github.com/jedib0t/go-pretty/v6/table` with `newTable()` helper in `internal/commands/get.go`.
- Server URL resolved from `--url` flag or from `eventdoctor.yml` servers list + `--env`.

### Frontend

- **API client:** `src/lib/api.ts` — all API calls. Base URL from `import.meta.env.VITE_API_URL` (baked at build time).
- **Types:** `src/types/index.ts` — TypeScript interfaces mirroring backend response types.
- **Hooks:** `src/hooks/` — one hook per resource (`useProducers`, `useConsumers`, `useEvents`, `useOverview`), built on a generic `useAsync` hook.
- **Components:** `src/components/` — page-level views (`OverviewView`, `ProducersView`, `ConsumersView`, `TopicsView`, `ConfigValidator`) + shared UI (`LoadingState`, `ErrorState`, `EmptyState`, `PageHeader`, `SearchInput`).
- **UI primitives:** `src/components/ui/` — Radix UI components (button, card, table, dialog, etc.) styled with TailwindCSS 4.
- **Stack:** React 19, Vite 7, TypeScript 5.9, TailwindCSS 4, Radix UI, Lucide icons, React Hook Form + Zod.
- **UI Patterns:** Maintain visual and functional consistency between similar entity views (e.g., Producers vs Consumers). Use shared wrapper components (like `ServiceDetailsDialog`) to enforce consistent layouts and reduce duplication.


## Database

Schema in `backend/migration.sql`. Tables: `topics`, `services`, `events`, `event_headers`, `producers`, `consumers`, `missing_topics`, `missing_events`. Applied automatically at startup via `db.NewSQLiteDB()`.

## Release

Two independent workflows (prefixed tags):

| Component | Tag format        | Trigger                                                  |
|-----------|-------------------|----------------------------------------------------------|
| Backend   | `backend/v1.2.0`  | `.github/workflows/release-backend.yml` → GoReleaser     |
| Frontend  | `frontend/v0.3.0` | `.github/workflows/release-frontend.yml` → Docker push   |

GoReleaser config: `backend/.goreleaser.yml`. Docker images published to `ghcr.io/nicolascb/`.

CI runs on path-filtered pushes/PRs (`backend/**` or `frontend/**`).
