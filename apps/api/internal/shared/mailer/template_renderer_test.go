package mailer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRenderTemplate_Defaults(t *testing.T) {
	// Ensure env vars are unset for default test
	t.Setenv("EMAIL_LOGO_URL", "")
	t.Setenv("CUREXAL_INTEGRATION_EMAIL_LOGO_URL", "")
	t.Setenv("EMAIL_APP_URL", "")
	t.Setenv("CUREXAL_INTEGRATION_EMAIL_APP_URL", "")

	data := TemplateData{
		UserName:         "Dr. Jane Doe",
		UserFirstName:    "Jane",
		VerificationCode: "654321",
	}

	html, err := RenderTemplate("verify-email", data)
	require.NoError(t, err)
	assert.NotEmpty(t, html)

	// Verify canonical default URLs
	assert.Contains(t, html, "https://cdn.curexal.space/email/full_logo.png")
	assert.Contains(t, html, "https://curexal.space")
}

func TestRenderTemplate_EnvironmentOverrides(t *testing.T) {
	customLogo := "https://cdn.curexal.space/email/custom-logo-v2.png"
	customApp := "https://curexal.space/app"

	t.Setenv("EMAIL_LOGO_URL", customLogo)
	t.Setenv("EMAIL_APP_URL", customApp)

	data := TemplateData{
		UserName:      "John Doe",
		UserFirstName: "John",
	}

	html, err := RenderTemplate("verify-email", data)
	require.NoError(t, err)
	assert.Contains(t, html, customLogo)
	assert.Contains(t, html, customApp)
}

func TestRenderTemplate_ExplicitDataFields(t *testing.T) {
	explicitLogo := "https://cdn.curexal.space/email/explicit-logo.png"
	explicitApp := "https://portal.curexal.space"

	data := TemplateData{
		UserName:      "Alice Smith",
		UserFirstName: "Alice",
		LogoURL:       explicitLogo,
		AppURL:        explicitApp,
	}

	html, err := RenderTemplate("verify-email", data)
	require.NoError(t, err)
	assert.Contains(t, html, explicitLogo)
	assert.Contains(t, html, explicitApp)
}

func TestRenderTemplate_LogoHyperlinkStructure(t *testing.T) {
	data := TemplateData{
		UserName:         "Test User",
		UserFirstName:    "Test",
		VerificationCode: "112233",
	}

	html, err := RenderTemplate("verify-email", data)
	require.NoError(t, err)

	// Ensure header logo is wrapped inside hyperlinked anchor tag
	assert.True(t, strings.Contains(html, "<a href=") && strings.Contains(html, "https://curexal.space"),
		"Expected email template to contain hyperlinked logo anchor")
}

func TestRenderTemplate_PasswordDelivery(t *testing.T) {
	data := TemplateData{
		UserName:          "Dr. Alex Carter",
		UserFirstName:     "Alex",
		GeneratedPassword: "P@ssw0rd!2026",
		LoginURL:          "https://curexal.space/auth/sign-in",
	}

	html, err := RenderTemplate("password-delivery", data)
	require.NoError(t, err)
	assert.NotEmpty(t, html)
	assert.Contains(t, html, "Alex")
	assert.Contains(t, html, "P@ssw0rd!2026")
	assert.Contains(t, html, "https://curexal.space/auth/sign-in")
}

