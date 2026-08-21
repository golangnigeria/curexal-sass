package domain

import (
	"context"

	"github.com/google/uuid"
)

type OrganizationRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Organization, error)
	Create(ctx context.Context, name, slug, plan string) (*Organization, error)
	Update(ctx context.Context, id uuid.UUID, name, slug, plan, customDomain *string, settings map[string]any) (*Organization, error)
	UpdateProfile(ctx context.Context, orgID uuid.UUID, payload *UpdateOrganizationProfilePayload, actorID uuid.UUID) (*Organization, error)
	UpdateSetupState(ctx context.Context, orgID uuid.UUID, newState SetupState, step int, actorID uuid.UUID) (*Organization, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetSettings(ctx context.Context, orgID uuid.UUID) (*OrganizationSettings, error)
	UpdateSettings(ctx context.Context, orgID uuid.UUID, logoURL, themeBranding, customDomain, supportEmail, supportPhone, cacNumber, tinNumber, taxNumber, businessType, address, timezone, currency, language *string) (*OrganizationSettings, error)
	List(ctx context.Context, userID string, isPlatformAdmin bool) ([]Organization, error)
}

type FacilityBranchRepository interface {
	ListBranches(ctx context.Context, orgID uuid.UUID) ([]FacilityBranch, error)
	GetBranchByID(ctx context.Context, orgID, branchID uuid.UUID) (*FacilityBranch, error)
	GetBranchByCode(ctx context.Context, orgID uuid.UUID, code string) (*FacilityBranch, error)
	CreateBranch(ctx context.Context, branch *FacilityBranch, actorID uuid.UUID) (*FacilityBranch, error)
	UpdateBranch(ctx context.Context, branch *FacilityBranch, actorID uuid.UUID) (*FacilityBranch, error)
	DeactivateBranch(ctx context.Context, orgID, branchID uuid.UUID, actorID uuid.UUID) error
	CountActiveBranches(ctx context.Context, orgID uuid.UUID) (int, error)
	HasActiveHeadquarters(ctx context.Context, orgID uuid.UUID) (bool, error)
	CheckFacilityTypeActive(ctx context.Context, facilityTypeID uuid.UUID) (bool, error)
}

type StaffMembershipRepository interface {
	ListMembers(ctx context.Context, orgID uuid.UUID) ([]StaffMemberDTO, error)
	GetMemberByID(ctx context.Context, orgID, membershipID uuid.UUID) (*StaffMemberDTO, error)
	CountActiveMembers(ctx context.Context, orgID uuid.UUID) (int, error)

	CreateInvitation(ctx context.Context, invite *StaffInvitation) (*StaffInvitation, error)
	ListInvitations(ctx context.Context, orgID uuid.UUID) ([]StaffInvitation, error)
	GetInvitationByHash(ctx context.Context, hash string) (*StaffInvitation, error)
	RevokeInvitation(ctx context.Context, orgID, inviteID uuid.UUID) error
	AcceptInvitation(ctx context.Context, inviteID uuid.UUID) error
	CheckPendingInviteExists(ctx context.Context, orgID uuid.UUID, email string) (bool, error)

	AssignBranch(ctx context.Context, membershipID, branchID uuid.UUID, actorID uuid.UUID) (*MembershipBranch, error)
	RemoveBranchAssignment(ctx context.Context, membershipID, branchID uuid.UUID) error
	AssignDepartment(ctx context.Context, membershipID, branchID uuid.UUID, deptCode string, actorID uuid.UUID) (*DepartmentalMembership, error)
	RemoveDepartmentAssignment(ctx context.Context, membershipID, branchID uuid.UUID, deptCode string) error
	UpdateMemberRole(ctx context.Context, orgID, membershipID uuid.UUID, role, roleTitle string, actorID uuid.UUID) (*StaffMemberDTO, error)
}

type OrganizationCatalogRepository interface {
	ListCatalogItems(ctx context.Context, orgID uuid.UUID, domainType string) ([]OrganizationCatalogItem, error)
	GetCatalogItemByID(ctx context.Context, orgID, itemID uuid.UUID) (*OrganizationCatalogItem, error)
	CreateCatalogItem(ctx context.Context, item *OrganizationCatalogItem, actorID uuid.UUID) (*OrganizationCatalogItem, error)
	UpdateCatalogItem(ctx context.Context, item *OrganizationCatalogItem, actorID uuid.UUID) (*OrganizationCatalogItem, error)

	SetBranchPriceOverride(ctx context.Context, override *BranchPriceOverride, actorID uuid.UUID) (*BranchPriceOverride, error)
	GetBranchPriceOverride(ctx context.Context, orgID, branchID, itemID uuid.UUID) (*BranchPriceOverride, error)
	ListBranchPriceOverrides(ctx context.Context, orgID, branchID uuid.UUID) ([]BranchPriceOverride, error)

	CreateInsuranceProvider(ctx context.Context, provider *InsuranceProvider, actorID uuid.UUID) (*InsuranceProvider, error)
	ListInsuranceProviders(ctx context.Context, orgID uuid.UUID) ([]InsuranceProvider, error)
}

