package service_test

import (
	"regexp"
	"testing"

	"github.com/golangnigeria/curexal/internal/modules/operations/service"
)

func TestAppointmentService_GenerateAppointmentNumber(t *testing.T) {
	svc := service.NewAppointmentService(nil, nil)
	aptRegex := regexp.MustCompile(`^APT-\d{4}-\d{5}$`)

	seen := make(map[string]bool)
	count := 500

	for i := 0; i < count; i++ {
		num := svc.GenerateAppointmentNumber()
		if !aptRegex.MatchString(num) {
			t.Fatalf("appointment number %s does not match expected pattern ^APT-YYYY-XXXXX$", num)
		}
		if seen[num] {
			t.Fatalf("collision detected on appointment number %s", num)
		}
		seen[num] = true
	}
}

func TestAppointmentService_DeliveryChannelValidation(t *testing.T) {
	validChannels := map[string]bool{
		"in_person":      true,
		"video":          true,
		"telephone":      true,
		"secure_message": true,
	}

	for ch := range validChannels {
		if !validChannels[ch] {
			t.Errorf("channel %s should be valid", ch)
		}
	}

	invalid := []string{"email", "carrier_pigeon", "drone"}
	for _, ch := range invalid {
		if validChannels[ch] {
			t.Errorf("channel %s should be invalid", ch)
		}
	}
}
