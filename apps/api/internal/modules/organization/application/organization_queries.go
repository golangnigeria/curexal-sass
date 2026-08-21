package application

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

func (s *OrganizationApplicationService) authorizeOrgMemberOrPlatformStaff(ctx echo.Context, orgID uuid.UUID) error {
	principal := middleware.GetPrincipal(ctx)
	userID := ""
	if principal != nil {
		userID = principal.UserID
		if principal.Platform.IsPlatformStaff {
			return nil
		}
	} else {
		userID = middleware.GetUserID(ctx)
	}

	if userID == "" {
		return errs.NewForbiddenError("Access denied: user authentication required", false)
	}

	dbExec := s.server.DB.Conn(ctx.Request().Context())
	var platformRole *string
	var isPlatAdmin bool
	err := dbExec.QueryRow(ctx.Request().Context(), `SELECT is_platform_admin, platform_role FROM identity.users WHERE id = $1`, userID).Scan(&isPlatAdmin, &platformRole)
	if err == nil && (isPlatAdmin || (platformRole != nil && (*platformRole == "super_admin" || *platformRole == "super_sales_staff"))) {
		return nil
	}

	var isMember bool
	err = dbExec.QueryRow(ctx.Request().Context(), `
		SELECT EXISTS(
			SELECT 1 FROM organization.organization_memberships 
			WHERE organization_id = $1 AND user_id = $2
		)
	`, orgID.String(), userID).Scan(&isMember)
	if err != nil || !isMember {
		return errs.NewForbiddenError("Access denied: you are not a member of this organization", false)
	}

	return nil
}

func (s *OrganizationApplicationService) authorizeOrgOwnerOrPlatformStaff(ctx echo.Context, orgID uuid.UUID) error {
	principal := middleware.GetPrincipal(ctx)
	userID := ""
	if principal != nil {
		userID = principal.UserID
		if principal.Platform.IsPlatformStaff {
			return nil
		}
	} else {
		userID = middleware.GetUserID(ctx)
	}

	if userID == "" {
		return errs.NewForbiddenError("Access denied: user authentication required", false)
	}

	dbExec := s.server.DB.Conn(ctx.Request().Context())
	var platformRole *string
	var isPlatAdmin bool
	err := dbExec.QueryRow(ctx.Request().Context(), `SELECT is_platform_admin, platform_role FROM identity.users WHERE id = $1`, userID).Scan(&isPlatAdmin, &platformRole)
	if err == nil && (isPlatAdmin || (platformRole != nil && (*platformRole == "super_admin" || *platformRole == "super_sales_staff"))) {
		return nil
	}

	var isOwner bool
	err = dbExec.QueryRow(ctx.Request().Context(), `
		SELECT EXISTS(
			SELECT 1 FROM organization.organization_memberships 
			WHERE organization_id = $1 AND user_id = $2 AND (role IN ('owner', 'admin', 'org_admin') OR role_title = 'owner')
		)
	`, orgID.String(), userID).Scan(&isOwner)
	if err != nil || !isOwner {
		return errs.NewForbiddenError("Access denied: you must be the owner of the organization or platform staff to perform this action", false)
	}

	return nil
}

func (s *OrganizationApplicationService) GetOrganizationByID(ctx echo.Context, id uuid.UUID) (*domain.Organization, error) {
	if err := s.authorizeOrgMemberOrPlatformStaff(ctx, id); err != nil {
		return nil, err
	}

	org, err := s.orgRepo.GetByID(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewNotFoundError("Organization not found", false, nil)
		}
		return nil, errs.NewInternalServerError()
	}

	return org, nil
}

func (s *OrganizationApplicationService) UpdateOrganization(ctx echo.Context, id uuid.UUID, name, slug, plan, customDomain *string, settings map[string]any) (*domain.Organization, error) {
	if err := s.authorizeOrgOwnerOrPlatformStaff(ctx, id); err != nil {
		return nil, err
	}

	org, err := s.orgRepo.Update(ctx.Request().Context(), id, name, slug, plan, customDomain, settings)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errs.NewNotFoundError("Organization not found", false, nil)
		}
		return nil, errs.NewInternalServerError()
	}

	return org, nil
}

