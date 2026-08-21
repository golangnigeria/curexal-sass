package mailer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

type Mailer struct {
	apiKey      string
	fromName    string
	fromAddress string
	logger      *zerolog.Logger
	httpClient  *http.Client
}

func New(apiKey, fromName, fromAddress string, logger *zerolog.Logger) *Mailer {
	if apiKey == "" {
		apiKey = ""
	}
	if fromName == "" {
		fromName = "Curexal"
	}
	if fromAddress == "" {
		fromAddress = "noreply@contact.curexal.space"
	}
	return &Mailer{
		apiKey:      apiKey,
		fromName:    fromName,
		fromAddress: fromAddress,
		logger:      logger,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
	}
}

type resendSendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	HTML    string   `json:"html"`
}

type resendSendResponse struct {
	ID    string `json:"id"`
	Error *struct {
		Message string `json:"message"`
		Name    string `json:"name"`
	} `json:"error,omitempty"`
}

func (m *Mailer) SendEmail(ctx context.Context, toEmail, subject, htmlBody string) error {
	if m.apiKey == "" {
		if m.logger != nil {
			m.logger.Warn().Str("to", toEmail).Msg("Resend API key missing; skipping email dispatch")
		}
		return nil
	}

	fromHeader := fmt.Sprintf("%s <%s>", m.fromName, m.fromAddress)
	reqBody := resendSendRequest{
		From:    fromHeader,
		To:      []string{toEmail},
		Subject: subject,
		HTML:    htmlBody,
	}

	jsonBytes, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("failed to marshal resend email payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return fmt.Errorf("failed to create resend request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+m.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		if m.logger != nil {
			m.logger.Error().Err(err).Str("to", toEmail).Msg("HTTP request to Resend failed")
		}
		return fmt.Errorf("failed to send email via resend: %w", err)
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		if m.logger != nil {
			m.logger.Error().
				Int("status", resp.StatusCode).
				Str("body", string(respBytes)).
				Str("to", toEmail).
				Msg("Resend API returned error")
		}
		return fmt.Errorf("resend API error (%d): %s", resp.StatusCode, string(respBytes))
	}

	var resendResp resendSendResponse
	_ = json.Unmarshal(respBytes, &resendResp)

	if m.logger != nil {
		m.logger.Info().
			Str("to", toEmail).
			Str("resend_id", resendResp.ID).
			Msg("Verification email successfully delivered via Resend")
	}

	return nil
}

func (m *Mailer) SendVerificationEmail(ctx context.Context, toEmail, userName, code string) error {
	subject := "Your Curexal Verification Code: " + code

	firstName := userName
	if parts := strings.Split(userName, " "); len(parts) > 0 {
		firstName = parts[0]
	}

	htmlBody, err := RenderTemplate("verify-email", TemplateData{
		UserName:         userName,
		UserFirstName:    firstName,
		VerificationCode: code,
		Code:             code,
		Title:            "Verify Your Email Address",
	})
	if err != nil || htmlBody == "" {
		htmlBody = renderFallbackHTML(TemplateData{
			UserName:         userName,
			UserFirstName:    firstName,
			VerificationCode: code,
			Code:             code,
			Title:            "Verify Your Email Address",
		})
	}

	return m.SendEmail(ctx, toEmail, subject, htmlBody)
}

func (m *Mailer) SendOwnerInvitationEmail(ctx context.Context, toEmail, userName, orgName, code string) error {
	subject := fmt.Sprintf("Welcome to Curexal — Owner Setup Code: %s", code)

	firstName := userName
	if parts := strings.Split(userName, " "); len(parts) > 0 {
		firstName = parts[0]
	}

	htmlBody, err := RenderTemplate("owner-invitation", TemplateData{
		UserName:         userName,
		UserFirstName:    firstName,
		OrgName:          orgName,
		VerificationCode: code,
		Code:             code,
		Title:            "Organization Owner Setup",
	})
	if err != nil || htmlBody == "" {
		htmlBody = renderFallbackHTML(TemplateData{
			UserName:         userName,
			UserFirstName:    firstName,
			OrgName:          orgName,
			VerificationCode: code,
			Code:             code,
			Title:            "Organization Owner Setup",
			Message:          fmt.Sprintf("You have been designated as the owner of %s on Curexal. Your verification code is: %s", orgName, code),
		})
	}

	return m.SendEmail(ctx, toEmail, subject, htmlBody)
}

func (m *Mailer) SendPasswordDeliveryEmail(ctx context.Context, toEmail, userName, generatedPassword, loginURL string) error {
	subject := "Your Curexal Account Password"

	firstName := userName
	if parts := strings.Split(userName, " "); len(parts) > 0 {
		firstName = parts[0]
	}

	htmlBody, err := RenderTemplate("password-delivery", TemplateData{
		UserName:          userName,
		UserFirstName:     firstName,
		GeneratedPassword: generatedPassword,
		Password:          generatedPassword,
		LoginURL:          loginURL,
		ActionURL:         loginURL,
		Title:             "Your Account Password",
		Message:           "Your account password has been generated upon your request. Please find your credentials below.",
	})
	if err != nil || htmlBody == "" {
		htmlBody = renderFallbackHTML(TemplateData{
			UserName:          userName,
			UserFirstName:     firstName,
			GeneratedPassword: generatedPassword,
			Password:          generatedPassword,
			LoginURL:          loginURL,
			ActionURL:         loginURL,
			Title:             "Your Account Password",
			Message:           "Your account password is: " + generatedPassword + ". Please use it to log in and update your password immediately.",
		})
	}

	return m.SendEmail(ctx, toEmail, subject, htmlBody)
}
