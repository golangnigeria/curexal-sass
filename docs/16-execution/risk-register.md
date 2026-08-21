# Curexal V2 Enterprise Risk Register

This document tracks system risks, architectural vulnerabilities, likelihood, impact ratings, and proactive mitigation strategies.

---

## Risk Register Matrix

| Risk ID | Risk Description | Likelihood | Impact | Proactive Mitigation Strategy | Owner |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **RISK-001** | **PostgreSQL DB Connection Exhaustion**: High tenant schema count causes database connection pool starvation. | Medium | High | Enforce `SET LOCAL search_path` within short-lived Bun transactions. Implement PgBouncer transaction-mode connection pooling. | Lead Architect |
| **RISK-002** | **Inter-Tenant State Desynchronization**: Async NATS referral events dropped due to worker restarts. | Low | High | Use NATS JetStream at-least-once persistent streams. Run daily reconciliation background daemon to scan inter-tenant ledgers. | Backend Lead |
| **RISK-003** | **Analyzer Offline Disconnection**: Network outage inside diagnostic lab prevents real-time test transmission. | High | Medium | Curexal Edge Agent caches ASTM telemetry strings in local SQLite database and replays via mTLS WebSocket stream upon reconnect. | Edge Engineer |
| **RISK-004** | **Unsettled Cash Referral Commissions**: Laboratories collect cash from patients but refuse/fail to pay clinic commissions. | Medium | Medium | Track wallet balance using double-entry ledgers. Auto-pause laboratory marketplace referral intake if wallet drops below credit limit. | Financial Lead |
| **RISK-005** | **Slow Dynamic Schema Migrations**: Migration execution across 1,000+ tenant schemas delays deployment pipeline. | Medium | High | Parallelize migration runner across concurrent Go worker routines with schema version tracking tables in each schema. | DevOps Lead |
