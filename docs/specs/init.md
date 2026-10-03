# Đặc tả Kỹ thuật: **Janus** — MCP Gateway bảo mật cho Database

> **Codename**: `ohjanus`
> **Ngôn ngữ**: Go 1.22+
> **Mục tiêu**: Gateway trung gian cho phép AI agent truy cập database qua MCP mà không bao giờ tiếp xúc với credentials thật, đồng thời kiểm soát chặt chẽ mọi truy vấn.

---

## 1. Tổng quan

### 1.1. Bài toán

AI agent (Cursor, Claude Desktop, VS Code, custom agent) cần truy vấn database để hỗ trợ người dùng, nhưng:
- Không được biết username/password/host thật của DB.
- Không được phép chạy các lệnh nguy hiểm (DROP, DELETE không kiểm soát).
- Mọi hành động phải được audit.
- Thao tác ghi phải có sự phê duyệt của con người.

### 1.2. Giải pháp

Janus là một **MCP Server** viết bằng Go, đóng vai trò gateway:
- Nhận yêu cầu từ MCP Client qua **stdio** (local) hoặc **HTTP/SSE** (remote).
- Xác thực bằng MCP Token.
- Validate SQL bằng AST parser (không dùng regex).
- Inject credentials từ vault (OS Keychain / file / env / HashiCorp Vault).
- Thực thi truy vấn qua connection pool.
- Ghi audit log có cấu trúc.
- Với thao tác ghi: yêu cầu preview + approval + single-use token.

### 1.3. Non-goals (những gì Janus KHÔNG làm)

- Không phải ORM hay query builder.
- Không tự sinh SQL từ ngôn ngữ tự nhiên (đó là việc của agent).
- Không thay thế quyền hạn DB — **vẫn yêu cầu DB user read-only riêng cho agent**.
- Không cache kết quả truy vấn (trừ khi bật tường minh).

---

## 2. Kiến trúc tổng thể

```
┌─────────────────────────────────────────────────────────────────┐
│                         MCP Client                              │
│              (Cursor / Claude Desktop / Custom)                 │
└──────────────────────────┬──────────────────────────────────────┘
                           │ stdio hoặc HTTP/SSE
                           │ kèm MCP Token
                           ▼
┌─────────────────────────────────────────────────────────────────┐
│                          JANUS                                  │
│                                                                 │
│  ┌──────────────┐   ┌──────────────┐   ┌──────────────────┐   │
│  │  Transport   │──▶│    Auth      │──▶│  Tool Router     │   │
│  │  Layer       │   │  (Token)     │   │  (dispatch)      │   │
│  └──────────────┘   └──────────────┘   └────────┬─────────┘   │
│                                                  │             │
│                       ┌──────────────────────────┼──────────┐  │
│                       ▼                          ▼          ▼  │
│              ┌────────────────┐        ┌──────────────┐ ┌────┐│
│              │ Query Validator│        │  Approval    │ │Audit│
│              │ (AST-based)    │        │  Engine      │ │Log ││
│              └───────┬────────┘        └──────┬───────┘ └────┘│
│                      │                        │               │
│                      ▼                        ▼               │
│              ┌────────────────┐        ┌──────────────┐       │
│              │  Policy Engine │        │Token Issuer  │       │
│              │ (allow/deny)   │        │(single-use)  │       │
│              └───────┬────────┘        └──────┬───────┘       │
│                      │                        │               │
│                      └────────┬───────────────┘               │
│                               ▼                               │
│                      ┌────────────────┐                       │
│                      │  Credential    │                       │
│                      │  Resolver      │                       │
│                      │(Keychain/File/ │                       │
│                      │ Vault/Env)     │                       │
│                      └───────┬────────┘                       │
│                              ▼                                │
│                      ┌────────────────┐                       │
│                      │  DB Connector  │                       │
│                      │  Pool + Driver │                       │
│                      └───────┬────────┘                       │
└──────────────────────────────┼────────────────────────────────┘
                               ▼
                        ┌─────────────┐
                        │  Database   │
                        └─────────────┘
```

### 2.1. Cấu trúc thư mục dự án

