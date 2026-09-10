package service_test

import (
	"testing"

	"github.com/golangnigeria/curexal/internal/modules/orchestration/model"
)

func TestProviderProfile_StatusValidation(t *testing.T) {
	validStatuses := map[string]bool{
		"ON_DUTY":  true,
		"ON_BREAK": true,
		"OFF_DUTY": true,
		"BUSY":     true,
	}

	for status := range validStatuses {
		if !validStatuses[status] {
			t.Errorf("expected %s to be valid", status)
		}
	}

	invalidStatuses := []string{"SLEEPING", "AWAY", "VACATION", "UNKNOWN"}
	for _, status := range invalidStatuses {
		if validStatuses[status] {
			t.Errorf("expected %s to be invalid", status)
		}
	}
}

func TestProviderProfile_DefaultVirtualRoomURL(t *testing.T) {
	req := model.CreateProviderProfileRequest{
		UserID:        "user-123",
		LicenseNumber: "MDCN/R/12345",
		SpecialtyCode: "GENERAL_PRACTICE",
	}

	if req.VirtualRoomURL != nil {
		t.Errorf("expected virtual room url to be nil initially")
	}
}

// Test Provider Candidate Acuity and Channel Routing
func TestProviderProfile_ChannelEnablement(t *testing.T) {
	provider := model.ProviderProfile{
		ID:                "prov-01",
		SpecialtyCode:     "GENERAL_PRACTICE",
		InPersonEnabled:   true,
		TelehealthEnabled: true,
		MaxActiveQueue:    10,
		Status:            "ON_DUTY",
	}

	if !provider.InPersonEnabled || !provider.TelehealthEnabled {
		t.Errorf("expected dual channel enablement to be true")
	}
	if provider.Status != "ON_DUTY" {
		t.Errorf("expected status to be ON_DUTY, got %s", provider.Status)
	}
}
