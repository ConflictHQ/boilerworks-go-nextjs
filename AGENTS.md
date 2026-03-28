# Agents -- Boilerworks Go + Next.js

Primary conventions doc: [`bootstrap.md`](bootstrap.md)

Read it before writing any code.

## Key Files

- `cmd/api/main.go` -- Go API entry point
- `internal/server/routes.go` -- all API route registrations
- `internal/handler/` -- REST handlers
- `internal/middleware/auth.go` -- session auth + permission middleware
- `frontend/app/` -- Next.js pages
- `frontend/lib/api.ts` -- REST API client
- `docker/docker-compose.yml` -- full stack
- `.github/workflows/ci.yml` -- CI pipeline
