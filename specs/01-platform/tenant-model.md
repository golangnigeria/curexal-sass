# CUREXAL PLATFORM — MULTI-TENANT HIERARCHY SPECIFICATION
**Document**: `specs/01-platform/tenant-model.md`  
**Status**: APPROVED BASELINE  

---

## 1. Canonical Hierarchy

```text
CUREXAL PLATFORM (Scope: platform)
       │
       └── ORGANIZATION / TENANT (Scope: organization)
               │
               ├── FACILITY / BRANCH A (Scope: workspace)
               │       ├── Outpatient Clinic
               │       └── Pharmacy
               │
               └── FACILITY / BRANCH B (Scope: workspace)
                       └── Diagnostic Laboratory
```

---

## 2. Invariants & Isolation Boundaries

1. **Organization Membership $\neq$ Executive Access**: Belongs to an organization does NOT grant access to `/organization/*`. Only explicit executive roles (`owner`, `org_admin`, `org_regional_manager`, `org_finance_manager`, `org_hr_manager`, `org_quality_manager`) possess organization-level privileges.
2. **Facility / Branch is a Strict Security Domain**: A user assigned to Branch A cannot access Branch B data or routes by modifying URL branch slugs or request headers.
3. **Server-Side Validation Only**: The server cryptographically validates identity and queries `organization.organization_memberships` and `organization.membership_branches` for every request.
4. **Fail Closed**: Inactive memberships, suspended organizations, or unassigned facility requests fail closed with HTTP 401/403.
