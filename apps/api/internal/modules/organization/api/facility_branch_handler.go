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

type FacilityBranchHandler struct {
	branchService *application.FacilityBranchService
}

func NewFacilityBranchHandler(branchService *application.FacilityBranchService) *FacilityBranchHandler {
	return &FacilityBranchHandler{
		branchService: branchService,
	}
}

func (h *FacilityBranchHandler) ListBranches(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	branches, err := h.branchService.ListBranches(c.Request().Context(), principal)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to list facility branches: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, branches)
}

func (h *FacilityBranchHandler) GetBranch(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	branchIDParam := c.Param("id")
	branchUUID, errParse := uuid.Parse(branchIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid branch ID format")
	}

	branch, err := h.branchService.GetBranchByID(c.Request().Context(), principal, branchUUID)
	if err != nil {
		if errors.Is(err, domain.ErrFacilityBranchNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, branch)
}

func (h *FacilityBranchHandler) CreateBranch(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.CreateFacilityBranchPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid branch payload")
	}

	created, err := h.branchService.CreateBranch(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrHeadquartersConflict) || errors.Is(err, domain.ErrDuplicateBranchCode) {
			return response.ErrorEcho(c, http.StatusConflict, "CONFLICT", err.Error())
		}
		if errors.Is(err, domain.ErrMaxBranchesExceeded) {
			return response.ErrorEcho(c, http.StatusPaymentRequired, "BRANCH_LIMIT_EXCEEDED", err.Error())
		}
		if errors.Is(err, domain.ErrInactiveFacilityType) || errors.Is(err, domain.ErrInvalidFacilityBranch) {
			return response.BadRequestEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to create facility branch: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusCreated, created)
}

func (h *FacilityBranchHandler) UpdateBranch(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	branchIDParam := c.Param("id")
	branchUUID, errParse := uuid.Parse(branchIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid branch ID format")
	}

	var payload domain.UpdateFacilityBranchPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid branch update payload")
	}

	updated, err := h.branchService.UpdateBranch(c.Request().Context(), principal, branchUUID, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrHeadquartersConflict) || errors.Is(err, domain.ErrOptimisticLockingConflict) {
			return response.ErrorEcho(c, http.StatusConflict, "CONFLICT", err.Error())
		}
		if errors.Is(err, domain.ErrFacilityBranchNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to update facility branch: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, updated)
}

func (h *FacilityBranchHandler) DeactivateBranch(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	branchIDParam := c.Param("id")
	branchUUID, errParse := uuid.Parse(branchIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid branch ID format")
	}

	err := h.branchService.DeactivateBranch(c.Request().Context(), principal, branchUUID)
	if err != nil {
		if errors.Is(err, domain.ErrFacilityBranchNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]string{"message": "Facility branch deactivated successfully"})
}
