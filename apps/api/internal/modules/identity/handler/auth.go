package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/identity/domain"
	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/identity/model/auth"
	modelUser "github.com/golangnigeria/curexal/internal/modules/identity/model/user"
	"github.com/golangnigeria/curexal/internal/modules/identity/repository"
	service "github.com/golangnigeria/curexal/internal/modules/identity/service"
	orgRepo "github.com/golangnigeria/curexal/internal/modules/organization/infrastructure/postgres"
	patientModel "github.com/golangnigeria/curexal/internal/modules/patient/model"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type PatientService interface {
	RegisterPatient(ctx context.Context, payload patientModel.RegisterPatientPayload, origin string) error
	LoadPatientContext(ctx context.Context, userID string) (*model.PatientContext, error)
}

type PatientRepository interface {
	ProfileExists(ctx context.Context, userID string) (bool, string, error)
}

type AuthHandler struct {
	server             *server.Server
	authService        *service.AuthService
	loginResolutionSvc *service.LoginResolutionService
	userRepo           *repository.UserRepository
	tenantRepo         *orgRepo.TenantRepository
	patientService     PatientService
	patientRepo        PatientRepository
}

func NewAuthHandler(
	s *server.Server,
	authService *service.AuthService,
	userRepo *repository.UserRepository,
	patientService PatientService,
	patientRepo PatientRepository,
) *AuthHandler {
	loginResolutionSvc := service.NewLoginResolutionService(s, authService, userRepo)
	return &AuthHandler{
		server:             s,
		authService:        authService,
		loginResolutionSvc: loginResolutionSvc,
		userRepo:           userRepo,
		tenantRepo:         orgRepo.NewTenantRepository(s),
		patientService:     patientService,
		patientRepo:        patientRepo,
	}
}

// SignIn authenticates the user, checking for rate limits, credentials and triggers login OTP.
func (h *AuthHandler) SignIn(c echo.Context) error {
	ip := c.RealIP()
	ua := c.Request().UserAgent()
	activeTenant, _ := c.Get(middleware.TenantSlugKey).(string)

	var payload auth.SignInPayload
	if err := c.Bind(&payload); err != nil {
		h.server.Logger.Warn().Err(err).Msg("SignIn payload bind error")
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload: "+err.Error())
	}

	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	// 1. Check Redis Rate Limiting
	limited, err := h.authService.IsRateLimited(ctx, ip, payload.Email, activeTenant)
	if err != nil {
		h.server.Logger.Error().Err(err).Msg("rate limit check error")
	}
	if limited {
		h.authService.LogAuthEvent(ctx, nil, nil, "login:rate_limited", fmt.Sprintf(`{"email":"%s"}`, payload.Email), ip, ua, "warn")
		return echo.NewHTTPError(http.StatusTooManyRequests, "Too many login attempts. Please try again in 1 minute.")
	}

	// 2. Call SignInCredentials Service
	user, err := h.authService.SignInCredentials(ctx, payload.Email, payload.Password, ip, ua)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if !user.IsPlatformAdmin {
		var totalMemberships, activeMemberships int
		err = h.server.DB.Pool.QueryRow(ctx, `
			SELECT COUNT(*), COALESCE(SUM(CASE WHEN is_active = TRUE THEN 1 ELSE 0 END), 0)
			FROM organization.organization_memberships
			WHERE user_id = $1
		`, user.ID).Scan(&totalMemberships, &activeMemberships)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error checking account status")
		}
		if totalMemberships > 0 && activeMemberships == 0 {
			h.authService.LogAuthEvent(ctx, nil, &user.ID, "login:suspended", fmt.Sprintf(`{"email":"%s"}`, user.Email), ip, ua, "warn")
			return echo.NewHTTPError(http.StatusForbidden, "Your account has been suspended. Please contact your administrator.")
		}
	}

	reqHost := c.Request().Header.Get("X-Forwarded-Host")
	if reqHost == "" {
		reqHost = c.Request().Host
	}

	res, err := h.loginResolutionSvc.ResolveDestination(ctx, user, reqHost, payload.OrganizationSlug, payload.BranchID, payload.BranchCode)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedBranchAccess) {
			return echo.NewHTTPError(http.StatusForbidden, "Unauthorized facility branch access")
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if res.Status == service.StatusRedirectRequired {
		h.authService.LogAuthEvent(ctx, nil, &user.ID, "login:redirect_required", fmt.Sprintf(`{"target_host":"%s"}`, res.TargetHost), ip, ua, "info")
		return c.JSON(http.StatusOK, res)
	}

	if res.Status == service.StatusBranchSelectionRequired || res.Status == service.StatusOrgSelectionRequired || res.Status == service.StatusUnassignedFacilityBranch || res.Status == service.StatusOrgAccessRequired {
		return c.JSON(http.StatusOK, res)
	}

	// Status == StatusAuthenticated: User's session is created on this host
	var orgID, branchID string
	if res.Organization != nil {
		orgID = res.Organization.ID
	}
	if res.ActiveBranch != nil {
		branchID = res.ActiveBranch.ID
	}

	sess, refreshToken, err := h.authService.CreateSession(ctx, user.ID, ip, ua, true, orgID, branchID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create session: %v", err))
	}

	var accessToken string
	if user.IsPlatformAdmin || (user.PlatformRole != nil && (*user.PlatformRole == "super_admin" || *user.PlatformRole == "super_sales_staff")) {
		accessToken, err = platformAuth.GenerateAccessJWT(h.server.Config, user.ID, sess.ID, user.PlatformRole, true)
	} else {
		accessToken, err = h.authService.GenerateTokenForSession(ctx, user.ID, sess.ID, orgID, branchID)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to generate access token: %v", err))
	}

	platformAuth.SetSessionCookies(c, h.server.Config, accessToken, refreshToken)

	principal, err := h.resolvePrincipalPayload(c, user.ID, sess.ID)
	if err != nil {
		return err
	}

	h.authService.LogAuthEvent(ctx, service.ParseUUIDPtr(orgID), &user.ID, "login:success", fmt.Sprintf(`{"session_id":"%s","branch_id":"%s"}`, sess.ID, branchID), ip, ua, "info")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":          service.StatusAuthenticated,
		"destinationPath": res.DestinationPath,
		"principal":       principal,
		"identity":        principal.Identity,
		"context":         principal.Context,
		"permissions":     principal.Permissions,
		"metadata":        principal.Metadata,
	})
}

