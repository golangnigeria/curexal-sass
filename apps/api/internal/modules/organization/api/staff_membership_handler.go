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

type StaffMembershipHandler struct {
	staffService *application.StaffMembershipService
}

func NewStaffMembershipHandler(staffService *application.StaffMembershipService) *StaffMembershipHandler {
	return &StaffMembershipHandler{
		staffService: staffService,
	}
}

func (h *StaffMembershipHandler) ListMembers(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	members, err := h.staffService.ListMembers(c.Request().Context(), principal)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to list staff members: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, members)
}

func (h *StaffMembershipHandler) CreateInvitation(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	var payload domain.CreateStaffInvitationPayload
	if err := c.Bind(&payload); err != nil {
		return response.BadRequestEcho(c, "Invalid staff invitation payload")
	}

	res, err := h.staffService.CreateInvitation(c.Request().Context(), principal, &payload)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrDuplicateStaffInvite) {
			return response.ErrorEcho(c, http.StatusConflict, "DUPLICATE_INVITATION", err.Error())
		}
		if errors.Is(err, domain.ErrMaxStaffExceeded) {
			return response.ErrorEcho(c, http.StatusPaymentRequired, "STAFF_LIMIT_EXCEEDED", err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to create staff invitation: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusCreated, res)
}

func (h *StaffMembershipHandler) ListInvitations(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	invitations, err := h.staffService.ListInvitations(c.Request().Context(), principal)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, "Failed to list staff invitations: "+err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, invitations)
}

func (h *StaffMembershipHandler) RevokeInvitation(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	inviteIDParam := c.Param("id")
	inviteUUID, errParse := uuid.Parse(inviteIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid invitation ID format")
	}

	err := h.staffService.RevokeInvitation(c.Request().Context(), principal, inviteUUID)
	if err != nil {
		if errors.Is(err, domain.ErrStaffInvitationNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]string{"message": "Staff invitation revoked successfully"})
}

type assignBranchRequest struct {
	BranchID string `json:"branchId"`
}

func (h *StaffMembershipHandler) AssignBranch(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	membershipIDParam := c.Param("id")
	memUUID, errParse := uuid.Parse(membershipIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid membership ID format")
	}

	var req assignBranchRequest
	if err := c.Bind(&req); err != nil || req.BranchID == "" {
		return response.BadRequestEcho(c, "branchId is required")
	}

	branchUUID, errBranchParse := uuid.Parse(req.BranchID)
	if errBranchParse != nil {
		return response.BadRequestEcho(c, "Invalid branch ID format")
	}

	res, err := h.staffService.AssignBranch(c.Request().Context(), principal, memUUID, branchUUID)
	if err != nil {
		if errors.Is(err, domain.ErrFacilityBranchNotFound) || errors.Is(err, domain.ErrStaffMemberNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusCreated, res)
}

type assignDepartmentRequest struct {
	BranchID       string `json:"branchId"`
	DepartmentCode string `json:"departmentCode"`
}

func (h *StaffMembershipHandler) AssignDepartment(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	membershipIDParam := c.Param("id")
	memUUID, errParse := uuid.Parse(membershipIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid membership ID format")
	}

	var req assignDepartmentRequest
	if err := c.Bind(&req); err != nil || req.BranchID == "" || req.DepartmentCode == "" {
		return response.BadRequestEcho(c, "branchId and departmentCode are required")
	}

	branchUUID, errBranchParse := uuid.Parse(req.BranchID)
	if errBranchParse != nil {
		return response.BadRequestEcho(c, "Invalid branch ID format")
	}

	res, err := h.staffService.AssignDepartment(c.Request().Context(), principal, memUUID, branchUUID, req.DepartmentCode)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidDepartmentCode) {
			return response.BadRequestEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrFacilityBranchNotFound) || errors.Is(err, domain.ErrStaffMemberNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusCreated, res)
}

type updateRoleRequest struct {
	Role      string `json:"role"`
	RoleTitle string `json:"roleTitle"`
}

func (h *StaffMembershipHandler) UpdateRole(c echo.Context) error {
	principal := middleware.GetPrincipal(c)
	if principal == nil {
		return response.UnauthorizedEcho(c, "Authentication required")
	}

	membershipIDParam := c.Param("id")
	memUUID, errParse := uuid.Parse(membershipIDParam)
	if errParse != nil {
		return response.BadRequestEcho(c, "Invalid membership ID format")
	}

	var req updateRoleRequest
	if err := c.Bind(&req); err != nil || req.Role == "" {
		return response.BadRequestEcho(c, "role is required")
	}

	roleTitleVal := "member"
	if req.RoleTitle != "" {
		roleTitleVal = req.RoleTitle
	}

	updated, err := h.staffService.UpdateMemberRole(c.Request().Context(), principal, memUUID, req.Role, roleTitleVal)
	if err != nil {
		if errors.Is(err, domain.ErrStaffMemberNotFound) {
			return response.NotFoundEcho(c, err.Error())
		}
		if errors.Is(err, domain.ErrUnauthorizedTenantAccess) {
			return response.ForbiddenEcho(c, err.Error())
		}
		return response.InternalErrorEcho(c, err.Error())
	}

	return response.SuccessEcho(c, http.StatusOK, updated)
}
