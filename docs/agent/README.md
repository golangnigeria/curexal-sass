# Curexal V2 Enterprise AI Agent Operating Playbook

> **Purpose**: Master operating instructions and Standard Operating Procedures (SOPs) for AI coding agents (and human pair-programmers) contributing to Curexal Healthcare Operating System (`curexalV2`).  
> **Owner**: Principal Systems Architect  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27  
> **Verification Criteria**: Enforced on 100% of pull requests, agent tasks, and feature implementations.

---

## 1. Core Operating Principles for AI Agents

1. **Source Code is Ground Truth**: Never assume a feature or API exists based on memory or outdated design docs. Inspect source code (`internal/modules/*`, `apps/*`, `packages/*`).
2. **Strict Layering & DDD Boundaries**: Every backend module MUST follow `domain/`, `app/`, `infra/`, `api/`, `module.go`. Handlers NEVER execute raw database queries. Repositories NEVER execute HTTP network calls.
3. **Zero Mock Data Policy**: Production frontend applications (`apps/*`) consume live Hertz backend endpoints exclusively via `@curexal/api-sdk`.
4. **Schema-per-Tenant Multi-Tenancy**: Clinical and financial data belong to tenant schemas (`tenant_<slug>`). Global platform governance data belongs to `public`.
5. **Cookie-First Session Strategy**: HttpOnly access cookies (15-minute TTL) with automated background refresh token rotation (7-day TTL).
6. **Decoupled Platform Management**: Platform Management & Bootstrap (`internal/modules/platform`) are isolated from Authentication (`internal/modules/auth`).

---

## 2. Agent Playbook Index

```text
docs/agent/
├── README.md                           # This Master Operating Guide
├── FEATURE_DEVELOPMENT_WORKFLOW.md     # 15-Step Pipeline from Requirement to Release
├── FEATURE_REQUEST_TEMPLATE.md         # Reusable Prompt Specification Template
├── FEATURE_BREAKDOWN_GUIDE.md          # Decomposing Large Epics (LIMS, RIS, EMR)
├── BACKEND_IMPLEMENTATION_GUIDE.md     # Layered Go Monolith Implementation Order
├── FRONTEND_IMPLEMENTATION_GUIDE.md    # React, SDK, Route & Component Order
├── API_INTEGRATION_GUIDE.md            # End-to-End Tracing (UI -> SDK -> Hertz -> DB)
├── DATABASE_WORKFLOW.md                # Migration-First Schema-per-Tenant Workflow
├── DEBUGGING_GUIDE.md                  # 8-Step Log-Driven Diagnostic Protocol
├── CODE_REVIEW_GUIDE.md                # 8-Dimension Architectural Review
├── REFACTORING_GUIDE.md                # Behavior-Preserving Refactoring Rules
├── TESTING_GUIDE.md                    # Healthcare QA, Unit, Integration & API Tests
├── DOCUMENTATION_GUIDE.md              # Living System Sync & Changelog Rules
├── RELEASE_GUIDE.md                    # Verification of 13 Quality Production Gates
├── PROMPT_LIBRARY.md                   # 14 Reusable Enterprise AI Prompts
└── CHECKLISTS.md                       # Task Execution Checklists
```

---

## 3. How to Execute an Implementation Task

1. **Read Request & Inspect Code**: Check `FEATURE_REQUEST_TEMPLATE.md` to clarify requirements. Inspect current module files before writing code.
2. **Execute Step-by-Step**: Follow `FEATURE_DEVELOPMENT_WORKFLOW.md`. Skip no layers.
3. **Verify Execution**: Run `go test ./...`, `go build ./...`, and `npm run build`. Confirm zero compiler or type errors.
4. **Update Documentation & Changelog**: Update `docs/project/CHANGELOG.md` and feature matrices upon completion.
 