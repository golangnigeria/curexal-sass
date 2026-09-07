package service_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/identity/service"
	"github.com/golangnigeria/curexal/internal/shared/config"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoginResolution_TokenFormat(t *testing.T) {
	token, err := service.GenerateExchangeToken()
	require.NoError(t, err)
	assert.Len(t, token, 64, "32-byte hex token must be exactly 64 characters")

	hash1 := service.HashToken(token)
	hash2 := service.HashToken(token)
	assert.Equal(t, hash1, hash2, "SHA-256 hash must be deterministic")
	assert.Len(t, hash1, 64, "SHA-256 hex string must be 64 characters")
}

func TestLoginResolution_HostResolution(t *testing.T) {
	svc := service.NewLoginResolutionService(nil, nil, nil)

	// Clean host extraction
	clean, port := svc.ExtractCleanHost("app.localhost:5002")
	assert.Equal(t, "app.localhost", clean)
	assert.Equal(t, ":5002", port)

	clean2, port2 := svc.ExtractCleanHost("curexal.space")
	assert.Equal(t, "curexal.space", clean2)
	assert.Equal(t, "", port2)

	// Platform host resolution
	assert.Equal(t, "app.localhost:5002", svc.ResolvePlatformHost("curexal-clinic.localhost:5002"))
	assert.Equal(t, "app.curexal.space", svc.ResolvePlatformHost("curexal-clinic.curexal.space"))

	// Target workspace host resolution
	assert.Equal(t, "curexal-clinic.localhost:5002", svc.ResolveTargetHost("app.localhost:5002", "curexal-clinic", ""))
	assert.Equal(t, "curexal-clinic.curexal.space", svc.ResolveTargetHost("app.curexal.space", "curexal-clinic", ""))
	assert.Equal(t, "portal.customclinic.org", svc.ResolveTargetHost("app.curexal.space", "curexal-clinic", "portal.customclinic.org"))

	// Target URL construction
	url1 := svc.BuildTargetURL("app.localhost:5002", "curexal-clinic.localhost:5002", "/auth/exchange?token=abc123")
	assert.Equal(t, "http://curexal-clinic.localhost:5002/auth/exchange?token=abc123", url1)

	url2 := svc.BuildTargetURL("app.curexal.space", "curexal-clinic.curexal.space", "/auth/exchange?token=abc123")
	assert.Equal(t, "https://curexal-clinic.curexal.space/auth/exchange?token=abc123", url2)
}

func TestLoginResolution_Integration(t *testing.T) {
	cfg, err := config.LoadConfig()
	if err != nil || cfg.Database.DSN() == "" {
		t.Skip("Skipping DB integration test; Database DSN is not set")
	}

	logger := zerolog.Nop()
	srv, err := server.New(cfg, &logger, nil)
	if err != nil {
		t.Skipf("Skipping integration test; failed to connect to database/redis: %v", err)
	}

	authSvc := service.NewAuthService(srv)
	resSvc := service.NewLoginResolutionService(srv, authSvc, nil)
	ctx := context.Background()

	t.Run("Single-Use Atomic Exchange Token Redemption", func(t *testing.T) {
		payload := &service.ExchangePayload{
			UserID:          uuid.New().String(),
			OrganizationID:  uuid.New().String(),
			BranchID:        uuid.New().String(),
			Role:            "doctor",
			TargetHost:      "curexal-clinic.localhost:5002",
			DestinationPath: "/workspace/dashboard",
			CreatedAt:       time.Now().Unix(),
		}

		rawToken, err := resSvc.CreateExchangeToken(ctx, payload)
		require.NoError(t, err)
		assert.NotEmpty(t, rawToken)

		// 1st redemption must succeed
		redeemed, err := resSvc.RedeemExchangeToken(ctx, rawToken)
		require.NoError(t, err)
		assert.Equal(t, payload.UserID, redeemed.UserID)
		assert.Equal(t, payload.OrganizationID, redeemed.OrganizationID)
		assert.Equal(t, payload.BranchID, redeemed.BranchID)
		assert.Equal(t, payload.TargetHost, redeemed.TargetHost)

		// 2nd redemption MUST fail (atomic GETDEL single-use guarantee)
		secondRedeem, err := resSvc.RedeemExchangeToken(ctx, rawToken)
		assert.Error(t, err, "Redeeming the same token twice must fail")
		assert.Nil(t, secondRedeem)
		assert.Contains(t, err.Error(), "already used")
	})

	t.Run("Platform Admin Destination Resolution", func(t *testing.T) {
		superAdmin := &model.User{
			ID:              uuid.New().String(),
			Email:           "superadmin@curexal.internal",
			IsPlatformAdmin: true,
		}

		// On platform host -> Authenticated directly
		res, err := resSvc.ResolveDestination(ctx, superAdmin, "app.localhost:5002", nil, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, service.StatusAuthenticated, res.Status)
		assert.Equal(t, "/platform/audit", res.DestinationPath)

		// On clinic host -> Redirect required to app.localhost:5002
		resRedirect, err := resSvc.ResolveDestination(ctx, superAdmin, "curexal-clinic.localhost:5002", nil, nil, nil)
		require.NoError(t, err)
		assert.Equal(t, service.StatusRedirectRequired, resRedirect.Status)
		assert.Equal(t, "app.localhost:5002", resRedirect.TargetHost)
		assert.True(t, strings.HasPrefix(resRedirect.TargetURL, "http://app.localhost:5002/auth/exchange?token="))
	})
}