```
janus/
├── cmd/
│   └── janus/
│       └── main.go                 # Entry point
├── internal/
│   ├── config/                     # Load & validate config
│   ├── transport/                  # stdio, http, sse
│   ├── auth/                       # Token verification
│   ├── mcp/                        # MCP protocol implementation
│   ├── tools/                      # Tool definitions & handlers
│   ├── validator/                  # AST SQL validation
│   ├── policy/                     # Allow/deny rules engine
│   ├── approval/                   # Write approval workflow
│   ├── token/                      # Single-use token issuer
│   ├── credential/                 # Credential resolvers
│   │   ├── keyring/
│   │   ├── file/
│   │   ├── env/
│   │   └── vault/
│   ├── db/                         # DB connectors
│   │   ├── postgres/
│   │   ├── mysql/
│   │   └── sqlite/
│   ├── audit/                      # Audit logger
│   ├── redact/                     # PII/secret redaction
│   └── observability/              # Metrics, tracing, health
├── pkg/
│   └── api/                        # Public types (nếu cần export)
├── configs/
│   └── janus.example.yaml
├── docs/
├── scripts/
├── go.mod
├── go.sum
├── Dockerfile
└── Makefile
```

---

## 3. Cấu hình (Configuration)

### 3.1. Định dạng

Hỗ trợ **YAML** là chính, **env override** cho mọi field. Ví dụ:

```yaml
# configs/janus.yaml
server:
  name: janus
  version: 0.1.0
  transport: stdio            # stdio | http | sse
  http:
    listen: "127.0.0.1:8787"
    tls:
      enabled: false
      cert_file: ""
      key_file: ""

auth:
  mode: token                 # token | none (chỉ dùng cho dev)
  tokens:
    - id: "tok_dev_001"
      hash: "sha256:abc123..."   # KHÔNG lưu plaintext token
      scopes: ["read"]
      expires_at: "2026-12-31T23:59:59Z"
    - id: "tok_dev_002"
      hash: "sha256:def456..."
      scopes: ["read", "write_preview", "write_execute"]
      expires_at: "2026-12-31T23:59:59Z"

connections:
  - name: "analytics"         # alias agent nhìn thấy
    driver: postgres          # postgres | mysql | sqlite
    dsn_ref: "keyring://janus/analytics"   # KHÔNG chứa password
    pool:
      max_open: 10
      max_idle: 5
      conn_max_lifetime: "1h"
    readonly: true            # bắt buộc mở transaction read-only
    allowed_schemas: ["public", "analytics"]
    denied_tables: ["users", "secrets", "api_keys"]
    allowed_tables: []        # rỗng = tất cả trừ denied
    row_limit: 1000           # cứng, không thể vượt

policy:
  default_action: deny        # deny | allow
  rules:
    - name: "allow-select"
      match:
        statement: ["SELECT", "WITH", "EXPLAIN", "SHOW"]
      action: allow
    - name: "deny-ddl"
      match:
        statement: ["CREATE", "ALTER", "DROP", "TRUNCATE", "GRANT", "REVOKE"]
      action: deny
      reason: "DDL bị chặn hoàn toàn"
    - name: "warn-write"
      match:
        statement: ["INSERT", "UPDATE", "DELETE", "MERGE"]
      action: require_approval

approval:
  enabled: true
  method: "cli"               # cli | http | webhook
  cli:
    prompt: true
  token_ttl: "5m"             # single-use token hết hạn sau 5 phút

credential:
  resolvers:
    - scheme: "keyring"
      service_prefix: "janus"
    - scheme: "file"
      base_path: "/etc/janus/secrets"
      file_mode: 0600
    - scheme: "env"
      prefix: "JANUS_"
    - scheme: "vault"
      address: "https://vault.internal:8200"
      auth_method: "approle"
      role_id_ref: "env:VAULT_ROLE_ID"
      secret_id_ref: "env:VAULT_SECRET_ID"

audit:
  enabled: true
  sinks:
    - type: "file"
      path: "/var/log/janus/audit.jsonl"
      format: "jsonl"
      rotate:
        max_size_mb: 100
        max_backups: 10
    - type: "stdout"
      format: "json"

redaction:
  enabled: true
  patterns:
    - name: "email"
      regex: "[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}"
      replacement: "[EMAIL]"
    - name: "credit_card"
      regex: "\\b(?:\\d[ -]*?){13,16}\\b"
      replacement: "[CARD]"
    - name: "vn_phone"
      regex: "\\b(0|\\+84)\\d{9,10}\\b"
      replacement: "[PHONE]"
  columns:                    # redact theo tên cột
    - "password"
    - "password_hash"
    - "api_key"
    - "secret"
    - "token"

observability:
  log_level: "info"           # debug | info | warn | error
  metrics:
    enabled: true
    listen: "127.0.0.1:9090"
    path: "/metrics"
  tracing:
    enabled: false
    otlp_endpoint: ""

limits:
  max_query_length: 10000     # bytes
  query_timeout: "30s"
  max_concurrent_queries: 20
  rate_limit:
    enabled: true
    requests_per_minute: 60
```

