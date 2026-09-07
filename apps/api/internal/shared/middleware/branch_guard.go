package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

// BranchAssignmentValidator validates whether a user has active facility access in a branch.
type BranchAssignmentValidator interface {
	VerifyUserFacilityAccess(ctx context.Context, orgID, branchID uuid.UUID, userID string, isOrgAdmin bool) (bool, error)
}

// PostgresBranchValidator implements BranchAssignmentValidator directly against Postgres.
type PostgresBranchValidator struct {
	pool *pgxpool.Pool
}

// NewPostgresBranchValidator creates a new PostgresBranchValidator.
func NewPostgresBranchValidator(pool *pgxpool.Pool) *PostgresBranchValidator {
	return &PostgresBranchValidator{pool: pool}
}

// VerifyUserFacilityAccess checks if a branch belongs to the org and whether the user is authorized.
func (v *PostgresBranchValidator) VerifyUserFacilityAccess(ctx context.Context, orgID, branchID uuid.UUID, userID string, isOrgAdmin bool) (bool, error) {
	if v == nil || v.pool == nil {
		return false, nil
	}

	// 1. Verify branch exists, belongs to target organization, and is ACTIVE
	branchCheckStmt := `
		SELECT EXISTS (
			SELECT 1 FROM organization.facility_branches
			WHERE id = $1 AND organization_id = $2 AND status = 'ACTIVE'
		)
	`
	var branchActive bool
	if err := v.pool.QueryRow(ctx, branchCheckStmt, branchID, orgID).Scan(&branchActive); err != nil {
		return false, fmt.Errorf("failed to verify branch status: %w", err)
	}
	if !branchActive {
		return false, nil
	}

	// 2. Organization governance administrators (owner / org_admin) have access to all branches within their organization
	if isOrgAdmin {
		return true, nil
	}

	// 3. For clinical and operational staff, verify explicit branch assignment in organization.membership_branches
	assignmentStmt := `
		SELECT EXISTS (
			SELECT 1
			FROM organization.membership_branches mb
			JOIN organization.organization_memberships om ON om.id = mb.membership_id
			WHERE om.user_id = $1::uuid
			  AND om.organization_id = $2
			  AND mb.facility_branch_id = $3
			  AND om.is_active = TRUE
		)
	`
	var isAssigned bool
	if err := v.pool.QueryRow(ctx, assignmentStmt, userID, orgID, branchID).Scan(&isAssigned); err != nil {
		return false, fmt.Errorf("failed to verify staff branch assignment: %w", err)
	}

	return isAssigned, nil
}

// RequireActiveBranch enforces Day 1 Feature #3 invariants:
// 1. Zero Implicit Fallback Policy: Returns HTTP 428 Precondition Required if active branch context is missing.
// 2. Staff Assignment Guard: Returns HTTP 403 Forbidden if non-governance staff attempts cross-branch access.
func RequireActiveBranch(validator BranchAssignmentValidator, auditLogger SecurityAuditLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			p := GetPrincipal(c)
			if p == nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
			}

			// 1. Extract branch ID from headers, principal claims, route params, or query params
			branchIDStr := p.ActiveBranchID
			if branchIDStr == "" {
				branchIDStr = c.Request().Header.Get("X-Branch-ID")
			}
			if branchIDStr == "" {
				branchIDStr = c.Request().Header.Get("X-Active-Branch-ID")
			}
			if branchIDStr == "" {
				branchIDStr = c.Param("branchId")
			}
			if branchIDStr == "" {
				branchIDStr = c.Param("branch_id")
			}
			if branchIDStr == "" {
				branchIDStr = c.QueryParam("branch_id")
			}

			// 2. Zero Implicit Fallback Policy: Never default to HQ or guess branch context
			if branchIDStr == "" {
				return echo.NewHTTPError(http.StatusPreconditionRequired, map[string]string{
					"code":    "missing_branch_context",
					"message": "Active facility branch context is required for this operation. Please select an active clinic location.",
				})
			}

			branchUUID, err := uuid.Parse(branchIDStr)
			if err != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid branch ID format: must be a valid UUID")
			}

			// 3. Resolve Organization Context
			orgIDStr := p.TenantID
			if orgIDStr == "" {
				orgIDStr = p.Organization.ActiveOrganizationID
			}
			if orgIDStr == "" {
				orgIDStr = p.OrganizationID
			}
			if orgIDStr == "" {
				orgIDStr = resolveRequestTenantID(c)
			}

			orgUUID, errOrg := uuid.Parse(orgIDStr)
			if errOrg != nil {
				return echo.NewHTTPError(http.StatusBadRequest, "Invalid organization context: must be a valid UUID")
			}

			// 4. Platform administrators bypass branch assignment restrictions
			if p.Platform.IsSuperAdmin || p.Platform.IsPlatformStaff || p.Platform.IsPlatformAdmin {
				c.Set("branch_id", branchIDStr)
				return next(c)
			}

			// 5. Check if user possesses organization governance authority
			isOrgAdmin := p.Role == "owner" || p.Role == "org_admin" ||
				p.HasPermission("organization:manage") || p.HasPermission("organization:branch:manage")

			// 6. Validate branch belonging and staff assignment
			if validator != nil {
				hasAccess, errAccess := validator.VerifyUserFacilityAccess(
					c.Request().Context(),
					orgUUID,
					branchUUID,
					p.UserID,
					isOrgAdmin,
				)

				if errAccess != nil || !hasAccess {
					if auditLogger != nil {
						go func(uID, oID, bID, uri, method, ip, ua string) {
							ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
							defer cancel()
							_ = auditLogger.LogSecurityViolation(ctx, SecurityViolationEvent{
								UserID:       uID,
								TenantID:     oID,
								BranchID:     bID,
								RequiredPerm: "workspace:branch:assignment",
								ActionURI:    uri,
								HTTPMethod:   method,
								ClientIP:     ip,
								UserAgent:    ua,
								Reason:       "Unauthorized cross-branch access attempt: staff not assigned to facility branch",
								Timestamp:    time.Now().UTC(),
							})
						}(p.UserID, orgIDStr, branchIDStr, c.Request().RequestURI, c.Request().Method, c.RealIP(), c.Request().UserAgent())
					}

					return echo.NewHTTPError(http.StatusForbidden, "You do not have an active staff assignment in this facility branch")
				}
			}

			// 7. Inject validated branch context into Echo Context
			c.Set("branch_id", branchIDStr)
			return next(c)
		}
	}
}
