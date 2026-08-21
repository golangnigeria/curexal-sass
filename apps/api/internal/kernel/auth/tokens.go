package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golangnigeria/curexal/internal/shared/config"
)

type UserClaims struct {
	SessionID        string  `json:"sid,omitempty"`
	PlatformRole     *string `json:"platform_role,omitempty"`
	OrganizationRole *string `json:"org_role,omitempty"`
	IsPlatformAdmin  bool    `json:"is_platform_admin,omitempty"`
	jwt.RegisteredClaims
}

// GenerateAccessJWT signs a short-lived access token using configuration options.
func GenerateAccessJWT(cfg *config.Config, userID, sessionID string, platformRole *string, isPlatformAdmin bool, orgRole ...*string) (string, error) {
	if cfg.Auth.SecretKey == "" {
		return "", errors.New("missing auth secret key in configuration")
	}

	var organizationRole *string
	if len(orgRole) > 0 {
		organizationRole = orgRole[0]
	}

	claims := UserClaims{
		SessionID:        sessionID,
		PlatformRole:     platformRole,
		OrganizationRole: organizationRole,
		IsPlatformAdmin:  isPlatformAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.Auth.JWTExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Auth.SecretKey))
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
