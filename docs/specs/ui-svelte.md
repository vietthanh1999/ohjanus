# Đặc tả Kỹ thuật: **Janus UI** (Svelte Edition)

> **Codename**: `janus-ui`
> **Framework**: Svelte 5 + SvelteKit, tự build core component với Tailwind (không dùng UI kit ngoài)
> **Mục tiêu**: Giao diện web tối giản cho approval workflow và audit viewer, tận dụng hiệu năng và sự đơn giản của Svelte 5.
> **Repo**: riêng biệt với `janus` core. Liên quan: `docs/specs/ui.md` (bản React, tham khảo API contract đầy đủ), Admin API port 8788.
> **Ngày**: 2026-10-03
> **Quyết định**: chốt Svelte 5 + SvelteKit, KHÔNG dùng sv5ui/luna-plus/Bits UI — tự build Button, Input, Modal, Table, Badge, Toast... bằng Svelte + Tailwind để kiểm soát bundle và a11y.

## 0. Nguyên tắc thiết kế

1. **FE là client, không phải core**: chỉ gọi Admin API. Không business logic trong FE.
2. **Minimal surface area**: chỉ approval + audit + tokens + connections. Không SQL editor, schema browser.
3. **Runes-first**: dùng Svelte 5 runes (`$state`, `$derived`, `$effect`, `$props`), không store cũ.
4. **Type-safe end-to-end**: TypeScript strict, Zod validation, typed API client.
5. **Deploy riêng**: static SPA hoặc SvelteKit adapter-static.

## 1. Tech Stack đầy đủ

### 1.1. Core

| Thành phần | Lựa chọn | Lý do |
| :--- | :--- | :--- |
| Framework | Svelte 5 | Runes fine-grained, bundle nhỏ, hiệu năng cao |
| Meta-framework | SvelteKit | Routing, SSR optional, adapter-static cho SPA |
| Ngôn ngữ | TypeScript 5 (strict) | Type safety |
| Build | Vite 5 | Mặc định SvelteKit, HMR nhanh |

### 1.2. UI & Styling (tự build, không UI kit)

| Thành phần | Lựa chọn | Ghi chú |
| :--- | :--- | :--- |
| CSS | Tailwind CSS 4 | Utility-first duy nhất, không CSS file rời |
| Components | Tự build (`lib/components/ui/*.svelte`) | Button, Input, Textarea, Select, Checkbox, Badge, Card, Modal/Dialog, Drawer, Tabs, Table, Pagination, Skeleton, Toast, Tooltip, Countdown — Svelte 5 runes + Tailwind, a11y (focus-visible, aria-labels, keyboard Esc) |
| Icons | Lucide Svelte | Icon nhất quán, tree-shakeable |
| Toast | Tự build hoặc svelte-sonner | Ưu tiên tự build `Toast.svelte` với `$state` queue; chỉ dùng svelte-sonner nếu cần nhanh ở MVP |

> Không dùng sv5ui / luna-plus / Bits UI để tránh lock-in, giữ bundle nhỏ và chủ động a11y. Chỉ kéo thêm lib ngoài cho domain khó (SQL editor, charts, SSE).

### 1.3. Data & State

| Thành phần | Lựa chọn | Ghi chú |
| :--- | :--- | :--- |
| Server state | @tanstack/svelte-query v6 | Query/cache/mutation, runes native |
| Client state | Svelte 5 runes | `$state` UI state, không Zustand |
| URL state | SvelteKit `$page` | Filter/pagination query string |
| Form | svelte-form-hook | Giống RHF, Zod resolver |
| Validation | Zod | Schema + type inference |

### 1.4. Domain-specific

| Thành phần | Lựa chọn | Ghi chú |
| :--- | :--- | :--- |
| SQL highlight | @agnosticeng/editor | Svelte 5 wrap CodeMirror 6, PG/MySQL dialect, autocomplete |
| Charts | @faintshadow/flarecharts | Svelte 5 native, runes-first, SVG, a11y, composable |
| SSE | @sourceregistry/sveltekit-eventsource | Typed end-to-end, custom channels, debug hooks |
| Date | date-fns | Nhẹ, tree-shakeable |

### 1.5. Dev & Test

ESLint + typescript-eslint, Prettier + prettier-plugin-svelte, Vitest, @testing-library/svelte, Playwright, MSW.

### 1.6. Dependencies

```json
{
  "dependencies": {
    "@tanstack/svelte-query": "^6.0.0",
    "svelte-form-hook": "^1.1.8",
    "zod": "^3.23.0",
    "@agnosticeng/editor": "^0.0.6",
    "@faintshadow/flarecharts": "^26.3.1",
    "@sourceregistry/sveltekit-eventsource": "^1.1.2",
    "lucide-svelte": "^0.400.0",
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

## 2. Kiến trúc

```
Browser (SvelteKit SPA: Svelte 5 + Query + Tailwind, components tự build)
  │ HTTPS, session cookie HttpOnly Secure
  ▼
