package domain

import (
	"context"
	"errors"
	"time"
)

var (
	ErrUnauthorizedScope = errors.New("unauthorized navigation context scope")
	ErrInvalidScope      = errors.New("invalid navigation context scope")
)

type NavigationStatus string

const (
	NavigationStatusActive   NavigationStatus = "active"
	NavigationStatusPending  NavigationStatus = "pending"
	NavigationStatusDisabled NavigationStatus = "disabled"
)

// NavigationItem represents a canonical database-driven navigation item
type NavigationItem struct {
	ID                 string           `json:"id"`
	Key                string           `json:"key"`
	ContextScope       string           `json:"contextScope"`
	ModuleCode         *string          `json:"moduleCode,omitempty"`
	Title              string           `json:"title"`
	Description        *string          `json:"description,omitempty"`
	Icon               string           `json:"icon"`
	Path               string           `json:"path"`
	Order              int              `json:"order"`
	ParentID           *string          `json:"parentId,omitempty"`
	RequiredPermission *string          `json:"requiredPermission,omitempty"`
	RequiredCapability *string          `json:"requiredCapability,omitempty"`
	Status             NavigationStatus `json:"status"`
	IsVisible          bool             `json:"isVisible"`
	IsActive           bool             `json:"isActive"`
	BadgeKey           *string          `json:"badgeKey,omitempty"`
	Children           []NavigationItem `json:"children,omitempty"`
	CreatedAt          time.Time        `json:"createdAt"`
}

// NavigationContext specifies the resolved scope of a navigation request
type NavigationContext struct {
	Type             string   `json:"type"` // "platform" | "organization" | "workspace" | "patient"
	OrganizationID   *string  `json:"organizationId,omitempty"`
	OrganizationSlug *string  `json:"organizationSlug,omitempty"`
	BranchID         *string  `json:"branchId,omitempty"`
	BranchSlug       *string  `json:"branchSlug,omitempty"`
	BranchType       *string  `json:"branchType,omitempty"`
	Role             *string  `json:"role,omitempty"`
	Capabilities     []string `json:"capabilities,omitempty"`
}

// NavigationResponse is the canonical API response contract
type NavigationResponse struct {
	Context NavigationContext `json:"context"`
	Items   []NavigationItem  `json:"items"`
}

// NavigationRepository defines the database access contract for navigation items
type NavigationRepository interface {
	GetNavigationItemsByScope(ctx context.Context, scope string, enabledModules []string, userPermissions []string, isSuperAdminOrOwner bool) ([]NavigationItem, error)
	GetAllActiveNavigationItems(ctx context.Context) ([]NavigationItem, error)
}
