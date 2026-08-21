package api

import (
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/clinical/application"
	"github.com/golangnigeria/curexal/internal/modules/clinical/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type CatalogHandler struct {
	server     *server.Server
	appService *application.ClinicalApplicationService
}

func NewCatalogHandler(s *server.Server, appService *application.ClinicalApplicationService) *CatalogHandler {
	return &CatalogHandler{
		server:     s,
		appService: appService,
	}
}

func (h *CatalogHandler) ListCatalog(c echo.Context) error {
	ctx := c.Request().Context()
	items, err := h.appService.ListCatalog(ctx)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch catalog: "+err.Error())
	}
	return c.JSON(http.StatusOK, items)
}

func (h *CatalogHandler) CreateCatalogItem(c echo.Context) error {
	if !middleware.IsPlatformStaff(c) {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied: Platform Staff only")
	}

	var payload domain.CreateCatalogItemPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload: "+err.Error())
	}

	ctx := c.Request().Context()
	item, err := h.appService.CreateCatalogItem(ctx, &payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create catalog item: "+err.Error())
	}

	return c.JSON(http.StatusCreated, item)
}

func (h *CatalogHandler) ImportCatalog(c echo.Context) error {
	if !middleware.IsPlatformStaff(c) {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied: Platform Staff only")
	}

	var payload []domain.CreateCatalogItemPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload: "+err.Error())
	}

	ctx := c.Request().Context()
	count, err := h.appService.ImportCatalog(ctx, payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to import catalog items: "+err.Error())
	}

	return c.JSON(http.StatusCreated, map[string]any{
		"success": true,
		"count":   count,
	})
}

func (h *CatalogHandler) UpdateCatalogMetadata(c echo.Context) error {
	if !middleware.IsPlatformStaff(c) {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied: Platform Staff only")
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	ctx := c.Request().Context()
	existing, err := h.appService.GetCatalogItemByID(ctx, id)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Catalog item not found")
	}

	var patchPayload domain.CreateCatalogItemPayload
	if err := c.Bind(&patchPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload: "+err.Error())
	}

	if patchPayload.Name == "" {
		patchPayload.Name = existing.Name
	}
	if patchPayload.Code == "" {
		patchPayload.Code = existing.Code
	}
	if patchPayload.Type == "" {
		patchPayload.Type = existing.Type
	}
	if patchPayload.DisplayName == nil || *patchPayload.DisplayName == "" {
		patchPayload.DisplayName = existing.DisplayName
	}
	if patchPayload.ShortName == nil || *patchPayload.ShortName == "" {
		patchPayload.ShortName = existing.ShortName
	}
	if patchPayload.RecoveryTime == nil || *patchPayload.RecoveryTime == "" {
		patchPayload.RecoveryTime = existing.RecoveryTime
	}
	if patchPayload.DepartmentID == nil {
		patchPayload.DepartmentID = existing.DepartmentID
	}
	if patchPayload.TestGroup == nil || *patchPayload.TestGroup == "" {
		patchPayload.TestGroup = existing.TestGroup
	}
	if patchPayload.TestCategory == nil || *patchPayload.TestCategory == "" {
		patchPayload.TestCategory = existing.TestCategory
	}
	if patchPayload.TatHours == 0 {
		patchPayload.TatHours = existing.TatHours
	}
	if len(patchPayload.ChildServiceIDs) == 0 {
		patchPayload.ChildServiceIDs = existing.ChildServiceIDs
	}

	patchPayload.BasePrice = existing.BasePrice
	patchPayload.UrgencyPrice = existing.UrgencyPrice
	patchPayload.CommissionValue = existing.CommissionValue
	patchPayload.CommissionPercentage = existing.CommissionPercentage
	patchPayload.DiscountAmount = existing.DiscountAmount
	patchPayload.DiscountPercentage = existing.DiscountPercentage

	updated, err := h.appService.UpdateCatalogMetadata(ctx, id, &patchPayload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update catalog metadata: "+err.Error())
	}

	return c.JSON(http.StatusOK, updated)
}

func (h *CatalogHandler) UpdateCatalogPricing(c echo.Context) error {
	userRole, _ := c.Get(middleware.UserRoleKey).(string)
	isPlatformStaff := middleware.IsPlatformStaff(c)
	isBranchAdmin := userRole == "owner" || userRole == "admin" || userRole == "branch_admin" || userRole == "org_regional_manager" || userRole == "Administrator"

	if !isPlatformStaff && !isBranchAdmin {
		return echo.NewHTTPError(http.StatusForbidden, "Access denied: Platform Staff or Branch Admin only")
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid ID format")
	}

	var payload domain.UpdatePricingPayload
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request payload: "+err.Error())
	}

	ctx := c.Request().Context()
	updated, err := h.appService.UpdateCatalogPricing(ctx, id, &payload)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update catalog pricing: "+err.Error())
	}

	return c.JSON(http.StatusOK, updated)
}
