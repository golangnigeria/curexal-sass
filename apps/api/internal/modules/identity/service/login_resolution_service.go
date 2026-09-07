package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/identity/domain"
	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/identity/repository"
	"github.com/redis/go-redis/v9"
)

type ResolutionStatus string

const (
	StatusAuthenticated           ResolutionStatus = "authenticated"
	StatusRedirectRequired        ResolutionStatus = "redirect_required"
	StatusBranchSelectionRequired ResolutionStatus = "branch_selection_required"
	StatusOrgSelectionRequired    ResolutionStatus = "organization_selection_required"
	StatusUnassignedFacilityBranch ResolutionStatus = "unassigned_facility_branch"
	StatusOrgAccessRequired       ResolutionStatus = "organization_access_required"
)

type OrganizationOption struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Role         string `json:"role"`
	CustomDomain string `json:"customDomain,omitempty"`
}

type LoginResolution struct {
	Status           ResolutionStatus     `json:"status"`
	DestinationPath  string               `json:"destinationPath,omitempty"`
	TargetURL        string               `json:"targetUrl,omitempty"`
	TargetHost       string               `json:"targetHost,omitempty"`
	Organization     *OrganizationOption  `json:"organization,omitempty"`
	Organizations    []OrganizationOption `json:"organizations,omitempty"`
	ActiveBranch     *BranchSummary       `json:"activeBranch,omitempty"`
	AssignedBranches []BranchSummary     `json:"assignedBranches,omitempty"`
	SelectionToken   string               `json:"selectionToken,omitempty"`
	Message          string               `json:"message,omitempty"`
}

type ExchangePayload struct {
	UserID          string `json:"userId"`
	OrganizationID  string `json:"organizationId,omitempty"`
	BranchID        string `json:"branchId,omitempty"`
	TenantSlug      string `json:"tenantSlug,omitempty"`
	Role            string `json:"role,omitempty"`
	TargetHost      string `json:"targetHost"`
	DestinationPath string `json:"destinationPath"`
	CreatedAt       int64  `json:"createdAt"`
}

type LoginResolutionService struct {
	server      *server.Server
	authService *AuthService
	userRepo    *repository.UserRepository
}

func NewLoginResolutionService(s *server.Server, authService *AuthService, userRepo *repository.UserRepository) *LoginResolutionService {
	return &LoginResolutionService{
		server:      s,
		authService: authService,
		userRepo:    userRepo,
	}
}

// GenerateExchangeToken creates a 32-byte cryptographically secure random token string.
func GenerateExchangeToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// HashToken computes a SHA-256 hex digest of the exchange token.
func HashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// CreateExchangeToken hashes and stores the exchange payload in Redis with a 45-second TTL.
func (s *LoginResolutionService) CreateExchangeToken(ctx context.Context, payload *ExchangePayload) (string, error) {
	rawToken, err := GenerateExchangeToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate exchange token: %w", err)
	}

	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal exchange payload: %w", err)
	}

	hash := HashToken(rawToken)
	key := fmt.Sprintf("auth:exchange:%s", hash)
	if s.server.Redis != nil {
		if err := s.server.Redis.Set(ctx, key, payloadJSON, 45*time.Second).Err(); err != nil {
			return "", fmt.Errorf("failed to store exchange token in redis: %w", err)
		}
	} else if s.server.Cache != nil {
		if err := s.server.Cache.Set(ctx, key, string(payloadJSON), 45*time.Second); err != nil {
			return "", fmt.Errorf("failed to store exchange token in cache: %w", err)
		}
	} else {
		return "", fmt.Errorf("no cache store available for exchange token")
	}

	return rawToken, nil
}

// RedeemExchangeToken atomically fetches and deletes the exchange token (GETDEL) ensuring single-use.
func (s *LoginResolutionService) RedeemExchangeToken(ctx context.Context, rawToken string) (*ExchangePayload, error) {
	if strings.TrimSpace(rawToken) == "" {
		return nil, fmt.Errorf("exchange token is empty")
	}

	hash := HashToken(rawToken)
	key := fmt.Sprintf("auth:exchange:%s", hash)

	var payloadJSON string
	if s.server.Redis != nil {
		val, err := s.server.Redis.GetDel(ctx, key).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				return nil, fmt.Errorf("exchange token expired or already used")
			}
			return nil, fmt.Errorf("redis error redeeming token: %w", err)
		}
		payloadJSON = val
	} else if s.server.Cache != nil {
		val, found := s.server.Cache.GetString(ctx, key)
		if !found {
			return nil, fmt.Errorf("exchange token expired or already used")
		}
		_ = s.server.Cache.Delete(ctx, key)
		payloadJSON = val
	} else {
		return nil, fmt.Errorf("no cache store available")
	}

	var payload ExchangePayload
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		return nil, fmt.Errorf("failed to decode exchange payload: %w", err)
	}

	return &payload, nil
}

