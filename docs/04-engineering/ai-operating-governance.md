# Curexal V2 AI Operating Governance & Daily Workflow

This document defines the mandatory daily operating rhythm, 10-step prompt cycle, execution hierarchy, and architectural governance rules for building Curexal V2.

---

## 1. Execution Hierarchy

To maintain clarity and reviewability, work is structured across 5 levels of granularity:

```text
Roadmap (Long-term product vision: Phases 1 to 6)
   │
   ▼
Epic (Major capability work package: EPIC-001 Identity Platform)
   │
   ▼
Sprint (2 to 3-day deliverable unit: Sprint 1 - User Directory)
   │
   ▼
Feature (Single daily work item: Feature - Registration Endpoint)
   │
   ▼
Task (Specific atomic code change right now)
```

- **Roadmap**: Tells us *where we are going*.
- **Epic**: Tells us *what major capability we are building*.
- **Sprint**: Tells us *what we will finish next (in 2-3 days)*.
- **Feature**: Tells us *what single item to implement today* (**ONE feature at a time**).
- **Task**: Tells the AI coding assistant *what exact code block to write right now*.

---

## 2. The Daily Operating Rhythm (9:00 AM – 2:30 PM Workflow)

```text
09:00 AM ──► Pick ONE Single Feature (From active Sprint schedule)
   │
09:10 AM ──► Step 1: Execute Planning Prompt (Generate implementation_plan.md)
   │
09:30 AM ──► Step 2: Execute CTO Architecture Review (Audit plan against Clean Architecture & DDD)
   │
09:45 AM ──► Step 3: Approve Implementation Plan
   │
10:00 AM ──► Step 4: Execute Implementation Prompt (Write production code & tests)
   │
12:00 PM ──► Step 5: Execute Senior Staff Code Review (Review diff for security & correctness)
   │
12:30 PM ──► Step 6: Apply Refactorings & Fixes
   │
01:00 PM ──► Step 7: Execute Verification Suite (go test ./...)
   │
01:30 PM ──► Step 8: Commit & Merge Feature
   │
02:00 PM ──► Step 9: Synchronize Documentation & Update ENGINEERING_STATUS.md
   │
02:30 PM ──► Step 10: Plan Tomorrow's Feature & Stop
```

---

## 3. Step-by-Step Prompt Templates with Separation of Concerns

### STEP 1 — Planning Prompt (Pre-Code Execution)
- **Role**: Lead Product Owner & Technical Architect
- **Objective**: Establish business rules, dependency graph, risk register, and implementation plan before writing code.
- **Inputs**: Master Blueprint (`/docs`), Capability Map, Domain Models, OpenAPI specs.
- **Deliverable**: `implementation_plan.md` artifact.
- **Constraints**: **DO NOT WRITE CODE**. Do not modify files outside the artifact directory.
```text
We are implementing Feature: [Feature Name] under Sprint [X] of EPIC-[XXX].

Before writing any code:
1. Read every relevant document under /docs.
2. Identify every dependency and prerequisite.
3. Identify affected bounded contexts and Go packages.
4. Produce a detailed implementation plan artifact explicitly detailing:
   - **Transaction Boundaries**: What must be atomic within a database transaction.
   - **Idempotency Strategy**: Safe, repeatable retries for background operations & handlers.
   - **Failure Recovery**: How partial failures are handled gracefully.
   - **Observability**: Structured logs, metrics, correlation IDs (X-Correlation-ID), and events.
   - **Rollback / Compensation Strategy**: Reversion or compensation steps for failed workflows.
5. Verify Golden Bounded Context Rule for every new database table:
   - Which bounded context owns it?
   - Which aggregate is responsible for it?
   - Which service is the source of truth?
   - Which events can modify it?
   - Which events can observe it?
6. List technical and operational risks.
7. List architectural assumptions.
8. Explain how this implementation complies with the documented architecture.

Do NOT write code yet. Wait for approval.
```

---