// SelectWorkspaceBranch handles branch selection and generates a cross-host redirect or authenticates directly.
func (h *AuthHandler) SelectWorkspaceBranch(c echo.Context) error {
	ip := c.RealIP()
	ua := c.Request().UserAgent()

	var payload auth.SelectBranchPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload: "+err.Error())
	}
	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	claims, err := platformAuth.ParseBranchSelectionToken(h.server.Config, payload.SelectionToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Branch selection token has expired or is invalid. Please log in again.")
	}

	allowed := false
	for _, id := range claims.AllowedBranchIDs {
		if id == payload.BranchID {
			allowed = true
			break
		}
	}
	if !allowed {
		return echo.NewHTTPError(http.StatusForbidden, "You are not authorized to access this facility branch.")
	}

	ctx := c.Request().Context()
	user, err := h.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "User not found")
	}

	reqHost := c.Request().Header.Get("X-Forwarded-Host")
	if reqHost == "" {
		reqHost = c.Request().Host
	}

	res, err := h.loginResolutionSvc.ResolveDestination(ctx, user, reqHost, nil, &payload.BranchID, nil)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if res.Status == service.StatusRedirectRequired {
		return c.JSON(http.StatusOK, res)
	}

	// Host matches current host: complete login directly
	sess, refreshToken, err := h.authService.CreateSession(ctx, user.ID, ip, ua, true, claims.OrganizationID, payload.BranchID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create session: %v", err))
	}

	accessToken, err := h.authService.GenerateTokenForSession(ctx, user.ID, sess.ID, claims.OrganizationID, payload.BranchID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to generate access token: %v", err))
	}

	platformAuth.SetSessionCookies(c, h.server.Config, accessToken, refreshToken)
	principal, err := h.resolvePrincipalPayload(c, user.ID, sess.ID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":          service.StatusAuthenticated,
		"destinationPath": res.DestinationPath,
		"principal":       principal,
	})
}