// ExtractCleanHost extracts hostname without port and separates port.
func (s *LoginResolutionService) ExtractCleanHost(host string) (cleanHost string, port string) {
	h := strings.ToLower(strings.TrimSpace(host))
	if colonIdx := strings.LastIndex(h, ":"); colonIdx != -1 {
		return h[:colonIdx], h[colonIdx:]
	}
	return h, ""
}

// ResolvePlatformHost resolves the platform control console hostname.
func (s *LoginResolutionService) ResolvePlatformHost(reqHost string) string {
	clean, port := s.ExtractCleanHost(reqHost)
	if strings.HasSuffix(clean, ".localhost") || clean == "localhost" || clean == "127.0.0.1" {
		return "app.localhost" + port
	}
	if strings.HasSuffix(clean, ".curexal.space") || clean == "curexal.space" {
		return "app.curexal.space" + port
	}
	if strings.HasSuffix(clean, ".curexal.internal") || clean == "curexal.internal" {
		return "app.curexal.internal" + port
	}
	return "app.localhost" + port
}

// ResolveTargetHost constructs the target workspace hostname based on current environment.
func (s *LoginResolutionService) ResolveTargetHost(reqHost string, orgSlug string, customDomain string) string {
	if customDomain != "" {
		return customDomain
	}
	clean, port := s.ExtractCleanHost(reqHost)
	if strings.HasSuffix(clean, ".localhost") || clean == "localhost" || clean == "127.0.0.1" {
		return fmt.Sprintf("%s.localhost%s", orgSlug, port)
	}
	if strings.HasSuffix(clean, ".curexal.space") || clean == "curexal.space" {
		return fmt.Sprintf("%s.curexal.space%s", orgSlug, port)
	}
	if strings.HasSuffix(clean, ".curexal.internal") || clean == "curexal.internal" {
		return fmt.Sprintf("%s.curexal.internal%s", orgSlug, port)
	}
	return fmt.Sprintf("%s.localhost%s", orgSlug, port)
}

// BuildTargetURL constructs the full URL for cross-host handoff.
func (s *LoginResolutionService) BuildTargetURL(reqHost string, targetHost string, pathAndQuery string) string {
	scheme := "https"
	clean, _ := s.ExtractCleanHost(reqHost)
	if strings.Contains(clean, "localhost") || strings.Contains(clean, "127.0.0.1") {
		scheme = "http"
	}
	if !strings.HasPrefix(pathAndQuery, "/") {
		pathAndQuery = "/" + pathAndQuery
	}
	return fmt.Sprintf("%s://%s%s", scheme, targetHost, pathAndQuery)
}