### STEP 2 — Architecture Review Prompt (CTO Review)
- **Role**: Chief Technology Officer & Principal Architect
- **Objective**: Critically evaluate the implementation plan against Clean Architecture, DDD, multi-tenant isolation, and security standards.
- **Inputs**: The generated `implementation_plan.md` artifact and `/docs/03-architecture/`.
- **Deliverable**: Architectural Approval or Revision Recommendations.
- **Constraints**: **DO NOT WRITE CODE**. Focus purely on structural integrity.
```text
Review your implementation plan as if you are the Chief Architect.

Look for:
- Architecture violations
- Clean Architecture / Layer isolation violations
- DDD & Aggregate boundary violations
- SOLID principles violations
- Multi-tenant schema isolation risks
- Security & RBAC gaps
- Performance & connection pool concerns
- Future scalability bottlenecks
- Duplicate responsibilities across packages
- Missing documentation updates

Suggest concrete improvements. Do not write code yet.
```

---

### STEP 3 — Implementation Prompt (Coding Phase)
- **Role**: Staff Software Engineer
- **Objective**: Execute ONLY the approved implementation plan, producing production-ready Go backend and React/TS frontend code.
- **Inputs**: Approved `implementation_plan.md` artifact and `PLATFORM_MANIFEST.md`.
- **Deliverable**: Clean code files, versioned database migrations, unit & integration tests, OpenAPI specs.
- **Constraints**: Implement ONLY the approved plan scope. **ZERO PLACEHOLDERS** (`// TODO`), zero mock data, zero unannounced breaking changes.
```text
Implement ONLY the approved plan for Feature: [Feature Name].

Requirements:
- Follow Clean Architecture (`domain`, `app`, `api`, `infra`).
- Follow existing folder structure and package naming.
- Use `TransactionManager.RunInTenantTx` for all tenant database mutations.
- Do not introduce new unapproved patterns.
- Do not modify unrelated code files.
- Keep code production-ready with zero TODOs or placeholders.
- Add comprehensive unit and integration tests.
- Add structured JSON logging and OpenTelemetry context tracing.
- Format all error responses as RFC 7807 Problem Details.

Do not implement anything outside this task scope.
```

---

### STEP 4 — Self-Review & Code Audit Prompt
- **Role**: Senior Staff Code Auditor
- **Objective**: Perform an objective line-by-line code review of newly generated or modified files.
- **Inputs**: Git diff of modified/created files and `docs/04-engineering/development-guide.md`.
- **Deliverable**: Audit Findings Table with Severity (High/Medium/Low) and Remediation Actions.
- **Constraints**: Do not gloss over edge cases, unhandled errors, or missing test assertions.
```text
Perform a Senior Staff Engineer code review on the newly implemented code.

Review for:
1. Correctness & Business Invariant Enforcement
2. Layer Isolation & Clean Architecture Compliance
3. Tenant Schema Search Path Isolation & Connection Pool Safety
4. Memory Allocation & Performance Bottlenecks
5. Security (RBAC, OWASP, Input Validation, Encryption)
6. Concurrency & Data Race Conditions (`go test -race`)
7. Error Handling & RFC 7807 Compliance
8. Test Coverage & Assertion Rigor

List every issue found, assign severity (High / Medium / Low), and recommend exact fixes.
```

---

### STEP 5 — Debug & Incident Response Prompt
- **Role**: Production Incident Response Lead
- **Objective**: Diagnose test failures or runtime errors using empirical log evidence rather than guessing.
- **Inputs**: Full, un-truncated error log, stack trace, or test output.
- **Deliverable**: Root Cause Analysis (RCA) and Verified Fix.
- **Constraints**: **NEVER GUESS OR MUTATE CODE IMMEDIATELY**. Form ranked hypotheses first.
```text
Act as a Production Incident Response Engineer.

Do not immediately modify code.

First:
1. Explain the exact observed failing behavior.
2. Form multiple hypotheses.
3. Rank hypotheses by likelihood based on stack trace evidence.
4. Determine which logs, metrics, or DB states to inspect.
5. Determine which unit/integration tests to execute.
6. Isolate the exact root cause.

Only then recommend a targeted code fix. Never guess.
```

---

