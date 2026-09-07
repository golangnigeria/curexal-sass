# CUREXAL PLATFORM — RBAC & AUTHORIZATION SPECIFICATION
**Document**: `specs/01-platform/rbac.md`  
**Status**: APPROVED BASELINE  

---

## 1. Core Principle: Role $\neq$ Permission

Authorization in Curexal is **granular and permission-based**. Roles are simply named collections of granular permissions.

### Authorization Decision Chain
$$\text{Principal (User)} + \text{Domain (Org/Branch)} + \text{Resource (Object)} + \text{Action} + \text{Credentials} \implies \text{ALLOW / DENY}$$

---

## 2. Granular Permissions Dictionary

### Core / MPI Module
- `core.patient.create`: Register a new patient in the Master Patient Index.
- `core.patient.read`: View patient profile and demographic history.
- `core.patient.search`: Search MPI by name, phone, email, or MRN.
- `core.patient.update`: Edit patient demographic and emergency contact data.

### Clinic (HMS) Module
- `hms.appointment.create`: Schedule an outpatient appointment.
- `hms.queue.manage`: Check in patients and advance queue stages.
- `hms.triage.create`: Record patient intake vitals and allergies.
- `hms.consultation.create`: Open and write clinical encounter notes (SOAP).
- `hms.consultation.sign`: Legally sign and close an encounter note (requires `physician` credential).
- `hms.prescription.create`: Issue digital prescription orders.

### Billing Module
- `billing.invoice.read`: View patient billable itemization.
- `billing.payment.create`: Process cash, card, and POS checkout transactions.
- `billing.refund.create`: Process authorized financial adjustments.

### Organization & Governance
- `organization.governance`: Access `/organization/*` corporate settings.
- `organization.staff.manage`: Invite, assign, and update employee roles.
- `organization.billing.manage`: Manage corporate SaaS subscription and plan tiers.

---

## 3. Casbin Policy Engine Model

The Casbin model (`model.conf`) enforces multi-tenant domain authorization:
```ini
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && r.dom == p.dom && keyMatch2(r.obj, p.obj) && r.act == p.act
```
