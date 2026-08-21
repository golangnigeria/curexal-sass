package channels

import (
	"context"
	"fmt"

	identityModel "github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/notification/model"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/mailer"
)

type EmailChannelAdapter struct {
	server *server.Server
}

func NewEmailChannelAdapter(s *server.Server) *EmailChannelAdapter {
	return &EmailChannelAdapter{server: s}
}

func (a *EmailChannelAdapter) ChannelType() string {
	return "email"
}

func (a *EmailChannelAdapter) Send(ctx context.Context, notif *model.Notification, u *identityModel.User) error {
	if u == nil || u.Email == "" {
		return fmt.Errorf("recipient user email is missing")
	}
	if a.server == nil || a.server.Mailer == nil {
		return fmt.Errorf("mailer service not initialized")
	}

	actionURL := "http://patient.localhost:5001/login"
	if notif.LinkURL != nil && *notif.LinkURL != "" {
		actionURL = *notif.LinkURL
	}

	templateName := "notification"
	if notif.Type == string(model.TypePatientResultReady) {
		templateName = "patient-result"
	}

	data := mailer.TemplateData{
		UserName:       u.Name,
		PatientName:    u.Name,
		Title:          notif.Title,
		Message:        notif.Message,
		ActionURL:      actionURL,
		PortalURL:      actionURL,
		ActionText:     "View in Patient Vault",
		TestName:       "Diagnostic Result",
		LaboratoryName: "Curexal Partner Lab",
		OrderNumber:    "ORD-2026",
	}

	htmlBody, err := mailer.RenderTemplate(templateName, data)
	if err != nil {
		return fmt.Errorf("failed to render HTML template %s: %w", templateName, err)
	}

	return a.server.Mailer.SendEmail(ctx, u.Email, notif.Title, htmlBody)
}
