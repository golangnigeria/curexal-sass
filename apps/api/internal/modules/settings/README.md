# Settings Bounded Context (`internal/modules/settings`)

## 1. Purpose & Scope
The **Settings Module** manages branch operational settings, custom organization branding and themes, user preference flags, and configuration history tracking.

---

## 2. Architectural Layers

- **`domain/`**: Entities (`BranchSetting`, `OrganizationSetting`, `UserPreference`), Value Objects, Domain Errors (`errors.go`), and `repository.go` contracts.
- **`application/`**: Use cases & application services (`SettingsApplicationService`).
- **`infrastructure/postgres/`**: Database repository (`BranchSettingsRepository`) executing queries on `tenant_<slug>` and `public` schemas.
- **`api/http/`**: HTTP controllers (`settings_handler.go`) and Echo route bindings.
- **`module.go`**: Dependency container exporting `NewModule(server)`.

---

## 3. Database Schema & Tables Scope

- **Schema**: `public` & `tenant_<slug>`
- **Tables**: `organization_settings`, `branch_settings`, `user_preferences`, `settings_history`.

---

## 4. Exported Public APIs

- `GET /api/v1/settings/branch`: Get operational branch settings
- `PUT /api/v1/settings/branch`: Update operational branch settings

---

## 5. Testing
```bash
go test ./apps/backend/internal/modules/settings/...
```
