package api

import (
	"errors"
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/organization/application"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/golangnigeria/curexal/internal/shared/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrganizationCatalogHandler struct {
	catalogService *application.OrganizationCatalogService
}

func NewOrganizationCatalogHandler(catalogService *application.OrganizationCatalogService) *OrganizationCatalogHandler {
	return &OrganizationCatalogHandler{
		catalogService: catalogService,
	}
}

func (h *OrganizationCatalogHandler) ListCatalogItems(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	domainType := c.QueryParam("domainType")
	items, err := h.catalogService.ListCatalogItems(c.Request().Context(), principal, domainType)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to list organization catalog items: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, items)
}

func (h *OrganizationCatalogHandler) GetCatalogItem(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	itemIDParam := c.Param("id")
	itemUUID, errParse := uuid.Parse(itemIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid catalog item ID format")
	}

	item, err := h.catalogService.GetCatalogItemByID(c.Request().Context(), principal, itemUUID)
	if err != nil {
		if errors.Is(err, domain.ErrCatalogItemNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, item)
}

func (h *OrganizationCatalogHandler) CreateCatalogItem(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.CreateCatalogItemPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid catalog item payload")
	}

	created, err := h.catalogService.CreateCatalogItem(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrInvalidCatalogDomain) {
			return response.BadRequestEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrDuplicateCatalogCode) {
			return response.ErrorEcho(c, http.StatusConflict, "DUPLICATE_CODE", err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to create catalog item: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusCreated, created)
}

func (h *OrganizationCatalogHandler) UpdateCatalogItem(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	itemIDParam := c.Param("id")
	itemUUID, errParse := uuid.Parse(itemIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid catalog item ID format")
	}

	var payload domain.UpdateCatalogItemPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid update payload")
	}

	updated, err := h.catalogService.UpdateCatalogItem(c.Request().Context(), principal, itemUUID, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrCatalogItemNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to update catalog item: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, updated)
}

func (h *OrganizationCatalogHandler) SetBranchPriceOverride(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	itemIDParam := c.Param("id")
	itemUUID, errParse := uuid.Parse(itemIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid catalog item ID format")
	}

	var payload domain.SetBranchPricePayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid branch price payload")
	}

	res, err := h.catalogService.SetBranchPriceOverride(c.Request().Context(), principal, itemUUID, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrFacilityBranchNotFound) || errors.Is(err, domain.ErrCatalogItemNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to set branch price override: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, res)
}

func (h *OrganizationCatalogHandler) ListInsuranceProviders(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	providers, err := h.catalogService.ListInsuranceProviders(c.Request().Context(), principal)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to list insurance providers: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, providers)
}

func (h *OrganizationCatalogHandler) CreateInsuranceProvider(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.CreateInsuranceProviderPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid insurance provider payload")
	}

	created, err := h.catalogService.CreateInsuranceProvider(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrDuplicateInsuranceCode) {
			return response.ErrorEcho(c, http.StatusConflict, "DUPLICATE_CODE", err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to create insurance provider: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusCreated, created)
}
