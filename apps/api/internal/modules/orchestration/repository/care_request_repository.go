package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/golangnigeria/curexal/internal/modules/orchestration/model"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type CareRequestRepository struct {
	server *server.Server
}

func NewCareRequestRepository(s *server.Server) *CareRequestRepository {
	return &CareRequestRepository{server: s}
}

// CreateCareRequest inserts a new care request into the database
func (r *CareRequestRepository) CreateCareRequest(ctx context.Context, req *model.CareRequest) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	if req.ID == "" {
		req.ID = uuid.New().String()
	}

	symptomsJSON, _ := json.Marshal(req.SymptomsJSON)
	if len(req.SymptomsJSON) == 0 {
		symptomsJSON = []byte("[]")
	}

	timeWindowJSON, _ := json.Marshal(req.PreferredTimeWindow)
	if req.PreferredTimeWindow == nil {
		timeWindowJSON = []byte("{}")
	}

	query := `
		INSERT INTO orchestration.care_requests (
			id, tenant_id, patient_id, request_number, service_type,
			preferred_mode, urgency, status, chief_complaint,
			symptoms_json, preferred_time_window, assigned_care_agent_id,
			matched_provider_id, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12,
			$13, NOW(), NOW()
		)
	`

	_, err := db.Exec(
		ctx, query,
		req.ID, req.TenantID, req.PatientID, req.RequestNumber, req.ServiceType,
		req.PreferredMode, req.Urgency, req.Status, req.ChiefComplaint,
		symptomsJSON, timeWindowJSON, req.AssignedCareAgentID,
		req.MatchedProviderID,
	)
	return err
}

