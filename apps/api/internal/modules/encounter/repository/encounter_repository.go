package repository

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/encounter/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type EncounterRepository struct {
	server *server.Server
}

func NewEncounterRepository(s *server.Server) *EncounterRepository {
	return &EncounterRepository{server: s}
}

// CreateEncounter stores a new clinical encounter
func (r *EncounterRepository) CreateEncounter(ctx context.Context, enc *model.Encounter) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	if enc.ID == "" {
		enc.ID = uuid.New().String()
	}

	channel := enc.EncounterChannel
	if channel == "" {
		channel = "in_person"
	}

	status := enc.Status
	if status == "" {
		status = "in_progress"
	}

	secDiagJSON, _ := json.Marshal(enc.SecondaryDiagnoses)
	if len(enc.SecondaryDiagnoses) == 0 {
		secDiagJSON = []byte("[]")
	}

	query := `
		INSERT INTO encounter.encounters (
			id, tenant_id, care_request_id, appointment_id, patient_id, provider_id,
			encounter_channel, status, chief_complaint, subjective,
			objective, assessment, plan, primary_diagnosis_code,
			primary_diagnosis_name, secondary_diagnoses, started_at,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10,
			$11, $12, $13, $14,
			$15, $16, NOW(),
			NOW(), NOW()
		)
	`
	_, err := db.Exec(
		ctx, query,
		enc.ID, enc.TenantID, enc.CareRequestID, enc.AppointmentID, enc.PatientID, enc.ProviderID,
		channel, status, enc.ChiefComplaint, enc.Subjective,
		enc.Objective, enc.Assessment, enc.Plan, enc.PrimaryDiagnosisCode,
		enc.PrimaryDiagnosisName, secDiagJSON,
	)
	return err
}

