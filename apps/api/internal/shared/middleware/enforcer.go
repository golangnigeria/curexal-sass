package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// MembershipValidator validates that the user is an active member of the target organization/tenant.
type MembershipValidator interface {
	IsMembershipActive(ctx context.Context, userID, tenantID string) (bool, error)
}

// SecurityAuditLogger records security violations asynchronously.
type SecurityAuditLogger interface {
	LogSecurityViolation(ctx context.Context, event SecurityViolationEvent) error
}

// SecurityViolationEvent captures context around an unauthorized or forbidden action attempt.
type SecurityViolationEvent struct {
	UserID       string    `json:"userId"`
	TenantID     string    `json:"tenantId"`
	BranchID     string    `json:"branchId,omitempty"`
	RequiredPerm string    `json:"requiredPerm"`
	ActionURI    string    `json:"actionUri"`
	HTTPMethod   string    `json:"httpMethod"`
	ClientIP     string    `json:"clientIp"`
	UserAgent    string    `json:"userAgent"`
	Reason       string    `json:"reason"`
	Timestamp    time.Time `json:"timestamp"`
}

// RBACEnforcer evaluates permissions and verifies tenant membership before allowing requests.
type RBACEnforcer struct {
	resolver    platformAuth.PermissionResolver
	validator   MembershipValidator
	auditLogger SecurityAuditLogger
}

// NewRBACEnforcer constructs a hardened RBAC enforcer.
func NewRBACEnforcer(resolver platformAuth.PermissionResolver, validator MembershipValidator, auditLogger SecurityAuditLogger) *RBACEnforcer {
	if resolver == nil {
		resolver = platformAuth.NewMemoryPermissionResolver()
	}
	return &RBACEnforcer{
		resolver:    resolver,
		validator:   validator,
		auditLogger: auditLogger,
	}
}

// RequirePermission enforces that the authenticated principal possesses the required permission.
// Hardcoded role string checks are strictly forbidden; all access decisions use granular permission codes.
func (e *RBACEnforcer) RequirePermission(permCode string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			p := GetPrincipal(c)
			if p == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
			}

			// 1. Platform Super Admin & Platform Staff have system-wide bypass
			if p.Platform.IsSuperAdmin || p.Platform.IsPlatformStaff || p.Platform.IsPlatformAdmin {
				e.injectContext(c, p)
				return next(c)
			}

			// 2. Identify Tenant & Branch context
			tenantID := p.TenantID
			if tenantID == "" {
				tenantID = p.Organization.ActiveOrganizationID
			}
			if tenantID == "" {
				tenantID = p.OrganizationID
			}
			if tenantID == "" {
				tenantID = GetActiveTenantID(c)
			}

			// 3. Verify Active Tenant Membership if tenant context is present and validator is configured
			if e.validator != nil && tenantID != "" && p.UserID != "" {
				active, err := e.validator.IsMembershipActive(c.Request().Context(), p.UserID, tenantID)
				if err != nil || !active {
					reason := "Tenant membership is inactive, suspended, or revoked"
					if err != nil {
						reason = fmt.Sprintf("Membership check failed: %v", err)
					}
					e.recordViolation(c, p, tenantID, permCode, reason)
					return echo.NewHTTPError(http.StatusForbidden, "Tenant membership is inactive or suspended")
				}
			}

			// 4. Permission Check: O(1) Fast path via principal's resolved permissions slice
			if p.HasPermission(permCode) {
				e.injectContext(c, p)
				return next(c)
			}

			// 5. Fallback check via permission resolver
			if e.resolver != nil {
				hasPerm, err := e.resolver.HasPermission(c.Request().Context(), p, permCode)
				if err == nil && hasPerm {
					e.injectContext(c, p)
					return next(c)
				}
			}

			// 6. Access Denied: Record audit violation and return 403 Forbidden
			e.recordViolation(c, p, tenantID, permCode, "Principal lacks required permission")
			return echo.NewHTTPError(http.StatusForbidden, "Insufficient permissions to perform this action")
		}
	}
}

// injectContext sets validated tenant_id, organization_id, and branch_id onto the Echo context.
func (e *RBACEnforcer) injectContext(c echo.Context, p *platformAuth.AuthenticatedPrincipal) {
	tenantID := p.TenantID
	if tenantID == "" {
		tenantID = p.Organization.ActiveOrganizationID
	}
	if tenantID == "" {
		tenantID = p.OrganizationID
	}
	if tenantID == "" {
		tenantID = GetActiveTenantID(c)
	}

	if tenantID != "" {
		c.Set("tenant_id", tenantID)
		c.Set("organization_id", tenantID)
	}

	branchID := p.ActiveBranchID
	if branchID == "" {
		branchID = c.Request().Header.Get("X-Branch-ID")
	}
	if branchID == "" {
		branchID = c.Request().Header.Get("X-Active-Branch-ID")
	}
	if branchID != "" {
		c.Set("branch_id", branchID)
	}
}

