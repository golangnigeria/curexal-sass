package domain

import (
	"context"

	"github.com/google/uuid"
)

type BranchSettingsRepository interface {
	GetByTenantID(ctx context.Context, tenantID uuid.UUID) (*BranchSettings, error)
	UpsertSection(ctx context.Context, tenantID uuid.UUID, section string, payload map[string]any) (*BranchSettings, error)
	ResetToDefaults(ctx context.Context, tenantID uuid.UUID, section *string) (*BranchSettings, error)
}
