package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/encounter/model"
	"github.com/golangnigeria/curexal/internal/modules/encounter/repository"
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

// StartEncounter initiates a consultation session (unified for in-person and telehealth)
func (s *EncounterService) StartEncounter(
	ctx context.Context,
	tenantID string,
	payload model.StartEncounterPayload,
) (*model.Encounter, error) {
	if payload.PatientID == "" || payload.ProviderID == "" {
		return nil, errors.New("patientId and providerId are required")
	}

	channel := payload.EncounterChannel
	if channel == "" {
		channel = "in_person"
	}

	enc := &model.Encounter{
		TenantID:         tenantID,
		CareRequestID:    payload.CareRequestID,
		AppointmentID:    payload.AppointmentID,
		PatientID:        payload.PatientID,
		ProviderID:       payload.ProviderID,
		EncounterChannel: channel,
		Status:           "in_progress",
		ChiefComplaint:   payload.ChiefComplaint,
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
		fmt.Sprintf("Consultation active via %s", channel),
	)

	s.server.Logger.Info().
		Str("encounterId", enc.ID).
		Str("patientId", payload.PatientID).
		Str("channel", channel).
		Msg("Clinical encounter started")

	return enc, nil
}

// GetEncounterByID retrieves encounter by ID
func (s *EncounterService) GetEncounterByID(ctx context.Context, tenantID, encounterID string) (*model.Encounter, error) {
	return s.repo.GetEncounterByID(ctx, tenantID, encounterID)
}

// ListEncounters retrieves a list of encounters with optional filters
func (s *EncounterService) ListEncounters(
	ctx context.Context,
	tenantID string,
	status *string,
	providerID *string,
	patientID *string,
) ([]model.Encounter, error) {
	return s.repo.ListEncounters(ctx, tenantID, status, providerID, patientID)
}

// SaveSOAPNotes saves physician clinical notes and primary diagnosis
func (s *EncounterService) SaveSOAPNotes(
	ctx context.Context,
	tenantID, encounterID string,
	payload model.UpdateSOAPNotesPayload,
) error {
	return s.repo.UpdateSOAPNotes(ctx, tenantID, encounterID, payload)
}

// AddDiagnosis records an ICD-10 diagnosis on the encounter
func (s *EncounterService) AddDiagnosis(
	ctx context.Context,
	tenantID, encounterID string,
	payload model.AddDiagnosisPayload,
) (*model.Diagnosis, error) {
	enc, err := s.repo.GetEncounterByID(ctx, tenantID, encounterID)
	if err != nil || enc == nil {
		return nil, errors.New("encounter not found")
	}

	clinicalStatus := payload.ClinicalStatus
	if clinicalStatus == "" {
		clinicalStatus = "ACTIVE"
	}

	verificationStatus := payload.VerificationStatus
	if verificationStatus == "" {
		verificationStatus = "CONFIRMED"
	}

	diag := &model.Diagnosis{
		EncounterID:        encounterID,
		PatientID:          enc.PatientID,
		ICD10Code:          payload.ICD10Code,
		ICD10Title:         payload.ICD10Title,
		IsPrimary:          payload.IsPrimary,
		ClinicalStatus:     clinicalStatus,
		VerificationStatus: verificationStatus,
		Notes:              payload.Notes,
	}

	if err := s.repo.AddDiagnosis(ctx, diag); err != nil {
		return nil, fmt.Errorf("failed to record diagnosis: %w", err)
	}

	return diag, nil
}

// ListDiagnoses returns all diagnoses recorded for an encounter
func (s *EncounterService) ListDiagnoses(ctx context.Context, encounterID string) ([]model.Diagnosis, error) {
	return s.repo.ListDiagnoses(ctx, encounterID)
}

// CreatePrescription creates an e-prescription order with items
func (s *EncounterService) CreatePrescription(
	ctx context.Context,
	tenantID, encounterID, prescriberID string,
	payload model.CreatePrescriptionPayload,
) (*model.Prescription, error) {
	enc, err := s.repo.GetEncounterByID(ctx, tenantID, encounterID)
	if err != nil || enc == nil {
		return nil, errors.New("encounter not found")
	}

	if len(payload.Items) == 0 {
		return nil, errors.New("at least one prescription item is required")
	}

	rx := &model.Prescription{
		EncounterID:  encounterID,
		PatientID:    enc.PatientID,
		PrescriberID: prescriberID,
		Status:       "PENDING_DISPENSE",
		Notes:        payload.Notes,
	}

	for _, itm := range payload.Items {
		rx.Items = append(rx.Items, model.PrescriptionItem{
			DrugName:           itm.DrugName,
			DosageForm:         itm.DosageForm,
			Strength:           itm.Strength,
			Route:              itm.Route,
			Frequency:          itm.Frequency,
			DurationDays:       itm.DurationDays,
			QuantityPrescribed: itm.QuantityPrescribed,
			Instructions:       itm.Instructions,
		})
	}

	if err := s.repo.CreatePrescription(ctx, rx); err != nil {
		return nil, fmt.Errorf("failed to create prescription: %w", err)
	}

	// Update PHARMACY_DISPENSE milestone in Care Journey
	_ = s.repo.CreateOrUpdateMilestone(
		ctx,
		tenantID,
		enc.PatientID,
		enc.ID,
		"PHARMACY_DISPENSE",
		"Pharmacy Dispensary & Fulfillment",
		"IN_PROGRESS",
		fmt.Sprintf("%d Medication(s) e-prescribed for dispensary", len(rx.Items)),
	)

	return rx, nil
}

// ListPrescriptions returns all prescriptions for an encounter
func (s *EncounterService) ListPrescriptions(ctx context.Context, encounterID string) ([]model.Prescription, error) {
	return s.repo.ListPrescriptions(ctx, encounterID)
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

// CompleteEncounter finalizes consultation, advances Care Journey to settlement,
// and automatically creates an invoice for cashier POS checkout.
func (s *EncounterService) CompleteEncounter(
	ctx context.Context,
	tenantID, encounterID string,
	consultationFee float64,
) (*model.CompleteEncounterResponse, error) {
	enc, err := s.repo.GetEncounterByID(ctx, tenantID, encounterID)
	if err != nil || enc == nil {
		return nil, errors.New("encounter not found")
	}

	if err := s.repo.CompleteEncounter(ctx, tenantID, encounterID); err != nil {
		return nil, fmt.Errorf("failed to complete encounter: %w", err)
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
		"Clinical consultation concluded and signed",
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
		"Billing reconciliation and cashier settlement",
	)

	// Auto-generate invoice for cashier settlement
	fee := consultationFee
	if fee <= 0 {
		fee = 5000.00
	}
	invID, invNum, invErr := s.repo.AutoGenerateEncounterInvoice(ctx, tenantID, enc.PatientID, encounterID, fee)
	if invErr != nil {
		s.server.Logger.Warn().Err(invErr).Str("encounterId", encounterID).Msg("Failed to auto-generate invoice upon encounter completion")
	}

	res := &model.CompleteEncounterResponse{
		EncounterID: encounterID,
		Status:      "signed",
		TotalAmount: fee,
		SignedAt:    time.Now().UTC().Format(time.RFC3339),
	}
	if invID != "" {
		res.InvoiceID = &invID
		res.InvoiceNumber = &invNum
	}

	return res, nil
}
