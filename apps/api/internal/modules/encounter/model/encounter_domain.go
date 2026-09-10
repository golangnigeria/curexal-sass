package model

import "time"

// Encounter represents a unified clinical consultation session
type Encounter struct {
	ID                   string     `json:"id" db:"id"`
	TenantID             string     `json:"tenantId" db:"tenant_id"`
	PatientID            string     `json:"patientId" db:"patient_id"`
	AppointmentID        *string    `json:"appointmentId,omitempty" db:"appointment_id"`
	CareRequestID        *string    `json:"careRequestId,omitempty" db:"care_request_id"`
	ProviderID           string     `json:"providerId" db:"provider_id"`
	EncounterChannel     string     `json:"encounterChannel" db:"encounter_channel"` // in_person, video, telephone, secure_message
	Status               string     `json:"status" db:"status"`                     // open, in_progress, awaiting_documentation, signed, closed, amended
	ChiefComplaint       *string    `json:"chiefComplaint,omitempty" db:"chief_complaint"`
	Subjective           *string    `json:"subjective,omitempty" db:"subjective"`
	Objective            *string    `json:"objective,omitempty" db:"objective"`
	Assessment           *string    `json:"assessment,omitempty" db:"assessment"`
	Plan                 *string    `json:"plan,omitempty" db:"plan"`
	PrimaryDiagnosisCode *string    `json:"primaryDiagnosisCode,omitempty" db:"primary_diagnosis_code"`
	PrimaryDiagnosisName *string    `json:"primaryDiagnosisName,omitempty" db:"primary_diagnosis_name"`
	SecondaryDiagnoses   []string   `json:"secondaryDiagnoses" db:"secondary_diagnoses"`
	StartedAt            time.Time  `json:"startedAt" db:"started_at"`
	SignedAt             *time.Time `json:"signedAt,omitempty" db:"signed_at"`
	SignedBy             *string    `json:"signedBy,omitempty" db:"signed_by"`
	ClosedAt             *time.Time `json:"closedAt,omitempty" db:"closed_at"`
	CreatedAt            time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt            time.Time  `json:"updatedAt" db:"updated_at"`

	// Enriched fields for consultation workspace
	PatientName  *string             `json:"patientName,omitempty"`
	MRN          *string             `json:"mrn,omitempty"`
	Gender       *string             `json:"gender,omitempty"`
	AgeYears     *int                `json:"ageYears,omitempty"`
	ProviderName *string             `json:"providerName,omitempty"`
	Observations []Observation       `json:"observations,omitempty"`
	Diagnoses    []Diagnosis         `json:"diagnoses,omitempty"`
	Prescriptions []Prescription     `json:"prescriptions,omitempty"`
}

// Observation represents a vital sign or clinical measurement
type Observation struct {
	ID           string    `json:"id" db:"id"`
	EncounterID  string    `json:"encounterId" db:"encounter_id"`
	PatientID    string    `json:"patientId" db:"patient_id"`
	Code         string    `json:"code" db:"code"`
	ValueNumeric *float64  `json:"valueNumeric,omitempty" db:"value_numeric"`
	ValueText    *string   `json:"valueText,omitempty" db:"value_text"`
	Unit         *string   `json:"unit,omitempty" db:"unit"`
	Source       string    `json:"source" db:"source"` // staff_measured, patient_reported, device_reported, provider_observed
	ObservedAt   time.Time `json:"observedAt" db:"observed_at"`
	RecordedBy   *string   `json:"recordedBy,omitempty" db:"recorded_by"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

// Diagnosis represents an ICD-10 diagnostic condition linked to an encounter
type Diagnosis struct {
	ID                 string    `json:"id" db:"id"`
	EncounterID        string    `json:"encounterId" db:"encounter_id"`
	PatientID          string    `json:"patientId" db:"patient_id"`
	ICD10Code          string    `json:"icd10Code" db:"icd10_code"`
	ICD10Title         string    `json:"icd10Title" db:"icd10_title"`
	IsPrimary          bool      `json:"isPrimary" db:"is_primary"`
	ClinicalStatus     string    `json:"clinicalStatus" db:"clinical_status"` // ACTIVE, RECURRENCE, RESOLVED
	VerificationStatus string    `json:"verificationStatus" db:"verification_status"` // PROVISIONAL, DIFFERENTIAL, CONFIRMED
	Notes              *string   `json:"notes,omitempty" db:"notes"`
	DiagnosedBy        *string   `json:"diagnosedBy,omitempty" db:"diagnosed_by"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
}

