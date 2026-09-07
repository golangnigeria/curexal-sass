package auth

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// TenantMembershipVerifier checks database membership for authenticated users before granting tenant context.
type TenantMembershipVerifier interface {
	VerifyMembership(ctx context.Context, userID, tenantID string) (bool, string, error)
}

type PostgresTenantVerifier struct {
	pool *pgxpool.Pool
}

func NewPostgresTenantVerifier(pool *pgxpool.Pool) *PostgresTenantVerifier {
	return &PostgresTenantVerifier{pool: pool}
}

func (v *PostgresTenantVerifier) VerifyMembership(ctx context.Context, userID, tenantID string) (bool, string, error) {
	if v == nil || v.pool == nil || userID == "" || tenantID == "" {
		return false, "", nil
	}

	var role string
	stmt := `
		SELECT COALESCE(NULLIF(m.role, ''), NULLIF(m.role_title, ''), 'member')
		FROM organization.organization_memberships m
		WHERE m.user_id = $1::uuid 
		  AND (
		      m.organization_id = $2::uuid 
		      OR EXISTS (
		          SELECT 1 FROM organization.facility_branches fb 
		          WHERE fb.id = $2::uuid AND fb.organization_id = m.organization_id
		      )
		  ) 
		  AND m.is_active = TRUE
		LIMIT 1
	`

	err := v.pool.QueryRow(ctx, stmt, userID, tenantID).Scan(&role)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, "", nil
		}
		return false, "", fmt.Errorf("membership verification database error: %w", err)
	}

	return true, role, nil
}
