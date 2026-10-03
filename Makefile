BINARY := janus
PKG := ./...
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -s -w -X github.com/vietthanh1999/ohjanus/internal/adapter/in/cli.Version=$(VERSION)

.PHONY: build test vet lint run validate clean

build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) ./cmd/janus

test:
	go test $(PKG)

vet:
	go vet $(PKG)

lint:
	golangci-lint run

run: build
	./$(BINARY) serve --config ./janus.yaml

validate:
	go run ./cmd/janus validate --config ./janus.yaml

clean:
	rm -f $(BINARY)
