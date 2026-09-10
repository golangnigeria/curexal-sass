package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/operations/model"
	"github.com/golangnigeria/curexal/internal/modules/operations/repository"
	"github.com/google/uuid"
)

var aptSequence uint64

func init() {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	aptSequence = (uint64(b[0])<<56 | uint64(b[1])<<48 | uint64(b[2])<<40 | uint64(b[3])<<32 |
		uint64(b[4])<<24 | uint64(b[5])<<16 | uint64(b[6])<<8 | uint64(b[7]))
}

type AppointmentService struct {
	server *server.Server
	repo   *repository.AppointmentRepository
}

func NewAppointmentService(s *server.Server, repo *repository.AppointmentRepository) *AppointmentService {
	return &AppointmentService{
		server: s,
		repo:   repo,
	}
}

// GenerateAppointmentNumber produces formatted, collision-free numbers (e.g. APT-2026-00481)
func (s *AppointmentService) GenerateAppointmentNumber() string {
	seq := atomic.AddUint64(&aptSequence, 1)
	val := (seq % 90000) + 10000
	year := time.Now().Year()
	return fmt.Sprintf("APT-%d-%05d", year, val)
}

// IsValidDeliveryChannel validates whether a given delivery channel is supported
func (s *AppointmentService) IsValidDeliveryChannel(channel string) bool {
	switch channel {
	case "in_person", "video", "telephone", "secure_message":
		return true
	default:
		return false
	}
}

// CreateAppointment schedules a new appointment with double-booking prevention
func (s *AppointmentService) CreateAppointment(
	ctx context.Context,
	tenantID string,
	createdBy *string,
	req model.CreateAppointmentRequest,
) (*model.Appointment, error) {
	if tenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	if strings.TrimSpace(req.PatientID) == "" {
		return nil, errors.New("patientId is required")
	}
	if strings.TrimSpace(req.ProviderID) == "" {
		return nil, errors.New("providerId is required")
	}

	startTime, err := time.Parse(time.RFC3339, strings.TrimSpace(req.StartTime))
	if err != nil {
		return nil, fmt.Errorf("invalid startTime format (RFC3339 expected): %w", err)
	}

	endTime, err := time.Parse(time.RFC3339, strings.TrimSpace(req.EndTime))
	if err != nil {
		return nil, fmt.Errorf("invalid endTime format (RFC3339 expected): %w", err)
	}

	if !endTime.After(startTime) {
		return nil, errors.New("endTime must be after startTime")
	}

	channel := strings.ToLower(strings.TrimSpace(req.DeliveryChannel))
	if channel == "" {
		channel = "in_person"
	}
	validChannels := map[string]bool{
		"in_person":      true,
		"video":          true,
		"telephone":      true,
		"secure_message": true,
	}
	if !validChannels[channel] {
		return nil, fmt.Errorf("invalid deliveryChannel '%s'; must be in_person, video, telephone, or secure_message", channel)
	}

	// 1. Anti-Double-Booking conflict check
	conflict, err := s.repo.CheckProviderConflict(ctx, req.ProviderID, startTime, endTime, nil)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, errors.New("provider has an existing booking overlapping with this time slot")
	}

	serviceType := req.ServiceType
	if strings.TrimSpace(serviceType) == "" {
		serviceType = "CONSULTATION"
	}

	aptID := uuid.New().String()
	aptNumber := s.GenerateAppointmentNumber()

	var meetingURL *string
	if channel == "video" {
		url := fmt.Sprintf("https://telehealth.curexal.com/join/%s?token=jwt_secure_room_token", aptID)
		meetingURL = &url
	}

	apt := &model.Appointment{
		ID:                 aptID,
		TenantID:           tenantID,
		PatientID:          req.PatientID,
		ProviderID:         req.ProviderID,
		AppointmentNumber:  aptNumber,
		ServiceType:        serviceType,
		DeliveryChannel:    channel,
		Status:             "BOOKED",
		StartTime:          startTime,
		EndTime:            endTime,
		ReasonForVisit:     req.ReasonForVisit,
		VirtualMeetingURL:  meetingURL,
		CreatedBy:          createdBy,
	}

	if err := s.repo.CreateAppointment(ctx, apt); err != nil {
		return nil, fmt.Errorf("failed to save appointment: %w", err)
	}

	return s.repo.GetAppointmentByID(ctx, tenantID, aptID)
}

// ListAppointments retrieves schedule
func (s *AppointmentService) ListAppointments(
	ctx context.Context,
	tenantID string,
	filter model.AppointmentFilter,
) ([]model.Appointment, error) {
	if tenantID == "" {
		return nil, errors.New("tenantId is required")
	}
	return s.repo.ListAppointments(ctx, tenantID, filter)
}

// GetAppointmentByID retrieves single appointment
func (s *AppointmentService) GetAppointmentByID(
	ctx context.Context,
	tenantID, id string,
) (*model.Appointment, error) {
	if tenantID == "" || id == "" {
		return nil, errors.New("tenantId and id are required")
	}
	return s.repo.GetAppointmentByID(ctx, tenantID, id)
}

// UpdateAppointmentStatus updates booking status
func (s *AppointmentService) UpdateAppointmentStatus(
	ctx context.Context,
	tenantID, id string,
	req model.UpdateAppointmentStatusRequest,
) error {
	if tenantID == "" || id == "" {
		return errors.New("tenantId and id are required")
	}
	return s.repo.UpdateAppointmentStatus(ctx, tenantID, id, req.Status, req.CancellationReason)
}
