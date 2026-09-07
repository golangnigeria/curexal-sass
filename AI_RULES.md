# CUREXAL AI DEVELOPMENT CONTRACT & ENGINEERING RULES
**Document**: `AI_RULES.md`  
**Status**: CANONICAL AND IMMUTABLE  
**Scope**: All AI agents, contributors, and developers working on Curexal  

---

## 1. Specification Supremacy
`/specs` is the single source of product, domain, data, and security truth. Code, tests, and comments must strictly conform to `/specs`.

## 2. No Invented Business Rules
Never invent business logic, clinical rules, financial calculations, or tenant semantics. If requirements are not defined in `/specs`, stop and report the gap.

## 3. No Authorization Downgrades
Never weaken authorization, bypass permission checks, or loosen security boundaries to make a feature or build pass.

## 4. Absolute Tenant Isolation
Tenant isolation is a strict security boundary. Never trust client-supplied `organization_id`, `facility_id`, `branch_id`, or `user_id` without server-side cryptographic and database membership validation.

## 5. PHI Data Protection
Protected Health Information (PHI) must never be logged in plain text, leaked across tenants, exposed in URLs, or transmitted to unapproved external services.

## 6. Single Domain Representation
Never create duplicate representations of the same domain concept (e.g. `ClinicPatient` vs `LabPatient`). All domain entities must belong to their canonical module or shared kernel.

## 7. No Speculative Abstraction
Never add an abstraction simply because it appears elegant or might be useful in the distant future. Prefer simple, direct, and explicit code.

## 8. Minimal Compliant Implementation
Prefer the smallest, simplest compliant implementation that completely satisfies the specification and passes all acceptance criteria.

## 9. Delete Rather than Layer Over
Delete obsolete, dead, or bypassed code rather than layering new wrappers or adapters over bad logic.

## 10. Explicit Reversible Migrations
Database migrations must be explicit, strictly typed, relational, and reversible where practical.

## 11. Testing is Mandatory
Every feature, bug fix, or refactor must have corresponding unit, integration, or tenant-isolation regression tests.

## 12. Security Review for Sensitive Paths
All authorization, authentication, billing, patient identity, and clinical note modifications require explicit security verification.

## 13. Human-Governed Clinical Workflows
Clinical workflows require human-defined rules. Clinical decisions must always originate from authorized, credentialed practitioners.

## 14. No Silent AI Clinical Decision-Making
AI-assisted features (e.g. ICD-10 suggestions, triage scoring) must always be presented as assistive recommendations requiring explicit clinician acceptance, never silent autonomous actions.

## 15. Disclose Specification vs Code Conflicts
If existing code and the specification disagree, stop and identify the conflict, explain the risk, and present the resolution before making changes.

## 16. Ask Rather than Guess
If requirements, domain boundaries, or clinical state transitions are ambiguous, ask for clarification.

## 17. Intentional Backwards Compatibility Only
Preserve backwards compatibility only when it is explicitly required by active production clients or documented database migration policies.

## 18. No Sacred Legacy Code
No legacy code is protected merely because it exists. All code is subject to discovery, classification (KEEP / REFACTOR / REPLACE / DELETE), and systematic purging.

## 19. Strict Module Boundaries
No module may directly reach into another module's private database tables or internal services. Communication must occur through documented contracts, shared kernel models, or domain events.

## 20. Documented Architectural Decisions
Every consequential architectural, database, or security decision must have a documented Architecture Decision Record (ADR) in `/specs/decisions/`.
