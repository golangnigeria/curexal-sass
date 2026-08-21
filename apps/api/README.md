# CUREXAL_BACKEND

Enterprise Platform Kernel for Curexal Healthcare Platform.

## Architecture

```
                 PostgreSQL
           (Single Source of Truth)
                     │
                     ▼
        Go Domain Application Services
                     │
                     ▼
            REST/OpenAPI Contract
                     │
                     ▼
          Generated Type-safe SDK
                     │
                     ▼
      Web / Mobile / Desktop Clients
```

## Documentation
- **Engineering Constitution**: [.agents/AGENTS.md](file://.agents/AGENTS.md)
- **API Engineering & Audit Guide**: [docs/backend/API_ENGINEERING_GUIDE.md](file://docs/backend/API_ENGINEERING_GUIDE.md)
- **OpenAPI 3.0 Specification**: [static/openapi.json](file://static/openapi.json)

## Quick Start
```bash
# Run server
go run ./cmd/CUREXAL
# OR via Taskfile
task run
```

## Governance
This repository operates under strict enterprise governance rules:
- **Backend as Product Kernel**: Backend owns 100% of state, business logic, multi-tenancy, and permissions.
- **Clean Bounded Context Isolation**: Independent module contexts (`identity`, `organization`, `workspace`, `authorization`, `patient`, `laboratory`, `clinical`, `pharmacy`, `inventory`, `radiology`, `billing`, `audit`).