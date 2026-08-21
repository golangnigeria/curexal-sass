# Curexal V2 REST API Endpoint Reference

> **Purpose**: Authoritative reference of all Hertz REST endpoints, authentication requirements, path parameters, request/response DTO payloads, and RFC7807 error formats.  
> **Owner**: Lead API Engineer  
> **Status**: APPROVED / VERIFIED  
> **Last Updated**: 2026-07-27  
> **Verification Criteria**: Audited from Hertz HTTP handlers across `internal/modules/*/api`.

---

## 1. RFC 7807 Problem Details Error Spec

All error responses return standard `application/json+problem` format:

```json
{
  "type": "https://curexal.com/errors/conflict",
  "title": "Resource Conflict",
  "status": 409,
  "detail": "Platform has already been bootstrapped. Initialization can only execute once.",
  "instance": "/api/v1/platform/bootstrap"
}
```

---

## 2. Platform Management Endpoints (`/api/v1/platform`)

| Method | Endpoint Path | Auth Required | Required Role / Permission | Description |
| :--- | :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/platform/dashboard` | Yes (Session) | `platform_owner` / `platform_admin` | Aggregated executive telemetry & metrics |
| `GET` | `/api/v1/platform/health` | Yes (Session) | `platform_owner` / `platform_devops` | Cluster health probes (DB, Redis, Workers) |

---

## 3. Identity & Authentication Endpoints (`/api/v1/auth`)

| Method | Endpoint Path | Auth Required | Description |
| :--- | :--- | :--- | :--- |
| `POST` | `/api/v1/auth/login` | No | Credential authentication & session cookie issuance |
| `POST` | `/api/v1/auth/register` | No | Universal SSO User Registration |
| `POST` | `/api/v1/auth/refresh` | No (Refresh Cookie) | Session rotation & access token refresh |
| `POST` | `/api/v1/auth/logout` | Yes | Session revocation & cookie eviction |
| `GET` | `/api/v1/auth/me` | Yes | Identity profile & `default_destination` routing |
| `POST` | `/api/v1/auth/password-reset` | No | Request password reset token & email dispatch |
| `POST` | `/api/v1/auth/password-reset/confirm` | No | Confirm reset token & set new password |
| `POST` | `/api/v1/auth/email-verification` | No | Request email verification token |
| `POST` | `/api/v1/auth/email-verification/confirm` | No | Confirm email token & activate user |

---

## 4. Commercial CRM Endpoints (`/api/v1/leads`)

| Method | Endpoint Path | Auth Required | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/leads` | Yes | Retrieve commercial leads pipeline |
| `POST` | `/api/v1/leads` | No | Book a Demo / Request Information submission |
| `PUT` | `/api/v1/leads/:id/status` | Yes | Update lead pipeline status stage |
| `PUT` | `/api/v1/leads/:id/convert` | Yes | Convert lead to customer organization |

---

## 5. Organization & Provisioning Endpoints (`/api/v1/organizations`, `/api/v1/provisioning`)

| Method | Endpoint Path | Auth Required | Description |
| :--- | :--- | :--- | :--- |
| `GET` | `/api/v1/organizations` | Yes | List registered customer organizations |
| `POST` | `/api/v1/organizations` | Yes | Create customer organization |
| `POST` | `/api/v1/organizations/:id/verify` | Yes | Approve compliance verification |
| `GET` | `/api/v1/provisioning/jobs` | Yes | List PostgreSQL schema runner jobs |
