# Frontend Implementation Guide & Component Rules

> **Purpose**: Standard Operating Procedure defining the implementation order for React applications (`apps/*`) and shared packages (`packages/*`).  
> **Owner**: Principal Frontend Architect  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. Required Component Implementation Sequence

```text
Step 1: TypeScript Interfaces (`packages/types`)
            │
            ▼
Step 2: SDK API Client Method (`packages/api-sdk`)
            │
            ▼
Step 3: React Custom Hook / TanStack Query (`apps/*/src/hooks/`)
            │
            ▼
Step 4: Atomic UI Components (`packages/ui-core` or `packages/ui-healthcare`)
            │
            ▼
Step 5: View Container Component (`apps/*/src/views/`)
            │
            ▼
Step 6: Route Registration & Permission Guards (`App.tsx`)
```

---

## 2. Mandatory Component State Requirements

Every view component MUST explicitly render:
1. **Loading State**: `<LoadingSkeleton rows={5} />` matching exact container bounds.
2. **Error State**: `<ErrorState error={error} onRetry={refetch} />` formatting RFC 7807 problem details.
3. **Empty State**: `<EmptyState title="No Records Found" icon="📂" />` when an empty array is returned.
4. **Data State**: Clean table/grid rendering consuming live `@curexal/api-sdk` responses.

---

## 3. Zero Mock Data Policy

- NEVER write static JSON fallbacks or dummy arrays directly inside production component views.
- If an API endpoint is not yet available on the backend, create the `@curexal/api-sdk` client method, mark the view as `BLOCKED` in `docs/implementation/PLATFORM_UI_IMPLEMENTATION.md`, and document the missing endpoint.
