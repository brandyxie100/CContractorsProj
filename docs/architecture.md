# Architecture

## Stack rationale (SME fit)

### Backend: Go + Gin + GORM + PostgreSQL

- Single static binary — low ops cost on one VM or container
- Gin is thin enough for eight REST modules without framework weight
- GORM covers CRUD for a modest schema; PostgreSQL handles overlap queries well
- Explicit concurrency for availability / dispatch scoring
- Fits small teams: readable domain services, one language for API + workers later

### Frontend: React + TypeScript + Vite + Tailwind

- Dashboard / calendar SPA — Vite + React Router (not Next.js)
- TanStack Query for list/calendar refetch; React Hook Form + Zod for forms
- Local UI kit (Button, Card, Badge, Modal, Table, Input, Select) — faster than a heavy design system
- Recharts for Phase 1–2 KPIs without a BI stack

### Constraints

One monorepo, one Postgres, JWT + RBAC, no microservices or message bus until Phase 3.

## High-level diagram

```
┌─────────────────────────────────────┐
│  frontend (Vite React SPA)          │
│  Pages · UI kit · TanStack Query    │
└─────────────────┬───────────────────┘
                  │ HTTPS / JSON
┌─────────────────▼───────────────────┐
│  backend (Go Gin)                   │
│  Handlers → Services → Repositories │
│  JWT + RBAC middleware              │
└─────────────────┬───────────────────┘
                  │
┌─────────────────▼───────────────────┐
│  PostgreSQL                         │
└─────────────────────────────────────┘
```

## Monorepo layout

```
backend/
  cmd/api/main.go
  internal/
    config/
    middleware/
    models/
    repository/
    service/
    handler/
    dto/
  tests/
frontend/
  src/
    components/ui/
    components/layout/
    pages/
    features/*/hooks/
    lib/
    types/
docs/
docker-compose.yml
```

## Layering (backend)

| Layer | Responsibility |
|-------|----------------|
| Handler | HTTP bind/validate, status codes, call service |
| Service | Business rules: availability, conflicts, dispatch ranking |
| Repository | GORM queries only — no business decisions |
| Models | Persistence structs |

Handlers stay thin. Domain algorithms live in `internal/service`.

## Auth flow

1. `POST /api/v1/auth/login` — email/password → access JWT (short-lived) + refresh token (httpOnly cookie)
2. Frontend keeps access token **in memory** only
3. `POST /api/v1/auth/refresh` — cookie → new access token
4. `POST /api/v1/auth/logout` — clear refresh cookie
5. Middleware: `Authorization: Bearer <access>` + role check on mutating routes

## Roles

| Role | Read | Write assignments | Manage users | Settings |
|------|------|-------------------|--------------|----------|
| admin | yes | yes | yes | yes |
| manager | yes | yes | no | locations |
| operator | yes | limited | no | no |
| viewer | yes | no | no | no |

## Single-tenant MVP

One company (Clements). No `tenant_id` until Phase 3 multi-branch.
