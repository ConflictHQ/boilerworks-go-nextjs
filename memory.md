# Boilerworks Memory

This file is the **AI context seed** for the Boilerworks Go + Next.js template. It captures decisions, constraints, and non-obvious facts that are not derivable from reading a single file.

For conventions and patterns, see [`bootstrap.md`](bootstrap.md).

---

## Template purpose

Catalogue template (not a support repo): API-first Go backend (Chi router, REST/JSON) + Next.js 16 frontend, for teams that want Go's performance and small footprint behind a modern React frontend.

---

## Key architectural decisions

| Decision | Why |
|---|---|
| Session auth, not JWT | bcrypt password hashing + SHA256-hashed session tokens in httpOnly cookies; the cookie is shared with Next.js via `credentials: "include"` |
| MutationResult pattern | Every mutation returns `{ok: true}` or `{ok: false, errors: [...]}` |
| UUID PKs + soft deletes | UUIDs are the external identifiers in API routes; `deleted_at` on all content tables, never hard-delete |
| goose-format SQL migrations | `db/migrations/` (init + seed) |
| CORS driven by `FRONTEND_URL` | `internal/server/server.go` sets `AllowedOrigins` from `cfg.FrontendURL` — one origin only |

---

## Non-obvious facts

- **Port split**: local-dev defaults are API :8088, frontend :3004, Postgres :5447, Redis :6390 (`.env.example`, `frontend/package.json`); the conventional :8000 / :3000 / :5432 / :6379 exist only as Docker Compose host mappings (`docker/docker-compose.yml`).
- **Redis is provisioned but unused**: `internal/config` loads `REDIS_URL` and compose runs a redis service the API health-depends on, but `go.mod` has no Redis client and `Config.RedisURL` has no consumers yet.
- **The frontend reaches the API two ways**: browser calls use `NEXT_PUBLIC_API_URL` (`frontend/lib/api.ts`; empty means same-origin, proxied through Next.js rewrites), while the server-side rewrite target is `API_URL` (`frontend/next.config.ts`; defaults to `http://localhost:8088`).
- **Seeded admin**: `admin@boilerworks.dev` / `password` (from `db/migrations/002_seed.sql`) — change before any deploy.
