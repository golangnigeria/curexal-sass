package testing

import (
	"context"
	"testing"

	"github.com/golangnigeria/curexal/internal/modules/subscription/application"
	"github.com/golangnigeria/curexal/internal/modules/subscription/domain"
	"github.com/google/uuid"
)

// MockEntitlementRepository for unit testing capability resolution logic
type MockEntitlementRepository struct {
	planCapabilities map[string][]string
	orgAddOns        map[string][]string
	dependencies     map[string][]string
	entitlements     map[string][]domain.OrganizationEntitlement
}

func NewMockEntitlementRepository() *MockEntitlementRepository {
	return &MockEntitlementRepository{
		planCapabilities: map[string][]string{
			"smart": {
				"core.organization", "core.patient", "core.customer_care", "core.billing", "core.audit", "core.documents",
				"laboratory.basic", "radiology.basic", "clinical.basic", "pharmacy.basic", "inventory.basic", "qms.basic",
			},
			"optimize": {
				"core.organization", "core.patient", "core.customer_care", "core.billing", "core.audit", "core.documents",
				"laboratory.basic", "radiology.basic", "clinical.basic", "pharmacy.basic", "inventory.basic", "qms.basic",
				"laboratory.analyzer_integration", "laboratory.advanced_qc",
			},
			"pro": {
				"core.organization", "core.patient", "core.customer_care", "core.billing", "core.audit", "core.documents",
				"laboratory.basic", "radiology.basic", "clinical.basic", "pharmacy.basic", "inventory.basic", "qms.basic",
				"laboratory.analyzer_integration", "laboratory.advanced_qc", "clinical.inpatient_wards", "pharmacy.advanced_inventory",
			},
		},
		orgAddOns: make(map[string][]string),
		dependencies: map[string][]string{
			"laboratory.analyzer_integration": {"laboratory.basic"},
			"laboratory.advanced_qc":          {"laboratory.basic"},
			"radiology.pacs_dicom":            {"radiology.basic"},
			"clinical.inpatient_wards":        {"clinical.basic"},
		},
		entitlements: make(map[string][]domain.OrganizationEntitlement),
	}
}

func (m *MockEntitlementRepository) GetPlanBaseCapabilities(ctx context.Context, planCode string) ([]string, error) {
	return m.planCapabilities[planCode], nil
}

func (m *MockEntitlementRepository) GetOrganizationAddOnCapabilities(ctx context.Context, orgID uuid.UUID) ([]string, error) {
	return m.orgAddOns[orgID.String()], nil
}

func (m *MockEntitlementRepository) GetCapabilityDependencies(ctx context.Context, capabilityCodes []string) ([]string, error) {
	var deps []string
	for _, code := range capabilityCodes {
		if dList, ok := m.dependencies[code]; ok {
			deps = append(deps, dList...)
		}
	}
	return deps, nil
}

func (m *MockEntitlementRepository) GetOrganizationEntitlements(ctx context.Context, orgID uuid.UUID) ([]domain.OrganizationEntitlement, error) {
	return m.entitlements[orgID.String()], nil
}

func (m *MockEntitlementRepository) GrantOrganizationEntitlement(ctx context.Context, entitlement *domain.OrganizationEntitlement) error {
	orgKey := entitlement.OrganizationID.String()
	m.orgAddOns[orgKey] = append(m.orgAddOns[orgKey], entitlement.CapabilityCode)
	m.entitlements[orgKey] = append(m.entitlements[orgKey], *entitlement)
	return nil
}

func (m *MockEntitlementRepository) RevokeOrganizationEntitlement(ctx context.Context, orgID uuid.UUID, capabilityCode string) error {
	orgKey := orgID.String()
	current := m.orgAddOns[orgKey]
	var updated []string
	for _, c := range current {
		if c != capabilityCode {
			updated = append(updated, c)
		}
	}
	m.orgAddOns[orgKey] = updated
	return nil
}

func (m *MockEntitlementRepository) GetCapabilityByCode(ctx context.Context, code string) (*domain.Capability, error) {
	return &domain.Capability{
		ID:       uuid.New(),
		Code:     code,
		Module:   "test",
		IsActive: true,
	}, nil
}

