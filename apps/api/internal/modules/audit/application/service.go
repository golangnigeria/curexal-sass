package application

import (
	"context"

	"github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
)

type AuditApplicationService struct {
	server    *server.Server
	auditRepo domain.AuditRepository
}

func NewAuditApplicationService(server *server.Server, auditRepo domain.AuditRepository) *AuditApplicationService {
	return &AuditApplicationService{
		server:    server,
		auditRepo: auditRepo,
	}
}

func (s *AuditApplicationService) LogEvent(
	ctx context.Context,
	payload *domain.CreateAuditLogPayload,
) (*domain.AuditLog, error) {
	return s.auditRepo.Create(ctx, payload)
}

func (s *AuditApplicationService) GetLogByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error) {
	return s.auditRepo.GetByID(ctx, id)
}

func (s *AuditApplicationService) ListTenantLogs(ctx context.Context, tenantID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]domain.AuditLog, error) {
	return s.auditRepo.ListTenantLogs(ctx, tenantID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search, limit, offset)
}

func (s *AuditApplicationService) ListPlatformLogs(ctx context.Context, orgID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]domain.AuditLog, error) {
	return s.auditRepo.ListPlatformLogs(ctx, orgID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search, limit, offset)
}

func (s *AuditApplicationService) GetAdminStats(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID) (*domain.AdminStats, error) {
	return s.auditRepo.GetStats(ctx, tenantID, orgID)
}
