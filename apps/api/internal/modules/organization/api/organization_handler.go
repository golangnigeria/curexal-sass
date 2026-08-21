package api

import (
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/organization/application"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type OrganizationHandler struct {
	server     *server.Server
	appService *application.OrganizationApplicationService
}

func NewOrganizationHandler(s *server.Server, appService *application.OrganizationApplicationService) *OrganizationHandler {
	return &OrganizationHandler{
		server:     s,
		appService: appService,
	}
}

func (h *OrganizationHandler) CreateOrganization(c echo.Context) error {
	var payload CreateOrganizationPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("invalid request payload", false, nil, nil, nil)
	}
	if err := payload.Validate(); err != nil {
		return errs.NewBadRequestError(err.Error(), false, nil, nil, nil)
	}

	userID := middleware.GetUserID(c)
	ownerEmail := payload.OwnerEmail
	ownerName := payload.OwnerName
	if payload.Owner != nil {
		ownerEmail = &payload.Owner.Email
		ownerName = payload.Owner.Name
	}

	cmd := &application.CreateOrganizationCommand{
		Name:               payload.Name,
		Slug:               payload.Slug,
		Plan:               payload.Plan,
		Address:            payload.Address,
		City:               payload.City,
		State:              payload.State,
		LGA:                payload.LGA,
		Country:            payload.Country,
		Phone:              payload.Phone,
		Email:              payload.Email,
		RegistrationNumber: payload.RegistrationNumber,
		LicenseNumber:      payload.LicenseNumber,
		TaxID:              payload.TaxID,
		OwnerEmail:         ownerEmail,
		OwnerName:          ownerName,
	}
	result, err := h.appService.CreateOrganization(c, userID, cmd)
	if err != nil {
		return err
	}

	resp := CreateOrganizationResponse{
		Message:      "Organization created successfully. An invitation has been sent to the organization owner.",
		Organization: result.Organization,
		Invitation: OrganizationInvitationInfo{
			Sent:  result.InvitationSent,
			Email: result.OwnerEmail,
		},
	}

	return c.JSON(http.StatusCreated, resp)
}

func (h *OrganizationHandler) TransferOwnership(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errs.NewBadRequestError("invalid organization id format", false, nil, nil, nil)
	}

	var payload TransferOrganizationOwnershipPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("invalid request payload", false, nil, nil, nil)
	}
	if err := payload.Validate(); err != nil {
		return errs.NewBadRequestError(err.Error(), false, nil, nil, nil)
	}

	userID := middleware.GetUserID(c)
	cmd := &application.TransferOwnershipCommand{
		OrganizationID: id,
		NewOwnerEmail:  payload.NewOwnerEmail,
		NewOwnerName:   payload.NewOwnerName,
		ActorUserID:    userID,
		Notes:          payload.Notes,
	}

	if err := h.appService.TransferOwnership(c, cmd); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Organization ownership transferred successfully",
	})
}

func (h *OrganizationHandler) ResendOwnerInvite(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errs.NewBadRequestError("invalid organization id format", false, nil, nil, nil)
	}

	_, err = h.appService.ResendOwnerInvite(c, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Owner verification setup code has been resent successfully",
	})
}

func (h *OrganizationHandler) ListOrganizations(c echo.Context) error {
	userID := middleware.GetUserID(c)
	platformRole := middleware.GetUserRole(c)
	orgs, err := h.appService.ListOrganizations(c, userID, platformRole)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, orgs)
}

func (h *OrganizationHandler) GetOrganization(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errs.NewBadRequestError("invalid organization id format", false, nil, nil, nil)
	}

	org, err := h.appService.GetOrganizationByID(c, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, org)
}

func (h *OrganizationHandler) UpdateOrganization(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errs.NewBadRequestError("invalid organization id format", false, nil, nil, nil)
	}

	var payload UpdateOrganizationPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("invalid request payload", false, nil, nil, nil)
	}

	org, err := h.appService.UpdateOrganization(c, id, payload.Name, payload.Slug, payload.Plan, payload.CustomDomain, payload.Settings)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, org)
}



func (h *OrganizationHandler) GetOrganizationSettings(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errs.NewBadRequestError("invalid organization id format", false, nil, nil, nil)
	}

	settings, err := h.appService.GetOrganizationSettings(c, id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, settings)
}

func (h *OrganizationHandler) UpdateOrganizationSettings(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return errs.NewBadRequestError("invalid organization id format", false, nil, nil, nil)
	}

	var payload UpdateOrganizationSettingsPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("invalid request payload", false, nil, nil, nil)
	}

	settings, err := h.appService.UpdateOrganizationSettings(c, id, payload.LogoURL, payload.ThemeBranding, payload.CustomDomain, payload.SupportEmail, payload.SupportPhone, payload.CACNumber, payload.TINNumber, payload.TaxNumber, payload.BusinessType, payload.Address, payload.Timezone, payload.Currency, payload.Language)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, settings)
}


