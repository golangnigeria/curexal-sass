# Billing Bounded Context (`internal/modules/billing`)

## 1. Purpose & Scope
The **Billing Module** handles clinical patient invoicing, service item pricing tariffs, payment processing, enterprise plan billing, and payment gateway integration.

---

## 2. Architectural Layers

- **`domain/`**: Entities (`Invoice`, `Payment`, `Tariff`), Value Objects, Domain Errors (`errors.go`), and `repository.go` contracts.
- **`application/`**: Use cases & application services (`BillingService`).
- **`infrastructure/postgres/`**: Database repository (`BillingRepository`) querying tenant invoices and public subscriptions.
- **`api/http/`**: HTTP controllers (`billing_handler.go`) and Echo route bindings.
- **`module.go`**: Dependency container exporting `NewModule(server)`.

---

## 3. Database Schema & Tables Scope

- **Schema**: `tenant_<slug>` & `public`
- **Tables**: `invoices`, `invoice_items`, `payments`, `pricing_tariffs`.

---

## 4. Exported Public APIs

- `GET /api/v1/invoices`: List patient invoices (requires `billing:read`)
- `POST /api/v1/invoices`: Generate patient invoice (requires `billing:write`)

---

## 5. Testing
```bash
go test ./apps/backend/internal/modules/billing/...
```
