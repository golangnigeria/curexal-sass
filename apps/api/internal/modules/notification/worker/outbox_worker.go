package worker

import (
	"context"
	"fmt"
	"time"

	identityRepo "github.com/golangnigeria/curexal/internal/modules/identity/repository"
	"github.com/golangnigeria/curexal/internal/modules/notification/channels"
	"github.com/golangnigeria/curexal/internal/modules/notification/model"
	notificationRepo "github.com/golangnigeria/curexal/internal/modules/notification/repository"
	"github.com/golangnigeria/curexal/internal/kernel/server"
)

type OutboxWorker struct {
	server           *server.Server
	notificationRepo *notificationRepo.NotificationRepository
	userRepo         *identityRepo.UserRepository
	dispatcher       *channels.NotificationDispatcher
	stopChan         chan struct{}
}

func NewOutboxWorker(
	s *server.Server,
	notifRepo *notificationRepo.NotificationRepository,
	userRepo *identityRepo.UserRepository,
	dispatcher *channels.NotificationDispatcher,
) *OutboxWorker {
	return &OutboxWorker{
		server:           s,
		notificationRepo: notifRepo,
		userRepo:         userRepo,
		dispatcher:       dispatcher,
		stopChan:         make(chan struct{}),
	}
}

func (w *OutboxWorker) Start() {
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				w.ProcessPendingOutboxEvents()
			case <-w.stopChan:
				ticker.Stop()
				return
			}
		}
	}()
	w.server.Logger.Info().Msg("Transactional Outbox Worker started")
}

func (w *OutboxWorker) Stop() {
	close(w.stopChan)
}

func (w *OutboxWorker) ProcessPendingOutboxEvents() {
	if w.server == nil || w.server.DB == nil || w.notificationRepo == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	events, err := w.notificationRepo.FetchPendingOutboxEvents(ctx, 20)
	if err != nil || len(events) == 0 {
		return
	}

	for _, event := range events {
		w.processEvent(ctx, event)
	}
}

func (w *OutboxWorker) processEvent(ctx context.Context, event *model.OutboxEvent) {
	notifID := event.AggregateID
	notif, err := w.notificationRepo.GetNotificationByID(ctx, notifID)
	if err != nil || notif == nil {
		errMsg := fmt.Sprintf("notification %s not found", notifID)
		_ = w.notificationRepo.UpdateOutboxStatus(ctx, event.ID, "failed", event.Attempts+1, &errMsg, time.Now())
		return
	}

	user, err := w.userRepo.GetByID(ctx, notif.UserID)
	if err != nil || user == nil {
		errMsg := fmt.Sprintf("user %s not found for notification %s", notif.UserID, notifID)
		_ = w.notificationRepo.UpdateOutboxStatus(ctx, event.ID, "failed", event.Attempts+1, &errMsg, time.Now())
		return
	}

	prefs, _ := w.notificationRepo.GetNotificationPreference(ctx, notif.UserID)

	// Dispatch notification via configured channel adapter
	dispatchErr := w.dispatcher.Dispatch(ctx, notif, user, prefs)

	if dispatchErr != nil {
		attempts := event.Attempts + 1
		errMsg := dispatchErr.Error()
		if attempts >= event.MaxAttempts {
			_ = w.notificationRepo.UpdateOutboxStatus(ctx, event.ID, "failed", attempts, &errMsg, time.Now())
			_ = w.notificationRepo.UpdateNotificationStatus(ctx, notifID, model.StatusFailed)
		} else {
			// Exponential backoff retry: 5s, 15s, 45s, 135s
			retryDelay := time.Duration(1<<attempts*5) * time.Second
			nextRetry := time.Now().Add(retryDelay)
			_ = w.notificationRepo.UpdateOutboxStatus(ctx, event.ID, "pending", attempts, &errMsg, nextRetry)
		}
	} else {
		_ = w.notificationRepo.UpdateOutboxStatus(ctx, event.ID, "completed", event.Attempts+1, nil, time.Now())
		_ = w.notificationRepo.UpdateNotificationStatus(ctx, notifID, model.StatusSent)
	}
}
