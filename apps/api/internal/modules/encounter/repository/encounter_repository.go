package repository

import (
	"context"
	"encoding/json"
	"errors"
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

	secDiagJSON, _ := json.Marshal(enc.SecondaryDiagnoses)
	if len(enc.SecondaryDiagnoses) == 0 {
		secDiagJSON = []byte("[]")
	}

	query := `
		INSERT INTO encounter.encounters (
			id, tenant_id, care_request_id, patient_id, provider_id,
			encounter_type, mode, status, chief_complaint, subjective,
			objective, assessment, plan, primary_diagnosis_code,
			primary_diagnosis_name, secondary_diagnoses, started_at,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13, $14,
			$15, $16, NOW(),
			NOW(), NOW()
		)
	`
	_, err := db.Exec(
		ctx, query,
		enc.ID, enc.TenantID, enc.CareRequestID, enc.PatientID, enc.ProviderID,
		enc.EncounterType, enc.Mode, enc.Status, enc.ChiefComplaint, enc.Subjective,
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
		SELECT e.id, e.tenant_id, e.care_request_id, e.patient_id, e.provider_id,
		       e.encounter_type, e.mode, e.status, e.chief_complaint, e.subjective,
		       e.objective, e.assessment, e.plan, e.primary_diagnosis_code,
		       e.primary_diagnosis_name, e.secondary_diagnoses, e.started_at,
		       e.completed_at, e.created_at, e.updated_at,
		       COALESCE(p.first_name || ' ' || p.last_name, 'Patient') AS patient_name,
		       p.mrn, p.gender
		FROM encounter.encounters e
		LEFT JOIN patient.patients p ON p.id = e.patient_id
		WHERE e.tenant_id = $1 AND e.id = $2
	`
	var enc model.Encounter
	var secDiagRaw []byte
	err := db.QueryRow(ctx, query, tenantID, encounterID).Scan(
		&enc.ID, &enc.TenantID, &enc.CareRequestID, &enc.PatientID, &enc.ProviderID,
		&enc.EncounterType, &enc.Mode, &enc.Status, &enc.ChiefComplaint, &enc.Subjective,
		&enc.Objective, &enc.Assessment, &enc.Plan, &enc.PrimaryDiagnosisCode,
		&enc.PrimaryDiagnosisName, &secDiagRaw, &enc.StartedAt,
		&enc.CompletedAt, &enc.CreatedAt, &enc.UpdatedAt,
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

	return &enc, nil
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

// CompleteEncounter marks an encounter as completed
func (r *EncounterRepository) CompleteEncounter(ctx context.Context, tenantID, encounterID string) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		UPDATE encounter.encounters
		SET status = 'COMPLETED', completed_at = NOW(), updated_at = NOW()
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
		INSERT INTO encounter.care_journey_milestones (
			id, tenant_id, patient_id, encounter_id, stage_code,
			title, description, status, completed_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, NOW()
		)
	`
	_, err := db.Exec(
		ctx, query,
		milestoneID, tenantID, patientID, encounterID, stageCode,
		title, desc, status, completedAt,
	)
	return err
}
