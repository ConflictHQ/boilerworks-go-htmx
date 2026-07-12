# Calliope — Boilerworks Go + HTMX
<!-- Agent shim for https://github.com/calliopeai/calliope-cli -->

Primary conventions doc: [`bootstrap.md`](bootstrap.md)

Read it before writing any code.

---

## Project-specific notes

- Go 1.25+ with the Chi router; HTMX 2.0 + Templ (type-safe Go HTML templates) + Tailwind (CDN); PostgreSQL 16 via pgx/v5; goose-format SQL migrations; Docker Compose stack (api + postgres + redis).
- Layout: `cmd/web/main.go` (entry), `internal/{config,database,server,middleware,handler,model,service}`, `templates/` (Templ), `db/migrations/`.
- HTMX: handlers return full pages (regular) or fragments (`HX-Request` header); `LayoutWithContent()` wraps page content in the base layout.
- Permissions group-based via `middleware.RequirePermission()`; CSRF cookie-based (form field or `X-CSRF-Token`); soft deletes via `deleted_at`, session auth (bcrypt + SHA256 token hashing).
- Run `make generate` (compile `.templ`) before `make build`/`run`; `make test` runs `go test -race`, `make lint` runs golangci-lint.
- Ports: API 8000, Postgres 5432, Redis 6379.
