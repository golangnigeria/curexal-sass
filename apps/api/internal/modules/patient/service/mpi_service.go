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

// soundex calculates the standard American soundex code for phonetic name matching
func soundex(s string) string {
	s = strings.ToUpper(strings.TrimSpace(s))
	if len(s) == 0 {
		return ""
	}

	charMap := func(c byte) byte {
		switch c {
		case 'B', 'F', 'P', 'V':
			return '1'
		case 'C', 'G', 'J', 'K', 'Q', 'S', 'X', 'Z':
			return '2'
		case 'D', 'T':
			return '3'
		case 'L':
			return '4'
		case 'M', 'N':
			return '5'
		case 'R':
			return '6'
		default:
			return '0'
		}
	}

	res := []byte{s[0]}
	lastCode := charMap(s[0])

	for i := 1; i < len(s) && len(res) < 4; i++ {
		c := s[i]
		if c < 'A' || c > 'Z' {
			continue
		}
		code := charMap(c)
		if code != '0' && code != lastCode {
			res = append(res, code)
		}
		lastCode = code
	}

	for len(res) < 4 {
		res = append(res, '0')
	}

	return string(res[:4])
}

// ScoreCandidates evaluates duplicate match confidence for a list of patient candidates
func (s *MPIService) ScoreCandidates(
	req patientModel.DuplicateEvaluationRequest,
	candidates []patientModel.Patient,
) []patientModel.DuplicateMatchCandidate {
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

	var results []patientModel.DuplicateMatchCandidate

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

		// 4. Exact Last Name matching (+20 pts) or Soundex matching (+10 pts)
		lastNameExact := false
		if strings.TrimSpace(req.LastName) != "" && strings.EqualFold(strings.TrimSpace(cand.LastName), strings.TrimSpace(req.LastName)) {
			score += 20
			signals = append(signals, "LAST_NAME")
			lastNameExact = true
		} else if strings.TrimSpace(req.LastName) != "" {
			reqSx := soundex(req.LastName)
			candSx := soundex(cand.LastName)
			if reqSx != "" && reqSx == candSx {
				score += 10
				signals = append(signals, "SOUNDEX")
			}
		}

		// 5. First Name matching (+15 pts) or Soundex matching (+10 pts if Soundex not already applied)
		if strings.TrimSpace(req.FirstName) != "" && strings.EqualFold(strings.TrimSpace(cand.FirstName), strings.TrimSpace(req.FirstName)) {
			score += 15
			signals = append(signals, "FIRST_NAME")
		} else if strings.TrimSpace(req.FirstName) != "" && !lastNameExact {
			reqSx := soundex(req.FirstName)
			candSx := soundex(cand.FirstName)
			hasSoundex := false
			for _, sig := range signals {
				if sig == "SOUNDEX" {
					hasSoundex = true
					break
				}
			}
			if !hasSoundex && reqSx != "" && reqSx == candSx {
				score += 10
				signals = append(signals, "SOUNDEX")
			}
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

		lastVisit := cand.UpdatedAt.UTC().Format(time.RFC3339)
		regChannel := cand.RegistrationChannel
		if regChannel == "" {
			regChannel = "RECEPTION"
		}
		results = append(results, patientModel.DuplicateMatchCandidate{
			PatientID:           cand.ID,
			MRN:                 cand.MRN,
			FirstName:           cand.FirstName,
			LastName:            cand.LastName,
			DateOfBirth:         cand.DateOfBirth.Format("2006-01-02"),
			Gender:              cand.Gender,
			MatchedSignals:      signals,
			ConfidenceScore:     score,
			Score:               score,
			ConfidenceLevel:     level,
			LastVisitAt:         &lastVisit,
			RegistrationChannel: &regChannel,
		})
	}

	return results
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

	scoredCandidates := s.ScoreCandidates(req, candidates)
	maxScore := 0

	for _, cand := range scoredCandidates {
		if cand.ConfidenceScore > maxScore {
			maxScore = cand.ConfidenceScore
		}
		if cand.ConfidenceScore >= 20 {
			response.Candidates = append(response.Candidates, cand)
		}
	}

	switch {
	case maxScore >= 80:
		response.MatchStatus = patientModel.MatchExact
	case maxScore >= 50:
		response.MatchStatus = patientModel.MatchProbableDuplicate
	case maxScore >= 20:
		response.MatchStatus = patientModel.MatchLow
	default:
		response.MatchStatus = patientModel.MatchNone
	}

	return response, nil
}
