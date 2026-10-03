# ADR-001: Ports & Adapters thay vì Clean Architecture đầy đủ

> Trạng thái: Accepted (sửa đổi: tách in/out)
> Ngày: 2026-10-03
> Liên quan: `docs/specs/init.md` §2, `docs/plans/implementation-plan.md`

**Quyết định ngắn: Không dùng full Clean Architecture. Dùng Ports & Adapters (Hexagonal) với tách rõ inbound (in / driving) và outbound (out / driven).**

## 1. Tại sao full Clean Architecture là quá nặng với Janus

Clean Architecture (Uncle Bob) được thiết kế cho ứng dụng có **domain logic phức tạp** — ngân hàng, ERP, y tế. Nó giả định:

- Tầng **Entities** với business rules phong phú.
- Tầng **Use Cases** điều phối entities.
- Framework, DB, UI chỉ là chi tiết thay thế được.

Nhưng Janus **không có domain logic kiểu đó**. "Domain" của Janus thực chất là:

| Thành phần | Bản chất thật |
| :--- | :--- |
| SQL Validator | Kỹ thuật (parse AST) |
| Policy Engine | Cấu hình + matching |
| Token Store | Kỹ thuật (TTL, hash) |
| Approval Engine | Workflow đơn giản |
| Credential Resolver | Kỹ thuật (I/O) |
| DB Connector | Kỹ thuật (I/O) |

Đây là **infrastructure concerns**, không phải business domain. Cố nhồi vào "Entities" sẽ tạo ra:

- **Boilerplate vô nghĩa**: `SQLStatement` entity, `PolicyDecision` entity... chỉ là data holder.
- **Indirection thừa**: thao tác đơn giản đi qua 3-4 lớp interface.
- **Mapping hell**: `dbModel ↔ domainModel ↔ dtoModel` cho struct gần giống hệt nhau.
- **Test khó hơn**: phải mock nhiều lớp hơn cho cùng một logic.

Triết lý Go: "a little copying is better than a little dependency", "clear is better than clever".

## 2. Ports & Adapters in/out là gì và tại sao phù hợp

Hexagonal chuẩn chia hai chiều:

- **Port in (inbound / driving / primary)**: Interface mà bên ngoài gọi vào. Do `core/service` implement. VD: `ReadUseCase`, `WriteUseCase`, `SchemaUseCase`.
- **Port out (outbound / driven / secondary)**: Interface mà core gọi ra ngoài. Do `adapter/out` implement. VD: `Validator`, `PolicyEngine`, `Pool`, `AuditSink`.
- **Adapter in (driving)**: Nhận request từ thế giới ngoài, gọi port in. VD: stdio, HTTP/SSE, MCP server, CLI.
- **Adapter out (driven)**: Bị core gọi, đi ra hạ tầng. VD: postgres, keyring, vault, file audit, redis token.

```
        ┌─ adapter/in ─┐     ┌─ core ──┐     ┌─ adapter/out ─┐
MCP ──▶ │ stdio/http   │ ──▶ │ port/in │ ──▶ │ service ──▶ │ port/out │ ──▶ │ postgres/vault/audit │
Client  │ mcp/server   │     │ use case│     │             │          │     │ ...                    │
        └──────────────┘     └─────────┘     └──────────────┘     └────────────────┘
```

Phù hợp với Janus vì:

1. Janus **sinh ra đã là boundary** — một chiều nhận MCP, một chiều chạm DB/vault/audit.
2. Không tách in/out sẽ lẫn lộn: transport lẫn với DB connector trong cùng `adapter/`, use case lẫn với validator trong cùng `port/`.
3. Tách ra giúp test đúng chỗ: test service bằng fake `port/out`, test `adapter/in` bằng mock `port/in`.
4. Thay transport (stdio → HTTP) không chạm service; thay DB (postgres → mysql) không chạm inbound.

## 3. Cấu trúc thư mục + quy ước suffix theo vai trò

