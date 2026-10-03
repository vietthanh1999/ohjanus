# Kế hoạch Triển khai Janus (ohjanus)

> Nguồn: `docs/specs/init.md` (v0.1 → v1.0)
> Kiến trúc: `docs/architecture/01-ports-and-adapters.md` (ADR-001), `docs/architecture/02-libraries.md` (ADR-002)
> Ngôn ngữ: Go 1.22+
> Cập nhật: 2026-10-03

## 0. Nguyên tắc

- Bám sát Definition of Done (§18 spec): unit + integration test, audit event, metric, docs, không log credentials, pass `go vet` / `staticcheck` / `golangci-lint`, benchmark không regression >10%.
- Fail-closed: không resolve được tên bảng / AST lỗi → deny.
- Không dùng regex để validate SQL, bắt buộc AST.
- Credentials không bao giờ rời process, không xuất hiện trong log / response / error.
- Kiến trúc: Ports & Adapters tách in/out. Quy tắc: `adapter/in → port/in ← service → port/out ← adapter/out`, `domain` không import gì. Chi tiết xem ADR-001 §4.
- Anti-patterns: không entity có logic, không `UseCaseImpl`, không mapper thừa, `adapter/in` chỉ gọi `port/in` (không gọi service trực tiếp), `adapter/out` chỉ import `port/out`, wiring thủ công trong `main.go`. Xem ADR-001 §6.

---

## Phase 0: Bootstrap (1–2 ngày)

| # | Task | File / Package | DoD |
|---|------|----------------|-----|
| 0.1 | `go mod init`, Go 1.22, `Makefile`, `Dockerfile`, `.golangci.yml` + pin dependencies ADR-002 (xem bảng dưới) | `go.mod`, `Makefile`, `Dockerfile` | `go build ./...` pass, `go mod tidy` sạch |
| 0.2 | Dựng cấu trúc in/out (thay §2.1 spec cũ) | xem cây bên dưới | đúng cây ADR-001 §3 |
| 0.3 | CLI skeleton theo §14 spec (`cobra`) | `cmd/janus/main.go`, `internal/adapter/in/cli/*.adapter.go` | `janus serve/validate/token/connection/version --help` chạy được |
| 0.4 | Config mẫu đầy đủ (`koanf`, fallback `viper`) | `configs/janus.example.yaml` | khớp §3.1 spec |
| 0.5 | Định nghĩa `core/domain` + `core/port/in` + `core/port/out` rỗng | `internal/core/domain/*.domain.go`, `internal/core/port/in/*.usecase.go`, `internal/core/port/out/*.port.go` | `go vet` pass, compile được |
| 0.6 | Spike chốt 2 điểm mở ADR-002: `finemcp` vs `mcp-go`, `pg_query_go/v6` vs v5 | `internal/adapter/in/transport/stdio/`, `internal/adapter/out/validator/postgres/` | chọn 1, ghi kết quả vào ADR-002 |

Cây thư mục chuẩn (ADR-001 §3, suffix theo vai trò):

```
janus/
├── cmd/janus/main.go              # Wiring duy nhất
├── internal/
│   ├── core/
│   │   ├── domain/*.domain.go     # query, credential, policy, token, auth, audit
│   │   ├── port/in/*.usecase.go   # ReadUseCase, WriteUseCase, SchemaUseCase, HealthUseCase
│   │   ├── port/out/*.port.go     # Validator, PolicyEngine, Pool, AuditSink, TokenStore...
│   │   └── service/*.service.go   # implement port/in, phụ thuộc port/out
│   ├── adapter/
│   │   ├── in/transport/*/stdio.adapter.go,http.adapter.go,sse.adapter.go
│   │   ├── in/mcp/protocol.go,server.adapter.go,tools.adapter.go
│   │   ├── in/cli/cli.adapter.go
│   │   └── out/validator,connector,credential,policy,audit,approval,token,redact/*/*.adapter.go
│   └── config/config.go
```

CLI:
```
janus [command] [flags]
  serve         Chạy gateway (mặc định)
  validate      Validate config file
  token         Quản lý token (create, list, revoke)
  connection    Quản lý connection (add, list, test)
  version       In version
Flags: --config/-c, --log-level, --transport
```

### Dependencies chốt Phase 0 (từ ADR-002)

