# Đặc tả Kỹ thuật: **ohjanus-ui** — Trang admin cho MCP Gateway

> **Codename**: `ohjanus-ui`
> **Repo**: **riêng biệt** với `janus` core
> **Framework**: Svelte 5 + SvelteKit + `@ohjanus/ui` kit (xem `ui-kit.md`)
> **Mục tiêu**: Giao diện web tối giản cho approval workflow và audit viewer.
> **Backend yêu cầu**: Admin API riêng port 8788 (không dùng chung MCP).
> **Ngày**: 2026-10-03

## 0. Nguyên tắc thiết kế

1. **FE là client, không phải core**: chỉ gọi Admin API. Không business logic trong FE.
2. **Minimal surface area**: chỉ approval + audit + tokens + connections. Không SQL editor, schema browser.
3. **Read-only mặc định**: destructive action phải có confirmation + audit trail.
4. **Không lưu secret**: không thấy credentials, MCP token gốc, DSN. Chỉ metadata.
5. **Runes-first**: Svelte 5 runes (`$state`, `$derived`, `$effect`, `$props`), không store cũ.
6. **Type-safe end-to-end**: TypeScript strict, Zod validation, typed API client.
7. **Deploy riêng**: static SPA (adapter-static), Vercel/Netlify/S3 hoặc Go embed.

## 1. Phạm vi (Scope)

### 1.1. Trong phạm vi (v1.0)

| Màn hình | Mục đích | Ưu tiên |
| :--- | :--- | :--- |
| Login | SSO hoặc local auth | P0 |
| Approval Queue | Write request chờ duyệt | P0 |
| Approval Detail | Chi tiết 1 request + approve/reject | P0 |
| Audit Log Viewer | Xem, filter, search audit | P0 |
| Connections | Danh sách connection + trạng thái | P1 |
| Tokens | Xem, tạo, revoke MCP tokens | P1 |
| Dashboard | Metrics tổng quan | P2 |
| Settings | User, notification prefs | P2 |

### 1.2. Non-goals

- Không SQL editor / query runner, schema browser, config editor, dashboard thay Grafana, user management phức tạp, real-time collaboration. Responsive tới tablet, không ưu tiên mobile.

## 2. Kiến trúc

```
Browser (SvelteKit SPA: Svelte 5 + Query + Tailwind + @ohjanus/ui)
  │ HTTPS, session cookie HttpOnly Secure
  ▼
Reverse Proxy (Nginx/Caddy): / → static, /api/* → Admin API :8788, /stream → SSE (buffering off)
  ▼
Janus Core: MCP Server (agent) | Admin API :8788 (UI) | Approval Engine
```

- Admin API tách khỏi MCP server: khác port, khác auth (MCP opaque token, Admin session cookie/JWT).
- Không endpoint nào trả DSN/password/MCP token gốc.
- FE static, không SSR.

Cấu trúc `ohjanus-ui/`:

```
src/
  app.html, app.d.ts (App.Events SSE typing), app.css
  lib/api/client.ts (fetch credentials:include), auth.ts, approvals.ts, audit.ts, connections.ts, tokens.ts, types.ts
  lib/schemas/index.ts (Zod)
  lib/components/layout/AppShell,Sidebar,TopBar.svelte (dùng kit primitives + Tailwind)
  lib/components/approval/ApprovalCard,ApprovalDetail,SqlPreview,ParamsList,ApprovalActions.svelte (dùng `@ohjanus/ui`)
  lib/components/audit/AuditTable,AuditFilters,AuditDetail.svelte (dùng `@ohjanus/ui`)
  lib/components/common/EmptyState,LoadingSkeleton,ErrorState.svelte (wrap kit)
  lib/hooks/useAuth,useApprovals,useAudit,useTokens,useSSE.svelte.ts
  lib/utils/format.ts, sql.ts
  routes/+layout.svelte, +layout.ts (QueryClientProvider), +page.svelte (→/approvals),
    login/+page.svelte, approvals/+page.svelte, approvals/[id]/+page.svelte,
    audit/+page.svelte, connections/+page.svelte, tokens/+page.svelte, dashboard/+page.svelte
tests/unit,integration,e2e — static/, svelte.config.js, vite.config.ts, playwright.config.ts, .env.example
```

