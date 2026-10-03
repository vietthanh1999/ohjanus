# OhJanus Production Runbook

Single-artifact deployment: the `janus` binary serves MCP (stdio/http),
the Admin API, and the embedded Admin UI (same origin, no extra web server).

## 1. Build

```bash
make build          # pnpm build UI -> embed into Go binary -> ./janus
```

`make build` runs `sync-ui` (copies `ui/dist` for `go:embed`). A binary built
without the UI still works but serves API-only (`/`, non-API paths → 503).

Docker (multi-stage, distroless non-root):

```bash
docker build -t ohjanus:0.1.0 .
```

## 2. Configure (production checklist)

Start from `configs/janus.example.yaml` → `./janus.yaml` (git-ignored).

- [ ] `auth.mode: token` (never `none`; `serve` refuses `none` unless
      `JANUS_ALLOW_NO_AUTH=1`, which must never be set in production).
- [ ] Tokens declared by **hash** only (`janus token create` prints the
      YAML snippet; the secret is shown once and never stored).
- [ ] DSNs via `dsn_ref: "env:..."` or `file:` refs — never plaintext DSN.
- [ ] `approval.method: api` for headless/container use (`cli` prompts on
      the terminal and hangs without a TTY). Humans approve in the UI
      (Approvals Queue) or via `POST /api/v1/approvals/{id}/approve`.
- [ ] `policy.default_action: deny`, least-privilege `allowed_schemas`,
      `denied_tables` for sensitive tables.
- [ ] `limits.query_timeout` / `max_query_length` / rate limits set.
- [ ] `redaction.enabled: true` for PII columns.

## 3. Runtime token persistence

```yaml
auth:
  store: sqlite              # memory (default) | sqlite
  sqlite_path: "./janus-tokens.db"
```

- `memory`: tokens created at runtime (UI / `POST /api/v1/tokens`) are
  lost on restart; only `janus.yaml` tokens survive.
- `sqlite`: runtime tokens persist across restarts (only hashes stored,
  file created `0600`, WAL mode). Config-file tokens are merged on every
  boot and never overwrite runtime rows.
- The db path is stable identity: changing it starts from an empty store.
  In containers, mount it as a volume so restarts keep tokens.
- Monitor `audit` event `token.created` to track issuance.

## 4. TLS termination (remote operation)

Janus listens plaintext on localhost by default. Terminate TLS at a reverse
proxy; the Admin UI calls the API same-origin so no CORS work is needed.

### Caddy (automatic HTTPS, recommended)

`Caddyfile`:

```
janus.example.com {
    reverse_proxy 127.0.0.1:8788
}
```

`docker-compose.yml` (see root example):

```yaml
services:
  janus:
    build: .
    volumes:
      - ./janus.yaml:/etc/janus/janus.yaml:ro
    environment:
      JANUS_ANALYTICS_DSN: ${JANUS_ANALYTICS_DSN}
    expose: ["8788"]
    stdin_open: true          # keep stdin open for stdio transport
    tty: true
  caddy:
    image: caddy:2-alpine
    ports: ["80:80", "443:443"]
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - caddy_data:/data
      - caddy_config:/config
volumes:
  caddy_data:
  caddy_config:
```

### nginx (manual certificates)

```nginx
server {
    listen 443 ssl;
    server_name janus.example.com;
    ssl_certificate     /etc/ssl/janus.crt;
    ssl_certificate_key /etc/ssl/janus.key;
    location / {
        proxy_pass http://127.0.0.1:8788;
        proxy_set_header Host $host;
    }
}
```

## 5. Observability

- Health: `GET /healthz`, readiness: `GET /readyz` (public).
- Metrics: set `observability.metrics.enabled: true`, scrape
  `127.0.0.1:9090/metrics` with Prometheus:

```yaml
scrape_configs:
  - job_name: janus
    static_configs:
      - targets: ["127.0.0.1:9090"]
```

- Audit: `stdout` JSON sink → collect via Docker/journald/Filebeat into
  your SIEM. Prefer the `file` sink (`audit.sinks[].path`) for local
  retention; the in-memory ring buffer only feeds the live UI view.

## 6. Operational notes

- **stdio transport ties process lifetime to stdin**: EOF on stdin shuts
  the server down (Admin API included). In containers use `stdin_open/tty`
  or switch `server.transport` to `http` (MCP over `POST /mcp`).
- **Graceful shutdown**: SIGTERM drains with a 5s timeout for both MCP
  and Admin servers.
- **SQLite persistent token store**: pending work (v0.2) — see runbook §3
  until then.