type OrganizationBrandingRepository interface {
	GetBranding(ctx context.Context, orgID uuid.UUID) (*BrandingConfig, error)
	UpdateBranding(ctx context.Context, orgID uuid.UUID, payload *UpdateBrandingPayload, actorID uuid.UUID) (*BrandingConfig, error)

	SaveNotificationConfig(ctx context.Context, config *NotificationConfig, actorID uuid.UUID) (*NotificationConfig, error)
	ListNotificationConfigs(ctx context.Context, orgID uuid.UUID) ([]NotificationConfig, error)
	GetNotificationConfig(ctx context.Context, orgID uuid.UUID, channel, provider string) (*NotificationConfig, error)

	SaveNotificationTemplate(ctx context.Context, template *NotificationTemplate, actorID uuid.UUID) (*NotificationTemplate, error)
	ListNotificationTemplates(ctx context.Context, orgID uuid.UUID) ([]NotificationTemplate, error)
	GetNotificationTemplate(ctx context.Context, orgID uuid.UUID, templateKey, channel string) (*NotificationTemplate, error)

	CreateUserNotification(ctx context.Context, notif *UserNotification) (*UserNotification, error)
	ListUserNotifications(ctx context.Context, orgID, userID uuid.UUID, limit int) ([]UserNotification, error)
	MarkNotificationRead(ctx context.Context, orgID, userID, notifID uuid.UUID) error
	MarkAllNotificationsRead(ctx context.Context, orgID, userID uuid.UUID) error

	CreateNotificationDelivery(ctx context.Context, delivery *NotificationDelivery) (*NotificationDelivery, error)
	ListNotificationDeliveries(ctx context.Context, orgID uuid.UUID, limit int) ([]NotificationDelivery, error)
}

type OrganizationIntegrationRepository interface {
	CreateAPIKey(ctx context.Context, apiKey *APIKey, keyHash string, actorID uuid.UUID) (*APIKey, error)
	GetAPIKeyByID(ctx context.Context, orgID, keyID uuid.UUID) (*APIKey, error)
	GetAPIKeyByHash(ctx context.Context, keyHash string) (*APIKey, error)
	ListAPIKeys(ctx context.Context, orgID uuid.UUID) ([]APIKey, error)
	RevokeAPIKey(ctx context.Context, orgID, keyID uuid.UUID, actorID uuid.UUID) error

	CreateWebhookSubscription(ctx context.Context, sub *WebhookSubscription, actorID uuid.UUID) (*WebhookSubscription, error)
	GetWebhookSubscriptionByID(ctx context.Context, orgID, subID uuid.UUID) (*WebhookSubscription, error)
	ListWebhookSubscriptions(ctx context.Context, orgID uuid.UUID) ([]WebhookSubscription, error)
	DeleteWebhookSubscription(ctx context.Context, orgID, subID uuid.UUID, actorID uuid.UUID) error

	CreateWebhookDelivery(ctx context.Context, delivery *WebhookDelivery) (*WebhookDelivery, error)
	ListWebhookDeliveries(ctx context.Context, orgID uuid.UUID, limit int) ([]WebhookDelivery, error)
}

type TenantRepository interface {
	CheckSlugExists(ctx context.Context, slug string) (bool, error)
	CreateTenant(ctx context.Context, userID string, name, slug, orgID, location, phone, address string, logoURL, currency *string, modules []string) (*Tenant, error)
	GetTenantByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	GetTenantBySlug(ctx context.Context, slug string) (*Tenant, error)
	UpdateTenant(ctx context.Context, id uuid.UUID, name, slug, logoURL, currency *string, settingsJSON *string) (*Tenant, error)
	DeleteTenant(ctx context.Context, id uuid.UUID) error
	ListTenants(ctx context.Context) ([]Tenant, error)
	CountActiveMembers(ctx context.Context, tenantID uuid.UUID) (int, error)
	ListBranchesByOrgID(ctx context.Context, orgID uuid.UUID) ([]Tenant, error)
}

type SubscriptionRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*Subscription, error)
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) (*Subscription, error)
	Create(ctx context.Context, sub *Subscription) (*Subscription, error)
	Update(ctx context.Context, id uuid.UUID, sub *Subscription) (*Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type DemoRepository interface {
	Create(ctx context.Context, labName, contactName, email string, phone, dailyVolume, notes *string) (*DemoRequest, error)
	List(ctx context.Context) ([]DemoRequest, error)
	Update(ctx context.Context, id uuid.UUID, status, notes *string) (*DemoRequest, error)
}
