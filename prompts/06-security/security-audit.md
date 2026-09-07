# SECURITY AUDIT & PHI PROTECTION PROMPT
**File**: `prompts/06-security/security-audit.md`  

Perform a full security and authorization audit:
1. Verify cryptographic JWT signature parsing and session invalidation.
2. Confirm that Casbin multi-tenant domain policies (`r = sub, dom, obj, act`) protect every sensitive endpoint.
3. Validate that clinical records (SOAP notes, diagnoses, e-prescriptions) require verified professional credentials (`physician`).
4. Ensure audit logging records all security-sensitive events without exposing sensitive PHI.
