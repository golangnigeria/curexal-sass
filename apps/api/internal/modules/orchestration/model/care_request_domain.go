package model

import "time"

// CareRequest represents a demand for clinical or diagnostic healthcare services
type CareRequest struct {
	ID                  string                 `json:"id" db:"id"`
	TenantID            string                 `json:"tenantId" db:"tenant_id"`
	PatientID           string                 `json:"patientId" db:"patient_id"`
	RequestNumber       string                 `json:"requestNumber" db:"request_number"`
	ServiceType         string                 `json:"serviceType" db:"service_type"`     // GENERAL_CONSULTATION, SPECIALIST, LAB_TEST, REFILL, TELEHEALTH
	PreferredMode       string                 `json:"preferredMode" db:"preferred_mode"` // IN_PERSON, VIDEO, AUDIO, ASYNC_CHAT
	Urgency             string                 `json:"urgency" db:"urgency"`              // ROUTINE, URGENT, EMERGENCY
	Status              string                 `json:"status" db:"status"`                // SUBMITTED, TRIAGED, ASSIGNED_AGENT, MATCHED, IN_PROGRESS, COMPLETED, CANCELLED
	ChiefComplaint      *string                `json:"chiefComplaint,omitempty" db:"chief_complaint"`
	SymptomsJSON        []string               `json:"symptoms" db:"symptoms_json"`
	PreferredTimeWindow map[string]interface{} `json:"preferredTimeWindow,omitempty" db:"preferred_time_window"`
	AssignedCareAgentID *string                `json:"assignedCareAgentId,omitempty" db:"assigned_care_agent_id"`
	MatchedProviderID   *string                `json:"matchedProviderId,omitempty" db:"matched_provider_id"`
	CreatedAt           time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time              `json:"updatedAt" db:"updated_at"`

	// Enriched Patient Demographics (when queried)
	PatientName *string `json:"patientName,omitempty"`
	MRN         *string `json:"mrn,omitempty"`
}

// CareJourneyMilestone represents a step along the Curexal Care Journey
type CareJourneyMilestone struct {
	ID             string     `json:"id" db:"id"`
	TenantID       string     `json:"tenantId" db:"tenant_id"`
	PatientID      string     `json:"patientId" db:"patient_id"`
	EncounterID    *string    `json:"encounterId,omitempty" db:"encounter_id"`
	StageCode      string     `json:"stageCode" db:"stage_code"` // INTAKE, TRIAGE, CONSULTATION, LAB_WORKLIST, RADIOLOGY_STUDY, PHARMACY_DISPENSE, SETTLEMENT, FOLLOW_UP
	Title          string     `json:"title" db:"title"`
	Description    *string    `json:"description,omitempty" db:"description"`
	Status         string     `json:"status" db:"status"` // PENDING, IN_PROGRESS, COMPLETED, SKIPPED, BLOCKED
	BlockingReason *string    `json:"blockingReason,omitempty" db:"blocking_reason"`
	CompletedAt    *time.Time `json:"completedAt,omitempty" db:"completed_at"`
	CreatedAt      time.Time  `json:"createdAt" db:"created_at"`
}

// CreateCareRequestPayload payload for creating a new care request
type CreateCareRequestPayload struct {
	PatientID           string                 `json:"patientId"`
	ServiceType         string                 `json:"serviceType" validate:"required"` // GENERAL_CONSULTATION, SPECIALIST, LAB_TEST, REFILL, TELEHEALTH
	PreferredMode       string                 `json:"preferredMode" validate:"required,oneof=IN_PERSON VIDEO AUDIO ASYNC_CHAT"`
	Urgency             string                 `json:"urgency" validate:"omitempty,oneof=ROUTINE URGENT EMERGENCY"`
	ChiefComplaint      string                 `json:"chiefComplaint" validate:"required"`
	Symptoms            []string               `json:"symptoms"`
	PreferredTimeWindow map[string]interface{} `json:"preferredTimeWindow"`
}

// CareRequestFilter filter parameters for querying care requests
type CareRequestFilter struct {
	Status    string `query:"status"`
	Urgency   string `query:"urgency"`
	PatientID string `query:"patientId"`
	Limit     int    `query:"limit"`
	Offset    int    `query:"offset"`
}

// CareRequestListResponse paginated list of care requests
type CareRequestListResponse struct {
	Items  []CareRequest `json:"items"`
	Total  int           `json:"total"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}