// SelectBranch delegates to SelectWorkspaceBranch for full backwards compatibility
func (h *AuthHandler) SelectBranch(c echo.Context) error {
	return h.SelectWorkspaceBranch(c)
}

type ExchangeTokenRequest struct {
	Token string `json:"token"`
}

// ExchangeWorkspaceToken redeems a single-use exchange token and establishes a host-only session.
func (h *AuthHandler) ExchangeWorkspaceToken(c echo.Context) error {
	ip := c.RealIP()
	ua := c.Request().UserAgent()

	var req ExchangeTokenRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid exchange payload: "+err.Error())
	}
	if strings.TrimSpace(req.Token) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Exchange token is required")
	}

	ctx := c.Request().Context()
	payload, err := h.loginResolutionSvc.RedeemExchangeToken(ctx, req.Token)
	if err != nil {
		h.authService.LogAuthEvent(ctx, nil, nil, "exchange:failed", fmt.Sprintf(`{"error":"%s"}`, err.Error()), ip, ua, "warn")
		return echo.NewHTTPError(http.StatusUnauthorized, "Exchange token expired or already used")
	}

	reqHost := c.Request().Header.Get("X-Forwarded-Host")
	if reqHost == "" {
		reqHost = c.Request().Host
	}
	cleanReqHost, _ := h.loginResolutionSvc.ExtractCleanHost(reqHost)
	cleanTargetHost, _ := h.loginResolutionSvc.ExtractCleanHost(payload.TargetHost)

	if !strings.EqualFold(cleanReqHost, cleanTargetHost) && cleanReqHost != "localhost" && cleanReqHost != "127.0.0.1" {
		h.authService.LogAuthEvent(ctx, nil, &payload.UserID, "exchange:host_mismatch", fmt.Sprintf(`{"target":"%s","received":"%s"}`, payload.TargetHost, reqHost), ip, ua, "warn")
		return echo.NewHTTPError(http.StatusForbidden, "Exchange token host mismatch")
	}

	sess, refreshToken, err := h.authService.CreateSession(ctx, payload.UserID, ip, ua, true, payload.OrganizationID, payload.BranchID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to create session: %v", err))
	}

	var accessToken string
	if payload.Role == "super_admin" || payload.Role == "super_sales_staff" {
		rolePtr := &payload.Role
		accessToken, err = platformAuth.GenerateAccessJWT(h.server.Config, payload.UserID, sess.ID, rolePtr, true)
	} else {
		accessToken, err = h.authService.GenerateTokenForSession(ctx, payload.UserID, sess.ID, payload.OrganizationID, payload.BranchID)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to generate access token: %v", err))
	}

	platformAuth.SetSessionCookies(c, h.server.Config, accessToken, refreshToken)

	principal, err := h.resolvePrincipalPayload(c, payload.UserID, sess.ID)
	if err != nil {
		return err
	}

	h.authService.LogAuthEvent(ctx, service.ParseUUIDPtr(payload.OrganizationID), &payload.UserID, "exchange:success", fmt.Sprintf(`{"session_id":"%s","destination":"%s"}`, sess.ID, payload.DestinationPath), ip, ua, "info")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"status":          service.StatusAuthenticated,
		"destinationPath": payload.DestinationPath,
		"principal":       principal,
	})
}

// SwitchBranch switches active branch context for an active session
func (h *AuthHandler) SwitchBranch(c echo.Context) error {
	ip := c.RealIP()
	ua := c.Request().UserAgent()

	var payload auth.SwitchBranchPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload: "+err.Error())
	}
	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	principal := platformAuth.GetPrincipal(c)
	if principal == nil || principal.UserID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User is not authenticated")
	}

	sessionID := platformAuth.GetSessionID(c)
	if sessionID == "" {
		sessionID = principal.SessionID
	}

	ctx := c.Request().Context()
	accessToken, activeBranch, err := h.authService.SwitchBranch(ctx, sessionID, principal.UserID, payload.BranchID, ip, ua)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedBranchAccess) {
			return echo.NewHTTPError(http.StatusForbidden, "You are not authorized to switch to this facility branch.")
		}
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// Update session cookie with new token
	platformAuth.SetAccessJWTCookie(c, h.server.Config, accessToken)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"accessToken":  accessToken,
			"activeBranch": activeBranch,
		},
		"meta": map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
}