### 3.2. Quy tắc override bằng ENV

Mọi field có thể override bằng `JANUS_<SECTION>_<FIELD>` (dấu `.` → `_`, uppercase). Ví dụ:
- `JANUS_SERVER_TRANSPORT=http`
- `JANUS_AUTH_MODE=none`
- `JANUS_OBSERVABILITY_LOG_LEVEL=debug`

---

## 4. MCP Interface (Lớp giao tiếp)

### 4.1. Transport

#### 4.1.1. stdio (mặc định, cho local)

- Đọc JSON-RPC messages từ `stdin`, ghi ra `stdout`.
- Log đi ra `stderr` (không được lẫn vào stdout).
- Mỗi message là một dòng JSON (JSON-RPC 2.0).

#### 4.1.2. HTTP + SSE

- `POST /mcp` — nhận JSON-RPC request, trả về response (hoặc stream qua SSE).
- `GET /mcp/sse` — mở SSE stream cho server-initiated messages.
- Header bắt buộc: `Authorization: Bearer <MCP_TOKEN>`.
- Trả `401` nếu token thiếu/sai, `403` nếu token hết hạn hoặc thiếu scope.

### 4.2. MCP Methods hỗ trợ

| Method | Bắt buộc | Mô tả |
| :--- | :--- | :--- |
| `initialize` | ✅ | Handshake, trả về capabilities |
| `initialized` | ✅ | Notification từ client |
| `tools/list` | ✅ | Liệt kê các tool có sẵn (lọc theo scope của token) |
| `tools/call` | ✅ | Gọi tool |
| `ping` | ✅ | Health check |
| `resources/list` | ⚪ | (Optional) liệt kê resources — có thể dùng để expose schema |
| `prompts/list` | ❌ | Không hỗ trợ ở v1 |

### 4.3. Handshake `initialize`

