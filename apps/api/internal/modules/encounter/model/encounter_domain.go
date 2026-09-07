package model

import "time"

// Encounter represents a unified clinical consultation session
type Encounter struct {
	ID                   string     `json:"id" db:"id"`
	TenantID             string     `json:"tenantId" db:"tenant_id"`
	CareRequestID        *string    `json:"careRequestId,omitempty" db:"care_request_id"`
	PatientID            string     `json:"patientId" db:"patient_id"`
	ProviderID           string     `json:"providerId" db:"provider_id"`
	EncounterType        string     `json:"encounterType" db:"encounter_type"` // OUTPATIENT, EMERGENCY, TELEHEALTH, INPATIENT_ROUND
	Mode                 string     `json:"mode" db:"mode"`                    // IN_PERSON, VIDEO, AUDIO, ASYNC_CHAT
	Status               string     `json:"status" db:"status"`                // WAITING, IN_PROGRESS, ON_HOLD, COMPLETED, CANCELLED
	ChiefComplaint       *string    `json:"chiefComplaint,omitempty" db:"chief_complaint"`
	Subjective           *string    `json:"subjective,omitempty" db:"subjective"`
	Objective            *string    `json:"objective,omitempty" db:"objective"`
	Assessment           *string    `json:"assessment,omitempty" db:"assessment"`
	Plan                 *string    `json:"plan,omitempty" db:"plan"`
	PrimaryDiagnosisCode *string    `json:"primaryDiagnosisCode,omitempty" db:"primary_diagnosis_code"`
	PrimaryDiagnosisName *string    `json:"primaryDiagnosisName,omitempty" db:"primary_diagnosis_name"`
	SecondaryDiagnoses   []string   `json:"secondaryDiagnoses" db:"secondary_diagnoses"`
	StartedAt            time.Time  `json:"startedAt" db:"started_at"`
	CompletedAt          *time.Time `json:"completedAt,omitempty" db:"completed_at"`
	CreatedAt            time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt            time.Time  `json:"updatedAt" db:"updated_at"`

	// Enriched fields for consultation workspace
	PatientName *string `json:"patientName,omitempty"`
	MRN         *string `json:"mrn,omitempty"`
	Gender      *string `json:"gender,omitempty"`
	AgeYears    *int    `json:"ageYears,omitempty"`
}

// StartEncounterPayload initiates a new clinical encounter
type StartEncounterPayload struct {
	CareRequestID  *string `json:"careRequestId"`
	PatientID      string  `json:"patientId" validate:"required"`
	ProviderID     string  `json:"providerId" validate:"required"`
	EncounterType  string  `json:"encounterType" validate:"required,oneof=OUTPATIENT EMERGENCY TELEHEALTH INPATIENT_ROUND"`
	Mode           string  `json:"mode" validate:"required,oneof=IN_PERSON VIDEO AUDIO ASYNC_CHAT"`
	ChiefComplaint *string `json:"chiefComplaint"`
}

// UpdateSOAPNotesPayload updates clinical notes and primary diagnosis
type UpdateSOAPNotesPayload struct {
	ChiefComplaint       *string  `json:"chiefComplaint"`
	Subjective           *string  `json:"subjective"`
	Objective            *string  `json:"objective"`
	Assessment           *string  `json:"assessment"`
	Plan                 *string  `json:"plan"`
	PrimaryDiagnosisCode *string  `json:"primaryDiagnosisCode"`
	PrimaryDiagnosisName *string  `json:"primaryDiagnosisName"`
	SecondaryDiagnoses   []string `json:"secondaryDiagnoses"`
}

// LabOrderRequest item for diagnostic test ordering
type LabOrderRequest struct {
	TestCode string `json:"testCode"`
	TestName string `json:"testName"`
	Urgency  string `json:"urgency"` // ROUTINE, STAT, URGENT
	Notes    string `json:"notes"`
}

// RadiologyOrderRequest item for imaging study ordering
type RadiologyOrderRequest struct {
	Modality  string `json:"modality"` // XRAY, CT, MRI, ULTRASOUND
	StudyName string `json:"studyName"`
	Urgency   string `json:"urgency"`
	Notes     string `json:"notes"`
}

// PrescriptionItem item for medication ordering
type PrescriptionItem struct {
	MedicationName string `json:"medicationName"`
	Dosage         string `json:"dosage"`
	Frequency      string `json:"frequency"`
	DurationDays   int    `json:"durationDays"`
	Instructions   string `json:"instructions"`
}

// DispatchOrdersPayload clinical order dispatch request
type DispatchOrdersPayload struct {
	LabOrders        []LabOrderRequest       `json:"labOrders"`
	RadiologyOrders  []RadiologyOrderRequest `json:"radiologyOrders"`
	PrescriptionList []PrescriptionItem      `json:"prescriptionList"`
}

// TelehealthSession represents virtual consultation connection state
type TelehealthSession struct {
	ID               string     `json:"id" db:"id"`
	EncounterID      string     `json:"encounterId" db:"encounter_id"`
	RoomSID          string     `json:"roomSid" db:"room_sid"`
	ChannelType      string     `json:"channelType" db:"channel_type"`
	Status           string     `json:"status" db:"status"`
	PatientJoinedAt  *time.Time `json:"patientJoinedAt,omitempty" db:"patient_joined_at"`
	ProviderJoinedAt *time.Time `json:"providerJoinedAt,omitempty" db:"provider_joined_at"`
	EndedAt          *time.Time `json:"endedAt,omitempty" db:"ended_at"`
	CreatedAt        time.Time  `json:"createdAt" db:"created_at"`
}
