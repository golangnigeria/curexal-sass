# End-to-End API Integration Guide

> **Purpose**: Standardized procedure for wiring frontend React applications to backend Hertz Go REST endpoints through `@curexal/api-sdk`.  
> **Owner**: Lead API Engineer  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. End-to-End Integration Trace Pipeline

```text
React View Component (e.g. DashboardView.tsx)
          │
          ▼
Custom Hook / TanStack Query (e.g. usePlatformDashboard.ts)
          │
          ▼
API SDK Client Method (platformApi.getDashboard())
          │
          ▼
Fetch Client HTTP Request (with credentials: 'include')
          │
          ▼
Hertz Middleware (RequirePlatformSession / SearchPathMiddleware)
          │
          ▼
Hertz Handler (PlatformHandler.GetDashboard)
          │
          ▼
Application Service (PlatformDashboardService.GetAggregatedMetrics)
          │
          ▼
Bun ORM Repository (BunPlatformRepository.GetSettings)
          │
          ▼
PostgreSQL Database Schema (`public` or `tenant_<slug>`)
```

---

## 2. Step-by-Step API Integration Example

### Step A: Define SDK Method (`packages/api-sdk/src/platform.ts`)
```typescript
export class PlatformApi {
  async getDashboard(): Promise<ApiResponse<DashboardMetricsResponse>> {
    return this.client.get<DashboardMetricsResponse>('/api/v1/platform/dashboard');
  }
}
```

### Step B: Consume in React View (`apps/web-admin/src/views/DashboardView.tsx`)
```typescript
const [metrics, setMetrics] = useState<DashboardMetricsResponse | null>(null);

useEffect(() => {
  platformApi.getDashboard().then((res) => {
    if (res.success && res.data) setMetrics(res.data);
  });
}, []);
```
