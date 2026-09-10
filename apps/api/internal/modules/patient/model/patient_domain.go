package patient

import (
	"time"
)

// Patient represents a canonical patient identity within a tenant boundary
type Patient struct {
	ID                  string                 `json:"id" db:"id"`
	UserID              *string                `json:"userId,omitempty" db:"user_id"`
	TenantID            string                 `json:"tenantId" db:"tenant_id"`
	MRN                 string                 `json:"mrn" db:"mrn"`
	FirstName           string                 `json:"firstName" db:"first_name"`
	MiddleName          *string                `json:"middleName,omitempty" db:"middle_name"`
	LastName            string                 `json:"lastName" db:"last_name"`
	PreferredName       *string                `json:"preferredName,omitempty" db:"preferred_name"`
	Gender              string                 `json:"gender" db:"gender"`
	DateOfBirth         time.Time              `json:"dateOfBirth" db:"date_of_birth"`
	BloodGroup          *string                `json:"bloodGroup,omitempty" db:"blood_group"`
	Genotype            *string                `json:"genotype,omitempty" db:"genotype"`
	MaritalStatus       *string                `json:"maritalStatus,omitempty" db:"marital_status"`
	Occupation          *string                `json:"occupation,omitempty" db:"occupation"`
	NIN                 *string                `json:"nin,omitempty" db:"nin"`
	ResidentialAddress  *string                `json:"residentialAddress,omitempty" db:"residential_address"`
	City                *string                `json:"city,omitempty" db:"city"`
	State               *string                `json:"state,omitempty" db:"state"`
	Country             string                 `json:"country" db:"country"`
	PreferredLanguage   string                 `json:"preferredLanguage" db:"preferred_language"`
	Status              string                 `json:"status" db:"status"`                            // REGISTERED, ACTIVE, INACTIVE, DECEASED
	RegistrationChannel string                 `json:"registrationChannel" db:"registration_channel"` // RECEPTION, PORTAL, TELEHEALTH, EMERGENCY
	Metadata            map[string]interface{} `json:"metadata" db:"metadata"`
	CreatedAt           time.Time              `json:"createdAt" db:"created_at"`
	UpdatedAt           time.Time              `json:"updatedAt" db:"updated_at"`

	Contacts  []PatientContact  `json:"contacts,omitempty"`
	Guardians []PatientGuardian `json:"guardians,omitempty"`
	Portal    *PortalAccount    `json:"portal,omitempty"`
}

// PatientGuardian represents next of kin or emergency contact
type PatientGuardian struct {
	ID                 string    `json:"id" db:"id"`
	PatientID          string    `json:"patientId" db:"patient_id"`
	RelationshipType   string    `json:"relationshipType" db:"relationship_type"`
	FullName           string    `json:"fullName" db:"full_name"`
	Phone              string    `json:"phone" db:"phone"`
	Email              *string   `json:"email,omitempty" db:"email"`
	Address            *string   `json:"address,omitempty" db:"address"`
	IsEmergencyContact bool      `json:"isEmergencyContact" db:"is_emergency_contact"`
	CreatedAt          time.Time `json:"createdAt" db:"created_at"`
}

// PatientContact represents a phone, email, or messaging address for a patient
type PatientContact struct {
	ID         string     `json:"id" db:"id"`
	PatientID  string     `json:"patientId" db:"patient_id"`
	System     string     `json:"system" db:"system"` // PHONE, EMAIL, WHATSAPP
	Value      string     `json:"value" db:"value"`
	UseType    string     `json:"useType" db:"use_type"` // MOBILE, HOME, WORK, EMERGENCY
	IsPrimary  bool       `json:"isPrimary" db:"is_primary"`
	VerifiedAt *time.Time `json:"verifiedAt,omitempty" db:"verified_at"`
	CreatedAt  time.Time  `json:"createdAt" db:"created_at"`
}

// PortalAccount represents the patient portal authentication identity
type PortalAccount struct {
	ID           string     `json:"id" db:"id"`
	PatientID    string     `json:"patientId" db:"patient_id"`
	TenantID     string     `json:"tenantId" db:"tenant_id"`
	Identifier   string     `json:"identifier" db:"identifier"` // Phone or Email
	PasswordHash *string    `json:"-" db:"password_hash"`
	PINHash      *string    `json:"-" db:"pin_hash"`
	Status       string     `json:"status" db:"status"` // INVITED, PENDING_VERIFICATION, ACTIVE, LOCKED, SUSPENDED
	MFAEnabled   bool       `json:"mfaEnabled" db:"mfa_enabled"`
	LastLoginAt  *time.Time `json:"lastLoginAt,omitempty" db:"last_login_at"`
	CreatedAt    time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time  `json:"updatedAt" db:"updated_at"`
}

// Consent represents a patient privacy/consent record
type Consent struct {
	ID          string                 `json:"id" db:"id"`
	PatientID   string                 `json:"patientId" db:"patient_id"`
	ConsentType string                 `json:"consentType" db:"consent_type"` // TELEHEALTH, DATA_SHARING, RESEARCH, PROXY_ACCESS
	Status      string                 `json:"status" db:"status"`            // ACTIVE, REVOKED, EXPIRED
	GrantedBy   string                 `json:"grantedBy" db:"granted_by"`     // SELF, LEGAL_GUARDIAN, PROXY
	GrantedAt   time.Time              `json:"grantedAt" db:"granted_at"`
	ExpiresAt   *time.Time             `json:"expiresAt,omitempty" db:"expires_at"`
	Metadata    map[string]interface{} `json:"metadata" db:"metadata"`
}