| Adapter | Thư viện pin ở `go.mod` | Dùng ở Phase |
| :--- | :--- | :--- |
| MCP framework | `finemcp/finemcp` (spike so với `mark3labs/mcp-go`) | 1.3, 2.1 |
| Validator PG | `pganalyze/pg_query_go/v6` (spike v6 vs v5) | 1.6 |
| Validator MySQL | `vitess.io/vitess/go/vt/sqlparser` | 2.6 |
| Connector PG | `jackc/pgx/v5` (+ `pgxpool`) | 1.5 |
| Connector MySQL | `go-sql-driver/mysql` | 2.6 |
| Connector SQLite | `modernc.org/sqlite` | 4 |
| Keyring | `zalando/go-keyring` | 2.7 |
| Vault | `hashicorp/vault-client-go` | 3 |
| Logging | `log/slog` stdlib (dự phòng `rs/zerolog`) | 1.10 |
| Config | `knadh/koanf` (fallback `spf13/viper`) | 1.1 |
| CLI | `spf13/cobra` | 0.3 |
| Test | `stretchr/testify`, `testcontainers-go` | 1.x, 2.x |
| Metrics | `prometheus/client_golang` | 3 |
| Tracing | `go.opentelemetry.io/otel` | 4 |

---

## Phase 1: v0.1 MVP — Read-only Gateway (trọng tâm)

Mục tiêu: Cursor / Claude Desktop SELECT an toàn qua stdio.

### 1.0 Core domain + ports in/out (làm trước)

- [ ] `core/domain/*.domain.go`: `Query`, `ResultSet`, `StatementType`, `ValidatedQuery`, `Credentials`, `PoolConfig`, `QueryOpts`, `PolicyDecision`/`Action`, `PreviewToken`, `AuditEvent`, error codes §12.1.
- [ ] `core/port/in/*.usecase.go`: `ReadUseCase`, `SchemaUseCase` (`ListConnections`, `GetSchema`), `HealthUseCase`. Chỉ import `domain` + `context`.
- [ ] `core/port/out/*.port.go`: `Validator`, `PolicyEngine`, `CredentialResolver`, `Connector`/`Pool`, `AuditSink`, `Redactor`. Chỉ import `domain` + `context`. Không import `port/in`.
- [ ] Quy ước: `var _ in.ReadUseCase = (*ReadService)(nil)` ở `*.service.go`, `var _ out.Validator = (*Validator)(nil)` ở `*.adapter.go`.
- [ ] Test: domain chỉ data holder, không logic.

### 1.1 Config (`internal/config/`)

- [ ] Load YAML + validate struct (server, auth, connections, policy, approval, credential, audit, redaction, observability, limits).
- [ ] ENV override: `JANUS_<SECTION>_<FIELD>`, `.` → `_`, uppercase.
- [ ] `janus validate --config` báo lỗi rõ ràng theo field.
- [ ] Test: thiếu field, sai driver, sai YAML, override ENV.

### 1.2 Auth (`core/port/out/*.port.go` + `internal/adapter/out/auth/*.adapter.go` + scope check ở `core/service/gateway.service.go`)

- [ ] Opaque token `jn_<random 32 bytes base62>`, lưu SHA-256 hash + `id`, `scopes`, `expires_at`.
- [ ] Flow: extract (adapter/in) → hash → lookup → check expiry → gắn scopes vào context → check scope per-tool ở `service/gateway.go`.
- [ ] Scopes v0.1: `read` → `db_list_connections`, `db_schema`, `db_read`, `db_explain`.
- [ ] Lỗi: `UNAUTHENTICATED`, `TOKEN_INVALID`, `TOKEN_EXPIRED`, `FORBIDDEN`.
- [ ] Test: token đúng/sai/hết hạn/thiếu scope, thu hồi. Mock ở tầng port/out.

### 1.3 Adapter/in: transport stdio + MCP (`internal/adapter/in/transport/stdio/stdio.adapter.go`, `internal/adapter/in/mcp/server.adapter.go,tools.adapter.go`)

- [ ] Implement `port` inbound: đọc JSON-RPC 2.0 từng dòng `stdin`, ghi `stdout`, log `stderr`.
- [ ] Methods: `initialize`, `initialized`, `tools/list` (lọc theo scope), `tools/call`, `ping`.
- [ ] Handshake `initialize` theo §4.3 spec. Error format `{code: -32000, message: "CODE: msg", data: {code, rule, request_id}}`.
- [ ] `mcp.Server` chỉ phụ thuộc `port/in` (`ReadUseCase`, `SchemaUseCase`), không import `core/service` trực tiếp.
- [ ] Test: handshake, list tools `read` vs `admin`, ping, message malformed. Mock `port/in`.

