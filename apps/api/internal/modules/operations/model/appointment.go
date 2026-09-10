package model

import "time"

// Appointment represents a scheduled clinical visit across in-person or telehealth delivery channels
type Appointment struct {
	ID                 string    `json:"id" db:"id"`
	TenantID           string    `json:"tenantId" db:"tenant_id"`
	PatientID          string    `json:"patientId" db:"patient_id"`
	ProviderID         string    `json:"providerId" db:"provider_id"`
	AppointmentNumber  string    `json:"appointmentNumber" db:"appointment_number"`
	ServiceType        string    `json:"serviceType" db:"service_type"`
	DeliveryChannel    string    `json:"deliveryChannel" db:"delivery_channel"` // in_person, video, telephone, secure_message
	Status             string    `json:"status" db:"status"`                     // BOOKED, CHECKED_IN, IN_PROGRESS, COMPLETED, CANCELLED, NO_SHOW
	StartTime          time.Time `json:"startTime" db:"start_time"`
	EndTime            time.Time `json:"endTime" db:"end_time"`
	ReasonForVisit     *string   `json:"reasonForVisit,omitempty" db:"reason_for_visit"`
	CancellationReason *string   `json:"cancellationReason,omitempty" db:"cancellation_reason"`
	VirtualMeetingURL  *string   `json:"virtualMeetingUrl,omitempty" db:"virtual_meeting_url"`
	CreatedBy          *string   `json:"createdBy,omitempty" db:"created_by"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time `json:"updatedAt" db:"updated_at"`

	// Enriched fields
	PatientName   *string `json:"patientName,omitempty"`
	PatientMRN    *string `json:"patientMrn,omitempty"`
	ProviderName  *string `json:"providerName,omitempty"`
	SpecialtyCode *string `json:"specialtyCode,omitempty"`
	RoomNumber    *string `json:"roomNumber,omitempty"`
}

// CreateAppointmentRequest payload to book an appointment
type CreateAppointmentRequest struct {
	PatientID       string  `json:"patientId" validate:"required"`
	ProviderID      string  `json:"providerId" validate:"required"`
	StartTime       string  `json:"startTime" validate:"required"` // RFC3339
	EndTime         string  `json:"endTime" validate:"required"`   // RFC3339
	ServiceType     string  `json:"serviceType"`                   // Default: CONSULTATION
	DeliveryChannel string  `json:"deliveryChannel"`               // in_person, video, telephone, secure_message
	ReasonForVisit  *string `json:"reasonForVisit"`
}

// UpdateAppointmentStatusRequest payload to transition appointment lifecycle
type UpdateAppointmentStatusRequest struct {
	Status             string  `json:"status" validate:"required"` // BOOKED, CHECKED_IN, IN_PROGRESS, COMPLETED, CANCELLED, NO_SHOW
	CancellationReason *string `json:"cancellationReason,omitempty"`
}

// AppointmentFilter query options
type AppointmentFilter struct {
	ProviderID      *string `query:"providerId"`
	PatientID       *string `query:"patientId"`
	Status          *string `query:"status"`
	DeliveryChannel *string `query:"deliveryChannel"`
	Date            *string `query:"date"` // YYYY-MM-DD
	Limit           int     `query:"limit"`
	Offset          int     `query:"offset"`
}
