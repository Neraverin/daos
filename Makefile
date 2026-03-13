.PHONY: all build daemon tui clean test install-dependencies generate-api generate-sql deb rpm e2e-test

VERSION := 0.1.0
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-s -w -X main.version=${VERSION} -X main.commit=${COMMIT} -X main.date=${DATE}"

SQLC := $(HOME)/go/bin/sqlc

all: build deb rpm

build: daemon tui

daemon: generate-api generate-sql
	@go build ${LDFLAGS} -o bin/daemon ./cmd/daemon

tui: generate-api
	@go build ${LDFLAGS} -o bin/tui ./cmd/tui

clean:
	rm -rf bin/
	rm -rf packaging/deb
	rm -rf packaging/rpm
	rm -rf pkg/api/openapi.gen.go

test:
	@go test -v $(go list ./... | grep -v tests/e2e)

install-dependencies:
	@go mod download
	@go mod tidy

generate-api:
	@echo "--- Generating OpenAPI server and types ---"
	@go tool oapi-codegen -generate=types,gin-server,client -package=api -o pkg/api/openapi.gen.go docs/openapi.yaml

generate-sql:
	@echo "--- Generating sqlc ---"
	@go tool sqlc generate

prepare-packaging:
	mkdir -p packaging/deb/opt/daos
	mkdir -p packaging/rpm/build
	mkdir -p packaging/rpm/SOURCES

deb: prepare-packaging
	cp bin/daemon packaging/deb/opt/daos
	cp bin/tui packaging/deb/opt/daos
	cp -r packaging/systemd packaging/deb/opt/daos/
	dpkg-deb --build packaging/deb bin/daos_${VERSION}_amd64.deb

rpm: prepare-packaging
	cp bin/daemon packaging/rpm/SOURCES
	cp bin/tui packaging/rpm/SOURCES
	cp -r packaging/systemd packaging/rpm/SOURCES
	cd packaging/rpm/SOURCES && tar -czvf ../SOURCES/daos-${VERSION}.tar.gz daemon tui systemd
	cd packaging/rpm && rpmbuild -bb daos.spec --define "_version ${VERSION}" --define "_topdir $(shell pwd)/packaging/rpm"
	cp packaging/rpm/RPMS/x86_64/daos-${VERSION}-1.x86_64.rpm bin/

e2e-test: build deb rpm
	sudo go test -v ./tests/e2e/...

help:
	@echo "Available targets:"
	@echo "  all                    - Build daemon and TUI (default)"
	@echo "  daemon                 - Build daemon binary"
	@echo "  tui                   - Build TUI binary"
	@echo "  build                  - Build both binaries"
	@echo "  clean                  - Remove build artifacts"
	@echo "  test                   - Run unit tests (excludes e2e)"
	@echo "  install-dependencies   - Install Go dependencies"
	@echo "  generate-api           - Generate Go types from OpenAPI spec"
	@echo "  generate-sql           - Generate Go code from SQL queries"
	@echo "  deb                    - Build deb package"
	@echo "  rpm                    - Build rpm package"
	@echo "  e2e-test               - Build packages and run e2e tests (requires root)"
