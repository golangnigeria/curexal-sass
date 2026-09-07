package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

const (
	ResolvedOrgIDKey        = "resolved_organization_id"
	ResolvedOrgSlugKey      = "resolved_organization_slug"
	ResolvedBranchSlugKey   = "resolved_branch_slug"
	ResolvedIsCustomDomainKey = "resolved_is_custom_domain"
	ResolvedHostTypeKey     = "resolved_host_type"
)

type HostType string

const (
	HostTypePlatform     HostType = "platform"
	HostTypeOrganization HostType = "organization"
	HostTypePublic       HostType = "public"
)

type OrganizationDomainResolver struct {
	dbPool *pgxpool.Pool
}

func NewOrganizationDomainResolver(pool *pgxpool.Pool) *OrganizationDomainResolver {
	return &OrganizationDomainResolver{dbPool: pool}
}

// ExtractCleanHost extracts the hostname from X-Forwarded-Host or Host header and strips port numbers.
func ExtractCleanHost(req *http.Request) string {
	host := req.Header.Get("X-Forwarded-Host")
	if host == "" {
		host = req.Host
	}
	if colonIdx := strings.Index(host, ":"); colonIdx != -1 {
		host = host[:colonIdx]
	}
	return strings.ToLower(strings.TrimSpace(host))
}

type ResolvedDomainContext struct {
	HostType       HostType
	OrganizationID uuid.UUID
	OrgSlug        string
	IsCustomDomain bool
	Hostname       string
}

// ResolveHost resolves an incoming HTTP request hostname to an organization context.
func (r *OrganizationDomainResolver) ResolveHost(ctx context.Context, host string) (*ResolvedDomainContext, error) {
	cleanHost := strings.ToLower(strings.TrimSpace(host))

	// 1. Platform Console Check
	if cleanHost == "app.curexal.space" || cleanHost == "app.curexal.internal" || cleanHost == "app.localhost" {
		return &ResolvedDomainContext{
			HostType:       HostTypePlatform,
			Hostname:       cleanHost,
			IsCustomDomain: false,
		}, nil
	}

	// 2. Public / Root Platform Check
	if cleanHost == "curexal.space" || cleanHost == "www.curexal.space" || cleanHost == "curexal.com" || cleanHost == "localhost" || cleanHost == "127.0.0.1" {
		return &ResolvedDomainContext{
			HostType:       HostTypePublic,
			Hostname:       cleanHost,
			IsCustomDomain: false,
		}, nil
	}

	// 3. Localhost Subdomain Check (e.g. everight.localhost)
	if strings.HasSuffix(cleanHost, ".localhost") {
		subdomain := strings.TrimSuffix(cleanHost, ".localhost")
		if subdomain != "app" && subdomain != "api" && subdomain != "public" && r.dbPool != nil {
			var orgID uuid.UUID
			var slug string
			err := r.dbPool.QueryRow(ctx, `
				SELECT id, slug FROM organization.organizations WHERE slug = $1 LIMIT 1
			`, subdomain).Scan(&orgID, &slug)
			if err == nil {
				return &ResolvedDomainContext{
					HostType:       HostTypeOrganization,
					OrganizationID: orgID,
					OrgSlug:        slug,
					IsCustomDomain: false,
					Hostname:       cleanHost,
				}, nil
			}
		}
	}

	// 4. Curexal Subdomain Check (e.g. everight.curexal.space)
	if strings.HasSuffix(cleanHost, ".curexal.space") || strings.HasSuffix(cleanHost, ".curexal.internal") {
		subdomain := strings.TrimSuffix(cleanHost, ".curexal.space")
		subdomain = strings.TrimSuffix(subdomain, ".curexal.internal")
		if subdomain != "app" && subdomain != "api" && subdomain != "public" && subdomain != "admin" {
			if r.dbPool != nil {
				var orgID uuid.UUID
				var slug string
				// Check organization_domains first or organizations table
				err := r.dbPool.QueryRow(ctx, `
					SELECT o.id, o.slug 
					FROM organization.organizations o
					LEFT JOIN organization.organization_domains od ON od.organization_id = o.id
					WHERE od.hostname = $1 OR o.slug = $2
					LIMIT 1
				`, cleanHost, subdomain).Scan(&orgID, &slug)
				if err == nil {
					return &ResolvedDomainContext{
						HostType:       HostTypeOrganization,
						OrganizationID: orgID,
						OrgSlug:        slug,
						IsCustomDomain: false,
						Hostname:       cleanHost,
					}, nil
				}
			}
		}
	}

	// 5. Custom Domain Check (e.g. everightdiagnostics.com)
	if r.dbPool != nil {
		var orgID uuid.UUID
		var slug string
		err := r.dbPool.QueryRow(ctx, `
			SELECT o.id, o.slug 
			FROM organization.organization_domains od
			JOIN organization.organizations o ON o.id = od.organization_id
			WHERE od.hostname = $1 AND od.is_verified = TRUE AND od.status = 'ACTIVE'
			LIMIT 1
		`, cleanHost).Scan(&orgID, &slug)
		if err == nil {
			return &ResolvedDomainContext{
				HostType:       HostTypeOrganization,
				OrganizationID: orgID,
				OrgSlug:        slug,
				IsCustomDomain: true,
				Hostname:       cleanHost,
			}, nil
		}
	}

	return &ResolvedDomainContext{
		HostType:       HostTypePublic,
		Hostname:       cleanHost,
		IsCustomDomain: false,
	}, nil
}

// DomainResolverMiddleware attaches the resolved organization and domain context to each request.
func DomainResolverMiddleware(resolver *OrganizationDomainResolver) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if resolver != nil {
				cleanHost := ExtractCleanHost(c.Request())
				resCtx, err := resolver.ResolveHost(c.Request().Context(), cleanHost)
				if err == nil && resCtx != nil {
					c.Set(ResolvedHostTypeKey, string(resCtx.HostType))
					if resCtx.OrganizationID != uuid.Nil {
						c.Set(ResolvedOrgIDKey, resCtx.OrganizationID.String())
						c.Set(ResolvedOrgSlugKey, resCtx.OrgSlug)
						c.Set(ResolvedIsCustomDomainKey, resCtx.IsCustomDomain)
					}
				}
			}
			return next(c)
		}
	}
}

// GetResolvedOrgID retrieves the resolved organization ID from context.
func GetResolvedOrgID(c echo.Context) string {
	if val := c.Get(ResolvedOrgIDKey); val != nil {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}

// GetResolvedOrgSlug retrieves the resolved organization slug from context.
func GetResolvedOrgSlug(c echo.Context) string {
	if val := c.Get(ResolvedOrgSlugKey); val != nil {
		if s, ok := val.(string); ok {
			return s
		}
	}
	return ""
}
