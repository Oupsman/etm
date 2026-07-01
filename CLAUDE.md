# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

ETM (Eisenhower Task Manager) is a multi-user task manager organized around the Eisenhower matrix (priority × urgency). It consists of:
- **Backend**: Go REST API using Gin + GORM + PostgreSQL
- **Frontend**: Vue 3 SPA (`yaftm/`) using Vuetify 3, Pinia, and Vue Router, served statically by the backend

## Build Commands

### Backend
```sh
go build -o main ./cmd/etm/main.go
go test ./...                          # run all tests
go test ./pkg/models/...               # run a specific package's tests
```

### Frontend (`yaftm/`)
```sh
npm install       # or: pnpm install
npm run dev       # dev server with hot reload (proxies API to backend)
npm run build     # production build → yaftm/dist/ (served by Go backend)
npm run lint      # eslint --fix
npm run type-check
```

### Docker
```sh
docker build . -t oupsman/etm
docker compose up
```

## Environment Variables

Copy `.env.template` to `.env` before running locally. Key variables:

| Variable | Default | Description |
|---|---|---|
| `HOST` | `""` | Bind address (empty = all interfaces) |
| `PORT` | `8080` | HTTP listen port |
| `DB_HOST` | `127.0.0.1` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DATABASE` | `etm` | Database name |
| `DB_USERNAME` | `etm` | Database user |
| `DB_PASSWORD` | `etmpass` | Database password |
| `SECRET_KEY` | — | JWT signing key (generate with `dd if=/dev/random bs=512 count=1 \| sha256sum`) |
| `TOKEN_DURATION` | `120` | JWT token lifetime in minutes |
| `OIDC_ENABLED` | `false` | Set to `"true"` to enable OIDC login |
| `OIDC_ISSUER_URL` | — | OIDC provider URL (e.g. Keycloak realm URL) |
| `OIDC_CLIENT_ID` | — | OIDC client ID |
| `OIDC_CLIENT_SECRET` | — | OIDC client secret |
| `OIDC_REDIRECT_URL` | — | Callback URL registered with the provider |

The database must exist before first run; GORM auto-migrates all tables on startup.

## Architecture

### Backend Package Layout

```
cmd/etm/main.go          Entry point: loads vars, creates App, initializes OIDC, starts server
pkg/vars/vars.go         Reads all env vars into package-level globals
pkg/app/app.go           App struct (holds DB + http.Client + Logger); NewApp() connects DB and runs migrations
pkg/models/              GORM models + DB methods (all as methods on *DB)
  database.go            DB type wrapping gorm.DB; ConnectToDB(), CreateOrMigrate()
  users.go / tasks.go / categories.go / devices.go / keys.go / tokens.go / oidc.go
pkg/types/types.go       Plain request-body structs shared between controllers and models
pkg/utils/utils.go       JWT parsing (ParseToken, GetUserID, GetUserUUID), bcrypt helpers
pkg/controllers/         Gin handler functions
  controllers.go         IsAuthorized() middleware — validates Bearer JWT, sets "userID" in context
  users.go / tasks.go / categories.go / tokens.go / notifications.go / oidc.go
pkg/webserver/router.go  Registers all routes under /api/v1, serves yaftm/dist as SPA
```

### Request Flow

1. `IsAuthorized()` middleware validates the Bearer JWT and writes `userID` (float64) and `uuid` into the Gin context.
2. Controllers retrieve the `*app.App` from context via `c.MustGet("App").(*app.App)` and call methods on `App.DB`.
3. All DB operations are methods on `models.DB` (which embeds `gorm.DB`).

### JWT Claims

Tokens are HS256-signed with `SECRET_KEY`. Standard claims used:
- `sub` — numeric user ID (stored as float64 in MapClaims)
- `uuid` — user UUID string
- `exp` — 30-minute expiry
- `iss` — `"etm"`

### OIDC Flow

When `OIDC_ENABLED=true`, `InitOIDC()` (called from main) sets up the provider. Login redirects to the provider; the callback exchanges the code, verifies the ID token, finds-or-creates the local user (via `OIDCSubject` field), and redirects to `/#/auth/callback?token=<jwt>` which the frontend `OIDCCallback.vue` picks up.

### Frontend Architecture

The SPA lives in `yaftm/src/`:
- **Pinia stores** (`src/stores/`): `auth` (JWT + localStorage), `user`, `task`, `category`, `device`, `snackbar`
- **Pages** (`src/pages/`): `index.vue` (main matrix view), `login.vue`, `signup.vue`
- **Views** (`src/views/`): `OIDCCallback.vue` (extracts token from URL fragment after OIDC redirect)
- **Axios plugin** (`src/plugins/axios.ts`): injects `Authorization: Bearer` header from auth store; retries once on 401 after token refresh
- All API calls use the shared `axiosInstance` from `src/plugins/axios.ts`

The built frontend (`yaftm/dist/`) is served at `/` by the Go backend via `gin-contrib/static`. During development, run both the Go backend and `npm run dev` separately; Vite proxies API calls to the backend.