**Request:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "method": "initialize",
  "params": {
    "protocolVersion": "2024-11-05",
    "capabilities": { "roots": { "listChanged": true } },
    "clientInfo": { "name": "cursor", "version": "0.42.0" }
  }
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "protocolVersion": "2024-11-05",
    "capabilities": {
      "tools": { "listChanged": false },
      "logging": {}
    },
    "serverInfo": { "name": "janus", "version": "0.1.0" }
  }
}
```

---

## 5. Tools (Định nghĩa chi tiết)

### 5.1. Danh sách tools

| Tool | Scope yêu cầu | Mô tả |
| :--- | :--- | :--- |
| `db_list_connections` | `read` | Liệt kê các connection alias agent được phép thấy |
| `db_schema` | `read` | Liệt kê schema, table, column của một connection |
| `db_read` | `read` | Chạy SELECT/WITH/EXPLAIN/SHOW |
| `db_write_preview` | `write_preview` | Sinh preview + token cho lệnh ghi |
| `db_write_execute` | `write_execute` | Thực thi lệnh ghi với single-use token |
| `db_explain` | `read` | Chạy EXPLAIN (không thực thi) |

### 5.2. Tool: `db_list_connections`

**Input schema:**
```json
{
  "type": "object",
  "properties": {},
  "additionalProperties": false
}
```

**Output:**
```json
{
  "connections": [
    { "name": "analytics", "driver": "postgres", "readonly": true }
  ]
}
```

### 5.3. Tool: `db_schema`

**Input:**
```json
{
  "type": "object",
  "properties": {
    "connection": { "type": "string" },
    "schema": { "type": "string" },
    "table": { "type": "string" }
  },
  "required": ["connection"]
}
```

**Output:**
```json
{
  "schemas": [
    {
      "name": "public",
      "tables": [
        {
          "name": "orders",
          "columns": [
            { "name": "id", "type": "bigint", "nullable": false },
            { "name": "total", "type": "numeric", "nullable": true }
          ],
          "primary_key": ["id"],
          "foreign_keys": []
        }
      ]
    }
  ]
}
```

**Ràng buộc:**
- Chỉ trả về schema/table nằm trong `allowed_schemas` và không nằm trong `denied_tables`.
- Nếu agent yêu cầu schema bị cấm → trả lỗi `SCHEMA_NOT_ALLOWED`.

### 5.4. Tool: `db_read`

**Input:**
```json
{
  "type": "object",
  "properties": {
    "connection": { "type": "string" },
    "sql": { "type": "string" },
    "params": {
      "type": "array",
      "items": {}
    },
    "limit": {
      "type": "integer",
      "minimum": 1,
      "maximum": 1000
    }
  },
  "required": ["connection", "sql"]
}
```

**Output (thành công):**
```json
{
  "columns": ["id", "total"],
  "rows": [[1, 100.5], [2, 200.0]],
  "row_count": 2,
  "truncated": false,
  "duration_ms": 12
}
```

**Output (bị chặn):**
```json
{
  "error": {
    "code": "QUERY_DENIED",
    "message": "Statement type UPDATE is not allowed by policy",
    "rule": "deny-ddl"
  }
}
```

### 5.5. Tool: `db_write_preview`

**Input:**
```json
{
  "type": "object",
  "properties": {
    "connection": { "type": "string" },
    "sql": { "type": "string" },
    "params": { "type": "array" }
  },
  "required": ["connection", "sql"]
}
```

**Output:**
```json
{
  "preview_token": "pvw_01HXYZ...",
  "expires_at": "2026-10-03T10:35:00Z",
  "statement_type": "UPDATE",
  "affected_estimate": 42,
  "plan": "...",
  "warnings": ["This will update 42 rows in table 'orders'"],
  "sql_normalized": "UPDATE orders SET status = $1 WHERE created_at < $2"
}
```

**Hành vi:**
- Chạy `EXPLAIN` để ước lượng.
- Không thực thi.
- Sinh `preview_token` gắn với `connection`, `sql`, `params` qua SHA-256.
- Token TTL = `approval.token_ttl` (mặc định 5 phút).

### 5.6. Tool: `db_write_execute`

**Input:**
```json
{
  "type": "object",
  "properties": {
    "connection": { "type": "string" },
    "sql": { "type": "string" },
    "params": { "type": "array" },
    "preview_token": { "type": "string" }
  },
  "required": ["connection", "sql", "params", "preview_token"]
}
```

**Output:**
```json
{
  "rows_affected": 42,
  "duration_ms": 87
}
```

**Hành vi:**
- Xác thực `preview_token` khớp với `connection + sql + params` (SHA-256).
- Nếu token đã dùng → `TOKEN_ALREADY_USED`.
- Nếu token hết hạn → `TOKEN_EXPIRED`.
- Nếu SQL bị sửa → `TOKEN_MISMATCH`.
- Nếu không khớp → từ chối.
- Nếu `approval.method=cli`, in prompt ra stderr và chờ y/n trước khi thực thi.
- Đánh dấu token đã dùng (in-memory + persist nếu multi-instance).
- Mở transaction, thực thi, commit/rollback.

### 5.7. Tool: `db_explain`

Tương tự `db_read` nhưng chỉ chạy `EXPLAIN <sql>`, không bao giờ thực thi.

---

## 6. Security Layer

### 6.1. Authentication

#### 6.1.1. Token format

Khuyến nghị dùng **opaque token** (không phải JWT) để dễ thu hồi:
- Format: `jn_<random 32 bytes base62>`
- Lưu trong config dưới dạng SHA-256 hash + metadata.

#### 6.1.2. Token verification flow

```
1. Extract token từ header/params
2. Hash token → tra cứu trong token store
3. Nếu không tìm thấy → 401 UNAUTHENTICATED
4. Nếu hết hạn → 401 TOKEN_EXPIRED
5. Lấy scopes → gắn vào context
6. Kiểm tra scope cho tool sắp gọi → 403 nếu thiếu
```

#### 6.1.3. Scopes

| Scope | Cho phép |
| :--- | :--- |
| `read` | `db_list_connections`, `db_schema`, `db_read`, `db_explain` |
| `write_preview` | `db_write_preview` |
| `write_execute` | `db_write_execute` |
| `admin` | Tất cả + reload config |

### 6.2. Query Validation (AST-based)

**Đây là lớp quan trọng nhất. Không dùng regex.**

#### 6.2.1. Thư viện

- **PostgreSQL**: dùng `github.com/pganalyze/pg_query_go/v5` (bindings của libpg_query).
- **MySQL**: dùng `github.com/pingcap/tidb/pkg/parser` hoặc `vitess.io/vitess/go/vt/sqlparser`.
- **SQLite**: dùng `modernc.org/sqlite` parser hoặc `github.com/auxten/postgresql-parser` (không lý tưởng).

#### 6.2.2. Validation pipeline

```
SQL raw
  │
  ▼
