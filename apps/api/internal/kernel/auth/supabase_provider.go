package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// SupabaseConfig holds the configuration needed for the Supabase authentication provider.
type SupabaseConfig struct {
	ProjectURL     string
	AnonKey        string
	ServiceRoleKey string
	JWTSecret      string
}

// SupabaseIdentityProvider integrates Curexal's kernel auth with Supabase Auth (GoTrue).
type SupabaseIdentityProvider struct {
	config     SupabaseConfig
	httpClient *http.Client
}

// NewSupabaseIdentityProvider creates a new instance of SupabaseIdentityProvider.
func NewSupabaseIdentityProvider(cfg SupabaseConfig) *SupabaseIdentityProvider {
	return &SupabaseIdentityProvider{
		config: cfg,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (p *SupabaseIdentityProvider) Name() string {
	return "supabase"
}

// SupabaseClaims represents the JWT claims standard in Supabase Auth.
type SupabaseClaims struct {
	jwt.RegisteredClaims
	Email        string                 `json:"email"`
	Phone        string                 `json:"phone,omitempty"`
	Role         string                 `json:"role,omitempty"`
	UserMetadata map[string]interface{} `json:"user_metadata,omitempty"`
	AppMetadata  map[string]interface{} `json:"app_metadata,omitempty"`
}

// Authenticate validates a Supabase JWT token either locally using JWTSecret or by verifying against Supabase GoTrue.
func (p *SupabaseIdentityProvider) Authenticate(ctx context.Context, sessionToken string) (*IdentitySession, error) {
	cleanToken := strings.TrimSpace(strings.TrimPrefix(sessionToken, "Bearer "))
	if cleanToken == "" {
		return nil, ErrSessionNotFound
	}

	// 1. If JWT secret is configured, perform local cryptographic verification
	if p.config.JWTSecret != "" {
		token, err := jwt.ParseWithClaims(cleanToken, &SupabaseClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(p.config.JWTSecret), nil
		})

		if err == nil && token.Valid {
			if claims, ok := token.Claims.(*SupabaseClaims); ok && claims.Subject != "" {
				name := ""
				avatar := ""
				if claims.UserMetadata != nil {
					if n, ok := claims.UserMetadata["name"].(string); ok {
						name = n
					} else if fn, ok := claims.UserMetadata["full_name"].(string); ok {
						name = fn
					}
					if a, ok := claims.UserMetadata["avatar_url"].(string); ok {
						avatar = a
					}
				}

				expiresAt := time.Now().Add(1 * time.Hour)
				if claims.ExpiresAt != nil {
					expiresAt = claims.ExpiresAt.Time
				}

				return &IdentitySession{
					ID:         "sess_" + claims.Subject,
					IdentityID: claims.Subject,
					Active:     true,
					Traits: IdentityTraits{
						Email:     claims.Email,
						Name:      name,
						AvatarURL: avatar,
					},
					ExpiresAt:       expiresAt,
					AuthenticatedAt: time.Now(),
				}, nil
			}
		}
	}

	// 2. Fallback: Verify token against Supabase Auth API endpoint
	if p.config.ProjectURL != "" {
		url := fmt.Sprintf("%s/auth/v1/user", strings.TrimRight(p.config.ProjectURL, "/"))
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", "Bearer "+cleanToken)
		if p.config.AnonKey != "" {
			req.Header.Set("apikey", p.config.AnonKey)
		}

		resp, err := p.httpClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("failed to contact Supabase Auth: %w", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, ErrSessionNotFound
		}

		var supaUser struct {
			ID           string                 `json:"id"`
			Email        string                 `json:"email"`
			UserMetadata map[string]interface{} `json:"user_metadata"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&supaUser); err != nil {
			return nil, err
		}

		name := ""
		avatar := ""
		if supaUser.UserMetadata != nil {
			if n, ok := supaUser.UserMetadata["name"].(string); ok {
				name = n
			} else if fn, ok := supaUser.UserMetadata["full_name"].(string); ok {
				name = fn
			}
			if a, ok := supaUser.UserMetadata["avatar_url"].(string); ok {
				avatar = a
			}
		}

		return &IdentitySession{
			ID:         "sess_" + supaUser.ID,
			IdentityID: supaUser.ID,
			Active:     true,
			Traits: IdentityTraits{
				Email:     supaUser.Email,
				Name:      name,
				AvatarURL: avatar,
			},
			ExpiresAt:       time.Now().Add(1 * time.Hour),
			AuthenticatedAt: time.Now(),
		}, nil
	}

	return nil, ErrSessionNotFound
}

// GetSession retrieves an active session by session ID.
func (p *SupabaseIdentityProvider) GetSession(ctx context.Context, sessionID string) (*IdentitySession, error) {
	identityID := strings.TrimPrefix(sessionID, "sess_")
	traits, err := p.GetIdentity(ctx, identityID)
	if err != nil {
		return nil, err
	}

	return &IdentitySession{
		ID:              sessionID,
		IdentityID:      identityID,
		Active:          true,
		Traits:          *traits,
		ExpiresAt:       time.Now().Add(1 * time.Hour),
		AuthenticatedAt: time.Now(),
	}, nil
}

// GetIdentity fetches user traits from Supabase Auth via Admin API.
func (p *SupabaseIdentityProvider) GetIdentity(ctx context.Context, identityID string) (*IdentityTraits, error) {
	if p.config.ProjectURL == "" || p.config.ServiceRoleKey == "" {
		return nil, errors.New("supabase project URL and service role key required for admin operations")
	}

	url := fmt.Sprintf("%s/auth/v1/admin/users/%s", strings.TrimRight(p.config.ProjectURL, "/"), identityID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+p.config.ServiceRoleKey)
	req.Header.Set("apikey", p.config.ServiceRoleKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrIdentityNotFound
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase admin request failed with status: %d", resp.StatusCode)
	}

	var supaUser struct {
		ID           string                 `json:"id"`
		Email        string                 `json:"email"`
		UserMetadata map[string]interface{} `json:"user_metadata"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&supaUser); err != nil {
		return nil, err
	}

	name := ""
	avatar := ""
	if supaUser.UserMetadata != nil {
		if n, ok := supaUser.UserMetadata["name"].(string); ok {
			name = n
		} else if fn, ok := supaUser.UserMetadata["full_name"].(string); ok {
			name = fn
		}
		if a, ok := supaUser.UserMetadata["avatar_url"].(string); ok {
			avatar = a
		}
	}

	return &IdentityTraits{
		Email:     supaUser.Email,
		Name:      name,
		AvatarURL: avatar,
	}, nil
}

// CreateIdentity provisions a new user in Supabase Auth via Admin API.
func (p *SupabaseIdentityProvider) CreateIdentity(ctx context.Context, params *CreateIdentityParams) (string, error) {
	if params == nil {
		return "", errors.New("params cannot be nil")
	}
	if p.config.ProjectURL == "" || p.config.ServiceRoleKey == "" {
		return "", errors.New("supabase project URL and service role key required")
	}

	url := fmt.Sprintf("%s/auth/v1/admin/users", strings.TrimRight(p.config.ProjectURL, "/"))
	bodyPayload := map[string]interface{}{
		"email":         params.Traits.Email,
		"email_confirm": true,
		"user_metadata": map[string]interface{}{
			"name":       params.Traits.Name,
			"avatar_url": params.Traits.AvatarURL,
		},
	}
	if params.Password != "" {
		bodyPayload["password"] = params.Password
	}

	jsonBytes, err := json.Marshal(bodyPayload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.ServiceRoleKey)
	req.Header.Set("apikey", p.config.ServiceRoleKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("supabase user creation failed (status %d): %s", resp.StatusCode, string(body))
	}

	var res struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.ID, nil
}

// UpdateIdentity modifies user traits in Supabase Auth via Admin API.
func (p *SupabaseIdentityProvider) UpdateIdentity(ctx context.Context, identityID string, params *UpdateIdentityParams) error {
	if params == nil {
		return errors.New("params cannot be nil")
	}
	if p.config.ProjectURL == "" || p.config.ServiceRoleKey == "" {
		return errors.New("supabase project URL and service role key required")
	}

	url := fmt.Sprintf("%s/auth/v1/admin/users/%s", strings.TrimRight(p.config.ProjectURL, "/"), identityID)
	bodyPayload := map[string]interface{}{
		"user_metadata": map[string]interface{}{
			"name":       params.Traits.Name,
			"avatar_url": params.Traits.AvatarURL,
		},
	}
	if params.Traits.Email != "" {
		bodyPayload["email"] = params.Traits.Email
	}

	jsonBytes, err := json.Marshal(bodyPayload)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(jsonBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+p.config.ServiceRoleKey)
	req.Header.Set("apikey", p.config.ServiceRoleKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrIdentityNotFound
	}
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase user update failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}

// DeleteIdentity removes a user in Supabase Auth via Admin API.
func (p *SupabaseIdentityProvider) DeleteIdentity(ctx context.Context, identityID string) error {
	if p.config.ProjectURL == "" || p.config.ServiceRoleKey == "" {
		return errors.New("supabase project URL and service role key required")
	}

	url := fmt.Sprintf("%s/auth/v1/admin/users/%s", strings.TrimRight(p.config.ProjectURL, "/"), identityID)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+p.config.ServiceRoleKey)
	req.Header.Set("apikey", p.config.ServiceRoleKey)

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil // idempotent deletion
	}
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("supabase user deletion failed (status %d): %s", resp.StatusCode, string(body))
	}

	return nil
}
