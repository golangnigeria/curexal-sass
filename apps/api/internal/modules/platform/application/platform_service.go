package application

import (
	"github.com/golangnigeria/curexal/internal/modules/platform/domain"
)

type PlatformApplicationService struct {
	platformDomain *domain.PlatformDomainProvider
}

func NewPlatformApplicationService() *PlatformApplicationService {
	return &PlatformApplicationService{
		platformDomain: domain.NewPlatformDomainProvider(),
	}
}

func (s *PlatformApplicationService) ResolveContext(isStaff bool, role string, email string) bool {
	return s.platformDomain.IsPlatformContext(isStaff, role, email)
}

func (s *PlatformApplicationService) GetNavigation() []domain.PlatformNavigationItem {
	return s.platformDomain.GetPlatformNavigation()
}