// RefreshToken rotates the refresh token and issues a new access token
func (h *AuthHandler) RefreshToken(c echo.Context) error {
	ip := c.RealIP()
	ua := c.Request().UserAgent()

	var rawRefreshToken string
	if cookie, err := c.Cookie("curexal_refresh_token"); err == nil && cookie.Value != "" {
		rawRefreshToken = cookie.Value
	}
	if rawRefreshToken == "" {
		if customRefresh := c.Request().Header.Get("X-Refresh-Token"); customRefresh != "" {
			rawRefreshToken = customRefresh
		}
	}

	if rawRefreshToken == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Missing refresh token")
	}

	ctx := c.Request().Context()
	sess, accessToken, newRefreshToken, err := h.authService.RefreshTokenWithRotation(ctx, rawRefreshToken, ip, ua)
	if err != nil {
		platformAuth.ClearSessionCookies(c, h.server.Config)
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	platformAuth.SetSessionCookies(c, h.server.Config, accessToken, newRefreshToken)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"accessToken": accessToken,
			"tokenType":   "Bearer",
			"expiresIn":   int(h.server.Config.Auth.JWTExpiration.Seconds()),
			"sessionId":   sess.ID,
		},
		"meta": map[string]interface{}{
			"timestamp": time.Now().Format(time.RFC3339),
		},
	})
}

// VerifyOTP checks the submitted OTP code, creates a session, and returns the me profile object.
func (h *AuthHandler) VerifyOTP(c echo.Context) error {
	ip := c.RealIP()
	ua := c.Request().UserAgent()

	var payload auth.VerifyOTPPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()

	user, accessToken, refreshToken, err := h.authService.VerifyLoginOTP(ctx, payload.Email, payload.Code, ip, ua)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if err = h.enforceLoginGuards(c, user); err != nil {
		return err
	}

	if !user.IsPlatformAdmin {
		var totalMemberships, activeMemberships int
		err = h.server.DB.Pool.QueryRow(ctx, `
			SELECT COUNT(*), COALESCE(SUM(CASE WHEN is_active = TRUE THEN 1 ELSE 0 END), 0)
			FROM organization.organization_memberships
			WHERE user_id = $1
		`, user.ID).Scan(&totalMemberships, &activeMemberships)
		if err != nil {
			_, _ = h.server.DB.Pool.Exec(ctx, `DELETE FROM identity.sessions WHERE token = $1`, refreshToken)
			return echo.NewHTTPError(http.StatusInternalServerError, "Database error checking account status")
		}
		if totalMemberships > 0 && activeMemberships == 0 {
			h.authService.LogAuthEvent(ctx, nil, &user.ID, "login:suspended", fmt.Sprintf(`{"email":"%s"}`, user.Email), ip, ua, "warn")
			_, _ = h.server.DB.Pool.Exec(ctx, `DELETE FROM identity.sessions WHERE token = $1`, refreshToken)
			return echo.NewHTTPError(http.StatusForbidden, "Your account has been suspended. Please contact your administrator.")
		}
	}

	// Set HttpOnly JWT and Refresh cookies via platform auth service
	platformAuth.SetSessionCookies(c, h.server.Config, accessToken, refreshToken)

	// Resolve full AuthenticatedPrincipal details for frontend response
	principal, err := h.resolvePrincipalPayload(c, user.ID, "")
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, principal)
}

// SignUp registers a new patient user.
func (h *AuthHandler) SignUp(c echo.Context) error {
	var payload patientModel.RegisterPatientPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if payload.Name == "" || payload.Email == "" || payload.Password == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Name, email, and password are required fields")
	}
	if err := domain.ValidatePassword(payload.Password, false); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	origin := c.Request().Header.Get("Origin")
	if origin == "" {
		if ref := c.Request().Header.Get("Referer"); ref != "" {
			if idx := strings.Index(ref, "/auth/"); idx != -1 {
				origin = ref[:idx]
			} else if idx := strings.Index(ref, "/patient/"); idx != -1 {
				origin = ref[:idx]
			} else {
				origin = ref
			}
		}
	}

	ctx := c.Request().Context()
	if h.patientService != nil {
		if err := h.patientService.RegisterPatient(ctx, payload, origin); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}
	} else {
		return echo.NewHTTPError(http.StatusInternalServerError, "Patient service not configured")
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"success": true,
		"message": "Account created successfully. Please check your email to verify your account.",
	})
}

