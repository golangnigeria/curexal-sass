package patient

import "time"

// PatientProfile is the database representation of the patient_profile table.
type PatientProfile struct {
	ID                   string     `json:"id"                   db:"id"`
	UserID               string     `json:"userId"               db:"user_id"`
	Phone                *string    `json:"phone"                db:"phone"`
	DateOfBirth          *time.Time `json:"dateOfBirth"          db:"date_of_birth"`
	Gender               *string    `json:"gender"               db:"gender"`
	BloodGroup           *string    `json:"bloodGroup"           db:"blood_group"`
	Genotype             *string    `json:"genotype"             db:"genotype"`
	Address              *string    `json:"address"              db:"address"`
	City                 *string    `json:"city"                 db:"city"`
	State                *string    `json:"state"                db:"state"`
	Country              string     `json:"country"              db:"country"`
	EmergencyContactName *string    `json:"emergencyContactName" db:"emergency_contact_name"`
	EmergencyContactPhone *string   `json:"emergencyContactPhone" db:"emergency_contact_phone"`
	CreatedAt            time.Time  `json:"createdAt"            db:"created_at"`
	UpdatedAt            time.Time  `json:"updatedAt"            db:"updated_at"`
}

// RegisterPatientPayload is the request payload for POST /auth/sign-up (patient).
type RegisterPatientPayload struct {
	Name     string  `json:"name"     validate:"required,min=2,max=120"`
	Email    string  `json:"email"    validate:"required,email"`
	Password string  `json:"password" validate:"required,min=8"`
	Phone    *string `json:"phone"    validate:"omitempty,min=7,max=20"`
}

// UpdatePatientProfilePayload is the request payload for PUT /patient/profile.
type UpdatePatientProfilePayload struct {
	Phone                 *string `json:"phone"`
	DateOfBirth           *string `json:"dateOfBirth"` // "YYYY-MM-DD"
	Gender                *string `json:"gender" validate:"omitempty,oneof=male female other"`
	BloodGroup            *string `json:"bloodGroup" validate:"omitempty,oneof=A+ A- B+ B- O+ O- AB+ AB-"`
	Genotype              *string `json:"genotype" validate:"omitempty,oneof=AA AS SS AC"`
	Address               *string `json:"address"`
	City                  *string `json:"city"`
	State                 *string `json:"state"`
	Country               *string `json:"country"`
	EmergencyContactName  *string `json:"emergencyContactName"`
	EmergencyContactPhone *string `json:"emergencyContactPhone"`
}