// GetEncounterByID retrieves an encounter with enriched patient info
func (r *EncounterRepository) GetEncounterByID(ctx context.Context, tenantID, encounterID string) (*model.Encounter, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT e.id, e.tenant_id, e.care_request_id, e.appointment_id, e.patient_id, e.provider_id,
		       e.encounter_channel, e.status, e.chief_complaint, e.subjective,
		       e.objective, e.assessment, e.plan, e.primary_diagnosis_code,
		       e.primary_diagnosis_name, e.secondary_diagnoses, e.started_at,
		       e.signed_at, e.signed_by, e.closed_at, e.created_at, e.updated_at,
		       COALESCE(p.first_name || ' ' || p.last_name, 'Patient') AS patient_name,
		       p.mrn, p.gender
		FROM encounter.encounters e
		LEFT JOIN patient.patients p ON p.id = e.patient_id
		WHERE e.tenant_id = $1 AND e.id = $2
	`
	var enc model.Encounter
	var secDiagRaw []byte
	err := db.QueryRow(ctx, query, tenantID, encounterID).Scan(
		&enc.ID, &enc.TenantID, &enc.CareRequestID, &enc.AppointmentID, &enc.PatientID, &enc.ProviderID,
		&enc.EncounterChannel, &enc.Status, &enc.ChiefComplaint, &enc.Subjective,
		&enc.Objective, &enc.Assessment, &enc.Plan, &enc.PrimaryDiagnosisCode,
		&enc.PrimaryDiagnosisName, &secDiagRaw, &enc.StartedAt,
		&enc.SignedAt, &enc.SignedBy, &enc.ClosedAt, &enc.CreatedAt, &enc.UpdatedAt,
		&enc.PatientName, &enc.MRN, &enc.Gender,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if len(secDiagRaw) > 0 {
		_ = json.Unmarshal(secDiagRaw, &enc.SecondaryDiagnoses)
	}

	// Load diagnoses and prescriptions
	enc.Diagnoses, _ = r.ListDiagnoses(ctx, enc.ID)
	enc.Prescriptions, _ = r.ListPrescriptions(ctx, enc.ID)

	return &enc, nil
}

// ListEncounters retrieves list of encounters for a tenant
func (r *EncounterRepository) ListEncounters(
	ctx context.Context,
	tenantID string,
	status *string,
	providerID *string,
	patientID *string,
) ([]model.Encounter, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT e.id, e.tenant_id, e.care_request_id, e.appointment_id, e.patient_id, e.provider_id,
		       e.encounter_channel, e.status, e.chief_complaint, e.subjective,
		       e.objective, e.assessment, e.plan, e.primary_diagnosis_code,
		       e.primary_diagnosis_name, e.started_at, e.created_at, e.updated_at,
		       COALESCE(p.first_name || ' ' || p.last_name, 'Patient') AS patient_name,
		       p.mrn, p.gender
		FROM encounter.encounters e
		LEFT JOIN patient.patients p ON p.id = e.patient_id
		WHERE e.tenant_id = $1
		  AND ($2::text IS NULL OR e.status = $2)
		  AND ($3::uuid IS NULL OR e.provider_id = $3::uuid)
		  AND ($4::uuid IS NULL OR e.patient_id = $4::uuid)
		ORDER BY e.started_at DESC
		LIMIT 50
	`
	rows, err := db.Query(ctx, query, tenantID, status, providerID, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var encounters []model.Encounter
	for rows.Next() {
		var enc model.Encounter
		err := rows.Scan(
			&enc.ID, &enc.TenantID, &enc.CareRequestID, &enc.AppointmentID, &enc.PatientID, &enc.ProviderID,
			&enc.EncounterChannel, &enc.Status, &enc.ChiefComplaint, &enc.Subjective,
			&enc.Objective, &enc.Assessment, &enc.Plan, &enc.PrimaryDiagnosisCode,
			&enc.PrimaryDiagnosisName, &enc.StartedAt, &enc.CreatedAt, &enc.UpdatedAt,
			&enc.PatientName, &enc.MRN, &enc.Gender,
		)
		if err != nil {
			return nil, err
		}
		encounters = append(encounters, enc)
	}

	return encounters, nil
}

// UpdateSOAPNotes saves physician notes and diagnoses
func (r *EncounterRepository) UpdateSOAPNotes(
	ctx context.Context,
	tenantID, encounterID string,
	notes model.UpdateSOAPNotesPayload,
) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	secDiagJSON, _ := json.Marshal(notes.SecondaryDiagnoses)
	if len(notes.SecondaryDiagnoses) == 0 {
		secDiagJSON = []byte("[]")
	}

	query := `
		UPDATE encounter.encounters
		SET subjective = COALESCE($1, subjective),
		    objective = COALESCE($2, objective),
		    assessment = COALESCE($3, assessment),
		    plan = COALESCE($4, plan),
		    primary_diagnosis_code = COALESCE($5, primary_diagnosis_code),
		    primary_diagnosis_name = COALESCE($6, primary_diagnosis_name),
		    secondary_diagnoses = $7,
		    updated_at = NOW()
		WHERE tenant_id = $8 AND id = $9
	`
	_, err := db.Exec(
		ctx, query,
		notes.Subjective, notes.Objective, notes.Assessment, notes.Plan,
		notes.PrimaryDiagnosisCode, notes.PrimaryDiagnosisName,
		secDiagJSON, tenantID, encounterID,
	)
	return err
}

// AddDiagnosis records an ICD-10 diagnosis
func (r *EncounterRepository) AddDiagnosis(ctx context.Context, diag *model.Diagnosis) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	if diag.ID == "" {
		diag.ID = uuid.New().String()
	}

	query := `
		INSERT INTO encounter.diagnoses (
			id, encounter_id, patient_id, icd10_code, icd10_title,
			is_primary, clinical_status, verification_status, notes,
			diagnosed_by, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, NOW()
		)
	`
	_, err := db.Exec(
		ctx, query,
		diag.ID, diag.EncounterID, diag.PatientID, diag.ICD10Code, diag.ICD10Title,
		diag.IsPrimary, diag.ClinicalStatus, diag.VerificationStatus, diag.Notes,
		diag.DiagnosedBy,
	)
	return err
}