[1] Parse thành AST
  │  └─ lỗi cú pháp → PARSE_ERROR
  ▼
[2] Xác định statement type (top-level)
  │  └─ SELECT/WITH/EXPLAIN/SHOW → read path
  │  └─ INSERT/UPDATE/DELETE/MERGE → write path (cần approval)
  │  └─ DDL khác → DENY ngay
  ▼
[3] Walk AST để kiểm tra:
  │  - Có truy cập bảng bị denied?
  │  - Có gọi function nguy hiểm (pg_read_file, lo_import...)?
  │  - Có subquery chứa DDL/DML?
  │  - Có tham chiếu schema bị cấm?
  ▼
[4] Nếu là read: inject LIMIT nếu chưa có
  │  └─ Nếu đã có LIMIT > row_limit → rewrite xuống row_limit
  ▼
[5] Trả về normalized SQL + metadata
```

#### 6.2.3. Danh sách function bị chặn (PostgreSQL)

```
pg_read_file, pg_read_binary_file, pg_ls_dir,
lo_import, lo_export, pg_sleep (chống DoS),
dblink, dblink_exec,
COPY (khi dùng với PROGRAM),
pg_execute_server_program
```

#### 6.2.4. Table/Schema allow-deny

- Kiểm tra ở cấp AST, không phải string match.
- Nếu query dùng alias hoặc CTE → vẫn phải resolve được tên bảng gốc.
- Nếu không resolve được (dynamic SQL) → **deny**.

### 6.3. Policy Engine

Rule-based, đánh giá theo thứ tự, dừng ở rule match đầu tiên.

**Ngôn ngữ rule:**
```yaml
- name: "string"
  match:
    statement: [list of types]
    connection: [list of names]
    tables: [list of patterns]      # glob
    functions: [list of patterns]
  action: allow | deny | require_approval
  reason: "string"
```

**Kết quả:**
- `allow` → tiếp tục pipeline.
- `deny` → trả lỗi `QUERY_DENIED`, không chạm DB.
- `require_approval` → chuyển sang nhánh write preview.

### 6.4. Read-only Enforcement

Với connection có `readonly: true`:
- Mở transaction với `SET TRANSACTION READ ONLY` (PostgreSQL) / `START TRANSACTION READ ONLY` (MySQL).
- Nếu policy cho phép write nhưng connection readonly → từ chối.

### 6.5. Row Limiting

- Mọi SELECT đều bị inject `LIMIT row_limit` nếu chưa có.
- Nếu có LIMIT lớn hơn → rewrite.
- Response có field `truncated: true` nếu bị cắt.

### 6.6. Rate Limiting

- Token bucket per token ID.
- Mặc định: 60 req/phút.
- Vượt → trả lỗi `RATE_LIMITED`.

### 6.7. Concurrency Limit

- Semaphore giới hạn số query đồng thời (mặc định 20).
- Vượt → `TOO_MANY_CONCURRENT_QUERIES`.

### 6.8. Query Timeout

- Context với timeout `query_timeout` (mặc định 30s).
- Hết hạn → cancel query, trả `QUERY_TIMEOUT`.

---

## 7. Credential Management

### 7.1. Resolver Interface

```go
type Resolver interface {
    Scheme() string
    Resolve(ctx context.Context, ref string) (Credentials, error)
}

type Credentials struct {
    DSN      string  // đã đầy đủ, dùng để mở connection
    Username string
    Password string
    Host     string
    Port     int
    Database string
    // ...
}
```

### 7.2. Các resolver hỗ trợ

| Scheme | Ví dụ ref | Backend |
| :--- | :--- | :--- |
| `env` | `env:JANUS_ANALYTICS_DSN` | Biến môi trường |
| `file` | `file:/etc/janus/secrets/analytics.dsn` | File, chmod 0600 |
| `keyring` | `keyring://janus/analytics` | OS Keychain (macOS/Linux/Windows) |
| `vault` | `vault://secret/data/janus/analytics#dsn` | HashiCorp Vault |

### 7.3. Quy tắc bảo mật

- **Không bao giờ** log credentials.
- **Không bao giờ** trả credentials trong MCP response.
- **Không bao giờ** để credentials trong error message.
- Resolver chỉ được gọi **một lần** khi mở connection, sau đó connection được giữ trong pool.
- Nếu dùng file: kiểm tra `file_mode` là `0600` và owner là user chạy janus.

