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

---

## 3. Curexal Edge Agent (LIS Instrument Integration)

Laboratory analyzers communicate raw scientific outputs using legacy ASTM E1394 or HL7 v2 protocol strings. To link these instruments to our cloud platform securely:

- **Deployment**: Lightweight Go binary running on a local gateway PC or Raspberry Pi within the laboratory network.
- **Dual-Buffer Caching**: Uses a local SQLite database engine. If the internet drops, raw analyzer outputs are stored locally, mapped to barcode reads, and synchronized via WebSocket streams when online.
- **Security**: Authenticates using Mutual TLS (mTLS) with client certificates rotated automatically by the platform PKI.
