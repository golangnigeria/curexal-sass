package service_test

import (
	"testing"

	"github.com/golangnigeria/curexal/internal/modules/orchestration/model"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/service"
)

func TestEvaluateAcuity(t *testing.T) {
	svc := service.NewCareRequestService(nil, nil)

	tempHigh := 39.8
	tempMod := 38.4
	tempNormal := 36.8

	spo2Low := 89
	spo2Mod := 93
	spo2Normal := 98

	painSevere := 9
	painMod := 6
	painMild := 2

	sysHigh := 185
	sysNormal := 115

	tests := []struct {
		name     string
		payload  model.SubmitTriagePayload
		expected string
	}{
		{
			name: "Severe fever & low oxygen -> RED",
			payload: model.SubmitTriagePayload{
				Temperature: &tempHigh,
				SpO2:        &spo2Low,
			},
			expected: "RED",
		},
		{
			name: "Severe pain score 9 -> RED",
			payload: model.SubmitTriagePayload{
				PainScore: &painSevere,
			},
			expected: "RED",
		},
		{
			name: "Hypertensive crisis (Sys BP 185) -> RED",
			payload: model.SubmitTriagePayload{
				SystolicBP: &sysHigh,
			},
			expected: "RED",
		},
		{
			name: "Moderate fever 38.4 -> YELLOW",
			payload: model.SubmitTriagePayload{
				Temperature: &tempMod,
				SpO2:        &spo2Normal,
			},
			expected: "YELLOW",
		},
		{
			name: "Moderate SpO2 93 -> YELLOW",
			payload: model.SubmitTriagePayload{
				SpO2: &spo2Mod,
			},
			expected: "YELLOW",
		},
		{
			name: "Moderate pain 6 -> YELLOW",
			payload: model.SubmitTriagePayload{
				PainScore: &painMod,
			},
			expected: "YELLOW",
		},
		{
			name: "Normal vitals -> GREEN",
			payload: model.SubmitTriagePayload{
				Temperature: &tempNormal,
				SpO2:        &spo2Normal,
				PainScore:   &painMild,
				SystolicBP:  &sysNormal,
			},
			expected: "GREEN",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := svc.EvaluateAcuity(tc.payload)
			if result != tc.expected {
				t.Errorf("expected acuity %s, got %s", tc.expected, result)
			}
		})
	}
}
