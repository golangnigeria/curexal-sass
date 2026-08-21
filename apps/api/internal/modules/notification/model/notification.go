package model

import (
	"time"
)

type NotificationPriority string

const (
	PriorityLow      NotificationPriority = "LOW"
	PriorityNormal   NotificationPriority = "NORMAL"
	PriorityHigh     NotificationPriority = "HIGH"
	PriorityCritical NotificationPriority = "CRITICAL"
)

type DeliveryStatus string

const (
	StatusPending    DeliveryStatus = "pending"
	StatusQueued     DeliveryStatus = "queued"
	StatusProcessing DeliveryStatus = "processing"
	StatusSent       DeliveryStatus = "sent"
	StatusDelivered  DeliveryStatus = "delivered"
	StatusFailed     DeliveryStatus = "failed"
	StatusBounced    DeliveryStatus = "bounced"
)

type NotificationType string

const (
	TypePatientResultReady NotificationType = "PATIENT_RESULT_READY"
	TypePatientRegistered  NotificationType = "PATIENT_REGISTERED"
	TypeOrderCompleted     NotificationType = "ORDER_COMPLETED"
	TypeSpecimenReceived   NotificationType = "SPECIMEN_RECEIVED"
	TypePasswordReset      NotificationType = "PASSWORD_RESET"
	TypeLoginAlert         NotificationType = "LOGIN_ALERT"
	TypeSubscriptionRenew  NotificationType = "SUBSCRIPTION_RENEWED"
	TypeSystemAlert        NotificationType = "SYSTEM_NOTIFICATION"
)

type Notification struct {
	ID               string                 `json:"id"`
	UserID           string                 `json:"user_id"`
	TenantID         *string                `json:"tenant_id,omitempty"`
	Title            string                 `json:"title"`
	Message          string                 `json:"message"`
	Type             string                 `json:"type"` // e.g. 'PATIENT_RESULT_READY'
	Channel          string                 `json:"channel"` // 'in_app', 'email', 'sms'
	Priority         NotificationPriority   `json:"priority"`
	DeliveryStatus   DeliveryStatus         `json:"delivery_status"`
	IsRead           bool                   `json:"is_read"`
	ReadAt           *time.Time             `json:"read_at,omitempty"`
	Metadata         map[string]interface{} `json:"metadata,omitempty"`
	LinkURL          *string                `json:"link_url,omitempty"`
	PatientID        *string                `json:"patient_id,omitempty"`
	VisitID          *string                `json:"visit_id,omitempty"`
	OrderID          *string                `json:"order_id,omitempty"`
	SpecimenID       *string                `json:"specimen_id,omitempty"`
	ResultID         *string                `json:"result_id,omitempty"`
	TriggeredBy      *string                `json:"triggered_by,omitempty"`
	CreatedAt        time.Time              `json:"created_at"`
	UpdatedAt        time.Time              `json:"updated_at"`
}

type CreateNotificationDTO struct {
	UserID         string                 `json:"user_id" validate:"required"`
	TenantID       *string                `json:"tenant_id,omitempty"`
	Title          string                 `json:"title" validate:"required"`
	Message        string                 `json:"message" validate:"required"`
	Type           NotificationType       `json:"type"`
	Channel        string                 `json:"channel"`
	Priority       NotificationPriority   `json:"priority"`
	Metadata       map[string]interface{} `json:"metadata,omitempty"`
	LinkURL        *string                `json:"link_url,omitempty"`
	PatientID      *string                `json:"patient_id,omitempty"`
	VisitID        *string                `json:"visit_id,omitempty"`
	OrderID        *string                `json:"order_id,omitempty"`
	SpecimenID     *string                `json:"specimen_id,omitempty"`
	ResultID       *string                `json:"result_id,omitempty"`
	TriggeredBy    *string                `json:"triggered_by,omitempty"`
	SendEmail      bool                   `json:"send_email"`
}

type OutboxEvent struct {
	ID            string                 `json:"id"`
	AggregateType string                 `json:"aggregate_type"`
	AggregateID   string                 `json:"aggregate_id"`
	EventType     string                 `json:"event_type"`
	Payload       map[string]interface{} `json:"payload"`
	Status        string                 `json:"status"` // 'pending', 'processing', 'completed', 'failed'
	Attempts      int                    `json:"attempts"`
	MaxAttempts   int                    `json:"max_attempts"`
	ErrorMessage  *string                `json:"error_message,omitempty"`
	NextRetryAt   time.Time              `json:"next_retry_at"`
	CreatedAt     time.Time              `json:"created_at"`
	ProcessedAt   *time.Time             `json:"processed_at,omitempty"`
}

type NotificationPreference struct {
	UserID           string     `json:"user_id"`
	EmailEnabled     bool       `json:"email_enabled"`
	SMSEnabled       bool       `json:"sms_enabled"`
	PushEnabled      bool       `json:"push_enabled"`
	WhatsAppEnabled  bool       `json:"whatsapp_enabled"`
	QuietHoursStart  *string    `json:"quiet_hours_start,omitempty"`
	QuietHoursEnd    *string    `json:"quiet_hours_end,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type Message struct {
	ID          string     `json:"id"`
	SenderID    string     `json:"sender_id"`
	RecipientID string     `json:"recipient_id"`
	Subject     *string    `json:"subject,omitempty"`
	Body        string     `json:"body"`
	IsRead      bool       `json:"is_read"`
	ReadAt      *time.Time `json:"read_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

type SendMessageDTO struct {
	RecipientID string  `json:"recipient_id" validate:"required"`
	Subject     *string `json:"subject,omitempty"`
	Body        string  `json:"body" validate:"required"`
}

type UnreadCountResponse struct {
	UnreadNotifications int `json:"unread_notifications"`
	UnreadMessages      int `json:"unread_messages"`
}
