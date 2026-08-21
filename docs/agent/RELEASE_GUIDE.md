# Production Release Guide & Quality Gates Checklist

> **Purpose**: Master production release protocol for certifying Release Candidates of Curexal Healthcare Operating System.  
> **Owner**: Principal Systems Architect  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. 13 Production Quality Gates Sign-Off Protocol

Before tagging a release or marking a milestone complete, verify:

```text
[✓] Gate 1:  Backend compiles cleanly (go build ./...)
[✓] Gate 2:  Frontend applications compile cleanly (npm run build)
[✓] Gate 3:  Turbo build pipeline executes without errors (turbo run build)
[✓] Gate 4:  TypeScript strict mode checking passes (tsc --noEmit)
[✓] Gate 5:  Go backend tests pass (go test -v ./...)
[✓] Gate 6:  Frontend component & SDK tests pass
[✓] Gate 7:  Zero mock data or static JSON fallbacks in production paths
[✓] Gate 8:  Live end-to-end API integration verified (Component -> SDK -> Hertz -> DB)
[✓] Gate 9:  Documentation synchronized with actual source code implementation
[✓] Gate 10: Environment configuration (.env.example) synchronized across services
[✓] Gate 11: Security audit passed (Argon2id, HTTP-only cookies, Casbin RBAC, tenant isolation)
[✓] Gate 12: Bounded context architecture compliant (domain/, app/, infra/, api/)
[✓] Gate 13: Performance baseline met (zero N+1 queries, optimized bundle sizes)
```