// GetCareRequestByID retrieves a care request with enriched patient details
func (r *CareRequestRepository) GetCareRequestByID(ctx context.Context, tenantID, requestID string) (*model.CareRequest, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT cr.id, cr.tenant_id, cr.patient_id, cr.request_number, cr.service_type,
		       cr.preferred_mode, cr.urgency, cr.status, cr.chief_complaint,
		       cr.symptoms_json, cr.preferred_time_window, cr.assigned_care_agent_id,
		       cr.matched_provider_id, cr.created_at, cr.updated_at,
		       (p.first_name || ' ' || p.last_name) AS patient_name,
		       p.mrn
		FROM orchestration.care_requests cr
		JOIN patient.patients p ON p.id = cr.patient_id
		WHERE cr.tenant_id = $1 AND cr.id = $2
	`

	var req model.CareRequest
	var symptomsRaw []byte
	var timeWindowRaw []byte

	err := db.QueryRow(ctx, query, tenantID, requestID).Scan(
		&req.ID, &req.TenantID, &req.PatientID, &req.RequestNumber, &req.ServiceType,
		&req.PreferredMode, &req.Urgency, &req.Status, &req.ChiefComplaint,
		&symptomsRaw, &timeWindowRaw, &req.AssignedCareAgentID,
		&req.MatchedProviderID, &req.CreatedAt, &req.UpdatedAt,
		&req.PatientName, &req.MRN,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	_ = json.Unmarshal(symptomsRaw, &req.SymptomsJSON)
	_ = json.Unmarshal(timeWindowRaw, &req.PreferredTimeWindow)

	return &req, nil
}

// ListCareRequests lists care requests within tenant with filters and pagination
func (r *CareRequestRepository) ListCareRequests(
	ctx context.Context,
	tenantID string,
	filter model.CareRequestFilter,
) ([]model.CareRequest, int, error) {
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

	whereClause := "WHERE cr.tenant_id = $1"
	args := []interface{}{tenantID}
	argIndex := 2

	if filter.Status != "" {
		whereClause += fmt.Sprintf(" AND cr.status = $%d", argIndex)
		args = append(args, filter.Status)
		argIndex++
	}

	if filter.Urgency != "" {
		whereClause += fmt.Sprintf(" AND cr.urgency = $%d", argIndex)
		args = append(args, filter.Urgency)
		argIndex++
	}

	if filter.PatientID != "" {
		whereClause += fmt.Sprintf(" AND cr.patient_id = $%d", argIndex)
		args = append(args, filter.PatientID)
		argIndex++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM orchestration.care_requests cr %s", whereClause)
	var total int
	if err := db.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT cr.id, cr.tenant_id, cr.patient_id, cr.request_number, cr.service_type,
		       cr.preferred_mode, cr.urgency, cr.status, cr.chief_complaint,
		       cr.symptoms_json, cr.preferred_time_window, cr.assigned_care_agent_id,
		       cr.matched_provider_id, cr.created_at, cr.updated_at,
		       (p.first_name || ' ' || p.last_name) AS patient_name,
		       p.mrn
		FROM orchestration.care_requests cr
		JOIN patient.patients p ON p.id = cr.patient_id
		%s
		ORDER BY cr.created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIndex, argIndex+1)

	args = append(args, limit, offset)
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []model.CareRequest
	for rows.Next() {
		var req model.CareRequest
		var symptomsRaw []byte
		var timeWindowRaw []byte

		err := rows.Scan(
			&req.ID, &req.TenantID, &req.PatientID, &req.RequestNumber, &req.ServiceType,
			&req.PreferredMode, &req.Urgency, &req.Status, &req.ChiefComplaint,
			&symptomsRaw, &timeWindowRaw, &req.AssignedCareAgentID,
			&req.MatchedProviderID, &req.CreatedAt, &req.UpdatedAt,
			&req.PatientName, &req.MRN,
		)
		if err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(symptomsRaw, &req.SymptomsJSON)
		_ = json.Unmarshal(timeWindowRaw, &req.PreferredTimeWindow)

		items = append(items, req)
	}

	return items, total, nil
}

// ListCareRequestsByPatient lists all requests for a specific patient
func (r *CareRequestRepository) ListCareRequestsByPatient(ctx context.Context, tenantID, patientID string) ([]model.CareRequest, error) {
	res, _, err := r.ListCareRequests(ctx, tenantID, model.CareRequestFilter{
		PatientID: patientID,
		Limit:     50,
	})
	return res, err
}

// CreateCareJourneyMilestone inserts a new Care Journey Milestone
func (r *CareRequestRepository) CreateCareJourneyMilestone(ctx context.Context, m *model.CareJourneyMilestone) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	if m.ID == "" {
		m.ID = uuid.New().String()
	}

	query := `
		INSERT INTO encounter.care_journey_milestones (
			id, tenant_id, patient_id, encounter_id, stage_code,
			title, description, status, blocking_reason,
			completed_at, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, NOW()
		)
	`
	_, err := db.Exec(
		ctx, query,
		m.ID, m.TenantID, m.PatientID, m.EncounterID, m.StageCode,
		m.Title, m.Description, m.Status, m.BlockingReason,
		m.CompletedAt,
	)
	return err
}

// GetPatientCareJourney retrieves all milestone stages for a patient in chronological order
func (r *CareRequestRepository) GetPatientCareJourney(ctx context.Context, tenantID, patientID string) ([]model.CareJourneyMilestone, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT id, tenant_id, patient_id, encounter_id, stage_code,
		       title, description, status, blocking_reason,
		       completed_at, created_at
		FROM encounter.care_journey_milestones
		WHERE tenant_id = $1 AND patient_id = $2
		ORDER BY created_at ASC
	`
	rows, err := db.Query(ctx, query, tenantID, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var milestones []model.CareJourneyMilestone
	for rows.Next() {
		var m model.CareJourneyMilestone
		err := rows.Scan(
			&m.ID, &m.TenantID, &m.PatientID, &m.EncounterID, &m.StageCode,
			&m.Title, &m.Description, &m.Status, &m.BlockingReason,
			&m.CompletedAt, &m.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		milestones = append(milestones, m)
	}

	return milestones, nil
}

// SaveTriageAssessment stores clinical vitals and updates request triage status
func (r *CareRequestRepository) SaveTriageAssessment(ctx context.Context, t *model.TriageAssessment) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	if t.ID == "" {
		t.ID = uuid.New().String()
	}

	query := `
		INSERT INTO orchestration.triage_assessments (
			id, care_request_id, patient_id, assessor_id, acuity_level,
			systolic_bp, diastolic_bp, pulse_rate, temperature,
			spo2, respiratory_rate, pain_score, triage_notes, created_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9,
			$10, $11, $12, $13, NOW()
		)
	`
	_, err := db.Exec(
		ctx, query,
		t.ID, t.CareRequestID, t.PatientID, t.AssessorID, t.AcuityLevel,
		t.SystolicBP, t.DiastolicBP, t.PulseRate, t.Temperature,
		t.SpO2, t.RespiratoryRate, t.PainScore, t.TriageNotes,
	)
	if err != nil {
		return err
	}

	// Update care request status to TRIAGED
	updateQuery := `
		UPDATE orchestration.care_requests
		SET status = 'TRIAGED', urgency = CASE WHEN $1 = 'RED' THEN 'EMERGENCY' WHEN $1 = 'YELLOW' THEN 'URGENT' ELSE urgency END, updated_at = NOW()
		WHERE id = $2
	`
	_, err = db.Exec(ctx, updateQuery, t.AcuityLevel, t.CareRequestID)
	return err
}

// GetTriageAssessmentByRequestID retrieves triage vitals for a care request
func (r *CareRequestRepository) GetTriageAssessmentByRequestID(ctx context.Context, careRequestID string) (*model.TriageAssessment, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT id, care_request_id, patient_id, assessor_id, acuity_level,
		       systolic_bp, diastolic_bp, pulse_rate, temperature,
		       spo2, respiratory_rate, pain_score, triage_notes, created_at
		FROM orchestration.triage_assessments
		WHERE care_request_id = $1
		ORDER BY created_at DESC
		LIMIT 1
	`
	var t model.TriageAssessment
	err := db.QueryRow(ctx, query, careRequestID).Scan(
		&t.ID, &t.CareRequestID, &t.PatientID, &t.AssessorID, &t.AcuityLevel,
		&t.SystolicBP, &t.DiastolicBP, &t.PulseRate, &t.Temperature,
		&t.SpO2, &t.RespiratoryRate, &t.PainScore, &t.TriageNotes, &t.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

// ListProviderProfiles queries qualified provider candidates in the tenant
func (r *CareRequestRepository) ListProviderProfiles(ctx context.Context, tenantID, specialtyCode string) ([]model.ProviderProfile, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT pp.id, pp.user_id, pp.tenant_id, pp.specialty_code, pp.sub_specialties,
		       pp.telehealth_enabled, pp.in_person_enabled, pp.max_active_queue,
		       pp.current_active_queue, pp.status, pp.consultation_languages,
		       pp.rating, pp.created_at, pp.updated_at,
		       COALESCE(u.name, 'Doctor On Duty') AS provider_name
		FROM orchestration.provider_profiles pp
		LEFT JOIN identity.users u ON u.id = pp.user_id
		WHERE pp.tenant_id = $1
		  AND ($2 = '' OR pp.specialty_code = $2 OR 'GENERAL_CONSULTATION' = $2)
		ORDER BY pp.current_active_queue ASC, pp.rating DESC
		LIMIT 20
	`
	rows, err := db.Query(ctx, query, tenantID, specialtyCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var providers []model.ProviderProfile
	for rows.Next() {
		var p model.ProviderProfile
		err := rows.Scan(
			&p.ID, &p.UserID, &p.TenantID, &p.SpecialtyCode, &p.SubSpecialties,
			&p.TelehealthEnabled, &p.InPersonEnabled, &p.MaxActiveQueue,
			&p.CurrentActiveQueue, &p.Status, &p.ConsultationLanguages,
			&p.Rating, &p.CreatedAt, &p.UpdatedAt,
			&p.ProviderName,
		)
		if err != nil {
			return nil, err
		}
		providers = append(providers, p)
	}

	return providers, nil
}

// AssignProviderToCareRequest assigns a doctor and increments their active queue count
func (r *CareRequestRepository) AssignProviderToCareRequest(ctx context.Context, tenantID, careRequestID, providerID string) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	updateReq := `
		UPDATE orchestration.care_requests
		SET matched_provider_id = $1, status = 'MATCHED', updated_at = NOW()
		WHERE tenant_id = $2 AND id = $3
	`
	_, err := db.Exec(ctx, updateReq, providerID, tenantID, careRequestID)
	if err != nil {
		return err
	}

	// Increment provider active load
	updateProvider := `
		UPDATE orchestration.provider_profiles
		SET current_active_queue = current_active_queue + 1, updated_at = NOW()
		WHERE id = $1
	`
	_, err = db.Exec(ctx, updateProvider, providerID)
	return err
}

