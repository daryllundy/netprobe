BINARY_NAME := netprobe
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
GO := go
LDFLAGS := -ldflags="-s -w -X main.Version=$(VERSION) -X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT)"

.PHONY: all
all: clean test build

.PHONY: build
build:
	$(GO) build $(LDFLAGS) -o build/$(BINARY_NAME) ./cmd/netprobe

.PHONY: test
test:
	$(GO) test -v -race ./...

.PHONY: clean
clean:
	rm -rf build/

.PHONY: install
install: build
	install -m 755 build/$(BINARY_NAME) /usr/local/bin/

.PHONY: run
run: build
	./build/$(BINARY_NAME)

.PHONY: lint
lint:
	golangci-lint run

.PHONY: deps
deps:
	$(GO) mod download
	$(GO) mod tidy

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build    - Build the binary"
	@echo "  test     - Run tests"
	@echo "  clean    - Clean build artifacts"
	@echo "  install  - Install binary"
	@echo "  run      - Build and run"
	@echo "  lint     - Run linters"
	@echo "  deps     - Download dependencies"
