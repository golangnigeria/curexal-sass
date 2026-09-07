# Curexal CI/CD & GitHub Actions Workflow

This document outlines the Continuous Integration (CI) and Continuous Deployment (CD) pipeline architecture for the Curexal healthcare platform.

---

## 1. Overview & Core Philosophy

Curexal operates under strict clinical reliability, regulatory compliance, and multi-tenant data isolation standards. The CI pipeline is designed around the following principles:

1. **Deterministic Quality Gates**: Formatting, static analysis, type checking, contract verification, and unit testing must pass before any code can merge.
2. **Ephemeral Service Isolation**: Database migrations and tests run against dedicated, disposable container services (`postgres:16-alpine`, `redis:7-alpine`) inside the CI runner. CI never connects to, touches, or alters staging, production, or real patient infrastructure.
3. **No Automated Production Deployment Without Staging Sign-Off**: Automated direct-to-production continuous deployment is intentionally withheld. In a healthcare domain, changes must undergo staging verification, compliance checks, and explicit promotion gates.
4. **Least-Privilege Execution**: All workflows operate under strict `permissions: contents: read` permissions.

---

## 2. GitHub Actions Workflow Architecture

The primary workflow file is located at [`.github/workflows/ci.yml`](../../.github/workflows/ci.yml).

### Triggers & Concurrency

- **Push Triggers**: Runs automatically on commits pushed to `main`, `master`, and `develop`.
- **Pull Request Triggers**: Runs automatically on PRs targeting `main`, `master`, and `develop`.
- **Manual Dispatch**: Can be triggered on demand via GitHub Actions UI (`workflow_dispatch`).
- **Concurrency Control**: 
  ```yaml
  concurrency:
    group: curexal-ci-${{ github.workflow }}-${{ github.ref }}
    cancel-in-progress: true
  ```
  Redundant runs on updated pull-request branches are cancelled immediately to conserve runner compute.

---

## 3. Pipeline Jobs & Stages

### Job 1: `backend-ci` (Go Backend API)

- **Environment**: `ubuntu-latest`
- **Service Containers**:
  - `postgres:16-alpine` (healthy check on `pg_isready`)
  - `redis:7-alpine` (healthy check on `redis-cli ping`)
- **Execution Steps**:
  1. **Source Checkout**: `actions/checkout@v4`
  2. **Go Toolchain Setup**: `actions/setup-go@v5` reading version from `apps/api/go.mod` with dependency caching.
  3. **Go Formatting Check**: Fails fast if any `.go` file violates `gofmt` standard formatting:
     ```bash
     UNFORMATTED=$(gofmt -l .)
     if [ -n "$UNFORMATTED" ]; then exit 1; fi
     ```
  4. **Static Analysis (`go vet`)**: Examines Go source code for suspicious constructs, structural mistakes, and shadowing.
  5. **Go Modules Verification**: Runs `go mod verify` to ensure hashes match `go.sum`.
  6. **Ephemeral Database Migration Validation**:
     ```bash
     CUREXAL_DB_DSN="postgres://curexal:password123@localhost:5432/curexal_test?sslmode=disable" go run ./cmd/migrate
     ```
     Executes all Goose platform and tenant migrations sequentially against the disposable Postgres service container. Ensures migrations are backwards-compatible and execute without schema syntax errors.
  7. **Test Suite Execution**:
     ```bash
     go test -v -race ./...
     ```
     Runs unit and integration tests with Go race detection enabled.
  8. **Binary Compilation**:
     ```bash
     go build -v -ldflags="-w -s" ./cmd/CUREXAL
     ```
     Validates that production binaries compile cleanly without linker errors.

### Job 2: `frontend-ci` (Vite / React 19 SPAs)

- **Environment**: `ubuntu-latest`
- **Execution Steps**:
  1. **Source Checkout**: `actions/checkout@v4`
  2. **Bun Runtime Setup**: `oven-sh/setup-bun@v2` (version: `latest`)
  3. **Dependency Installation**:
     ```bash
     bun install --frozen-lockfile
     ```
     Guarantees zero silent dependency drift by locking strictly to `bun.lock`.
  4. **TypeScript Typechecking**:
     - `apps/web-platform`: `bun run typecheck` (`tsc --noEmit`)
     - `apps/web-public`: `bun run typecheck` (`tsc --noEmit`)
     - `apps/web-patient`: `bun run typecheck` (`tsc --noEmit`)
  5. **Frontend Unit Tests**:
     - `apps/web-platform`: `bun test`
     - `apps/web-public`: `bun test`
  6. **Production Bundle Compilation**:
     - `apps/web-platform`: `bun run build`
     - `apps/web-public`: `bun run build`
     - `apps/web-patient`: `bun run build`

### Job 3: `security-audit` (Vulnerability Scanning)

- **Tool**: `golang/govulncheck-action@v1`
- **Scope**: Analyzes `apps/api` dependencies and call graphs for known CVEs reported in the Go vulnerability database.

---

## 4. How to Reproduce CI Checks Locally

Developers should always run CI checks locally prior to opening or pushing to pull requests.

### Using `Taskfile` (Recommended)

```bash
# Run full local CI suite (formatting, linting, tests, typechecking, builds)
task ci

# Run individual checks
task fmt         # Format Go codebase
task lint        # Run go vet on backend
task test        # Run backend test suite
task typecheck   # Typecheck all frontend apps
task build       # Build backend binary and all frontend SPAs
```

### Using `Makefile`

```bash
# Run all backend checks
make check

# Individual commands
make fmt         # Format Go files
make lint        # Run go vet
make test        # Run backend tests
make build       # Compile backend binary
```

### Direct CLI Commands

#### Backend:
```bash
cd apps/api

# 1. Format check
gofmt -l .

# 2. Static analysis
go vet ./...

# 3. Test execution
go test -v ./internal/kernel/...

# 4. Build binary
go build -v -ldflags="-w -s" -o bin/curexal-backend ./cmd/CUREXAL
```

#### Frontend:
```bash
# In repo root
bun install --frozen-lockfile

# Typechecks
cd apps/web-platform && bun run typecheck
cd ../web-public && bun run typecheck
cd ../web-patient && bun run typecheck

# Builds
cd apps/web-platform && bun run build
cd ../web-public && bun run build
cd ../web-patient && bun run build
```

---

## 5. Clinical Safety & Deployment Strategy

### Why Production Auto-Deploy is Held Back

In multi-tenant healthcare software, data integrity and high availability are mission-critical. Automated CD straight to production on branch push is intentionally omitted for the following reasons:

1. **Schema Migration Integrity**: In production, migrations may execute on databases containing gigabytes of patient records. Zero-downtime migrations (expand/contract pattern) must be vetted against staging replicas before production application.
2. **Tenant Boundary Verification**: Cross-tenant isolation and host-based workspace routing must be validated in staging environments with realistic subdomains (`curexal.space`).
3. **Audit Compliance**: Healthcare regulations (HIPAA, GDPR, local health authorities) require auditable deployment records with explicit human approvals.
4. **Promotion Pipeline**: Deployment to production follows a staged path:
   - CI automated verification (PR / branch)
   - Merge to `develop` $\rightarrow$ Automated deployment to Staging/Dev environment
   - Merge or Tag to `main` $\rightarrow$ Production deployment requires manual approval and change management ticket.