> Suffix giúp nhìn tên file biết ngay vai trò, không cần mở file. In/out phân biệt bằng **thư mục** (`port/in`, `port/out`, `adapter/in`, `adapter/out`), không lặp lại trong suffix.

| Vai trò | Thư mục | Suffix | Ví dụ |
| :--- | :--- | :--- | :--- |
| Domain | `core/domain/` | `*.domain.go` | `query.domain.go` |
| Port in (use case) | `core/port/in/` | `*.usecase.go` | `read.usecase.go` |
| Port out | `core/port/out/` | `*.port.go` | `validator.port.go` |
| Service | `core/service/` | `*.service.go` | `read.service.go` |
| Adapter in/out | `adapter/in/`, `adapter/out/` | `*.adapter.go` | `stdio.adapter.go`, `postgres.adapter.go` |
| Config/MCP protocol | `config/`, `adapter/in/mcp/` | không suffix, giữ `*.go` | `config.go`, `protocol.go` |

```
janus/
├── cmd/
│   └── janus/
│       └── main.go                        # Wiring duy nhất
├── internal/
│   ├── core/
│   │   ├── domain/                        # Struct thuần, không logic
│   │   │   ├── query.domain.go            # Query, ResultSet, ValidatedQuery, StatementType
│   │   │   ├── credential.domain.go
│   │   │   ├── policy.domain.go           # PolicyDecision, Action
│   │   │   ├── token.domain.go            # PreviewToken, TokenState
│   │   │   ├── auth.domain.go             # Token, Scopes
│   │   │   └── audit.domain.go            # AuditEvent
│   │   ├── port/
│   │   │   ├── in/                        # Driven by outside, implemented by service
│   │   │   │   ├── read.usecase.go        # ReadUseCase
│   │   │   │   ├── write.usecase.go       # WriteUseCase (Preview + Execute)
│   │   │   │   ├── schema.usecase.go      # SchemaUseCase, ListConnectionsUseCase
│   │   │   │   └── health.usecase.go      # HealthUseCase
│   │   │   └── out/                       # Driven by core, implemented by adapter/out
│   │   │       ├── validator.port.go      # Validator
│   │   │       ├── policy.port.go         # PolicyEngine
│   │   │       ├── credential.port.go     # Resolver, Chain
│   │   │       ├── connector.port.go      # Connector, Pool
│   │   │       ├── audit.port.go          # AuditSink
│   │   │       ├── approval.port.go       # ApprovalEngine
│   │   │       ├── token.port.go          # TokenStore
│   │   │       ├── redact.port.go         # Redactor
│   │   │       └── clock.port.go          # Clock (test TTL dễ)
│   │   └── service/                       # Implement port/in, chỉ phụ thuộc port/out
│   │       ├── read.service.go
│   │       ├── write.service.go
│   │       ├── schema.service.go
│   │       └── gateway.service.go         # Scope check + rate limit + concurrency + timeout
│   ├── adapter/
│   │   ├── in/                            # Driving: gọi port/in
│   │   │   ├── transport/
│   │   │   │   ├── stdio/stdio.adapter.go
│   │   │   │   ├── http/http.adapter.go
│   │   │   │   └── sse/sse.adapter.go
│   │   │   ├── mcp/                       # protocol, server, tools registry
│   │   │   │   ├── protocol.go
│   │   │   │   ├── server.adapter.go
│   │   │   │   └── tools.adapter.go
│   │   │   └── cli/cli.adapter.go         # janus serve/validate/token/connection
│   │   └── out/                           # Driven: implement port/out
│   │       ├── validator/postgres/postgres.adapter.go
│   │       ├── validator/mysql/mysql.adapter.go
│   │       ├── validator/sqlite/sqlite.adapter.go
│   │       ├── connector/postgres/postgres.adapter.go
│   │       ├── connector/mysql/mysql.adapter.go
│   │       ├── connector/sqlite/sqlite.adapter.go
│   │       ├── credential/env/env.adapter.go
│   │       ├── credential/file/file.adapter.go
│   │       ├── credential/keyring/keyring.adapter.go
│   │       ├── credential/vault/vault.adapter.go
│   │       ├── policy/yaml/yaml.adapter.go
│   │       ├── audit/file/file.adapter.go
│   │       ├── audit/stdout/stdout.adapter.go
│   │       ├── audit/multi/multi.adapter.go
│   │       ├── approval/cli/cli.adapter.go
│   │       ├── approval/http/http.adapter.go
│   │       ├── token/memory/memory.adapter.go
│   │       ├── token/redis/redis.adapter.go
│   │       └── redact/regex/regex.adapter.go
│   └── config/config.go
```

