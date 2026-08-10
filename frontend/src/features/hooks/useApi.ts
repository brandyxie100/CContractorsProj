import { useQuery } from '@tanstack/react-query';
import { apiFetch } from '@/lib/api';
import type { Asset, Assignment, DashboardSummary, DispatchResult, Location, MaintenanceRecord, Project, User } from '@/types';

export function useDashboard() {
  return useQuery({
    queryKey: ['dashboard'],
    queryFn: () => apiFetch<DashboardSummary>('/dashboard/summary'),
  });
}

export function useAssets(params: Record<string, string> = {}) {
  const qs = new URLSearchParams(params).toString();
  return useQuery({
    queryKey: ['assets', params],
    queryFn: () => apiFetch<Asset[]>(`/assets${qs ? `?${qs}` : ''}`),
  });
}

export function useProjects(status = '') {
  const qs = status ? `?status=${encodeURIComponent(status)}` : '';
  return useQuery({
    queryKey: ['projects', status],
    queryFn: () => apiFetch<Project[]>(`/projects${qs}`),
  });
}

export function useAssignments(from: string, to: string) {
  return useQuery({
    queryKey: ['assignments', from, to],
    queryFn: () =>
      apiFetch<Assignment[]>(
        `/assignments?from=${encodeURIComponent(from)}&to=${encodeURIComponent(to)}`,
      ),
  });
}

export function useMaintenance() {
  return useQuery({
    queryKey: ['maintenance'],
    queryFn: () => apiFetch<MaintenanceRecord[]>('/maintenance'),
  });
}

export function useLocations() {
  return useQuery({
    queryKey: ['locations'],
    queryFn: () => apiFetch<Location[]>('/locations'),
  });
}

export function useUsers() {
  return useQuery({
    queryKey: ['users'],
    queryFn: () => apiFetch<User[]>('/users'),
  });
}

export async function searchDispatch(body: Record<string, unknown>) {
  return apiFetch<{ results: DispatchResult[] }>('/dispatch/search', {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export async function createAssignment(body: Record<string, unknown>) {
  return apiFetch<Assignment>('/assignments', {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export async function createMaintenance(body: Record<string, unknown>) {
  return apiFetch<MaintenanceRecord>('/maintenance', {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export async function createAsset(body: Record<string, unknown>) {
  return apiFetch<Asset>('/assets', {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export async function createProject(body: Record<string, unknown>) {
  return apiFetch<Project>('/projects', {
    method: 'POST',
    body: JSON.stringify(body),
  });
}

export async function createLocation(body: Record<string, unknown>) {
  return apiFetch<Location>('/locations', {
    method: 'POST',
    body: JSON.stringify(body),
  });
}
