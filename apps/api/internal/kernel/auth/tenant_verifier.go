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

	var roleTitle string
	stmt := `
		SELECT role_title
		FROM organization.organization_memberships
		WHERE user_id = $1::uuid AND (organization_id = $2::uuid OR id = $2::uuid) AND is_active = TRUE
		LIMIT 1
	`

	err := v.pool.QueryRow(ctx, stmt, userID, tenantID).Scan(&roleTitle)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Also check workspace memberships if organization lookup misses
			stmtWorkspace := `
				SELECT role_title
				FROM workspace.workspace_memberships
				WHERE user_id = $1::uuid AND workspace_id = $2::uuid AND is_active = TRUE
				LIMIT 1
			`
			errWs := v.pool.QueryRow(ctx, stmtWorkspace, userID, tenantID).Scan(&roleTitle)
			if errWs == nil {
				return true, roleTitle, nil
			}
			return false, "", nil
		}
		return false, "", fmt.Errorf("membership verification database error: %w", err)
	}

	return true, roleTitle, nil
}
