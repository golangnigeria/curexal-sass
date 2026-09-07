package authz

import (
	"context"
	"fmt"
	"strings"
)

// Service provides high-level centralized authorization evaluations.
type Service struct {
	casbin *CasbinEngine
}

func NewService(engine *CasbinEngine) *Service {
	return &Service{
		casbin: engine,
	}
}

// Authorize evaluates a comprehensive authorization request against the principal context.
func (s *Service) Authorize(ctx context.Context, principal *Principal, req AuthzRequest) error {
	if principal == nil {
		return ErrUnauthorized
	}

	// 1. Platform Staff Boundary
	if req.Scope.Scope == ScopePlatform {
		if !principal.IsPlatformStaff {
			return ErrForbidden
		}
		return nil
	}

	// 2. Organization Boundary
	if req.Scope.Scope == ScopeOrganization {
		if principal.Context.Scope != ScopeOrganization && !principal.IsPlatformStaff {
			hasOrgRole := false
			for _, r := range principal.Roles {
				lr := strings.ToLower(r)
				if lr == "owner" || lr == "org_admin" || lr == "org_regional_manager" || lr == "org_quality_manager" || lr == "org_finance_manager" || lr == "org_hr_manager" {
					hasOrgRole = true
					break
				}
			}
			if !hasOrgRole {
				return ErrForbidden
			}
		}
	}

	// 3. Professional Credential Verification
	if req.Credential != "" && !principal.HasCredential(req.Credential) {
		return ErrCredentialRequired
	}

	// 4. Permission / Casbin Evaluation
	if req.Permission != "" {
		if !principal.HasPermission(req.Permission) {
			// Also evaluate via Casbin roles
			allowed := false
			dom := "*"
			if req.Scope.OrganizationID != nil {
				dom = req.Scope.OrganizationID.String()
			}
			for _, role := range principal.Roles {
				ok, _ := s.casbin.Enforce(role, dom, req.Resource, req.Action)
				if ok {
					allowed = true
					break
				}
			}
			if !allowed {
				return fmt.Errorf("%w: required permission %s", ErrForbidden, req.Permission)
			}
		}
	}

	return nil
}

// Can checks if the principal is authorized for a specific permission string within a scope.
func (s *Service) Can(ctx context.Context, principal *Principal, permission string, scope ScopeContext) bool {
	err := s.Authorize(ctx, principal, AuthzRequest{
		Permission: permission,
		Scope:      scope,
	})
	return err == nil
}
