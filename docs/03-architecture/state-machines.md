# Curexal V2 Domain State Machine Specifications

This document defines the formal state machines and valid transitions for key domain entities in Curexal V2.

---

## 1. Referral Order State Machine

```mermaid
stateDiagram-v2
    [*] --> Draft
    Draft --> Dispatched: Dispatch Order
    Dispatched --> Accepted: Target Lab Accepts
    Dispatched --> Rejected: Target Lab Rejects
    Dispatched --> Cancelled: Clinician Cancels
    Dispatched --> Expired: 24h Timeout
    Accepted --> SpecimenCollected: Phlebotomist Scans Barcode
    SpecimenCollected --> Testing: Ingest Analyzer Telemetry
    Testing --> Completed: Consultant Releases Results
    Completed --> [*]
    Rejected --> [*]
    Cancelled --> [*]
    Expired --> [*]
```

---

## 2. Specimen State Machine

```mermaid
stateDiagram-v2
    [*] --> Collected: Phlebotomist Logs Specimen
    Collected --> InTransit: Dispatched to Central Lab
    InTransit --> ReceivedAtBench: Received by Lab Tech
    Collected --> ReceivedAtBench: Direct Local Collection
    ReceivedAtBench --> Processing: Assigned to Worklist
    Processing --> Analyzed: Analyzer Ingests Data
    Analyzed --> Rejected: Delta Check Failure / Clotted Sample
    Analyzed --> Disposed: Verification Complete
    Disposed --> [*]
    Rejected --> [*]
```

---

## 3. Financial Wallet State Machine

```mermaid
stateDiagram-v2
    [*] --> Active: Provisioned with Workspace
    Active --> Overdrawn: Balance Drops Below (-1 * CreditLimit)
    Overdrawn --> Active: Top-up Payment Cleared
    Active --> Frozen: Security / Dispute Hold
    Frozen --> Active: Dispute Resolved
    Frozen --> Terminated: Workspace Cancelled
    Terminated --> [*]
```