## 3. Tech Stack

### 3.1. Core

| Thành phần | Lựa chọn | Lý do |
| :--- | :--- | :--- |
| Framework | Svelte 5 | Runes fine-grained, bundle nhỏ |
| Meta-framework | SvelteKit | Routing, adapter-static cho SPA |
| Ngôn ngữ | TypeScript 5 (strict) | Type safety |
| Build | Vite 5 | Mặc định SvelteKit |

### 3.2. UI & Styling (dùng kit `ui-kit.md`)

| Thành phần | Lựa chọn | Ghi chú |
| :--- | :--- | :--- |
| CSS | Tailwind CSS 4 | Utility-first, tokens từ `@ohjanus/tokens` |
| Components | `@ohjanus/ui` | Button, Input, Modal, Drawer, Tabs, Table, Badge, Toast... — không tự build lẻ |
| Icons | `@ohjanus/icons` | Nhất quán kit |
| Toast | `@ohjanus/ui` Toast/Toaster | Queue `$state` |

### 3.3. Data & State

| Thành phần | Lựa chọn | Ghi chú |
| :--- | :--- | :--- |
| Server state | @tanstack/svelte-query v6 | Query/cache/mutation, runes native |
| Client state | Svelte 5 runes | `$state` UI state |
| URL state | SvelteKit `$page` | Filter/pagination query string |
| Form | svelte-form-hook | Zod resolver |
| Validation | Zod | Schema + type inference |

### 3.4. Domain-specific

| Thành phần | Lựa chọn | Ghi chú |
| :--- | :--- | :--- |
| SQL highlight | @agnosticeng/editor | Svelte 5 wrap CodeMirror 6, PG/MySQL dialect |
| Charts | @faintshadow/flarecharts | Svelte 5 native, runes-first, SVG, a11y |
| SSE | @sourceregistry/sveltekit-eventsource | Typed end-to-end, custom channels |
| Date | date-fns | Nhẹ, tree-shakeable |

### 3.5. Dev & Test

ESLint + typescript-eslint, Prettier + prettier-plugin-svelte, Vitest, @testing-library/svelte, Playwright, MSW.

### 3.6. Dependencies

```json
{
  "dependencies": {
    "@ohjanus/ui": "workspace:*",
    "@ohjanus/tokens": "workspace:*",
    "@ohjanus/icons": "workspace:*",
    "@tanstack/svelte-query": "^6.0.0",
    "svelte-form-hook": "^1.1.8",
    "zod": "^3.23.0",
    "@agnosticeng/editor": "^0.0.6",
    "@faintshadow/flarecharts": "^26.3.1",
    "@sourceregistry/sveltekit-eventsource": "^1.1.2",
    "date-fns": "^3.6.0",
    "tailwindcss": "^4.0.0"
  },
  "devDependencies": {
    "@sveltejs/kit": "^2.0.0",
    "@sveltejs/adapter-static": "^3.0.0",
    "svelte": "^5.0.0",
    "vite": "^5.0.0",
    "vitest": "^2.0.0",
    "@testing-library/svelte": "^5.0.0",
    "@playwright/test": "^1.45.0",
    "msw": "^2.0.0",
    "typescript": "^5.5.0",
    "eslint": "^9.0.0",
    "prettier": "^3.3.0",
    "prettier-plugin-svelte": "^3.2.0"
  }
}
```

## 4. API Contract (Admin API — yêu cầu Janus Core)

- Base `/api/v1`, JSON, auth cookie `janus_session` (HttpOnly, Secure, SameSite=Lax). CSRF header `X-Requested-With` cho mutating.

### Auth

- `POST /api/v1/auth/login` `{email,password}` → `{user:{id,email,name,role}}`
- `POST /api/v1/auth/logout` → 204
- `GET /api/v1/auth/me` → `{user:{id,email,name,role,permissions[]}}`
- SSO OIDC: `/api/v1/auth/oidc/authorize` → IdP → `/callback` → set cookie → FE gọi `/me`.

### Approvals

