package application

import (
	"context"

	"github.com/golangnigeria/curexal/internal/modules/settings/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
)

type SettingsApplicationService struct {
	server *server.Server
	repo   domain.BranchSettingsRepository
}

func NewSettingsApplicationService(server *server.Server, repo domain.BranchSettingsRepository) *SettingsApplicationService {
	return &SettingsApplicationService{
		server: server,
		repo:   repo,
	}
}

func (s *SettingsApplicationService) GetByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.BranchSettings, error) {
	return s.repo.GetByTenantID(ctx, tenantID)
}

func (s *SettingsApplicationService) UpdateSection(ctx context.Context, tenantID uuid.UUID, section string, payload map[string]any) (*domain.BranchSettings, error) {
	return s.repo.UpsertSection(ctx, tenantID, section, payload)
}

func (s *SettingsApplicationService) ResetToDefaults(ctx context.Context, tenantID uuid.UUID, section *string) (*domain.BranchSettings, error) {
	return s.repo.ResetToDefaults(ctx, tenantID, section)
}
