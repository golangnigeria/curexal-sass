# 8-Dimension Code Review & Quality Assurance Guide

> **Purpose**: Code review standards for evaluating AI-generated code or human contributions prior to merging.  
> **Owner**: Lead Quality Assurance Architect  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. The 8 Dimensions of Code Review

1. **Architecture & DDD Alignment**: Are bounded context layers (`domain/`, `app/`, `infra/`, `api/`) strictly respected? Are repositories calling HTTP?
2. **Security & Authorization**: Is every endpoint protected by session authentication and Casbin RBAC? Are passwords hashed using Argon2id?
3. **Database Efficiency**: Are schema-per-tenant switches applied via `search_path` middleware? Are queries index-backed to prevent full table scans?
4. **Error Handling**: Are errors mapped to RFC 7807 problem details? Are internal database errors sanitized before reaching clients?
5. **Zero Mock Data Policy**: Does the frontend view consume live endpoints strictly via `@curexal/api-sdk`?
6. **Design System Tokens**: Are UI components using Curexal theme tokens (`#266210`, `#90B800`, `#00E1E1`, `#063B00`) rather than ad-hoc inline colors?
7. **Test Coverage**: Does the code include corresponding unit/integration tests?
8. **Documentation**: Is `docs/project/CHANGELOG.md` updated?
