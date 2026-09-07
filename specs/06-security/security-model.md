# CUREXAL PLATFORM — SECURITY & PHI PROTECTION SPECIFICATION
**Document**: `specs/06-security/security-model.md`  
**Status**: APPROVED BASELINE  

---

## 1. Zero-Trust Tenant Isolation

1. **Authentication Verification**: All protected endpoints require a verified cryptographic JWT or secure session cookie.
2. **Context Resolution**: Server validates that the user is an active member of the requested organization and branch.
3. **Cross-Tenant Prevention**:
   - `Organization A` user attempting to access `Organization B` data returns **HTTP 403 Forbidden**.
   - `Branch A` user attempting to access `Branch B` data returns **HTTP 403 Forbidden**.
4. **Fail-Closed Principle**: If a principal's roles, permissions, or branch memberships cannot be securely verified, the request is rejected immediately.

---

## 2. Protected Health Information (PHI) Invariants

1. **No PHI in URLs**: Never place patient names, phone numbers, or clinical diagnoses in query strings.
2. **No PHI in Audit Logs**: Audit events record Actor, Timestamp, Action, and Resource IDs without plain-text medical notes.
3. **Encrypted Transport**: All communication requires TLS 1.3 in production.
