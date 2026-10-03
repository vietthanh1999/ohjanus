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

const BASE_URL = (import.meta.env.VITE_API_BASE_URL as string | undefined) ?? 'http://127.0.0.1:8788';
const ADMIN_TOKEN = (import.meta.env.VITE_ADMIN_TOKEN as string | undefined) ?? '';

export function apiBaseUrl(): string {
  return BASE_URL.replace(/\/$/, '');
}

export function adminToken(): string {
  return ADMIN_TOKEN;
}

interface FetchOptions {
  method?: string;
  body?: unknown;
  signal?: AbortSignal;
}

export async function apiFetch<T>(path: string, opts: FetchOptions = {}): Promise<T> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' };
  if (ADMIN_TOKEN) headers['Authorization'] = `Bearer ${ADMIN_TOKEN}`;

  let resp: Response;
  try {
    resp = await fetch(`${apiBaseUrl()}${path}`, {
      method: opts.method ?? 'GET',
      headers,
      body: opts.body === undefined ? undefined : JSON.stringify(opts.body),
      signal: opts.signal
    });
  } catch (e) {
    throw new ApiError(0, 'NETWORK_ERROR', `Cannot reach Admin API at ${apiBaseUrl()}. Is janus serve running?`, '');
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
    throw new ApiError(
      resp.status,
      body?.error?.code ?? 'UNKNOWN',
      body?.error?.message ?? `Request failed with status ${resp.status}`,
      body?.error?.request_id ?? ''
    );
  }
  return data as T;
}

/** GET that downloads a blob (used for audit export, where <a href> cannot send headers). */
export async function apiDownload(path: string, filename: string): Promise<void> {
  const headers: Record<string, string> = {};
  if (ADMIN_TOKEN) headers['Authorization'] = `Bearer ${ADMIN_TOKEN}`;
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
