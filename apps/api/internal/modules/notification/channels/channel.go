package channels

import (
	"context"

	identityModel "github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/notification/model"
)

type ChannelAdapter interface {
	ChannelType() string
	Send(ctx context.Context, notif *model.Notification, u *identityModel.User) error
}
