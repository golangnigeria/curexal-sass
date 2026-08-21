# Curexal V2 Definition of Done (DoD) Quality Standard

This document establishes the mandatory Definition of Done (DoD) checklist that every feature must satisfy before code is merged into `main` or deployed to production.

---

## The 14-Point Definition of Done Checklist

Every feature pull request must satisfy all 14 quality criteria:

1. **Business Rules Verified**: All business invariants and constraints declared in the module PRD (`docs/05-product-specifications/`) are strictly enforced.
2. **Clean Architecture Layering**: Domain models have zero external dependencies. Services orchestrate use cases cleanly.
3. **Database & Migrations**: PostgreSQL tables created via versioned SQL migrations. Schema-per-tenant isolation handled using `RunInTenantTx`.
4. **API & RFC 7807 Error Standard**: Endpoint adheres to standard REST payload structure and returns `application/problem+json` for error states.
5. **OpenAPI Specs Updated**: `docs/openapi.yaml` updated with request schemas, response examples, and error codes.
6. **NATS JetStream Events**: Domain events published with standard envelope JSON structures.
7. **Immutable Audit Logging**: Entity mutations emit hash-chained audit log entries in `internal/core/audit`.
8. **Casbin RBAC Permissions**: Endpoint protected by middleware authorization check for active user role and tenant context.
9. **Curexal Enterprise UI Complete**: All 5 UI states handled in frontend component (`Loading`, `Empty`, `Error`, `Success`, `Permission Denied`) adhering to Linear/Stripe design standards.
10. **Feature Flags Configured**: Module toggles defined in `docs/15-release-plans/feature-flags.md` and enforced in code.
11. **Comprehensive Tests Provided**:
    - Go Unit & Integration Tests (`go test ./...`) passing cleanly.
    - Postman collection updated with positive and negative test cases.
12. **Observability & Tracing**: Handler propagates `X-Correlation-ID` header and logs structured JSON.
13. **ADR Updated**: If new architectural patterns were introduced, an ADR is logged in `docs/14-adr/`.
14. **No Placeholders / Todo Comments**: Zero `// TODO`, mock fallbacks, or dummy data anywhere in the codebase.