func (m *MockEntitlementRepository) GetAllCapabilities(ctx context.Context) ([]domain.Capability, error) {
	return []domain.Capability{
		{ID: uuid.New(), Code: "laboratory.basic", Name: "Basic Laboratory Operations", Module: "laboratory", IsActive: true},
		{ID: uuid.New(), Code: "laboratory.analyzer_integration", Name: "Analyzer Integration", Module: "laboratory", IsActive: true},
		{ID: uuid.New(), Code: "radiology.pacs_dicom", Name: "PACS DICOM", Module: "radiology", IsActive: true},
	}, nil
}

func (m *MockEntitlementRepository) GetCapabilityPrices(ctx context.Context, capabilityID uuid.UUID) ([]domain.CapabilityPrice, error) {
	return []domain.CapabilityPrice{
		{ID: uuid.New(), CapabilityID: capabilityID, Currency: "NGN", BillingPeriod: "monthly", Price: 25000, IsActive: true},
		{ID: uuid.New(), CapabilityID: capabilityID, Currency: "USD", BillingPeriod: "monthly", Price: 20, IsActive: true},
	}, nil
}

func (m *MockEntitlementRepository) CreateCapabilitySubscription(ctx context.Context, sub *domain.CapabilitySubscription) error {
	return nil
}

func (m *MockEntitlementRepository) GetCapabilitySubscription(ctx context.Context, subID uuid.UUID) (*domain.CapabilitySubscription, error) {
	return &domain.CapabilitySubscription{
		ID:     subID,
		Status: "active",
	}, nil
}

func (m *MockEntitlementRepository) UpdateCapabilitySubscriptionStatus(ctx context.Context, subID uuid.UUID, status string) error {
	return nil
}

// TestCapabilityResolutionEngine tests deterministic capability resolution
func TestCapabilityResolutionEngine(t *testing.T) {
	ctx := context.Background()
	repo := NewMockEntitlementRepository()
	svc := application.NewEntitlementService(nil, repo)

	orgID := uuid.New()

	// 1. Grant LIS add-on to Smart organization
	errGrantLIS := svc.GrantCapabilityAddOn(ctx, "", orgID, "laboratory.analyzer_integration", "purchase", nil)
	if errGrantLIS != nil {
		t.Fatalf("failed to grant LIS add-on: %v", errGrantLIS)
	}

	// 2. Grant RIS add-on to Smart organization
	errGrantRIS := svc.GrantCapabilityAddOn(ctx, "", orgID, "radiology.pacs_dicom", "purchase", nil)
	if errGrantRIS != nil {
		t.Fatalf("failed to grant RIS add-on: %v", errGrantRIS)
	}

	// 3. Verify HasCapability
	hasLIS, _ := svc.HasCapability(ctx, orgID, "laboratory.analyzer_integration")
	if !hasLIS {
		t.Errorf("expected organization to have LIS analyzer integration capability")
	}

	hasRIS, _ := svc.HasCapability(ctx, orgID, "radiology.pacs_dicom")
	if !hasRIS {
		t.Errorf("expected organization to have RIS PACS DICOM capability")
	}

	// 4. Test Revocation
	errRevoke := svc.RevokeCapabilityAddOn(ctx, "", orgID, "radiology.pacs_dicom")
	if errRevoke != nil {
		t.Fatalf("failed to revoke RIS add-on: %v", errRevoke)
	}

	hasRISAfterRevocation, _ := svc.HasCapability(ctx, orgID, "radiology.pacs_dicom")
	if hasRISAfterRevocation {
		t.Errorf("expected RIS capability to be false after revocation")
	}
}

// TestTrialCapabilityExpiration tests trial activation
func TestTrialCapabilityExpiration(t *testing.T) {
	ctx := context.Background()
	repo := NewMockEntitlementRepository()
	svc := application.NewEntitlementService(nil, repo)

	orgID := uuid.New()
	errTrial := svc.StartTrialCapability(ctx, "", orgID, "clinical.inpatient_wards", 30)
	if errTrial != nil {
		t.Fatalf("failed to start trial capability: %v", errTrial)
	}

	entitlements, errEnt := svc.GetOrganizationEntitlements(ctx, orgID)
	if errEnt != nil || len(entitlements) == 0 {
		t.Fatalf("expected active trial entitlement")
	}

	if entitlements[0].Source != "trial" {
		t.Errorf("expected source to be trial, got: %s", entitlements[0].Source)
	}
}

