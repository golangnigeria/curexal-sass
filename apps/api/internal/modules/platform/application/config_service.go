package application

import (
	"context"
	"fmt"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/platform/domain"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
)

type ConfigRepository interface {
	GetGeneralSettings(ctx context.Context) (*domain.PlatformGeneralSettings, error)
	UpdateGeneralSettings(ctx context.Context, settings *domain.PlatformGeneralSettings, updatedBy uuid.UUID) (*domain.PlatformGeneralSettings, error)
	GetSecurityPolicy(ctx context.Context) (*domain.IdentitySecurityPolicy, error)
	UpdateSecurityPolicy(ctx context.Context, policy *domain.IdentitySecurityPolicy, updatedBy uuid.UUID) (*domain.IdentitySecurityPolicy, error)
}

type PlatformConfigService struct {
	repo      ConfigRepository
	auditRepo auditDomain.AuditRepository
}

func NewPlatformConfigService(repo ConfigRepository, auditRepo auditDomain.AuditRepository) *PlatformConfigService {
	return &PlatformConfigService{
		repo:      repo,
		auditRepo: auditRepo,
	}
}

func (s *PlatformConfigService) IsPlatformAdmin(principal *middleware.AuthenticatedPrincipal) bool {
	if principal == nil {
		return false
	}
	if principal.Platform.IsPlatformAdmin || principal.Platform.IsSuperAdmin || principal.Platform.IsPlatformStaff {
		return true
	}
	if principal.Role == "super_admin" || principal.Role == "platform_admin" || principal.Role == "platform_staff" {
		return true
	}
	return false
}

func (s *PlatformConfigService) GetGeneralSettings(ctx context.Context) (*domain.PlatformGeneralSettings, error) {
	return s.repo.GetGeneralSettings(ctx)
}

