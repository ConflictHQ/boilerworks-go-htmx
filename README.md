# Boilerworks Go + HTMX

> Full-stack Go application with HTMX, Templ, and Tailwind CSS. Session auth, permissions, CRUD, forms engine, and workflow engine with the Boilerworks dark theme.

## Quick Start

```bash
cd docker
docker compose up -d --build
```

Open [http://localhost:8000](http://localhost:8000) and sign in:

- **Email:** `admin@boilerworks.dev`
- **Password:** `password`

## Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.22+ with Chi router |
| Templates | Templ (type-safe Go HTML) |
| Frontend | HTMX 2.0 + Tailwind CSS |
| Database | PostgreSQL 16 (pgx/v5) |
| Migrations | goose format |
| Auth | Session-based (bcrypt + SHA256) |

## Features

- **Session Authentication** -- Register, login, logout with bcrypt password hashing and httpOnly session cookies
- **Group-based Permissions** -- Admin, editor, viewer roles with granular permission checks on every route
- **Products CRUD** -- Full create/read/update/delete with categories, status, and pricing
- **Categories CRUD** -- Organize products into categories
- **Forms Engine** -- Define forms with JSON schema, render dynamically, collect and validate submissions
- **Workflow Engine** -- Define state machines with states and transitions, create instances, track transition history
- **HTMX Integration** -- Full-page loads for standard requests, HTML fragment swaps for HTMX requests
- **CSRF Protection** -- Cookie-based tokens validated on all mutating requests
- **Dark Theme** -- Boilerworks dark theme (gray-950 background, gray-900 cards, indigo accents)

## Development

```bash
# Install templ
go install github.com/a-h/templ/cmd/templ@latest

# Generate templates
~/go/bin/templ generate

# Build and run
go build -o bin/web ./cmd/web
./bin/web

# Run tests
go test -v -race ./...
```

## Project Structure

```
cmd/web/main.go              -- entry point
internal/
  config/                    -- env-based configuration
  database/                  -- pgx pool connection
  database/queries/          -- typed query functions
  server/                    -- Chi router + route registration
  middleware/                -- auth, CSRF, permission middleware
  handler/                   -- HTTP handlers
  model/                     -- Go structs
  service/                   -- business logic
templates/                   -- Templ files
  layout.templ               -- base layout with sidebar
  components/                -- flash, pagination
  pages/                     -- auth, dashboard, products, categories, forms, workflows
db/migrations/               -- goose SQL migrations
docker/                      -- Docker Compose stack
```

## Ports

| Service | Port |
|---------|------|
| API | 8084 |
| PostgreSQL | 5441 |
| Redis | 6384 |

## Want to help build this?

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

---

Boilerworks is a [Conflict](https://weareconflict.com) brand. CONFLICT is a registered trademark of Conflict LLC.
