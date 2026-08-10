# Availability & Smart Dispatch

## Conflict detection

Before creating or updating an assignment:

1. Load non-cancelled assignments for `asset_id` where intervals overlap the requested window.
2. Load maintenance records with status in `scheduled | in_progress | overdue` that overlap.
3. If asset status is `retired` or `sold` → reject.
4. If any overlap → HTTP 409 CONFLICT with conflicting IDs.

Overlap:

```
existing.start_time < requested.end_time
AND existing.end_time > requested.start_time
```

## IsAvailable

```
is_available(asset_id, start, end):
  if asset.status in (retired, sold): return false
  if overlapping assignments: return false
  if overlapping maintenance: return false
  return true
```

## NextAvailableDate

```
next_available_date(asset_id):
  now = utc_now()
  current = assignment where start <= now < end (not cancelled)
  if current is nil:
    if in open maintenance: return maintenance.scheduled_end
    return now
  after = current.end_time
  maint = maintenance starting at/after current.end that abuts or overlaps
  if maint: after = max(after, maint.scheduled_end)
  return after
```

Also expose the next scheduled assignment start so the UI can show gaps.

## Smart Dispatch Finder ranking

Input: equipment type, time window, job coordinates (or project location), optional min specs, `include_alternatives`.

### Steps

1. Filter assets by `type` (exact). If alternatives enabled, also load related types (e.g. excavator → larger excavator).
2. For each candidate, compute:
   - **availability**: available for full window → 1.0; else score by how soon `next_available` is vs requested start (decay)
   - **proximity**: Haversine distance from `asset.current_location` to job; score = `1 / (1 + km/10)`
   - **specs match**: compare requested min weight/capacity to `specs` JSON; exact/over → 1.0; under → partial or flag alternative
3. **rank_score** = `0.50 * availability + 0.30 * proximity + 0.20 * specs`
4. Sort descending by `rank_score`; available-now + close first.
5. Alternatives (`is_alternative: true`) listed after exact type matches unless they score higher and exact type is empty.

### Ranking weights (MVP constants)

| Factor | Weight |
|--------|--------|
| Availability | 0.50 |
| Proximity | 0.30 |
| Specs match | 0.20 |

Tune later via config if needed — avoid premature configurability.

### Transport estimate

```
transport_minutes ≈ distance_km / 40 * 60   // assume ~40 km/h plant transport
```

No live GPS until Phase 3 — use Location lat/lng only.

## Status color mapping (UI)

| Status | Color |
|--------|-------|
| available | green |
| assigned | blue |
| in_transit | yellow |
| maintenance | orange |
| reserved | purple |
