package repository

import (
	"context"
	"fmt"
	"time"

	identityModel "github.com/golangnigeria/curexal/internal/modules/identity/model"
	patientModel "github.com/golangnigeria/curexal/internal/modules/patient/model"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/jackc/pgx/v5"
)

type PatientRepository struct {
	server *server.Server
}

func NewPatientRepository(s *server.Server) *PatientRepository {
	return &PatientRepository{server: s}
}

func (r *PatientRepository) CreateProfile(ctx context.Context, userID string, phone *string) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `
		INSERT INTO patient.patient_profiles (user_id, phone)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO NOTHING
	`, userID, phone)
	return err
}

func (r *PatientRepository) GetProfileByUserID(ctx context.Context, userID string) (*patientModel.PatientProfile, error) {
	db := r.server.DB.Conn(ctx)
	p := &patientModel.PatientProfile{}
	err := db.QueryRow(ctx, `
		SELECT id, user_id, phone, date_of_birth, gender, blood_group, genotype, address, city, state, country, emergency_contact_name, emergency_contact_phone, created_at, updated_at
		FROM patient.patient_profiles
		WHERE user_id = $1
	`, userID).Scan(
		&p.ID, &p.UserID, &p.Phone, &p.DateOfBirth, &p.Gender, &p.BloodGroup, &p.Genotype,
		&p.Address, &p.City, &p.State, &p.Country, &p.EmergencyContactName, &p.EmergencyContactPhone,
		&p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return p, nil
}

func (r *PatientRepository) GetPatientContext(ctx context.Context, userID string) (*identityModel.PatientContext, error) {
	db := r.server.DB.Conn(ctx)
	var profileID string
	var phone, gender, bloodGroup, genotype, city, state *string
	var dob *time.Time
	var country string

	err := db.QueryRow(ctx, `
		SELECT id, phone, date_of_birth, gender, blood_group, genotype, city, state, country
		FROM patient.patient_profiles
		WHERE user_id = $1
	`, userID).Scan(&profileID, &phone, &dob, &gender, &bloodGroup, &genotype, &city, &state, &country)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	var dobStr *string
	if dob != nil {
		s := dob.Format("2006-01-02")
		dobStr = &s
	}

	return &identityModel.PatientContext{
		ProfileID:   profileID,
		Phone:       phone,
		DateOfBirth: dobStr,
		Gender:      gender,
		BloodGroup:  bloodGroup,
		Genotype:    genotype,
		City:        city,
		State:       state,
		Country:     country,
	}, nil
}

func (r *PatientRepository) ProfileExists(ctx context.Context, userID string) (bool, string, error) {
	db := r.server.DB.Conn(ctx)
	var id string
	err := db.QueryRow(ctx, `SELECT id FROM patient.patient_profiles WHERE user_id = $1`, userID).Scan(&id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, "", nil
		}
		return false, "", err
	}
	return true, id, nil
}

func (r *PatientRepository) UpdateProfile(ctx context.Context, userID string, payload patientModel.UpdatePatientProfilePayload) error {
	db := r.server.DB.Conn(ctx)

	var dobTime *time.Time
	if payload.DateOfBirth != nil && *payload.DateOfBirth != "" {
		t, err := time.Parse("2006-01-02", *payload.DateOfBirth)
		if err == nil {
			dobTime = &t
		}
	}

	_, err := db.Exec(ctx, `
		UPDATE patient.patient_profiles
		SET phone = COALESCE($1, phone),
		    date_of_birth = COALESCE($2, date_of_birth),
		    gender = COALESCE($3, gender),
		    blood_group = COALESCE($4, blood_group),
		    genotype = COALESCE($5, genotype),
		    address = COALESCE($6, address),
		    city = COALESCE($7, city),
		    state = COALESCE($8, state),
		    country = COALESCE($9, country),
		    emergency_contact_name = COALESCE($10, emergency_contact_name),
		    emergency_contact_phone = COALESCE($11, emergency_contact_phone),
		    updated_at = NOW()
		WHERE user_id = $12
	`, payload.Phone, dobTime, payload.Gender, payload.BloodGroup, payload.Genotype,
		payload.Address, payload.City, payload.State, payload.Country,
		payload.EmergencyContactName, payload.EmergencyContactPhone, userID)

	if err != nil {
		return fmt.Errorf("failed to update patient profile: %w", err)
	}

	return nil
}
