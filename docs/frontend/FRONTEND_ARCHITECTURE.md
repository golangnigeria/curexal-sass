# Frontend Applications Architecture & Shared Packages

> **Purpose**: Technical specification of all 4 monorepo single-page web applications (`apps/*`) and 10 shared npm packages (`packages/*`).  
> **Owner**: Principal Frontend Architect  
> **Status**: APPROVED / VERIFIED  
> **Last Updated**: 2026-07-27  
> **Verification Criteria**: Audited from ground truth source code across `/apps` and `/packages`.

---

## 1. Web Applications Architecture Overview

```text
apps/
├── web-admin       # Platform Admin Portal (admin.curexal.com / Port 3003)
├── web-workspace   # Multi-Tenant Workspace (slug.curexal.com / Port 3002)
├── web-patient     # Patient Self-Service Portal (patient.curexal.com / Port 3001)
└── web-public      # Commercial Landing & Onboarding Portal (curexal.com / Port 3000)
```

---

## 2. 100% `@curexal/api-sdk` Live API Integration

All 4 frontend applications consume live Hertz Go backend endpoints strictly via `@curexal/api-sdk`:

```text
React View Component (e.g. DashboardView.tsx)
          │
          ▼
@curexal/api-sdk Client (platformApi.getDashboard())
          │
          ▼
Fetch Request (credentials: 'include', HTTP-only cookies)
          │
          ▼
Hertz Go Backend Router (/api/v1/platform/dashboard)
```

---

## 3. Shared Monorepo Packages Breakdown

1. `@curexal/api-sdk` (`packages/api-sdk`): Typed HTTP client class with automatic 401 unauthenticated session refresh interceptor and RFC7807 error parsing.
2. `@curexal/auth` (`packages/auth`): Auth React context provider managing login states and user profile sessions.
3. `@curexal/authorization` (`packages/authorization`): Client-side Casbin permission hooks (`useHasPermission`).
4. `@curexal/ui-core` (`packages/ui-core`): Shared atomic UI components (Buttons, Modals, StatCards, LoadingSkeletons, ErrorStates).
5. `@curexal/ui-healthcare` (`packages/ui-healthcare`): Clinical UI components (Specimen tracker, DICOM viewer, LOINC concept picker).
6. `@curexal/design-tokens` (`packages/design-tokens`): Curexal enterprise theme token definitions (`Primary #266210`, `Accent #90B800`, `Success #00E1E1`, `Dark #063B00`).
7. `@curexal/types` (`packages/types`): Centralized TypeScript DTOs, API responses, and problem detail interfaces.
8. `@curexal/i18n` (`packages/i18n`): Multi-language translation dictionaries (English, French, Spanish, Swahili, Hausa).
9. `@curexal/state` (`packages/state`): Global Zustand client state stores.
10. `@curexal/config` (`packages/config`): Shared monorepo build configs (Vite, ESLint, Tailwind CSS, TypeScript).
