package repository

import (
	"context"
	"fmt"

	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type InviteRepository struct {
	server *server.Server
}

func NewInviteRepository(server *server.Server) *InviteRepository {
	return &InviteRepository{server: server}
}

// CreateInvitation inserts a new invitation record.
func (r *InviteRepository) CreateInvitation(ctx context.Context, inv *model.Invitation) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		INSERT INTO invitations (id, tenant_id, email, role_id, token, invited_by, status, expires_at)
		VALUES (@id, @tenant_id, @email, @role_id, @token, @invited_by, @status, @expires_at)
	`
	_, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{
		"id":         inv.ID,
		"tenant_id":  inv.TenantID,
		"email":      inv.Email,
		"role_id":    inv.RoleID,
		"token":      inv.Token,
		"invited_by": inv.InvitedBy,
		"status":     inv.Status,
		"expires_at": inv.ExpiresAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create invitation: %w", err)
	}
	return nil
}

// GetInvitationByToken retrieves an invitation by its unique token.
func (r *InviteRepository) GetInvitationByToken(ctx context.Context, token string) (*model.Invitation, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, tenant_id, email, role_id, token, invited_by, status, expires_at, accepted_at, created_at, updated_at
		FROM invitations
		WHERE token = @token
	`
	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"token": token})
	if err != nil {
		return nil, fmt.Errorf("failed to query invitation by token: %w", err)
	}

	inv, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Invitation])
	if err != nil {
		return nil, fmt.Errorf("failed to collect invitation row by token: %w", err)
	}

	return &inv, nil
}

// GetPendingInvitationByEmail checks for an existing pending invitation for this email in the tenant.
func (r *InviteRepository) GetPendingInvitationByEmail(ctx context.Context, tenantID uuid.UUID, email string) (*model.Invitation, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, tenant_id, email, role_id, token, invited_by, status, expires_at, accepted_at, created_at, updated_at
		FROM invitations
		WHERE tenant_id = @tenant_id AND email = @email AND status = 'pending'
	`
	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{
		"tenant_id": tenantID,
		"email":     email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to query pending invitation: %w", err)
	}

	inv, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[model.Invitation])
	if err != nil {
		return nil, err
	}

	return &inv, nil
}

// MarkInvitationAccepted updates the invitation status to accepted.
func (r *InviteRepository) MarkInvitationAccepted(ctx context.Context, invitationID uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE invitations
		SET status = 'accepted', accepted_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = @id
	`
	_, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{"id": invitationID})
	if err != nil {
		return fmt.Errorf("failed to mark invitation as accepted: %w", err)
	}
	return nil
}

// RevokeInvitation sets the invitation status to revoked.
func (r *InviteRepository) RevokeInvitation(ctx context.Context, invitationID uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		UPDATE invitations
		SET status = 'revoked', updated_at = CURRENT_TIMESTAMP
		WHERE id = @id AND status = 'pending'
	`
	_, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{"id": invitationID})
	if err != nil {
		return fmt.Errorf("failed to revoke invitation: %w", err)
	}
	return nil
}

// ListPendingByTenant returns all pending invitations for a tenant.
func (r *InviteRepository) ListPendingByTenant(ctx context.Context, tenantID uuid.UUID) ([]model.Invitation, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, tenant_id, email, role_id, token, invited_by, status, expires_at, accepted_at, created_at, updated_at
		FROM invitations
		WHERE tenant_id = @tenant_id AND status = 'pending'
		ORDER BY created_at DESC
	`
	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"tenant_id": tenantID})
	if err != nil {
		return nil, fmt.Errorf("failed to list pending invitations: %w", err)
	}

	invitations, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.Invitation])
	if err != nil {
		return nil, fmt.Errorf("failed to collect invitation rows: %w", err)
	}

	if invitations == nil {
		invitations = []model.Invitation{}
	}
	return invitations, nil
}

// DeleteExpiredInvitation removes or updates an expired pending invitation for a tenant+email
// so a fresh invitation can be issued via the UNIQUE(tenant_id, email) constraint.
func (r *InviteRepository) DeleteExpiredInvitation(ctx context.Context, tenantID uuid.UUID, email string) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		DELETE FROM invitations
		WHERE tenant_id = @tenant_id AND email = @email AND status IN ('pending', 'expired')
	`
	_, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{
		"tenant_id": tenantID,
		"email":     email,
	})
	if err != nil {
		return fmt.Errorf("failed to delete expired invitation: %w", err)
	}
	return nil
}
