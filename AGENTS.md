# Agents -- Boilerworks Go + HTMX

Primary conventions doc: [`bootstrap.md`](bootstrap.md)

Read it before writing any code.

## Key Commands

- `~/go/bin/templ generate` -- compile .templ files to Go
- `go build ./cmd/web` -- build the server binary
- `go test -v -race ./...` -- run all tests
- `cd docker && docker compose up -d --build` -- start the stack

## Architecture Notes

- Handlers in `internal/handler/` return full pages or HTMX fragments based on `HX-Request` header
- Templates in `templates/` use Templ's type-safe component system
- Auth middleware in `internal/middleware/auth.go` validates sessions and injects user + permissions into context
- Permission checks via `middleware.RequirePermission("resource.action")`
- CSRF middleware in `internal/middleware/csrf.go` protects all mutating requests

## Ports

API: 8084, PostgreSQL: 5441, Redis: 6384