### 7.4. Redaction trong log

Mọi log line phải đi qua `redact.Apply()`:
- Che password trong DSN: `postgres://user:***@host/db`
- Che token trong header.
- Che pattern PII (email, phone, card) nếu bật.

---

## 8. Approval Engine

### 8.1. Luồng approval (CLI mode)

```
Agent gọi db_write_preview
   │
   ▼
Janus chạy EXPLAIN, sinh preview_token
   │
   ▼
Trả preview cho agent + in ra stderr:
   ┌─────────────────────────────────────┐
   │ ⚠️  WRITE REQUEST                   │
   │ Connection: analytics               │
   │ Statement: UPDATE                   │
   │ SQL: UPDATE orders SET status=$1... │
   │ Params: ['shipped', '2026-01-01']   │
   │ Affected: ~42 rows                  │
   │                                     │
   │ Approve? [y/N]                      │
   └─────────────────────────────────────┘
   │
   ▼ (người dùng gõ y)
Đánh dấu token là "approved"
   │
   ▼
Agent gọi db_write_execute với preview_token
   │
   ▼
Janus kiểm tra token approved + chưa dùng
   │
   ▼
Thực thi, mark token used, trả kết quả
```

### 8.2. Token store

- In-memory map với TTL.
- Nếu chạy multi-instance → dùng Redis hoặc SQLite chung.
- Mỗi token lưu: `id`, `hash`, `connection`, `sql_hash`, `params_hash`, `created_at`, `expires_at`, `state` (`pending`/`approved`/`used`), `approved_by`.

### 8.3. HTTP approval mode (optional)

- Endpoint `POST /approve/{token}` với body `{ "decision": "approve" | "reject", "approver": "..." }`.
- Cần auth riêng (basic auth hoặc token admin).

---

## 9. Audit Logging

### 9.1. Format

Mỗi event là một JSON object một dòng (JSONL):

```json
{
  "ts": "2026-10-03T10:30:00.123Z",
  "event": "query.executed",
  "request_id": "req_01HXYZ",
  "token_id": "tok_dev_001",
  "client": { "name": "cursor", "version": "0.42.0" },
  "connection": "analytics",
  "tool": "db_read",
  "sql_hash": "sha256:...",
  "sql_normalized": "SELECT id, total FROM orders WHERE created_at > $1 LIMIT 1000",
  "statement_type": "SELECT",
  "tables": ["orders"],
  "policy_decision": "allow",
  "policy_rule": "allow-select",
  "row_count": 42,
  "truncated": false,
  "duration_ms": 12,
  "status": "success",
  "error": null
}
```

### 9.2. Events

| Event | Khi nào |
| :--- | :--- |
| `session.initialized` | Client handshake xong |
| `auth.success` / `auth.failure` | Xác thực token |
| `query.validated` | AST validation xong |
| `query.denied` | Policy từ chối |
| `query.executed` | Query chạy xong |
| `query.failed` | Query lỗi |
| `write.preview_created` | Sinh preview token |
| `write.approved` / `write.rejected` | Người dùng quyết định |
| `write.executed` | Ghi thành công |
| `token.reused` | Cố dùng lại token |
| `config.reloaded` | Reload config |

### 9.3. Redaction trong audit

- `sql_normalized` được giữ nguyên (không chứa params).
- `params` **không** được log (có thể chứa PII).
- Nếu cần log params → hash chúng.

### 9.4. Rotation

- File sink hỗ trợ rotate theo size.
- Giữ tối đa `max_backups`.

---

## 10. Database Connector

### 10.1. Interface

```go
type Connector interface {
    Driver() string
    Open(ctx context.Context, creds Credentials, pool PoolConfig) (Pool, error)
}

type Pool interface {
    Query(ctx context.Context, sql string, params []any, opts QueryOpts) (ResultSet, error)
    Exec(ctx context.Context, sql string, params []any) (ExecResult, error)
    Explain(ctx context.Context, sql string, params []any) (Plan, error)
    Schema(ctx context.Context, schema string) ([]Table, error)
    Close() error
}

type QueryOpts struct {
    ReadOnly bool
    Timeout  time.Duration
    RowLimit int
}
```

### 10.2. Drivers hỗ trợ (v1)

