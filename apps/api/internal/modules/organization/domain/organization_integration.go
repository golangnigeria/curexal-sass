package domain

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAPIKeyNotFound              = errors.New("organization API key record not found")
	ErrWebhookSubscriptionNotFound = errors.New("organization webhook subscription not found")
	ErrDuplicateWebhookTarget      = errors.New("a webhook subscription with this target URL already exists")
	ErrInvalidIPWhitelist          = errors.New("invalid IP CIDR address in whitelist constraint")
	ErrSSRFURLForbidden            = errors.New("webhook target URL violates SSRF security boundary: internal, loopback, and private network addresses are forbidden")
	ErrInvalidWebhookEvent         = errors.New("unrecognized or invalid webhook event type")
)

var CanonicalScopes = map[string]bool{
	"patients:read":      true,
	"patients:write":     true,
	"appointments:read":  true,
	"appointments:write": true,
	"laboratory:read":    true,
	"laboratory:write":   true,
	"billing:read":       true,
	"billing:write":      true,
	"webhooks:read":      true,
	"webhooks:write":     true,
}

func IsValidScope(scope string) bool {
	return CanonicalScopes[strings.ToLower(strings.TrimSpace(scope))]
}

// ValidateSSRFSafeURL enforces strict outbound network SSRF boundaries.
func ValidateSSRFSafeURL(targetURL string) error {
	parsed, err := url.Parse(targetURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return ErrSSRFURLForbidden
	}

	hostname := strings.ToLower(parsed.Hostname())
	if hostname == "" || hostname == "localhost" || strings.HasSuffix(hostname, ".local") || strings.HasSuffix(hostname, ".internal") {
		return ErrSSRFURLForbidden
	}

	ip := net.ParseIP(hostname)
	if ip != nil {
		if isForbiddenIP(ip) {
			return ErrSSRFURLForbidden
		}
	}

	return nil
}

func isForbiddenIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() || ip.IsUnspecified() {
		return true
	}
	return false
}

type APIKey struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organizationId"`
	Name           string          `json:"name"`
	KeyPrefix      string          `json:"keyPrefix"`
	KeyHash        string          `json:"-"` // Never returned in JSON
	Scopes         json.RawMessage `json:"scopes,omitempty"`
	IPWhitelist    json.RawMessage `json:"ipWhitelist,omitempty"`
	RateLimitRPM   int             `json:"rateLimitRpm"`
	ExpiresAt      *time.Time      `json:"expiresAt,omitempty"`
	LastUsedAt     *time.Time      `json:"lastUsedAt,omitempty"`
	IsActive       bool            `json:"isActive"`
	Version        int             `json:"version"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	CreatedBy      *uuid.UUID      `json:"createdBy,omitempty"`
}

type APIKeyCreateResult struct {
	APIKey APIKey `json:"apiKey"`
	RawKey string `json:"rawKey"` // Shown ONCE at creation time
}

type WebhookSubscription struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organizationId"`
	Name           string          `json:"name"`
	TargetURL      string          `json:"targetUrl"`
	EventTypes     json.RawMessage `json:"eventTypes,omitempty"`
	SigningSecret  *string         `json:"signingSecret,omitempty"` // Redacted as "••••••••"
	IsActive       bool            `json:"isActive"`
	Version        int             `json:"version"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
	CreatedBy      *uuid.UUID      `json:"createdBy,omitempty"`
}

type WebhookDelivery struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organizationId"`
	SubscriptionID *uuid.UUID      `json:"subscriptionId,omitempty"`
	EventType      string          `json:"eventType"`
	Payload        json.RawMessage `json:"payload,omitempty"`
	ResponseStatus *int            `json:"responseStatus,omitempty"`
	ResponseBody   *string         `json:"responseBody,omitempty"`
	AttemptCount   int             `json:"attemptCount"`
	LastError      *string         `json:"lastError,omitempty"`
	Status         string          `json:"status"` // 'PENDING', 'SUCCESS', 'FAILED'
	CreatedAt      time.Time       `json:"createdAt"`
}

type CreateAPIKeyPayload struct {
	Name         string          `json:"name"`
	Scopes       json.RawMessage `json:"scopes"`
	IPWhitelist  json.RawMessage `json:"ipWhitelist"`
	RateLimitRPM *int            `json:"rateLimitRpm"`
	ExpiresAt    *time.Time      `json:"expiresAt"`
}

type CreateWebhookSubscriptionPayload struct {
	Name          string          `json:"name"`
	TargetURL     string          `json:"targetUrl"`
	EventTypes    json.RawMessage `json:"eventTypes"`
	SigningSecret *string         `json:"signingSecret"`
}

// ValidateIPWhitelist validates CIDR notation for listed IP whitelist rules.
func ValidateIPWhitelist(ips []string) error {
	for _, ipStr := range ips {
		trimmed := strings.TrimSpace(ipStr)
		if trimmed == "" {
			continue
		}
		if !strings.Contains(trimmed, "/") {
			trimmed = trimmed + "/32"
		}
		_, _, err := net.ParseCIDR(trimmed)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidIPWhitelist, ipStr)
		}
	}
	return nil
}
