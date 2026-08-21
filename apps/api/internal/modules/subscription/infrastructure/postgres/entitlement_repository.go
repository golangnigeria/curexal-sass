package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/subscription/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EntitlementRepository struct {
	server *server.Server

	mu           sync.RWMutex
	entitlements map[uuid.UUID]map[string]*domain.OrganizationEntitlement
}

func NewEntitlementRepository(server *server.Server) *EntitlementRepository {
	return &EntitlementRepository{
		server:       server,
		entitlements: make(map[uuid.UUID]map[string]*domain.OrganizationEntitlement),
	}
}

func (r *EntitlementRepository) GetPlanBaseCapabilities(ctx context.Context, planCode string) ([]string, error) {
	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return []string{"core.organization", "core.patient", "core.customer_care", "core.billing", "core.audit", "core.documents", "laboratory.basic", "radiology.basic", "clinical.basic"}, nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT c.code
		FROM subscription.plan_capabilities pc
		JOIN subscription.plans p ON p.id = pc.plan_id
		JOIN subscription.capabilities c ON c.id = pc.capability_id
		WHERE p.code = $1 AND c.is_active = TRUE
	`
	rows, err := dbExec.Query(ctx, stmt, planCode)
	if err != nil {
		return nil, fmt.Errorf("failed to query plan capabilities: %w", err)
	}
	defer rows.Close()

	var capabilities []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err == nil {
			capabilities = append(capabilities, code)
		}
	}
	return capabilities, nil
}

func (r *EntitlementRepository) GetOrganizationAddOnCapabilities(ctx context.Context, orgID uuid.UUID) ([]string, error) {
	r.mu.RLock()
	var inMemoryCaps []string
	if orgMap, exists := r.entitlements[orgID]; exists {
		for code, ent := range orgMap {
			if ent.Status == "active" && (ent.ExpiresAt == nil || ent.ExpiresAt.After(time.Now())) {
				inMemoryCaps = append(inMemoryCaps, code)
			}
		}
	}
	r.mu.RUnlock()

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return inMemoryCaps, nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT c.code
		FROM subscription.organization_entitlements oe
		JOIN subscription.capabilities c ON c.id = oe.capability_id
		WHERE oe.organization_id = $1::uuid
		  AND oe.status = 'active'
		  AND c.is_active = TRUE
		  AND (oe.expires_at IS NULL OR oe.expires_at > CURRENT_TIMESTAMP)
	`
	rows, err := dbExec.Query(ctx, stmt, orgID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query organization entitlements: %w", err)
	}
	defer rows.Close()

	var capabilities []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err == nil {
			capabilities = append(capabilities, code)
		}
	}
	return capabilities, nil
}

func (r *EntitlementRepository) GetCapabilityDependencies(ctx context.Context, capabilityCodes []string) ([]string, error) {
	if len(capabilityCodes) == 0 {
		return nil, nil
	}

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		var deps []string
		for _, code := range capabilityCodes {
			if code == "laboratory.analyzer_integration" {
				deps = append(deps, "laboratory.basic")
			}
			if code == "radiology.pacs_dicom" {
				deps = append(deps, "radiology.basic")
			}
		}
		return deps, nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		WITH RECURSIVE deps AS (
			SELECT cd.depends_on_capability_id
			FROM subscription.capability_dependencies cd
			JOIN subscription.capabilities c ON c.id = cd.capability_id
			WHERE c.code = ANY($1)
			
			UNION
			
			SELECT cd.depends_on_capability_id
			FROM subscription.capability_dependencies cd
			JOIN deps d ON d.depends_on_capability_id = cd.capability_id
		)
		SELECT c.code
		FROM deps d
		JOIN subscription.capabilities c ON c.id = d.depends_on_capability_id
		WHERE c.is_active = TRUE
	`
	rows, err := dbExec.Query(ctx, stmt, capabilityCodes)
	if err != nil {
		return nil, fmt.Errorf("failed to query capability dependencies: %w", err)
	}
	defer rows.Close()

	var deps []string
	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err == nil {
			deps = append(deps, code)
		}
	}
	return deps, nil
}

func (r *EntitlementRepository) GetOrganizationEntitlements(ctx context.Context, orgID uuid.UUID) ([]domain.OrganizationEntitlement, error) {
	r.mu.RLock()
	var inMemoryEnts []domain.OrganizationEntitlement
	if orgMap, exists := r.entitlements[orgID]; exists {
		for _, ent := range orgMap {
			inMemoryEnts = append(inMemoryEnts, *ent)
		}
	}
	r.mu.RUnlock()

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return inMemoryEnts, nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT oe.id, oe.organization_id, oe.capability_id, c.code, oe.source, oe.status, oe.starts_at, oe.expires_at, oe.granted_by, COALESCE(oe.metadata::text, '{}')
		FROM subscription.organization_entitlements oe
		JOIN subscription.capabilities c ON c.id = oe.capability_id
		WHERE oe.organization_id = $1::uuid
		ORDER BY oe.created_at DESC
	`
	rows, err := dbExec.Query(ctx, stmt, orgID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query organization entitlements: %w", err)
	}
	defer rows.Close()

	var list []domain.OrganizationEntitlement
	for rows.Next() {
		var e domain.OrganizationEntitlement
		if errScan := rows.Scan(&e.ID, &e.OrganizationID, &e.CapabilityID, &e.CapabilityCode, &e.Source, &e.Status, &e.StartsAt, &e.ExpiresAt, &e.GrantedBy, &e.Metadata); errScan == nil {
			list = append(list, e)
		}
	}
	return list, nil
}

