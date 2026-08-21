package application

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	auditDomain "github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
)

type OrganizationIntegrationService struct {
	integrationRepo domain.OrganizationIntegrationRepository
	orgRepo         domain.OrganizationRepository
	auditRepo       auditDomain.AuditRepository
	masterKey       string
}

func NewOrganizationIntegrationService(
	integrationRepo domain.OrganizationIntegrationRepository,
	orgRepo domain.OrganizationRepository,
	auditRepo auditDomain.AuditRepository,
	masterKey string,
) *OrganizationIntegrationService {
	if masterKey == "" {
		masterKey = "CUREXAL_DEFAULT_PLATFORM_AES_MASTER_KEY_2026"
	}
	return &OrganizationIntegrationService{
		integrationRepo: integrationRepo,
		orgRepo:         orgRepo,
		auditRepo:       auditRepo,
		masterKey:       masterKey,
	}
}

func (s *OrganizationIntegrationService) isPlatformAdmin(principal *middleware.AuthenticatedPrincipal) bool {
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

func (s *OrganizationIntegrationService) resolveActiveOrgUUID(principal *middleware.AuthenticatedPrincipal) (uuid.UUID, error) {
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

func generateRandomToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate secure random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

func computeSHA256(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// ComputeHMACSignature generates standard HMAC-SHA256 signature for outbound webhooks.
func ComputeHMACSignature(secret, timestamp, payload string) string {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(timestamp + "." + payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func (s *OrganizationIntegrationService) CreateAPIKey(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	payload *domain.CreateAPIKeyPayload,
) (*domain.APIKeyCreateResult, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	// Validate scopes
	if len(payload.Scopes) > 0 {
		var scopeList []string
		if errUnmarshal := json.Unmarshal(payload.Scopes, &scopeList); errUnmarshal == nil {
			for _, sc := range scopeList {
				if !domain.IsValidScope(sc) {
					return nil, fmt.Errorf("invalid API key scope '%s'", sc)
				}
			}
		}
	}

	// Validate IP Whitelist
	if len(payload.IPWhitelist) > 0 {
		var ipList []string
		if errUnmarshal := json.Unmarshal(payload.IPWhitelist, &ipList); errUnmarshal == nil {
			if errVal := domain.ValidateIPWhitelist(ipList); errVal != nil {
				return nil, errVal
			}
		}
	}

	// Generate secure raw key e.g. "cx_live_<32 hex chars>"
	randomHex, errToken := generateRandomToken(16)
	if errToken != nil {
		return nil, errToken
	}
	rawKey := fmt.Sprintf("cx_live_%s", randomHex)
	keyPrefix := rawKey[:12] // e.g. "cx_live_abcd"
	keyHash := computeSHA256(rawKey)

	rpm := 60
	if payload.RateLimitRPM != nil && *payload.RateLimitRPM > 0 {
		rpm = *payload.RateLimitRPM
	}

	apiKeyEntity := &domain.APIKey{
		OrganizationID: orgUUID,
		Name:           strings.TrimSpace(payload.Name),
		KeyPrefix:      keyPrefix,
		Scopes:         payload.Scopes,
		IPWhitelist:    payload.IPWhitelist,
		RateLimitRPM:   rpm,
		ExpiresAt:      payload.ExpiresAt,
	}

	created, errCreate := s.integrationRepo.CreateAPIKey(ctx, apiKeyEntity, keyHash, actorUUID)
	if errCreate != nil {
		return nil, errCreate
	}

	if s.auditRepo != nil {
		action := "API_KEY_CREATED"
		resType := "organization.api_keys"
		resID := created.ID.String()
		eventCat := "ORGANIZATION_SECURITY_GOVERNANCE"
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

	return &domain.APIKeyCreateResult{
		APIKey: *created,
		RawKey: rawKey,
	}, nil
}

func (s *OrganizationIntegrationService) ListAPIKeys(ctx context.Context, principal *middleware.AuthenticatedPrincipal) ([]domain.APIKey, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	return s.integrationRepo.ListAPIKeys(ctx, orgUUID)
}

func (s *OrganizationIntegrationService) RevokeAPIKey(ctx context.Context, principal *middleware.AuthenticatedPrincipal, keyID uuid.UUID) error {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	errRev := s.integrationRepo.RevokeAPIKey(ctx, orgUUID, keyID, actorUUID)
	if errRev != nil {
		return errRev
	}

	if s.auditRepo != nil {
		action := "API_KEY_REVOKED"
		resType := "organization.api_keys"
		resID := keyID.String()
		eventCat := "ORGANIZATION_SECURITY_GOVERNANCE"
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

	return nil
}

func (s *OrganizationIntegrationService) CreateWebhookSubscription(
	ctx context.Context,
	principal *middleware.AuthenticatedPrincipal,
	payload *domain.CreateWebhookSubscriptionPayload,
) (*domain.WebhookSubscription, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return nil, fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	// Validate SSRF Protection URL boundary
	if errSSRF := domain.ValidateSSRFSafeURL(payload.TargetURL); errSSRF != nil {
		return nil, errSSRF
	}

	secretVal := ""
	if payload.SigningSecret != nil && *payload.SigningSecret != "" {
		secretVal = *payload.SigningSecret
	} else {
		genSecret, errGen := generateRandomToken(16)
		if errGen != nil {
			return nil, errGen
		}
		secretVal = fmt.Sprintf("whsec_%s", genSecret)
	}

	encSecret, errEnc := crypto.EncryptAEAD(secretVal, s.masterKey)
	if errEnc != nil {
		return nil, fmt.Errorf("failed to encrypt webhook signing secret: %w", errEnc)
	}

	subEntity := &domain.WebhookSubscription{
		OrganizationID: orgUUID,
		Name:           strings.TrimSpace(payload.Name),
		TargetURL:      strings.TrimSpace(payload.TargetURL),
		EventTypes:     payload.EventTypes,
		SigningSecret:  &encSecret,
		IsActive:       true,
	}

	created, errCreate := s.integrationRepo.CreateWebhookSubscription(ctx, subEntity, actorUUID)
	if errCreate != nil {
		return nil, errCreate
	}

	if s.auditRepo != nil {
		action := "WEBHOOK_SUBSCRIPTION_CREATED"
		resType := "organization.webhook_subscriptions"
		resID := created.ID.String()
		eventCat := "ORGANIZATION_SECURITY_GOVERNANCE"
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

	redactedMask := redactedSecretMask
	created.SigningSecret = &redactedMask
	return created, nil
}

func (s *OrganizationIntegrationService) ListWebhookSubscriptions(ctx context.Context, principal *middleware.AuthenticatedPrincipal) ([]domain.WebhookSubscription, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	subs, err := s.integrationRepo.ListWebhookSubscriptions(ctx, orgUUID)
	if err != nil {
		return nil, err
	}

	for i := range subs {
		redactedMask := redactedSecretMask
		subs[i].SigningSecret = &redactedMask
	}

	return subs, nil
}

func (s *OrganizationIntegrationService) DeleteWebhookSubscription(ctx context.Context, principal *middleware.AuthenticatedPrincipal, subID uuid.UUID) error {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return err
	}

	actorUUID, errParse := uuid.Parse(principal.UserID)
	if errParse != nil {
		return fmt.Errorf("invalid principal user ID: %w", errParse)
	}

	errDel := s.integrationRepo.DeleteWebhookSubscription(ctx, orgUUID, subID, actorUUID)
	if errDel != nil {
		return errDel
	}

	if s.auditRepo != nil {
		action := "WEBHOOK_SUBSCRIPTION_DELETED"
		resType := "organization.webhook_subscriptions"
		resID := subID.String()
		eventCat := "ORGANIZATION_SECURITY_GOVERNANCE"
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

	return nil
}

func (s *OrganizationIntegrationService) ListWebhookDeliveries(ctx context.Context, principal *middleware.AuthenticatedPrincipal, limit int) ([]domain.WebhookDelivery, error) {
	orgUUID, err := s.resolveActiveOrgUUID(principal)
	if err != nil {
		return nil, err
	}

	return s.integrationRepo.ListWebhookDeliveries(ctx, orgUUID, limit)
}
