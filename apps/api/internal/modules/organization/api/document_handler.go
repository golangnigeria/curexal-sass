package api

import (
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/organization/application"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
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

func resolveOrgUUID(c echo.Context) (uuid.UUID, error) {
	orgIDParam := c.Param("id")
	if orgIDParam != "" && orgIDParam != "org_default" {
		if parsed, err := uuid.Parse(orgIDParam); err == nil {
			return parsed, nil
		}
	}

	if orgContextID := middleware.GetOrganizationID(c); orgContextID != "" && orgContextID != "org_default" {
		if parsed, err := uuid.Parse(orgContextID); err == nil {
			return parsed, nil
		}
	}

	if p := middleware.GetPrincipal(c); p != nil {
		if p.Organization.ActiveOrganizationID != "" && p.Organization.ActiveOrganizationID != "org_default" {
			if parsed, err := uuid.Parse(p.Organization.ActiveOrganizationID); err == nil {
				return parsed, nil
			}
		}
		if p.OrganizationID != "" && p.OrganizationID != "org_default" {
			if parsed, err := uuid.Parse(p.OrganizationID); err == nil {
				return parsed, nil
			}
		}
	}

	return uuid.MustParse("00000000-0000-0000-0000-000000000001"), nil
}

func (h *DocumentHandler) UploadDocument(c echo.Context) error {
	orgID, _ := resolveOrgUUID(c)

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
	orgID, _ := resolveOrgUUID(c)

	callerID := middleware.GetUserID(c)
	docs, errList := h.docAppService.ListDocuments(c.Request().Context(), callerID, orgID)
	if errList != nil {
		return errList
	}

	return response.SuccessEcho(c, http.StatusOK, docs)
}

func (h *DocumentHandler) PreviewDocument(c echo.Context) error {
	return h.streamDocument(c, false)
}

func (h *DocumentHandler) DownloadDocument(c echo.Context) error {
	return h.streamDocument(c, true)
}

func (h *DocumentHandler) streamDocument(c echo.Context, isDownload bool) error {
	docIDParam := c.Param("docID")
	if docIDParam == "" {
		docIDParam = c.Param("id")
	}
	docID, errParse := uuid.Parse(docIDParam)
	if errParse != nil {
		return errs.NewBadRequestError("invalid document ID format")
	}

	orgIDParam := c.Param("id")
	var orgID uuid.UUID
	if orgIDParam != "" && orgIDParam != docIDParam {
		orgID, _ = uuid.Parse(orgIDParam)
	}
	if orgID == uuid.Nil {
		if orgContextID := middleware.GetOrganizationID(c); orgContextID != "" {
			orgID, _ = uuid.Parse(orgContextID)
		}
	}
	if orgID == uuid.Nil {
		if p := middleware.GetPrincipal(c); p != nil && p.Organization.ActiveOrganizationID != "" {
			orgID, _ = uuid.Parse(p.Organization.ActiveOrganizationID)
		}
	}

	callerID := middleware.GetUserID(c)
	doc, rc, fileSize, mimeType, err := h.docAppService.GetDocumentForViewing(c.Request().Context(), callerID, orgID, docID, isDownload)
	if err != nil {
		return err
	}
	defer rc.Close()

	filename := doc.OriginalFilename
	if filename == "" {
		filename = "document"
	}

	disposition := "inline"
	if isDownload {
		disposition = "attachment"
	}

	c.Response().Header().Set("Content-Type", mimeType)
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=%q", disposition, filename))
	c.Response().Header().Set("Cache-Control", "private, no-store")
	c.Response().Header().Set("Accept-Ranges", "bytes")
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")
	if fileSize > 0 {
		c.Response().Header().Set("Content-Length", strconv.FormatInt(fileSize, 10))
	}

	return c.Stream(http.StatusOK, mimeType, rc)
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
