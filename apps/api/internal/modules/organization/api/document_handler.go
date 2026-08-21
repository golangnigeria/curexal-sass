package api

import (
	"io"
	"net/http"

	"github.com/golangnigeria/curexal/internal/modules/organization/application"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/golangnigeria/curexal/internal/shared/middleware"
	"github.com/golangnigeria/curexal/internal/shared/response"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type DocumentHandler struct {
	server        *server.Server
	docAppService *application.OrganizationDocumentApplicationService
}

func NewDocumentHandler(
	s *server.Server,
	docAppService *application.OrganizationDocumentApplicationService,
) *DocumentHandler {
	return &DocumentHandler{
		server:        s,
		docAppService: docAppService,
	}
}

func (h *DocumentHandler) UploadDocument(c echo.Context) error {
	orgIDParam := c.Param("id")
	orgID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid organization ID format")
	}

	docType := c.FormValue("document_type")
	if docType == "" {
		docType = c.FormValue("documentType")
	}
	if docType == "" {
		return errs.NewBadRequestError("document_type form field is required")
	}

	fileHeader, errFile := c.FormFile("file")
	if errFile != nil {
		return errs.NewBadRequestError("uploaded file is required in 'file' form field")
	}

	fileSrc, errOpen := fileHeader.Open()
	if errOpen != nil {
		return errs.NewBadRequestError("failed to open uploaded file stream")
	}
	defer fileSrc.Close()

	content, errRead := io.ReadAll(fileSrc)
	if errRead != nil {
		return errs.NewBadRequestError("failed to read file content")
	}

	callerID := middleware.GetUserID(c)
	doc, errUpload := h.docAppService.UploadDocument(c.Request().Context(), callerID, orgID, docType, fileHeader.Filename, content)
	if errUpload != nil {
		return errUpload
	}

	return response.CreatedEcho(c, doc)
}

func (h *DocumentHandler) ListDocuments(c echo.Context) error {
	orgIDParam := c.Param("id")
	orgID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid organization ID format")
	}

	callerID := middleware.GetUserID(c)
	docs, errList := h.docAppService.ListDocuments(c.Request().Context(), callerID, orgID)
	if errList != nil {
		return errList
	}

	return response.SuccessEcho(c, http.StatusOK, docs)
}

func (h *DocumentHandler) ReviewDocument(c echo.Context) error {
	docIDParam := c.Param("docID")
	docID, errParse := uuid.Parse(docIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid document ID format")
	}

	var payload ReviewDocumentPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("invalid request payload")
	}
	if err := payload.Validate(); err != nil {
		return errs.NewBadRequestError(err.Error())
	}

	reviewerID := middleware.GetUserID(c)
	status := domain.OrganizationDocumentStatus(payload.Status)
	errReview := h.docAppService.ReviewDocument(c.Request().Context(), reviewerID, docID, status, payload.RejectionReason)
	if errReview != nil {
		return errReview
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]string{
		"message": "Document review status updated successfully",
	})
}

func (h *DocumentHandler) ApproveOrganization(c echo.Context) error {
	orgIDParam := c.Param("id")
	orgID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid organization ID format")
	}

	reviewerID := middleware.GetUserID(c)
	errApprove := h.docAppService.ApproveOrganization(c.Request().Context(), reviewerID, orgID)
	if errApprove != nil {
		return errApprove
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]string{
		"message": "Organization verified and activated successfully",
	})
}

func (h *DocumentHandler) RejectOrganization(c echo.Context) error {
	orgIDParam := c.Param("id")
	orgID, errParse := uuid.Parse(orgIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid organization ID format")
	}

	var payload RejectOrganizationPayload
	if err := c.Bind(&payload); err != nil {
		return errs.NewBadRequestError("invalid request payload")
	}
	if err := payload.Validate(); err != nil {
		return errs.NewBadRequestError(err.Error())
	}

	reviewerID := middleware.GetUserID(c)
	errReject := h.docAppService.RejectOrganization(c.Request().Context(), reviewerID, orgID, payload.Reason)
	if errReject != nil {
		return errReject
	}

	return response.SuccessEcho(c, http.StatusOK, map[string]string{
		"message": "Organization verification rejected",
	})
}
