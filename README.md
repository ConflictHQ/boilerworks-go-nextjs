# Boilerworks Go + Next.js

> API-first Go backend with Next.js 16 frontend for performance-sensitive applications.

**Status:** Complete

Go with Chi router as a REST API backend paired with Next.js 16 for teams that want Go's raw performance, small memory footprint, and straightforward concurrency model behind a modern React frontend.

## Quick Start

```bash
# Docker (full stack)
cd docker && docker compose up -d --build
# API: http://localhost:8000
# Frontend: http://localhost:3000

# Local development
make run              # Start Go API on :8000
make frontend-dev     # Start Next.js on :3000
```

**Default credentials:** `admin@boilerworks.dev` / `password`

## Stack

| Layer | Technology |
|-------|-----------|
| Backend | Go 1.25, Chi router |
| Frontend | Next.js 16, React 19, TypeScript |
| API | REST (JSON) |
| Database | PostgreSQL 16 (pgx/v5) |
| Cache | Redis 7 |
| Auth | Session-based (bcrypt + SHA256) |
| Styling | Tailwind CSS 4 (dark theme) |
| Testing | Go testing + Vitest |

## Features

- Session-based authentication with httpOnly cookies
- Group-based permissions (admin, editor, viewer)
- Items + Categories CRUD with soft deletes
- Forms engine (dynamic schema, validation, submissions)
- Workflow engine (state machine, transitions, audit log)
- Dark admin theme
- Docker Compose stack
- CI pipeline (GitHub Actions)

## Ports

| Service | Port |
|---------|------|
| Go API | 8000 |
| Next.js | 3000 |
| PostgreSQL | 5432 |
| Redis | 6379 |

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) and the [stack primer](../primers/go-nextjs/PRIMER.md).

---

Boilerworks is a [CONFLICT](https://weareconflict.com) brand. CONFLICT is a registered trademark of CONFLICT LLC.
