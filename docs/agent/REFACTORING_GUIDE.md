# Refactoring Guide & Legacy Code Migration Rules

> **Purpose**: Guidelines for refactoring existing codebase modules without introducing breaking API changes or regressions.  
> **Owner**: Principal Systems Architect  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. Refactoring Rules

1. **Zero External Behavioral Changes**: Refactoring MUST preserve existing REST request/response contracts and RFC 7807 problem details payloads.
2. **Execute Tests Before & After**: Run unit tests (`go test ./...`) prior to refactoring and immediately after to confirm zero regressions.
3. **Atomic Commits**: Keep refactoring commits separate from new feature implementations.
4. **Identify Affected Callers**: Use ripgrep code search (`grep_search`) to locate all invocation sites before modifying interface method signatures.
