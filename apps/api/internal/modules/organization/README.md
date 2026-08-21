# Organization Bounded Context (`internal/modules/organization`)

## 1. Purpose & Scope
The **Organization Module** manages the multi-tenant organization hierarchy, operational branch tenant creation, database schema provisioning (`CREATE SCHEMA tenant_<slug>`), enterprise billing subscriptions, and public demo request processing.

---

## 2. Architectural Layers

- **`domain/`**: Entities (`Organization`, `Tenant`, `Subscription`, `DemoRequest`), Value Objects, Domain Errors (`errors.go`), and `repository.go` contracts.
- **`application/`**: Use cases & application services (`TenantService`, `OrganizationService`, `DemoService`).
- **`infrastructure/postgres/`**: Database repositories (`TenantRepository`, `OrganizationRepository`, `SubscriptionRepository`, `DemoRepository`) executing SQL queries on `public` schema.
- **`api/http/`**: HTTP controllers (`tenant_handler.go`, `organization_handler.go`, `demo_handler.go`) and Echo route bindings.
- **`module.go`**: Dependency container exporting `NewModule(server)`.

---

## 3. Database Schema & Tables Scope

- **Schema**: `public`
- **Tables**: `organizations`, `organization_settings`, `tenants`, `subscriptions`, `demo_requests`, `tenant_modules`.

---

## 4. Exported Public APIs

- `POST /api/v1/demo-requests`: Public sales lead submission
- `GET /api/v1/demo-requests`: List demo requests (requires `demo:read`)
- `POST /api/v1/tenants`: Create new branch tenant facility
- `GET /api/v1/organizations`: List parent laboratory networks

---

## 5. Testing
```bash
go test ./apps/backend/internal/modules/organization/...
```
