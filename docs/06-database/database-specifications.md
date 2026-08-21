# Curexal V2 Database Specifications & Schema Strategy

This document establishes the technical design, entity-relationship structures, dynamic schema partitioning, indexing rules, and migration pipelines for Curexal V2.

---

## 1. Schema-per-Tenant Database Topology

Curexal utilizes a hybrid multi-tenant structure in PostgreSQL:
1. **`public` Schema**: Shared global registry. Stores user credentials, MFA secrets, tenant metadata, public marketplace discovery catalogs, and inter-tenant referral routing logs.
2. **`tenant_<slug>` Schemas**: Isolated tenant operational schemas. Stores patient EMR records, consultation notes, phlebotomy benches, local billing invoices, and branch settings.

```
                      +────────────────────────────────────────+
                      |         PostgreSQL Database            |
                      +───────────────────┬────────────────────+
                                          |
                ┌─────────────────────────┴─────────────────────────┐
                ▼                                                   ▼
+───────────────────────────────+                   +───────────────────────────────+
|         public Schema         |                   |      tenant_everight_labs     |
|                               |                   |                               |
| - users                       |                   | - patients                    |
| - user_credentials            |                   | - specimens                   |
| - organizations               |                   | - test_orders                 |
| - organization_memberships    |                   | - diagnostic_results          |
| - marketplace_catalogs        |                   | - invoices                    |
| - inter_tenant_referrals      |                   | - branch_settings             |
+───────────────────────────────+                   +───────────────────────────────+
```

---

## 2. Dynamic Schema Isolation & Transaction Safety

To guarantee zero cross-schema connection pool pollution in Go:

```go
func (mgr *DBManager) RunInTenantTx(ctx context.Context, tenantSlug string, fn func(ctx context.Context, tx bun.Tx) error) error {
	tx, err := mgr.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	schema := "tenant_" + tenantSlug
	_, err = tx.ExecContext(ctx, "SET LOCAL search_path TO ?, public", bun.Ident(schema))
	if err != nil {
		return fmt.Errorf("failed to set search_path: %w", err)
	}

	if err := fn(ctx, tx); err != nil {
		return err
	}

	return tx.Commit()
}
```

---

## 3. Core Database Table Definitions

### Global `public.organizations` Table
```sql
CREATE TABLE public.organizations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(100) UNIQUE NOT NULL,
    facility_type VARCHAR(50) NOT NULL, -- 'clinic', 'laboratory', 'pharmacy', 'radiology'
    status VARCHAR(50) NOT NULL DEFAULT 'Provisioning', -- 'Provisioning', 'Ready', 'Suspended'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_organizations_slug ON public.organizations(slug);
```

### Tenant `tenant_<slug>.specimens` Table
```sql
CREATE TABLE specimens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    referral_order_id UUID,
    barcode VARCHAR(100) UNIQUE NOT NULL,
    sample_type VARCHAR(50) NOT NULL,
    container_type VARCHAR(50) NOT NULL,
    collected_by UUID NOT NULL,
    collected_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status VARCHAR(50) NOT NULL DEFAULT 'Collected'
);

CREATE INDEX idx_specimens_barcode ON specimens(barcode);
```

---

## 4. Migration Strategy

- Migrations are versioned SQL scripts located in `internal/modules/[module]/infra/migrations/`.
- The dynamic schema provisioner runs migrations across all active tenant schemas sequentially or concurrently during upgrades.
