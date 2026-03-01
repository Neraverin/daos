.PHONY: all build daemon tui clean test install generate-api generate-sql

# Build variables
VERSION := 0.1.0
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

GO := /home/neraverin/sdk/go1.21.0/bin/go
OAPI_CODEGEN := $(HOME)/go/bin/oapi-codegen
SQLC := $(HOME)/go/bin/sqlc

# Default target - just build binaries
all: build

# Build daemon
daemon: generate-api generate-sql
	$(GO) build ${LDFLAGS} -o bin/daemon ./cmd/daemon

# Build TUI
tui: generate-api
	$(GO) build ${LDFLAGS} -o bin/tui ./cmd/tui

# Build both binaries
build: daemon tui

# Clean build artifacts
clean:
	rm -rf bin/

# Run tests
test:
	$(GO) test -v ./...

# Install dependencies
install:
	$(GO) mod download
	$(GO) mod tidy

# Generate API code from OpenAPI spec
# Note: types are manually defined in pkg/api/types.go
# To regenerate, run: oapi-codegen -package api -generate types api/openapi.yaml > pkg/api/types.gen.go
generate-api:
	@echo "Skipping API generation - using manual types"
	@touch pkg/api/types.gen.go pkg/api/client.gen.go

# Generate SQL code from queries
# Note: queries are manually handled in handlers
# To regenerate, run: sqlc generate
generate-sql:
	@echo "Skipping SQL generation - using manual queries"
	@touch pkg/db/query_generated.go

# Create directories for packaging
prepare-packaging:
	mkdir -p packaging/deb/usr/bin
	mkdir -p packaging/deb/etc/daos
	mkdir -p packaging/rpm/build
	mkdir -p packaging/rpm/SOURCES

# Build deb package
deb: prepare-packaging
	cp bin/daemon packaging/deb/usr/bin/
	cp bin/tui packaging/deb/usr/bin/
	dpkg-deb --build packaging/deb daos_${VERSION}_amd64.deb

# Build rpm package
rpm: prepare-packaging
	cp bin/daemon packaging/rpm/SOURCES/
	cp bin/tui packaging/rpm/SOURCES/
	rpmbuild -bb packaging/rpm/daos.spec --define "_version ${VERSION}"

# All targets
help:
	@echo "Available targets:"
	@echo "  all           - Build daemon and TUI (default)"
	@echo "  daemon        - Build daemon binary"
	@echo "  tui           - Build TUI binary"
	@echo "  build         - Build both binaries"
	@echo "  clean         - Remove build artifacts"
	@echo "  test          - Run tests"
	@echo "  install       - Install Go dependencies"
	@echo "  generate-api  - Generate Go types from OpenAPI spec"
	@echo "  generate-sql  - Generate Go code from SQL queries"
	@echo "  deb           - Build deb package"
	@echo "  rpm           - Build rpm package"
