.PHONY: help fmt lint tidy test build check clean

help:
	@echo "Available Curexal local developer commands:"
	@echo "  make fmt        - Format Go backend codebase"
	@echo "  make lint       - Run static analysis (go vet) on backend"
	@echo "  make tidy       - Tidy and verify Go module dependencies"
	@echo "  make test       - Run backend kernel test suite"
	@echo "  make build      - Build native Go backend binary"
	@echo "  make check      - Run all local quality checks (fmt, lint, test, build)"
	@echo "  make clean      - Clean build artifacts"

fmt:
	@cd apps/api && go fmt ./...

lint:
	@cd apps/api && go vet ./...

tidy:
	@cd apps/api && go mod tidy && go mod verify

test:
	@cd apps/api && go test -v ./internal/kernel/...

build:
	@cd apps/api && go build -v -ldflags="-w -s" -o bin/curexal-backend ./cmd/CUREXAL

check: fmt lint tidy test build
	@echo "All local checks passed successfully!"

clean:
	@rm -rf apps/api/bin apps/web-public/dist apps/web-platform/dist apps/web-patient/dist