- `GET /api/v1/approvals?state=pending|approved|rejected|expired&connection&requested_by&from&to&cursor&limit` → `{items[], next_cursor, total}`
  - item: `id, state, created_at, expires_at, connection, statement_type, sql, params, affected_estimate, requested_by{token_id,client}, warnings[]`
- `GET /api/v1/approvals/{id}` → full: thêm `sql_hash, params_hash, plan, decided_by, decided_at, decision_reason`
- `POST /api/v1/approvals/{id}/approve` `{reason}` → `{id,state:approved,decided_by,decided_at,decision_reason}`
- `POST /api/v1/approvals/{id}/reject` `{reason*}` → `{id,state:rejected,...}`
- `GET /api/v1/approvals/stream` SSE: `approval.created`, `approval.expired`.

### Audit

- `GET /api/v1/audit?event&connection&token_id&request_id&status&from&to&q&cursor&limit` → `{items[], next_cursor, total}`
  - item: `id, ts, event, request_id, token_id, client, connection, tool, sql_normalized, statement_type, tables[], policy_decision, policy_rule, row_count, truncated, duration_ms, status, error`
- `GET /api/v1/audit/{id}` (kèm params nếu có `audit:read:params`).
- `GET /api/v1/audit/export` CSV/JSONL, rate limit.

### Connections

- `GET /api/v1/connections` → `{items:{name,driver,readonly,status,last_ping_at,pool{open,idle,in_use}}}` — không DSN/host/user/password.
- `POST /api/v1/connections/{name}/test` → latency + DB version.

### Tokens

- `GET /api/v1/tokens` → `{items:{id,name,scopes,created_at,expires_at,last_used_at,state}}`
- `POST /api/v1/tokens` `{name,scopes,ttl_hours}` → `{id,token,...}` — token chỉ trả 1 lần, FE cảnh báo copy ngay.
- `DELETE /api/v1/tokens/{id}` revoke.

### Dashboard (P2)

- `GET /api/v1/dashboard/summary` → `{requests_24h,denials_24h,pending_approvals,p95_latency_ms,active_tokens,connections_healthy,connections_total}`

### Error format

```json
{"error": {"code": "APPROVAL_NOT_FOUND", "message": "...", "request_id": "..."}}
```

400/401/403/404/409 (đã xử lý)/410 (hết hạn)/429/500.

## 5. Màn hình chi tiết

### 5.1. Login `/login`

Form email+password (Zod + svelte-form-hook), SSO button, submit → login, success → `/approvals`, fail không lộ email tồn tại. Rate limit 5/phút/IP.

### 5.2. Approval Queue `/approvals`

Tabs Pending/Approved/Rejected/Expired/All, card hiển thị statement type màu (SELECT xanh dương, INSERT xanh lá, UPDATE vàng, DELETE cam, DROP đỏ), SQL preview 3 dòng readonly (@agnosticeng/editor), affected estimate, requester, countdown. Click → detail. SSE toast + badge. Empty: "Không có yêu cầu nào 🎉".

### 5.3. Approval Detail `/approvals/[id]`

Thông tin chung (connection, requester, request ID, tạo/hết hạn + countdown), SQL readonly PG dialect, params list, ước lượng tác động + warnings, EXPLAIN plan collapsible, form quyết định (lý do optional khi approve, bắt buộc khi reject). Approve → confirm dialog. Hết hạn → disable. 409 → báo đã xử lý bởi người khác.

### 5.4. Audit `/audit`

Filters (event, connection, status, date, search debounce 300ms, shareable query string), table time/event/conn/tool/rows/ms, click row → drawer chi tiết (SQL, kết quả, policy). Export CSV/JSONL giới hạn 10k rows.

### 5.5. Connections `/connections`

Card: name, driver, status healthy/degraded, readonly, pool stats, last ping, nút Test + drawer (allowed schemas, denied tables, query 24h, link Grafana).

### 5.6. Tokens `/tokens`

List id/name/scopes/tạo/hết hạn/last used + Revoke (confirm). Modal tạo (name, scopes, TTL). Sau tạo: modal cảnh báo hiển thị 1 lần + Copy + checkbox "đã lưu".

### 5.7. Dashboard `/dashboard` (P2)

KPI cards, line chart requests over time (flarecharts), top denials, top connections, link Grafana.

