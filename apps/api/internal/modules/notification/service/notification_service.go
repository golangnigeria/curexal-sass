package service

import (
	"context"
	"fmt"
	"time"

	identityRepo "github.com/golangnigeria/curexal/internal/modules/identity/repository"
	"github.com/golangnigeria/curexal/internal/modules/notification/model"
	notificationRepo "github.com/golangnigeria/curexal/internal/modules/notification/repository"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
)

type NotificationService struct {
	server           *server.Server
	notificationRepo *notificationRepo.NotificationRepository
	userRepo         *identityRepo.UserRepository
}

func NewNotificationService(s *server.Server, notifRepo *notificationRepo.NotificationRepository, userRepo *identityRepo.UserRepository) *NotificationService {
	return &NotificationService{
		server:           s,
		notificationRepo: notifRepo,
		userRepo:         userRepo,
	}
}

func (s *NotificationService) CreateNotification(ctx context.Context, dto *model.CreateNotificationDTO) (*model.Notification, error) {
	if dto.Type == "" {
		dto.Type = model.TypeSystemAlert
	}
	if dto.Channel == "" {
		dto.Channel = "in_app"
	}
	if dto.Priority == "" {
		dto.Priority = model.PriorityNormal
	}

	notifID := uuid.New().String()

	n := &model.Notification{
		ID:             notifID,
		UserID:         dto.UserID,
		TenantID:       dto.TenantID,
		Title:          dto.Title,
		Message:        dto.Message,
		Type:           string(dto.Type),
		Channel:        dto.Channel,
		Priority:       dto.Priority,
		DeliveryStatus: model.StatusPending,
		IsRead:         false,
		Metadata:       dto.Metadata,
		LinkURL:        dto.LinkURL,
		PatientID:      dto.PatientID,
		VisitID:        dto.VisitID,
		OrderID:        dto.OrderID,
		SpecimenID:     dto.SpecimenID,
		ResultID:       dto.ResultID,
		TriggeredBy:    dto.TriggeredBy,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	var outbox *model.OutboxEvent
	if dto.SendEmail || dto.Channel == "email" {
		n.Channel = "email"
		outbox = &model.OutboxEvent{
			ID:            uuid.New().String(),
			AggregateType: "notification",
			AggregateID:   notifID,
			EventType:     string(dto.Type),
			Payload: map[string]interface{}{
				"notification_id": notifID,
				"user_id":         dto.UserID,
				"channel":         "email",
			},
			Status:      "pending",
			Attempts:    0,
			MaxAttempts: 5,
			NextRetryAt: time.Now(),
			CreatedAt:   time.Now(),
		}
	} else {
		// In-App notifications transition directly to delivered
		n.DeliveryStatus = model.StatusDelivered
	}

	// Transactionally save Notification and Outbox Event
	if err := s.notificationRepo.CreateNotificationWithOutbox(ctx, n, outbox); err != nil {
		return nil, fmt.Errorf("failed to save notification outbox transaction: %w", err)
	}

	return n, nil
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID string, unreadOnly bool, page, limit int) ([]*model.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return s.notificationRepo.GetUserNotifications(ctx, userID, unreadOnly, limit, offset)
}

func (s *NotificationService) MarkAsRead(ctx context.Context, id, userID string) error {
	return s.notificationRepo.MarkAsRead(ctx, id, userID)
}

func (s *NotificationService) MarkAllAsRead(ctx context.Context, userID string) error {
	return s.notificationRepo.MarkAllAsRead(ctx, userID)
}

func (s *NotificationService) GetUnreadCount(ctx context.Context, userID string) (*model.UnreadCountResponse, error) {
	notifCount, msgCount, err := s.notificationRepo.GetUnreadCount(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &model.UnreadCountResponse{
		UnreadNotifications: notifCount,
		UnreadMessages:      msgCount,
	}, nil
}

func (s *NotificationService) GetPreferences(ctx context.Context, userID string) (*model.NotificationPreference, error) {
	return s.notificationRepo.GetNotificationPreference(ctx, userID)
}

func (s *NotificationService) UpdatePreferences(ctx context.Context, pref *model.NotificationPreference) error {
	return s.notificationRepo.UpsertNotificationPreference(ctx, pref)
}

func (s *NotificationService) SendMessage(ctx context.Context, senderID string, dto *model.SendMessageDTO) (*model.Message, error) {
	msg := &model.Message{
		ID:          uuid.New().String(),
		SenderID:    senderID,
		RecipientID: dto.RecipientID,
		Subject:     dto.Subject,
		Body:        dto.Body,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	if err := s.notificationRepo.CreateMessage(ctx, msg); err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	}

	subjectStr := "New Direct Message"
	if dto.Subject != nil && *dto.Subject != "" {
		subjectStr = *dto.Subject
	}

	_, _ = s.CreateNotification(ctx, &model.CreateNotificationDTO{
		UserID:      dto.RecipientID,
		Title:       subjectStr,
		Message:     dto.Body,
		Type:        model.TypeSystemAlert,
		Channel:     "in_app",
		TriggeredBy: &senderID,
		SendEmail:   true,
	})

	return msg, nil
}

func (s *NotificationService) GetUserMessages(ctx context.Context, userID string, page, limit int) ([]*model.Message, error) {
	if limit <= 0 {
		limit = 20
	}
	if page <= 0 {
		page = 1
	}
	offset := (page - 1) * limit
	return s.notificationRepo.GetUserMessages(ctx, userID, limit, offset)
}
