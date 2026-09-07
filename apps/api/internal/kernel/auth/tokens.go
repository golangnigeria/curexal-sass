package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golangnigeria/curexal/internal/shared/config"
)

// UserClaims conforms to the Day 1 Production Specification for Access Tokens.
type UserClaims struct {
	Email            string   `json:"email,omitempty"`
	OrganizationID   string   `json:"org_id,omitempty"`
	ActiveBranchID   string   `json:"branch_id,omitempty"`
	Roles            []string `json:"roles,omitempty"`
	PermHash         string   `json:"perm_hash,omitempty"`
	SessionID        string   `json:"sid,omitempty"`
	PlatformRole     *string  `json:"platform_role,omitempty"`
	OrganizationRole *string  `json:"org_role,omitempty"`
	IsPlatformAdmin  bool     `json:"is_platform_admin,omitempty"`
	jwt.RegisteredClaims
}

// BranchSelectionClaims is used for short-lived (5 min) tokens when a user has multiple assigned branches.
type BranchSelectionClaims struct {
	UserID           string   `json:"sub"`
	OrganizationID   string   `json:"org_id"`
	AllowedBranchIDs []string `json:"allowed_branch_ids"`
	TokenType        string   `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateTokenWithClaims generates an access JWT with full custom claims.
func GenerateTokenWithClaims(cfg *config.Config, claims UserClaims) (string, error) {
	if cfg.Auth.SecretKey == "" {
		return "", errors.New("missing auth secret key in configuration")
	}

	expDuration := cfg.Auth.JWTExpiration
	if expDuration <= 0 {
		expDuration = 15 * time.Minute
	}

	now := time.Now()
	if claims.IssuedAt == nil {
		claims.IssuedAt = jwt.NewNumericDate(now)
	}
	if claims.ExpiresAt == nil {
		claims.ExpiresAt = jwt.NewNumericDate(now.Add(expDuration))
	}
	if claims.Issuer == "" {
		claims.Issuer = "curexal-auth-engine"
	}
	if len(claims.Audience) == 0 {
		claims.Audience = jwt.ClaimStrings{"curexal-clinic-os"}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Auth.SecretKey))
}

// GenerateAccessJWT signs a short-lived access token using standard options (backward-compatible).
func GenerateAccessJWT(cfg *config.Config, userID, sessionID string, platformRole *string, isPlatformAdmin bool, orgRole ...*string) (string, error) {
	var organizationRole *string
	if len(orgRole) > 0 {
		organizationRole = orgRole[0]
	}

	var roles []string
	if platformRole != nil && *platformRole != "" {
		roles = append(roles, *platformRole)
	}
	if organizationRole != nil && *organizationRole != "" {
		roles = append(roles, *organizationRole)
	}

	claims := UserClaims{
		SessionID:        sessionID,
		PlatformRole:     platformRole,
		OrganizationRole: organizationRole,
		IsPlatformAdmin:  isPlatformAdmin,
		Roles:            roles,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userID,
		},
	}

	return GenerateTokenWithClaims(cfg, claims)
}

// GenerateBranchSelectionToken issues a 5-minute token allowing the user to select an operational branch.
func GenerateBranchSelectionToken(cfg *config.Config, userID, orgID string, allowedBranchIDs []string) (string, error) {
	if cfg.Auth.SecretKey == "" {
		return "", errors.New("missing auth secret key in configuration")
	}

	claims := BranchSelectionClaims{
		UserID:           userID,
		OrganizationID:   orgID,
		AllowedBranchIDs: allowedBranchIDs,
		TokenType:        "branch_selection",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    "curexal-auth-engine",
			Audience:  jwt.ClaimStrings{"curexal-clinic-os"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Auth.SecretKey))
}

// ParseBranchSelectionToken validates and extracts claims from a branch selection token.
func ParseBranchSelectionToken(cfg *config.Config, tokenStr string) (*BranchSelectionClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("branch selection token string is empty")
	}
	if cfg.Auth.SecretKey == "" {
		return nil, errors.New("missing auth secret key in configuration")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &BranchSelectionClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.Auth.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired branch selection token")
	}

	claims, ok := token.Claims.(*BranchSelectionClaims)
	if !ok || claims.TokenType != "branch_selection" {
		return nil, errors.New("invalid token payload type")
	}

	return claims, nil
}

// ParseAccessJWT validates and extracts claims from an incoming access JWT string.
func ParseAccessJWT(cfg *config.Config, tokenStr string) (*UserClaims, error) {
	if tokenStr == "" {
		return nil, errors.New("token string is empty")
	}
	if cfg.Auth.SecretKey == "" {
		return nil, errors.New("missing auth secret key in configuration")
	}

	token, err := jwt.ParseWithClaims(tokenStr, &UserClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.Auth.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired JWT token")
	}

	claims, ok := token.Claims.(*UserClaims)
	if !ok {
		return nil, errors.New("invalid token claims format")
	}

	return claims, nil
}
