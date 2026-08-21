# Facility Config Bounded Context (`internal/modules/facility_config`)

## 1. Purpose & Scope
The **Facility Config Module** manages medical branch facility configurations including department structures, room allocations, bed availability, equipment registries, and duty roster shifts.

---

## 2. Architectural Layers

- **`domain/`**: Entities (`Department`, `Room`, `Bed`, `Equipment`, `Shift`), Value Objects, Domain Errors (`errors.go`), and `repository.go` contracts.
- **`application/`**: Use cases & application services (`FacilityConfigService`).
- **`infrastructure/postgres/`**: Database repository (`FacilityConfigRepository`) executing SQL queries on `tenant_<slug>` schema.
- **`api/http/`**: HTTP controllers (`facility_config_handler.go`) and Echo route bindings.
- **`module.go`**: Dependency container exporting `NewModule(server)`.

---

## 3. Database Schema & Tables Scope

- **Schema**: `tenant_<slug>`
- **Tables**: `departments`, `rooms`, `beds`, `equipment`, `shifts`.

---

## 4. Exported Public APIs

- `GET /api/v1/facility/departments`: List branch departments
- `POST /api/v1/facility/departments`: Create department
- `GET /api/v1/facility/rooms`: List rooms and bed allocations

---

## 5. Testing
```bash
go test ./apps/backend/internal/modules/facility_config/...
```
