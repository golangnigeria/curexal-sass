package domain

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrNotificationConfigNotFound   = errors.New("organization notification configuration not found")
	ErrNotificationTemplateNotFound = errors.New("organization notification template not found")
	ErrInvalidNotificationChannel   = errors.New("invalid notification channel: must be EMAIL, SMS, WHATSAPP, or IN_APP")
	ErrInvalidNotificationProvider  = errors.New("invalid notification provider: must be SMTP, RESEND, TWILIO, TERMII, META_WHATSAPP, or INTERNAL")
	ErrDuplicateCustomDomain        = errors.New("this custom domain is already registered to another organization")
	ErrUserNotificationNotFound     = errors.New("in-app notification record not found")
)

var AllowedNotificationChannels = map[string]bool{
	"EMAIL":    true,
	"SMS":      true,
	"WHATSAPP": true,
	"IN_APP":   true,
}

var AllowedNotificationProviders = map[string]bool{
	"SMTP":          true,
	"RESEND":        true,
	"TWILIO":        true,
	"TERMII":        true,
	"META_WHATSAPP": true,
	"INTERNAL":      true,
}

func IsValidNotificationChannel(channel string) bool {
	return AllowedNotificationChannels[strings.ToUpper(strings.TrimSpace(channel))]
}

func IsValidNotificationProvider(provider string) bool {
	return AllowedNotificationProviders[strings.ToUpper(strings.TrimSpace(provider))]
}

type BrandingConfig struct {
	OrganizationID     uuid.UUID       `json:"organizationId"`
	LogoURL            *string         `json:"logoUrl,omitempty"`
	PrimaryColor       string          `json:"primaryColor"`
	ThemeBranding      json.RawMessage `json:"themeBranding,omitempty"`
	CustomDomain       *string         `json:"customDomain,omitempty"`
	CustomDomainStatus string          `json:"customDomainStatus"` // 'PENDING', 'VERIFIED', 'ACTIVE', 'DISABLED'
	Version            int             `json:"version"`
	UpdatedAt          time.Time       `json:"updatedAt"`
	UpdatedBy          *uuid.UUID      `json:"updatedBy,omitempty"`
}

type NotificationConfig struct {
	ID             uuid.UUID       `json:"id"`
	OrganizationID uuid.UUID       `json:"organizationId"`
	Channel        string          `json:"channel"`
	Provider       string          `json:"provider"`
	SenderEmail    *string         `json:"senderEmail,omitempty"`
	SenderName     *string         `json:"senderName,omitempty"`
	Host           *string         `json:"host,omitempty"`
	Port           *int            `json:"port,omitempty"`
	Username       *string         `json:"username,omitempty"`
	Password       *string         `json:"password,omitempty"` // Redacted as "••••••••"
	APIKey         *string         `json:"apiKey,omitempty"`   // Redacted as "••••••••"
	ConfigMetadata json.RawMessage `json:"configMetadata,omitempty"`
	IsActive       bool            `json:"isActive"`
	Version        int             `json:"version"`
	CreatedAt      time.Time       `json:"createdAt"`
	UpdatedAt      time.Time       `json:"updatedAt"`
}

type NotificationTemplate struct {
	ID               uuid.UUID       `json:"id"`
	OrganizationID   uuid.UUID       `json:"organizationId"`
	TemplateKey      string          `json:"templateKey"`
	Channel          string          `json:"channel"`
	Subject          string          `json:"subject"`
	BodyTemplate     string          `json:"bodyTemplate"`
	AllowedVariables json.RawMessage `json:"allowedVariables,omitempty"`
	IsActive         bool            `json:"isActive"`
	Version          int             `json:"version"`
	CreatedAt        time.Time       `json:"createdAt"`
	UpdatedAt        time.Time       `json:"updatedAt"`
}

type UserNotification struct {
	ID               uuid.UUID       `json:"id"`
	OrganizationID   uuid.UUID       `json:"organizationId"`
	UserID           uuid.UUID       `json:"userId"`
	NotificationType string          `json:"notificationType"`
	Title            string          `json:"title"`
	Body             string          `json:"body"`
	Data             json.RawMessage `json:"data,omitempty"`
	ReadAt           *time.Time      `json:"readAt,omitempty"`
	CreatedAt        time.Time       `json:"createdAt"`
	ExpiresAt        *time.Time      `json:"expiresAt,omitempty"`
}

type NotificationDelivery struct {
	ID                uuid.UUID  `json:"id"`
	OrganizationID    uuid.UUID  `json:"organizationId"`
	NotificationID    *uuid.UUID `json:"notificationId,omitempty"`
	Channel           string     `json:"channel"`
	Provider          string     `json:"provider"`
	Recipient         string     `json:"recipient"`
	TemplateKey       string     `json:"templateKey"`
	Status            string     `json:"status"` // 'PENDING', 'QUEUED', 'SENDING', 'SENT', 'DELIVERED', 'FAILED', 'CANCELLED'
	ProviderMessageID *string    `json:"providerMessageId,omitempty"`
	AttemptCount      int        `json:"attemptCount"`
	LastError         *string    `json:"lastError,omitempty"`
	QueuedAt          *time.Time `json:"queuedAt,omitempty"`
	SentAt            *time.Time `json:"sentAt,omitempty"`
	DeliveredAt       *time.Time `json:"deliveredAt,omitempty"`
	FailedAt          *time.Time `json:"failedAt,omitempty"`
	CreatedAt         time.Time  `json:"createdAt"`
	UpdatedAt         time.Time  `json:"updatedAt"`
}

type UpdateBrandingPayload struct {
	LogoURL       *string         `json:"logoUrl"`
	PrimaryColor  *string         `json:"primaryColor"`
	ThemeBranding json.RawMessage `json:"themeBranding"`
	CustomDomain  *string         `json:"customDomain"`
	Version       int             `json:"version"`
}

type SaveNotificationConfigPayload struct {
	Channel        string          `json:"channel"`
	Provider       string          `json:"provider"`
	SenderEmail    *string         `json:"senderEmail"`
	SenderName     *string         `json:"senderName"`
	Host           *string         `json:"host"`
	Port           *int            `json:"port"`
	Username       *string         `json:"username"`
	Password       *string         `json:"password"`
	APIKey         *string         `json:"apiKey"`
	ConfigMetadata json.RawMessage `json:"configMetadata"`
	IsActive       *bool           `json:"isActive"`
}

type SaveNotificationTemplatePayload struct {
	Channel          string          `json:"channel"`
	Subject          string          `json:"subject"`
	BodyTemplate     string          `json:"bodyTemplate"`
	AllowedVariables json.RawMessage `json:"allowedVariables"`
}
