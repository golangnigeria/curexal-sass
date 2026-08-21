package application

import (
	"context"
	"fmt"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
)

type OrganizationBrandingService struct {
	brandingRepo domain.OrganizationBrandingRepository
	orgRepo      domain.OrganizationRepository
	auditRepo    auditDomain.AuditRepository
	masterKey    string
}

func NewOrganizationBrandingService(
	brandingRepo domain.OrganizationBrandingRepository,
	orgRepo domain.OrganizationRepository,
	auditRepo auditDomain.AuditRepository,
	masterKey string,
) *OrganizationBrandingService {
	if masterKey == "" {
		masterKey = "CUREXAL_DEFAULT_PLATFORM_AES_MASTER_KEY_2026"
	}
	return &OrganizationBrandingService{
		brandingRepo: brandingRepo,
		orgRepo:      orgRepo,
		auditRepo:    auditRepo,
		masterKey:    masterKey,
	}
}

func (s *OrganizationBrandingService) isPlatformAdmin(principal *middleware.AuthenticatedPrincipal) bool {
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

func (s *OrganizationBrandingService) resolveActiveOrgUUID(principal *middleware.AuthenticatedPrincipal) (uuid.UUID, error) {
	if principal == nil {
		return uuid.Nil, domain.ErrUnauthorizedTenantAccess
	}

	orgIDStr := principal.Organization.ActiveOrganizationID
	if orgIDStr == "" {
		orgIDStr = principal.OrganizationID
	}
	if orgIDStr == "" {
		orgIDStr = principal.TenantID
	}

	if orgIDStr == "" {
		return uuid.Nil, domain.ErrUnauthorizedTenantAccess
	}

	parsed, err := uuid.Parse(orgIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid active organization ID: %w", err)
	}

	return parsed, nil
}

const redactedSecretMask = "••••••••"

func (s *OrganizationBrandingService) redactConfigSecrets(cfg *domain.NotificationConfig) *domain.NotificationConfig {
	if cfg == nil {
		return nil
	}
	cp := *cfg
	if cp.Password != nil && *cp.Password != "" {
		mask := redactedSecretMask
		cp.Password = &mask
	}
	if cp.APIKey != nil && *cp.APIKey != "" {
		mask := redactedSecretMask
		cp.APIKey = &mask
	}
	return &cp
}

func (s *OrganizationBrandingService) GetBranding(ctx context.Context, principal *middleware.AuthenticatedPrincipal) (*domain.BrandingConfig, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	return s.brandingRepo.GetBranding(ctx, orgUUID)
}

func (s *OrganizationBrandingService) UpdateBranding(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	payload *domain.UpdateBrandingPayload,
) (*domain.BrandingConfig, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	updated, errUp := s.brandingRepo.UpdateBranding(ctx, orgUUID, payload, actorUUID)
	if errUp != nil {
		return nil, errUp
	}

	if s.auditRepo != nil {
		action := "BRANDING_UPDATED"
		resType := "organization.organizations"
		resID := orgUUID.String()
		eventCat := "ORGANIZATION_BRANDING"
		severity := "HIGH"
		status := "SUCCESS"
		orgIDStr := orgUUID.String()

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    s.isPlatformAdmin(principal),
			TenantID:      &orgIDStr,
			ActorID:       &principal.UserID,
			ActorName:     &principal.Identity.FullName,
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

func (s *OrganizationBrandingService) SaveNotificationConfig(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	payload *domain.SaveNotificationConfigPayload,
) (*domain.NotificationConfig, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	if !domain.IsValidNotificationChannel(payload.Channel) {
		return nil, domain.ErrInvalidNotificationChannel
	}
	if !domain.IsValidNotificationProvider(payload.Provider) {
		return nil, domain.ErrInvalidNotificationProvider
	}

	var (
		encPassword *string
		encAPIKey   *string
	)

	if payload.Password != nil && *payload.Password != "" && *payload.Password != redactedSecretMask {
		enc, errEnc := crypto.EncryptAEAD(*payload.Password, s.masterKey)
		if errEnc != nil {
			return nil, fmt.Errorf("failed to encrypt notification password: %w", errEnc)
		}
		encPassword = &enc
	}

	if payload.APIKey != nil && *payload.APIKey != "" && *payload.APIKey != redactedSecretMask {
		enc, errEnc := crypto.EncryptAEAD(*payload.APIKey, s.masterKey)
		if errEnc != nil {
			return nil, fmt.Errorf("failed to encrypt notification API key: %w", errEnc)
		}
		encAPIKey = &enc
	}

	isActiveVal := true
	if payload.IsActive != nil {
		isActiveVal = *payload.IsActive
	}

	configEntity := &domain.NotificationConfig{
		OrganizationID: orgUUID,
		Channel:        payload.Channel,
		Provider:       payload.Provider,
		SenderEmail:    payload.SenderEmail,
		SenderName:     payload.SenderName,
		Host:           payload.Host,
		Port:           payload.Port,
		Username:       payload.Username,
		Password:       encPassword,
		APIKey:         encAPIKey,
		ConfigMetadata: payload.ConfigMetadata,
		IsActive:       isActiveVal,
	}

	saved, errSave := s.brandingRepo.SaveNotificationConfig(ctx, configEntity, actorUUID)
	if errSave != nil {
		return nil, errSave
	}

	if s.auditRepo != nil {
		action := "NOTIFICATION_CONFIG_SAVED"
		resType := "organization.notification_configs"
		resID := saved.ID.String()
		eventCat := "ORGANIZATION_NOTIFICATION_GOVERNANCE"
		severity := "HIGH"
		status := "SUCCESS"
		orgIDStr := orgUUID.String()

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    s.isPlatformAdmin(principal),
			TenantID:      &orgIDStr,
			ActorID:       &principal.UserID,
			ActorName:     &principal.Identity.FullName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	return s.redactConfigSecrets(saved), nil
}

func (s *OrganizationBrandingService) ListNotificationConfigs(ctx context.Context, principal *middleware.AuthenticatedPrincipal) ([]domain.NotificationConfig, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	configs, err := s.brandingRepo.ListNotificationConfigs(ctx, orgUUID)
	if err != nil {
		return nil, err
	}

	redactedList := make([]domain.NotificationConfig, len(configs))
	for i := range configs {
		redactedList[i] = *s.redactConfigSecrets(&configs[i])
	}

	return redactedList, nil
}

func (s *OrganizationBrandingService) SaveNotificationTemplate(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	templateKey string,
	payload *domain.SaveNotificationTemplatePayload,
) (*domain.NotificationTemplate, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	if !domain.IsValidNotificationChannel(payload.Channel) {
		return nil, domain.ErrInvalidNotificationChannel
	}

	templateEntity := &domain.NotificationTemplate{
		OrganizationID:   orgUUID,
		TemplateKey:      templateKey,
		Channel:          payload.Channel,
		Subject:          payload.Subject,
		BodyTemplate:     payload.BodyTemplate,
		AllowedVariables: payload.AllowedVariables,
		IsActive:         true,
	}

	saved, errSave := s.brandingRepo.SaveNotificationTemplate(ctx, templateEntity, actorUUID)
	if errSave != nil {
		return nil, errSave
	}

	if s.auditRepo != nil {
		action := "NOTIFICATION_TEMPLATE_SAVED"
		resType := "organization.notification_templates"
		resID := saved.ID.String()
		eventCat := "ORGANIZATION_NOTIFICATION_GOVERNANCE"
		severity := "MEDIUM"
		status := "SUCCESS"
		orgIDStr := orgUUID.String()

		_, _ = s.auditRepo.Create(ctx, &auditDomain.CreateAuditLogPayload{
			IsPlatform:    s.isPlatformAdmin(principal),
			TenantID:      &orgIDStr,
			ActorID:       &principal.UserID,
			ActorName:     &principal.Identity.FullName,
			ActorRole:     &principal.Role,
			Action:        action,
			ResourceType:  &resType,
			ResourceID:    &resID,
			EventCategory: &eventCat,
			Severity:      severity,
			Status:        status,
		})
	}

	return saved, nil
}

func (s *OrganizationBrandingService) ListNotificationTemplates(ctx context.Context, principal *middleware.AuthenticatedPrincipal) ([]domain.NotificationTemplate, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	return s.brandingRepo.ListNotificationTemplates(ctx, orgUUID)
}

func (s *OrganizationBrandingService) ListUserNotifications(ctx context.Context, principal *middleware.AuthenticatedPrincipal, limit int) ([]domain.UserNotification, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	userUUID, errUser := uuid.Parse(principal.UserID)
	if errUser != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errUser)
	}

	return s.brandingRepo.ListUserNotifications(ctx, orgUUID, userUUID, limit)
}

func (s *OrganizationBrandingService) MarkNotificationRead(ctx context.Context, principal *middleware.AuthenticatedPrincipal, notifID uuid.UUID) error {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return err
	}

	userUUID, errUser := uuid.Parse(principal.UserID)
	if errUser != nil {
		return fmt.Errorf("invalid principal user ID: %w", errUser)
	}

	return s.brandingRepo.MarkNotificationRead(ctx, orgUUID, userUUID, notifID)
}

func (s *OrganizationBrandingService) MarkAllNotificationsRead(ctx context.Context, principal *middleware.AuthenticatedPrincipal) error {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return err
	}

	userUUID, errUser := uuid.Parse(principal.UserID)
	if errUser != nil {
		return fmt.Errorf("invalid principal user ID: %w", errUser)
	}

	return s.brandingRepo.MarkAllNotificationsRead(ctx, orgUUID, userUUID)
}

func (s *OrganizationBrandingService) ListNotificationDeliveries(ctx context.Context, principal *middleware.AuthenticatedPrincipal, limit int) ([]domain.NotificationDelivery, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	return s.brandingRepo.ListNotificationDeliveries(ctx, orgUUID, limit)
}
