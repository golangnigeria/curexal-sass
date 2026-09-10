package patient

import "time"

// EmergencyContactPayload for next of kin and emergency contact
type EmergencyContactPayload struct {
	FullName     string  `json:"fullName" validate:"required"`
	Relationship string  `json:"relationship" validate:"required"` // SPOUSE, PARENT, CHILD, SIBLING, GUARDIAN, NEXT_OF_KIN, OTHER
	Phone        string  `json:"phone" validate:"required"`
	Email        *string `json:"email,omitempty"`
	Address      *string `json:"address,omitempty"`
}

// RegisterCanonicalPatientPayload is the request payload for registering a new patient
type RegisterCanonicalPatientPayload struct {
	FirstName           string                   `json:"firstName" validate:"required,min=2,max=100"`
	MiddleName          *string                  `json:"middleName,omitempty"`
	LastName            string                   `json:"lastName" validate:"required,min=2,max=100"`
	PreferredName       *string                  `json:"preferredName,omitempty"`
	Gender              string                   `json:"gender" validate:"required,oneof=MALE FEMALE OTHER male female other"`
	DateOfBirth         string                   `json:"dateOfBirth" validate:"required"` // "YYYY-MM-DD"
	Phone               string                   `json:"phone" validate:"required,min=7,max=20"`
	Email               *string                  `json:"email,omitempty" validate:"omitempty,email"`
	BloodGroup          *string                  `json:"bloodGroup,omitempty"`
	Genotype            *string                  `json:"genotype,omitempty"`
	MaritalStatus       *string                  `json:"maritalStatus,omitempty"`
	Occupation          *string                  `json:"occupation,omitempty"`
	NIN                 *string                  `json:"nin,omitempty"`
	Address             *string                  `json:"address,omitempty"`
	City                *string                  `json:"city,omitempty"`
	State               *string                  `json:"state,omitempty"`
	Country             *string                  `json:"country,omitempty"`
	PreferredLanguage   *string                  `json:"preferredLanguage,omitempty"`
	EmergencyContact    *EmergencyContactPayload `json:"emergencyContact,omitempty"`
	RegistrationChannel string                   `json:"registrationChannel"` // RECEPTION, PORTAL, TELEHEALTH, EMERGENCY
	ForceRegistration   bool                     `json:"forceRegistration"`   // Proceed even if probable duplicates exist
}

// SendPortalOTPPayload request payload to send a login OTP
type SendPortalOTPPayload struct {
	Identifier string `json:"identifier" validate:"required"` // Phone or Email
}

// VerifyPortalOTPPayload request payload to verify OTP
type VerifyPortalOTPPayload struct {
	Identifier string `json:"identifier" validate:"required"`
	Code       string `json:"code" validate:"required,len=6"`
}

// SetPortalPINPayload request payload to set 4-6 digit quick login PIN
type SetPortalPINPayload struct {
	PIN string `json:"pin" validate:"required,min=4,max=6"`
}

// PortalLoginPayload request payload for direct PIN or password login
type PortalLoginPayload struct {
	Identifier string `json:"identifier" validate:"required"`
	Secret     string `json:"secret" validate:"required"`
}

// PortalAuthResult returned after successful portal authentication
type PortalAuthResult struct {
	Success      bool      `json:"success"`
	Message      string    `json:"message"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken,omitempty"`
	ExpiresAt    time.Time `json:"expiresAt"`
	Patient      *Patient  `json:"patient"`
}

// PortalAuthResponse returned after successful portal authentication
type PortalAuthResponse struct {
	Token     string    `json:"token"`
	PatientID string    `json:"patientId"`
	MRN       string    `json:"mrn"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Status    string    `json:"status"`
	HasPIN    bool      `json:"hasPin"`
	ExpiresAt time.Time `json:"expiresAt"`
}

// PatientListFilter query parameters for patient directory search
type PatientListFilter struct {
	Query  string `query:"q"`
	Status string `query:"status"`
	Gender string `query:"gender"`
	Limit  int    `query:"limit"`
	Offset int    `query:"offset"`
}

// PatientListResponse paginated list of canonical patients
type PatientListResponse struct {
	Items  []Patient `json:"items"`
	Total  int       `json:"total"`
	Limit  int       `json:"limit"`
	Offset int       `json:"offset"`
}