| Driver | Thư viện |
| :--- | :--- |
| PostgreSQL | `github.com/jackc/pgx/v5` |
| MySQL | `github.com/go-sql-driver/mysql` |
| SQLite | `modernc.org/sqlite` (pure Go) |

### 10.3. Connection pool

- Mỗi connection alias có một pool riêng.
- Cấu hình `max_open`, `max_idle`, `conn_max_lifetime`.
- Health check định kỳ (ping).

### 10.4. Read-only transaction

```go
tx, _ := pool.BeginTx(ctx, &sql.TxOptions{ReadOnly: true})
defer tx.Rollback()
// ... query ...
tx.Commit()  // no-op với read-only
```

---

## 11. Redaction (PII / Secrets)

### 11.1. Áp dụng ở đâu

- **Log**: mọi log line.
- **Audit**: field `params`, `error`.
- **Response**: cột có tên trong `redaction.columns`.
- **Error messages**: che DSN, token.

### 11.2. Cách hoạt động

```go
func Apply(input string, cfg RedactionConfig) string {
    out := input
    for _, p := range cfg.Patterns {
        out = p.Regex.ReplaceAllString(out, p.Replacement)
    }
    return out
}
```

### 11.3. Redact theo cột

Khi build ResultSet, nếu tên cột match `redaction.columns` → thay giá trị bằng `"[REDACTED]"`.

---

## 12. Error Handling

### 12.1. Error codes

| Code | HTTP | Mô tả |
| :--- | :--- | :--- |
| `UNAUTHENTICATED` | 401 | Thiếu token |
| `TOKEN_INVALID` | 401 | Token sai |
| `TOKEN_EXPIRED` | 401 | Token hết hạn |
| `FORBIDDEN` | 403 | Thiếu scope |
| `CONNECTION_NOT_FOUND` | 404 | Alias không tồn tại |
| `SCHEMA_NOT_ALLOWED` | 403 | Schema bị cấm |
| `TABLE_NOT_ALLOWED` | 403 | Bảng bị cấm |
| `PARSE_ERROR` | 400 | SQL không parse được |
| `QUERY_DENIED` | 403 | Policy từ chối |
| `QUERY_TOO_LONG` | 400 | Vượt max_query_length |
| `QUERY_TIMEOUT` | 408 | Hết thời gian |
| `RATE_LIMITED` | 429 | Vượt rate limit |
| `TOO_MANY_CONCURRENT_QUERIES` | 429 | Vượt concurrency |
| `TOKEN_MISMATCH` | 400 | preview_token không khớp |
| `TOKEN_ALREADY_USED` | 400 | Token đã dùng |
| `APPROVAL_REQUIRED` | 428 | Cần approval trước |
| `DB_ERROR` | 500 | Lỗi từ DB |
| `INTERNAL` | 500 | Lỗi nội bộ |

### 12.2. Error response format (MCP)

```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "error": {
    "code": -32000,
    "message": "QUERY_DENIED: Statement type DROP is not allowed",
    "data": {
      "code": "QUERY_DENIED",
      "rule": "deny-ddl",
      "request_id": "req_01HXYZ"
    }
  }
}
```

---

## 13. Observability

### 13.1. Metrics (Prometheus)

| Metric | Type | Labels |
| :--- | :--- | :--- |
| `janus_requests_total` | counter | `tool`, `status` |
| `janus_query_duration_seconds` | histogram | `connection`, `statement_type` |
| `janus_policy_denials_total` | counter | `rule`, `connection` |
| `janus_active_connections` | gauge | `connection` |
| `janus_approval_pending` | gauge | — |
| `janus_audit_events_total` | counter | `event` |

### 13.2. Tracing (OpenTelemetry)

- Span cho mỗi request MCP.
- Child span cho: auth, validate, policy, resolve credential, query, audit.
- Export OTLP nếu bật.

### 13.3. Health check

- `GET /healthz` → 200 nếu process sống.
- `GET /readyz` → 200 nếu tất cả connection pool ping được.

---

## 14. CLI

```
janus [command] [flags]

Commands:
  serve         Chạy gateway (mặc định)
  validate      Validate config file
  token         Quản lý token (create, list, revoke)
  connection    Quản lý connection (add, list, test)
  version       In version

Flags:
  --config, -c  Đường dẫn config (default: ./janus.yaml)
  --log-level   debug|info|warn|error
  --transport   stdio|http|sse

Ví dụ:
  janus serve --config /etc/janus/janus.yaml
  janus token create --scopes read,write_preview --ttl 720h
  janus connection test analytics
```

