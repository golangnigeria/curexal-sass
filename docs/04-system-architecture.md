# Curexal V2 Technical Architecture Specifications

This document defines the complete technical architecture, infrastructure specifications, data pipelines, security protocols, and monorepo structural patterns for the Curexal V2 platform.

---

## 1. Monorepo & Frontend Applications Topology

Curexal V2 is built as a Turborepo monorepo connecting four distinct frontend web applications to a Go backend (Hertz HTTP framework).

```
Curexal V2 Monorepo Structure
│
├── apps/
│   ├── portal/          # Public marketing, pricing, B2B organization onboarding wizard
│   ├── workspace/       # Unified SaaS dashboard (Clinics, Labs, Pharmacies, RIS)
│   ├── patient/         # Patient self-service portal (Results, Appointments, Payments)
│   └── admin/           # Platform control center (Schema status, Metrics, ICD-10)
│
├── packages/
│   ├── ui/              # Shared Enterprise Design System (shadcn/ui + Tailwind CSS)
│   ├── api/             # Auto-generated TypeScript API client bindings & types
│   └── utils/           # Shared validation schemas, date parsers, and financial helpers
│
└── internal/            # Backend Clean Architecture & Modular Monolith in Go
```

### Frontend State & Offline Strategy
- **TanStack Query + IndexedDB**: API requests are cached using TanStack Query. During network drops, mutations are saved to client-side `IndexedDB` with UUID keys and replayed chronologically when connectivity resumes.
- **CORS & Cookie Session Handling**: All web applications pass credentials enabled by the Hertz backend CORS middleware via HTTP-Only `SameSite=Lax` secure cookies.

---

## 2. Tenancy & Database Isolation Model

Curexal utilizes a hybrid multi-tenant model in PostgreSQL:
1. **Shared Metadata Schema (`public`)**: Contains the global discovery index, tenant registry, user credentials directory, and inter-tenant referral routing logs.
2. **Dynamic Isolated Schemas (`tenant_<slug>`)**: Contains tenant-specific data including patient EHR files, local billing invoices, phlebotomy benches, inventory catalogs, and branch configurations.

### Schema Switching Mechanism (Bun ORM & Go)
The application handles multi-tenancy dynamically. It intercepts the HTTP requests, resolves the tenant slug, and runs all queries inside a database transaction scoped using PostgreSQL's `SET LOCAL search_path`.

```go
func (mgr *DBManager) RunInTenantTx(ctx context.Context, tenantSlug string, fn func(ctx context.Context, tx bun.Tx) error) error {
	tx, err := mgr.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	schema := "tenant_" + tenantSlug

	// Set search_path local to this database transaction
	_, err = tx.ExecContext(ctx, "SET LOCAL search_path TO ?, public", bun.Ident(schema))
	if err != nil {
		return fmt.Errorf("failed to set search_path to tenant schema: %w", err)
	}

	if err := fn(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}
```

### Database Sharding & Connection Pool Limits
- Connection pool limits per node: `SetMaxOpenConns(25)`, `SetMaxIdleConns(5)`, `SetConnMaxLifetime(5 * time.Minute)`.
- **Horizontal Multi-DB Sharding**: As tenant volume grows past 1,000 schemas, the global tenant registry maps `TenantSlug -> { ShardConnectionString, SchemaName }`, enabling database traffic to be routed across sharded PostgreSQL database clusters seamlessly.

---

## 3. Dynamic Schema Provisioning Pipeline

When a new organization registers:
1. An entry is created in the global `organizations` registry with status `Provisioning`.
2. An asynchronous background worker executes the dynamic setup sequence:
   - Executes `CREATE SCHEMA tenant_<slug>;`
   - Runs all system DDL migrations against `tenant_<slug>`.
   - Populates initial seeds (default roles, departments, tax rates, invoice counters).
   - Updates organization status to `Ready`.

---

## 4. Universal SSO & Session Resolution

- All user credential details exist in the `public` schema's global directories.
- Users authenticate and receive a 15-minute JWT access token and a 7-day refresh token.
- To access tenant resources, requests must pass the `X-Tenant-Slug` header.
- The **Tenant Resolver Middleware** checks if the user has an active membership record in the target tenant registry. If valid, the middleware attaches the schema path to the request context.

---

## 5. Role-Based Access Control (RBAC) via Casbin