### STEP 6 — Architecture Audit Prompt
- **Role**: Governance Lead & Enterprise Auditor
- **Objective**: Audit the entire codebase against the master blueprint in `/docs` to catch architectural drift.
- **Inputs**: Whole codebase structure, `PLATFORM_MANIFEST.md`, `/docs/`.
- **Deliverable**: Architecture Compliance Scorecard (scored out of 10 per category).
```text
Audit the entire codebase against the documentation in /docs.

Check:
- Folder structure & package layout
- Naming conventions across Go packages and DB tables
- Bounded context boundaries & package responsibilities
- Import direction rules (Domain -> App -> Infra)
- API response consistency (RFC 7807)
- NATS event envelope consistency
- Database migration safety & search path isolation
- Security & RBAC enforcement

Generate a report. Score every category out of 10. Recommend structural improvements.
```

---

### STEP 7 — Technical Debt Audit Prompt
- **Role**: Principal Refactoring Engineer
- **Objective**: Identify code smells, duplicate logic, god objects, or circular dependencies across the repository.
- **Inputs**: Project source code files (`internal/` and `frontend/`).
- **Deliverable**: Technical Debt Remediation Backlog.
```text
Review the repository for technical debt.

Identify:
- Duplicate code or logic across packages
- Large files (>300 lines) or God objects
- Circular package dependencies
- Unused code, dead functions, or unreferenced DTOs
- Premature abstractions or over-engineering
- Missing error checks or swallowed exceptions

Suggest a prioritized remediation plan with low-risk refactoring steps.
```

---

### STEP 8 — End of Epic Completion Report Prompt
- **Role**: Engineering Lead & Delivery Manager
- **Objective**: Summarize deliverables, files created/modified, tests added, API changes, and verification checklist upon epic completion.
- **Inputs**: Completed epic codebase, test run results, OpenAPI updates.
- **Deliverable**: Epic Completion Report.
```text
Generate an Epic Completion Report for EPIC-[XXX].

Include:
1. Objectives achieved vs original scope
2. Files created (with file:// links)
3. Files modified (with file:// links)
4. Database migrations & schema changes
5. API endpoints added / updated
6. NATS domain events registered
7. Unit & Integration test coverage summary
8. Known issues & technical debt logged
9. Deployment readiness rating & rollback plan

Score the Epic completion (out of 10) and recommend the next Epic in dependency order.
```

---

### STEP 9 — Weekly CTO Review Prompt
- **Role**: Chief Technology Officer (CTO)
- **Objective**: Provide high-level executive review of product progress, team cadence, technical debt, and business risks.
- **Inputs**: `ENGINEERING_STATUS.md`, recent commit logs, open issues.
- **Deliverable**: Weekly CTO Assessment & Strategic Directives.
```text
You are now the CTO of Curexal.

Review the entire project status and answer:
1. What are we doing exceptionally well?
2. What technical or architectural debt is growing?
3. What security, multi-tenancy, or compliance risks exist?
4. What scalability bottlenecks exist?
5. What should we STOP doing immediately?
6. What should we IMPROVE in our workflow?
7. What is our next execution priority?

Be brutally honest.
```

---

### STEP 10 — Monthly 10 Million User Scalability Review Prompt
- **Role**: Principal Distributed Systems Architect
- **Objective**: Stress-test current system architecture against 10 million active users and high B2B referral throughput.
- **Inputs**: System Architecture Specs (`docs/03-architecture/`), C4 Models, Database Sharding Specs.
- **Deliverable**: 10M User Infrastructure Vulnerability & Scaling Plan.
```text
Assume Curexal has scaled to 10 million active patients and 50,000 healthcare provider facilities processing 1 million referrals per day.

Review the system architecture:
1. Would the current modular monolith and database design still work?
2. What component or bottleneck breaks FIRST under this load?
3. What operations become prohibitively expensive (compute/storage/bandwidth)?
4. What bounded contexts must be extracted into standalone microservices?
5. What should remain inside the modular monolith?
6. What architectural changes must we plan for the next 12 months?
```

---

## 4. The Golden Operating Rule

> **"Before answering or writing code, ask yourself: Does this move Curexal closer to the documented architecture, or is it merely solving today's problem? If it only solves today's problem, stop. Find the architectural solution instead. Always optimize for the platform, not the feature."**