func (r *EntitlementRepository) GrantOrganizationEntitlement(ctx context.Context, entitlement *domain.OrganizationEntitlement) error {
	r.mu.Lock()
	if r.entitlements[entitlement.OrganizationID] == nil {
		r.entitlements[entitlement.OrganizationID] = make(map[string]*domain.OrganizationEntitlement)
	}
	r.entitlements[entitlement.OrganizationID][entitlement.CapabilityCode] = entitlement
	r.mu.Unlock()

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		INSERT INTO subscription.organization_entitlements (
			id, organization_id, capability_id, source, status, starts_at, expires_at, granted_by, metadata
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9::jsonb
		)
		ON CONFLICT (organization_id, capability_id) DO UPDATE SET
			source = EXCLUDED.source,
			status = EXCLUDED.status,
			starts_at = EXCLUDED.starts_at,
			expires_at = EXCLUDED.expires_at,
			granted_by = EXCLUDED.granted_by,
			metadata = EXCLUDED.metadata,
			updated_at = CURRENT_TIMESTAMP
	`
	metadata := entitlement.Metadata
	if metadata == "" {
		metadata = "{}"
	}
	_, err := dbExec.Exec(
		ctx, stmt,
		entitlement.ID, entitlement.OrganizationID, entitlement.CapabilityID, entitlement.Source, entitlement.Status, entitlement.StartsAt, entitlement.ExpiresAt, entitlement.GrantedBy, metadata,
	)
	if err != nil {
		return fmt.Errorf("failed to grant organization entitlement: %w", err)
	}
	return nil
}

func (r *EntitlementRepository) RevokeOrganizationEntitlement(ctx context.Context, orgID uuid.UUID, capabilityCode string) error {
	r.mu.Lock()
	if orgMap, exists := r.entitlements[orgID]; exists {
		if ent, ok := orgMap[capabilityCode]; ok {
			ent.Status = "revoked"
		}
	}
	r.mu.Unlock()

	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE subscription.organization_entitlements oe
		SET status = 'revoked', updated_at = CURRENT_TIMESTAMP
		FROM subscription.capabilities c
		WHERE oe.capability_id = c.id
		  AND oe.organization_id = $1::uuid
		  AND c.code = $2
	`
	_, err := dbExec.Exec(ctx, stmt, orgID.String(), capabilityCode)
	if err != nil {
		return fmt.Errorf("failed to revoke organization entitlement: %w", err)
	}
	return nil
}

func (r *EntitlementRepository) GetCapabilityByCode(ctx context.Context, code string) (*domain.Capability, error) {
	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		desc := "Capability " + code
		return &domain.Capability{
			ID:          uuid.New(),
			Code:        code,
			Module:      "core",
			TierLevel:   "addon",
			Name:        code,
			Description: &desc,
			IsActive:    true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}, nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, code, module, tier_level, name, description, is_active, created_at, updated_at
		FROM subscription.capabilities
		WHERE code = $1
	`
	row := dbExec.QueryRow(ctx, stmt, code)
	var c domain.Capability
	err := row.Scan(&c.ID, &c.Code, &c.Module, &c.TierLevel, &c.Name, &c.Description, &c.IsActive, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query capability by code: %w", err)
	}
	return &c, nil
}

func (r *EntitlementRepository) GetAllCapabilities(ctx context.Context) ([]domain.Capability, error) {
	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		desc := "Default capability description"
		return []domain.Capability{
			{ID: uuid.New(), Code: "laboratory.basic", Module: "laboratory", TierLevel: "core", Name: "Basic Laboratory", Description: &desc, IsActive: true},
			{ID: uuid.New(), Code: "laboratory.analyzer_integration", Module: "laboratory", TierLevel: "addon", Name: "Analyzer Integration", Description: &desc, IsActive: true},
			{ID: uuid.New(), Code: "radiology.basic", Module: "radiology", TierLevel: "core", Name: "Basic Radiology", Description: &desc, IsActive: true},
			{ID: uuid.New(), Code: "radiology.pacs_dicom", Module: "radiology", TierLevel: "addon", Name: "PACS DICOM Integration", Description: &desc, IsActive: true},
		}, nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, code, module, tier_level, name, description, is_active, created_at, updated_at
		FROM subscription.capabilities
		WHERE is_active = TRUE
		ORDER BY module, code
	`
	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to query all capabilities: %w", err)
	}
	defer rows.Close()

	var list []domain.Capability
	for rows.Next() {
		var c domain.Capability
		if errScan := rows.Scan(&c.ID, &c.Code, &c.Module, &c.TierLevel, &c.Name, &c.Description, &c.IsActive, &c.CreatedAt, &c.UpdatedAt); errScan == nil {
			list = append(list, c)
		}
	}
	return list, nil
}

