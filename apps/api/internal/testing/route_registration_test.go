package testing

import (
	"testing"

	"github.com/golangnigeria/curexal/internal/bootstrap"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

func TestCriticalRoutesExistOnEchoEngine(t *testing.T) {
	e := echo.New()
	logger := zerolog.Nop()
	s := &server.Server{
		Echo:   e,
		Logger: &logger,
	}

	reg := bootstrap.InitModules(s)
	if reg == nil {
		t.Fatalf("expected module registry to initialize")
	}

	routes := e.Routes()
	registeredMap := make(map[string]map[string]bool)
	for _, r := range routes {
		if registeredMap[r.Path] == nil {
			registeredMap[r.Path] = make(map[string]bool)
		}
		registeredMap[r.Path][r.Method] = true
	}

	criticalRoutes := []struct {
		method string
		path   string
	}{
		{"POST", "/api/v1/auth/sign-in"},
		{"POST", "/api/v1/organizations"},
		{"GET", "/api/v1/tenant/active"},
		{"POST", "/api/v1/lims/orders"},
		{"POST", "/api/v1/lims/specimens/accession"},
		{"POST", "/api/v1/lims/results"},
		{"POST", "/api/v1/clinical/patient-visits"},
		{"GET", "/api/v1/marketplace/capabilities"},
		{"GET", "/api/v1/organizations/:id/capabilities"},
		{"POST", "/api/v1/organizations/:id/documents"},
		{"PATCH", "/api/v1/platform/documents/:docID/review"},
		{"POST", "/api/v1/platform/organizations/:id/approve"},
		{"GET", "/api/v1/bootstrap"},
		{"GET", "/api/v1/audit-logs/tenant"},
		{"POST", "/api/v1/organizations/:id/marketplace/orders"},
		{"POST", "/api/v1/billing/webhooks/:provider"},

		// Day 2 Clinical Intake, MPI, Provider & Queue Routes
		{"POST", "/api/v1/patients/canonical"},
		{"POST", "/api/v1/patients/mpi/evaluate"},
		{"GET", "/api/v1/providers/profiles"},
		{"PUT", "/api/v1/providers/profiles/:id/status"},
		{"POST", "/api/v1/appointments"},
		{"GET", "/api/v1/appointments"},
		{"POST", "/api/v1/orchestration/requests"},
		{"POST", "/api/v1/orchestration/requests/:id/triage"},
		{"GET", "/api/v1/orchestration/requests"},
	}

	for _, cr := range criticalRoutes {
		if methods, exists := registeredMap[cr.path]; !exists || !methods[cr.method] {
			t.Errorf("critical route missing from Echo router: %s %s", cr.method, cr.path)
		}
	}
}
