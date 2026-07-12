# Calliope — Boilerworks Go + Next.js
<!-- Agent shim for https://github.com/calliopeai/calliope-cli -->

Primary conventions doc: [`bootstrap.md`](bootstrap.md)
Context seed: [`memory.md`](memory.md)

Read both before writing any code.

---

## Project-specific notes

- Go 1.25+ (Chi) backend + Next.js 16 (App Router, React 19, TypeScript) frontend; Postgres 16 (pgx/v5) + Redis 7; Docker Compose stack (api + ui + postgres + redis).
- Session-based auth (bcrypt + SHA256 token hashing, httpOnly cookies) shared between Go API and Next.js via `credentials: "include"`.
- REST API under `/api/*` returning JSON; `MutationResult` pattern `{ok: true}` / `{ok: false, errors: [...]}`; group-based permissions via `middleware.RequirePermission()`.
- UUID primary keys everywhere; soft deletes via `deleted_at`.
- Ports differ Docker vs local: Go API 8000/8088, Next.js 3000/3004, Postgres 5432/5447, Redis 6379/6390.
- Default admin credentials: `admin@boilerworks.dev` / `password`.
