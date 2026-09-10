package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/model"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/repository"
	"github.com/google/uuid"
)

type ProviderProfileService struct {
	server       *server.Server
	providerRepo *repository.ProviderProfileRepository
}

func NewProviderProfileService(s *server.Server, providerRepo *repository.ProviderProfileRepository) *ProviderProfileService {
	return &ProviderProfileService{
		server:       s,
		providerRepo: providerRepo,
	}
}

// ListProviders retrieves provider directory for a tenant
func (s *ProviderProfileService) ListProviders(
	ctx context.Context,
	tenantID string,
	status *string,
	specialty *string,
) ([]model.ProviderProfile, error) {
	if tenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	return s.providerRepo.ListProviderProfiles(ctx, tenantID, status, specialty)
}

// GetProviderByID retrieves single provider profile
func (s *ProviderProfileService) GetProviderByID(
	ctx context.Context,
	tenantID, id string,
) (*model.ProviderProfile, error) {
	if tenantID == "" || id == "" {
		return nil, errors.New("tenantId and id are required")
	}
	return s.providerRepo.GetProviderProfileByID(ctx, tenantID, id)
}

// CreateOrLinkProviderProfile registers or updates clinical credentials for a user
func (s *ProviderProfileService) CreateOrLinkProviderProfile(
	ctx context.Context,
	tenantID string,
	req model.CreateProviderProfileRequest,
) (*model.ProviderProfile, error) {
	if tenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	if strings.TrimSpace(req.UserID) == "" {
		return nil, errors.New("userId is required")
	}
	if strings.TrimSpace(req.LicenseNumber) == "" {
		return nil, errors.New("licenseNumber is required")
	}
	if strings.TrimSpace(req.SpecialtyCode) == "" {
		return nil, errors.New("specialtyCode is required")
	}

	telehealth := true
	if req.TelehealthEnabled != nil {
		telehealth = *req.TelehealthEnabled
	}
	inPerson := true
	if req.InPersonEnabled != nil {
		inPerson = *req.InPersonEnabled
	}
	maxQueue := 10
	if req.MaxActiveQueue != nil && *req.MaxActiveQueue > 0 {
		maxQueue = *req.MaxActiveQueue
	}
	status := "ON_DUTY"
	if strings.TrimSpace(req.Status) != "" {
		status = strings.ToUpper(strings.TrimSpace(req.Status))
	}
	issuer := "MDCN"
	if strings.TrimSpace(req.LicenseIssuer) != "" {
		issuer = strings.TrimSpace(req.LicenseIssuer)
	}

	now := time.Now()
	id := uuid.New().String()
	virtualURL := req.VirtualRoomURL
	if virtualURL == nil || *virtualURL == "" {
		defaultURL := fmt.Sprintf("https://telehealth.curexal.com/room/prov_%s", id[:8])
		virtualURL = &defaultURL
	}

	profile := &model.ProviderProfile{
		ID:                    id,
		UserID:                strings.TrimSpace(req.UserID),
		TenantID:              tenantID,
		LicenseNumber:         strings.TrimSpace(req.LicenseNumber),
		LicenseIssuer:         issuer,
		LicenseVerifiedAt:     &now,
		SpecialtyCode:         strings.ToUpper(strings.TrimSpace(req.SpecialtyCode)),
		SubSpecialties:        req.SubSpecialties,
		RoomNumber:            req.RoomNumber,
		RoomName:              req.RoomName,
		VirtualRoomURL:        virtualURL,
		TelehealthEnabled:     telehealth,
		InPersonEnabled:       inPerson,
		MaxActiveQueue:        maxQueue,
		CurrentActiveQueue:    0,
		Status:                status,
		ConsultationLanguages: req.ConsultationLanguages,
		Rating:                5.0,
		CreatedAt:             now,
		UpdatedAt:             now,
	}

	if err := s.providerRepo.CreateProviderProfile(ctx, profile); err != nil {
		return nil, err
	}

	return s.providerRepo.GetProviderProfileByID(ctx, tenantID, profile.ID)
}

// UpdateProviderStatus modifies duty status
func (s *ProviderProfileService) UpdateProviderStatus(
	ctx context.Context,
	tenantID, id string,
	req model.UpdateProviderStatusRequest,
) (*time.Time, error) {
	if tenantID == "" || id == "" {
		return nil, errors.New("tenantId and id are required")
	}
	return s.providerRepo.UpdateProviderStatus(ctx, tenantID, id, req.Status)
}