// ListDiagnoses returns diagnoses for an encounter
func (r *EncounterRepository) ListDiagnoses(ctx context.Context, encounterID string) ([]model.Diagnosis, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT id, encounter_id, patient_id, icd10_code, icd10_title,
		       is_primary, clinical_status, verification_status, notes, created_at
		FROM encounter.diagnoses
		WHERE encounter_id = $1
		ORDER BY is_primary DESC, created_at ASC
	`
	rows, err := db.Query(ctx, query, encounterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Diagnosis
	for rows.Next() {
		var d model.Diagnosis
		err := rows.Scan(
			&d.ID, &d.EncounterID, &d.PatientID, &d.ICD10Code, &d.ICD10Title,
			&d.IsPrimary, &d.ClinicalStatus, &d.VerificationStatus, &d.Notes, &d.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, d)
	}

	return list, nil
}

// CreatePrescription persists a complete prescription with its line items
func (r *EncounterRepository) CreatePrescription(ctx context.Context, rx *model.Prescription) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}

	if rx.ID == "" {
		rx.ID = uuid.New().String()
	}
	if rx.PrescriptionNumber == "" {
		n, _ := rand.Int(rand.Reader, big.NewInt(90000))
		rx.PrescriptionNumber = fmt.Sprintf("RX-%d-%05d", time.Now().Year(), n.Int64()+10000)
	}

	return r.server.DB.RunInTx(ctx, func(txCtx context.Context) error {
		db := r.server.DB.Conn(txCtx)

		// 1. Insert prescription header
		query := `
			INSERT INTO encounter.prescriptions (
				id, encounter_id, patient_id, prescriber_id, prescription_number,
				status, notes, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, NOW(), NOW()
			)
		`
		_, err := db.Exec(
			txCtx, query,
			rx.ID, rx.EncounterID, rx.PatientID, rx.PrescriberID, rx.PrescriptionNumber,
			rx.Status, rx.Notes,
		)
		if err != nil {
			return err
		}

		// 2. Insert items
		itemQuery := `
			INSERT INTO encounter.prescription_items (
				id, prescription_id, drug_name, dosage_form, strength,
				route, frequency, duration_days, quantity_prescribed, instructions,
				created_at
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, $7, $8, $9, $10,
				NOW()
			)
		`
		for _, item := range rx.Items {
			itemID := item.ID
			if itemID == "" {
				itemID = uuid.New().String()
			}
			_, err := db.Exec(
				txCtx, itemQuery,
				itemID, rx.ID, item.DrugName, item.DosageForm, item.Strength,
				item.Route, item.Frequency, item.DurationDays, item.QuantityPrescribed, item.Instructions,
			)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

// ListPrescriptions returns prescriptions for an encounter
func (r *EncounterRepository) ListPrescriptions(ctx context.Context, encounterID string) ([]model.Prescription, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT id, encounter_id, patient_id, prescriber_id, prescription_number,
		       status, notes, created_at, updated_at
		FROM encounter.prescriptions
		WHERE encounter_id = $1
		ORDER BY created_at DESC
	`
	rows, err := db.Query(ctx, query, encounterID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.Prescription
	for rows.Next() {
		var rx model.Prescription
		err := rows.Scan(
			&rx.ID, &rx.EncounterID, &rx.PatientID, &rx.PrescriberID, &rx.PrescriptionNumber,
			&rx.Status, &rx.Notes, &rx.CreatedAt, &rx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Load prescription items
		itemQuery := `
			SELECT id, prescription_id, drug_name, dosage_form, strength,
			       route, frequency, duration_days, quantity_prescribed, instructions, created_at
			FROM encounter.prescription_items
			WHERE prescription_id = $1
		`
		iRows, err := db.Query(ctx, itemQuery, rx.ID)
		if err == nil {
			for iRows.Next() {
				var itm model.PrescriptionItem
				_ = iRows.Scan(
					&itm.ID, &itm.PrescriptionID, &itm.DrugName, &itm.DosageForm, &itm.Strength,
					&itm.Route, &itm.Frequency, &itm.DurationDays, &itm.QuantityPrescribed, &itm.Instructions, &itm.CreatedAt,
				)
				rx.Items = append(rx.Items, itm)
			}
			iRows.Close()
		}

		list = append(list, rx)
	}

	return list, nil
}

// AutoGenerateEncounterInvoice creates billing invoice upon encounter completion
func (r *EncounterRepository) AutoGenerateEncounterInvoice(
	ctx context.Context,
	tenantID, patientID, encounterID string,
	consultationFee float64,
) (string, string, error) {
	if r.server.DB == nil {
		return "", "", errors.New("database pool not initialized")
	}

	invoiceID := uuid.New().String()
	n, _ := rand.Int(rand.Reader, big.NewInt(90000))
	invoiceNumber := fmt.Sprintf("INV-%d-%05d", time.Now().Year(), n.Int64()+10000)

	totalAmount := consultationFee
	if totalAmount <= 0 {
		totalAmount = 5000.00 // standard 5,000 NGN default consultation fee
	}

	err := r.server.DB.RunInTx(ctx, func(txCtx context.Context) error {
		db := r.server.DB.Conn(txCtx)

		// 1. Insert invoice
		invQuery := `
			INSERT INTO billing.invoices (
				id, tenant_id, patient_id, encounter_id, invoice_number,
				subtotal, discount_amount, tax_amount, total_amount, amount_paid,
				balance_due, status, created_at, updated_at
			) VALUES (
				$1, $2, $3, $4, $5,
				$6, 0.00, 0.00, $7, 0.00,
				$8, 'UNPAID', NOW(), NOW()
			)
		`
		_, err := db.Exec(
			txCtx, invQuery,
			invoiceID, tenantID, patientID, encounterID, invoiceNumber,
			totalAmount, totalAmount, totalAmount,
		)
		if err != nil {
			return err
		}

		// 2. Insert invoice line item
		itemQuery := `
			INSERT INTO billing.invoice_items (
				id, invoice_id, description, category, unit_price,
				quantity, total_price, created_at
			) VALUES (
				gen_random_uuid(), $1, 'Physician Clinical Consultation', 'CONSULTATION', $2,
				1, $3, NOW()
			)
		`
		_, err = db.Exec(txCtx, itemQuery, invoiceID, totalAmount, totalAmount)
		return err
	})

	if err != nil {
		return "", "", err
	}

	return invoiceID, invoiceNumber, nil
}

// CompleteEncounter marks an encounter as completed
func (r *EncounterRepository) CompleteEncounter(ctx context.Context, tenantID, encounterID string) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		UPDATE encounter.encounters
		SET status = 'signed', closed_at = NOW(), signed_at = NOW(), updated_at = NOW()
		WHERE tenant_id = $1 AND id = $2
	`
	_, err := db.Exec(ctx, query, tenantID, encounterID)
	return err
}

// CreateOrUpdateMilestone manages milestone stage progress in the Care Journey
func (r *EncounterRepository) CreateOrUpdateMilestone(
	ctx context.Context,
	tenantID, patientID, encounterID, stageCode, title, status, desc string,
) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	milestoneID := uuid.New().String()
	var completedAt *time.Time
	if status == "COMPLETED" {
		now := time.Now()
		completedAt = &now
	}

	query := `
		INSERT INTO orchestration.care_journey_milestones (
			id, tenant_id, patient_id, encounter_id, stage_code,
			title, description, status, completed_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, NOW()
		) ON CONFLICT DO NOTHING
	`
	_, err := db.Exec(
		ctx, query,
		milestoneID, tenantID, patientID, encounterID, stageCode,
		title, desc, status, completedAt,
	)
	return err
}
