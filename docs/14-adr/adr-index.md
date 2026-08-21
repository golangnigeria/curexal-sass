# Curexal V2 Architecture Decision Records (ADR Index)

This document indexes all major architectural decisions governing Curexal V2.

---

## Index of Architectural Decisions

| ADR ID | Title | Status | Date | Decision Summary |
| :--- | :--- | :--- | :--- | :--- |
| **ADR-0001** | Record Architecture Decisions | Accepted | 2026-07-21 | Use ADRs to track all significant technical choices. |
| **ADR-0002** | Schema-per-Tenant Isolation Model | Accepted | 2026-07-21 | Use PostgreSQL dynamic schemas with `SET LOCAL search_path` for data privacy. |
| **ADR-0003** | Choice of Backend Stack (Go + Hertz + Bun ORM) | Accepted | 2026-07-21 | Use Go 1.25+ with Hertz for high-performance HTTP APIs and Bun for fluent ORM. |
| **ADR-0004** | NATS JetStream Event Broker | Accepted | 2026-07-21 | Use NATS JetStream for asynchronous at-least-once inter-tenant event dispatch. |
| **ADR-0005** | On-Premises Curexal Edge Agent | Accepted | 2026-07-21 | Deploy lightweight Go binary locally in labs to parse ASTM/HL7 analyzer strings over serial/TCP. |
| **ADR-0006** | Turborepo Monorepo Layout | Accepted | 2026-07-21 | Organize frontend apps (`portal`, `workspace`, `patient`, `admin`) and packages (`ui`, `api`) in single Turborepo. |
