# Data Model

## Entity relationship

```
User ─────────────┐
                  │ assigned_by
Location ◄── Asset ──► Assignment ──► Project
   ▲                │
   │                └──► MaintenanceRecord
   │
Project.location (optional link via address / location_id)
```

## Entities

### User

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| name | string | |
| email | string unique | login |
| password_hash | string | bcrypt |
| role | enum | admin, manager, operator, viewer |
| phone | string nullable | |
| avatar_url | string nullable | |
| created_at / updated_at | timestamptz | |

### Location

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| name | string | e.g. Main Depot |
| type | enum | depot, job_site, workshop |
| address | string | |
| lat / lng | float8 nullable | Haversine for dispatch |
| description | text nullable | |

### Asset

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| asset_code | string unique | human-readable EXC-007 |
| name | string | |
| model | string | |
| type | string | excavator, truck, roller, … |
| category | string | heavy_equipment, trucks, tools |
| specs | jsonb | weight_t, capacity, attachments |
| status | enum | see Status |
| current_location_id | UUID FK → Location | |
| purchase_date | date nullable | |
| purchase_cost | numeric nullable | |
| hourly_rate | numeric nullable | internal cost |
| photo_url | string nullable | |
| qr_code | string nullable | Phase 3 |
| created_at / updated_at | timestamptz | |

### Project

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| name | string | |
| address | string | |
| location_id | UUID FK nullable | |
| start_date / end_date | date | |
| status | enum | planning, active, completed |
| project_manager_id | UUID FK → User nullable | |
| description | text | |

### Assignment

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| asset_id | UUID FK | |
| project_id | UUID FK | |
| start_time / end_time | timestamptz | half-open overlap: start < other.end AND end > other.start |
| status | enum | scheduled, active, completed, cancelled |
| assigned_by | UUID FK → User | |
| notes | text | |
| created_at | timestamptz | |

### MaintenanceRecord

| Column | Type | Notes |
|--------|------|-------|
| id | UUID PK | |
| asset_id | UUID FK | |
| type | enum | preventive, repair, inspection |
| description | text | |
| scheduled_start / scheduled_end | timestamptz | blocks availability |
| completed_at | timestamptz nullable | |
| status | enum | scheduled, in_progress, completed, overdue |
| cost | numeric nullable | |
| performed_by | string nullable | |
| notes | text | |

## Enums

### Asset status (5-state)

| Value | Meaning | Assignable now? |
|-------|---------|-----------------|
| available | Ready at depot | yes |
| assigned | On a job | schedule after current window |
| in_transit | Moving between sites | no (temporary) |
| maintenance | Out of service | no |
| reserved | Booked upcoming | if no conflict |
| retired / sold | Terminal | never |

### Roles

`admin` | `manager` | `operator` | `viewer`

## Indexes (overlap & filters)

```sql
CREATE INDEX idx_assignments_asset_time ON assignments (asset_id, start_time, end_time)
  WHERE status NOT IN ('cancelled', 'completed');

CREATE INDEX idx_maintenance_asset_time ON maintenance_records (asset_id, scheduled_start, scheduled_end)
  WHERE status IN ('scheduled', 'in_progress', 'overdue');

CREATE INDEX idx_assets_type_status ON assets (type, status);
CREATE INDEX idx_assets_location ON assets (current_location_id);
```

## Overlap predicate

Two intervals `[A_start, A_end)` and `[B_start, B_end)` overlap when:

```
A_start < B_end AND A_end > B_start
```
