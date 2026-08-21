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
│   │   ├── storage/documents/        # Local filesystem storage fallback
│   │   └── Taskfile.yml              # API-specific tasks (migrations, db, tidy)
│   ├── web-public/                   # Static Vite + React marketing, demo & waitlist site
│   └── web-platform/                 # Static Vite + React multi-tenant workspace dashboard
├── docs/                             # Architecture specs, PRDs & API manuals
├── .env.example                      # Complete environment configuration template
├── Taskfile.yml                      # Unified Task runner commands
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

## 3. Quick Start & Task Commands

### Prerequisites
- Go 1.23+
- Bun (or Node.js 20+)
- Task CLI (`task`)

### Available Tasks
```bash
# List all workspace tasks
task

# Build all production artifacts (Native Go binary + Static SPAs)
task build

# Or build individual applications:
task build:api
task build:web-public
task build:web-platform

# Run backend test suites
task test

# Start applications in development:
task dev:api          # Native Go API on port 8080
task dev:web-public   # Public marketing site on port 5001
task dev:web-platform # Platform portal dashboard on port 5002

# Database & migration tasks (namespaced under api)
task api:migrations:up
task api:migrations:new name=add_feature_table
task api:db:reset
```