---

## 15. Deployment

### 15.1. Binary

```bash
go build -o janus ./cmd/janus
```

Binary tĩnh, không cần runtime.

### 15.2. Docker

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -o /janus ./cmd/janus

FROM gcr.io/distroless/static-debian12
COPY --from=builder /janus /janus
COPY configs/janus.yaml /etc/janus/janus.yaml
USER nonroot:nonroot
ENTRYPOINT ["/janus", "serve", "--config", "/etc/janus/janus.yaml"]
```

### 15.3. systemd

```ini
[Unit]
Description=Janus MCP Gateway
After=network.target

[Service]
Type=simple
User=janus
ExecStart=/usr/local/bin/janus serve --config /etc/janus/janus.yaml
Restart=on-failure
Environment=JANUS_OBSERVABILITY_LOG_LEVEL=info

[Install]
WantedBy=multi-user.target
```

---

## 16. Testing

### 16.1. Unit tests

- `validator`: test AST parsing với hàng trăm mẫu SQL (SELECT, CTE, subquery, comment injection, unicode).
- `policy`: test rule matching.
- `token`: test TTL, single-use.
- `redact`: test pattern matching.
- `credential`: test từng resolver.

### 16.2. Integration tests

- Dùng `testcontainers-go` để spin up Postgres/MySQL.
- Test end-to-end: handshake → read → write preview → approve → execute.
- Test negative: DROP bị chặn, token reuse bị chặn, schema bị cấm.

### 16.3. Security tests

- Fuzz test AST validator với SQL random.
- Test SQL injection vectors.
- Test privilege escalation.

### 16.4. Benchmark

- So sánh latency khi dùng AST vs regex.
- So sánh throughput với 1, 10, 100 concurrent clients.

---

## 17. Roadmap

### v0.1 (MVP)
- [x] stdio transport
- [x] Token auth
- [x] PostgreSQL connector
- [x] AST validation (Postgres)
- [x] `db_read`, `db_schema`, `db_list_connections`
- [x] File audit log
- [x] Config YAML

### v0.2
- [ ] HTTP/SSE transport
- [ ] `db_write_preview` + `db_write_execute`
- [ ] CLI approval
- [ ] MySQL connector
- [ ] Keyring resolver

### v0.3
- [ ] Vault resolver
- [ ] Redaction engine
- [ ] Prometheus metrics
- [ ] Rate limiting
- [ ] Row limit injection

### v0.4
- [ ] OpenTelemetry tracing
- [ ] Multi-instance token store (Redis)
- [ ] Web UI cho approval
- [ ] SQLite connector
- [ ] Config hot reload

### v1.0
- [ ] Stable API
- [ ] Full test coverage > 80%
- [ ] Docker + Helm chart
- [ ] Documentation site

---

## 18. Tiêu chí hoàn thành (Definition of Done)

Một feature được coi là xong khi:
1. Có unit test + integration test.
2. Có audit event tương ứng.
3. Có metric tương ứng.
4. Có tài liệu trong `docs/`.
5. Không log credentials (kiểm tra bằng test).
6. Pass `go vet`, `staticcheck`, `golangci-lint`.
7. Benchmark không regression > 10%.

---

## 19. Quyết định thiết kế quan trọng (ADR tóm tắt)

| # | Quyết định | Lý do |
| :--- | :--- | :--- |
| 1 | Dùng AST parser thay vì regex | Regex dễ bypass bằng comment, encoding |
| 2 | Opaque token thay vì JWT | Dễ thu hồi, không cần verify signature mỗi request |
| 3 | Single-use token cho write | Chống replay, buộc approval từng lần |
| 4 | Credentials không bao giờ rời process | Giảm blast radius khi agent bị compromise |
| 5 | Read-only transaction bắt buộc | Lớp phòng vệ cuối cùng nếu validator có lỗi |
| 6 | Audit log JSONL | Dễ parse, dễ ship vào ELK/Loki |
| 7 | Go + binary tĩnh | Triển khai đơn giản, hiệu năng cao |

---

Đây là spec đầy đủ cho v0.1 → v1.0. Bạn muốn tôi đi sâu vào phần nào tiếp theo? Ví dụ:

- **Code skeleton** cho `cmd/janus/main.go` + `internal/mcp/`
- **AST validator** chi tiết với `pg_query_go`
- **Approval flow** với CLI prompt
- **Test plan** cụ thể cho security
- **Schema database** cho token store (nếu multi-instance)