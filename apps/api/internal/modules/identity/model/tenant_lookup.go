package model

import (
	"context"

	"github.com/google/uuid"
)

type TenantInfo struct {
	ID   uuid.UUID
	Name string
	Slug string
}

type TenantLookup interface {
	GetTenantByID(ctx context.Context, id uuid.UUID) (*TenantInfo, error)
}
