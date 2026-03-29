# Claude -- Boilerworks Go + HTMX

Primary conventions doc: [`bootstrap.md`](bootstrap.md)

Read it before writing any code.

## Stack

- **Backend**: Go 1.22+ with Chi router
- **Frontend**: HTMX 2.0 + Templ (type-safe Go HTML templates)
- **Styling**: Tailwind CSS (CDN)
- **Database**: PostgreSQL 16 with pgx/v5
- **Migrations**: goose format SQL files
- **Auth**: Session-based (bcrypt + SHA256 token hashing)
- **Docker**: Compose stack (api + postgres + redis)

## Commands

```bash
# Development
make generate          # Compile .templ files
make build             # Build binary
make run               # Build and run
make test              # Run tests
make lint              # Run golangci-lint

# Docker
make docker-up         # Start stack (port 8000)
make docker-down       # Stop stack
make docker-reset      # Reset with fresh volumes

# Manual
~/go/bin/templ generate
go build ./cmd/web
go test -v -race ./...
```

## Architecture

- `cmd/web/main.go` -- entry point, server bootstrap
- `internal/config/` -- env-based configuration
- `internal/database/` -- pgx pool + query functions
- `internal/server/` -- Chi router, route registration
- `internal/middleware/` -- auth (session), CSRF, permission checks
- `internal/handler/` -- HTTP handlers (health, auth, CRUD, forms, workflows)
- `internal/model/` -- Go structs
- `internal/service/` -- business logic (auth, form validation, workflow state machine)
- `templates/` -- Templ files (layout, components, pages)
- `db/migrations/` -- goose-format SQL migrations

## Key Patterns

- HTMX: handlers return full pages (regular) or fragments (HX-Request header)
- Templ: `LayoutWithContent()` wraps page content in the base layout
- Permissions: group-based, checked via `middleware.RequirePermission()`
- CSRF: cookie-based, validated via form field or X-CSRF-Token header
- Soft deletes: `deleted_at` column on all content tables

## Default Credentials

- Email: `admin@boilerworks.dev`
- Password: `password`

## Ports

- API: 8000
- PostgreSQL: 5432
- Redis: 6379