Phân loại chi tiết:

| Nhóm | Đặt ở đâu | Gồm gì |
| :--- | :--- | :--- |
| port in | `core/port/in` | `ReadUseCase`, `WriteUseCase`, `SchemaUseCase`, `HealthUseCase` |
| port out | `core/port/out` | `Validator`, `PolicyEngine`, `CredentialResolver`, `Connector`/`Pool`, `AuditSink`, `TokenStore`, `ApprovalEngine`, `Redactor` |
| adapter in | `adapter/in` | `transport/stdio,http,sse`, `mcp/*`, `cli/*` |
| adapter out | `adapter/out` | `validator/*`, `connector/*`, `credential/*`, `policy/yaml`, `audit/*`, `approval/*`, `token/*`, `redact/*` |
| service | `core/service` | implement port in, phụ thuộc port out |

Lưu ý: `approval/cli` tuy prompt qua stderr (trông giống inbound) nhưng bản chất là **outbound**: service gọi ra để xin quyết định con người, nên đặt ở `adapter/out/approval/cli`. Tương tự `approval/http`.

## 4. Quy tắc phụ thuộc (in/out)

```
adapter/in ──▶ port/in ◀── service ──▶ port/out ◀── adapter/out
                  │           │              │
                  ▼           ▼              ▼
                 domain ◀─────┴──────────────┘
                 (không import gì)
```

- `core/domain`: không import gì từ internal.
- `core/port/in`: chỉ import `domain` + `context`. Không import `port/out`.
- `core/port/out`: chỉ import `domain` + `context`. Không import `port/in`.
- `core/service`: import `port/in` (để implement) + `port/out` (để gọi) + `domain`. **Không import adapter.**
- `adapter/in`: import `port/in` (+ `domain` cho DTO boundary). **Không import service trực tiếp**, chỉ qua interface. Không import `adapter/out`.
- `adapter/out`: import `port/out` + `domain`. Không import `port/in`, `service`, `adapter/in`.
- `cmd/janus`: wiring duy nhất, được import tất cả.

Check nhanh bằng import graph: nếu thấy `core/service → adapter/...` hoặc `adapter/out → port/in` là sai.

## 5. Ví dụ: Read flow (đã tách in/out)

```go
// core/port/in/read.usecase.go — inbound, service implement
package in

type ReadUseCase interface {
    Read(ctx context.Context, req domain.ReadRequest) (*domain.ResultSet, error)
}

// core/port/out/validator.port.go — outbound, adapter/out implement
package out

type Validator interface {
    Validate(ctx context.Context, conn string, sql string) (*domain.ValidatedQuery, error)
}

// core/port/out/connector.port.go
package out

type Pool interface {
    Query(ctx context.Context, q domain.Query, opts domain.QueryOpts) (*domain.ResultSet, error)
    Explain(ctx context.Context, q domain.Query) (*domain.Plan, error)
    Schema(ctx context.Context, schema string) ([]domain.Table, error)
    Ping(ctx context.Context) error
    Close() error
}
```

