package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/model"
	"github.com/golangnigeria/curexal/internal/modules/orchestration/repository"
	"github.com/google/uuid"
)

type CareRequestService struct {
	server *server.Server
	repo   *repository.CareRequestRepository
}

func NewCareRequestService(s *server.Server, repo *repository.CareRequestRepository) *CareRequestService {
	return &CareRequestService{
		server: s,
		repo:   repo,
	}
}

// GenerateRequestNumber generates a tracking number for care requests (e.g. REQ-2026-92831)
func (s *CareRequestService) GenerateRequestNumber() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(900000))
	val := n.Int64() + 100000
	year := time.Now().Year()
	return fmt.Sprintf("REQ-%d-%d", year, val)
}

// CreateCareRequest initializes a care request and provisions its live Care Journey milestones
func (s *CareRequestService) CreateCareRequest(
	ctx context.Context,
	tenantID string,
	payload model.CreateCareRequestPayload,
) (*model.CareRequest, error) {
	if tenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	if payload.PatientID == "" {
		return nil, errors.New("patientId is required")
	}

	reqID := uuid.New().String()
	reqNumber := s.GenerateRequestNumber()

	urgency := payload.Urgency
	if urgency == "" {
		urgency = "ROUTINE"
	}

	channel := strings.ToLower(strings.TrimSpace(payload.DeliveryChannel))
	if channel == "" {
		if payload.PreferredMode != "" {
			switch strings.ToUpper(payload.PreferredMode) {
			case "VIDEO":
				channel = "video"
			case "AUDIO":
				channel = "telephone"
			case "ASYNC_CHAT":
				channel = "secure_message"
			default:
				channel = "in_person"
			}
		} else {
			channel = "in_person"
		}
	}

	prefMode := payload.PreferredMode
	if prefMode == "" {
		switch channel {
		case "video":
			prefMode = "VIDEO"
		case "telephone":
			prefMode = "AUDIO"
		case "secure_message":
			prefMode = "ASYNC_CHAT"
		default:
			prefMode = "IN_PERSON"
		}
	}

	serviceType := payload.ServiceType
	if serviceType == "" {
		serviceType = "GENERAL_CONSULTATION"
	}

	status := "WAITING_TRIAGE"
	if channel == "video" {
		status = "WAITING_VIRTUAL_ROOM"
	}

	careReq := &model.CareRequest{
		ID:                  reqID,
		TenantID:            tenantID,
		PatientID:           payload.PatientID,
		RequestNumber:       reqNumber,
		ServiceType:         serviceType,
		PreferredMode:       prefMode,
		Urgency:             urgency,
		Status:              status,
		ChiefComplaint:      &payload.ChiefComplaint,
		SymptomsJSON:        payload.Symptoms,
		PreferredTimeWindow: payload.PreferredTimeWindow,
	}

	if err := s.repo.CreateCareRequest(ctx, careReq); err != nil {
		return nil, fmt.Errorf("failed to create care request: %w", err)
	}

	// Initialize Care Journey Milestones
	now := time.Now()
	descIntake := fmt.Sprintf("Care demand submitted for %s via %s", payload.ServiceType, payload.PreferredMode)
	_ = s.repo.CreateCareJourneyMilestone(ctx, &model.CareJourneyMilestone{
		TenantID:    tenantID,
		PatientID:   payload.PatientID,
		StageCode:   "INTAKE",
		Title:       "Care Intake Submitted",
		Description: &descIntake,
		Status:      "COMPLETED",
		CompletedAt: &now,
	})

	descTriage := "Clinical triage assessment and qualified doctor allocation"
	_ = s.repo.CreateCareJourneyMilestone(ctx, &model.CareJourneyMilestone{
		TenantID:    tenantID,
		PatientID:   payload.PatientID,
		StageCode:   "TRIAGE",
		Title:       "Triage & Doctor Matching",
		Description: &descTriage,
		Status:      "IN_PROGRESS",
	})

	s.server.Logger.Info().
		Str("tenantId", tenantID).
		Str("patientId", payload.PatientID).
		Str("requestNumber", reqNumber).
		Msg("Care Request created and Care Journey initialized")

	return careReq, nil
}

