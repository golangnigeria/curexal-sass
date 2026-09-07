package model

import "time"

// TriageAssessment represents clinical vitals intake & urgency categorization
type TriageAssessment struct {
	ID              string     `json:"id" db:"id"`
	CareRequestID   string     `json:"careRequestId" db:"care_request_id"`
	PatientID       string     `json:"patientId" db:"patient_id"`
	AssessorID      *string    `json:"assessorId,omitempty" db:"assessor_id"`
	AcuityLevel     string     `json:"acuityLevel" db:"acuity_level"` // RED (Immediate), YELLOW (Urgent), GREEN (Standard)
	SystolicBP      *int       `json:"systolicBp,omitempty" db:"systolic_bp"`
	DiastolicBP     *int       `json:"diastolicBp,omitempty" db:"diastolic_bp"`
	PulseRate       *int       `json:"pulseRate,omitempty" db:"pulse_rate"`
	Temperature     *float64   `json:"temperature,omitempty" db:"temperature"`
	SpO2            *int       `json:"spo2,omitempty" db:"spo2"`
	RespiratoryRate *int       `json:"respiratoryRate,omitempty" db:"respiratory_rate"`
	PainScore       *int       `json:"painScore,omitempty" db:"pain_score"`
	TriageNotes     *string    `json:"triageNotes,omitempty" db:"triage_notes"`
	CreatedAt       time.Time  `json:"createdAt" db:"created_at"`
}

// SubmitTriagePayload payload for submitting triage vitals
type SubmitTriagePayload struct {
	SystolicBP      *int     `json:"systolicBp"`
	DiastolicBP     *int     `json:"diastolicBp"`
	PulseRate       *int     `json:"pulseRate"`
	Temperature     *float64 `json:"temperature"`
	SpO2            *int     `json:"spo2"`
	RespiratoryRate *int     `json:"respiratoryRate"`
	PainScore       *int     `json:"painScore"`
	TriageNotes     *string  `json:"triageNotes"`
	AcuityOverride  *string  `json:"acuityOverride"` // Optional nurse override: RED, YELLOW, GREEN
}

// ProviderProfile represents doctor availability, clinical specialty and active capacity
type ProviderProfile struct {
	ID                    string    `json:"id" db:"id"`
	UserID                string    `json:"userId" db:"user_id"`
	TenantID              string    `json:"tenantId" db:"tenant_id"`
	SpecialtyCode         string    `json:"specialtyCode" db:"specialty_code"`
	SubSpecialties        []string  `json:"subSpecialties" db:"sub_specialties"`
	TelehealthEnabled     bool      `json:"telehealthEnabled" db:"telehealth_enabled"`
	InPersonEnabled       bool      `json:"inPersonEnabled" db:"in_person_enabled"`
	MaxActiveQueue        int       `json:"maxActiveQueue" db:"max_active_queue"`
	CurrentActiveQueue    int       `json:"currentActiveQueue" db:"current_active_queue"`
	Status                string    `json:"status" db:"status"` // ON_DUTY, ON_BREAK, OFF_DUTY, BUSY
	ConsultationLanguages []string  `json:"consultationLanguages" db:"consultation_languages"`
	Rating                float64   `json:"rating" db:"rating"`
	CreatedAt             time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt             time.Time `json:"updatedAt" db:"updated_at"`

	// Enriched provider name
	ProviderName *string `json:"providerName,omitempty"`
}

// MatchedProviderCandidate scored candidate for care request assignment
type MatchedProviderCandidate struct {
	ProviderID         string   `json:"providerId"`
	UserID             string   `json:"userId"`
	ProviderName       string   `json:"providerName"`
	SpecialtyCode      string   `json:"specialtyCode"`
	Status             string   `json:"status"`
	CurrentActiveQueue int      `json:"currentActiveQueue"`
	MaxActiveQueue     int      `json:"maxActiveQueue"`
	MatchScore         int      `json:"matchScore"` // 0 - 100
	MatchingReasons    []string `json:"matchingReasons"`
}

// ProviderMatchingResponse returned by matching algorithm
type ProviderMatchingResponse struct {
	CareRequestID string                     `json:"careRequestId"`
	RequiredMode  string                     `json:"requiredMode"`
	Specialty     string                     `json:"specialty"`
	Candidates    []MatchedProviderCandidate `json:"candidates"`
}

// AssignProviderPayload payload to reserve slot and assign doctor
type AssignProviderPayload struct {
	ProviderID string `json:"providerId" validate:"required"`
}