Curexal uses Casbin for access control:
- Casbin policy rules are stored dynamically inside tenant-isolated database tables.
- The `casbin_adapter.go` evaluates permissions mapped to roles like `doctor`, `scientist`, or `receptionist`.
- Permissions are strictly evaluated at the Hertz middleware layer before any service execution occurs.

---

## 6. Event-Driven Messaging (NATS JetStream)

Inter-tenant collaborations must never execute cross-schema database joins. All module integrations are performed asynchronously via NATS JetStream:

```json
{
  "event_id": "018e-4aef-92b1",
  "event_type": "referral.dispatched",
  "producer_tenant": "everight-clinic",
  "timestamp": "2026-07-21T17:34:00Z",
  "schema_version": "v1",
  "data": {
    "referral_id": "ref-99201",
    "target_lab_slug": "synlab-main",
    "patient_checksum": "sha256-a9b8..."
  }
}
```

- **Delivery Guarantees**: NATS JetStream is configured with at-least-once delivery.
- **Dead-Letter Queues & Reconciliation**: Unprocessed events are routed to a DLQ after 5 retries. A daily reconciliation daemon scans inter-tenant ledger logs to fix state anomalies.

---

## 7. Curexal Edge Agent (LIS Instrument Integration)

Laboratory analyzers communicate raw scientific outputs using legacy ASTM E1394 or HL7 v2 protocol strings. To link these instruments to our cloud platform securely:

```
+──────────────────────────────+                +──────────────────────────────+
│     Laboratory LAN (Edge)    │                │         Cloud Platform       │
│                              │                │                              │
│ +──────────────────────────+ │                │      +─────────────────+     │
│ │   Laboratory Analyzer    │ │                │      │   API Gateway   │     │
│ │ (Sysmex / Mindray / etc) │ │                │      │   (HTTP/WS)     │     │
│ +────────────┬─────────────+ │                │      +────────┬────────+     │
│              │ ASTM / HL7  │                │               ▲                │
│              ▼             │                │               │ mTLS           │
│ +──────────────────────────+ │                │               │ WebSockets     │
│ │   Curexal Edge Agent     ├─┼────────────────┼───────────────┘              │
│ │   (SQLite Local Cache)   │ │                │                              │
│ +──────────────────────────+ │                │                              │
+──────────────────────────────+                +──────────────────────────────+
```

- **Deployment**: Lightweight Go binary running on a local gateway PC or Raspberry Pi within the laboratory network.
- **Dual-Buffer Caching**: Uses a local SQLite database engine. If the internet drops, raw analyzer outputs are stored locally, mapped to barcode reads, and synchronized via WebSocket streams when online.
- **Security**: Authenticates using Mutual TLS (mTLS) with client certificates rotated automatically by the platform PKI.

---

## 8. Object Storage & Signed URL Protocol

All diagnostic PDF reports, scanned medical documents, and PACS DICOM files are stored in S3-compatible storage (MinIO for local dev, AWS S3 for production):
- **Path Structure**: `vault/{tenant_slug}/{year}/{month}/{file_uuid}.pdf`
- **Private Access Control**: Buckets block public access entirely. Document links are served via short-lived AWS S3 pre-signed URLs valid for 15 minutes.

---

## 9. Payments & B2B Settlement Ledger

- **Split Payments**: Native Paystack and Flutterwave webhooks verify payments using cryptographic signatures (`X-Paystack-Signature`). Payments are automatically split between platform fee, lab wallet, and clinic commission.
- **Virtual Balance Ledger**: Cash payments at reception trigger a double-entry ledger entry (`Debit` Lab Wallet / `Credit` Clinic Wallet). If a tenant's wallet drops below their credit limit, incoming marketplace referrals are paused automatically.

---

## 10. Security & Observability Architecture

- **Correlation Tracing**: Every request is assigned a unique `X-Correlation-ID` header, propagated across Hertz HTTP handlers, Bun database queries, and NATS event messages.
- **Observability**: Exposes `/metrics` endpoints for Prometheus scraping, OpenTelemetry tracing for Jaeger, and structured JSON logs for Loki aggregation.
- **Cryptographic Audit Log**: Every patient record mutation creates an immutable audit entry whose hash relies on the hash of the preceding entry, making log tampering immediately detectable.
