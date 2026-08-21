# Identity Bounded Context (`internal/modules/identity`)

## 1. Purpose & Scope
The **Identity Module** is responsible for user authentication, session management, Multi-Factor Authentication (MFA), password credentials, workspace invitations, user profiles, professional credentials, and Role-Based Access Control (RBAC) permissions.

---

## 2. Architectural Layers

- **`domain/`**: Entities (`User`, `Session`, `Membership`, `Role`, `Permission`, `Invitation`, `UserProfile`), Value Objects, Domain Errors (`errors.go`), and `repository.go` contracts.
- **`application/`**: Use cases & application services (`AuthService`, `InviteService`, `UserRoleService`).
- **`infrastructure/postgres/`**: Database repositories (`UserRepository`, `InviteRepository`) executing SQL queries on `public` schema.
- **`api/http/`**: HTTP controllers (`auth_handler.go`, `invite_handler.go`, `user_role_handler.go`) and Echo route bindings.
- **`module.go`**: Dependency container exporting `NewModule(server)`.

---

## 3. Database Schema & Tables Scope

- **Schema**: `public`
- **Tables**: `user`, `user_profiles`, `user_mfa`, `password_reset_tokens`, `accounts`, `session`, `verifications`, `memberships`, `roles`, `permissions`, `role_permissions`, `permission_overrides`, `invitations`.

---

## 4. Exported Public APIs

- `GET /api/v1/auth/csrf`: CSRF token generation
- `GET /api/v1/users/me`: Current user profile & active tenant context
- `POST /api/v1/users/active-tenant`: Switch active operational tenant
- `POST /api/v1/auth/accept-invite`: Public workspace invitation acceptance
- `GET /api/v1/users`: List tenant memberships (requires `users:read`)
- `POST /api/v1/memberships`: Add tenant member (requires `users:write`)

---

## 5. Testing
```bash
go test ./apps/backend/internal/modules/identity/...
```