```go
// core/service/read.service.go — implement port/in, phụ thuộc port/out
package service

import (
    "janus/internal/core/domain"
    "janus/internal/core/port/in"
    "janus/internal/core/port/out"
)

type ReadService struct {
    validator out.Validator
    policy    out.PolicyEngine
    pools     map[string]out.Pool
    audit     out.AuditSink
}

var _ in.ReadUseCase = (*ReadService)(nil)

func (s *ReadService) Read(ctx context.Context, req domain.ReadRequest) (*domain.ResultSet, error) {
    vq, err := s.validator.Validate(ctx, req.Connection, req.SQL)
    if err != nil { /* audit query.denied */ return nil, err }
    d := s.policy.Evaluate(ctx, vq)
    if d.Action == domain.ActionDeny { /* audit */ return nil, domain.ErrQueryDenied(d.Reason) }
    pool, ok := s.pools[req.Connection]
    if !ok { return nil, domain.ErrConnectionNotFound(req.Connection) }
    return pool.Query(ctx, vq.Query, domain.QueryOpts{ReadOnly: true, RowLimit: 1000})
}
```

```go
// adapter/out/validator/postgres/postgres.adapter.go
package postgres

import (
    "janus/internal/core/domain"
    "janus/internal/core/port/out"
)

var _ out.Validator = (*Validator)(nil)
// Validate dùng pg_query_go, walk AST, check denied tables/functions...

// adapter/in/mcp/server.adapter.go
package mcp

import "janus/internal/core/port/in"

type Server struct {
    read   in.ReadUseCase
    write  in.WriteUseCase
    schema in.SchemaUseCase
}
// tools/call → s.read.Read(ctx, ...) — chỉ biết port/in, không biết service.
```

Wiring trong `cmd/janus/main.go`: tạo `adapter/out/*` → inject vào `service.NewReadService(...)` (thỏa `port/in`) → inject `port/in` vào `adapter/in/mcp.NewServer(...)`.

## 6. Anti-patterns cần tránh

| ❌ Không nên | ✅ Nên |
| :--- | :--- |
| `domain.SQLStatement` entity với method `Validate()` | `out.Validator` xử lý, domain chỉ data holder |
| `usecase/ReadUseCaseImpl` | `service.ReadService` implement `in.ReadUseCase`, không thêm interface cho chính nó |
| Mapper giữa các layer giống nhau | Dùng chung domain, chỉ map ở boundary driver |
| Đặt interface cạnh implementation | `in`/`out` ở `core/port`, consumer sở hữu |
| `service` import `adapter/...` | Chỉ import `port/out` |
| `adapter/out` import `port/in` | Chỉ import `port/out` |
| `adapter/in` gọi `service` struct trực tiếp | Chỉ gọi qua `port/in` interface |
| DTO riêng từng layer | Dùng `domain`, chỉ DTO ở `adapter/in` boundary |
| DI framework | Wiring thủ công trong `main.go` |

## 7. Khi nào mới cần full Clean Architecture?

Khi Janus thành nền tảng governance phức tạp:

- Multi-tenant, policy/quota/billing riêng.
- Approval đa cấp, escalation.
- Compliance rules (VD: không truy vấn dữ liệu EU từ IP ngoài EU).
- Rich domain model với state machine.

Lúc đó tách thêm `usecase` + `entity`. **Đừng làm sớm.**

## 8. Tóm tắt

| Khía cạnh | Quyết định |
| :--- | :--- |
| Kiến trúc | Ports & Adapters với tách in/out |
| port in | `core/port/in/*.usecase.go`: use cases, service implement |
| port out | `core/port/out/*.port.go`: validator, policy, pool, audit, token... adapter/out implement |
| adapter in/out | `adapter/in/**/*.adapter.go`, `adapter/out/**/*.adapter.go` — gọi/implement ports |
| domain | `core/domain/*.domain.go`: struct thuần |
| service | `core/service/*.service.go`: implement port/in |
| Dependency | `adapter/in → port/in ← service → port/out ← adapter/out` |
| Wiring | Thủ công trong `main.go` |
| Testing | Service: fake port/out. adapter/in: mock port/in |

---

## Action tiếp theo

- [x] Cập nhật `docs/plans/implementation-plan.md` theo cấu trúc in/out + suffix này.
- [ ] Viết skeleton `core/port/in/*.usecase.go`, `core/port/out/*.port.go` + 1 `*.in.adapter.go` mẫu + 1 `*.out.adapter.go` mẫu + `main.go` wiring.
