package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/notification/model"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/jackc/pgx/v5"
)

type NotificationRepository struct {
	server *server.Server
}

func NewNotificationRepository(s *server.Server) *NotificationRepository {
	return &NotificationRepository{server: s}
}

// CreateNotificationWithOutbox performs an atomic transaction inserting notification and outbox event
func (r *NotificationRepository) CreateNotificationWithOutbox(ctx context.Context, n *model.Notification, outbox *model.OutboxEvent) error {
	return r.server.DB.RunInTx(ctx, func(txCtx context.Context) error {
		db := r.server.DB.Conn(txCtx)

		// 1. Insert Notification
		_, err := db.Exec(txCtx, `
			INSERT INTO notification.notifications (
				id, user_id, tenant_id, title, message, type, channel, priority, delivery_status,
				notification_type, is_read, metadata, link_url, patient_id, visit_id, order_id,
				specimen_id, result_id, triggered_by
			) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		`, n.ID, n.UserID, n.TenantID, n.Title, n.Message, n.Type, n.Channel, n.Priority, n.DeliveryStatus,
			n.Type, n.IsRead, n.Metadata, n.LinkURL, n.PatientID, n.VisitID, n.OrderID,
			n.SpecimenID, n.ResultID, n.TriggeredBy)
		if err != nil {
			return err
		}

		// 2. Insert Outbox Event if provided
		if outbox != nil {
			payloadBytes, _ := json.Marshal(outbox.Payload)
			_, err = db.Exec(txCtx, `
				INSERT INTO notification.outbox_events (
					id, aggregate_type, aggregate_id, event_type, payload, status, attempts, max_attempts, next_retry_at
				) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			`, outbox.ID, outbox.AggregateType, outbox.AggregateID, outbox.EventType, payloadBytes, outbox.Status, outbox.Attempts, outbox.MaxAttempts, outbox.NextRetryAt)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *NotificationRepository) UpdateNotificationStatus(ctx context.Context, id string, status model.DeliveryStatus) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `
		UPDATE notification.notifications
		SET delivery_status = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`, id, status)
	return err
}

func (r *NotificationRepository) GetNotificationByID(ctx context.Context, id string) (*model.Notification, error) {
	db := r.server.DB.Conn(ctx)
	n := &model.Notification{}
	err := db.QueryRow(ctx, `
		SELECT id, user_id, tenant_id, title, message, type, channel, priority, delivery_status,
		       is_read, read_at, metadata, link_url, patient_id, visit_id, order_id, specimen_id, result_id, triggered_by, created_at, updated_at
		FROM notification.notifications
		WHERE id = $1
	`, id).Scan(
		&n.ID, &n.UserID, &n.TenantID, &n.Title, &n.Message, &n.Type, &n.Channel, &n.Priority, &n.DeliveryStatus,
		&n.IsRead, &n.ReadAt, &n.Metadata, &n.LinkURL, &n.PatientID, &n.VisitID, &n.OrderID, &n.SpecimenID, &n.ResultID, &n.TriggeredBy, &n.CreatedAt, &n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (r *NotificationRepository) GetUserNotifications(ctx context.Context, userID string, unreadOnly bool, limit, offset int) ([]*model.Notification, error) {
	db := r.server.DB.Conn(ctx)
	query := `
		SELECT id, user_id, tenant_id, title, message, type, channel, priority, delivery_status,
		       is_read, read_at, metadata, link_url, patient_id, visit_id, order_id, specimen_id, result_id, triggered_by, created_at, updated_at
		FROM notification.notifications
		WHERE user_id = $1
	`
	if unreadOnly {
		query += ` AND is_read = FALSE`
	}
	query += ` ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*model.Notification
	for rows.Next() {
		n := &model.Notification{}
		err := rows.Scan(
			&n.ID, &n.UserID, &n.TenantID, &n.Title, &n.Message, &n.Type, &n.Channel, &n.Priority, &n.DeliveryStatus,
			&n.IsRead, &n.ReadAt, &n.Metadata, &n.LinkURL, &n.PatientID, &n.VisitID, &n.OrderID, &n.SpecimenID, &n.ResultID, &n.TriggeredBy, &n.CreatedAt, &n.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, n)
	}
	return result, nil
}

func (r *NotificationRepository) MarkAsRead(ctx context.Context, id, userID string) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `
		UPDATE notification.notifications
		SET is_read = TRUE, read_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
	`, id, userID)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, userID string) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `
		UPDATE notification.notifications
		SET is_read = TRUE, read_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND is_read = FALSE
	`, userID)
	return err
}

func (r *NotificationRepository) GetUnreadCount(ctx context.Context, userID string) (int, int, error) {
	db := r.server.DB.Conn(ctx)

	var unreadNotifs int
	err := db.QueryRow(ctx, `SELECT COUNT(*) FROM notification.notifications WHERE user_id = $1 AND is_read = FALSE`, userID).Scan(&unreadNotifs)
	if err != nil {
		return 0, 0, err
	}

	var unreadMsgs int
	err = db.QueryRow(ctx, `SELECT COUNT(*) FROM notification.messages WHERE recipient_id = $1 AND is_read = FALSE`, userID).Scan(&unreadMsgs)
	if err != nil && !errorsIsPgx(err) {
		unreadMsgs = 0
	}

	return unreadNotifs, unreadMsgs, nil
}

// Outbox Operations
func (r *NotificationRepository) FetchPendingOutboxEvents(ctx context.Context, limit int) ([]*model.OutboxEvent, error) {
	db := r.server.DB.Conn(ctx)
	rows, err := db.Query(ctx, `
		SELECT id, aggregate_type, aggregate_id, event_type, payload, status, attempts, max_attempts, error_message, next_retry_at, created_at, processed_at
		FROM notification.outbox_events
		WHERE status IN ('pending', 'processing') AND next_retry_at <= CURRENT_TIMESTAMP
		ORDER BY created_at ASC
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.OutboxEvent
	for rows.Next() {
		e := &model.OutboxEvent{}
		var payloadBytes []byte
		if err := rows.Scan(&e.ID, &e.AggregateType, &e.AggregateID, &e.EventType, &payloadBytes, &e.Status, &e.Attempts, &e.MaxAttempts, &e.ErrorMessage, &e.NextRetryAt, &e.CreatedAt, &e.ProcessedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(payloadBytes, &e.Payload)
		events = append(events, e)
	}
	return events, nil
}

func (r *NotificationRepository) UpdateOutboxStatus(ctx context.Context, id, status string, attempts int, errMsg *string, nextRetryAt time.Time) error {
	db := r.server.DB.Conn(ctx)
	var processedAt *time.Time
	if status == "completed" || status == "failed" {
		now := time.Now()
		processedAt = &now
	}
	_, err := db.Exec(ctx, `
		UPDATE notification.outbox_events
		SET status = $2, attempts = $3, error_message = $4, next_retry_at = $5, processed_at = $6
		WHERE id = $1
	`, id, status, attempts, errMsg, nextRetryAt, processedAt)
	return err
}

// Preference Operations
func (r *NotificationRepository) GetNotificationPreference(ctx context.Context, userID string) (*model.NotificationPreference, error) {
	db := r.server.DB.Conn(ctx)
	pref := &model.NotificationPreference{}
	err := db.QueryRow(ctx, `
		SELECT user_id, email_enabled, sms_enabled, push_enabled, whatsapp_enabled, quiet_hours_start, quiet_hours_end, updated_at
		FROM notification.preferences
		WHERE user_id = $1
	`, userID).Scan(&pref.UserID, &pref.EmailEnabled, &pref.SMSEnabled, &pref.PushEnabled, &pref.WhatsAppEnabled, &pref.QuietHoursStart, &pref.QuietHoursEnd, &pref.UpdatedAt)
	if err != nil {
		if errorsIsPgx(err) {
			// Default preferences
			return &model.NotificationPreference{
				UserID:       userID,
				EmailEnabled: true,
				PushEnabled:  true,
				UpdatedAt:    time.Now(),
			}, nil
		}
		return nil, err
	}
	return pref, nil
}

func (r *NotificationRepository) UpsertNotificationPreference(ctx context.Context, pref *model.NotificationPreference) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `
		INSERT INTO notification.preferences (user_id, email_enabled, sms_enabled, push_enabled, whatsapp_enabled, quiet_hours_start, quiet_hours_end, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id) DO UPDATE SET
			email_enabled = EXCLUDED.email_enabled,
			sms_enabled = EXCLUDED.sms_enabled,
			push_enabled = EXCLUDED.push_enabled,
			whatsapp_enabled = EXCLUDED.whatsapp_enabled,
			quiet_hours_start = EXCLUDED.quiet_hours_start,
			quiet_hours_end = EXCLUDED.quiet_hours_end,
			updated_at = CURRENT_TIMESTAMP
	`, pref.UserID, pref.EmailEnabled, pref.SMSEnabled, pref.PushEnabled, pref.WhatsAppEnabled, pref.QuietHoursStart, pref.QuietHoursEnd)
	return err
}

func (r *NotificationRepository) CreateMessage(ctx context.Context, msg *model.Message) error {
	db := r.server.DB.Conn(ctx)
	_, err := db.Exec(ctx, `
		INSERT INTO notification.messages (id, sender_id, recipient_id, subject, body, is_read)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, msg.ID, msg.SenderID, msg.RecipientID, msg.Subject, msg.Body, msg.IsRead)
	return err
}

func (r *NotificationRepository) GetUserMessages(ctx context.Context, userID string, limit, offset int) ([]*model.Message, error) {
	db := r.server.DB.Conn(ctx)
	rows, err := db.Query(ctx, `
		SELECT id, sender_id, recipient_id, subject, body, is_read, read_at, created_at
		FROM notification.messages
		WHERE recipient_id = $1 OR sender_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		m := &model.Message{}
		if err := rows.Scan(&m.ID, &m.SenderID, &m.RecipientID, &m.Subject, &m.Body, &m.IsRead, &m.ReadAt, &m.CreatedAt); err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func errorsIsPgx(err error) bool {
	return err == pgx.ErrNoRows
}
