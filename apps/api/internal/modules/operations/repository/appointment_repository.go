package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/operations/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type AppointmentRepository struct {
	server *server.Server
}

func NewAppointmentRepository(s *server.Server) *AppointmentRepository {
	return &AppointmentRepository{server: s}
}

// CheckProviderConflict checks if doctor has overlapping active booking
func (r *AppointmentRepository) CheckProviderConflict(
	ctx context.Context,
	providerID string,
	startTime, endTime time.Time,
	excludeAppointmentID *string,
) (bool, error) {
	if r.server.DB == nil {
		return false, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT COUNT(*) 
		FROM operations.appointments
		WHERE provider_id = $1
		  AND status IN ('BOOKED', 'CHECKED_IN', 'IN_PROGRESS')
		  AND (start_time < $3 AND end_time > $2)
		  AND ($4::uuid IS NULL OR id != $4)
	`
	var count int
	var exclID interface{}
	if excludeAppointmentID != nil && *excludeAppointmentID != "" {
		exclID = *excludeAppointmentID
	}

	err := db.QueryRow(ctx, query, providerID, startTime, endTime, exclID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("error checking provider slot availability: %w", err)
	}

	return count > 0, nil
}

// CreateAppointment persists a validated appointment
func (r *AppointmentRepository) CreateAppointment(
	ctx context.Context,
	apt *model.Appointment,
) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	if apt.ID == "" {
		apt.ID = uuid.New().String()
	}
	if apt.ServiceType == "" {
		apt.ServiceType = "CONSULTATION"
	}
	if apt.DeliveryChannel == "" {
		apt.DeliveryChannel = "in_person"
	}
	if apt.Status == "" {
		apt.Status = "BOOKED"
	}

	query := `
		INSERT INTO operations.appointments (
			id, tenant_id, patient_id, provider_id, appointment_number,
			service_type, delivery_channel, status, start_time, end_time,
			reason_for_visit, cancellation_reason, virtual_meeting_url,
			created_by, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5,
			$6, $7, $8, $9, $10,
			$11, $12, $13,
			$14, NOW(), NOW()
		)
	`

	_, err := db.Exec(
		ctx, query,
		apt.ID, apt.TenantID, apt.PatientID, apt.ProviderID, apt.AppointmentNumber,
		apt.ServiceType, apt.DeliveryChannel, apt.Status, apt.StartTime, apt.EndTime,
		apt.ReasonForVisit, apt.CancellationReason, apt.VirtualMeetingURL,
		apt.CreatedBy,
	)
	return err
}

// ListAppointments retrieves appointments with optional filters
func (r *AppointmentRepository) ListAppointments(
	ctx context.Context,
	tenantID string,
	filter model.AppointmentFilter,
) ([]model.Appointment, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT a.id, a.tenant_id, a.patient_id, a.provider_id, a.appointment_number,
		       a.service_type, a.delivery_channel, a.status, a.start_time, a.end_time,
		       a.reason_for_visit, a.cancellation_reason, a.virtual_meeting_url,
		       a.created_by, a.created_at, a.updated_at,
		       (p.first_name || ' ' || p.last_name) AS patient_name,
		       p.mrn AS patient_mrn,
		       COALESCE(u.name, 'Doctor') AS provider_name,
		       pp.specialty_code,
		       pp.room_number
		FROM operations.appointments a
		JOIN patient.patients p ON p.id = a.patient_id
		JOIN orchestration.provider_profiles pp ON pp.id = a.provider_id
		LEFT JOIN identity.users u ON u.id = pp.user_id
		WHERE a.tenant_id = $1
		  AND ($2::uuid IS NULL OR a.provider_id = $2)
		  AND ($3::uuid IS NULL OR a.patient_id = $3)
		  AND ($4::text IS NULL OR a.status = $4)
		  AND ($5::text IS NULL OR a.delivery_channel = $5)
		  AND ($6::date IS NULL OR a.start_time::date = $6)
		ORDER BY a.start_time ASC
		LIMIT $7 OFFSET $8
	`

	limit := filter.Limit
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	var provID, patID, statusPtr, channelPtr, datePtr interface{}
	if filter.ProviderID != nil && *filter.ProviderID != "" {
		provID = *filter.ProviderID
	}
	if filter.PatientID != nil && *filter.PatientID != "" {
		patID = *filter.PatientID
	}
	if filter.Status != nil && *filter.Status != "" {
		statusPtr = strings.ToUpper(*filter.Status)
	}
	if filter.DeliveryChannel != nil && *filter.DeliveryChannel != "" {
		channelPtr = strings.ToLower(*filter.DeliveryChannel)
	}
	if filter.Date != nil && *filter.Date != "" {
		datePtr = *filter.Date
	}

	rows, err := db.Query(ctx, query, tenantID, provID, patID, statusPtr, channelPtr, datePtr, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query appointments: %w", err)
	}
	defer rows.Close()

	appointments := make([]model.Appointment, 0)
	for rows.Next() {
		var a model.Appointment
		err := rows.Scan(
			&a.ID, &a.TenantID, &a.PatientID, &a.ProviderID, &a.AppointmentNumber,
			&a.ServiceType, &a.DeliveryChannel, &a.Status, &a.StartTime, &a.EndTime,
			&a.ReasonForVisit, &a.CancellationReason, &a.VirtualMeetingURL,
			&a.CreatedBy, &a.CreatedAt, &a.UpdatedAt,
			&a.PatientName, &a.PatientMRN, &a.ProviderName, &a.SpecialtyCode, &a.RoomNumber,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan appointment: %w", err)
		}
		appointments = append(appointments, a)
	}

	return appointments, nil
}

// GetAppointmentByID retrieves single appointment by ID
func (r *AppointmentRepository) GetAppointmentByID(
	ctx context.Context,
	tenantID, id string,
) (*model.Appointment, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	var tenantParam interface{}
	if strings.TrimSpace(tenantID) != "" {
		tenantParam = tenantID
	}

	query := `
		SELECT a.id, a.tenant_id, a.patient_id, a.provider_id, a.appointment_number,
		       a.service_type, a.delivery_channel, a.status, a.start_time, a.end_time,
		       a.reason_for_visit, a.cancellation_reason, a.virtual_meeting_url,
		       a.created_by, a.created_at, a.updated_at,
		       (p.first_name || ' ' || p.last_name) AS patient_name,
		       p.mrn AS patient_mrn,
		       COALESCE(u.name, 'Doctor') AS provider_name,
		       pp.specialty_code,
		       pp.room_number
		FROM operations.appointments a
		JOIN patient.patients p ON p.id = a.patient_id
		JOIN orchestration.provider_profiles pp ON pp.id = a.provider_id
		LEFT JOIN identity.users u ON u.id = pp.user_id
		WHERE ($1::uuid IS NULL OR a.tenant_id = $1::uuid) AND a.id = $2::uuid
		LIMIT 1
	`

	var a model.Appointment
	err := db.QueryRow(ctx, query, tenantParam, id).Scan(
		&a.ID, &a.TenantID, &a.PatientID, &a.ProviderID, &a.AppointmentNumber,
		&a.ServiceType, &a.DeliveryChannel, &a.Status, &a.StartTime, &a.EndTime,
		&a.ReasonForVisit, &a.CancellationReason, &a.VirtualMeetingURL,
		&a.CreatedBy, &a.CreatedAt, &a.UpdatedAt,
		&a.PatientName, &a.PatientMRN, &a.ProviderName, &a.SpecialtyCode, &a.RoomNumber,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan appointment: %w", err)
	}

	return &a, nil
}

// UpdateAppointmentStatus updates appointment status
func (r *AppointmentRepository) UpdateAppointmentStatus(
	ctx context.Context,
	tenantID, id, status string,
	cancellationReason *string,
) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	validStatuses := map[string]bool{
		"BOOKED":      true,
		"CHECKED_IN":  true,
		"IN_PROGRESS": true,
		"COMPLETED":   true,
		"CANCELLED":   true,
		"NO_SHOW":     true,
	}
	cleanStatus := strings.ToUpper(strings.TrimSpace(status))
	if !validStatuses[cleanStatus] {
		return fmt.Errorf("invalid appointment status '%s'", status)
	}

	query := `
		UPDATE operations.appointments
		SET status = $1, cancellation_reason = COALESCE($2, cancellation_reason), updated_at = NOW()
		WHERE tenant_id = $3 AND id = $4
	`
	res, err := db.Exec(ctx, query, cleanStatus, cancellationReason, tenantID, id)
	if err != nil {
		return err
	}
	if res.RowsAffected() == 0 {
		return errors.New("appointment not found")
	}

	return nil
}
