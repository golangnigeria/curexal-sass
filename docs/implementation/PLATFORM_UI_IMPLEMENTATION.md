# Curexal Platform Administration Portal — Enterprise UI/UX Implementation Matrix

This document is the authoritative single source of truth for the development, backend API integration, component status, and verification of the **Curexal Platform Administration Portal (`apps/web-admin`)**.

---

## 1. System Architecture & Standards

- **Application Target**: `apps/web-admin` (`admin.curexal.com` / `http://localhost:3003`)
- **Theme**: Curexal Enterprise Dark/Light Design System
  - Primary: `#266210`
  - Accent: `#90B800`
  - Success: `#00E1E1`
  - Dark: `#063B00`
- **Grid & Spacing**: 8px baseline spacing system with subtle shadows and 12px border radii.
- **Layout Architecture**: Desktop-First Workspace (`Top Navigation` + `12-Section Sidebar` + `Workspace Grid` + `Status Bar`).
- **Data Integration**: 100% `@curexal/api-sdk` connected to live Hertz Go backend endpoints. Zero fake data, zero mock JSON files.

---

## 2. Implementation Progress Tracking Matrix

| Module / Screen | Backend Endpoint(s) | Backend Status | UI Built | SDK Connected | Status |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Authentication & Routing** | `GET /api/v1/auth/me`, `POST /auth/login` | ✅ Ready | ✅ Completed | ✅ Connected | ✅ VERIFIED |
| **Platform Dashboard** | `GET /api/v1/platform/dashboard`, `GET /platform/health` | ✅ Ready | ✅ Completed | ✅ Connected | ✅ VERIFIED |
| **CRM Lead Pipeline** | `GET /api/v1/leads`, `POST /leads`, `PUT /leads/:id/status` | ✅ Ready | ✅ Completed | ✅ Connected | ✅ VERIFIED |
| **Customer Conversion Modal** | `POST /leads/:id/convert`, `POST /auth/invitation` | ✅ Ready | ✅ Completed | ✅ Connected | ✅ VERIFIED |
| **Verification Queue** | `GET /api/v1/organizations`, `PUT /verify` | ✅ Ready | ✅ Completed | ✅ Connected | ✅ VERIFIED |
| **Organization Details** | `GET /api/v1/organizations` | ✅ Ready | ✅ Completed | ✅ Connected | ✅ VERIFIED |
| **Provisioning Monitor** | `GET /api/v1/provisioning/jobs` | ✅ Ready | ✅ Completed | ✅ Connected | ✅ VERIFIED |
| **Platform Infrastructure Monitoring** | `GET /api/v1/platform/health` | ✅ Ready | ✅ Completed | ✅ Connected | ✅ VERIFIED |
| **Enterprise Audit Explorer** | `GET /api/v1/platform/audit` | ✅ Ready | ✅ Completed | ✅ Connected | ✅ VERIFIED |
| **Platform User Management** | `GET /api/v1/platform/users` | ✅ Ready | ✅ Completed | ✅ Connected | ✅ VERIFIED |
| **Marketplace Configurations** | `GET /api/v1/platform/marketplace` | ✅ Ready | ✅ Completed | ✅ Connected | ✅ VERIFIED |
| **Platform Settings & DevTools** | `GET /api/v1/platform/health` | ✅ Ready | ✅ Completed | ✅ Connected | ✅ VERIFIED |

---

## 3. Global Reusable UI Components (`apps/web-admin/src/components`)

- [x] **TopNav**: Enterprise header with search bar (`Ctrl+K`), system status indicator, notification bell, user profile drawer.
- [x] **Sidebar**: 12-section navigation with active indicator, keyboard shortcuts, collapsible state, and badges.
- [x] **StatusBar**: Real-time platform status bar displaying API latency, DB status, and active environment.
- [x] **StatCard**: High-density metric display featuring current value, trend indicator, status pill, loading skeleton, and error state.
- [x] **StatusBadge**: Accessible color-coded badge for health and verification statuses (Operational, Active, Pending, Degraded, Failed).
- [x] **LoadingSkeleton**: Skeleton loaders matching exact layout dimensions.
- [x] **ErrorState**: RFC7807 problem details viewer with retry call-to-action button.
- [x] **EmptyState**: Clean illustrations and actionable prompts when zero data is returned.
- [x] **DocumentViewerModal**: Medical license and accreditation document inspection modal with zoom/download tools.
- [x] **CustomerConversionModal**: 5-step wizard layout for customer conversion and invitation link generation.
