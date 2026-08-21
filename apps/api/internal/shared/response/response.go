package response

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// APIEnvelope represents the standardized JSON API response contract across all endpoints.
type APIEnvelope struct {
	Data   interface{}            `json:"data"`
	Meta   map[string]interface{} `json:"meta"`
	Links  map[string]interface{} `json:"links"`
	Errors []interface{}          `json:"errors"`
}

func BuildEnvelope(data interface{}, errors []interface{}) APIEnvelope {
	meta := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"version":   "1.0",
	}

	links := map[string]interface{}{}

	if errors == nil {
		errors = []interface{}{}
	}

	return APIEnvelope{
		Data:   data,
		Meta:   meta,
		Links:  links,
		Errors: errors,
	}
}

func JSON(w http.ResponseWriter, status int, data interface{}, errors []interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(BuildEnvelope(data, errors))
}

func Success(w http.ResponseWriter, status int, data interface{}) {
	JSON(w, status, data, nil)
}

func Error(w http.ResponseWriter, status int, message string) {
	errPayload := []interface{}{
		map[string]interface{}{
			"message": message,
			"code":    http.StatusText(status),
		},
	}
	JSON(w, status, nil, errPayload)
}

// Echo Helpers
func SuccessEcho(c echo.Context, status int, data interface{}) error {
	return c.JSON(status, BuildEnvelope(data, nil))
}

func CreatedEcho(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, BuildEnvelope(data, nil))
}

func ErrorEcho(c echo.Context, status int, code, message string) error {
	errPayload := []interface{}{
		map[string]interface{}{
			"code":    code,
			"message": message,
		},
	}
	return c.JSON(status, BuildEnvelope(nil, errPayload))
}

func BadRequestEcho(c echo.Context, message string) error {
	return ErrorEcho(c, http.StatusBadRequest, "BAD_REQUEST", message)
}

func UnauthorizedEcho(c echo.Context, message string) error {
	return ErrorEcho(c, http.StatusUnauthorized, "UNAUTHORIZED", message)
}

func ForbiddenEcho(c echo.Context, message string) error {
	return ErrorEcho(c, http.StatusForbidden, "FORBIDDEN", message)
}

func NotFoundEcho(c echo.Context, message string) error {
	return ErrorEcho(c, http.StatusNotFound, "NOT_FOUND", message)
}

func InternalErrorEcho(c echo.Context, message string) error {
	return ErrorEcho(c, http.StatusInternalServerError, "INTERNAL_ERROR", message)
}
