# Documentation Synchronization & Changelog Maintenance Guide

> **Purpose**: SOP for updating project documentation whenever source code, API contracts, or environment variables are modified.  
> **Owner**: Technical Writer & Principal Architect  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. Documentation Update Triggers

| Codebase Action | Required Documentation Files to Update |
| :--- | :--- |
| **New Hertz API Endpoint Added** | `docs/api/API_REFERENCE.md`, `@curexal/api-sdk`, `docs/project/CHANGELOG.md` |
| **New Environment Variable Added** | `.env.example`, `docs/environment/ENVIRONMENT.md`, `internal/core/config/config.go` |
| **New Monorepo Package Created** | `docs/packages/PACKAGES.md`, `docs/audit/REPOSITORY_MAP.md` |
| **New Bounded Context Created** | `docs/backend/BACKEND_ARCHITECTURE.md`, `docs/architecture/BOUNDED_CONTEXTS.md` |
| **Feature Completed / Milestone Reached** | `docs/project/PROJECT_STATUS.md`, `docs/audit/FEATURE_MATRIX.md`, `docs/project/CHANGELOG.md` |
