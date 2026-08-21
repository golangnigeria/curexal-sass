# Curexal V2 C4 Architectural Model

This document defines the C4 Model diagrams (System Context, Containers, and Components) for Curexal V2.

---

## Level 1: System Context Diagram

```mermaid
graph TD
    UserPatient[Patient / Consumer] -->|Uses Patient Portal / WhatsApp| CurexalSystem[Curexal V2 Healthcare Operating System & Marketplace]
    UserDoctor[Clinician / Doctor] -->|Uses Clinic Workspace| CurexalSystem
    UserLab[Lab Scientist / Director] -->|Uses LIMS Workspace| CurexalSystem
    UserPharm[Pharmacist] -->|Uses Pharmacy POS Workspace| CurexalSystem
    
    CurexalSystem -->|Processes Payments| Paystack[Paystack / Flutterwave Gateway]
    CurexalSystem -->|Dispatches Notifications| WhatsApp[WhatsApp Business API]
    CurexalSystem -->|Submits Claims| NHIA[National Health Insurance NHIA]
```

---

## Level 2: Container Diagram

```mermaid
graph TD
    subgraph Monorepo Frontend Applications
        AppPortal[Public Portal - React/Vite]
        AppWorkspace[Organization Workspace - React/Vite]
        AppPatient[Patient Portal - React/Vite]
        AppAdmin[Platform Admin - React/Vite]
    end

    subgraph Go Hertz Backend Modular Monolith
        APIGateway[Hertz HTTP API Gateway]
        AuthModule[Identity Module]
        TenantModule[Tenant Provisioner]
        ReferralModule[Referral Engine]
        LIMSModule[LIMS Module]
        FinanceModule[Financial Ledger Engine]
    end

    subgraph On-Premises Lab LAN
        EdgeAgent[Curexal Edge Agent - Go]
        Analyzer[Sysmex/Mindray Instrument]
    end

    subgraph Infrastructure Persistence Layer
        PostgreSQL[(PostgreSQL - Schema per Tenant)]
        Redis[(Redis Cache & Session Store)]
        MinIO[(MinIO / AWS S3 Storage)]
        NATS[(NATS JetStream Event Broker)]
    end

    AppWorkspace -->|REST / HTTPS| APIGateway
    AppPatient -->|REST / HTTPS| APIGateway
    Analyzer -->|ASTM/HL7| EdgeAgent
    EdgeAgent -->|mTLS WebSockets| APIGateway

    APIGateway --> AuthModule
    APIGateway --> TenantModule
    APIGateway --> ReferralModule
    APIGateway --> LIMSModule

    LIMSModule -->|Publish Events| NATS
    ReferralModule -->|Publish Events| NATS
    NATS --> FinanceModule

    AuthModule --> Redis
    LIMSModule --> MinIO
    TenantModule --> PostgreSQL
    LIMSModule --> PostgreSQL
```

---

## Level 3: Component Diagram (Laboratory LIMS Module)

```mermaid
graph TD
    subgraph LIMS Module Boundary
        PhlebotomyHandler[Phlebotomy Handler]
        WorklistService[Bench Worklist Service]
        ValidationEngine[Scientific Validation Engine]
        PDFGenerator[PDF Report Generator]
        LIMSRepo[LIMS Bun Repository]
    end

    PhlebotomyHandler --> WorklistService
    WorklistService --> ValidationEngine
    ValidationEngine --> PDFGenerator
    WorklistService --> LIMSRepo
    LIMSRepo -->|SET LOCAL search_path| PostgreSQL[(Tenant Schema DB)]
```
