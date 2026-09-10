package handler

import (
	"net/http"
	"strings"
	"time"

	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/model"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/service"
	"github.com/labstack/echo/v4"
)

type ProviderProfileHandler struct {
	server  *server.Server
	service *service.ProviderProfileService
}

func NewProviderProfileHandler(s *server.Server, svc *service.ProviderProfileService) *ProviderProfileHandler {
	return &ProviderProfileHandler{
		server:  s,
		service: svc,
	}
}

// resolveTenantID extracts active branch/tenant ID from context or header
func (h *ProviderProfileHandler) resolveTenantID(c echo.Context) string {
	if tid, ok := c.Get("tenant_id").(string); ok && tid != "" {
		return tid
	}
	if bid, ok := c.Get("branch_id").(string); ok && bid != "" {
		return bid
	}
	if tid := c.Request().Header.Get("X-Tenant-ID"); tid != "" && tid != "undefined" {
		return tid
	}
	if bid := c.Request().Header.Get("X-Branch-ID"); bid != "" && bid != "undefined" {
		return bid
	}
	if p := platformAuth.GetPrincipal(c); p != nil {
		if p.TenantID != "" {
			return p.TenantID
		}
		if p.ActiveBranchID != "" {
			return p.ActiveBranchID
		}
	}
	return ""
}

// ListProviders handles GET /api/v1/providers/profiles
func (h *ProviderProfileHandler) ListProviders(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Active facility branch tenant context is required")
	}
	var statusPtr, specialtyPtr *string
	if s := c.QueryParam("status"); strings.TrimSpace(s) != "" {
		statusPtr = &s
	}
	if sp := c.QueryParam("specialty"); strings.TrimSpace(sp) != "" {
		specialtyPtr = &sp
	}

	providers, err := h.service.ListProviders(c.Request().Context(), tenantID, statusPtr, specialtyPtr)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to query provider profiles")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": providers,
		"meta": map[string]interface{}{
			"total":     len(providers),
			"timestamp": time.Now().UTC().Format(time.RFC3339),
		},
	})
}

// GetProviderByID handles GET /api/v1/providers/profiles/:id
func (h *ProviderProfileHandler) GetProviderByID(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Active facility branch tenant context is required")
	}
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Provider ID is required")
	}

	provider, err := h.service.GetProviderByID(c.Request().Context(), tenantID, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to retrieve provider profile")
	}
	if provider == nil {
		return echo.NewHTTPError(http.StatusNotFound, "Provider profile not found")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": provider,
	})
}

// CreateProvider handles POST /api/v1/providers/profiles
func (h *ProviderProfileHandler) CreateProvider(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Active facility branch tenant context is required")
	}
	var req model.CreateProviderProfileRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid provider profile payload")
	}

	profile, err := h.service.CreateOrLinkProviderProfile(c.Request().Context(), tenantID, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"data":    profile,
		"message": "Provider profile created successfully",
	})
}

// UpdateStatus handles PUT /api/v1/providers/profiles/:id/status
func (h *ProviderProfileHandler) UpdateStatus(c echo.Context) error {
	tenantID := h.resolveTenantID(c)
	if tenantID == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Active facility branch tenant context is required")
	}
	id := c.Param("id")
	if strings.TrimSpace(id) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "Provider ID is required")
	}

	var req model.UpdateProviderStatusRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid provider status payload")
	}

	updatedAt, err := h.service.UpdateProviderStatus(c.Request().Context(), tenantID, id, req)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"id":        id,
			"status":    req.Status,
			"updatedAt": updatedAt.UTC().Format(time.RFC3339),
		},
	})
}