func (s *OrganizationApplicationService) GetOrganizationSettings(ctx echo.Context, orgID uuid.UUID) (*domain.OrganizationSettings, error) {
	if err := s.authorizeOrgMemberOrPlatformStaff(ctx, orgID); err != nil {
		return nil, err
	}

	settings, err := s.orgRepo.GetSettings(ctx.Request().Context(), orgID)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}

	return settings, nil
}

func (s *OrganizationApplicationService) UpdateOrganizationSettings(ctx echo.Context, orgID uuid.UUID, logoURL, themeBranding, customDomain, supportEmail, supportPhone, cacNumber, tinNumber, taxNumber, businessType, address, timezone, currency, language *string) (*domain.OrganizationSettings, error) {
	if err := s.authorizeOrgOwnerOrPlatformStaff(ctx, orgID); err != nil {
		return nil, err
	}

	settings, err := s.orgRepo.UpdateSettings(ctx.Request().Context(), orgID, logoURL, themeBranding, customDomain, supportEmail, supportPhone, cacNumber, tinNumber, taxNumber, businessType, address, timezone, currency, language)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}

	return settings, nil
}

func (s *OrganizationApplicationService) CreateBranch(ctx echo.Context, orgID uuid.UUID, passedUserID string, name, slug, location, phone, address string, logoURL, currency *string, modules []string) (*domain.Tenant, error) {
	principal := middleware.GetPrincipal(ctx)
	userID := passedUserID
	if userID == "" {
		if principal != nil {
			userID = principal.UserID
		} else {
			userID = middleware.GetUserID(ctx)
		}
	}
	dbExec := s.server.DB.Conn(ctx.Request().Context())

	isPlatformStaff := principal != nil && principal.Platform.IsPlatformStaff
	if !isPlatformStaff && userID != "" {
		var platformRole *string
		var isPlatAdmin bool
		err := dbExec.QueryRow(ctx.Request().Context(), `SELECT is_platform_admin, platform_role FROM identity.users WHERE id = $1`, userID).Scan(&isPlatAdmin, &platformRole)
		isPlatformStaff = err == nil && (isPlatAdmin || (platformRole != nil && (*platformRole == "super_admin" || *platformRole == "super_sales_staff")))
	}

	if !isPlatformStaff {
		var isOwner bool
		err := dbExec.QueryRow(ctx.Request().Context(), `
			SELECT EXISTS(
				SELECT 1 FROM organization.organization_memberships 
				WHERE organization_id = $1 AND user_id = $2 AND (role IN ('owner', 'admin', 'org_admin') OR role_title = 'owner')
			)
		`, orgID.String(), userID).Scan(&isOwner)
		if err != nil || !isOwner {
			return nil, errs.NewForbiddenError("Access denied: you must be the owner of the organization or platform staff to create branches under it", false)
		}
	}

	var orgPlan string
	err := dbExec.QueryRow(ctx.Request().Context(), `SELECT plan FROM organization.organizations WHERE id = $1`, orgID.String()).Scan(&orgPlan)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}

	var currentBranchCount int
	err = dbExec.QueryRow(ctx.Request().Context(), `SELECT COUNT(*) FROM organization.facility_branches WHERE organization_id = $1`, orgID.String()).Scan(&currentBranchCount)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}

	maxBranches := 1
	var limitsJSON []byte
	_ = dbExec.QueryRow(ctx.Request().Context(), `SELECT limits FROM subscription.plans WHERE code = $1`, strings.ToLower(orgPlan)).Scan(&limitsJSON)
	if len(limitsJSON) > 0 {
		var limitsMap map[string]int
		if errU := json.Unmarshal(limitsJSON, &limitsMap); errU == nil {
			if lim, ok := limitsMap["maxBranches"]; ok && lim > 0 {
				maxBranches = lim
			}
		}
	}

	if currentBranchCount >= maxBranches {
		planDisplayName := strings.ToUpper(orgPlan[:1]) + orgPlan[1:]
		reasonMsg := fmt.Sprintf("Branch limit reached. Your subscription plan (%s) allows a maximum of %d branches. Please contact support or purchase additional branch entitlement.", planDisplayName, maxBranches)
		return nil, errs.NewForbiddenError(reasonMsg, false)
	}

	cleanSlug := strings.ToLower(strings.TrimSpace(slug))

	reserved := map[string]bool{"admin": true, "api": true, "app": true, "www": true, "billing": true, "status": true}
	if reserved[cleanSlug] {
		return nil, errs.NewBadRequestError("Reserved subdomains cannot be used as workspace URLs", false, nil, nil, nil)
	}

	exists, err := s.tenantRepo.CheckSlugExists(ctx.Request().Context(), cleanSlug)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}
	if exists {
		return nil, errs.NewBadRequestError("Workspace URL (slug) is already taken", false, nil, nil, nil)
	}

	var createdTenant *domain.Tenant
	err = s.server.DB.RunInTx(ctx.Request().Context(), func(txCtx context.Context) error {
		var txErr error
		createdTenant, txErr = s.tenantRepo.CreateTenant(txCtx, userID, name, cleanSlug, orgID.String(), location, phone, address, logoURL, currency, modules)
		if txErr != nil {
			return txErr
		}
		return s.server.DB.ProvisionTenant(txCtx, createdTenant.Slug)
	})

	if err != nil {
		s.server.Logger.Error().Err(err).Str("org_id", orgID.String()).Str("slug", slug).Msg("failed to create branch tenant")
		return nil, errs.NewInternalServerError()
	}

	return createdTenant, nil
}

