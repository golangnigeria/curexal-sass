package api

import (
	"fmt"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/storage"
	"github.com/golangnigeria/curexal/internal/shared/errs"
	"github.com/labstack/echo/v4"
)

type StorageHandler struct {
	storageService storage.ObjectStorage
}

func NewStorageHandler(storageService storage.ObjectStorage) *StorageHandler {
	return &StorageHandler{storageService: storageService}
}

func (h *StorageHandler) DownloadFile(c echo.Context) error {
	ctx := c.Request().Context()
	key := c.QueryParam("key")
	if key == "" {
		return errs.NewBadRequestError("storage key query parameter 'key' is required")
	}

	// Validate expiration timestamp if provided
	expiresParam := c.QueryParam("expires")
	if expiresParam != "" {
		if expiresUnix, err := strconv.ParseInt(expiresParam, 10, 64); err == nil {
			if expiresUnix > 0 && time.Now().Unix() > expiresUnix {
				return echo.NewHTTPError(http.StatusUnauthorized, "Presigned download link has expired")
			}
		}
	}

	rc, err := h.storageService.GetObject(ctx, key)
	if err != nil {
		return errs.NewNotFoundError(fmt.Sprintf("document object %q not found in storage", key))
	}
	defer rc.Close()

	// Determine MIME type and clean filename from key
	ext := strings.ToLower(filepath.Ext(key))
	contentType := "application/octet-stream"
	switch ext {
	case ".pdf":
		contentType = "application/pdf"
	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	default:
		if m := mime.TypeByExtension(ext); m != "" {
			contentType = m
		}
	}

	filename := filepath.Base(key)

	// Set headers for clean inline browser preview
	c.Response().Header().Set("Content-Type", contentType)
	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", filename))
	c.Response().Header().Set("X-Content-Type-Options", "nosniff")

	return c.Stream(http.StatusOK, contentType, rc)
}
