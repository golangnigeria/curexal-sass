# Curexal V2 API Contracts & Open Specification

This document details the REST API guidelines, authentication headers, error handling conventions, and JSON payload contracts for Curexal V2.

---

## 1. General Principles & Headers

All HTTP API endpoints follow RESTful standards and return JSON payloads.

### Standard Request Headers
- `Authorization`: `Bearer <jwt_access_token>`
- `X-Tenant-Slug`: `<tenant_slug>` (Required for all tenant-scoped endpoints)
- `X-Correlation-ID`: `<uuid>` (Passed or generated for end-to-end tracing)
- `Content-Type`: `application/json`

---

## 2. Error Response Standard (RFC 7807 Problem Details)

All API errors return `application/problem+json`:

```json
{
  "type": "https://curexal.com/errors/unauthorized",
  "title": "Permission Denied",
  "status": 403,
  "detail": "User lacks the 'referrals:accept' permission for the specified tenant branch.",
  "instance": "/api/v1/referrals/ref-99201/accept"
}
```

---

## 3. Key Endpoint Specifications

### 3.1. Create B2B Referral Order
- **HTTP Method**: `POST /api/v1/referrals`
- **Security**: Bearer Auth + `X-Tenant-Slug`
- **Request Body**:
```json
{
  "target_organization_id": "018e-4aef-92b1",
  "patient_id": "018e-4aef-8810",
  "test_codes": ["FBC", "LIPID"],
  "clinical_summary": "Persistent fatigue and elevated blood pressure."
}
```
- **Response (201 Created)**:
```json
{
  "referral_id": "ref_99201",
  "status": "Dispatched",
  "created_at": "2026-07-21T17:52:00Z"
}
```
