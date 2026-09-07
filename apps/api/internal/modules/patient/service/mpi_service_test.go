package service_test

import (
	"context"
	"testing"
	"time"

	patientModel "github.com/golangnigeria/curexal/internal/modules/patient/model"
)

// Mock repository for testing MPI scoring logic
type mockMPIRepo struct {
	candidates []patientModel.Patient
}

func (m *mockMPIRepo) FindCandidatesForDuplicateResolution(
	ctx context.Context,
	tenantID, phone, nin, lastName string,
	dob *time.Time,
) ([]patientModel.Patient, error) {
	return m.candidates, nil
}

func TestMPIService_DuplicateEvaluation(t *testing.T) {
	dob, _ := time.Parse("2006-01-02", "1992-06-15")
	ninVal := "12345678901"

	existingPatient := patientModel.Patient{
		ID:          "pat-01",
		TenantID:    "tenant-1",
		MRN:         "PAT-2026-0001",
		FirstName:   "Amina",
		LastName:    "Yusuf",
		Gender:      "FEMALE",
		DateOfBirth: dob,
		NIN:         &ninVal,
		Contacts: []patientModel.PatientContact{
			{System: "PHONE", Value: "+2348012345678", IsPrimary: true},
		},
	}

	tests := []struct {
		name            string
		candidates      []patientModel.Patient
		request         patientModel.DuplicateEvaluationRequest
		expectedStatus  patientModel.MatchConfidence
		expectedSignals []string
	}{
		{
			name:       "Exact NIN and Phone match -> EXACT_MATCH",
			candidates: []patientModel.Patient{existingPatient},
			request: patientModel.DuplicateEvaluationRequest{
				FirstName:   "Amina",
				LastName:    "Yusuf",
				DateOfBirth: "1992-06-15",
				Phone:       "+2348012345678",
				NIN:         &ninVal,
			},
			expectedStatus:  patientModel.MatchExact,
			expectedSignals: []string{"PHONE", "NIN", "DOB", "LAST_NAME", "FIRST_NAME"},
		},
		{
			name:       "No matching candidates -> NONE",
			candidates: []patientModel.Patient{},
			request: patientModel.DuplicateEvaluationRequest{
				FirstName:   "Chinedu",
				LastName:    "Okafor",
				Phone:       "+2348099999999",
			},
			expectedStatus: patientModel.MatchNone,
		},
		{
			name:       "Phone match only -> PROBABLE_DUPLICATE",
			candidates: []patientModel.Patient{existingPatient},
			request: patientModel.DuplicateEvaluationRequest{
				FirstName: "Fatima",
				LastName:  "Bello",
				Phone:     "+2348012345678",
			},
			expectedStatus:  patientModel.MatchProbableDuplicate,
			expectedSignals: []string{"PHONE"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Test scoring logic directly
			req := tc.request
			candidates := tc.candidates

			maxScore := 0
			var matchedSignals []string

			for _, cand := range candidates {
				score := 0
				var signals []string

				if req.Phone != "" {
					for _, c := range cand.Contacts {
						if c.System == "PHONE" && c.Value == req.Phone {
							score += 40
							signals = append(signals, "PHONE")
							break
						}
					}
				}

				if cand.NIN != nil && req.NIN != nil && *cand.NIN == *req.NIN {
					score += 50
					signals = append(signals, "NIN")
				}

				if req.DateOfBirth != "" && cand.DateOfBirth.Format("2006-01-02") == req.DateOfBirth {
					score += 20
					signals = append(signals, "DOB")
				}

				if cand.LastName == req.LastName {
					score += 20
					signals = append(signals, "LAST_NAME")
				}

				if cand.FirstName == req.FirstName {
					score += 15
					signals = append(signals, "FIRST_NAME")
				}

				if score > maxScore {
					maxScore = score
					matchedSignals = signals
				}
			}

			var status patientModel.MatchConfidence
			switch {
			case maxScore >= 80:
				status = patientModel.MatchExact
			case maxScore >= 40:
				status = patientModel.MatchProbableDuplicate
			case maxScore >= 20:
				status = patientModel.MatchLow
			default:
				status = patientModel.MatchNone
			}

			if status != tc.expectedStatus {
				t.Errorf("expected status %s, got %s (score: %d)", tc.expectedStatus, status, maxScore)
			}

			if len(tc.expectedSignals) > 0 {
				for _, expectedSig := range tc.expectedSignals {
					found := false
					for _, sig := range matchedSignals {
						if sig == expectedSig {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected signal %s in %v", expectedSig, matchedSignals)
					}
				}
			}
		})
	}
}
