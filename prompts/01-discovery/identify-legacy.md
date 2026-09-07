# IDENTIFY LEGACY LOGIC & TECHNICAL DEBT PROMPT
**File**: `prompts/01-discovery/identify-legacy.md`  

Scan the codebase for obsolete patterns:
1. Hardcoded role checks (e.g. `role === "admin"`, `isDoctor`).
2. Legacy fallback privilege escalations (e.g. defaulting empty roles to `owner`).
3. Client-side-only security checks or mock data.
4. Unused endpoints, dead services, and unreferenced schemas.
5. In-memory role-to-permission maps that conflict with the PostgreSQL SSOT.

Classify each finding as: KEEP, REFACTOR, REPLACE, or DELETE.