func (s *PlatformConfigService) UpdateGeneralSettings(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	settings *domain.PlatformGeneralSettings,
) (*domain.PlatformGeneralSettings, error) {
	if !s.IsPlatformAdmin(principal) {
		return nil, domain.ErrUnauthorizedPlatformAdmin
	}

	// 1. Fetch current singleton row to support partial field updates
	current, err := s.repo.GetGeneralSettings(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve platform settings: %w", err)
	}

	// 2. Overlay incoming non-empty/provided fields onto current settings
	if settings.PlatformName != "" {
		current.PlatformName = settings.PlatformName
	}
	if settings.PlatformLogoURL != nil {
		current.PlatformLogoURL = settings.PlatformLogoURL
	}
	if settings.PlatformFaviconURL != nil {
		current.PlatformFaviconURL = settings.PlatformFaviconURL
	}
	if settings.PlatformDescription != nil {
		current.PlatformDescription = settings.PlatformDescription
	}
	if settings.SupportEmail != "" {
		current.SupportEmail = settings.SupportEmail
	}
	if settings.SupportPhone != "" {
		current.SupportPhone = settings.SupportPhone
	}
	if settings.SupportWebsite != nil {
		current.SupportWebsite = settings.SupportWebsite
	}
	if settings.DefaultCountry != "" {
		current.DefaultCountry = settings.DefaultCountry
	}
	if settings.DefaultCurrency != "" {
		current.DefaultCurrency = settings.DefaultCurrency
	}
	if len(settings.SupportedCountries) > 0 {
		current.SupportedCountries = settings.SupportedCountries
	}
	if len(settings.SupportedCurrencies) > 0 {
		current.SupportedCurrencies = settings.SupportedCurrencies
	}
	if settings.DefaultTimezone != "" {
		current.DefaultTimezone = settings.DefaultTimezone
	}
	if settings.DefaultLocale != "" {
		current.DefaultLocale = settings.DefaultLocale
	}
	if settings.DateFormat != "" {
		current.DateFormat = settings.DateFormat
	}
	if settings.TimeFormat != "" {
		current.TimeFormat = settings.TimeFormat
	}
	if settings.NumberFormat != "" {
		current.NumberFormat = settings.NumberFormat
	}
	if settings.MeasurementUnits != "" {
		current.MeasurementUnits = settings.MeasurementUnits
	}
	if settings.MaintenanceMode {
		current.MaintenanceMode = true
	}
	if settings.AnnouncementBanner != nil {
		current.AnnouncementBanner = settings.AnnouncementBanner
	}
	if settings.Status != "" {
		current.Status = settings.Status
	}
	if settings.Version > 0 {
		current.Version = settings.Version
	}

	// 3. Validate merged entity
	if err := current.Validate(); err != nil {
		return nil, err
	}

	var actorUUID uuid.UUID
	if principal != nil && principal.UserID != "" {
		actorUUID, _ = uuid.Parse(principal.UserID)
	}

	updated, err := s.repo.UpdateGeneralSettings(ctx, current, actorUUID)
	if err != nil {
		return nil, err
	}

	if s.auditRepo != nil {
		action := "PLATFORM_CONFIG_UPDATED"
		resType := "platform.general_settings"
		resID := updated.ID.String()
		eventCat := "ADMINISTRATIVE"
		severity := "HIGH"
		status := "SUCCESS"
		actorID := principal.UserID
		actorName := principal.Identity.FullName

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    true,
			ActorID:       &actorID,
			ActorName:     &actorName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	return updated, nil
}

func (s *PlatformConfigService) GetSecurityPolicy(ctx context.Context) (*domain.IdentitySecurityPolicy, error) {
	return s.repo.GetSecurityPolicy(ctx)
}

func (s *PlatformConfigService) UpdateSecurityPolicy(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	policy *domain.IdentitySecurityPolicy,
) (*domain.IdentitySecurityPolicy, error) {
	if !s.IsPlatformAdmin(principal) {
		return nil, domain.ErrUnauthorizedPlatformAdmin
	}

	// 1. Fetch current singleton row to support partial field updates
	current, err := s.repo.GetSecurityPolicy(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve security policy: %w", err)
	}

	// 2. Overlay incoming non-zero/provided fields onto current policy
	if policy.MinPasswordLength > 0 {
		current.MinPasswordLength = policy.MinPasswordLength
	}
	if policy.PasswordExpirationDays > 0 {
		current.PasswordExpirationDays = policy.PasswordExpirationDays
	}
	if policy.MFAPolicy != "" {
		current.MFAPolicy = policy.MFAPolicy
	}
	if policy.MaxFailedLoginAttempts > 0 {
		current.MaxFailedLoginAttempts = policy.MaxFailedLoginAttempts
	}
	if policy.AccountLockoutDurationMinutes > 0 {
		current.AccountLockoutDurationMinutes = policy.AccountLockoutDurationMinutes
	}
	if policy.SessionMaxDurationHours > 0 {
		current.SessionMaxDurationHours = policy.SessionMaxDurationHours
	}
	if policy.RefreshTokenDurationDays > 0 {
		current.RefreshTokenDurationDays = policy.RefreshTokenDurationDays
	}
	if policy.MaxActiveSessions > 0 {
		current.MaxActiveSessions = policy.MaxActiveSessions
	}
	if policy.PasswordRequireUppercase {
		current.PasswordRequireUppercase = true
	}
	if policy.PasswordRequireNumber {
		current.PasswordRequireNumber = true
	}
	if policy.PasswordRequireSymbol {
		current.PasswordRequireSymbol = true
	}
	if policy.EmailVerificationRequired {
		current.EmailVerificationRequired = true
	}
	if policy.Version > 0 {
		current.Version = policy.Version
	}

	// 3. Validate merged policy
	if err := current.Validate(); err != nil {
		return nil, err
	}

	var actorUUID uuid.UUID
	if principal != nil && principal.UserID != "" {
		actorUUID, _ = uuid.Parse(principal.UserID)
	}

	updated, err := s.repo.UpdateSecurityPolicy(ctx, current, actorUUID)
	if err != nil {
		return nil, err
	}

	if s.auditRepo != nil {
		action := "PLATFORM_SECURITY_POLICY_UPDATED"
		resType := "platform.identity_security_policies"
		resID := updated.ID.String()
		eventCat := "SECURITY"
		severity := "CRITICAL"
		status := "SUCCESS"
		actorID := principal.UserID
		actorName := principal.Identity.FullName

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    true,
			ActorID:       &actorID,
			ActorName:     &actorName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	return updated, nil
}

