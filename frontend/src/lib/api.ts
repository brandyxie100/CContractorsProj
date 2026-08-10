import { getAccessToken, setAccessToken } from '@/lib/auth';

const API_BASE = '/api/v1';

export class ApiError extends Error {
  code: string;
  status: number;

  constructor(status: number, code: string, message: string) {
    super(message);
    this.status = status;
    this.code = code;
  }
}

async function refreshAccessToken(): Promise<boolean> {
  const res = await fetch(`${API_BASE}/auth/refresh`, {
    method: 'POST',
    credentials: 'include',
  });
  if (!res.ok) {
    setAccessToken(null);
    return false;
  }
  const data = (await res.json()) as { access_token: string };
  setAccessToken(data.access_token);
  return true;
}

export async function apiFetch<T>(path: string, init: RequestInit = {}, retry = true): Promise<T> {
  const headers = new Headers(init.headers);
  if (!headers.has('Content-Type') && init.body) {
    headers.set('Content-Type', 'application/json');
  }
  const token = getAccessToken();
  if (token) {
    headers.set('Authorization', `Bearer ${token}`);
  }

  const res = await fetch(`${API_BASE}${path}`, {
    ...init,
    headers,
    credentials: 'include',
  });

  if (res.status === 401 && retry && !path.startsWith('/auth/login')) {
    const ok = await refreshAccessToken();
    if (ok) {
      return apiFetch<T>(path, init, false);
    }
  }

  if (res.status === 204) {
    return undefined as T;
  }

  const text = await res.text();
  const data = text ? (JSON.parse(text) as unknown) : null;

  if (!res.ok) {
    const err = data as { error?: { code?: string; message?: string } } | null;
    throw new ApiError(
      res.status,
      err?.error?.code ?? 'ERROR',
      err?.error?.message ?? res.statusText,
    );
  }

  return data as T;
}
