# Frontend Architecture

## Stack

- Vite + React 18 + TypeScript (strict)
- React Router
- Tailwind CSS
- TanStack Query
- React Hook Form + Zod
- Recharts

## Directory layout

```
frontend/src/
  components/ui/          Button, Card, Badge, Modal, Table, Input, Select
  components/layout/      AppShell, Sidebar
  pages/                  one file per route page
  features/
    assets/hooks/
    assignments/hooks/
    dispatch/hooks/
    ...
  lib/api.ts              fetch wrapper + refresh retry
  lib/auth.ts             in-memory access token
  types/                  shared DTOs
  App.tsx
  main.tsx
```

## Pages (8)

| Route | Page | Phase |
|-------|------|-------|
| `/` | Dashboard | 1 |
| `/calendar` | Availability Calendar | 1 |
| `/dispatch` | Dispatch Finder | 2 |
| `/assets` | Assets | 1 |
| `/projects` | Projects | 1 |
| `/maintenance` | Maintenance | 2 |
| `/reports` | Reports | 2 |
| `/settings` | Settings (users, locations) | 0–1 |

## Patterns

### Container / Presenter

- `*Container` — Query hooks, mutations, wire handlers
- `*View` — pure UI props

### Data fetching

- All HTTP via `lib/api.ts` and feature hooks (`useAssets`, `useAssignments`, …)
- No direct `fetch` in presentational components
- Query keys: `['assets', filters]`, `['assignments', { from, to }]`, …

### Auth

- Access token in memory (`lib/auth.ts`)
- Refresh token httpOnly cookie set by backend
- On 401: try refresh once, then redirect to `/login`

### Forms

- React Hook Form + Zod schemas mirroring API DTOs
- Absolute imports `@/`

### UI kit

Minimal local components only — no heavy component library.

| Component | Use |
|-----------|-----|
| Button | primary / secondary / danger |
| Card | section containers (interaction surfaces only where needed) |
| Badge | status colors |
| Modal | create/edit dialogs |
| Table | list pages |
| Input / Select | forms |

### Calendar

Primary view: horizontal timeline (assets × time). Color-coded assignment bars. Filters by type/status/location. Day/week/month range via `from`/`to` query params.

### Dispatch Finder UI

Form: type, date range, project (or lat/lng), optional min specs → ranked result list with Assign action.

## Forbidden (project rules)

- `any` type
- `localStorage` for tokens
- Direct API calls in dumb components
- Inline styles — use Tailwind utilities
