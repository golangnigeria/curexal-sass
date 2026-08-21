package application

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/subscription/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/google/uuid"
)

type EntitlementService struct {
	server *server.Server
	repo   domain.EntitlementRepository
}

func NewEntitlementService(server *server.Server, repo domain.EntitlementRepository) *EntitlementService {
	return &EntitlementService{
		server: server,
		repo:   repo,
	}
}

// GetEffectiveCapabilities resolves the deterministic merged capabilities for an organization:
// EffectiveCapabilities(orgID) = PlanBaseCapabilities + ActiveAddOnEntitlements + ResolvedDependencies
func (s *EntitlementService) GetEffectiveCapabilities(ctx context.Context, orgID uuid.UUID) ([]string, error) {
	// 1. Fetch active subscription plan code for organization
	var planCode string
	stmtPlan := `
		SELECT COALESCE(s.plan, o.plan, 'smart')
		FROM organization.organizations o
		LEFT JOIN subscription.subscriptions s ON s.organization_id = o.id AND s.status = 'active'
		WHERE o.id = $1::uuid
		LIMIT 1
	`
	if s.server != nil && s.server.DB != nil && s.server.DB.Pool != nil {
		_ = s.server.DB.Pool.QueryRow(ctx, stmtPlan, orgID.String()).Scan(&planCode)
	}
	if planCode == "" {
		planCode = "smart"
	}

	// 2. Fetch base capabilities mapped to subscription plan
	baseCaps, errBase := s.repo.GetPlanBaseCapabilities(ctx, strings.ToLower(planCode))
	if errBase != nil {
		return nil, fmt.Errorf("failed to resolve plan base capabilities: %w", errBase)
	}

	// 3. Fetch active add-on entitlements
	addOnCaps, errAddOn := s.repo.GetOrganizationAddOnCapabilities(ctx, orgID)
	if errAddOn != nil {
		return nil, fmt.Errorf("failed to resolve add-on capabilities: %w", errAddOn)
	}

	// Merge unique capability codes
	capSet := make(map[string]bool)
	var allExplicit []string

	for _, c := range baseCaps {
		if !capSet[c] {
			capSet[c] = true
			allExplicit = append(allExplicit, c)
		}
	}
	for _, c := range addOnCaps {
		if !capSet[c] {
			capSet[c] = true
			allExplicit = append(allExplicit, c)
		}
	}

	// 4. Resolve capability dependencies
	deps, errDeps := s.repo.GetCapabilityDependencies(ctx, allExplicit)
	if errDeps == nil {
		for _, dep := range deps {
			if !capSet[dep] {
				capSet[dep] = true
				allExplicit = append(allExplicit, dep)
			}
		}
	}

	return allExplicit, nil
}

func (s *EntitlementService) HasCapability(ctx context.Context, orgID uuid.UUID, capabilityCode string) (bool, error) {
	effective, err := s.GetEffectiveCapabilities(ctx, orgID)
	if err != nil {
		return false, err
	}
	for _, c := range effective {
		if c == capabilityCode {
			return true, nil
		}
	}
	return false, nil
}

