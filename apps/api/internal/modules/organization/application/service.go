package application

import (
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
)

type OrganizationApplicationService struct {
	server     *server.Server
	orgRepo    domain.OrganizationRepository
	tenantRepo domain.TenantRepository
	subRepo    domain.SubscriptionRepository
	demoRepo   domain.DemoRepository
}

func NewOrganizationApplicationService(
	server *server.Server,
	orgRepo domain.OrganizationRepository,
	tenantRepo domain.TenantRepository,
	subRepo domain.SubscriptionRepository,
	demoRepo domain.DemoRepository,
) *OrganizationApplicationService {
	return &OrganizationApplicationService{
		server:     server,
		orgRepo:    orgRepo,
		tenantRepo: tenantRepo,
		subRepo:    subRepo,
		demoRepo:   demoRepo,
	}
}

func (s *OrganizationApplicationService) GetNavigation() []domain.OrganizationNavigationItem {
	return domain.NewOrganizationDomainProvider().GetOrganizationNavigation()
}

