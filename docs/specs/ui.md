# Đặc tả Kỹ thuật: **Janus UI** — Frontend cho MCP Gateway

> **Codename**: `janus-ui`
> **Repo**: **riêng biệt** với `janus` core
> **Mục tiêu**: Giao diện web tối giản cho approval workflow và audit viewer.
> **Backend yêu cầu**: Admin API riêng port 8788 (không dùng chung MCP).
> **Ngày**: 2026-10-03

## 0. Nguyên tắc thiết kế

1. **FE là client, không phải core**: chỉ gọi REST API. Không business logic trong FE.
2. **Minimal surface area**: chỉ build màn hình thực sự cần.
3. **Read-only mặc định**: destructive action phải có confirmation + audit trail.
4. **Không lưu secret**: FE không thấy credentials, MCP token gốc, DSN. Chỉ metadata.
5. **Deploy riêng**: static site (Vercel/Netlify/S3), không nhồi vào Go binary trừ khi muốn.

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
Browser (React+TS+Vite+TanStack Query)
  │ HTTPS, session cookie HttpOnly Secure
  ▼
Reverse Proxy (Nginx/Caddy): / → static, /api/* → Janus Admin API
  ▼
Janus Core: MCP Server (agent) | Admin API :8788 (UI) | Approval Engine
```

- Admin API tách khỏi MCP server: khác port, khác auth (MCP opaque token, Admin session cookie/JWT).
- Không endpoint nào trả DSN/password/MCP token gốc.
- FE static, không SSR.

Cấu trúc `janus-ui/`:

```
src/
  main.tsx, App.tsx
  api/client.ts, auth.ts, approvals.ts, audit.ts, connections.ts, tokens.ts, types.ts
  components/ui/ (shadcn), layout/AppShell,Sidebar,TopBar,
    approval/ApprovalCard,ApprovalDetail,SqlPreview, audit/AuditRow,AuditFilters
  pages/Login,ApprovalQueue,ApprovalDetail,Audit,Connections,Tokens,Dashboard,NotFound
  hooks/useAuth,useApprovals,useAudit,useSSE.ts
  lib/format,sql,validation.ts
  stores/authStore,uiStore.ts
  styles/globals.css
tests/unit,integration,e2e (Playwright)
vite.config, tailwind.config, playwright.config, Dockerfile, README
```

## 3. Tech Stack

- Core: React 18, TypeScript 5, Vite 5, React Router v6, TanStack Query v5, Zustand (UI state), React Hook Form + Zod, fetch native.
- UI: Tailwind 3, shadcn/ui, Lucide, Recharts (dashboard), Shiki (SQL highlight), date-fns, Sonner (toast).
- Dev/Test: ESLint, Prettier, Vitest, Testing Library, Playwright, MSW.
- Không Next.js/Remix: không cần SSR/SEO/API routes, Vite SPA build nhanh, deploy static dễ.

## 4. API Contract (Admin API — yêu cầu Janus Core)

- Base `/api/v1`, JSON, auth cookie `janus_session` (HttpOnly, Secure, SameSite=Lax).

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

Form email+password, SSO button, validation, submit → login, success → `/approvals`, fail không lộ email tồn tại. Rate limit 5/phút/IP.

### 5.2. Approval Queue `/approvals`

Tabs Pending/Approved/Rejected/Expired/All, card hiển thị statement type màu (SELECT xanh dương, INSERT xanh lá, UPDATE vàng, DELETE cam, DROP đỏ), SQL preview 3 dòng, affected estimate, requester, countdown. Click → detail. SSE toast + badge. Empty: "Không có yêu cầu nào 🎉".

### 5.3. Approval Detail `/approvals/{id}`

Thông tin chung (connection, requester, request ID, tạo/hết hạn + countdown), SQL (Shiki), params list, ước lượng tác động + warnings, EXPLAIN plan collapsible, form quyết định (lý do optional khi approve, bắt buộc khi reject). Approve → confirm dialog. Hết hạn → disable. 409 → báo đã xử lý bởi người khác.

### 5.4. Audit `/audit`

Filters (event, connection, status, date, search debounce 300ms, shareable query string), table time/event/conn/tool/rows/ms, click row → drawer chi tiết (SQL, kết quả, policy). Export CSV/JSONL giới hạn 10k rows.

### 5.5. Connections `/connections`

Card: name, driver, status healthy/degraded, readonly, pool stats, last ping, nút Test + drawer (allowed schemas, denied tables, query 24h, link Grafana).

### 5.6. Tokens `/tokens`

List id/name/scopes/tạo/hết hạn/last used + Revoke (confirm). Modal tạo (name, scopes, TTL). Sau tạo: modal cảnh báo hiển thị 1 lần + Copy + checkbox "đã lưu".

### 5.7. Dashboard `/dashboard` (P2)

KPI cards, chart requests over time (Recharts), top denials, top connections, link Grafana.

## 6. Auth & Phân quyền

- Local: login → cookie → user info memory (không localStorage) → 401 redirect login.
- Roles: `viewer(audit:read,connection:read)`, `approver(+approval:*)`, `admin(+token:*,connection:*,config:read)`, `auditor(audit:read,export)`. FE ẩn/hiện UI theo permission, server luôn validate.
- Session: TTL 8h sliding, CSRF `SameSite=Lax` + header `X-Requested-With`, logout xóa session.

## 7. Real-time (SSE)

`EventSource('/api/v1/approvals/stream')` listen `approval.created/expired` → invalidate TanStack Query + toast. Auto-reconnect exponential backoff. Fallback polling 10s khi tab active nếu SSE bị chặn.

## 8. UX

Không giấu SQL/params/affected rows; confirm mọi destructive; không undo được nên confirm kỹ; loading/error/empty states rõ; error actionable; shortcuts j/k/Enter/Esc. WCAG 2.1 AA, keyboard 100%, color kèm icon+text. Desktop-first, tablet dùng được. Dark mode theo `prefers-color-scheme` + toggle.

## 9. Security FE

CSP strict, không token ở storage, không log sensitive, sanitize input (tránh `dangerouslySetInnerHTML`), validate client chỉ UX, HTTPS+HSTS, SRI nếu CDN ngoài, Dependabot/Snyk, không expose sourcemap prod. SQL highlight Shiki (không eval HTML thô). Params dạng list, truncate string dài, binary → `[binary, N bytes]`.

## 10. Testing

- Unit (Vitest): utils, hooks với MSW, components. Target >70%.
- Integration (Testing Library): login, approval approve, audit filter.
- E2E (Playwright): happy path login→approve→audit, expired, reject, token create-once. Visual regression screenshots baseline.

## 11. Build & Deploy

Env: `VITE_API_BASE_URL, VITE_SSE_ENABLED, VITE_GRAFANA_URL, VITE_APP_NAME, VITE_SENTRY_DSN`.
Build `npm run build → dist/`. Dockerfile node:20 builder + nginx serve, SPA fallback, `/api/` proxy về `janus-core:8788` với SSE (`proxy_buffering off`, `read_timeout 24h`), security headers. Options: Docker+nginx, Vercel/Netlify (CORS), S3+CloudFront, hoặc `embed.FS` serve từ Go binary.

## 12. Roadmap

- v0.1 MVP: login local, approval queue+detail, audit viewer cơ bản.
- v0.2: SSE, toast, dark mode, shortcuts, audit filter nâng cao.
- v0.3: tokens, connections, SSO OIDC.
- v0.4: dashboard, export, link Grafana.
- v1.0: coverage, a11y audit, perf, docs.

## 13. DoD

Unit+integration happy/error path, E2E critical flow, screenshots `docs/screenshots/`, Lighthouse Perf>90 A11y>95 BP>90, không console warning, responsive 1280/1024/768, light+dark, loading/error/empty states, ESLint+Prettier+TS strict, PR review ≥1 người.

## 14. Tóm tắt

| Khía cạnh | Quyết định |
| :--- | :--- |
| Phạm vi | Tối giản: approval + audit |
| Repo | Riêng `janus-ui` |
| Stack | React+TS+Vite+TanStack Query+shadcn/ui |
| Backend | Admin API riêng :8788 |
| Auth | Session cookie + OIDC optional |
| Real-time | SSE |
| Deploy | Docker+nginx / Vercel / S3 / Go embed |
| Ưu tiên | Approval > Audit > Tokens/Connections > Dashboard |
| Không làm | SQL editor, schema browser, config editor, Grafana replacement |
