package api

import (
	"errors"
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/catalogs/application"
	"github.com/golangnigeria/curexal/internal/modules/catalogs/domain"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/golangnigeria/curexal/internal/shared/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type MasterCatalogHandler struct {
	catalogService *application.MasterCatalogService
}

func NewMasterCatalogHandler(service *application.MasterCatalogService) *MasterCatalogHandler {
	return &MasterCatalogHandler{
		catalogService: service,
	}
}

func parseDomainParam(d string) (domain.CatalogDomain, error) {
	switch d {
	case "clinical":
		return domain.ClinicalDomain, nil
	case "lab", "laboratory":
		return domain.LabDomain, nil
	case "radiology", "rad":
		return domain.RadiologyDomain, nil
	case "pharmacy", "phm":
		return domain.PharmacyDomain, nil
	default:
		return "", domain.ErrInvalidCatalogDomain
	}
}

func (h *MasterCatalogHandler) ListItems(c echo.Context) error {
	domainParam := c.Param("domain")
	catDomain, err := parseDomainParam(domainParam)
	if err != nil {
		return response.BadRequestEcho(c, err.Error())
	}

	category := c.QueryParam("category")
	activeOnly := c.QueryParam("active_only") == "true" || c.QueryParam("active_only") == "1"

	items, err := h.catalogService.ListItems(c.Request().Context(), catDomain, category, activeOnly)
	if err != nil {
		return response.InternalErrorEcho(c, "Failed to list master catalog items: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, items)
}

func (h *MasterCatalogHandler) SearchItems(c echo.Context) error {
	domainParam := c.Param("domain")
	catDomain, err := parseDomainParam(domainParam)
	if err != nil {
		return response.BadRequestEcho(c, err.Error())
	}

	query := c.QueryParam("q")
	items, err := h.catalogService.SearchItems(c.Request().Context(), catDomain, query)
	if err != nil {
		return response.InternalErrorEcho(c, "Failed to search master catalog items: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, items)
}

func (h *MasterCatalogHandler) CreateItem(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	domainParam := c.Param("domain")
	catDomain, err := parseDomainParam(domainParam)
	if err != nil {
		return response.BadRequestEcho(c, err.Error())
	}

	var payload domain.CatalogItem
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid request payload")
	}
	payload.Domain = catDomain

	created, err := h.catalogService.CreateItem(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedPlatformAdmin) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrInvalidCatalogItem) || errors.Is(err, domain.ErrInvalidCatalogDomain) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusCreated, created)
}

func (h *MasterCatalogHandler) UpdateItem(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	domainParam := c.Param("domain")
	catDomain, err := parseDomainParam(domainParam)
	if err != nil {
		return response.BadRequestEcho(c, err.Error())
	}

	idParam := c.Param("id")
	itemUUID, err := uuid.Parse(idParam)
	if err != nil {
		return response.BadRequestEcho(c, "Invalid catalog item ID format")
	}

	var payload domain.UpdateCatalogItemPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid request payload")
	}

	updated, err := h.catalogService.UpdateItem(c.Request().Context(), principal, itemUUID, catDomain, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedPlatformAdmin) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrCatalogItemNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrOptimisticLockingConflict) {
			return response.ErrorEcho(c, http.StatusConflict, "CONFLICT", err.Error())
		}
		if errors.Is(err, domain.ErrInvalidCatalogItem) || errors.Is(err, domain.ErrInvalidCatalogDomain) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, updated)
}

// Baseline legacy DTO compatibility handlers
func (h *MasterCatalogHandler) GetMasterCatalogs(c echo.Context) error {
	ctx := c.Request().Context()
	specimens, _ := h.catalogService.GetSpecimenTypes(ctx)
	categories, _ := h.catalogService.GetTestCategories(ctx)
	units, _ := h.catalogService.GetUnitsOfMeasure(ctx)
	specialties, _ := h.catalogService.GetSpecialties(ctx)

	return c.JSON(http.StatusOK, domain.CatalogDataResponse{
		SpecimenTypes:  specimens,
		TestCategories: categories,
		UnitsOfMeasure: units,
		Specialties:    specialties,
	})
}

func (h *MasterCatalogHandler) SearchICD10(c echo.Context) error {
	ctx := c.Request().Context()
	q := c.QueryParam("q")
	results, err := h.catalogService.SearchICD10(ctx, q)
	if err != nil {
		return response.InternalErrorEcho(c, "Failed to query ICD-10 codes: "+err.Error())
	}
	return c.JSON(http.StatusOK, results)
}
