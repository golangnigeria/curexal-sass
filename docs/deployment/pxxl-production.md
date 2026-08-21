# Curexal Production Deployment Guide (pxxl.app)

**Current Status:** `DEPLOYMENT READY — LIVE VERIFICATION PENDING`

---

## 1. System Architecture Overview

```text
                                curexal.space
                                      │
                        ┌─────────────┼─────────────┐
                        │             │             │
                        ▼             ▼             ▼
                   curexal.space  app.curexal   api.curexal
                        │             │             │
                     Static        Static        Native
                   React/Vite    React/Vite     Go/Echo
                  (web-public)  (web-platform)    (api)
                                                    │
                                             ┌──────┴──────┐
                                             ▼             ▼
                                          NeonDB          R2
                                       (PostgreSQL)    (Storage)
```

---

## 2. Services Configuration on pxxl.app

### Service A: Go API Backend (`api.curexal.space`)
- **Service Name:** `curexal-api`
- **Root Directory:** `apps/api`
- **Build Command:** `go build -ldflags="-w -s" -o bin/curexal-backend ./cmd/CUREXAL`
- **Start Command:** `./bin/curexal-backend`
- **Domain:** `api.curexal.space`
- **Health Check Endpoint:** `GET /health` (or `GET /status`) &rarr; HTTP 200
- **Environment Variables:**
  ```env
  PORT=8080
  CUREXAL_PRIMARY_ENV=production
  CUREXAL_SERVER_DOMAIN=curexal.space
  CUREXAL_SERVER_CORS_ALLOWED_ORIGINS=https://curexal.space,https://app.curexal.space

  # Database (NeonDB Serverless PostgreSQL)
  CUREXAL_DB_DSN=postgres://[user]:[password]@[endpoint].neon.tech/curexal?sslmode=require

  # Object Storage (Cloudflare R2)
  STORAGE_PROVIDER=s3
  S3_ENDPOINT=https://<ACCOUNT_ID>.r2.cloudflarestorage.com
  S3_BUCKET=curexal-documents
  S3_REGION=auto
  S3_ACCESS_KEY=<R2_ACCESS_KEY_ID>
  S3_SECRET_KEY=<R2_SECRET_ACCESS_KEY>
  S3_PUBLIC_URL=https://cdn.curexal.space

  # Ephemeral Cache
  CACHE_PROVIDER=memory
  CACHE_DEFAULT_TTL=10m

  # Authentication & Cross-Subdomain Cookies
  CUREXAL_AUTH_SECRET_KEY=<256_BIT_RANDOM_SECRET_KEY>
  CUREXAL_AUTH_COOKIE_DOMAIN=.curexal.space
  CUREXAL_AUTH_COOKIE_PATH=/
  CUREXAL_AUTH_COOKIE_SECURE=true
  CUREXAL_AUTH_COOKIE_HTTPONLY=true
  CUREXAL_AUTH_COOKIE_SAMESITE=Lax
  CUREXAL_AUTH_ALLOW_TEST_HEADERS=false

  # Integrations & Transactional Emails
  CUREXAL_INTEGRATION_RESEND_API_KEY=<RESEND_API_KEY>
  CUREXAL_INTEGRATION_EMAIL_FROM_NAME="Curexal Health"
  CUREXAL_INTEGRATION_EMAIL_FROM_ADDRESS=noreply@contact.curexal.space
  CUREXAL_INTEGRATION_EMAIL_APP_URL=https://curexal.space
  ```

---

### Service B: Public Marketing Site (`curexal.space`)
- **Service Name:** `curexal-web-public`
- **Root Directory:** `apps/web-public`
- **Build Command:** `bun run build` (or `npm run build`)
- **Publish / Output Directory:** `dist`
- **Domain:** `curexal.space` (and `www.curexal.space`)
- **Environment Variables:**
  ```env
  VITE_API_URL=https://api.curexal.space
  VITE_PORTAL_URL=https://app.curexal.space
  VITE_APP_DOMAIN=curexal.space
  VITE_ENV=production

  # Active Marketing Lead-Capture (Waitlist & Book Demo)
  VITE_SUPABASE_URL=https://xrqwupliqiuotzbkcrja.supabase.co
  VITE_SUPABASE_ANON_KEY=<SUPABASE_ANON_KEY>
  ```

---

### Service C: Platform Portal & Workspaces (`app.curexal.space`)
- **Service Name:** `curexal-web-platform`
- **Root Directory:** `apps/web-platform`
- **Build Command:** `bun run build` (or `npm run build`)
- **Publish / Output Directory:** `dist`
- **Domain:** `app.curexal.space`
- **Environment Variables:**
  ```env
  VITE_API_URL=https://api.curexal.space
  VITE_PUBLIC_URL=https://curexal.space
  VITE_APP_DOMAIN=curexal.space
  VITE_ENV=production
  ```

---

## 3. Database Configuration (NeonDB)

- **Provider:** Neon Serverless PostgreSQL
- **Connection String:**
  ```text
  postgres://[user]:[password]@[endpoint].neon.tech/curexal?sslmode=require
  ```
- **Migration Pipeline:**
  The Go binary embeds all 44 database migrations via Go `embed.FS` and executes them automatically upon server start via `database.Migrate()`.
- **Security Rule:** Never inject database credentials into frontend client environment variables.

---

## 4. Object Storage (Cloudflare R2)

