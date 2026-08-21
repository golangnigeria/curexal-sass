# 2. Choice of Core Backend Technology Stack

Date: 2026-07-17

## Status

Approved

## Context

Curexal requires a high-performance, scalable backend engine capable of handling secure multi-tenant diagnostic data, HL7/ASTM stream parsing, and real-time medical referral routing. We need to standardize on a backend programming language, HTTP framework, ORM, and database isolation strategy.

## Decision

We standardize on the following core backend stack:
1. **Programming Language:** Go 1.25+. Go offers compile-time safety, concurrency primitives (goroutines), fast startup, and low memory footprint.
2. **HTTP Framework:** CloudWeGo Hertz. Hertz is chosen over Gin or fiber for its superior high-performance network transport engine and customization options.
3. **ORM:** Bun ORM. Bun is selected over GORM because it is SQL-first, light, performs fast query mappings with minimal reflection overhead, and handles dynamic multi-tenant schema-prefixes natively.
4. **Database:** PostgreSQL with a **schema-per-tenant** design. Shared SSO/tenancy metadata is in the `public` schema; medical diagnostic records are isolated in individual tenant schemas (`tenant_<id>`).

## Consequences

- All backend code must be written in Go.
- Developers must use Bun's query builder instead of GORM's active record pattern.
- Database access middleware must dynamically resolve and set PostgreSQL `search_path` for every request to guarantee tenant data isolation.
