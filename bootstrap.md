# Boilerworks Go + Next.js -- Bootstrap

> Go backend with Chi router + Next.js 16 frontend. REST API with session-based auth,
> group-based permissions, Products/Categories CRUD, Forms engine, Workflow engine.

See the [Boilerworks Catalogue](../primers/CATALOGUE.md) for philosophy and universal patterns.

See the [stack primer](../primers/go-nextjs/PRIMER.md) for stack-specific conventions and build order.

## Quick Start

```bash
# Docker (full stack)
cd docker && docker compose up -d --build

# Local development
# Terminal 1: API
make run

# Terminal 2: Frontend
make frontend-dev
```

## Architecture

```
Go API (Chi router, port 8000)
  |-- REST JSON API (/api/*)
  |-- Session auth (httpOnly cookies)
  |-- PostgreSQL 16 (pgx/v5)
  |-- Redis 7 (cache)
  |
  +-- Next.js 16 (port 3000)
        |-- App Router + TypeScript
        |-- Tailwind CSS (dark theme)
        |-- Sonner toasts
        |-- REST API client
```

## Conventions

- UUID primary keys on every table
- Soft deletes via `deleted_at` column
- `MutationResult` pattern for all mutations: `{ok: true}` or `{ok: false, errors: [...]}`
- Auth check at the top of every handler
- Permission checks via middleware
- Go tests use `testing` + `httptest`
- Frontend tests use Vitest + Testing Library