- **Provider:** Cloudflare R2 (`STORAGE_PROVIDER=s3`)
- **Isolation:** Storage credentials reside exclusively in `curexal-api`. Frontends interact solely via backend proxy or time-limited presigned URLs.
- **Manual Live Smoke Test:**
  1. **Upload**: Issue `POST /api/v1/organization/documents/upload` with a test PDF/image.
  2. **Existence**: Check `GET /api/v1/organization/documents` to verify the document record and metadata.
  3. **Presigned URL**: Request download URL (`GET /api/v1/organization/documents/:id/download`).
  4. **Access Check**: Verify the signed URL returns `HTTP 200` with the file.
  5. **Protected Check**: Attempt to access the raw R2 object URL without SigV4 parameters and verify it returns `HTTP 403 / 404`.
  6. **Delete**: Issue `DELETE /api/v1/organization/documents/:id` and verify deletion from bucket.

---

## 5. Authentication, Cookies, CORS & CSRF

- **Cookie Scope:**
  - `CookieDomain: .curexal.space`
  - `Secure: true` (Requires HTTPS across all domains)
  - `HttpOnly: true` (Inaccessible to malicious JavaScript)
  - `SameSite: Lax`
- **CORS Allowed Origins:**
  - `https://curexal.space`
  - `https://app.curexal.space`
  - *(Wildcard `*` is strictly disallowed with credentials)*
- **Frontend API Client:** Configured with `withCredentials: true` and automatic `X-CSRF-Token` headers.

---

## 6. SPA Routing & Fallback

- **Rule:** All routes must fallback to `/index.html` with status `200` for client-side routing.
- **Buildpack Support:** [`apps/web-platform/public/_redirects`](file:///c:/Users/HomePC/Desktop/program/fullstack_Curexal/apps/web-platform/public/_redirects) is included (`/* /index.html 200`).
- **PXXL.APP RUNTIME VERIFICATION REQUIRED:**
  - If pxxl.app static hosting automatically parses `_redirects`, deep links (e.g. `https://app.curexal.space/organization/dashboard`) will work immediately.
  - If pxxl.app uses a custom static server setting, enable the **SPA Rewrite / Fallback** toggle in the pxxl.app UI (or add rewrite rule `try_files {path} /index.html`).

---

## 7. Recommended Deployment Order

```text
Step 1: Provision NeonDB PostgreSQL (Obtain DSN with ?sslmode=require)
  │
Step 2: Provision Cloudflare R2 Bucket & API Token (curexal-documents)
  │
Step 3: Deploy Go Backend Service (api.curexal.space) on pxxl.app
  │
Step 4: Verify Backend Health (https://api.curexal.space/health -> HTTP 200)
  │
Step 5: Deploy Public Web (curexal.space) on pxxl.app
  │
Step 6: Deploy Platform Web (app.curexal.space) on pxxl.app
  │
Step 7: Configure DNS Records (A / CNAME pointing to pxxl.app)
  │
Step 8: Verify SSL Certificates (HTTPS provisioned on all 3 domains)
  │
Step 9: Verify SPA Routing (Direct browser reload on subroutes)
  │
Step 10: Verify Authentication & Cross-Subdomain Cookies
  │
Step 11: Execute R2 Upload / Download Smoke Test
  │
Step 12: Complete End-to-End Production Verification Checklist
```

---

## 8. Production Smoke Test Checklist

Execute this checklist against the live deployment before marking the deployment as `PRODUCTION VERIFIED`:

- [ ] `https://curexal.space` loads successfully (Public landing page).
- [ ] `https://app.curexal.space` loads successfully (Platform login/portal).
- [ ] `https://api.curexal.space/health` returns `HTTP 200` with JSON status `"healthy"`.
- [ ] Login flow succeeds at `https://app.curexal.space/login`.
- [ ] `jwt` and `refresh_token` cookies have `HttpOnly`, `Secure`, and `Domain=.curexal.space`.
- [ ] `GET /api/v1/bootstrap` returns authenticated principal and organization profile.
- [ ] Direct page refresh on `https://app.curexal.space/organization/dashboard` returns `HTTP 200` (SPA routing check).
- [ ] Organization dashboard loads with correct theme variables.
- [ ] Specialized workspaces load correctly:
  - [ ] Laboratory workspace (`/workspace/laboratory`)
  - [ ] Clinical workspace (`/workspace/clinical`)
  - [ ] Hospital workspace (`/workspace/hospital`)
  - [ ] Radiology workspace (`/workspace/radiology`)
  - [ ] Pharmacy workspace (`/workspace/pharmacy`)
  - [ ] Billing workspace (`/workspace/billing`)
- [ ] CSRF protection header `X-CSRF-Token` is attached on mutating requests.
- [ ] R2 document upload succeeds and returns verified storage key.
- [ ] R2 document presigned download URL succeeds and downloads file.
- [ ] R2 raw storage endpoint rejects unauthenticated access.
- [ ] User logout clears cookies and redirects to login.
- [ ] NeonDB database migrations executed without errors.

---

## 9. Classification Status

| Milestone | Status | Criteria |
| :--- | :--- | :--- |
| **Monorepo Restructure** | **COMPLETE** | Clean directory hierarchy, single `.git`, Taskfile tooling |
| **Dockerless Architecture** | **COMPLETE** | Native Go binary + Static Vite SPAs |
| **Code & Unit Test Audit** | **COMPLETE** | Go storage/cache/auth tests pass, frontend builds pass |
| **Deployment Execution** | **READY** | All environment schemas, build scripts, and configs documented |
| **Live Verification** | **PENDING** | Requires live deployment on pxxl.app & DNS activation |
