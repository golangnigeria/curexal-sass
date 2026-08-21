package channels

import (
	"context"
	"fmt"
	"sync"

	identityModel "github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/notification/model"
	"github.com/golangnigeria/curexal/internal/kernel/server"
)

type NotificationDispatcher struct {
	server   *server.Server
	adapters map[string]ChannelAdapter
	mu       sync.RWMutex
}

func NewNotificationDispatcher(s *server.Server) *NotificationDispatcher {
	d := &NotificationDispatcher{
		server:   s,
		adapters: make(map[string]ChannelAdapter),
	}
	d.Register(NewInAppChannelAdapter())
	d.Register(NewEmailChannelAdapter(s))
	return d
}

func (d *NotificationDispatcher) Register(adapter ChannelAdapter) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.adapters[adapter.ChannelType()] = adapter
}

func (d *NotificationDispatcher) Dispatch(ctx context.Context, notif *model.Notification, u *identityModel.User, prefs *model.NotificationPreference) error {
	d.mu.RLock()
	adapter, exists := d.adapters[notif.Channel]
	d.mu.RUnlock()

	if !exists {
		return fmt.Errorf("unsupported notification channel adapter: %s", notif.Channel)
	}

	// Respect user channel preferences
	if prefs != nil {
		if notif.Channel == "email" && !prefs.EmailEnabled {
			d.server.Logger.Info().Str("user_id", u.ID).Msg("Email channel disabled by user preference; skipping delivery")
			return nil
		}
		if notif.Channel == "sms" && !prefs.SMSEnabled {
			d.server.Logger.Info().Str("user_id", u.ID).Msg("SMS channel disabled by user preference; skipping delivery")
			return nil
		}
	}

	return adapter.Send(ctx, notif, u)
}
