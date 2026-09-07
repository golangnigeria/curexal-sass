package patient

// MatchConfidence represents the resolution score category
type MatchConfidence string

const (
	MatchNone             MatchConfidence = "NONE"
	MatchLow              MatchConfidence = "LOW"
	MatchProbableDuplicate MatchConfidence = "PROBABLE_DUPLICATE"
	MatchExact            MatchConfidence = "EXACT_MATCH"
)

// DuplicateEvaluationRequest contains signals to check for duplicates
type DuplicateEvaluationRequest struct {
	FirstName   string  `json:"firstName"`
	MiddleName  *string `json:"middleName,omitempty"`
	LastName    string  `json:"lastName"`
	DateOfBirth string  `json:"dateOfBirth,omitempty"`
	Phone       string  `json:"phone,omitempty"`
	Email       string  `json:"email,omitempty"`
	NIN         *string `json:"nin,omitempty"`
}

// DuplicateMatchCandidate is a matched patient with scoring breakdown
type DuplicateMatchCandidate struct {
	PatientID       string          `json:"patientId"`
	MRN             string          `json:"mrn"`
	FirstName       string          `json:"firstName"`
	LastName        string          `json:"lastName"`
	DateOfBirth     string          `json:"dateOfBirth"`
	Gender          string          `json:"gender"`
	MatchedSignals  []string        `json:"matchedSignals"` // "PHONE", "NIN", "DOB", "NAME"
	ConfidenceScore int             `json:"confidenceScore"` // 0 - 100
	ConfidenceLevel MatchConfidence `json:"confidenceLevel"`
}

// DuplicateEvaluationResponse is returned by the MPI engine
type DuplicateEvaluationResponse struct {
	MatchStatus MatchConfidence            `json:"matchStatus"`
	Candidates  []DuplicateMatchCandidate `json:"candidates"`
}
