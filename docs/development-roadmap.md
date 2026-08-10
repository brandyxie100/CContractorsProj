# Development Roadmap

## Conventions

- Branch from `main`: `feature/<module>-…`, `fix/…`, `chore/…`, `docs/…`
- Conventional commits: `feat(asset):`, `fix(dispatch):`, scopes = `auth|user|asset|project|assignment|maintenance|dashboard|dispatch|shared|frontend`
- Test-first for domain services (availability, conflict, ranking)
- PR &lt; 900 lines; squash merge to `main`

## Phase 0 — Foundations

**Goals**

- Monorepo scaffold (`backend/`, `frontend/`)
- Docker Compose Postgres
- Config via env
- CI: Go lint/test + frontend lint/typecheck/test/build
- JWT login/refresh/logout + RBAC middleware
- Seed admin user + sample locations

**Verify**

- [ ] `docker compose up -d` starts Postgres
- [ ] `go test ./...` passes
- [ ] Login returns access token + refresh cookie
- [ ] Viewer cannot POST `/assets`

## Phase 1 — MVP availability

**Goals**

- Asset / Location / Project / Assignment CRUD
- Conflict detection + `next_available_at`
- Calendar, Assets, Projects, Dashboard pages
- Status badges (5-color system)

**Verify**

- [ ] Overlapping assignment returns 409
- [ ] Calendar loads assignments for date range
- [ ] Asset detail shows next available
- [ ] Dashboard shows status counts

## Phase 2 — Smart ops

**Goals**

- Maintenance CRUD; windows block availability
- `POST /dispatch/search` ranking + Dispatch page
- Utilization reports (basic)
- Responsive sidebar layout

**Verify**

- [ ] Maintenance overlap blocks assignment
- [ ] Dispatch returns ranked list with distance
- [ ] Alternatives appear when exact type unavailable
- [ ] Reports page shows utilization metrics

## Phase 3 — Advanced (later)

Telematics, QR check-in, cost tracking, bundles, external APIs, multi-branch.

## Success criteria (overall)

Replace spreadsheet/whiteboard scheduling with a digital calendar, conflict-safe assignments, and a Smart Dispatch Finder for Clements Contractors.
