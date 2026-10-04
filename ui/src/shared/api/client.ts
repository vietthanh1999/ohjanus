import type { ApiErrorBody } from './types';

export class ApiError extends Error {
  status: number;
  code: string;
  requestId: string;

  constructor(status: number, code: string, message: string, requestId: string) {
    super(message);
    this.name = 'ApiError';
    this.status = status;
    this.code = code;
    this.requestId = requestId;
  }
}

const ENV_BASE_URL = (import.meta.env.VITE_API_BASE_URL as string | undefined) ?? '';
const ENV_ADMIN_TOKEN = (import.meta.env.VITE_ADMIN_TOKEN as string | undefined) ?? '';

const LS_BASE_URL_KEY = 'ohjanus_api_base_url';
const LS_TOKEN_KEY = 'ohjanus_admin_token';

export function apiBaseUrl(): string {
  try {
    const ls = localStorage.getItem(LS_BASE_URL_KEY);
    if (ls) return ls.replace(/\/$/, '');
  } catch { /* ignore */ }
  if (ENV_BASE_URL) return ENV_BASE_URL.replace(/\/$/, '');
  // Same-origin by default so Vite proxy (/api -> :8788) works without CORS.
  return '';
}

export function setApiBaseUrl(url: string): void {
  try {
    if (url) localStorage.setItem(LS_BASE_URL_KEY, url);
    else localStorage.removeItem(LS_BASE_URL_KEY);
  } catch { /* ignore */ }
}

export function adminToken(): string {
  try {
    const ls = localStorage.getItem(LS_TOKEN_KEY);
    if (ls) return ls;
  } catch { /* ignore */ }
  return ENV_ADMIN_TOKEN;
}

export function setAdminToken(token: string): void {
  try {
    if (token) localStorage.setItem(LS_TOKEN_KEY, token);
    else localStorage.removeItem(LS_TOKEN_KEY);
  } catch { /* ignore */ }
}

function authHeaders(): Record<string, string> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  const token = adminToken();
  if (token) headers['Authorization'] = `Bearer ${token}`;
  return headers;
}

interface FetchOptions {
  method?: string;
  body?: unknown;
  signal?: AbortSignal;
}

export async function apiFetch<T>(path: string, opts: FetchOptions = {}): Promise<T> {
  const headers = authHeaders();
  const base = apiBaseUrl();

  let resp: Response;
  try {
    resp = await fetch(`${base}${path}`, {
      method: opts.method ?? 'GET',
      headers,
      body: opts.body === undefined ? undefined : JSON.stringify(opts.body),
      signal: opts.signal
    });
  } catch (e) {
    throw new ApiError(0, 'NETWORK_ERROR', `Cannot reach Admin API at ${base || 'same origin'}. Is janus serve running?`, '');
  }

  if (resp.status === 204) return undefined as T;

  let data: any = undefined;
  try {
    data = await resp.json();
  } catch {
    // Non-JSON body (should not happen for this API).
  }

  if (!resp.ok) {
    const body = data as ApiErrorBody | undefined;
    const code = body?.error?.code;
    // Vite proxy (or any gateway) answers 502/503/504 with an HTML body when
    // janus serve is down. Surface that as NETWORK_ERROR with the actionable
    // hint instead of a cryptic "Request failed with status 502".
    if (!code && (resp.status === 502 || resp.status === 503 || resp.status === 504)) {
      throw new ApiError(resp.status, 'NETWORK_ERROR', `Cannot reach Admin API at ${base || 'same origin'}. Is janus serve running?`, '');
    }
    throw new ApiError(
      resp.status,
      code ?? 'UNKNOWN',
      body?.error?.message ?? `Request failed with status ${resp.status}`,
      body?.error?.request_id ?? ''
    );
  }
  return data as T;
}

/** GET that downloads a blob (used for audit export, where <a href> cannot send headers). */
export async function apiDownload(path: string, filename: string): Promise<void> {
  const headers: Record<string, string> = {};
  const token = adminToken();
  if (token) headers['Authorization'] = `Bearer ${token}`;
  let resp: Response;
  try {
    resp = await fetch(`${apiBaseUrl()}${path}`, { headers });
  } catch {
    throw new ApiError(0, 'NETWORK_ERROR', `Cannot reach Admin API at ${apiBaseUrl()}.`, '');
  }
  if (!resp.ok) {
    let message = `Export failed with status ${resp.status}`;
    try {
      const body = (await resp.json()) as ApiErrorBody;
      message = body?.error?.message ?? message;
    } catch { /* keep default */ }
    throw new ApiError(resp.status, 'EXPORT_FAILED', message, '');
  }
  const blob = await resp.blob();
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  link.remove();
  URL.revokeObjectURL(url);
}