### 1.4 Credential (`core/port/out/credential.port.go` + `internal/adapter/out/credential/env/env.adapter.go`, `file/file.adapter.go`)

- [ ] Interface ở `core/port/out`, không cạnh implementation.
```go
type Resolver interface {
    Scheme() string
    Resolve(ctx context.Context, ref string) (domain.Credentials, error)
}
```
- [ ] v0.1: `env:`, `file:` (check chmod 0600 + owner). Chain resolver.
- [ ] Resolve 1 lần khi mở pool trong `main.go` wiring, không log DSN.
- [ ] Test: không log credentials, file sai permission bị từ chối.

### 1.5 DB Postgres (`core/port/out/connector.port.go` + `internal/adapter/out/connector/postgres/postgres.adapter.go`)

- [ ] Interfaces ở `core/port/out` theo §10.1 spec: `Connector`, `Pool` (`Query`, `Explain`, `Schema`, `Ping`, `Close`).
- [ ] Adapter: `pgx/v5`, pool per-alias, ping định kỳ. Chỉ map ở boundary driver → domain.
- [ ] Read-only: service truyền `QueryOpts{ReadOnly: true}` → adapter mở `SET TRANSACTION READ ONLY`.
- [ ] Test integration `testcontainers-go`, test service bằng fake Pool.

### 1.6 Validator AST (`core/port/out/validator.port.go` + `internal/adapter/out/validator/postgres/postgres.adapter.go`) — quan trọng nhất

- [ ] Adapter dùng `pg_query_go/v5`. Pipeline: Parse → statement type → walk AST (denied tables, banned functions, subquery DML/DDL, banned schema) → inject LIMIT → normalized SQL.
- [ ] Block functions: `pg_read_file`, `pg_read_binary_file`, `pg_ls_dir`, `lo_import`, `lo_export`, `pg_sleep`, `dblink`, `dblink_exec`, `COPY ... PROGRAM`, `pg_execute_server_program`.
- [ ] Resolve CTE/alias về bảng gốc, không resolve được → deny.
- [ ] Test: 100+ mẫu SELECT, CTE, subquery, comment injection, unicode, `WITH ... DELETE`, function cấm.

### 1.7 Policy (`core/port/out/policy.port.go` + `internal/adapter/out/policy/yaml/yaml.adapter.go`)

- [ ] Interface `Evaluate(ctx, ValidatedQuery) PolicyDecision` ở `port/out`.
- [ ] Rule YAML: `name`, `match {statement, connection, tables(glob), functions}`, `action`, `reason`. Đánh giá thứ tự, dừng match đầu tiên, `default_action: deny`.
- [ ] v0.1: `allow-select`, `deny-ddl`.
- [ ] Test rule matching, thứ tự ưu tiên, mock policy khi test service.

### 1.8 Services + Tools read-only (`core/service/read.service.go`, `schema.service.go`, `gateway.service.go` + `adapter/in/mcp/tools.adapter.go`)

- [ ] `ReadService` implement `in.ReadUseCase`, chỉ phụ thuộc `out.*`: validate → policy → lấy pool → `Query(ReadOnly:true)` → redact → audit.
- [ ] `SchemaService` implement `in.SchemaUseCase`: `ListConnections`, `GetSchema`, lọc `allowed_schemas`/`denied_tables`.
- [ ] Tools (adapter/in): `db_list_connections`, `db_schema`, `db_read`, `db_explain`. Row limit inject/rewrite, `truncated` flag. Limits: `max_query_length 10000`, `query_timeout 30s`, concurrency 20.
- [ ] Test service bằng fake `port/out`, test adapter/in bằng mock `port/in`. E2E: handshake → list → schema → read OK, DROP chặn.

### 1.9 Audit (`core/port/out/audit.port.go` + `internal/adapter/out/audit/file/file.adapter.go`, `stdout/stdout.adapter.go`, `multi/multi.adapter.go`)

- [ ] Interface `AuditSink.Emit(ctx, AuditEvent)` ở `port/out`.
- [ ] Adapter file JSONL + rotate, stdout, multi fan-out.
- [ ] Events v0.1: `session.initialized`, `auth.success/failure`, `query.validated/denied/executed/failed`. Không log `params` raw.
- [ ] Test: mỗi service call đều Emit.

