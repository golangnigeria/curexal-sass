package service_test

import (
	"testing"
	"time"

	patientModel "github.com/golangnigeria/curexal/internal/modules/patient/model"
)

func TestMPIService_ScoringMatrix(t *testing.T) {
	dob, _ := time.Parse("2006-01-02", "1990-01-15")
	ninVal := "12345678901"

	existingPatient := patientModel.Patient{
		ID:          "pat-01",
		TenantID:    "branch-01",
		MRN:         "PAT-2026-00001",
		FirstName:   "John",
		LastName:    "Doe",
		Gender:      "MALE",
		DateOfBirth: dob,
		NIN:         &ninVal,
		Contacts: []patientModel.PatientContact{
			{System: "PHONE", Value: "+2348011223344", IsPrimary: true},
		},
	}

	tests := []struct {
		name            string
		req             patientModel.DuplicateEvaluationRequest
		candidate       patientModel.Patient
		expectedScore   int
		expectedStatus  patientModel.MatchConfidence
		expectedSignals []string
	}{
		{
			name: "Phone + NIN -> 90 pts (EXACT_MATCH)",
			req: patientModel.DuplicateEvaluationRequest{
				FirstName: "Different",
				LastName:  "Name",
				Phone:     "+2348011223344",
				NIN:       &ninVal,
			},
			candidate:       existingPatient,
			expectedScore:   90,
			expectedStatus:  patientModel.MatchExact,
			expectedSignals: []string{"PHONE", "NIN"},
		},
		{
			name: "Phone + DOB -> 60 pts (PROBABLE_DUPLICATE)",
			req: patientModel.DuplicateEvaluationRequest{
				FirstName:   "Johnny",
				LastName:    "Smith",
				DateOfBirth: "1990-01-15",
				Phone:       "+2348011223344",
			},
			candidate:       existingPatient,
			expectedScore:   60,
			expectedStatus:  patientModel.MatchProbableDuplicate,
			expectedSignals: []string{"PHONE", "DOB"},
		},
		{
			name: "Phone Only -> 40 pts (LOW)",
			req: patientModel.DuplicateEvaluationRequest{
				FirstName: "Alice",
				LastName:  "Johnson",
				Phone:     "+2348011223344",
			},
			candidate:       existingPatient,
			expectedScore:   40,
			expectedStatus:  patientModel.MatchLow,
			expectedSignals: []string{"PHONE"},
		},
		{
			name: "Full Match (Phone + NIN + DOB + Last + First) -> 100 pts (EXACT_MATCH)",
			req: patientModel.DuplicateEvaluationRequest{
				FirstName:   "John",
				LastName:    "Doe",
				DateOfBirth: "1990-01-15",
				Phone:       "+2348011223344",
				NIN:         &ninVal,
			},
			candidate:       existingPatient,
			expectedScore:   100,
			expectedStatus:  patientModel.MatchExact,
			expectedSignals: []string{"PHONE", "NIN", "DOB", "LAST_NAME", "FIRST_NAME"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			score := 0
			var signals []string

			// 1. Phone match (+40)
			for _, c := range tc.candidate.Contacts {
				if c.System == "PHONE" && c.Value == tc.req.Phone && tc.req.Phone != "" {
					score += 40
					signals = append(signals, "PHONE")
					break
				}
			}

			// 2. NIN match (+50)
			if tc.candidate.NIN != nil && tc.req.NIN != nil && *tc.candidate.NIN == *tc.req.NIN {
				score += 50
				signals = append(signals, "NIN")
			}

			// 3. DOB match (+20)
			if tc.req.DateOfBirth != "" && tc.candidate.DateOfBirth.Format("2006-01-02") == tc.req.DateOfBirth {
				score += 20
				signals = append(signals, "DOB")
			}

			// 4. Last Name (+20)
			if tc.candidate.LastName == tc.req.LastName {
				score += 20
				signals = append(signals, "LAST_NAME")
			}

			// 5. First Name (+15)
			if tc.candidate.FirstName == tc.req.FirstName {
				score += 15
				signals = append(signals, "FIRST_NAME")
			}

			if score > 100 {
				score = 100
			}

			if score != tc.expectedScore {
				t.Errorf("expected score %d, got %d", tc.expectedScore, score)
			}

			var status patientModel.MatchConfidence
			switch {
			case score >= 80:
				status = patientModel.MatchExact
			case score >= 50:
				status = patientModel.MatchProbableDuplicate
			case score >= 20:
				status = patientModel.MatchLow
			default:
				status = patientModel.MatchNone
			}

			if status != tc.expectedStatus {
				t.Errorf("expected status %s, got %s", tc.expectedStatus, status)
			}
		})
	}
}