// PrescriptionItem item for medication ordering
type PrescriptionItem struct {
	ID                 string    `json:"id" db:"id"`
	PrescriptionID     string    `json:"prescriptionId" db:"prescription_id"`
	DrugName           string    `json:"drugName" db:"drug_name"`
	DosageForm         string    `json:"dosageForm" db:"dosage_form"` // TABLET, CAPSULE, SYRUP, INJECTION, OINTMENT
	Strength           string    `json:"strength" db:"strength"`
	Route              string    `json:"route" db:"route"`
	Frequency          string    `json:"frequency" db:"frequency"`
	DurationDays       int       `json:"durationDays" db:"duration_days"`
	QuantityPrescribed int       `json:"quantityPrescribed" db:"quantity_prescribed"`
	Instructions       string    `json:"instructions" db:"instructions"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
}

// Prescription represents a complete electronic prescription order
type Prescription struct {
	ID                 string             `json:"id" db:"id"`
	EncounterID        string             `json:"encounterId" db:"encounter_id"`
	PatientID          string             `json:"patientId" db:"patient_id"`
	PrescriberID       string             `json:"prescriberId" db:"prescriber_id"`
	PrescriptionNumber string             `json:"prescriptionNumber" db:"prescription_number"`
	Status             string             `json:"status" db:"status"` // PENDING_DISPENSE, PARTIALLY_DISPENSED, DISPENSED, CANCELLED
	Notes              *string            `json:"notes,omitempty" db:"notes"`
	Items              []PrescriptionItem `json:"items,omitempty"`
	CreatedAt          time.Time          `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time          `json:"updatedAt" db:"updated_at"`
}

// StartEncounterPayload initiates a new clinical encounter
type StartEncounterPayload struct {
	CareRequestID    *string `json:"careRequestId"`
	AppointmentID    *string `json:"appointmentId"`
	PatientID        string  `json:"patientId" validate:"required"`
	ProviderID       string  `json:"providerId" validate:"required"`
	EncounterChannel string  `json:"encounterChannel"` // in_person, video, telephone, secure_message
	ChiefComplaint   *string `json:"chiefComplaint"`
}

// UpdateSOAPNotesPayload updates physician clinical notes
type UpdateSOAPNotesPayload struct {
	ChiefComplaint       *string `json:"chiefComplaint"`
	Subjective           *string `json:"subjective"`
	Objective            *string `json:"objective"`
	Assessment           *string `json:"assessment"`
	Plan                 *string `json:"plan"`
	PrimaryDiagnosisCode *string `json:"primaryDiagnosisCode"`
	PrimaryDiagnosisName *string `json:"primaryDiagnosisName"`
	SecondaryDiagnoses   []string `json:"secondaryDiagnoses"`
}

// AddDiagnosisPayload adds an ICD-10 diagnosis to an encounter
type AddDiagnosisPayload struct {
	ICD10Code          string  `json:"icd10Code" validate:"required"`
	ICD10Title         string  `json:"icd10Title" validate:"required"`
	IsPrimary          bool    `json:"isPrimary"`
	ClinicalStatus     string  `json:"clinicalStatus"`
	VerificationStatus string  `json:"verificationStatus"`
	Notes              *string `json:"notes"`
}

// CreatePrescriptionPayload creates an e-prescription order
type CreatePrescriptionPayload struct {
	Notes *string                     `json:"notes"`
	Items []CreatePrescriptionItemDTO `json:"items" validate:"required,min=1"`
}

type CreatePrescriptionItemDTO struct {
	DrugName           string `json:"drugName" validate:"required"`
	DosageForm         string `json:"dosageForm" validate:"required"`
	Strength           string `json:"strength"`
	Route              string `json:"route"`
	Frequency          string `json:"frequency" validate:"required"`
	DurationDays       int    `json:"durationDays" validate:"required,min=1"`
	QuantityPrescribed int    `json:"quantityPrescribed" validate:"required,min=1"`
	Instructions       string `json:"instructions"`
}

// CompleteEncounterResponse returned when an encounter is finalized
type CompleteEncounterResponse struct {
	EncounterID   string  `json:"encounterId"`
	Status        string  `json:"status"`
	InvoiceID     *string `json:"invoiceId,omitempty"`
	InvoiceNumber *string `json:"invoiceNumber,omitempty"`
	TotalAmount   float64 `json:"totalAmount"`
	SignedAt      string  `json:"signedAt"`
}

// DispatchOrdersPayload routes lab, radiology, and prescription orders
type DispatchOrdersPayload struct {
	LabOrders        []string `json:"labOrders"`
	RadiologyOrders  []string `json:"radiologyOrders"`
	PrescriptionList []string `json:"prescriptionList"`
}
