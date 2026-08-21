# Curexal V2 Platform Kernel Manifest

This document catalogs every package within `internal/core/` (Platform Kernel), defining their core responsibilities, ownership, dependencies, and architectural rules.

---

## Objective Statement

> **EPIC-000 is not about building utilities. It is about establishing the engineering platform that every future Curexal module depends upon. Any shared concern that is likely to be reused across multiple bounded contexts belongs in the Platform Kernel unless there is a compelling architectural reason otherwise.**

---

## Master Kernel Package Catalog

| Package Path | Core Responsibility | Internal Dependencies | Downstream Consumers | Status |
| :--- | :--- | :--- | :--- | :--- |
| `internal/core/primitive` | Testable `Clock` interface, ULID & UUIDv7 generators. | *None* | All Modules | Baseline v1.0 |
| `internal/core/errors` | Centralized sentinel errors & RFC 7807 problem details mapper. | `primitive` | All Modules | Baseline v1.0 |
| `internal/core/logging` | Structured JSON logger wrapper (Zero-allocation). | `primitive` | All Modules | Baseline v1.0 |
| `internal/core/telemetry` | OpenTelemetry tracing abstractions & correlation context. | `primitive`, `logging` | All Modules | Baseline v1.0 |
| `internal/core/metrics` | Prometheus metrics counters, histograms, `/metrics` handler. | `primitive` | Hertz Server | Baseline v1.0 |
| `internal/core/config` | Environment profile loader & secret validator. | `errors` | All Modules | Baseline v1.0 |
| `internal/core/security` | Argon2 hashing, AES encryption, HMAC signatures, secure tokens. | `errors`, `primitive` | Auth, Payments | Baseline v1.0 |
| `internal/core/database` | PostgreSQL connection pool setup & Bun ORM dialect init. | `config`, `logging` | Transaction Manager | Baseline v1.0 |
| `internal/core/transaction` | Centralized `TransactionManager` (`RunInTenantTx`, retries, traces). | `database`, `errors` | Repositories | Baseline v1.0 |
| `internal/core/cache` | Redis connection pool, key caching, distributed locks. | `config`, `errors` | Sessions, Rate Limit | Baseline v1.0 |
| `internal/core/validation` | Centralized request DTO validation (`go-playground/validator`). | `errors` | Handlers | Baseline v1.0 |
| `internal/core/audit` | Hash-chained immutable audit log writer. | `database`, `primitive` | All Modules | Baseline v1.0 |
| `internal/core/event` | NATS JetStream event publisher & transactional outbox. | `database`, `logging` | All Modules | Baseline v1.0 |
| `internal/core/authz` | Casbin RBAC policy adapter & context enforcer. | `database`, `errors` | Handlers | Baseline v1.0 |
| `internal/core/storage` | Abstract object storage provider (MinIO / S3 signed URLs). | `config`, `errors` | LIMS, RIS, Documents| Baseline v1.0 |
| `internal/core/featureflag` | Feature flag evaluator across tenant plans & types. | `database`, `cache` | Application Services| Baseline v1.0 |
| `internal/core/notification` | Abstract notification dispatcher (Email, SMS, WhatsApp adapters). | `event`, `logging` | All Modules | Baseline v1.0 |
| `internal/core/health` | Liveness (`/healthz`) & Readiness (`/readyz`) probe checkers. | `database`, `cache` | Hertz Server | Baseline v1.0 |
| `internal/core/middleware` | Hertz HTTP middleware stack (Correlation, Auth, Rate limit). | `telemetry`, `errors` | Hertz Server | Baseline v1.0 |
| `internal/core/rest` | HTTP response helpers & OpenAPI Swagger UI handlers. | `errors` | Handlers | Baseline v1.0 |

---

## Architectural Layering Rules

```text
       Presentation Layer (api/ Handlers, DTOs)
                        │
                        ▼
      Application Layer (app/ Services, Use Cases)
                        │
                        ▼
         Domain Layer (domain/ Entities, Logic)
                        │
                        ▼
   Infrastructure & Core (infra/ & internal/core/)
```

### Strict Import Constraints
1. **Domain Layer**: Must have **zero dependencies** on external frameworks, Bun ORM, Hertz HTTP, or SQL drivers.
2. **Business Logic Prohibition**: Business code must never import `net/http` or HTTP-specific packages directly.
3. **Transaction Safety**: Repositories must consume `bun.IDB` or `transaction.Tx` interfaces rather than managing `sql.DB` connections directly.
4. **Kernel Isolation**: Packages within `internal/core/` must not import business modules (`internal/modules/*`).
