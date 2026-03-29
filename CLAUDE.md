# Claude -- Boilerworks Go + Next.js

Primary conventions doc: [`bootstrap.md`](bootstrap.md)

Read it before writing any code.

## Stack

- **Backend**: Go 1.25+ with Chi router
- **Frontend**: Next.js 16 (App Router, React 19, TypeScript)
- **API**: REST (JSON over HTTP)
- **Database**: PostgreSQL 16 with pgx/v5
- **Cache**: Redis 7
- **Auth**: Session-based (bcrypt + SHA256 token hashing, httpOnly cookies)
- **Docker**: Compose stack (api + ui + postgres + redis)

## Commands

```bash
# Backend
make build             # Build Go binary
make run               # Build and run
make test              # Run Go tests
make lint              # Run golangci-lint

# Frontend
make frontend-dev      # Start Next.js dev server
make frontend-build    # Build Next.js for production
make frontend-test     # Run Vitest tests
make frontend-lint     # Run ESLint

# Docker
make docker-up         # Start full stack
make docker-down       # Stop stack
make docker-reset      # Reset with fresh volumes
```

## Architecture

- `cmd/api/main.go` -- entry point, server bootstrap
- `internal/config/` -- env-based configuration
- `internal/database/` -- pgx pool + query functions
- `internal/server/` -- Chi router, route registration, CORS
- `internal/middleware/` -- auth (session), permission checks
- `internal/handler/` -- REST API handlers (health, auth, CRUD, forms, workflows)
- `internal/service/` -- business logic (auth, form validation, workflow state machine)
- `internal/model/` -- Go structs (UUID PKs, soft deletes)
- `db/migrations/` -- goose-format SQL migrations
- `frontend/` -- Next.js 16 app (dark admin theme, Tailwind CSS)

## Key Patterns

- REST API: all endpoints under `/api/` return JSON
- MutationResult pattern: `{ok: true}` or `{ok: false, errors: [...]}`
- Permissions: group-based, checked via `middleware.RequirePermission()`
- Soft deletes: `deleted_at` column on all content tables
- UUID primary keys everywhere
- CORS configured for Next.js frontend origin
- Session cookie shared between Go API and Next.js via `credentials: "include"`

## API Routes

- `POST /api/auth/login` -- login
- `POST /api/auth/register` -- register
- `POST /api/auth/logout` -- logout
- `GET /api/auth/me` -- current user + permissions
- `GET /api/dashboard` -- dashboard stats
- `GET/POST /api/items` -- list/create
- `GET/PUT/DELETE /api/items/:uuid` -- get/update/delete
- `GET/POST /api/categories` -- list/create
- `GET/PUT/DELETE /api/categories/:uuid` -- get/update/delete
- `GET/POST /api/forms` -- list/create definitions
- `GET/PUT/DELETE /api/forms/:uuid` -- get/update/delete definition
- `GET/POST /api/forms/:uuid/submissions` -- list/create submissions
- `GET/POST /api/workflows` -- list/create definitions
- `GET/PUT/DELETE /api/workflows/:uuid` -- get/update/delete definition
- `GET/POST /api/workflows/:uuid/instances` -- list/create instances
- `GET /api/workflows/instances/:uuid` -- get instance detail
- `POST /api/workflows/instances/:uuid/transition` -- transition instance

## Default Credentials

- Email: `admin@boilerworks.dev`
- Password: `password`

## Ports

- Go API: 8000
- Next.js: 3000
- PostgreSQL: 5432
- Redis: 6379
