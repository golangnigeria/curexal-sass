# 8-Step Log-Driven Debugging Protocol

> **Purpose**: Mandatory 8-step protocol for investigating, diagnosing, and fixing bugs, build failures, or runtime errors in Curexal V2.  
> **Owner**: Principal Systems Architect  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. The 8-Step Debugging Protocol

```text
Step 1: Reproduce Failure
            │
            ▼
Step 2: Collect Un-truncated Logs & Stack Tracing
            │
            ▼
Step 3: Trace HTTP & API Gateway Route
            │
            ▼
Step 4: Trace Database Query & Schema State
            │
            ▼
Step 5: Identify Verified Root Cause
            │
            ▼
Step 6: Execute Root-Cause Fix (No Symptom Hacking)
            │
            ▼
Step 7: Execute Regression Test Suite
            │
            ▼
Step 8: Update Documentation & Add Safeguards
```

---

## 2. Forbidden Debugging Hacks

> [!CAUTION]
> **Prohibited Workarounds**
> - NEVER wrap failing calls in silent `try/catch` or return fake 0-byte fallbacks to mask errors.
> - NEVER disable authentication middleware or bypass Casbin RBAC authorization to fix 401/403 errors.
> - NEVER comment out broken assertions or delete failing unit test files.
> - NEVER invent function names without viewing exact definitions in source files.
