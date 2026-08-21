# Audit Bounded Context (`internal/modules/audit`)

## 1. Purpose & Scope
The **Audit Module** provides high-performance asynchronous audit logging, security compliance trails, platform activity monitoring, and administrative access logging.

---

## 2. Architectural Layers

- **`domain/`**: Entities (`AuditLog`, `PlatformAuditLog`), Value Objects, Domain Errors (`errors.go`), and `repository.go` contracts.
- **`application/`**: Use cases & application services (`AuditApplicationService`).
- **`infrastructure/postgres/`**: Database repository (`AuditRepository`) with dual support for `public.platform_audit_logs` and `tenant_<slug>.audit_logs`.
- **`api/http/`**: HTTP controller (`audit_handler.go`) and Echo route bindings.
- **`module.go`**: Dependency container exporting `NewModule(server)`.

---

## 3. Database Schema & Tables Scope

- **Schema**: `public` & `tenant_<slug>`
- **Tables**: `public.audit_logs`, `public.platform_audit_logs`, `tenant_<slug>.audit_logs`.

---

## 4. Exported Public APIs

- `POST /api/v1/audit-logs`: Log tenant activity event
- `GET /api/v1/audit-logs/platform`: Query platform-wide audit log trail
- `GET /api/v1/admin/audit-logs`: Query admin audit log traces (requires `audit:read`)

---

## 5. Testing
```bash
go test ./apps/backend/internal/modules/audit/...
```
