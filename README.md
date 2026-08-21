# Curexal SASS — Production Monorepo

Unified healthcare operating system and clinical management platform for Africa and emerging markets.

---

## 1. Monorepo Architecture

```text
curexal-sass/
├── apps/
│   ├── api/                          # High-performance native Go Echo backend kernel
│   │   ├── cmd/CUREXAL/              # Main server entrypoint
│   │   ├── database/migrations/      # PostgreSQL migrations & seeders
│   │   ├── internal/                 # Hexagonal DDD modules (identity, org, clinical, etc.)
│   │   └── storage/documents/        # Local filesystem storage fallback
│   ├── web-public/                   # Static Vite + React marketing, demo & waitlist site
│   └── web-platform/                 # Static Vite + React multi-tenant workspace dashboard
├── docs/                             # Architecture specs, PRDs & API manuals
├── .env.example                      # Complete environment configuration template
├── Makefile                          # Unified build & testing commands
└── README.md
```

---

## 2. Infrastructure Principles

- **Dockerless Baseline**: The backend compiles directly into a native single binary (`curexal-backend`) with zero external daemons. Frontends compile to static SPAs.
- **Database**: PostgreSQL / NeonDB Serverless via standard `CUREXAL_DB_DSN`.
- **Object Storage**: S3-compatible persistent storage (Cloudflare R2 for production with zero egress fees; local filesystem for local dev).
- **Ephemeral Cache**: Zero-dependency Go in-memory TTL cache (Postgres remains authoritative for all durable data).
- **Domain Topology**: Served under single root domain `curexal.space` (`curexal.space` for public web, `app.curexal.space` for portal, `api.curexal.space` for REST API).

---

## 3. Quick Start & Build

### Prerequisites
- Go 1.23+
- Bun (or Node.js 20+)

### Building Production Artifacts
```bash
# Build everything (Native Go binary + Static SPAs)
make build

# Or build individual applications:
make api
make web-public
make web-platform
```

### Running Locally
```bash
# Start backend API (Port 8080)
make dev-api

# Start Public Marketing site (Port 5001)
make dev-public

# Start Platform Portal (Port 5002)
make dev-platform
```
