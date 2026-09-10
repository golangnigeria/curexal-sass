package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ProviderProfileRepository struct {
	server *server.Server
}

func NewProviderProfileRepository(s *server.Server) *ProviderProfileRepository {
	return &ProviderProfileRepository{server: s}
}

// ListProviderProfiles retrieves providers in a tenant with optional status and specialty filters
func (r *ProviderProfileRepository) ListProviderProfiles(
	ctx context.Context,
	tenantID string,
	status *string,
	specialty *string,
) ([]model.ProviderProfile, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT pp.id, pp.user_id, pp.tenant_id, 
		       COALESCE(pp.license_number, ''), COALESCE(pp.license_issuer, 'MDCN'), pp.license_verified_at,
		       pp.specialty_code, pp.sub_specialties,
		       pp.room_number, pp.room_name, pp.virtual_room_url,
		       pp.telehealth_enabled, pp.in_person_enabled, pp.max_active_queue,
		       pp.current_active_queue, pp.status, pp.consultation_languages,
		       pp.rating, pp.created_at, pp.updated_at,
		       COALESCE(u.name, 'Doctor On Duty') AS provider_name,
		       u.email, u.phone
		FROM orchestration.provider_profiles pp
		LEFT JOIN identity.users u ON u.id = pp.user_id
		WHERE pp.tenant_id = $1
		  AND ($2::text IS NULL OR pp.status = $2)
		  AND ($3::text IS NULL OR pp.specialty_code = $3)
		ORDER BY pp.status ASC, pp.current_active_queue ASC, pp.rating DESC
	`

	var statusFilter, specialtyFilter *string
	if status != nil && strings.TrimSpace(*status) != "" {
		s := strings.ToUpper(strings.TrimSpace(*status))
		statusFilter = &s
	}
	if specialty != nil && strings.TrimSpace(*specialty) != "" {
		sp := strings.ToUpper(strings.TrimSpace(*specialty))
		specialtyFilter = &sp
	}

	rows, err := db.Query(ctx, query, tenantID, statusFilter, specialtyFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to query provider profiles: %w", err)
	}
	defer rows.Close()

	providers := make([]model.ProviderProfile, 0)
	for rows.Next() {
		var p model.ProviderProfile
		err := rows.Scan(
			&p.ID, &p.UserID, &p.TenantID,
			&p.LicenseNumber, &p.LicenseIssuer, &p.LicenseVerifiedAt,
			&p.SpecialtyCode, &p.SubSpecialties,
			&p.RoomNumber, &p.RoomName, &p.VirtualRoomURL,
			&p.TelehealthEnabled, &p.InPersonEnabled, &p.MaxActiveQueue,
			&p.CurrentActiveQueue, &p.Status, &p.ConsultationLanguages,
			&p.Rating, &p.CreatedAt, &p.UpdatedAt,
			&p.ProviderName, &p.Email, &p.Phone,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan provider profile: %w", err)
		}
		providers = append(providers, p)
	}

	return providers, nil
}

// GetProviderProfileByID fetches single provider by profile ID
func (r *ProviderProfileRepository) GetProviderProfileByID(
	ctx context.Context,
	tenantID, id string,
) (*model.ProviderProfile, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	query := `
		SELECT pp.id, pp.user_id, pp.tenant_id, 
		       COALESCE(pp.license_number, ''), COALESCE(pp.license_issuer, 'MDCN'), pp.license_verified_at,
		       pp.specialty_code, pp.sub_specialties,
		       pp.room_number, pp.room_name, pp.virtual_room_url,
		       pp.telehealth_enabled, pp.in_person_enabled, pp.max_active_queue,
		       pp.current_active_queue, pp.status, pp.consultation_languages,
		       pp.rating, pp.created_at, pp.updated_at,
		       COALESCE(u.name, 'Doctor On Duty') AS provider_name,
		       u.email, u.phone
		FROM orchestration.provider_profiles pp
		LEFT JOIN identity.users u ON u.id = pp.user_id
		WHERE pp.tenant_id = $1 AND pp.id = $2
		LIMIT 1
	`

	var p model.ProviderProfile
	err := db.QueryRow(ctx, query, tenantID, id).Scan(
		&p.ID, &p.UserID, &p.TenantID,
		&p.LicenseNumber, &p.LicenseIssuer, &p.LicenseVerifiedAt,
		&p.SpecialtyCode, &p.SubSpecialties,
		&p.RoomNumber, &p.RoomName, &p.VirtualRoomURL,
		&p.TelehealthEnabled, &p.InPersonEnabled, &p.MaxActiveQueue,
		&p.CurrentActiveQueue, &p.Status, &p.ConsultationLanguages,
		&p.Rating, &p.CreatedAt, &p.UpdatedAt,
		&p.ProviderName, &p.Email, &p.Phone,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan provider profile: %w", err)
	}

	return &p, nil
}

// CreateProviderProfile inserts or upserts provider profile
func (r *ProviderProfileRepository) CreateProviderProfile(
	ctx context.Context,
	p *model.ProviderProfile,
) error {
	if r.server.DB == nil {
		return errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	if p.LicenseIssuer == "" {
		p.LicenseIssuer = "MDCN"
	}
	if p.Status == "" {
		p.Status = "ON_DUTY"
	}
	if p.MaxActiveQueue <= 0 {
		p.MaxActiveQueue = 10
	}
	if p.Rating <= 0 {
		p.Rating = 5.0
	}
	if len(p.ConsultationLanguages) == 0 {
		p.ConsultationLanguages = []string{"English"}
	}

	query := `
		INSERT INTO orchestration.provider_profiles (
			id, user_id, tenant_id, license_number, license_issuer, license_verified_at,
			specialty_code, sub_specialties, room_number, room_name, virtual_room_url,
			telehealth_enabled, in_person_enabled, max_active_queue, current_active_queue,
			status, consultation_languages, rating, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, NOW(), NOW()
		)
		ON CONFLICT (tenant_id, user_id) DO UPDATE SET
			license_number = EXCLUDED.license_number,
			specialty_code = EXCLUDED.specialty_code,
			sub_specialties = EXCLUDED.sub_specialties,
			room_number = COALESCE(EXCLUDED.room_number, orchestration.provider_profiles.room_number),
			room_name = COALESCE(EXCLUDED.room_name, orchestration.provider_profiles.room_name),
			virtual_room_url = COALESCE(EXCLUDED.virtual_room_url, orchestration.provider_profiles.virtual_room_url),
			telehealth_enabled = EXCLUDED.telehealth_enabled,
			in_person_enabled = EXCLUDED.in_person_enabled,
			max_active_queue = EXCLUDED.max_active_queue,
			status = EXCLUDED.status,
			consultation_languages = EXCLUDED.consultation_languages,
			updated_at = NOW()
	`

	_, err := db.Exec(ctx, query,
		p.ID, p.UserID, p.TenantID, p.LicenseNumber, p.LicenseIssuer, p.LicenseVerifiedAt,
		p.SpecialtyCode, p.SubSpecialties, p.RoomNumber, p.RoomName, p.VirtualRoomURL,
		p.TelehealthEnabled, p.InPersonEnabled, p.MaxActiveQueue, p.CurrentActiveQueue,
		p.Status, p.ConsultationLanguages, p.Rating,
	)
	return err
}

// UpdateProviderStatus updates duty status (ON_DUTY, ON_BREAK, OFF_DUTY, BUSY)
func (r *ProviderProfileRepository) UpdateProviderStatus(
	ctx context.Context,
	tenantID, id, status string,
) (*time.Time, error) {
	if r.server.DB == nil {
		return nil, errors.New("database pool not initialized")
	}
	db := r.server.DB.Conn(ctx)

	validStatuses := map[string]bool{
		"ON_DUTY":  true,
		"ON_BREAK": true,
		"OFF_DUTY": true,
		"BUSY":     true,
	}
	cleanStatus := strings.ToUpper(strings.TrimSpace(status))
	if !validStatuses[cleanStatus] {
		return nil, fmt.Errorf("invalid provider status '%s'; must be ON_DUTY, ON_BREAK, OFF_DUTY, or BUSY", status)
	}

	var updatedAt time.Time
	query := `
		UPDATE orchestration.provider_profiles
		SET status = $1, updated_at = NOW()
		WHERE tenant_id = $2 AND id = $3
		RETURNING updated_at
	`
	err := db.QueryRow(ctx, query, cleanStatus, tenantID, id).Scan(&updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("provider profile not found")
		}
		return nil, err
	}

	return &updatedAt, nil
}
