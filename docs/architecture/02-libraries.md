# ADR-002: Lựa chọn thư viện Go cho Janus

> Trạng thái: Proposed
> Ngày: 2026-10-03
> Liên quan: ADR-001 (`01-ports-and-adapters.md`), spec `docs/specs/init.md`
> Nguyên tắc: ổn định, hiệu năng, swappability của từng adapter (in/out).

## 1. MCP Protocol & Transport (adapter/in)

| Thư viện | Vai trò | Ghi chú |
| :--- | :--- | :--- |
| `finemcp/finemcp` | Framework MCP chính | Production-grade, 16 middleware, circuit breaker, OpenTelemetry, zero-network test harness giúp test logic MCP không cần server thật. |
| `mark3labs/mcp-go` | Thay thế | Phổ biến, có client+server, hỗ trợ stdio và SSE (đang phát triển). API đơn giản hơn. |

> Khuyến nghị: bắt đầu với `finemcp` vì có sẵn resilience + observability theo spec.

## 2. SQL Validation AST-based (adapter/out/validator)

Lớp bảo mật quan trọng nhất, cần parser chính xác từng hệ CSDL.

| Database | Thư viện | Ghi chú |
| :--- | :--- | :--- |
| PostgreSQL | `pganalyze/pg_query_go/v6` | Dùng chính parser của PostgreSQL server, chính xác nhất. Chuẩn cho phân tích SQL Postgres trong Go. Lưu ý: spec cũ ghi v5, chốt lên v6. |
| MySQL | `vitess.io/vitess/go/vt/sqlparser` | Bao phủ cú pháp MySQL rộng (kể cả `SPATIAL`), binary nhỏ hơn so với tidb parser. |
| SQLite | `modernc.org/sqlite` | Driver pure-Go (không CGO). Không phải parser độc lập — cần tự viết walker trên AST driver cung cấp, hoặc chấp nhận bảo vệ mức thấp hơn cho SQLite. |

## 3. Database Drivers & Pool (adapter/out/connector)

| Driver | Thư viện | Ghi chú |
| :--- | :--- | :--- |
| PostgreSQL | `jackc/pgx/v5` | Pure Go, hiệu năng cao, hỗ trợ tính năng đặc thù Postgres. Kèm `pgxpool`. |
| MySQL | `go-sql-driver/mysql` | Chuẩn, ổn định, dùng rộng rãi. |
| SQLite | `modernc.org/sqlite` | Pure-Go, build binary đơn giản, không CGO. |

## 4. Credential Resolution (adapter/out/credential)

Thiết kế theo `port/out/credential.port.go` để thay thế lẫn nhau.

| Backend | Thư viện | Ghi chú |
| :--- | :--- | :--- |
| OS Keychain | `zalando/go-keyring` | Cross-platform: macOS Keychain, Linux Secret Service/D-Bus, Windows Credential Manager. |
| HashiCorp Vault | `hashicorp/vault-client-go` | Client chính thức HashiCorp, full API đọc/ghi/list secret. |
| File & Env | stdlib (`os`, `io`) | Không cần lib ngoài. Bắt buộc check chmod 0600, không log secret. |

## 5. Logging & Audit (adapter/out/audit)

| Thư viện | Vai trò | Ghi chú |
| :--- | :--- | :--- |
| `log/slog` (stdlib, Go 1.21+) | Structured logger chính | JSON/Text, hiệu năng tốt, không dependency ngoài. Đủ cho hầu hết trường hợp. |
| `rs/zerolog` | High-perf logger (dự phòng) | Zero-allocation JSON. Chỉ dùng nếu đo được `slog` thành bottleneck. |

> Khuyến nghị: bắt đầu `log/slog`.

## 6. Configuration & CLI (adapter/in/cli + config)

| Thư viện | Vai trò | Ghi chú |
| :--- | :--- | :--- |
| `knadh/koanf` | Config (nhẹ) | Load + merge từ nhiều nguồn. Gọn nhẹ, hợp dự án ưu tiên tối giản. |
| `spf13/viper` | Config (thay thế) | Phổ biến, đa nguồn (YAML, ENV, flags) nhưng nặng. |
| `spf13/cobra` | CLI framework | Chuẩn (K8s, Docker dùng). Subcommand, flags, auto-completion, sinh docs. |

> Khuyến nghị: `koanf` + `cobra`. Quen `viper` thì dùng được nhưng chú ý kích thước dependency.

## 7. Testing

| Thư viện | Vai trò | Ghi chú |
| :--- | :--- | :--- |
| `testcontainers-go` | Integration test | Spin up Postgres/MySQL trong test, chạy với DB thật. |
| `stretchr/testify` | Assert & mock | De-facto: `assert`, `require`, `mock`. Dùng để mock `port/in` và fake `port/out`. |

## 8. Observability (metrics & tracing)

| Thư viện | Vai trò | Ghi chú |
| :--- | :--- | :--- |
| `prometheus/client_golang` | Metrics | Expose `janus_requests_total`, `janus_query_duration_seconds`, v.v. |
| `go.opentelemetry.io/otel` | Tracing | Span per request: auth, validate, policy, execute. |

## 9. Tóm tắt chốt

| Thành phần | Thư viện đề xuất |
| :--- | :--- |
| MCP Framework | `finemcp/finemcp` |
| SQL Parser (Postgres) | `pganalyze/pg_query_go/v6` |
| SQL Parser (MySQL) | `vitess.io/vitess/go/vt/sqlparser` |
| DB Driver (Postgres) | `jackc/pgx/v5` |
| DB Driver (MySQL) | `go-sql-driver/mysql` |
| DB Driver (SQLite) | `modernc.org/sqlite` |
| Keyring | `zalando/go-keyring` |
| Vault Client | `hashicorp/vault-client-go` |
| Logging | `log/slog` |
| Config | `knadh/koanf` (hoặc `spf13/viper`) |
| CLI | `spf13/cobra` |
| Testing | `testcontainers-go`, `stretchr/testify` |
| Metrics | `prometheus/client_golang` |
| Tracing | `go.opentelemetry.io/otel` |

## 10. Việc cần làm tiếp

- [ ] Pin version trong `go.mod` khi scaffold Phase 0 (đặc biệt `pg_query_go/v6` vs v5 trong spec cũ).
- [ ] Xác minh `finemcp` vs `mark3labs/mcp-go`: spike nhỏ adapter/in transport stdio với cả hai, đo test harness + SSE readiness rồi chốt.
- [ ] Cập nhật `docs/plans/implementation-plan.md` Phase 0.1 với danh sách dependency này.