### 1.10 Observability + Redact + Wiring

- [ ] `log_level`, mọi log qua redact (`adapter/out/redact/regex/regex.adapter.go` implement `port/out/redact.port.go`).
- [ ] `GET /healthz`, `GET /readyz` qua `in.HealthUseCase`.
- [ ] `cmd/janus/main.go` wiring: tạo `adapter/out/*` → inject vào `service.*` (thỏa `port/in`) → inject `port/in` vào `adapter/in/mcp.NewServer(...)`.

### 1.11 Exit criteria v0.1

- [ ] `janus serve` + Cursor qua stdio: SELECT được, DROP bị chặn.
- [ ] `go vet`, `staticcheck`, `golangci-lint` pass.
- [ ] Không có import sai chiều: `service → adapter`, `adapter/out → port/in`, `adapter/in → service` (check import graph).

---

## Phase 2: v0.2 — Write + Remote (2–3 tuần)

| # | Task | Port in/out / Adapter in/out | Ghi chú spec |
|---|------|------------------------------|--------------|
| 2.1 | Transport HTTP/SSE | `adapter/in/transport/http/http.adapter.go`, `sse/sse.adapter.go` (gọi `port/in`) | `POST /mcp`, `GET /mcp/sse`, `Bearer`, 401/403 |
| 2.2 | Approval + single-use token | `port/out/approval.port.go`, `token.port.go` + `adapter/out/approval/cli/cli.adapter.go`, `adapter/out/token/memory/memory.adapter.go` | `pvw_*`, TTL 5m, `pending/approved/used`, `TokenStore` interface để sau thay Redis |
| 2.3 | `db_write_preview` | `port/in/write.usecase.go: Preview()` + `service/write.service.go` | EXPLAIN estimate, không thực thi |
| 2.4 | `db_write_execute` | `port/in/write.usecase.go: Execute()` + `service/write.service.go` | Check `TOKEN_*`, prompt stderr `Approve? [y/N]`, tx commit/rollback |
| 2.5 | Policy `warn-write` | `adapter/out/policy/yaml/yaml.adapter.go` | `INSERT/UPDATE/DELETE/MERGE` → `require_approval` |
| 2.6 | MySQL | `adapter/out/connector/mysql/mysql.adapter.go` + `adapter/out/validator/mysql/mysql.adapter.go` | `go-sql-driver/mysql`, parser tidb/vitess |
| 2.7 | Keyring | `adapter/out/credential/keyring/keyring.adapter.go` | `keyring://...` |
| 2.8 | Audit bổ sung | `adapter/out/audit/*/*.adapter.go` | `write.preview_created/approved/rejected/executed`, `token.reused` |
| 2.9 | HTTP approval | `adapter/out/approval/http/http.adapter.go` (outbound, dù là HTTP) | `POST /approve/{token}` |
| 2.10 | Admin API approvals (phục vụ UI §4.3 `ui.md`) | `port/in/approval.usecase.go` + `service/approval.service.go` + `adapter/in/admin/approval.adapter.go` | `GET /api/v1/approvals`, `GET /{id}`, `POST /{id}/approve`, `POST /{id}/reject`, port :8788 tách khỏi MCP, auth session cookie |
| 2.11 | Admin API audit read (phục vụ UI §4.4) | `port/in/audit-query.usecase.go` + `adapter/in/admin/audit.adapter.go` | `GET /api/v1/audit`, `GET /{id}`, filter + cursor pagination, không trả params raw trừ `audit:read:params` |
| 2.12 | Admin SSE stream (phục vụ UI §7) | `adapter/in/admin/stream.adapter.go` | `GET /api/v1/approvals/stream` SSE `approval.created/expired`, fallback polling |

Test negative: DROP chặn, reuse token chặn, sửa SQL sau preview chặn, readonly + allow-write → từ chối. Admin API: 401 không cookie, 403 thiếu permission, 409 approve trùng, 410 hết hạn.

---

## Phase 3: v0.3 — Hardening + Admin quản trị (UI P1)

