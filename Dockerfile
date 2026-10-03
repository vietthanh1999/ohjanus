# OhJanus single-artifact image: Svelte Admin UI embedded in the Go binary.
# The Admin API serves the UI same-origin, so no separate web server is needed.

# ---- UI build ----
FROM node:22-alpine AS uibuild
WORKDIR /src
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml turbo.json ./
COPY packages ./packages
COPY ui ./ui
RUN corepack enable && corepack prepare pnpm@11.1.1 --activate && \
    pnpm install --frozen-lockfile && \
    pnpm build

# ---- Go build ----
FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=uibuild /src/ui/dist/. ./internal/adapter/in/admin/uistatic/dist/
RUN CGO_ENABLED=0 go build -ldflags "-s -w" -o /janus ./cmd/janus

# ---- Runtime (no shell, non-root) ----
FROM gcr.io/distroless/static-debian12
COPY --from=builder /janus /janus
COPY configs/janus.example.yaml /etc/janus/janus.yaml
USER nonroot:nonroot
ENTRYPOINT ["/janus", "serve", "--config", "/etc/janus/janus.yaml"]
