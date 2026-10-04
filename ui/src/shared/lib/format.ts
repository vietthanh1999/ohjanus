/** 'YYYY-MM-DD HH:mm:ss' stamp for the activity log. */
export function nowStamp(): string {
  return new Date().toISOString().replace('T', ' ').substring(0, 19);
}

/** Render any cell value as text; null/undefined become the explicit `<null>` marker. */
export function cellText(v: unknown): string {
  if (v === null || v === undefined) return '<null>';
  if (typeof v === 'object') return JSON.stringify(v);
  return String(v);
}

export function isNullCell(v: unknown): boolean {
  return v === null || v === undefined;
}

/** '586 ms' / '1 s 159 ms' — same shape as the gateway latency badges. */
export function formatLatency(ms: number): string {
  if (ms < 1000) return `${ms} ms`;
  const s = Math.floor(ms / 1000);
  return `${s} s ${ms - s * 1000} ms`;
}
