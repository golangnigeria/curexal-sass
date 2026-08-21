# Clinical Bounded Context (`internal/modules/clinical`)

## 1. Purpose & Scope
The **Clinical Module** is responsible for patient registration, clinical visit tracking, medical laboratory test catalog management, order accessioning, and diagnostic test results.

---

## 2. Architectural Layers

- **`domain/`**: Entities (`Patient`, `Visit`, `CatalogItem`, `Order`, `Result`), Value Objects, Domain Errors (`errors.go`), and `repository.go` contracts.
- **`application/`**: Use cases & application services (`ClinicalApplicationService`).
- **`infrastructure/postgres/`**: Database repositories (`CatalogRepository`, `PatientRepository`) executing SQL queries on isolated `tenant_<slug>` schemas.
- **`api/http/`**: HTTP controllers (`catalog_handler.go`, `patient_visit_handler.go`) and Echo route bindings.
- **`module.go`**: Dependency container exporting `NewModule(server)`.

---

## 3. Database Schema & Tables Scope

- **Schema**: `tenant_<slug>` (isolated per branch facility)
- **Tables**: `patients`, `catalog_items`, `orders`, `results`, `patient_visits`.

---

## 4. Exported Public APIs

- `POST /api/v1/patients/visit`: Register new clinical patient visit (requires `patients:write`)
- `GET /api/v1/catalog`: List medical test catalog items
- `POST /api/v1/catalog`: Create new catalog test item

---

## 5. Testing
```bash
go test ./apps/backend/internal/modules/clinical/...
```