// ListCareRequests queries care requests with filters
func (s *CareRequestService) ListCareRequests(
	ctx context.Context,
	tenantID string,
	filter model.CareRequestFilter,
) (*model.CareRequestListResponse, error) {
	items, total, err := s.repo.ListCareRequests(ctx, tenantID, filter)
	if err != nil {
		return nil, err
	}
	return &model.CareRequestListResponse{
		Items:  items,
		Total:  total,
		Limit:  filter.Limit,
		Offset: filter.Offset,
	}, nil
}

// GetCareRequestByID retrieves single request by ID
func (s *CareRequestService) GetCareRequestByID(ctx context.Context, tenantID, requestID string) (*model.CareRequest, error) {
	return s.repo.GetCareRequestByID(ctx, tenantID, requestID)
}

// ListPatientCareRequests retrieves all requests for a specific patient
func (s *CareRequestService) ListPatientCareRequests(ctx context.Context, tenantID, patientID string) ([]model.CareRequest, error) {
	return s.repo.ListCareRequestsByPatient(ctx, tenantID, patientID)
}

// GetPatientCareJourney retrieves milestones along the patient's Care Journey
func (s *CareRequestService) GetPatientCareJourney(ctx context.Context, tenantID, patientID string) ([]model.CareJourneyMilestone, error) {
	return s.repo.GetPatientCareJourney(ctx, tenantID, patientID)
}

// EvaluateAcuity computes clinical urgency based on vital parameters
func (s *CareRequestService) EvaluateAcuity(p model.SubmitTriagePayload) string {
	if p.AcuityOverride != nil && *p.AcuityOverride != "" {
		return *p.AcuityOverride
	}

	// 1. Emergency thresholds (RED)
	if (p.Temperature != nil && *p.Temperature >= 39.5) ||
		(p.SpO2 != nil && *p.SpO2 <= 91) ||
		(p.PainScore != nil && *p.PainScore >= 8) ||
		(p.SystolicBP != nil && (*p.SystolicBP >= 180 || *p.SystolicBP <= 85)) {
		return "RED"
	}

	// 2. Urgent thresholds (YELLOW)
	if (p.Temperature != nil && *p.Temperature >= 38.2) ||
		(p.SpO2 != nil && *p.SpO2 <= 94) ||
		(p.PainScore != nil && *p.PainScore >= 5) ||
		(p.SystolicBP != nil && (*p.SystolicBP >= 140 || *p.SystolicBP <= 95)) {
		return "YELLOW"
	}

	return "GREEN"
}

// SubmitTriage records clinical triage vitals and updates the Care Request status
func (s *CareRequestService) SubmitTriage(
	ctx context.Context,
	tenantID, careRequestID, assessorID string,
	payload model.SubmitTriagePayload,
) (*model.TriageAssessment, error) {
	careReq, err := s.repo.GetCareRequestByID(ctx, tenantID, careRequestID)
	if err != nil || careReq == nil {
		return nil, errors.New("care request not found")
	}

	acuity := s.EvaluateAcuity(payload)

	triage := &model.TriageAssessment{
		CareRequestID:   careRequestID,
		PatientID:       careReq.PatientID,
		AssessorID:      &assessorID,
		AcuityLevel:     acuity,
		SystolicBP:      payload.SystolicBP,
		DiastolicBP:     payload.DiastolicBP,
		PulseRate:       payload.PulseRate,
		Temperature:     payload.Temperature,
		SpO2:            payload.SpO2,
		RespiratoryRate: payload.RespiratoryRate,
		PainScore:       payload.PainScore,
		TriageNotes:     payload.TriageNotes,
	}

	if err := s.repo.SaveTriageAssessment(ctx, triage); err != nil {
		return nil, fmt.Errorf("failed to save triage assessment: %w", err)
	}

	s.server.Logger.Info().
		Str("careRequestId", careRequestID).
		Str("acuityLevel", acuity).
		Msg("Clinical triage vitals recorded")

	return triage, nil
}