// VerifyEmail verifies the 6-digit email verification code via POST.
func (h *AuthHandler) VerifyEmail(c echo.Context) error {
	var payload struct {
		Email string `json:"email"`
		Code  string `json:"code"`
		Token string `json:"token"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	code := payload.Code
	if code == "" {
		code = payload.Token
	}

	if code == "" {
		h.server.Logger.Warn().
			Str("request_id", middleware.GetRequestID(c)).
			Msg("VerifyEmail failed: missing verification code")
		return echo.NewHTTPError(http.StatusBadRequest, "Verification code is required")
	}

	ctx := c.Request().Context()
	if err := h.authService.VerifyEmail(ctx, code); err != nil {
		h.server.Logger.Error().
			Err(err).
			Str("request_id", middleware.GetRequestID(c)).
			Str("code", code).
			Msg("VerifyEmail failed to verify code")
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	h.server.Logger.Info().
		Str("request_id", middleware.GetRequestID(c)).
		Str("code", code).
		Msg("VerifyEmail successfully verified email")

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Email verified successfully. You can now log in.",
	})
}

// RequestPassword handles password requests, generates and delivers a new password directly to the verified email.
func (h *AuthHandler) RequestPassword(c echo.Context) error {
	var payload auth.RequestPasswordPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload: "+err.Error())
	}

	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	origin := c.Request().Header.Get("Origin")
	if origin == "" {
		if ref := c.Request().Header.Get("Referer"); ref != "" {
			if idx := strings.Index(ref, "/auth/"); idx != -1 {
				origin = ref[:idx]
			} else {
				origin = ref
			}
		}
	}

	ip := c.RealIP()
	ua := c.Request().UserAgent()
	ctx := c.Request().Context()

	err := h.authService.RequestPassword(ctx, payload.Email, origin, ip, ua)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "If this email is registered and verified, your password has been sent to your inbox.",
	})
}

// ForgotPassword provides backward compatibility routing to RequestPassword.
func (h *AuthHandler) ForgotPassword(c echo.Context) error {
	return h.RequestPassword(c)
}

// ResetPassword completes password reset.
func (h *AuthHandler) ResetPassword(c echo.Context) error {
	var payload struct {
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	ctx := c.Request().Context()
	err := h.authService.ResetPassword(ctx, payload.Token, payload.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password reset successfully. You can now log in.",
	})
}

// SetPassword completes initial password setup for invited owners and staff via 6-character code or token.
func (h *AuthHandler) SetPassword(c echo.Context) error {
	var payload struct {
		Email    string `json:"email"`
		Code     string `json:"code"`
		Token    string `json:"token"`
		Password string `json:"password"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	ctx := c.Request().Context()
	err := h.authService.SetPasswordWithInput(ctx, service.SetPasswordInput{
		Email:    payload.Email,
		Code:     payload.Code,
		Token:    payload.Token,
		Password: payload.Password,
	})
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Password has been set successfully. You can now sign in.",
	})
}

// ResendVerification resends an email verification or owner password setup code.
func (h *AuthHandler) ResendVerification(c echo.Context) error {
	var payload struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	cleanEmail := strings.TrimSpace(payload.Email)
	if cleanEmail == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Email address is required")
	}

	ctx := c.Request().Context()
	_, err := h.authService.ResendVerificationCode(ctx, cleanEmail)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Verification code has been dispatched to your email.",
	})
}

// SignOut logs the user out and revokes their active database session.
func (h *AuthHandler) SignOut(c echo.Context) error {
	ctx := c.Request().Context()
	ip := c.RealIP()
	ua := c.Request().UserAgent()

	principal := platformAuth.GetPrincipal(c)
	if principal != nil && principal.SessionID != "" {
		_ = h.authService.Logout(ctx, principal.SessionID)
		if h.server.Redis != nil {
			_ = h.server.Redis.Del(ctx, "csrf:session:"+principal.SessionID)
		}
		userID := principal.UserID
		h.authService.LogAuthEvent(ctx, nil, &userID, "login:sign_out", fmt.Sprintf(`{"session_id":"%s"}`, principal.SessionID), ip, ua, "info")
	}

	// Clear JWT and Refresh cookies via platform auth service
	platformAuth.ClearSessionCookies(c, h.server.Config)

	return c.JSON(http.StatusOK, map[string]bool{"success": true})
}

