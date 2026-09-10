package service_test

import (
	"regexp"
	"testing"

	"github.com/golangnigeria/curexal/internal/modules/patient/service"
)

func TestCanonicalPatientService_GenerateMRN(t *testing.T) {
	svc := service.NewCanonicalPatientService(nil, nil, nil)
	mrnRegex := regexp.MustCompile(`^PAT-\d{4}-\d{5}$`)

	seen := make(map[string]bool)
	count := 1000

	for i := 0; i < count; i++ {
		mrn := svc.GenerateMRN()
		if !mrnRegex.MatchString(mrn) {
			t.Fatalf("MRN %s does not match expected format ^PAT-YYYY-XXXXX$", mrn)
		}
		if seen[mrn] {
			t.Fatalf("MRN collision detected on iteration %d: %s", i, mrn)
		}
		seen[mrn] = true
	}

	if len(seen) != count {
		t.Errorf("expected %d unique MRNs, got %d", count, len(seen))
	}
}
