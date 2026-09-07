# CODEBASE DISCOVERY & INVENTORY PROMPT
**File**: `prompts/01-discovery/inspect-codebase.md`  

Execute a forensic discovery across the repository to map:
1. Applications (`apps/api`, `apps/web-platform`, `apps/web-patient`).
2. Packages (`packages/*`).
3. Database migrations (`migrations/*.sql`).
4. Endpoints and route handlers.
5. Domain services and repository layer.
6. RBAC guards, middlewares, and authorization engines.
7. Active test suites.

Produce a structured inventory classifying all modules against `/specs`.
