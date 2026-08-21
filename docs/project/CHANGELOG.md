# Curexal Enterprise Version Changelog

All notable changes to Curexal Healthcare Operating System are documented in this file in accordance with [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) and [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [2.4.0] - 2026-07-27

### Added
- **Decoupled Platform Management Context** (`internal/modules/platform`): Platform Settings, 7 Governance Roles (`platform_owner` through `platform_devops`), Reference Catalogs, Marketplace configurations, and Commercial Demo Plans.
- **10-Step Modular Bootstrap Pipeline**: CLI execution via `cmd/bootstrap` with pre-flight system validation.
- **Dual Isolated Middleware Chains**: Isolated `RequirePlatformSession` and `RequirePlatformRBAC` middlewares operating in global scope.
- **Expanded Identity Endpoint (`GET /api/v1/auth/me`)**: Returns structured identity payload including `user`, `session`, `platform`, `memberships`, and `default_destination`.
- **Platform Administration Portal UI/UX (`apps/web-admin`)**: Enterprise Control Center featuring 12-section sidebar, TopNav with command palette (`Ctrl+K`), real-time StatusBar, and 100% `@curexal/api-sdk` integration.
- **Platform Dashboard Micro-Aggregator** (`GET /api/v1/platform/dashboard`): Aggregates telemetry across platform health, leads, verifications, active provisioning jobs, MRR/ARR, and audit events.
- **5-10 Year Enterprise Documentation Architecture**: Rebuilt documentation system across 17 subdirectories under `docs/`, 7 Architecture Decision Records (ADRs), project tracking framework, and verified codebase audit reports.

### Changed
- Infrastructure-only bootstrap pipeline (removed REST route `POST /api/v1/platform/bootstrap`).
- Refactored `cmd/bootstrap/main.go` and `cmd/server/main.go` to inject `PlatformModule`.
