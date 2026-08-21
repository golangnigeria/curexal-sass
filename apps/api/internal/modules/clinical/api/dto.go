package api

import (
	"github.com/golangnigeria/curexal/internal/modules/clinical/domain"
	"github.com/google/uuid"
)

type CreateCatalogItemPayload = domain.CreateCatalogItemPayload
type UpdatePricingPayload = domain.UpdatePricingPayload

type VisitRegistrationPayload struct {
	TenantID            uuid.UUID                  `json:"tenant_id"`
	PatientDemographics PatientDemographicsPayload `json:"patient_demographics"`
	VisitDetails        VisitDetailsPayload        `json:"visit_details"`
	Investigations      []InvestigationPayload     `json:"investigations"`
	PaymentAllocations  []PaymentAllocationPayload `json:"payment_allocations"`
}

type PatientDemographicsPayload struct {
	IsExisting bool      `json:"is_existing"`
	PatientID  uuid.UUID `json:"patient_id"`
	Title      string    `json:"title"`
	FirstName  string    `json:"first_name"`
	MiddleName string    `json:"middle_name"`
	LastName   string    `json:"last_name"`
	Gender     string    `json:"gender"`
	DOB        string    `json:"dob"`
	Phone      string    `json:"phone"`
	Email      string    `json:"email"`
	Address    string    `json:"address"`
}

type VisitDetailsPayload struct {
	CollectionCenterID uuid.UUID        `json:"collection_center_id"`
	CollectorUserID    uuid.UUID        `json:"collector_user_id"`
	Referrals          ReferralsPayload `json:"referrals"`
}

type ReferralsPayload struct {
	DoctorPartnerID   *uuid.UUID `json:"doctor_partner_id"`
	FacilityPartnerID *uuid.UUID `json:"facility_partner_id"`
	CampaignID        *uuid.UUID `json:"campaign_id"`
}

type InvestigationPayload struct {
	InvestigationID string `json:"investigation_id"`
	Quantity        int    `json:"quantity"`
}

type PaymentAllocationPayload struct {
	GuarantorType        string     `json:"guarantor_type"`
	PayerPartnerID       *uuid.UUID `json:"payer_partner_id"`
	AllocatedAmount      float64    `json:"allocated_amount"`
	BillingPolicy        string     `json:"billing_policy"`
	PreAuthorizationCode string     `json:"pre_authorization_code"`
}
