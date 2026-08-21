package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	identityModel "github.com/golangnigeria/curexal/internal/modules/identity/model"
	patientModel "github.com/golangnigeria/curexal/internal/modules/patient/model"
	patientRepo "github.com/golangnigeria/curexal/internal/modules/patient/repository"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	crypto "github.com/golangnigeria/curexal/internal/shared/crypto"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type UserIdentityRepo interface {
	GetByEmail(ctx context.Context, email string) (*identityModel.User, error)
	CreateUser(ctx context.Context, u *identityModel.User, passwordHash string) error
	CreateVerificationToken(ctx context.Context, token, email, tokenType string, expiresAt time.Time) error
}

type PatientService struct {
	server      *server.Server
	patientRepo *patientRepo.PatientRepository
	userRepo    UserIdentityRepo
}

func NewPatientService(s *server.Server, patientRepo *patientRepo.PatientRepository, userRepo UserIdentityRepo) *PatientService {
	return &PatientService{
		server:      s,
		patientRepo: patientRepo,
		userRepo:    userRepo,
	}
}

func (s *PatientService) RegisterPatient(ctx context.Context, payload patientModel.RegisterPatientPayload, origin string) error {
	email := strings.ToLower(strings.TrimSpace(payload.Email))

	return s.server.DB.RunInTx(ctx, func(txCtx context.Context) error {
		// 1. Check if user already exists
		existingUser, err := s.userRepo.GetByEmail(txCtx, email)
		var userID string
		if err == nil && existingUser != nil {
			userID = existingUser.ID

			// Check if patient profile already exists
			exists, _, err := s.patientRepo.ProfileExists(txCtx, userID)
			if err != nil {
				return fmt.Errorf("error checking patient profile: %w", err)
			}
			if exists {
				return errors.New("An account with this email address already exists. Please sign in.")
			}
		} else if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("database query error: %w", err)
		} else {
			// Create new identity user
			userID = uuid.New().String()
			hash, err := crypto.HashPassword(payload.Password)
			if err != nil {
				return fmt.Errorf("failed to hash password: %w", err)
			}

			newUser := &identityModel.User{
				ID:              userID,
				Name:            strings.TrimSpace(payload.Name),
				Email:           email,
				EmailVerified:   false,
				IsPlatformAdmin: false,
			}

			if err := s.userRepo.CreateUser(txCtx, newUser, hash); err != nil {
				return fmt.Errorf("failed to create user identity: %w", err)
			}
		}

		// 2. Create patient profile persona
		if err := s.patientRepo.CreateProfile(txCtx, userID, payload.Phone); err != nil {
			return fmt.Errorf("failed to create patient profile: %w", err)
		}

		// 3. Create alphanumeric email verification code & send verification email
		code, err := crypto.GenerateAlphanumericCode(6)
		if err != nil {
			return fmt.Errorf("failed to generate verification code: %w", err)
		}
		expiresAt := time.Now().Add(24 * time.Hour)

		if err := s.userRepo.CreateVerificationToken(txCtx, code, email, "email_verify", expiresAt); err != nil {
			return fmt.Errorf("failed to store verification code: %w", err)
		}

		s.server.Logger.Info().
			Str("email", email).
			Str("verificationCode", code).
			Msg("Patient verification code generated")

		if s.server.Mailer != nil {
			if err := s.server.Mailer.SendVerificationEmail(txCtx, email, payload.Name, code); err != nil {
				s.server.Logger.Error().Err(err).Str("email", email).Msg("failed to send patient verification email via Resend")
			}
		}

		return nil
	})
}

func (s *PatientService) LoadPatientContext(ctx context.Context, userID string) (*identityModel.PatientContext, error) {
	return s.patientRepo.GetPatientContext(ctx, userID)
}

func (s *PatientService) GetProfile(ctx context.Context, userID string) (*patientModel.PatientProfile, error) {
	return s.patientRepo.GetProfileByUserID(ctx, userID)
}

func (s *PatientService) UpdateProfile(ctx context.Context, userID string, payload patientModel.UpdatePatientProfilePayload) error {
	return s.patientRepo.UpdateProfile(ctx, userID, payload)
}
