# Quality Assurance & Testing Strategy Guide

> **Purpose**: Standard Operating Procedure for backend unit/integration tests and frontend testing in Curexal V2.  
> **Owner**: Lead Quality Assurance Engineer  
> **Status**: APPROVED / MANDATORY  
> **Last Updated**: 2026-07-27

---

## 1. Testing Pyramid Architecture

```text
               / \
              /   \          End-to-End User Journeys (Playwright / Cypress)
             / E2E \         - Customer Onboarding & Verification Flow
            /-------\
           /   API   \       Hertz REST API Integration Tests (Postman / Go httptest)
          / Integration\    - Login, Session Rotation, Lead Conversion
         /---------------\
        /  Unit & Domain   \  Go Unit Tests & React Component Tests (Vitest)
       /    Logic Suites    \ - Domain Rules, Casbin Enforcer, SDK Methods
      /---------------------\
```

---

## 2. Go Unit Test Example (`app/service_test.go`)

```go
func TestPlatformBootstrapPipeline_Execute(t *testing.T) {
	ctx := context.Background()
	pCtx := &pipeline.PipelineContext{
		Logger: log.Default(),
	}
	step := &pipeline.PrerequisitesStep{}
	err := step.Execute(ctx, pCtx)
	if err == nil {
		t.Fatalf("expected error when DB connection is nil, got nil")
	}
}
```