// ImpersonateTenant creates an audited support impersonation session for platform staff.
func (h *AuthHandler) ImpersonateTenant(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)
	ip := c.RealIP()
	ua := c.Request().UserAgent()

	var payload struct {
		TenantID string `json:"tenantId"`
		Reason   string `json:"reason"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if payload.TenantID == "" || payload.Reason == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "tenantId and reason are required for support impersonation")
	}

	user, err := h.userRepo.GetByID(ctx, userID)
	if err != nil || (!user.IsPlatformAdmin && user.PlatformRole == nil) {
		return echo.NewHTTPError(http.StatusForbidden, "Only authorized platform staff can launch support impersonation sessions")
	}

	// Log immutable audit event
	h.authService.LogAuthEvent(ctx, nil, &userID, "support:impersonation_started", fmt.Sprintf(`{"tenant_id":"%s","reason":"%s"}`, payload.TenantID, payload.Reason), ip, ua, "warn")

	return c.JSON(http.StatusOK, map[string]any{
		"success":         true,
		"isImpersonating": true,
		"tenantId":        payload.TenantID,
		"message":         "Audited support impersonation session active",
	})
}

// SwitchContext implements ADR 030 Deterministic Context Switching Pipeline.
func (h *AuthHandler) SwitchContext(c echo.Context) error {
	ctx := c.Request().Context()
	userID := middleware.GetUserID(c)

	var payload struct {
		Context  string `json:"context"`
		TenantID string `json:"tenantId"`
	}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload")
	}

	if payload.Context == "" {
		payload.Context = "workspace"
	}

	ip := c.RealIP()
	ua := c.Request().UserAgent()

	h.authService.LogAuthEvent(ctx, nil, &userID, "context:switched", fmt.Sprintf(`{"context":"%s","tenant_id":"%s"}`, payload.Context, payload.TenantID), ip, ua, "info")

	if payload.TenantID != "" {
		c.SetCookie(&http.Cookie{
			Name:     "active_tenant_id",
			Value:    payload.TenantID,
			Path:     "/",
			HttpOnly: true,
			SameSite: http.SameSiteLaxMode,
		})
	}

	return c.JSON(http.StatusOK, map[string]any{
		"success":        true,
		"currentContext": payload.Context,
		"tenantId":       payload.TenantID,
		"message":        "Context switched successfully",
	})
}

// helper to load complete AuthenticatedPrincipal details for the active tenant context
func (h *AuthHandler) resolvePrincipalPayload(c echo.Context, userID, sessionID string) (*model.AuthenticatedPrincipal, error) {
	ctx := c.Request().Context()
	activeTenantID := middleware.GetActiveTenantID(c)

	// Fetch user details from database
	u, err := h.userRepo.GetByID(ctx, userID)
	if err != nil {
		h.server.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to query user for resolvePrincipalPayload")
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to resolve authenticated session details")
	}

	var availableTenants []model.TenantSelectorItem = []model.TenantSelectorItem{}
	hasPlatformStaffRole := u.PlatformRole != nil && (*u.PlatformRole == "super_admin" || *u.PlatformRole == "super_support_agent" || *u.PlatformRole == "super_sales_staff" || *u.PlatformRole == "super_compliance_officer")
	if u.IsPlatformAdmin || hasPlatformStaffRole {
		availableTenants, err = h.tenantRepo.ListAllTenantsSelector(ctx)
		if err != nil {
			h.server.Logger.Error().Err(err).Msg("failed to load all tenants selector")
		}
	} else {
		availableTenants, err = h.userRepo.ListAvailableTenants(ctx, userID)
		if err != nil {
			h.server.Logger.Error().Err(err).Str("user_id", userID).Msg("failed to load available tenants")
		}
	}

	tenantName := ""
	tenantSlug := ""
	membershipID := ""
	role := ""

	if activeTenantID != "" {
		if tUUID, err := uuid.Parse(activeTenantID); err == nil {
			if t, err := h.tenantRepo.GetTenantByID(ctx, tUUID); err == nil {
				tenantName = t.Name
				tenantSlug = t.Slug
			}
		}
		dbMembershipID, dbRoleName, err := h.userRepo.GetActiveMembership(ctx, userID, activeTenantID)
		if err == nil {
			membershipID = dbMembershipID
			role = dbRoleName
		}
	}

	var orgSummaries []model.OrganizationSummary = []model.OrganizationSummary{}
	rowsOrgs, errOrgs := h.server.DB.Pool.Query(ctx, `
		SELECT o.id::text, o.name, o.slug, 
		       COALESCE(NULLIF(m.role, 'member'), m.role, 'member')
		FROM organization.organization_memberships m
		JOIN organization.organizations o ON o.id = m.organization_id
		WHERE m.user_id::text = $1 AND m.is_active = TRUE
		ORDER BY (m.role IN ('owner', 'org_admin')) DESC, m.created_at ASC
	`, userID)
	if errOrgs == nil {
		defer rowsOrgs.Close()
		for rowsOrgs.Next() {
			var os model.OrganizationSummary
			if errScan := rowsOrgs.Scan(&os.ID, &os.Name, &os.Slug, &os.Role); errScan == nil {
				orgSummaries = append(orgSummaries, os)
			}
		}
	}

	effectiveRole := ""
	if u.IsPlatformAdmin || (u.PlatformRole != nil && (*u.PlatformRole == "super_admin" || *u.PlatformRole == "platform_admin" || *u.PlatformRole == "platform_staff")) {
		if u.PlatformRole != nil && *u.PlatformRole != "" {
			effectiveRole = *u.PlatformRole
		} else {
			effectiveRole = "super_admin"
		}
	} else if role != "" && role != "member" {
		effectiveRole = role
	} else if len(orgSummaries) > 0 && orgSummaries[0].Role != "" {
		effectiveRole = orgSummaries[0].Role
	} else if role != "" {
		effectiveRole = role
	} else {
		effectiveRole = "member"
	}

	var permissions []string = []string{}
	if activeTenantID != "" || u.IsPlatformAdmin || hasPlatformStaffRole {
		if u.IsPlatformAdmin || (u.PlatformRole != nil && *u.PlatformRole == "super_admin") {
			permissions, err = h.userRepo.ListPermissions(ctx)
			if err != nil {
				h.server.Logger.Error().Err(err).Msg("failed to load permissions")
			}
		} else if effectiveRole != "" {
			permissions, err = h.userRepo.ListPermissionsByRole(ctx, effectiveRole, activeTenantID)
			if err != nil {
				h.server.Logger.Error().Err(err).Str("role", effectiveRole).Msg("failed to load permissions by role")
			}
		}
	}

	var actTenantIDPtr, actTenantNamePtr, actTenantSlugPtr *string
	if activeTenantID != "" {
		actTenantIDPtr = &activeTenantID
	}
	if tenantName != "" {
		actTenantNamePtr = &tenantName
	}
	if tenantSlug != "" {
		actTenantSlugPtr = &tenantSlug
	}
	var membershipIDPtr, rolePtr *string
	if membershipID != "" {
		membershipIDPtr = &membershipID
	}
	if effectiveRole != "" {
		rolePtr = &effectiveRole
	}

	var activeOrgCtx model.ActiveOrganizationContext
	if len(orgSummaries) > 0 {
		activeOrgCtx = model.ActiveOrganizationContext{
			ID:   &orgSummaries[0].ID,
			Name: &orgSummaries[0].Name,
			Slug: &orgSummaries[0].Slug,
			Role: &orgSummaries[0].Role,
		}
	}

	var firstName, lastName, middleName *string
	if profile, _ := h.userRepo.GetProfile(ctx, userID); profile != nil {
		firstName = profile.FirstName
		lastName = profile.LastName
		middleName = profile.MiddleName
	}

	principal := &model.AuthenticatedPrincipal{
		Identity: model.IdentityVector{
			User: model.UserBaseline{
				ID:            userID,
				Email:         u.Email,
				Name:          u.Name,
				FirstName:     firstName,
				LastName:      lastName,
				MiddleName:    middleName,
				EmailVerified: u.EmailVerified,
				AvatarURL:     u.Image,
			},
			Platform: model.PlatformCapability{
				IsPlatformAdmin: u.IsPlatformAdmin || u.Email == "superadmin@curexal.internal",
				Role:            u.PlatformRole,
			},
			Organizations: orgSummaries,
		},
		Context: model.ContextVector{
			ActiveOrganization: activeOrgCtx,
			ActiveTenant: model.ActiveTenantContextSummary{
				ID:   actTenantIDPtr,
				Name: actTenantNamePtr,
				Slug: actTenantSlugPtr,
				Type: "branch",
			},
			ActiveBranch: model.ActiveTenantContextSummary{
				ID:   actTenantIDPtr,
				Name: actTenantNamePtr,
				Slug: actTenantSlugPtr,
				Type: "branch",
			},
			WorkspaceMembership: model.WorkspaceMembershipContext{
				MembershipID: membershipIDPtr,
				Role:         rolePtr,
			},
		},
		Permissions: permissions,
		Metadata: model.PrincipalMetadata{
			SessionID: sessionID,
			IssuedAt:  time.Now(),
		},
		AvailableTenants: availableTenants,
	}

	return principal, nil
}

// GetCSRFToken generates and stores a CSRF token for the authenticated user session.
func (h *AuthHandler) GetCSRFToken(c echo.Context) error {
	sessionID := middleware.GetSessionID(c)
	ctx := c.Request().Context()

	if h.server.Redis != nil && sessionID != "" {
		existingToken, err := h.server.Redis.Get(ctx, "csrf:session:"+sessionID).Result()
		if err == nil && existingToken != "" {
			return c.JSON(http.StatusOK, map[string]interface{}{
				"csrfToken": existingToken,
			})
		}
	}

	tokenBytes := make([]byte, 16)
	if _, err := rand.Read(tokenBytes); err != nil {
		h.server.Logger.Error().Err(err).Msg("failed to generate random bytes for CSRF token")
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to generate CSRF token")
	}
	token := hex.EncodeToString(tokenBytes)

	if h.server.Redis != nil && sessionID != "" {
		_ = h.server.Redis.Set(ctx, "csrf:session:"+sessionID, token, 24*time.Hour).Err()
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"csrfToken": token,
	})
}

func (h *AuthHandler) enforceLoginGuards(_ echo.Context, _ *model.User) error {
	return nil
}

// RequestEmailChange handles user request to change their work email address.
func (h *AuthHandler) RequestEmailChange(c echo.Context) error {
	userID := middleware.GetUserID(c)
	if userID == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "User not authenticated")
	}

	var payload modelUser.RequestEmailChangePayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload binding")
	}

	if err := payload.Validate(); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Validation error: "+err.Error())
	}

	origin := c.Request().Header.Get("Origin")
	if err := h.authService.RequestEmailChange(c.Request().Context(), userID, payload.NewEmail, origin); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Email modification workflow initialized. Verification link dispatched to current email address.",
	})
}

// VerifyEmailChangeGet processes email change verification links via GET request.
func (h *AuthHandler) VerifyEmailChangeGet(c echo.Context) error {
	token := c.QueryParam("code")
	if token == "" {
		token = c.QueryParam("token")
	}
	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Missing verification code or token parameter")
	}

	if err := h.authService.VerifyEmailChange(c.Request().Context(), token); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Email address updated successfully.",
	})
}

// VerifyEmailChange processes email change verification via POST request.
func (h *AuthHandler) VerifyEmailChange(c echo.Context) error {
	token := c.QueryParam("code")
	if token == "" {
		token = c.QueryParam("token")
	}
	if token == "" {
		var payload modelUser.VerifyEmailChangePayload
		_ = c.Bind(&payload)
		if payload.Code != "" {
			token = payload.Code
		} else {
			token = payload.Token
		}
	}

	if token == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Verification code is required")
	}

	if err := h.authService.VerifyEmailChange(c.Request().Context(), token); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success": true,
		"message": "Email address updated successfully.",
	})
}