Reverse Proxy: / → static, /api/* → Admin API :8788, /stream → SSE (buffering off)
  ▼
Janus Core: MCP Server | Admin API | Approval Engine
```

Cấu trúc `janus-ui/`:

```
src/
  app.html, app.d.ts (App.Events SSE typing), app.css
  lib/api/client.ts (fetch credentials:include), auth.ts, approvals.ts, audit.ts, connections.ts, tokens.ts, types.ts
  lib/schemas/index.ts (Zod)
  lib/components/ui/Button,Input,Textarea,Select,Checkbox,Badge,Card,Modal,Drawer,Tabs,Table,Pagination,Skeleton,Toast,Tooltip,Countdown.svelte
  lib/components/layout/AppShell,Sidebar,TopBar.svelte
  lib/components/approval/ApprovalCard,ApprovalDetail,SqlPreview,ParamsList,ApprovalActions.svelte
  lib/components/audit/AuditTable,AuditFilters,AuditDetail.svelte
  lib/components/common/EmptyState,LoadingSkeleton,ErrorState.svelte
  lib/hooks/useAuth,useApprovals,useAudit,useTokens,useSSE.svelte.ts
  lib/utils/format.ts, sql.ts
  routes/+layout.svelte, +layout.ts (QueryClientProvider), +page.svelte (→/approvals),
    login/+page.svelte, approvals/+page.svelte, approvals/[id]/+page.svelte,
    audit/+page.svelte, connections/+page.svelte, tokens/+page.svelte, dashboard/+page.svelte
tests/unit,integration,e2e — static/, svelte.config.js, vite.config.ts, playwright.config.ts, .env.example
```

## 3. API Contract (tóm tắt — chi tiết xem `ui.md` §4)

- Base `/api/v1`, JSON, session cookie `janus_session`, CSRF header `X-Requested-With` cho mutating.
- Auth: `POST /auth/login`, `POST /auth/logout`, `GET /auth/me`, `GET /auth/oidc/authorize`, `GET /auth/oidc/callback`.
- Approvals: `GET /approvals?state,connection,requested_by,from,to,cursor,limit`, `GET /approvals/{id}`, `POST /{id}/approve`, `POST /{id}/reject`, `GET /approvals/stream` (SSE).
- Audit: `GET /audit`, `GET /audit/{id}`, `GET /audit/export`.
- Connections: `GET /connections`, `POST /connections/{name}/test` (không DSN).
- Tokens: `GET /tokens`, `POST /tokens` (trả 1 lần), `DELETE /tokens/{id}`.
- Dashboard: `GET /dashboard/summary`.
- Error: `{"error":{"code","message","request_id"}}`.

## 4. SSE Integration

`app.d.ts`:

```typescript
declare global {
  namespace App {
    interface Events {
      'approval.created': ApprovalItem;
      'approval.expired': { id: string };
      'approval.decided': { id: string; state: string };
    }
  }
}
```

Server forward SSE từ Janus Core (`src/routes/api/v1/approvals/stream/+server.ts` dùng `EventSource` server). Client hook `useApprovalSSE` (`useSSE.svelte.ts`): `onMount` mở `EventSource('/api/v1/approvals/stream')`, `on('approval.created/expired')` → `invalidateQueries(['approvals'])` + toast, `onDestroy` close.

## 5. Màn hình

- Login `/login`: Zod + svelte-form-hook, SSO redirect, rate limit 5/phút/IP, không lộ email tồn tại.
- Approval queue `/approvals`: tabs Pending/Approved/Rejected/Expired/All, card badge màu (SELECT xanh dương, INSERT xanh lá, UPDATE vàng, DELETE cam, DROP đỏ), SQL preview readonly (@agnosticeng/editor), affected estimate, countdown, SSE update, empty state.
- Approval detail `/approvals/[id]`: info chung, SQL readonly PG dialect, params list, affected + warnings, EXPLAIN collapsible, textarea lý do + Approve/Reject, countdown → disable, confirm dialog.
- Audit `/audit`: filter bar + table + drawer chi tiết, export 10k rows.
- Connections `/connections`: list + Test + drawer (schemas, denied tables, 24h count).
- Tokens `/tokens`: list + modal tạo (name, scopes, TTL) + hiển thị 1 lần + copy.
- Dashboard `/dashboard`: KPI cards, line chart flarecharts, top denials/connections, link Grafana.

Xem mock layout chi tiết ở `ui.md` §5.

## 6. Testing

- Unit Vitest: utils, schemas, hooks với MSW, components với testing-library/svelte. Target >70%.
- E2E Playwright: login→approve→audit happy path, expired, reject với lý do, token once. Critical paths. Ví dụ `approval-flow.spec.ts` xem bản gốc.

## 7. Build & Deploy

- `svelte.config.js` adapter-static fallback `index.html` (SPA mode).
- Env: `VITE_API_BASE_URL, VITE_SSE_ENABLED, VITE_GRAFANA_URL, VITE_APP_NAME`.
- Dockerfile node:20 builder + nginx serve `build/`, SPA fallback, `/api/` proxy `janus-core:8788` + SSE (`proxy_buffering off`, `read_timeout 24h`), security headers (CSP, nosniff, DENY frame).
- Xem `ui.md` §11 cho nginx đầy đủ và option Vercel/S3/Go embed.

## 8. Roadmap

- v0.1: SvelteKit+Tailwind setup + `lib/components/ui/*` tự build (Button/Input/Modal/Table/Badge/Toast...), login local, approval queue+detail, audit cơ bản.
- v0.2: SSE eventsource, toast tự build, dark mode, audit filter nâng cao.
- v0.3: tokens, connections, SSO OIDC.
- v0.4: dashboard flarecharts, export, link Grafana.
- v1.0: coverage, a11y, perf.

## 9. Tóm tắt

| Khía cạnh | Quyết định |
| :--- | :--- |
| Framework | Svelte 5 + SvelteKit |
| UI | Tự build với Tailwind (không UI kit) |
| Data | @tanstack/svelte-query v6 |
| Form | svelte-form-hook + Zod |
| SQL | @agnosticeng/editor (CodeMirror 6) |
| Charts | @faintshadow/flarecharts |
| SSE | @sourceregistry/sveltekit-eventsource |
| Styling | Tailwind 4 |
| State | Runes |
| Test | Vitest + Playwright + MSW |
| Deploy | adapter-static + nginx |

Stack runes + Tailwind tự chủ — bundle nhỏ, không lock-in UI kit.
