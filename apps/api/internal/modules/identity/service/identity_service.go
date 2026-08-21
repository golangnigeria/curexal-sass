package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/identity/repository"
	platformAuth "github.com/golangnigeria/curexal/internal/kernel/auth"
	"github.com/golangnigeria/curexal/internal/shared/errs"
)

// IdentityService encapsulates identity provider interaction, user profile projection, and lazy identity linking.
type IdentityService struct {
	provider platformAuth.IdentityProvider
	userRepo *repository.UserRepository
}

// NewIdentityService initializes a new IdentityService instance.
func NewIdentityService(provider platformAuth.IdentityProvider, userRepo *repository.UserRepository) *IdentityService {
	return &IdentityService{
		provider: provider,
		userRepo: userRepo,
	}
}

// Provider returns the underlying IdentityProvider instance.
func (s *IdentityService) Provider() platformAuth.IdentityProvider {
	return s.provider
}

// ResolveUserSession validates an incoming session token/cookie with the IdentityProvider and maps it to a canonical Curexal User profile.
func (s *IdentityService) ResolveUserSession(ctx context.Context, sessionToken string) (*model.User, *platformAuth.IdentitySession, error) {
	if s == nil || s.userRepo == nil {
		return nil, nil, errors.New("identity service uninitialized")
	}

	if s.provider != nil {
		sess, err := s.provider.Authenticate(ctx, sessionToken)
		if err == nil && sess != nil && sess.Active {
			u, errUser := s.userRepo.GetByID(ctx, sess.IdentityID)
			if errUser == nil && u != nil {
				return u, sess, nil
			}
			if sess.Traits.Email != "" {
				existingUser, errEmail := s.userRepo.GetByEmail(ctx, sess.Traits.Email)
				if errEmail == nil && existingUser != nil {
					return existingUser, sess, nil
				}
			}
		}
	}

	return nil, nil, errs.NewUnauthorizedError("user is not authenticated")
}

// ProvisionIdentity registers a user in the IdentityProvider and provisions a corresponding Curexal User domain profile.
func (s *IdentityService) ProvisionIdentity(ctx context.Context, user *model.User, password string) (string, error) {
	if user == nil {
		return "", errors.New("user cannot be nil")
	}

	if s.provider != nil {
		traits := platformAuth.IdentityTraits{
			Email: user.Email,
			Name:  user.Name,
		}
		if user.Image != nil {
			traits.AvatarURL = *user.Image
		}

		providerID, err := s.provider.CreateIdentity(ctx, &platformAuth.CreateIdentityParams{
			Traits:   traits,
			Password: password,
		})
		if err != nil {
			return "", fmt.Errorf("failed to provision provider identity: %w", err)
		}
		return providerID, nil
	}

	return user.ID, nil
}