func (s *OrganizationApplicationService) ListBranches(ctx echo.Context, orgID uuid.UUID) ([]domain.Tenant, error) {
	branches, err := s.tenantRepo.ListBranchesByOrgID(ctx.Request().Context(), orgID)
	if err != nil {
		return nil, errs.NewInternalServerError()
	}

	return branches, nil
}

func (s *OrganizationApplicationService) ListOrganizations(ctx echo.Context, passedUserID string, platformRole string) ([]domain.Organization, error) {
	principal := middleware.GetPrincipal(ctx)
	userID := passedUserID
	if userID == "" {
		if principal != nil {
			userID = principal.UserID
		} else {
			userID = middleware.GetUserID(ctx)
		}
	}

	isPlatformAdmin := (principal != nil && principal.Platform.IsPlatformStaff) || platformRole == "super_admin" || platformRole == "super_sales_staff"

	if !isPlatformAdmin && userID != "" {
		dbExec := s.server.DB.Conn(ctx.Request().Context())
		var dbPlatformRole *string
		var isPlatAdmin bool
		_ = dbExec.QueryRow(ctx.Request().Context(), `SELECT is_platform_admin, platform_role FROM identity.users WHERE id = $1`, userID).Scan(&isPlatAdmin, &dbPlatformRole)
		if isPlatAdmin || (dbPlatformRole != nil && (*dbPlatformRole == "super_admin" || *dbPlatformRole == "super_sales_staff")) {
			isPlatformAdmin = true
		}
	}

	if !isPlatformAdmin && userID == "" {
		isPlatformAdmin = true
	}

	orgs, err := s.orgRepo.List(ctx.Request().Context(), userID, isPlatformAdmin)
	if err != nil {
		s.server.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to list organizations")
		return []domain.Organization{}, nil
	}
	if orgs == nil {
		return []domain.Organization{}, nil
	}

	return orgs, nil
}




func (s *OrganizationApplicationService) CreateDemoRequest(ctx context.Context, labName, contactName, email string, phone, dailyVolume, notes *string) (*domain.DemoRequest, error) {
	return s.demoRepo.Create(ctx, labName, contactName, email, phone, dailyVolume, notes)
}

func (s *OrganizationApplicationService) ListDemoRequests(ctx context.Context) ([]domain.DemoRequest, error) {
	return s.demoRepo.List(ctx)
}

func (s *OrganizationApplicationService) UpdateDemoRequest(ctx context.Context, id uuid.UUID, status, notes *string) (*domain.DemoRequest, error) {
	return s.demoRepo.Update(ctx, id, status, notes)
}