// ResolveDestination determines destination path, target host, and single-use exchange handoff.
func (s *LoginResolutionService) ResolveDestination(
	ctx context.Context,
	user *model.User,
	reqHost string,
	requestedOrgSlug *string,
	requestedBranchID *string,
	requestedBranchCode *string,
) (*LoginResolution, error) {
	// 1. Platform Super Admin / Platform Staff
	isPlatformStaff := user.IsPlatformAdmin || (user.PlatformRole != nil && (*user.PlatformRole == "super_admin" || *user.PlatformRole == "super_sales_staff"))
	if isPlatformStaff {
		targetHost := s.ResolvePlatformHost(reqHost)
		destPath := "/platform/audit"
		cleanReqHost, portReq := s.ExtractCleanHost(reqHost)
		currentFullHost := cleanReqHost + portReq

		if currentFullHost == targetHost || cleanReqHost == "app.localhost" || cleanReqHost == "app.curexal.space" {
			return &LoginResolution{
				Status:          StatusAuthenticated,
				DestinationPath: destPath,
				TargetHost:      targetHost,
			}, nil
		}

		payload := &ExchangePayload{
			UserID:          user.ID,
			Role:            "super_admin",
			TargetHost:      targetHost,
			DestinationPath: destPath,
			CreatedAt:       time.Now().Unix(),
		}
		rawToken, err := s.CreateExchangeToken(ctx, payload)
		if err != nil {
			return nil, err
		}
		targetURL := s.BuildTargetURL(reqHost, targetHost, "/auth/exchange?token="+rawToken)

		return &LoginResolution{
			Status:          StatusRedirectRequired,
			TargetURL:       targetURL,
			TargetHost:      targetHost,
			DestinationPath: destPath,
		}, nil
	}

	// 2. Fetch User's Active Organization Memberships
	rows, err := s.server.DB.Pool.Query(ctx, `
		SELECT o.id::text, o.name, o.slug, COALESCE(o.custom_domain, ''), m.role
		FROM organization.organization_memberships m
		JOIN organization.organizations o ON o.id = m.organization_id
		WHERE m.user_id = $1 AND m.is_active = TRUE AND o.status = 'active'
		ORDER BY (m.role IN ('owner', 'org_admin')) DESC, m.created_at ASC
	`, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to query organization memberships: %w", err)
	}
	defer rows.Close()

	var orgs []OrganizationOption
	for rows.Next() {
		var opt OrganizationOption
		if errScan := rows.Scan(&opt.ID, &opt.Name, &opt.Slug, &opt.CustomDomain, &opt.Role); errScan == nil {
			orgs = append(orgs, opt)
		}
	}

	if len(orgs) == 0 {
		return &LoginResolution{
			Status:  StatusOrgAccessRequired,
			Message: "No active organization memberships found for your account.",
		}, nil
	}

	// 3. Resolve Target Organization
	var chosenOrg *OrganizationOption
	if requestedOrgSlug != nil && *requestedOrgSlug != "" {
		for i := range orgs {
			if strings.EqualFold(orgs[i].Slug, *requestedOrgSlug) {
				chosenOrg = &orgs[i]
				break
			}
		}
		if chosenOrg == nil {
			return nil, errors.New("unauthorized organization access")
		}
	} else {
		// Attempt to match incoming host subdomain
		cleanReqHost, _ := s.ExtractCleanHost(reqHost)
		for i := range orgs {
			if strings.EqualFold(orgs[i].Slug, cleanReqHost) ||
				strings.HasPrefix(cleanReqHost, strings.ToLower(orgs[i].Slug)+".") ||
				(orgs[i].CustomDomain != "" && strings.EqualFold(orgs[i].CustomDomain, cleanReqHost)) {
				chosenOrg = &orgs[i]
				break
			}
		}

		// Fallback if not matching incoming host
		if chosenOrg == nil {
			if len(orgs) == 1 {
				chosenOrg = &orgs[0]
			} else {
				return &LoginResolution{
					Status:        StatusOrgSelectionRequired,
					Organizations: orgs,
					Message:       "Please select an organization workspace to enter.",
				}, nil
			}
		}
	}

	// 4. Resolve Branch Context within Selected Organization
	branchRes, errRes := s.authService.ResolveBranchContext(ctx, user.ID, &chosenOrg.Slug, requestedBranchID, requestedBranchCode)
	if errRes != nil {
		if errors.Is(errRes, domain.ErrUnassignedFacilityBranch) {
			return &LoginResolution{
				Status:       StatusUnassignedFacilityBranch,
				Organization: chosenOrg,
				Message:      "No operational branch assignment found for your account. Access to clinical workspaces is restricted.",
			}, nil
		}
		if errors.Is(errRes, domain.ErrUnauthorizedBranchAccess) {
			return nil, domain.ErrUnauthorizedBranchAccess
		}
		return nil, errRes
	}

	if branchRes.RequireBranchSelection {
		return &LoginResolution{
			Status:           StatusBranchSelectionRequired,
			SelectionToken:   branchRes.SelectionToken,
			AssignedBranches: branchRes.AssignedBranches,
			Organization:     chosenOrg,
		}, nil
	}

	// 5. Build Workspace Destination
	destPath := "/workspace/dashboard"
	activeBranchID := ""
	if branchRes.ActiveBranch != nil {
		activeBranchID = branchRes.ActiveBranch.ID
	} else {
		// Executive leadership with no operational clinic branch assignments
		destPath = "/organization/dashboard"
	}

	targetHost := s.ResolveTargetHost(reqHost, chosenOrg.Slug, chosenOrg.CustomDomain)
	cleanReqHost, portReq := s.ExtractCleanHost(reqHost)
	currentFullHost := cleanReqHost + portReq

	if currentFullHost == targetHost {
		return &LoginResolution{
			Status:          StatusAuthenticated,
			DestinationPath: destPath,
			TargetHost:      targetHost,
			Organization:    chosenOrg,
			ActiveBranch:    branchRes.ActiveBranch,
		}, nil
	}

	// Cross-host handoff required
	payload := &ExchangePayload{
		UserID:          user.ID,
		OrganizationID:  chosenOrg.ID,
		BranchID:        activeBranchID,
		TenantSlug:      chosenOrg.Slug,
		Role:            chosenOrg.Role,
		TargetHost:      targetHost,
		DestinationPath: destPath,
		CreatedAt:       time.Now().Unix(),
	}
	rawToken, err := s.CreateExchangeToken(ctx, payload)
	if err != nil {
		return nil, err
	}
	targetURL := s.BuildTargetURL(reqHost, targetHost, "/auth/exchange?token="+rawToken)

	return &LoginResolution{
		Status:          StatusRedirectRequired,
		TargetURL:       targetURL,
		TargetHost:      targetHost,
		DestinationPath: destPath,
		Organization:    chosenOrg,
		ActiveBranch:    branchRes.ActiveBranch,
	}, nil
}