// TestSmartOrgAddOnImmutabilityAndTrace verifies Plan Immutability and Entitlement Tracing
func TestSmartOrgAddOnImmutabilityAndTrace(t *testing.T) {
	ctx := context.Background()
	repo := NewMockEntitlementRepository()
	svc := application.NewEntitlementService(nil, repo)

	// Org A: Smart base plan + LIS & RIS add-ons
	orgA := uuid.New()
	_ = svc.GrantCapabilityAddOn(ctx, "", orgA, "laboratory.analyzer_integration", "add_on", nil)
	_ = svc.GrantCapabilityAddOn(ctx, "", orgA, "radiology.pacs_dicom", "promotional_trial", nil)

	// Verify plan code remains "smart"
	planA, _ := svc.GetOrganizationPlan(ctx, orgA)
	if planA != "smart" {
		t.Errorf("expected Org A plan to remain 'smart', got: %s", planA)
	}

	// Verify Org A possesses LIS and RIS capabilities
	hasLIS, _ := svc.HasCapability(ctx, orgA, "laboratory.analyzer_integration")
	hasRIS, _ := svc.HasCapability(ctx, orgA, "radiology.pacs_dicom")
	if !hasLIS || !hasRIS {
		t.Errorf("expected Org A to possess both LIS and RIS capabilities")
	}

	// Org B: Smart base plan only (no LIS/RIS add-ons)
	orgB := uuid.New()
	hasLIS_B, _ := svc.HasCapability(ctx, orgB, "laboratory.analyzer_integration")
	hasRIS_B, _ := svc.HasCapability(ctx, orgB, "radiology.pacs_dicom")
	if hasLIS_B || hasRIS_B {
		t.Errorf("expected Org B to be DENIED LIS and RIS capabilities")
	}

	// Test Entitlement Trace for Org A
	traceA, errTrace := svc.GetEntitlementTrace(ctx, orgA, "laboratory.analyzer_integration")
	if errTrace != nil {
		t.Fatalf("failed to generate entitlement trace: %v", errTrace)
	}
	if !traceA.Allowed {
		t.Errorf("expected trace decision to be ALLOWED")
	}
	if traceA.BasePlan != "smart" {
		t.Errorf("expected trace base plan to be 'smart', got: %s", traceA.BasePlan)
	}

	// Verify trace contains step for dependency resolution
	foundDepStep := false
	for _, s := range traceA.Steps {
		if s.Step == "Dependency Resolution" {
			foundDepStep = true
			break
		}
	}
	if !foundDepStep {
		t.Errorf("expected entitlement trace to include Dependency Resolution step")
	}

	// Test Entitlement Trace for Org B (Denied)
	traceB, _ := svc.GetEntitlementTrace(ctx, orgB, "radiology.pacs_dicom")
	if traceB.Allowed {
		t.Errorf("expected trace decision for Org B to be DENIED")
	}
}

func TestCommercialCapabilitySubscriptionAndCatalog(t *testing.T) {
	ctx := context.Background()
	repo := NewMockEntitlementRepository()
	svc := application.NewEntitlementService(nil, repo)

	orgID := uuid.New()

	// 1. Query Capability Catalog with Pricing
	catalog, errCatalog := svc.GetCapabilityCatalog(ctx, orgID)
	if errCatalog != nil {
		t.Fatalf("failed to fetch capability catalog: %v", errCatalog)
	}
	if len(catalog) == 0 {
		t.Fatalf("expected capability catalog items")
	}

	// 2. Purchase Capability Add-on Subscription (NGN monthly)
	sub, errPurchase := svc.PurchaseCapabilityAddOn(ctx, "", orgID, "laboratory.analyzer_integration", "NGN", "monthly")
	if errPurchase != nil {
		t.Fatalf("failed to purchase capability add-on: %v", errPurchase)
	}
	if sub.Status != "active" {
		t.Errorf("expected subscription status to be active, got: %s", sub.Status)
	}
	if sub.Price != 25000 {
		t.Errorf("expected NGN price to be 25000, got: %f", sub.Price)
	}

	// 3. Verify Base Plan Immutability
	plan, _ := svc.GetOrganizationPlan(ctx, orgID)
	if plan != "smart" {
		t.Errorf("expected base plan code to remain 'smart', got: %s", plan)
	}

	// 4. Verify Effective Capability Immediate Availability
	hasCap, _ := svc.HasCapability(ctx, orgID, "laboratory.analyzer_integration")
	if !hasCap {
		t.Errorf("expected purchased capability to be immediately active")
	}
}
