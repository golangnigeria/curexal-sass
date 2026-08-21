package mailer

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sync"
)

var (
	logoBase64Cache     string
	logoBase64CacheOnce sync.Once
)

func getLogoDataURI() string {
	logoBase64CacheOnce.Do(func() {
		possibleLogoPaths := []string{
			filepath.Join("static", "full_logo.png"),
			filepath.Join("static", "logo.png"),
			filepath.Join("apps", "backend", "static", "full_logo.png"),
			filepath.Join("..", "static", "full_logo.png"),
		}
		for _, p := range possibleLogoPaths {
			b, err := os.ReadFile(p)
			if err == nil && len(b) > 0 {
				logoBase64Cache = "data:image/png;base64," + base64.StdEncoding.EncodeToString(b)
				break
			}
		}
	})
	return logoBase64Cache
}

type TemplateTheme struct {
	BgClass     string
	BorderColor string
	TitleColor  string
	LabelText   string
}

type TemplateData struct {
	UserName         string
	UserFirstName    string
	Title            string
	Message          string
	ActionURL        string
	ActionText       string
	PatientName      string
	TestName         string
	LaboratoryName   string
	OrderNumber      string
	PortalURL        string
	VerificationCode string
	Code             string
	VerificationLink  string
	GeneratedPassword string
	Password          string
	LoginURL          string
	OrgName           string
	OrgType           string
	LogoURL           string
	AppURL            string
	Theme             TemplateTheme
}

// RenderTemplate renders an HTML email template from templates/emails/<templateName>.html with fallbacks
func RenderTemplate(templateName string, data TemplateData) (string, error) {
	if data.LogoURL == "" {
		if envLogo := os.Getenv("EMAIL_LOGO_URL"); envLogo != "" {
			data.LogoURL = envLogo
		} else if envLogo := os.Getenv("CUREXAL_INTEGRATION_EMAIL_LOGO_URL"); envLogo != "" {
			data.LogoURL = envLogo
		} else {
			data.LogoURL = "https://cdn.curexal.space/email/full_logo.png"
		}
	}
	if data.AppURL == "" {
		if envApp := os.Getenv("EMAIL_APP_URL"); envApp != "" {
			data.AppURL = envApp
		} else if envApp := os.Getenv("CUREXAL_INTEGRATION_EMAIL_APP_URL"); envApp != "" {
			data.AppURL = envApp
		} else {
			data.AppURL = "https://curexal.space"
		}
	}
	if data.Theme.BgClass == "" {
		data.Theme.BgClass = "rgb(240,253,250)"
		data.Theme.BorderColor = "rgb(15,118,110)"
		data.Theme.TitleColor = "rgb(15,118,110)"
		data.Theme.LabelText = "Curexal Platform"
	}
	if data.OrgName == "" {
		data.OrgName = "Curexal Workspace"
	}

	possiblePaths := []string{
		filepath.Join("templates", "emails", templateName+".html"),
		filepath.Join("apps", "backend", "templates", "emails", templateName+".html"),
		filepath.Join("..", "templates", "emails", templateName+".html"),
		filepath.Join("..", "..", "templates", "emails", templateName+".html"),
		filepath.Join("..", "..", "..", "templates", "emails", templateName+".html"),
	}

	var tmpl *template.Template
	var err error

	for _, path := range possiblePaths {
		if _, statErr := os.Stat(path); statErr == nil {
			tmpl, err = template.ParseFiles(path)
			if err == nil {
				break
			}
		}
	}

	if tmpl == nil {
		// Render dynamic default fallback template
		return renderFallbackHTML(data), nil
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute email template %s: %w", templateName, err)
	}

	return buf.String(), nil
}

func renderFallbackHTML(data TemplateData) string {
	title := data.Title
	if title == "" {
		title = "Notification"
	}
	userName := data.UserName
	if userName == "" {
		userName = data.PatientName
	}
	if userName == "" {
		userName = "User"
	}
	actionURL := data.ActionURL
	if actionURL == "" {
		actionURL = data.PortalURL
	}
	if actionURL == "" {
		if data.AppURL != "" {
			actionURL = data.AppURL
		} else {
			actionURL = "https://curexal.space"
		}
	}
	appURL := data.AppURL
	if appURL == "" {
		appURL = "https://curexal.space"
	}
	logoURL := data.LogoURL
	if logoURL == "" {
		logoURL = "https://cdn.curexal.space/email/full_logo.png"
	}

	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head><meta charset="utf-8"></head>
<body style="font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background-color: #f8fafc; margin: 0; padding: 40px 20px;">
  <div style="max-width: 540px; margin: 0 auto; background: #ffffff; border-radius: 16px; border: 1px solid #e2e8f0; padding: 40px;">
    <div style="margin-bottom: 24px;">
      <a href="%s" target="_blank" style="text-decoration:none;border:none;">
        <img src="%s" alt="Curexal" height="40" style="display:inline-block;vertical-align:middle;max-height:40px;width:auto;margin-right:8px;border:0;outline:none;" />
        <span style="color: #0F766E; font-weight: 800; font-size: 22px; vertical-align:middle;">Curexal Healthcare</span>
      </a>
    </div>
    <h1 style="font-size: 20px; font-weight: 700; color: #0f172a; margin-bottom: 12px;">%s</h1>
    <p style="font-size: 14px; color: #475569; line-height: 1.6; margin-bottom: 20px;">Hello %s,</p>
    <p style="font-size: 14px; color: #475569; line-height: 1.6; margin-bottom: 24px;">%s</p>
    <div style="text-align: center; margin: 32px 0;">
      <a href="%s" style="display: inline-block; background-color: #0F766E; color: #ffffff; font-weight: 600; font-size: 14px; padding: 12px 28px; border-radius: 10px; text-decoration: none;">View in Portal</a>
    </div>
    <div style="margin-top: 32px; font-size: 12px; color: #94a3b8; border-top: 1px solid #f1f5f9; padding-top: 20px;">
      This notification was generated automatically by Curexal.
    </div>
  </div>
</body>
</html>`, appURL, logoURL, title, userName, data.Message, actionURL)
}
