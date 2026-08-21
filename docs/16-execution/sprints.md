# Curexal V2 Master Sprint Schedule & Execution Hierarchy

This document defines the execution hierarchy (Roadmap $\rightarrow$ Epic $\rightarrow$ Sprint $\rightarrow$ Feature $\rightarrow$ Task) and Sprint schedule for Curexal V2.

---

## 1. Execution Granularity Hierarchy

```text
Roadmap (Long-term product vision: Phases 1 to 6)
   │
   ▼
Epic (Major capability work package: EPIC-001 Identity Platform)
   │
   ▼
Sprint (2 to 3-day deliverable unit: Sprint 1A - User Directory)
   │
   ▼
Feature (Single daily work item: Feature - Registration Endpoint)
   │
   ▼
Task (Specific atomic code change right now)
```

- **Roadmap**: Tells us _where we are going_.
- **Epic**: Tells us _what major capability we are building_.
- **Sprint**: Tells us _what we will finish next (in 2-3 days)_.
- **Feature**: Tells us _what to implement today_.
- **Task**: Tells the AI coding assistant _what exact code block to write right now_.

---

## 2. Sprint Schedule Breakdown

### EPIC-000: Platform Kernel & Shared Infrastructure (Baseline v1.0)

- [x] **Sprint 1 (Primitives & Errors)**: `primitive` (Clock, ULID), `errors` (RFC 7807 Mapper), `logging` (JSON logger).
- [x] **Sprint 2 (Database & Transactions)**: `database` (PostgreSQL/Bun pool), `transaction` (TransactionManager), `cache` (Redis).
- [x] **Sprint 3 (Security & Validation)**: `security` (Argon2id, AES, HMAC), `validation` (Struct Validator), `storage` (S3 Adapter).
- [x] **Sprint 4 (Health, Middleware & Manifest)**: `health` (Health Probes), `middleware` (Correlation, Headers, Recovery), `PLATFORM_MANIFEST.md`.

### EPIC-001: Identity Platform

- [ ] **Sprint 1 (User Directory & Credentials)**: User entity, Argon2 password hashing, registration DTOs & handler.
- [ ] **Sprint 2 (Sessions & JWT Token Rotation)**: Session repository, JWT access token generation, 15-min token rotation endpoint.
- [ ] **Sprint 3 (Password Reset & Verification)**: Crypto reset token generation, verification emails, password reset confirm endpoint.

### EPIC-002: Organization & Branch Management

- [ ] **Sprint 1 (Organization Registry)**: Organization domain model, slug validator, workspace registration API.
- [ ] **Sprint 2 (Multi-Branch & Settings)**: Branch directory, working hours, logo/branding configuration endpoints.

### EPIC-003: Tenant Provisioning Engine

- [ ] **Sprint 1 (Schema Creator)**: Dynamic `CREATE SCHEMA tenant_<slug>;` worker and versioned migration runner.
- [ ] **Sprint 2 (Provisioning Status Monitor)**: Provisioning progress monitor endpoint and background status updater.
