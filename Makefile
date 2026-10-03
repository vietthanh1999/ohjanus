BINARY := janus
PKG := ./...
CONFIG ?= ./janus.dev.yaml
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/vietthanh1999/ohjanus/internal/adapter/in/cli.Version=$(VERSION)

.PHONY: help build test vet lint run validate clean \
	setup dev dev-api dev-ui build-ui check-ui

help:
	@echo "Targets:"
	@echo "  setup    - pnpm install (root, ui, packages)"
	@echo "  dev-api  - run Go Admin API :8788 (CONFIG=$(CONFIG))"
	@echo "  dev-ui   - run Vite UI :5173 (proxies /api -> :8788)"
	@echo "  dev      - run api + ui together (Ctrl+C to stop)"
	@echo "  run      - build + serve with CONFIG=$(CONFIG)"
	@echo "  validate - validate CONFIG=$(CONFIG)"
	@echo "  build    - build Go binary"
	@echo "  build-ui - build UI (vite build)"
	@echo "  test     - go test ./..."
	@echo "  check-ui - svelte-check + tsc (ui)"
	@echo "  vet/lint - go vet / golangci-lint"

setup:
	pnpm install

# Backend only: fixes `vite http proxy error ... ECONNREFUSED 127.0.0.1:8788`
dev-api:
	go run ./cmd/janus serve --config $(CONFIG)

# Frontend only: expects API on :8788 (see ui/vite.config.ts)
dev-ui:
	pnpm --filter ui dev

# Both together, single terminal. Ctrl+C stops all.
dev:
	trap 'kill 0' INT TERM; \
	go run ./cmd/janus serve --config $(CONFIG) & \
	pnpm --filter ui dev & \
	wait

build-ui:
	pnpm --filter ui build

check-ui:
	pnpm --filter ui check

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/janus

test:
	go test $(PKG)

vet:
	go vet $(PKG)

lint:
	golangci-lint run

run: build
	./$(BINARY) serve --config $(CONFIG)

validate:
	go run ./cmd/janus validate --config $(CONFIG)

clean:
	rm -f $(BINARY)
