export type Role = 'admin' | 'manager' | 'operator' | 'viewer';

export type AssetStatus =
  | 'available'
  | 'assigned'
  | 'in_transit'
  | 'maintenance'
  | 'reserved'
  | 'retired'
  | 'sold';

export interface User {
  id: string;
  name: string;
  email: string;
  role: Role;
  phone?: string | null;
  avatar_url?: string | null;
}

export interface Location {
  id: string;
  name: string;
  type: string;
  address: string;
  lat?: number | null;
  lng?: number | null;
  description?: string | null;
}

export interface Asset {
  id: string;
  asset_code: string;
  name: string;
  model: string;
  type: string;
  category: string;
  specs: Record<string, unknown>;
  status: AssetStatus;
  current_location_id?: string | null;
  current_location?: Location | null;
  hourly_rate?: number | null;
}

export interface Project {
  id: string;
  name: string;
  address: string;
  location_id?: string | null;
  start_date: string;
  end_date: string;
  status: string;
  description: string;
  location?: Location | null;
}

export interface Assignment {
  id: string;
  asset_id: string;
  project_id: string;
  start_time: string;
  end_time: string;
  status: string;
  notes: string;
  asset?: Asset;
  project?: Project;
}

export interface MaintenanceRecord {
  id: string;
  asset_id: string;
  type: string;
  description: string;
  scheduled_start: string;
  scheduled_end: string;
  status: string;
  notes: string;
  asset?: Asset;
}

export interface DashboardSummary {
  fleet_by_status: Record<string, number>;
  overdue_maintenance: number;
  active_assignments: number;
  utilization_percent: number;
  upcoming_assignments: number;
}

export interface DispatchResult {
  asset: Asset;
  available: boolean;
  available_from: string;
  distance_km?: number;
  transport_minutes?: number;
  suitability_score: number;
  rank_score: number;
  is_alternative: boolean;
  warnings: string[];
}
