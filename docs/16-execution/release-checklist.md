# Curexal V2 Production Release Checklist

This document details the pre-release checklist required before tagging and deploying any release version to production environments.

---

## Release Readiness Gate Checklist

### 1. Security & Compliance Gate
- [ ] Static Application Security Testing (SAST) run cleanly (`gosec ./...`).
- [ ] Dependency security vulnerability scan completed without high/critical flags (`govulncheck ./...`).
- [ ] OWASP Top 10 security review completed (CSRF, XSS, SQL injection protection verified).
- [ ] HTTPS & TLS 1.3 enforced everywhere.

### 2. Database & Migration Gate
- [ ] Database migrations tested on staging replica schemas.
- [ ] Rollback SQL scripts executed and verified on staging.
- [ ] Database indexes verified using `EXPLAIN ANALYZE` on heavy query paths.

### 3. Quality & Test Gate
- [ ] All unit and integration tests passing (`go test -v ./...`).
- [ ] End-to-End Postman API collection verified against staging server.
- [ ] Zero mock fallbacks or placeholder code verified in release branch.

### 4. Operations & Monitoring Gate
- [ ] Prometheus metrics endpoints (`/metrics`) verified.
- [ ] OpenTelemetry Jaeger request correlation tracing verified.
- [ ] NATS JetStream consumer queues healthy with zero dead-letter accumulation.
- [ ] Feature flag matrix configured and deployed to environment.
