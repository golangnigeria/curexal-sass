package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	patientModel "github.com/golangnigeria/curexal/internal/modules/patient/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CanonicalPatientRepository struct {
	server *server.Server
}

func NewCanonicalPatientRepository(s *server.Server) *CanonicalPatientRepository {
	return &CanonicalPatientRepository{server: s}
}

// FindCandidatesForDuplicateResolution queries existing tenant patients that share phone, NIN, or last name + DOB
func (r *CanonicalPatientRepository) FindCandidatesForDuplicateResolution(
	ctx context.Context,
	tenantID string,
	phone, nin, lastName string,
	dob *time.Time,
) ([]patientModel.Patient, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT DISTINCT p.id, p.tenant_id, p.mrn, p.first_name, p.middle_name, p.last_name, 
		       p.gender, p.date_of_birth, p.blood_group, p.genotype, p.nin, p.status, 
		       p.registration_channel, p.metadata, p.created_at, p.updated_at
		FROM patient.patients p
		LEFT JOIN patient.patient_contacts c ON c.patient_id = p.id
		WHERE p.tenant_id = $1
		  AND (
		      ($2 != '' AND c.value = $2)
		      OR ($3 != '' AND p.nin = $3)
		      OR (LOWER(p.last_name) = LOWER($4) AND $5::date IS NOT NULL AND p.date_of_birth = $5::date)
		  )
		LIMIT 10
	`

	cleanPhone := strings.TrimSpace(phone)
	cleanNIN := strings.TrimSpace(nin)
	cleanLastName := strings.TrimSpace(lastName)

	var dobParam *time.Time
	if dob != nil && !dob.IsZero() {
		dobParam = dob
	}

	rows, err := db.Query(ctx, query, tenantID, cleanPhone, cleanNIN, cleanLastName, dobParam)
	if err != nil {
		return nil, fmt.Errorf("error querying duplicate candidates: %w", err)
	}
	defer rows.Close()

	var patients []patientModel.Patient
	for rows.Next() {
		var p patientModel.Patient
		err := rows.Scan(
			&p.ID, &p.TenantID, &p.MRN, &p.FirstName, &p.MiddleName, &p.LastName,
			&p.Gender, &p.DateOfBirth, &p.BloodGroup, &p.Genotype, &p.NIN, &p.Status,
			&p.RegistrationChannel, &p.Metadata, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning candidate: %w", err)
		}
		patients = append(patients, p)
	}

	return patients, nil
}

// CreatePatient inserts a canonical patient record along with contacts inside a transaction
func (r *CanonicalPatientRepository) CreatePatient(
	ctx context.Context,
	p *patientModel.Patient,
	contacts []patientModel.PatientContact,
) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}

	if p.ID == "" {
		p.ID = uuid.New().String()
	}

	return r.server.DB.RunInTx(ctx, func(txCtx context.Context) error {
		db := r.server.DB.Conn(txCtx)

		// 1. Insert patient
		insertPatientQuery := `
			INSERT INTO patient.patients (
				id, user_id, tenant_id, organization_id, mrn, first_name, middle_name, last_name,
				gender, date_of_birth, blood_group, genotype, nin, status,
				registration_channel, metadata, created_at, updated_at
			) VALUES (
				$1, $2, $3, (SELECT organization_id FROM organization.facility_branches WHERE id = $3 LIMIT 1), $4, $5, $6, $7,
				$8, $9, $10, $11, $12, $13,
				$14, $15, NOW(), NOW()
			)
		`
		_, err := db.Exec(
			txCtx, insertPatientQuery,
			p.ID, p.UserID, p.TenantID, p.MRN, p.FirstName, p.MiddleName, p.LastName,
			p.Gender, p.DateOfBirth, p.BloodGroup, p.Genotype, p.NIN, p.Status,
			p.RegistrationChannel, p.Metadata,
		)
		if err != nil {
			return fmt.Errorf("failed to insert patient: %w", err)
		}

		// 2. Insert contacts
		insertContactQuery := `
			INSERT INTO patient.patient_contacts (
				id, patient_id, system, value, use_type, is_primary, verified_at, created_at
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, NOW()
			)
		`
		for _, c := range contacts {
			cid := c.ID
			if cid == "" {
				cid = uuid.New().String()
			}
			_, err := db.Exec(
				txCtx, insertContactQuery,
				cid, p.ID, c.System, c.Value, c.UseType, c.IsPrimary, c.VerifiedAt,
			)
			if err != nil {
				return fmt.Errorf("failed to insert contact %s: %w", c.Value, err)
			}
		}

		return nil
	})
}

// GetPatientByID retrieves a patient and their contacts by ID within a tenant
func (r *CanonicalPatientRepository) GetPatientByID(ctx context.Context, tenantID, patientID string) (*patientModel.Patient, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT id, user_id, tenant_id, mrn, first_name, middle_name, last_name, 
		       gender, date_of_birth, blood_group, genotype, nin, status, 
		       registration_channel, metadata, created_at, updated_at
		FROM patient.patients
		WHERE (tenant_id = $1 OR $1 = '') AND id = $2
	`
	var p patientModel.Patient
	err := db.QueryRow(ctx, query, tenantID, patientID).Scan(
		&p.ID, &p.UserID, &p.TenantID, &p.MRN, &p.FirstName, &p.MiddleName, &p.LastName,
		&p.Gender, &p.DateOfBirth, &p.BloodGroup, &p.Genotype, &p.NIN, &p.Status,
		&p.RegistrationChannel, &p.Metadata, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	// Fetch contacts
	contactsQuery := `
		SELECT id, patient_id, system, value, use_type, is_primary, verified_at, created_at
		FROM patient.patient_contacts
		WHERE patient_id = $1
		ORDER BY is_primary DESC, created_at ASC
	`
	cRows, err := db.Query(ctx, contactsQuery, p.ID)
	if err == nil {
		defer cRows.Close()
		for cRows.Next() {
			var c patientModel.PatientContact
			if scanErr := cRows.Scan(&c.ID, &c.PatientID, &c.System, &c.Value, &c.UseType, &c.IsPrimary, &c.VerifiedAt, &c.CreatedAt); scanErr == nil {
				p.Contacts = append(p.Contacts, c)
			}
		}
	}

	return &p, nil
}

// GetPatientByUserID retrieves a patient profile by identity.users ID
func (r *CanonicalPatientRepository) GetPatientByUserID(ctx context.Context, userID string) (*patientModel.Patient, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT id, user_id, tenant_id, mrn, first_name, middle_name, last_name, 
		       gender, date_of_birth, blood_group, genotype, nin, status, 
		       registration_channel, metadata, created_at, updated_at
		FROM patient.patients
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var p patientModel.Patient
	err := db.QueryRow(ctx, query, userID).Scan(
		&p.ID, &p.UserID, &p.TenantID, &p.MRN, &p.FirstName, &p.MiddleName, &p.LastName,
		&p.Gender, &p.DateOfBirth, &p.BloodGroup, &p.Genotype, &p.NIN, &p.Status,
		&p.RegistrationChannel, &p.Metadata, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

// ListPatients lists and searches patients with pagination
func (r *CanonicalPatientRepository) ListPatients(
	ctx context.Context,
	tenantID string,
	filter patientModel.PatientListFilter,
) ([]patientModel.Patient, int, error) {
	if r.server.DB == nil {
		return nil, 0, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	whereClause := "WHERE p.tenant_id = $1"
	args := []interface{}{tenantID}
	argIndex := 2

	if filter.Query != "" {
		q := "%" + strings.ToLower(filter.Query) + "%"
		whereClause += fmt.Sprintf(" AND (LOWER(p.first_name) LIKE $%d OR LOWER(p.last_name) LIKE $%d OR LOWER(p.mrn) LIKE $%d)", argIndex, argIndex, argIndex)
		args = append(args, q)
		argIndex++
	}

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND p.status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.Gender != "" {
		whereClause += fmt.Sprintf(" AND p.gender = $%d", argIndex)
		args = append(args, strings.ToUpper(filter.Gender))
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM patient.patients p %s", whereClause)
	var total int
	if err := db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT p.id, p.user_id, p.tenant_id, p.mrn, p.first_name, p.middle_name, p.last_name, 
		       p.gender, p.date_of_birth, p.blood_group, p.genotype, p.nin, p.status, 
		       p.registration_channel, p.metadata, p.created_at, p.updated_at
		FROM patient.patients p
		%s
		ORDER BY p.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []patientModel.Patient
	for rows.Next() {
		var p patientModel.Patient
		err := rows.Scan(
			&p.ID, &p.UserID, &p.TenantID, &p.MRN, &p.FirstName, &p.MiddleName, &p.LastName,
			&p.Gender, &p.DateOfBirth, &p.BloodGroup, &p.Genotype, &p.NIN, &p.Status,
			&p.RegistrationChannel, &p.Metadata, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, p)
	}

	return items, total, nil
}

// GetPortalAccountByIdentifier retrieves a portal account by phone or email within tenant
func (r *CanonicalPatientRepository) GetPortalAccountByIdentifier(ctx context.Context, tenantID, identifier string) (*patientModel.PortalAccount, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT id, patient_id, tenant_id, identifier, password_hash, pin_hash, status, mfa_enabled, last_login_at, created_at, updated_at
		FROM patient.portal_accounts
		WHERE tenant_id = $1 AND identifier = $2
	`
	var acc patientModel.PortalAccount
	err := db.QueryRow(ctx, query, tenantID, strings.TrimSpace(identifier)).Scan(
		&acc.ID, &acc.PatientID, &acc.TenantID, &acc.Identifier, &acc.PasswordHash,
		&acc.PINHash, &acc.Status, &acc.MFAEnabled, &acc.LastLoginAt, &acc.CreatedAt, &acc.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil
}

// UpsertPortalAccount creates or updates a portal account
func (r *CanonicalPatientRepository) UpsertPortalAccount(ctx context.Context, acc *patientModel.PortalAccount) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	if acc.ID == "" {
		acc.ID = uuid.New().String()
	}
	query := `
		INSERT INTO patient.portal_accounts (
			id, patient_id, tenant_id, identifier, password_hash, pin_hash, status, mfa_enabled, last_login_at, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW()
		)
		ON CONFLICT (tenant_id, identifier) DO UPDATE SET
			pin_hash = COALESCE(EXCLUDED.pin_hash, patient.portal_accounts.pin_hash),
			status = EXCLUDED.status,
			last_login_at = COALESCE(EXCLUDED.last_login_at, patient.portal_accounts.last_login_at),
			updated_at = NOW()
	`
	_, err := db.Exec(
		ctx, query,
		acc.ID, acc.PatientID, acc.TenantID, acc.Identifier, acc.PasswordHash,
		acc.PINHash, acc.Status, acc.MFAEnabled, acc.LastLoginAt,
	)
	return err
}