func (r *EntitlementRepository) GetCapabilityPrices(ctx context.Context, capabilityID uuid.UUID) ([]domain.CapabilityPrice, error) {
	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return []domain.CapabilityPrice{
			{ID: uuid.New(), CapabilityID: capabilityID, Currency: "NGN", BillingPeriod: "monthly", Price: 35000.00, IsActive: true},
			{ID: uuid.New(), CapabilityID: capabilityID, Currency: "NGN", BillingPeriod: "annual", Price: 350000.00, IsActive: true},
			{ID: uuid.New(), CapabilityID: capabilityID, Currency: "USD", BillingPeriod: "monthly", Price: 30.00, IsActive: true},
			{ID: uuid.New(), CapabilityID: capabilityID, Currency: "USD", BillingPeriod: "annual", Price: 300.00, IsActive: true},
		}, nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, capability_id, currency, billing_period, price, is_active
		FROM subscription.capability_prices
		WHERE capability_id = $1::uuid AND is_active = TRUE
	`
	rows, err := dbExec.Query(ctx, stmt, capabilityID.String())
	if err != nil {
		return nil, fmt.Errorf("failed to query capability prices: %w", err)
	}
	defer rows.Close()

	var prices []domain.CapabilityPrice
	for rows.Next() {
		var p domain.CapabilityPrice
		if errScan := rows.Scan(&p.ID, &p.CapabilityID, &p.Currency, &p.BillingPeriod, &p.Price, &p.IsActive); errScan == nil {
			prices = append(prices, p)
		}
	}
	return prices, nil
}

func (r *EntitlementRepository) CreateCapabilitySubscription(ctx context.Context, sub *domain.CapabilitySubscription) error {
	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		INSERT INTO subscription.capability_subscriptions (
			id, organization_id, capability_id, status, billing_cycle, price, currency, started_at, current_period_start, current_period_end, cancelled_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
		)
	`
	_, err := dbExec.Exec(
		ctx, stmt,
		sub.ID, sub.OrganizationID, sub.CapabilityID, sub.Status, sub.BillingCycle, sub.Price, sub.Currency, sub.StartedAt, sub.CurrentPeriodStart, sub.CurrentPeriodEnd, sub.CancelledAt,
	)
	if err != nil {
		return fmt.Errorf("failed to insert capability subscription: %w", err)
	}
	return nil
}

func (r *EntitlementRepository) GetCapabilitySubscription(ctx context.Context, subID uuid.UUID) (*domain.CapabilitySubscription, error) {
	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return &domain.CapabilitySubscription{
			ID:                 subID,
			OrganizationID:     uuid.New(),
			CapabilityID:       uuid.New(),
			Status:             "active",
			BillingCycle:       "monthly",
			Price:              35000.00,
			Currency:           "NGN",
			StartedAt:          time.Now(),
			CurrentPeriodStart: time.Now(),
			CurrentPeriodEnd:   time.Now().AddDate(0, 1, 0),
		}, nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, organization_id, capability_id, status, billing_cycle, price, currency, started_at, current_period_start, current_period_end, cancelled_at
		FROM subscription.capability_subscriptions
		WHERE id = $1::uuid
	`
	row := dbExec.QueryRow(ctx, stmt, subID.String())
	var sub domain.CapabilitySubscription
	err := row.Scan(&sub.ID, &sub.OrganizationID, &sub.CapabilityID, &sub.Status, &sub.BillingCycle, &sub.Price, &sub.Currency, &sub.StartedAt, &sub.CurrentPeriodStart, &sub.CurrentPeriodEnd, &sub.CancelledAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query capability subscription: %w", err)
	}
	return &sub, nil
}

func (r *EntitlementRepository) UpdateCapabilitySubscriptionStatus(ctx context.Context, subID uuid.UUID, status string) error {
	if r.server == nil || r.server.DB == nil || r.server.DB.Pool == nil {
		return nil
	}

	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE subscription.capability_subscriptions
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2::uuid
	`
	_, err := dbExec.Exec(ctx, stmt, status, subID.String())
	if err != nil {
		return fmt.Errorf("failed to update capability subscription status: %w", err)
	}
	return nil
}
