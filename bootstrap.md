# Boilerworks Go + HTMX -- Bootstrap

## Stack

- **Go 1.22+** with Chi router
- **Templ** for type-safe HTML templates (compile with `~/go/bin/templ generate`)
- **HTMX 2.0** for HTML-over-the-wire dynamic behavior
- **Tailwind CSS** via CDN (no build step)
- **PostgreSQL 16** via pgx/v5 connection pool
- **goose** format SQL migrations
- **Docker Compose** for local development

## Architecture

Handlers check `HX-Request` header to return full pages or HTML fragments. Templates use Templ's type-safe component system with a shared layout. Auth is session-based with bcrypt password hashing and SHA256 token storage.

## Conventions

1. All models use UUID primary keys (`github.com/google/uuid`)
2. Soft deletes via `deleted_at` column
3. Group-based permissions checked via middleware
4. CSRF tokens in cookies, validated via form fields or X-CSRF-Token header
5. Query functions written manually (same pattern as sqlc output)
6. Services contain business logic, handlers handle HTTP concerns
7. Dark theme: gray-950 background, gray-900 cards, indigo-600 accents

## Quick Start

```bash
cd docker && docker compose up -d --build
# API at http://localhost:8000
# Login: admin@boilerworks.dev / password
```

## Development

```bash
~/go/bin/templ generate   # compile templates
go build ./cmd/web        # build
go test -v -race ./...    # test
```

See the [Boilerworks Catalogue](../primers/CATALOGUE.md) for philosophy and universal patterns.
