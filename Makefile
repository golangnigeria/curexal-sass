.PHONY: all build test clean api web-public web-platform

all: build

# ------------------------------------------------------------------------------
# Build Targets
# ------------------------------------------------------------------------------
build: api web-public web-platform

api:
	@echo "==> Building Curexal API (Native Go Binary)..."
	cd apps/api && go build -ldflags="-w -s" -o bin/curexal-backend ./cmd/CUREXAL

web-public:
	@echo "==> Building Curexal Public Web (Static Vite SPA)..."
	cd apps/web-public && bun run build

web-platform:
	@echo "==> Building Curexal Platform Portal (Static Vite SPA)..."
	cd apps/web-platform && bun run build

# ------------------------------------------------------------------------------
# Test Targets
# ------------------------------------------------------------------------------
test:
	@echo "==> Running Backend Unit & Contract Tests..."
	cd apps/api && go test -v ./internal/kernel/...

# ------------------------------------------------------------------------------
# Development Run Targets
# ------------------------------------------------------------------------------
dev-api:
	cd apps/api && go run ./cmd/CUREXAL

dev-public:
	cd apps/web-public && bun run dev

dev-platform:
	cd apps/web-platform && bun run dev

clean:
	rm -rf apps/api/bin apps/web-public/dist apps/web-platform/dist
