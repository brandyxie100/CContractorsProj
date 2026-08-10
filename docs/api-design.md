# API Design

Base path: `/api/v1`  
Content-Type: `application/json`  
Auth: `Authorization: Bearer <access_token>` (except login/refresh)

## Error shape

```json
{
  "error": {
    "code": "CONFLICT",
    "message": "Assignment overlaps existing booking for asset EXC-007",
    "details": { "conflicting_assignment_id": "..." }
  }
}
```

| HTTP | code examples |
|------|----------------|
| 400 | VALIDATION |
| 401 | UNAUTHORIZED |
| 403 | FORBIDDEN |
| 404 | NOT_FOUND |
| 409 | CONFLICT |
| 500 | INTERNAL |

## RBAC matrix

| Endpoint group | viewer | operator | manager | admin |
|----------------|--------|----------|---------|-------|
| GET * (read) | yes | yes | yes | yes |
| POST/PATCH assignments | no | yes | yes | yes |
| POST/PATCH assets, projects, maintenance | no | no | yes | yes |
| User CRUD / role change | no | no | no | yes |
| Dispatch search | yes | yes | yes | yes |
| Assign from dispatch | no | yes | yes | yes |

## Auth

| Method | Path | Body | Response |
|--------|------|------|----------|
| POST | `/auth/login` | `{ email, password }` | `{ access_token, user }` + Set-Cookie refresh |
| POST | `/auth/refresh` | (cookie) | `{ access_token }` |
| POST | `/auth/logout` | — | 204 + clear cookie |
| GET | `/auth/me` | — | `{ user }` |

## Users

| Method | Path | Notes |
|--------|------|-------|
| GET | `/users` | admin |
| POST | `/users` | admin |
| GET | `/users/:id` | admin/manager |
| PATCH | `/users/:id` | admin |
| DELETE | `/users/:id` | admin |

## Locations

| Method | Path | Notes |
|--------|------|-------|
| GET | `/locations` | list |
| POST | `/locations` | manager+ |
| GET | `/locations/:id` | |
| PATCH | `/locations/:id` | manager+ |
| DELETE | `/locations/:id` | manager+ |

## Assets

| Method | Path | Notes |
|--------|------|-------|
| GET | `/assets` | query: `type`, `status`, `location_id`, `q` |
| POST | `/assets` | manager+ |
| GET | `/assets/:id` | includes `next_available_at` |
| PATCH | `/assets/:id` | manager+ |
| DELETE | `/assets/:id` | manager+ (soft: status=retired preferred) |

## Projects

| Method | Path | Notes |
|--------|------|-------|
| GET | `/projects` | query: `status` |
| POST | `/projects` | manager+ |
| GET | `/projects/:id` | |
| PATCH | `/projects/:id` | manager+ |
| DELETE | `/projects/:id` | manager+ |

## Assignments

| Method | Path | Notes |
|--------|------|-------|
| GET | `/assignments` | query: `from`, `to`, `asset_id`, `project_id` (calendar range) |
| POST | `/assignments` | conflict check → 409 if overlap |
| GET | `/assignments/:id` | |
| PATCH | `/assignments/:id` | re-check conflicts |
| DELETE | `/assignments/:id` | cancel |

### Create assignment body

```json
{
  "asset_id": "uuid",
  "project_id": "uuid",
  "start_time": "2026-08-12T08:00:00Z",
  "end_time": "2026-08-14T17:00:00Z",
  "notes": "Smith Street earthworks"
}
```

## Maintenance

| Method | Path | Notes |
|--------|------|-------|
| GET | `/maintenance` | query: `asset_id`, `status` |
| POST | `/maintenance` | manager+; blocks availability |
| GET | `/maintenance/:id` | |
| PATCH | `/maintenance/:id` | |
| DELETE | `/maintenance/:id` | |

## Dashboard

| Method | Path | Response |
|--------|------|----------|
| GET | `/dashboard/summary` | fleet counts by status, overdue maintenance, utilization snapshot, upcoming assignments |

## Dispatch (Smart Dispatch Finder)

| Method | Path | Notes |
|--------|------|-------|
| POST | `/dispatch/search` | ranked results |

### Search body

```json
{
  "equipment_type": "excavator",
  "start_time": "2026-08-12T08:00:00Z",
  "end_time": "2026-08-14T17:00:00Z",
  "project_id": "uuid",
  "job_lat": -36.85,
  "job_lng": 174.76,
  "min_weight_t": 20,
  "include_alternatives": true
}
```

### Search response item

```json
{
  "asset": { "id": "...", "asset_code": "EXC-007", "type": "excavator", "specs": {} },
  "available": true,
  "available_from": "2026-08-12T08:00:00Z",
  "distance_km": 12.4,
  "transport_minutes": 18,
  "suitability_score": 0.95,
  "rank_score": 0.88,
  "is_alternative": false,
  "warnings": []
}
```
