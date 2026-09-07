# CUREXAL PLATFORM — COMMERCIAL ENTITLEMENTS & PLANS
**Document**: `specs/08-business/entitlements.md`  
**Status**: APPROVED BASELINE  

---

## 1. Commercial Entitlement Architecture

```text
Organization
     └── Commercial Subscription
             └── Plan (Smart, Optimize, Enterprise)
                     └── Entitlements & Licensed Capabilities
                             └── Usage Metering (SMS, Appointments, POS Interchange)
```

---

## 2. Invariants

1. **No Hardcoded Plan Strings in Business Logic**:
   - BAD: `if (organization.plan === "enterprise") { enableConsultation(); }`
   - GOOD: `if (organization.hasCapability("clinic.consultation")) { enableConsultation(); }`
2. **Capability Decoupling**: Business modules function independently of billing plan names.
3. **Usage Tracking**: Billable events (SMS appointment reminders, POS interchange) are recorded in `billing.usage_records`.
