package service

import (
	"context"
	"strings"
	"time"

	patientModel "github.com/golangnigeria/curexal/internal/modules/patient/model"
	patientRepo "github.com/golangnigeria/curexal/internal/modules/patient/repository"
)

type MPIService struct {
	patientRepo *patientRepo.CanonicalPatientRepository
}

func NewMPIService(patientRepo *patientRepo.CanonicalPatientRepository) *MPIService {
	return &MPIService{patientRepo: patientRepo}
}

// EvaluateDuplicates queries candidates and scores identity resolution signals
func (s *MPIService) EvaluateDuplicates(
	ctx context.Context,
	tenantID string,
	req patientModel.DuplicateEvaluationRequest,
) (*patientModel.DuplicateEvaluationResponse, error) {
	var dob *time.Time
	if strings.TrimSpace(req.DateOfBirth) != "" {
		if t, err := time.Parse("2006-01-02", req.DateOfBirth); err == nil {
			dob = &t
		} else if t, err := time.Parse(time.RFC3339, req.DateOfBirth); err == nil {
			dob = &t
		}
	}

	nin := ""
	if req.NIN != nil {
		nin = *req.NIN
	}

	candidates, err := s.patientRepo.FindCandidatesForDuplicateResolution(
		ctx, tenantID, req.Phone, nin, req.LastName, dob,
	)
	if err != nil {
		return nil, err
	}

	response := &patientModel.DuplicateEvaluationResponse{
		MatchStatus: patientModel.MatchNone,
		Candidates:  make([]patientModel.DuplicateMatchCandidate, 0),
	}

	if len(candidates) == 0 {
		return response, nil
	}

	maxScore := 0

	for _, cand := range candidates {
		score := 0
		var signals []string

		// 1. Phone matching (+40 pts)
		cleanReqPhone := strings.TrimSpace(req.Phone)
		for _, c := range cand.Contacts {
			if c.System == "PHONE" && strings.TrimSpace(c.Value) == cleanReqPhone && cleanReqPhone != "" {
				score += 40
				signals = append(signals, "PHONE")
				break
			}
		}

		// 2. NIN matching (+50 pts)
		if cand.NIN != nil && nin != "" && strings.EqualFold(strings.TrimSpace(*cand.NIN), strings.TrimSpace(nin)) {
			score += 50
			signals = append(signals, "NIN")
		}

		// 3. Date of Birth matching (+20 pts)
		if dob != nil && !dob.IsZero() {
			if cand.DateOfBirth.Format("2006-01-02") == dob.Format("2006-01-02") {
				score += 20
				signals = append(signals, "DOB")
			}
		}

		// 4. Exact Last Name matching (+20 pts)
		if strings.EqualFold(strings.TrimSpace(cand.LastName), strings.TrimSpace(req.LastName)) {
			score += 20
			signals = append(signals, "LAST_NAME")
		}

		// 5. First Name matching (+15 pts)
		if strings.EqualFold(strings.TrimSpace(cand.FirstName), strings.TrimSpace(req.FirstName)) {
			score += 15
			signals = append(signals, "FIRST_NAME")
		}

		if score > 100 {
			score = 100
		}

		var level patientModel.MatchConfidence
		switch {
		case score >= 80:
			level = patientModel.MatchExact
		case score >= 50:
			level = patientModel.MatchProbableDuplicate
		case score >= 20:
			level = patientModel.MatchLow
		default:
			level = patientModel.MatchNone
		}

		if score > maxScore {
			maxScore = score
		}

		if score >= 35 {
			response.Candidates = append(response.Candidates, patientModel.DuplicateMatchCandidate{
				PatientID:       cand.ID,
				MRN:             cand.MRN,
				FirstName:       cand.FirstName,
				LastName:        cand.LastName,
				DateOfBirth:     cand.DateOfBirth.Format("2006-01-02"),
				Gender:          cand.Gender,
				MatchedSignals:  signals,
				ConfidenceScore: score,
				ConfidenceLevel: level,
			})
		}
	}

	switch {
	case maxScore >= 80:
		response.MatchStatus = patientModel.MatchExact
	case maxScore >= 50:
		response.MatchStatus = patientModel.MatchProbableDuplicate
	case maxScore >= 35:
		response.MatchStatus = patientModel.MatchLow
	default:
		response.MatchStatus = patientModel.MatchNone
	}

	return response, nil
}
