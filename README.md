# nova

[![CI](https://github.com/githubuseradmin/nova-dashboard/actions/workflows/ci.yml/badge.svg)](https://github.com/githubuseradmin/nova-dashboard/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
![Go](https://img.shields.io/badge/Go-1.23%2B-00ADD8)
![Svelte](https://img.shields.io/badge/Svelte-5-FF3E00)

A clean, modern full-stack starter — **Go · Svelte · PostgreSQL**. A marketing landing page, session-based auth, a user dashboard and an admin panel — with a bespoke, hand-written UI and no framework bloat.

> 🇷🇺 Russian version: [README.ru.md](README.ru.md)

## Why this stack

| Layer | Choice | Why |
|---|---|---|
| **Backend** | Go — stdlib `net/http` + [`chi`](https://github.com/go-chi/chi) | Native `net/http` handlers, no `fasthttp`. `chi` adds routing + middleware only. |
| **Database** | PostgreSQL + [`pgx`](https://github.com/jackc/pgx) | Explicit SQL and migrations, no ORM magic. |
| **Frontend** | Svelte 5 + TypeScript + Vite | ~17 kB gzipped bundle, compiled, no virtual DOM. |
| **Styling** | Hand-written CSS + design tokens | No Tailwind — a bespoke light/dark theme via CSS custom properties. |
| **Auth** | Session cookie (`HttpOnly`/`Secure`/`SameSite`) + CSRF + argon2id | Safer than JWT-in-`localStorage`. |
| **Edge** | Nginx + Certbot | TLS; serves static, reverse-proxies `/api`. |

## Architecture

```
nova/
├── cmd/server/          # entrypoint: config, logger, HTTP server, graceful shutdown
├── internal/
│   ├── config/          # env-driven configuration
│   ├── http/            # transport: router (chi), middleware, handlers, responses
│   ├── auth/            # sessions, password hashing (argon2id), CSRF
│   ├── store/           # data access (pgx), migrations runner
│   └── models/          # domain types
├── migrations/          # SQL migrations (embedded)
└── web/
    ├── site/            # static marketing landing (own CSS)
    └── app/             # Svelte 5 + Vite SPA: login, dashboard, admin
```

Dependencies flow one way: `http → auth / store → models`. The transport layer never touches SQL directly — it goes through `store`.

## API

| Method | Path | Guard | Description |
|---|---|---|---|
| `GET` | `/healthz` | — | Liveness probe |
| `GET` | `/api/csrf` | — | Issue a CSRF token (cookie + body) |
| `POST` | `/api/auth/login` | CSRF | Verify credentials, start a session |
| `POST` | `/api/auth/logout` | CSRF | Destroy the session |
| `GET` | `/api/auth/me` | session | Current user |
| `GET` | `/api/dashboard` | session | Dashboard data |
| `GET` | `/api/admin/users` | admin | List users |
| `PUT` | `/api/admin/users/{id}/role` | admin + CSRF | Change a user's role |

## Quick start (development)

Requirements: **Go 1.23+**, **Node 20+**, **Docker**.

```sh
cp .env.example .env          # then edit if needed

make db-up                    # start Postgres (Docker)
make dev                      # API on http://localhost:8080  (migrates + seeds an admin)
make web-install              # once, to install frontend deps
make web-dev                  # Svelte dev server on http://localhost:5173 (proxies /api)
```

Open http://localhost:5173 and sign in with the seeded admin
(`NOVA_ADMIN_EMAIL` / `NOVA_ADMIN_PASSWORD`, defaults in `.env.example`).

## Testing

```sh
make test                     # unit tests always run
# store integration tests run only when a test database is configured:
NOVA_TEST_DATABASE_URL=postgres://nova:nova@localhost:5432/nova?sslmode=disable make test
```

CI (GitHub Actions) runs `go vet`, `go test -race` against a Postgres service, and
`svelte-check` + the production frontend build on every push.

## Production

```sh
make web-build                # builds web/app/dist with base /app/
make build                    # builds ./bin/server
```

The single Go binary serves the landing at `/`, the SPA at `/app/`, and the JSON API
at `/api`. Put Nginx + Certbot in front for TLS.

## License

[MIT](LICENSE).
