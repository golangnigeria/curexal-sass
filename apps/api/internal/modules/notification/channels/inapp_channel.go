package channels

import (
	"context"

	identityModel "github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/notification/model"
)

type InAppChannelAdapter struct{}

func NewInAppChannelAdapter() *InAppChannelAdapter {
	return &InAppChannelAdapter{}
}

func (a *InAppChannelAdapter) ChannelType() string {
	return "in_app"
}

func (a *InAppChannelAdapter) Send(ctx context.Context, notif *model.Notification, u *identityModel.User) error {
	// In-App notifications are persisted directly to database in SQL transaction
	return nil
}
