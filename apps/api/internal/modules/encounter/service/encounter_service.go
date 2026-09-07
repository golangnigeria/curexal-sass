package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/golangnigeria/curexal/internal/modules/encounter/model"
	"github.com/golangnigeria/curexal/internal/modules/encounter/repository"
	"github.com/golangnigeria/curexal/internal/kernel/server"
)

type EncounterService struct {
	server *server.Server
	repo   *repository.EncounterRepository
}

func NewEncounterService(s *server.Server, repo *repository.EncounterRepository) *EncounterService {
	return &EncounterService{
		server: s,
		repo:   repo,
	}
}

// StartEncounter initiates a consultation session
func (s *EncounterService) StartEncounter(
	ctx context.Context,
	tenantID string,
	payload model.StartEncounterPayload,
) (*model.Encounter, error) {
	if payload.PatientID == "" || payload.ProviderID == "" {
		return nil, errors.New("patientId and providerId are required")
	}

	enc := &model.Encounter{
		TenantID:       tenantID,
		CareRequestID:  payload.CareRequestID,
		PatientID:      payload.PatientID,
		ProviderID:     payload.ProviderID,
		EncounterType:  payload.EncounterType,
		Mode:           payload.Mode,
		Status:         "IN_PROGRESS",
		ChiefComplaint: payload.ChiefComplaint,
	}

	if err := s.repo.CreateEncounter(ctx, enc); err != nil {
		return nil, fmt.Errorf("failed to start encounter: %w", err)
	}

	// Update Care Journey to CONSULTATION (IN_PROGRESS)
	_ = s.repo.CreateOrUpdateMilestone(
		ctx,
		tenantID,
		payload.PatientID,
		enc.ID,
		"CONSULTATION",
		"Doctor Consultation Session",
		"IN_PROGRESS",
		fmt.Sprintf("Consultation active via %s", payload.Mode),
	)

	s.server.Logger.Info().
		Str("encounterId", enc.ID).
		Str("patientId", payload.PatientID).
		Str("mode", payload.Mode).
		Msg("Clinical encounter started")

	return enc, nil
}

// GetEncounterByID retrieves encounter by ID
func (s *EncounterService) GetEncounterByID(ctx context.Context, tenantID, encounterID string) (*model.Encounter, error) {
	return s.repo.GetEncounterByID(ctx, tenantID, encounterID)
}

// SaveSOAPNotes saves physician clinical notes and primary diagnosis
func (s *EncounterService) SaveSOAPNotes(
	ctx context.Context,
	tenantID, encounterID string,
	payload model.UpdateSOAPNotesPayload,
) error {
	return s.repo.UpdateSOAPNotes(ctx, tenantID, encounterID, payload)
}

// DispatchOrders routes clinical orders and automatically updates Care Journey milestones
func (s *EncounterService) DispatchOrders(
	ctx context.Context,
	tenantID, encounterID string,
	payload model.DispatchOrdersPayload,
) error {
	enc, err := s.repo.GetEncounterByID(ctx, tenantID, encounterID)
	if err != nil || enc == nil {
		return errors.New("encounter not found")
	}

	// 1. If Lab Orders present, launch LAB_WORKLIST milestone
	if len(payload.LabOrders) > 0 {
		desc := fmt.Sprintf("%d Diagnostic Lab Test(s) ordered", len(payload.LabOrders))
		_ = s.repo.CreateOrUpdateMilestone(
			ctx,
			tenantID,
			enc.PatientID,
			enc.ID,
			"LAB_WORKLIST",
			"Diagnostic Pathology & Laboratory",
			"IN_PROGRESS",
			desc,
		)
	}

	// 2. If Radiology Orders present, launch RADIOLOGY_STUDY milestone
	if len(payload.RadiologyOrders) > 0 {
		desc := fmt.Sprintf("%d Radiology Imaging Study(ies) ordered", len(payload.RadiologyOrders))
		_ = s.repo.CreateOrUpdateMilestone(
			ctx,
			tenantID,
			enc.PatientID,
			enc.ID,
			"RADIOLOGY_STUDY",
			"Radiology & Imaging Investigations",
			"IN_PROGRESS",
			desc,
		)
	}

	// 3. If Prescription items present, launch PHARMACY_DISPENSE milestone
	if len(payload.PrescriptionList) > 0 {
		desc := fmt.Sprintf("%d Medication(s) e-prescribed for dispensary", len(payload.PrescriptionList))
		_ = s.repo.CreateOrUpdateMilestone(
			ctx,
			tenantID,
			enc.PatientID,
			enc.ID,
			"PHARMACY_DISPENSE",
			"Pharmacy Dispensary & Fulfillment",
			"IN_PROGRESS",
			desc,
		)
	}

	s.server.Logger.Info().
		Str("encounterId", encounterID).
		Int("labOrders", len(payload.LabOrders)).
		Int("radiologyOrders", len(payload.RadiologyOrders)).
		Int("prescriptions", len(payload.PrescriptionList)).
		Msg("Clinical diagnostic & medication orders dispatched")

	return nil
}

// CompleteEncounter finalizes consultation and advances Care Journey to settlement
func (s *EncounterService) CompleteEncounter(ctx context.Context, tenantID, encounterID string) error {
	enc, err := s.repo.GetEncounterByID(ctx, tenantID, encounterID)
	if err != nil || enc == nil {
		return errors.New("encounter not found")
	}

	if err := s.repo.CompleteEncounter(ctx, tenantID, encounterID); err != nil {
		return fmt.Errorf("failed to complete encounter: %w", err)
	}

	// Mark CONSULTATION as COMPLETED
	_ = s.repo.CreateOrUpdateMilestone(
		ctx,
		tenantID,
		enc.PatientID,
		enc.ID,
		"CONSULTATION",
		"Doctor Consultation Session",
		"COMPLETED",
		"Clinical consultation concluded",
	)

	// Advance to SETTLEMENT
	_ = s.repo.CreateOrUpdateMilestone(
		ctx,
		tenantID,
		enc.PatientID,
		enc.ID,
		"SETTLEMENT",
		"Care Summary & Final Settlement",
		"IN_PROGRESS",
		"Billing reconciliation and discharge summary",
	)

	return nil
}
