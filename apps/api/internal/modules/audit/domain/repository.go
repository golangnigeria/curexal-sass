package domain

import (
	"context"

	"github.com/google/uuid"
)

type AuditRepository interface {
	GetByID(ctx context.Context, id uuid.UUID) (*AuditLog, error)
	Create(ctx context.Context, payload *CreateAuditLogPayload) (*AuditLog, error)
	ListTenantLogs(ctx context.Context, tenantID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]AuditLog, error)
	ListPlatformLogs(ctx context.Context, orgID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]AuditLog, error)
	GetStats(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID) (*AdminStats, error)
	ListAll(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID, limit int, offset int) ([]AuditLog, error)
}
