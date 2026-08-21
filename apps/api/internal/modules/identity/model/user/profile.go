package user

import (
	"time"
)

type UserProfile struct {
	UserID           string     `json:"userId"      db:"user_id"`
	FirstName        *string    `json:"firstName"   db:"first_name"`
	LastName         *string    `json:"lastName"    db:"last_name"`
	PhoneNumber      *string    `json:"phoneNumber" db:"phone_number"`
	Bio              *string    `json:"bio"         db:"bio"`
	MiddleName       *string    `json:"middleName"   db:"middle_name"`
	AvatarURL        *string    `json:"avatarUrl"    db:"avatar_url"`
	Gender           *string    `json:"gender"       db:"gender"`
	DateOfBirth      *time.Time `json:"dateOfBirth"  db:"date_of_birth"`
	Nationality      *string    `json:"nationality"  db:"nationality"`
	PrimaryPhone     *string    `json:"primaryPhone" db:"primary_phone"`
	Timezone         *string    `json:"timezone"     db:"timezone"`
	LanguageCode     *string    `json:"languageCode" db:"language_code"`
	EmergencyContact *string    `json:"emergencyContact" db:"emergency_contact"`
}

type ProfessionalProfile struct {
	ID                 string    `json:"id"                 db:"id"`
	UserID             string    `json:"userId"             db:"user_id"`
	Profession         string    `json:"profession"         db:"profession"`
	Specialty          string    `json:"specialty"          db:"specialty"`
	SubSpecialty       *string   `json:"subSpecialty"       db:"sub_specialty"`
	RegistrationNumber string    `json:"registrationNumber" db:"registration_number"`
	LicensingBody      string    `json:"licensingBody"      db:"licensing_body"`
	LicenseIssueDate   time.Time `json:"licenseIssueDate"   db:"license_issue_date"`
	LicenseExpiry      time.Time `json:"licenseExpiry"      db:"license_expiry"`
	Status             string    `json:"status"             db:"status"`
	CreatedAt          time.Time `json:"createdAt"          db:"created_at"`
	UpdatedAt          time.Time `json:"updatedAt"          db:"updated_at"`
}

type ProfessionalSignature struct {
	ID                string    `json:"id"                db:"id"`
	UserID            string    `json:"userId"            db:"user_id"`
	TenantID          string    `json:"tenantId"          db:"tenant_id"`
	SignatureImageUrl string    `json:"signatureImageUrl" db:"signature_image_url"`
	SignatureHash     string    `json:"signatureHash"     db:"signature_hash"`
	CertificateSerial *string   `json:"certificateSerial" db:"certificate_serial"`
	SigningAlgorithm  string    `json:"signingAlgorithm"  db:"signing_algorithm"`
	IsActive          bool      `json:"isActive"          db:"is_active"`
	ExpiresAt         time.Time `json:"expiresAt"         db:"expires_at"`
	CreatedAt         time.Time `json:"createdAt"         db:"created_at"`
}

type ProfessionalVerification struct {
	ID                    string    `json:"id"                    db:"id"`
	ProfessionalProfileID string    `json:"professionalProfileId" db:"professional_profile_id"`
	VerifiedBy            *string   `json:"verifiedBy"            db:"verified_by"`
	VerifiedAt            time.Time `json:"verifiedAt"            db:"verified_at"`
	VerificationSource    string    `json:"verificationSource"    db:"verification_source"`
	AuditNotes            *string   `json:"auditNotes"            db:"audit_notes"`
	ReferenceNumber       *string   `json:"referenceNumber"       db:"reference_number"`
	Status                string    `json:"status"                db:"status"`
}

type TenantStaff struct {
	ID               string    `json:"id"               db:"id"`
	TenantID         string    `json:"tenantId"         db:"tenant_id"`
	UserID           string    `json:"userId"           db:"user_id"`
	EmployeeNumber   string    `json:"employeeNumber"   db:"employee_number"`
	JobTitle         *string   `json:"jobTitle"         db:"job_title"`
	EmploymentType   string    `json:"employmentType"   db:"employment_type"`
	DateJoined       time.Time `json:"dateJoined"       db:"date_joined"`
	EmploymentStatus string    `json:"employmentStatus" db:"employment_status"`
	ManagerID        *string   `json:"managerId"        db:"manager_id"`
	BankName         *string   `json:"bankName"         db:"bank_name"`
	AccountNumber    *string   `json:"accountNumber"    db:"account_number"`
	TaxID            *string   `json:"taxId"            db:"tax_id"`
	SalaryRate       *float64  `json:"salaryRate"       db:"salary_rate"`
	CreatedAt        time.Time `json:"createdAt"        db:"created_at"`
	UpdatedAt        time.Time `json:"updatedAt"        db:"updated_at"`
}

type ClinicalPrivilege struct {
	ID            string     `json:"id"            db:"id"`
	StaffID       string     `json:"staffId"       db:"staff_id"`
	PrivilegeKey  string     `json:"privilegeKey"  db:"privilege_key"`
	EffectiveDate time.Time  `json:"effectiveDate" db:"effective_date"`
	ExpiryDate    *time.Time `json:"expiryDate"    db:"expiry_date"`
	ApprovedBy    *string    `json:"approvedBy"    db:"approved_by"`
}

type ProfessionalCompetency struct {
	ID                  string    `json:"id"                  db:"id"`
	StaffID             string    `json:"staffId"             db:"staff_id"`
	CompetencyType      string    `json:"competencyType"      db:"competency_type"`
	AssessmentDate      time.Time `json:"assessmentDate"      db:"assessment_date"`
	ReassessmentDueDate time.Time `json:"reassessmentDueDate" db:"reassessment_due_date"`
	AssessedBy          string    `json:"assessedBy"          db:"assessed_by"`
	Status              string    `json:"status"              db:"status"`
	EvidenceDocumentURL *string   `json:"evidenceDocumentUrl" db:"evidence_document_url"`
	CreatedAt           time.Time `json:"createdAt"           db:"created_at"`
	UpdatedAt           time.Time `json:"updatedAt"           db:"updated_at"`
}