- [ ] Vault (`adapter/out/credential/vault/vault.adapter.go`): AppRole.
- [ ] Redaction (`port/out/redact.port.go` + `adapter/out/redact/regex/regex.adapter.go`): email/card/phone + cột nhạy cảm → `[REDACTED]`. Áp dụng cả Admin API responses.
- [ ] Metrics: counters/histograms, `/metrics`. Nằm ở `service/gateway.service.go` (middleware) + adapter/out observability.
- [ ] Rate limiting (token bucket 60 req/m → `RATE_LIMITED`), concurrency (semaphore 20), timeout (30s → `QUERY_TIMEOUT`). Áp dụng riêng cho MCP port và Admin :8788.
- [ ] Error codes full §12.1 spec + error format Admin `{"error":{code,message,request_id}}` (§4.8 `ui.md`).
- [ ] Admin auth local (UI §6.1, P0 trễ): `POST /api/v1/auth/login/logout`, `GET /me`, cookie `janus_session` HttpOnly/Secure/SameSite=Lax, roles `viewer/approver/admin/auditor` — `adapter/in/admin/auth.adapter.go` + `port/in/auth.usecase.go`.
- [ ] Admin tokens (UI §4.6, P1): `GET/POST/DELETE /api/v1/tokens`, token chỉ trả 1 lần — `adapter/in/admin/token.adapter.go`.
- [ ] Admin connections (UI §4.5, P1): `GET /api/v1/connections`, `POST /{name}/test`, không trả DSN — `adapter/in/admin/connection.adapter.go`.

---

## Phase 4: v0.4 — Scale + Admin insights (UI P2)

- [ ] Tracing OTel: span per request (adapter/in → service → adapter/out), bao gồm cả Admin API.
- [ ] Token Redis (`adapter/out/token/redis/redis.adapter.go`): thay memory, không đổi `port/out`.
- [ ] Admin dashboard (UI §4.7, P2): `GET /api/v1/dashboard/summary` — `adapter/in/admin/dashboard.adapter.go`.
- [ ] Admin audit export (UI §4.4): `GET /api/v1/audit/export` CSV/JSONL, rate limit, cần `audit:export`.
- [ ] Admin SSO OIDC (UI §6.2, P1 trễ): `/auth/oidc/authorize`, `/callback` — `adapter/in/admin/oidc.adapter.go`.
- [ ] SQLite (`adapter/out/connector/sqlite/sqlite.adapter.go` + `adapter/out/validator/sqlite/sqlite.adapter.go`).
- [ ] Config hot reload + audit `config.reloaded`. Scopes `write_preview`, `write_execute`, `admin`.

---

## Phase 5: v1.0 — Stable

- [ ] Coverage >80% (ưu tiên `core/service`).
- [ ] Fuzz validator, injection vectors, privilege escalation.
- [ ] Benchmark AST vs regex, throughput 1/10/100.
- [ ] Docker distroless, systemd, Helm. Docs site.
- [ ] ADR spec §19 chốt.

---

## Phụ lục A: Thứ tự triển khai (in/out)

```
0. domain + port/in + port/out interfaces
→ config
→ adapter/out/auth + service/gateway scope check
→ adapter/out/credential env/file
→ adapter/out/connector postgres
→ adapter/out/validator postgres
→ adapter/out/policy yaml
→ service/read + service/schema + service/gateway (implement port/in)
→ adapter/in/transport stdio + adapter/in/mcp
→ adapter/out/audit + redact
→ wiring main.go → integration-test
→ adapter/in http/sse (MCP remote)
→ port/out approval/token + adapter/out memory/cli + service/write (port/in/write)
→ adapter/in/admin approvals + audit read + SSE stream (UI P0, :8788 tách khỏi MCP)
→ adapter/out mysql + keyring
→ vault/redact/metrics/limits + adapter/in/admin auth/tokens/connections (UI P1)
→ tracing/redis/dashboard/export/oidc/sqlite/reload (UI P2)
```

Wiring: `main.go` tạo adapter/out → inject service (port/in) → inject adapter/in. Không service nào tự `new` adapter, không adapter/in nào gọi service trực tiếp.

## Phụ lục B: Rủi ro

1. `pg_query_go` cgo — build Docker phức tạp.
2. Mỗi driver một AST — đã tách `adapter/out/validator/<driver>` từ Phase 1.
3. Token memory → Redis — đã có `port/out.TokenStore` từ Phase 2.
4. Read-only tx là phòng vệ cuối — `ReadService` luôn `ReadOnly: true`.
5. Lẫn in/out (VD: để approval/cli ở adapter/in, để transport ở adapter/out) — review import graph mỗi PR.
