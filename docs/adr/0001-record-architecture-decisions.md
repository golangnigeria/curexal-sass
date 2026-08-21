# 1. Record Architecture Decisions

Date: 2026-07-17

## Status

Approved

## Context

In setting up the next-generation enterprise healthcare operating system (Curexal), we require a structured method to document significant architectural decisions, reasons behind them, and historical context. This ensures that new team members and future maintainers understand the system design tradeoffs.

## Decision

We will use Architecture Decision Records (ADRs) to document all major design choices. 
ADR files will:
- Be located in `docs/adr/`.
- Use a sequential numerical prefix: `NNNN-title-in-kebab-case.md`.
- Use standard markdown formatting.
- Follow the structure: Title, Date, Status, Context, Decision, Consequences.

## Consequences

- All developers must document significant changes (such as database migrations strategies, key libraries addition, security mechanisms) as ADRs before implementation.
- History of choices remains immutable and version-controlled inside the repository.