func (s *EntitlementService) GrantCapabilityAddOn(
	ctx context.Context,
	grantedByUserID string,
	orgID uuid.UUID,
	capabilityCode string,
	source string,
	expiresAt *time.Time,
) error {
	if source == "" {
		source = "purchase"
	}

	cap, errCap := s.repo.GetCapabilityByCode(ctx, capabilityCode)
	if errCap != nil || cap == nil {
		return errs.NewNotFoundError(fmt.Sprintf("capability '%s' not found", capabilityCode))
	}

	var grantedByUUID *uuid.UUID
	if grantedByUserID != "" {
		parsed, errP := uuid.Parse(grantedByUserID)
		if errP == nil {
			grantedByUUID = &parsed
		}
	}

	entitlement := &domain.OrganizationEntitlement{
		ID:             uuid.New(),
		OrganizationID: orgID,
		CapabilityID:   cap.ID,
		CapabilityCode: cap.Code,
		Source:         source,
		Status:         "active",
		StartsAt:       time.Now(),
		ExpiresAt:      expiresAt,
		GrantedBy:      grantedByUUID,
		Metadata:       "{}",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if err := s.repo.GrantOrganizationEntitlement(ctx, entitlement); err != nil {
		return err
	}

	// Audit Trail Entry
	if s.server != nil && s.server.DB != nil && s.server.DB.Pool != nil {
		stmtAudit := `
			INSERT INTO audit.audit_events (id, action, resource_type, resource_id, actor_id, organization_id, payload)
			VALUES ($1, 'CAPABILITY_GRANTED', 'CAPABILITY', $2, $3, $4, jsonb_build_object('capability', $5, 'source', $6))
		`
		var actorIDStr *string
		if grantedByUserID != "" {
			actorIDStr = &grantedByUserID
		}
		_, _ = s.server.DB.Pool.Exec(ctx, stmtAudit, uuid.New().String(), cap.ID.String(), actorIDStr, orgID.String(), capabilityCode, source)
	}

	return nil
}

func (s *EntitlementService) RevokeCapabilityAddOn(ctx context.Context, actorID string, orgID uuid.UUID, capabilityCode string) error {
	if err := s.repo.RevokeOrganizationEntitlement(ctx, orgID, capabilityCode); err != nil {
		return err
	}

	// Audit Trail Entry
	if s.server != nil && s.server.DB != nil && s.server.DB.Pool != nil {
		stmtAudit := `
			INSERT INTO audit.audit_events (id, action, resource_type, resource_id, actor_id, organization_id, payload)
			VALUES ($1, 'CAPABILITY_REVOKED', 'CAPABILITY', $2, $3, $4, jsonb_build_object('capability', $5))
		`
		var actorIDStr *string
		if actorID != "" {
			actorIDStr = &actorID
		}
		_, _ = s.server.DB.Pool.Exec(ctx, stmtAudit, uuid.New().String(), capabilityCode, actorIDStr, orgID.String(), capabilityCode)
	}

	return nil
}

func (s *EntitlementService) StartTrialCapability(ctx context.Context, actorID string, orgID uuid.UUID, capabilityCode string, durationDays int) error {
	if durationDays <= 0 {
		durationDays = 30
	}
	expiresAt := time.Now().AddDate(0, 0, durationDays)
	return s.GrantCapabilityAddOn(ctx, actorID, orgID, capabilityCode, "trial", &expiresAt)
}

func (s *EntitlementService) GetOrganizationPlan(ctx context.Context, orgID uuid.UUID) (string, error) {
	var planCode string
	stmtPlan := `
		SELECT COALESCE(s.plan, o.plan, 'smart')
		FROM organization.organizations o
		LEFT JOIN subscription.subscriptions s ON s.organization_id = o.id AND s.status = 'active'
		WHERE o.id = $1::uuid
		LIMIT 1
	`
	if s.server != nil && s.server.DB != nil && s.server.DB.Pool != nil {
		_ = s.server.DB.Pool.QueryRow(ctx, stmtPlan, orgID.String()).Scan(&planCode)
	}
	if planCode == "" {
		planCode = "smart"
	}
	return strings.ToLower(planCode), nil
}

func (s *EntitlementService) GetEntitlementTrace(ctx context.Context, orgID uuid.UUID, capabilityCode string) (*domain.EntitlementTrace, error) {
	planCode, _ := s.GetOrganizationPlan(ctx, orgID)
	trace := &domain.EntitlementTrace{
		OrganizationID:      orgID,
		BasePlan:            planCode,
		RequestedCapability: capabilityCode,
		Allowed:             false,
		Steps:               []domain.EntitlementTraceStep{},
	}

	baseCaps, _ := s.repo.GetPlanBaseCapabilities(ctx, planCode)
	hasInBase := false
	for _, c := range baseCaps {
		if c == capabilityCode {
			hasInBase = true
			break
		}
	}

	if hasInBase {
		trace.Steps = append(trace.Steps, domain.EntitlementTraceStep{
			Step:    "Base Plan Check",
			Source:  "plan",
			Status:  "ACTIVE",
			Details: fmt.Sprintf("Capability '%s' is provided by base plan '%s'", capabilityCode, planCode),
		})
		trace.Allowed = true
		return trace, nil
	}

	trace.Steps = append(trace.Steps, domain.EntitlementTraceStep{
		Step:    "Base Plan Check",
		Source:  "plan",
		Status:  "NOT_PROVIDED",
		Details: fmt.Sprintf("Base plan '%s' does not provide capability '%s'", planCode, capabilityCode),
	})

	entitlements, errEnt := s.repo.GetOrganizationEntitlements(ctx, orgID)
	if errEnt == nil {
		for _, e := range entitlements {
			if e.CapabilityCode == capabilityCode {
				if e.Status == "active" && (e.ExpiresAt == nil || e.ExpiresAt.After(time.Now())) {
					trace.Steps = append(trace.Steps, domain.EntitlementTraceStep{
						Step:    "Organization Add-on Entitlement Check",
						Source:  e.Source,
						Status:  "ACTIVE",
						Details: fmt.Sprintf("Active entitlement found (source: %s)", e.Source),
					})
					trace.Allowed = true

					// Trace dependency resolution
					deps, errDeps := s.repo.GetCapabilityDependencies(ctx, []string{capabilityCode})
					if errDeps == nil && len(deps) > 0 {
						trace.Steps = append(trace.Steps, domain.EntitlementTraceStep{
							Step:    "Dependency Resolution",
							Source:  "dependency",
							Status:  "RESOLVED",
							Details: fmt.Sprintf("Dependencies auto-resolved: %s", strings.Join(deps, ", ")),
						})
					}
					return trace, nil
				} else {
					trace.Steps = append(trace.Steps, domain.EntitlementTraceStep{
						Step:    "Organization Add-on Entitlement Check",
						Source:  e.Source,
						Status:  e.Status,
						Details: fmt.Sprintf("Entitlement exists but is inactive/expired (status: %s)", e.Status),
					})
				}
			}
		}
	}

	trace.Steps = append(trace.Steps, domain.EntitlementTraceStep{
		Step:    "Final Decision",
		Status:  "DENIED",
		Details: fmt.Sprintf("Organization '%s' does not possess active capability '%s'", orgID, capabilityCode),
	})

	return trace, nil
}

func (s *EntitlementService) GetCapabilityCatalog(ctx context.Context, orgID uuid.UUID) ([]domain.CapabilityCatalogItem, error) {
	capabilities, errCaps := s.repo.GetAllCapabilities(ctx)
	if errCaps != nil {
		return nil, errCaps
	}

	planCode, _ := s.GetOrganizationPlan(ctx, orgID)
	baseCaps, _ := s.repo.GetPlanBaseCapabilities(ctx, planCode)
	baseSet := make(map[string]bool)
	for _, c := range baseCaps {
		baseSet[c] = true
	}

	effectiveCaps, _ := s.GetEffectiveCapabilities(ctx, orgID)
	effectiveSet := make(map[string]bool)
	for _, c := range effectiveCaps {
		effectiveSet[c] = true
	}

	var catalog []domain.CapabilityCatalogItem
	for _, cap := range capabilities {
		prices, _ := s.repo.GetCapabilityPrices(ctx, cap.ID)
		deps, _ := s.repo.GetCapabilityDependencies(ctx, []string{cap.Code})

		catalog = append(catalog, domain.CapabilityCatalogItem{
			Capability:      cap,
			Prices:          prices,
			Dependencies:    deps,
			AlreadyIncluded: baseSet[cap.Code],
			IsEffective:     effectiveSet[cap.Code],
		})
	}
	return catalog, nil
}

func (s *EntitlementService) PurchaseCapabilityAddOn(
	ctx context.Context,
	actorID string,
	orgID uuid.UUID,
	capabilityCode string,
	currency string,
	billingCycle string,
) (*domain.CapabilitySubscription, error) {
	if currency == "" {
		currency = "NGN"
	}
	if billingCycle == "" {
		billingCycle = "monthly"
	}

	cap, errCap := s.repo.GetCapabilityByCode(ctx, capabilityCode)
	if errCap != nil || cap == nil {
		return nil, errs.NewNotFoundError(fmt.Sprintf("capability '%s' not found", capabilityCode))
	}

	prices, errPrices := s.repo.GetCapabilityPrices(ctx, cap.ID)
	var chosenPrice float64 = 0.0
	if errPrices == nil {
		for _, p := range prices {
			if strings.EqualFold(p.Currency, currency) && strings.EqualFold(p.BillingPeriod, billingCycle) {
				chosenPrice = p.Price
				break
			}
		}
	}

	now := time.Now()
	sub := &domain.CapabilitySubscription{
		ID:                 uuid.New(),
		OrganizationID:     orgID,
		CapabilityID:       cap.ID,
		Status:             "active",
		BillingCycle:       billingCycle,
		Price:              chosenPrice,
		Currency:           currency,
		StartedAt:          now,
		CurrentPeriodStart: now,
		CurrentPeriodEnd:   now.AddDate(0, 1, 0),
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repo.CreateCapabilitySubscription(ctx, sub); err != nil {
		return nil, fmt.Errorf("failed to create capability subscription: %w", err)
	}

	// Grant capability entitlement with explicit source tracking
	// Preservation Rule: organization.plan is NEVER mutated; base plan remains smart/current
	expiresAt := sub.CurrentPeriodEnd
	errGrant := s.GrantCapabilityAddOn(ctx, actorID, orgID, capabilityCode, "add_on", &expiresAt)
	if errGrant != nil {
		return nil, fmt.Errorf("failed to grant entitlement for capability subscription: %w", errGrant)
	}

	return sub, nil
}

func (s *EntitlementService) GetOrganizationEntitlements(ctx context.Context, orgID uuid.UUID) ([]domain.OrganizationEntitlement, error) {
	return s.repo.GetOrganizationEntitlements(ctx, orgID)
}
