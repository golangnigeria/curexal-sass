package user

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

type UpdateUserProfilePayload struct {
	UserID           string  `param:"userId" validate:"required"`
	FirstName        *string `json:"firstName" validate:"omitempty,max=60"`
	LastName         *string `json:"lastName" validate:"omitempty,max=60"`
	MiddleName       *string `json:"middleName" validate:"omitempty,max=60"`
	PhoneNumber      *string `json:"phoneNumber" validate:"omitempty,max=30"`
	Bio              *string `json:"bio" validate:"omitempty,max=500"`
	AvatarURL        *string `json:"avatarUrl" validate:"omitempty"`
	Gender           *string `json:"gender" validate:"omitempty,oneof=male female other undisclosed"`
	DateOfBirth      *string `json:"dateOfBirth" validate:"omitempty"`
	Nationality      *string `json:"nationality" validate:"omitempty,max=100"`
	Timezone         *string `json:"timezone" validate:"omitempty"`
	LanguageCode     *string `json:"languageCode" validate:"omitempty"`
	EmergencyContact *string `json:"emergencyContact" validate:"omitempty"`
}

func (p *UpdateUserProfilePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type GetUserProfilePayload struct {
	UserID string `param:"userId" validate:"required"`
}

func (p *GetUserProfilePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type RequestEmailChangePayload struct {
	NewEmail string `json:"newEmail" validate:"required,email"`
}

func (p *RequestEmailChangePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type VerifyEmailChangePayload struct {
	Code  string `json:"code,omitempty"`
	Token string `json:"token,omitempty"`
}

func (p *VerifyEmailChangePayload) Validate() error {
	if strings.TrimSpace(p.Code) == "" && strings.TrimSpace(p.Token) == "" {
		return errors.New("verification code is required")
	}
	return nil
}

type CreateProfessionalProfilePayload struct {
	Profession         string `json:"profession" validate:"required,min=2"`
	Specialty          string `json:"specialty" validate:"required,min=2"`
	SubSpecialty       string `json:"subSpecialty,omitempty"`
	RegistrationNumber string `json:"registrationNumber" validate:"required,min=2"`
	LicensingBody      string `json:"licensingBody" validate:"required,min=2"`
	LicenseIssueDate   string `json:"licenseIssueDate" validate:"required"`
	LicenseExpiry      string `json:"licenseExpiry" validate:"required"`
}

func (p *CreateProfessionalProfilePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type CreateProfessionalSignaturePayload struct {
	SignatureImageUrl string `json:"signatureImageUrl" validate:"required,url"`
	CertificateSerial string `json:"certificateSerial,omitempty"`
	SigningAlgorithm  string `json:"signingAlgorithm,omitempty"`
	ExpiresAt         string `json:"expiresAt" validate:"required"`
}

func (p *CreateProfessionalSignaturePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type UpdateTenantStaffPayload struct {
	JobTitle       string    `json:"jobTitle,omitempty"`
	EmploymentType string    `json:"employmentType,omitempty" validate:"omitempty,oneof=full_time part_time locum contract"`
	DateJoined     string    `json:"dateJoined,omitempty"`
	ManagerID      string    `json:"managerId,omitempty" validate:"omitempty,uuid"`
	DepartmentIDs  []string  `json:"departmentIds,omitempty"`
	BankName       string    `json:"bankName,omitempty"`
	AccountNumber  string    `json:"accountNumber,omitempty"`
	TaxID          string    `json:"taxId,omitempty"`
	SalaryRate     *float64  `json:"salaryRate,omitempty"`
}

func (p *UpdateTenantStaffPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type CreateClinicalPrivilegePayload struct {
	PrivilegeKey string `json:"privilegeKey" validate:"required,min=1"`
	ExpiryDate   string `json:"expiryDate,omitempty"`
}

func (p *CreateClinicalPrivilegePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type CreateProfessionalCompetencyPayload struct {
	CompetencyType      string `json:"competencyType" validate:"required,min=2"`
	AssessmentDate      string `json:"assessmentDate" validate:"required"`
	ReassessmentDueDate string `json:"reassessmentDueDate" validate:"required"`
	AssessedBy          string `json:"assessedBy" validate:"required,uuid"`
	EvidenceDocumentURL string `json:"evidenceDocumentUrl,omitempty" validate:"omitempty,url"`
}

func (p *CreateProfessionalCompetencyPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type CreateProfessionalVerificationPayload struct {
	VerificationSource string `json:"verificationSource" validate:"required,min=2"`
	AuditNotes         string `json:"auditNotes,omitempty"`
	ReferenceNumber    string `json:"referenceNumber,omitempty"`
	Status             string `json:"status,omitempty" validate:"omitempty,oneof=verified rejected"`
}

func (p *CreateProfessionalVerificationPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type CreatePermissionOverridePayload struct {
	MembershipID    string `param:"id" validate:"required,uuid"`
	Permission      string `json:"permission" validate:"required,min=1,max=100"`
	OverrideType    string `json:"overrideType" validate:"required,oneof=grant deny"`
	DurationSeconds *int   `json:"durationSeconds,omitempty" validate:"omitempty,gt=0"`
}

func (p *CreatePermissionOverridePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type GetPermissionOverridesPayload struct {
	MembershipID string `param:"id" validate:"required,uuid"`
}

func (p *GetPermissionOverridesPayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}

type DeletePermissionOverridePayload struct {
	MembershipID string `param:"id" validate:"required,uuid"`
	OverrideID   string `param:"overrideId" validate:"required,uuid"`
}

func (p *DeletePermissionOverridePayload) Validate() error {
	validate := validator.New()
	return validate.Struct(p)
}