// MatchProviders scores provider candidates for a care request
func (s *CareRequestService) MatchProviders(ctx context.Context, tenantID, careRequestID string) (*model.ProviderMatchingResponse, error) {
	careReq, err := s.repo.GetCareRequestByID(ctx, tenantID, careRequestID)
	if err != nil || careReq == nil {
		return nil, errors.New("care request not found")
	}

	providers, err := s.repo.ListProviderProfiles(ctx, tenantID, careReq.ServiceType)
	if err != nil {
		return nil, err
	}

	response := &model.ProviderMatchingResponse{
		CareRequestID: careRequestID,
		RequiredMode:  careReq.PreferredMode,
		Specialty:     careReq.ServiceType,
		Candidates:    make([]model.MatchedProviderCandidate, 0),
	}

	for _, p := range providers {
		score := 0
		var reasons []string

		// 1. Status Check (+20 pts)
		if p.Status == "ON_DUTY" {
			score += 20
			reasons = append(reasons, "On Duty")
		} else if p.Status == "BUSY" {
			score += 10
			reasons = append(reasons, "Busy (Accepting Queue)")
		}

		// 2. Mode Compatibility (+30 pts)
		if (careReq.PreferredMode == "VIDEO" || careReq.PreferredMode == "AUDIO") && p.TelehealthEnabled {
			score += 30
			reasons = append(reasons, "Telehealth Certified")
		} else if careReq.PreferredMode == "IN_PERSON" && p.InPersonEnabled {
			score += 30
			reasons = append(reasons, "In-Person Clinic Available")
		}

		// 3. Specialty Fit (+30 pts)
		if p.SpecialtyCode == careReq.ServiceType || careReq.ServiceType == "GENERAL_CONSULTATION" {
			score += 30
			reasons = append(reasons, "Clinical Specialty Match")
		}

		// 4. Queue Capacity Load (+20 pts)
		if p.MaxActiveQueue > 0 {
			remainingCap := p.MaxActiveQueue - p.CurrentActiveQueue
			if remainingCap > 0 {
				ratio := float64(remainingCap) / float64(p.MaxActiveQueue)
				loadPoints := int(ratio * 20.0)
				score += loadPoints
				reasons = append(reasons, fmt.Sprintf("Queue Capacity Available (%d slots)", remainingCap))
			}
		}

		if score > 100 {
			score = 100
		}

		provName := "Consulting Doctor"
		if p.ProviderName != nil && *p.ProviderName != "" {
			provName = *p.ProviderName
		}

		response.Candidates = append(response.Candidates, model.MatchedProviderCandidate{
			ProviderID:         p.ID,
			UserID:             p.UserID,
			ProviderName:       provName,
			SpecialtyCode:      p.SpecialtyCode,
			Status:             p.Status,
			CurrentActiveQueue: p.CurrentActiveQueue,
			MaxActiveQueue:     p.MaxActiveQueue,
			MatchScore:         score,
			MatchingReasons:    reasons,
		})
	}

	return response, nil
}

// AssignProvider reserves provider and advances Care Journey to CONSULTATION
func (s *CareRequestService) AssignProvider(ctx context.Context, tenantID, careRequestID, providerID string) error {
	careReq, err := s.repo.GetCareRequestByID(ctx, tenantID, careRequestID)
	if err != nil || careReq == nil {
		return errors.New("care request not found")
	}

	if err := s.repo.AssignProviderToCareRequest(ctx, tenantID, careRequestID, providerID); err != nil {
		return fmt.Errorf("failed to assign provider: %w", err)
	}

	// Advance Care Journey Milestone to CONSULTATION
	now := time.Now()
	desc := "Provider matched and consultation room prepared"
	_ = s.repo.CreateCareJourneyMilestone(ctx, &model.CareJourneyMilestone{
		TenantID:    tenantID,
		PatientID:   careReq.PatientID,
		StageCode:   "CONSULTATION",
		Title:       "Clinical Doctor Consultation",
		Description: &desc,
		Status:      "IN_PROGRESS",
		CreatedAt:   now,
	})

	return nil
}
