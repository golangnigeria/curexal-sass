package middleware

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// PatientGuard checks that the loaded request context contains a patient context.
// The patient context is populated ONCE per request by the auth middleware — ZERO extra DB queries here.
func PatientGuard(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		pCtx := GetPatientContext(c)
		if pCtx == nil {
			return echo.NewHTTPError(http.StatusForbidden, "Access denied: Patient profile required to access this resource.")
		}
		return next(c)
	}
}
