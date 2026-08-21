package notification

import (
	identityRepo "github.com/golangnigeria/curexal/internal/modules/identity/repository"
	"github.com/golangnigeria/curexal/internal/modules/notification/channels"
	"github.com/golangnigeria/curexal/internal/modules/notification/handler"
	"github.com/golangnigeria/curexal/internal/modules/notification/repository"
	"github.com/golangnigeria/curexal/internal/modules/notification/service"
	"github.com/golangnigeria/curexal/internal/modules/notification/worker"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
)

type Module struct {
	Service    *service.NotificationService
	Handler    *handler.NotificationHandler
	Dispatcher *channels.NotificationDispatcher
	Worker     *worker.OutboxWorker
}

func NewModule(s *server.Server, userRepo *identityRepo.UserRepository) *Module {
	repo := repository.NewNotificationRepository(s)
	svc := service.NewNotificationService(s, repo, userRepo)
	hnd := handler.NewNotificationHandler(svc)
	dispatcher := channels.NewNotificationDispatcher(s)
	outboxWorker := worker.NewOutboxWorker(s, repo, userRepo, dispatcher)

	// Start asynchronous outbox worker for reliable delivery
	outboxWorker.Start()

	return &Module{
		Service:    svc,
		Handler:    hnd,
		Dispatcher: dispatcher,
		Worker:     outboxWorker,
	}
}

func (m *Module) RegisterRoutes(e *echo.Echo, authMiddleware echo.MiddlewareFunc) {
	v1 := e.Group("/api/v1", authMiddleware)

	// Notifications API
	notifs := v1.Group("/notifications")
	notifs.GET("", m.Handler.GetNotifications)
	notifs.POST("/:id/read", m.Handler.MarkAsRead)
	notifs.POST("/read-all", m.Handler.MarkAllAsRead)
	notifs.GET("/unread-count", m.Handler.GetUnreadCount)
	notifs.POST("", m.Handler.CreateNotification)
	notifs.GET("/preferences", m.Handler.GetPreferences)
	notifs.PUT("/preferences", m.Handler.UpdatePreferences)

	// Direct Messaging API
	msgs := v1.Group("/messages")
	msgs.GET("", m.Handler.GetMessages)
	msgs.POST("", m.Handler.SendMessage)
}