func (e *RBACEnforcer) recordViolation(c echo.Context, p *platformAuth.AuthenticatedPrincipal, tenantID, permCode, reason string) {
	branchID := p.ActiveBranchID
	if branchID == "" {
		branchID = c.Request().Header.Get("X-Branch-ID")
	}

	req := c.Request()
	event := SecurityViolationEvent{
		UserID:       p.UserID,
		TenantID:     tenantID,
		BranchID:     branchID,
		RequiredPerm: permCode,
		ActionURI:    req.RequestURI,
		HTTPMethod:   req.Method,
		ClientIP:     c.RealIP(),
		UserAgent:    req.UserAgent(),
		Reason:       reason,
		Timestamp:    time.Now().UTC(),
	}

	c.Logger().Warn(fmt.Sprintf(
		"RBAC Security Violation: user_id='%s' tenant_id='%s' required_perm='%s' reason='%s' path='%s'",
		event.UserID, event.TenantID, event.RequiredPerm, event.Reason, event.ActionURI,
	))

	if e.auditLogger != nil {
		go func(ev SecurityViolationEvent) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = e.auditLogger.LogSecurityViolation(ctx, ev)
		}(event)
	}
}

// PostgresMembershipValidator implements MembershipValidator against PostgreSQL.
type PostgresMembershipValidator struct {
	pool *pgxpool.Pool
}

// NewPostgresMembershipValidator returns a MembershipValidator using the database pool.
func NewPostgresMembershipValidator(pool *pgxpool.Pool) *PostgresMembershipValidator {
	return &PostgresMembershipValidator{pool: pool}
}

// IsMembershipActive checks if the user has an active membership in the organization or workspace.
func (v *PostgresMembershipValidator) IsMembershipActive(ctx context.Context, userID, tenantID string) (bool, error) {
	if v == nil || v.pool == nil || userID == "" || tenantID == "" {
		return false, nil
	}

	stmt := `
		SELECT is_active
		FROM organization.organization_memberships
		WHERE user_id = $1::uuid AND (organization_id = $2::uuid OR id = $2::uuid)
		LIMIT 1
	`
	var isActive bool
	err := v.pool.QueryRow(ctx, stmt, userID, tenantID).Scan(&isActive)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Fallback check workspace memberships
			stmtWs := `
				SELECT is_active
				FROM workspace.workspace_memberships
				WHERE user_id = $1::uuid AND workspace_id = $2::uuid
				LIMIT 1
			`
			errWs := v.pool.QueryRow(ctx, stmtWs, userID, tenantID).Scan(&isActive)
			if errWs == nil {
				return isActive, nil
			}
			return false, nil
		}
		return false, err
	}
	return isActive, nil
}

// DefaultSecurityAuditLogger logs security violation events to the database and Echo logger.
type DefaultSecurityAuditLogger struct {
	pool *pgxpool.Pool
}

// NewDefaultSecurityAuditLogger constructs an audit logger using Postgres.
func NewDefaultSecurityAuditLogger(pool *pgxpool.Pool) *DefaultSecurityAuditLogger {
	return &DefaultSecurityAuditLogger{pool: pool}
}

// LogSecurityViolation records an unauthorized access attempt to the audit table.
func (l *DefaultSecurityAuditLogger) LogSecurityViolation(ctx context.Context, event SecurityViolationEvent) error {
	if l == nil || l.pool == nil {
		return nil
	}

	query := `
		INSERT INTO audit.audit_events (
			event_type, user_id, organization_id, resource, action, status, metadata, created_at
		) VALUES (
			'security.access_denied',
			NULLIF($1, '')::uuid,
			NULLIF($2, '')::uuid,
			$3,
			$4,
			'FAILURE',
			jsonb_build_object(
				'required_perm', $5::text,
				'branch_id', $6::text,
				'client_ip', $7::text,
				'user_agent', $8::text,
				'reason', $9::text
			),
			$10
		)
	`
	_, err := l.pool.Exec(ctx, query,
		event.UserID,
		event.TenantID,
		event.ActionURI,
		event.HTTPMethod,
		event.RequiredPerm,
		event.BranchID,
		event.ClientIP,
		event.UserAgent,
		event.Reason,
		event.Timestamp,
	)
	if err != nil {
		return fmt.Errorf("failed to insert security audit event: %w", err)
	}
	return nil
}