## 6. Auth & Phân quyền

- Local: login → cookie → user info memory (không localStorage) → 401 redirect login.
- OIDC SSO: `/login` → SSO → `/auth/oidc/authorize?redirect_uri` → IdP → `/callback` → cookie → `/me`.
- Roles: `viewer(audit:read,connection:read)`, `approver(+approval:*)`, `admin(+token:*,connection:*,config:read)`, `auditor(audit:read,export)`. FE ẩn/hiện theo permission, server luôn validate.
- Session: TTL 8h sliding, CSRF `SameSite=Lax` + header `X-Requested-With`, logout xóa session.

## 7. Real-time (SSE)

`app.d.ts` typing `App.Events`: `approval.created`, `approval.expired`, `approval.decided`. Server forward SSE từ Janus Core (`src/routes/api/v1/approvals/stream/+server.ts`). Client `useApprovalSSE` (`useSSE.svelte.ts`): `onMount` mở `EventSource`, `on('approval.created/expired')` → invalidate queries + toast, `onDestroy` close. Auto-reconnect exponential backoff. Fallback polling 10s khi tab active nếu SSE bị chặn.

## 8. UX

Không giấu SQL/params/affected rows; confirm mọi destructive; không undo được nên confirm kỹ; loading/error/empty states rõ; error actionable; shortcuts j/k/Enter/Esc. WCAG 2.1 AA, keyboard 100%, color kèm icon+text. Desktop-first, tablet dùng được. Dark mode theo `prefers-color-scheme` + toggle.

## 9. Security FE

CSP strict, không token ở storage, không log sensitive, sanitize input (tránh `{@html}` không escape), validate client chỉ UX, HTTPS+HSTS, SRI nếu CDN ngoài, Dependabot/Snyk, không expose sourcemap prod. SQL highlight an toàn. Params dạng list, truncate string dài, binary → `[binary, N bytes]`.

## 10. Testing

- Unit Vitest: utils, schemas, hooks với MSW, components với testing-library/svelte. Target >70%.
- E2E Playwright: happy path login→approve→audit, expired, reject với lý do, token once. Critical paths + visual regression screenshots baseline.

## 11. Build & Deploy

- `svelte.config.js` adapter-static fallback `index.html` (SPA mode).
- Env: `VITE_API_BASE_URL, VITE_SSE_ENABLED, VITE_GRAFANA_URL, VITE_APP_NAME`.
- Dockerfile node:20 builder + nginx serve `build/`, SPA fallback, `/api/` proxy `janus-core:8788` + SSE (`proxy_buffering off`, `read_timeout 24h`), security headers. Options: Docker+nginx, Vercel/Netlify (CORS), S3+CloudFront, Go `embed.FS`.

## 12. Roadmap

- v0.1: SvelteKit + `@ohjanus/ui` setup, login local, approval queue+detail, audit cơ bản.
- v0.2: SSE, toast kit, dark mode, shortcuts, audit filter nâng cao.
- v0.3: tokens, connections, SSO OIDC.
- v0.4: dashboard flarecharts, export, link Grafana.
- v1.0: coverage, a11y audit, perf, docs.

## 13. DoD

Unit+integration happy/error path, E2E critical flow, screenshots `docs/screenshots/`, Lighthouse Perf>90 A11y>95 BP>90, không console warning, responsive 1280/1024/768, light+dark, loading/error/empty states, ESLint+Prettier+TS strict, PR review ≥1 người.

## 14. Tóm tắt

| Khía cạnh | Quyết định |
| :--- | :--- |
| Phạm vi | Tối giản: approval + audit |
| Repo/App | `ohjanus-ui` (Svelte 5 + SvelteKit) |
| UI kit | `@ohjanus/ui` — xem `ui-kit.md`, trang admin dùng trực tiếp |
| Backend | Admin API riêng :8788 |
| Auth | Session cookie + OIDC optional |
| Real-time | SSE typed end-to-end |
| Deploy | adapter-static + nginx / Vercel / S3 / Go embed |
| Ưu tiên | Approval > Audit > Tokens/Connections > Dashboard |
| Không làm | SQL editor, schema browser, config editor, Grafana replacement |
