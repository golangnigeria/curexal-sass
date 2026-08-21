package application

import (
	"fmt"

	"github.com/golangnigeria/curexal/internal/modules/billing/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
)

type BillingPolicyRegistry struct {
	policies map[string]domain.BillingPolicy
}

func NewBillingPolicyRegistry() *BillingPolicyRegistry {
	return &BillingPolicyRegistry{
		policies: map[string]domain.BillingPolicy{
			"CASH":    &domain.CashBillingPolicy{},
			"PREPAID": &domain.PrepaidBillingPolicy{},
			"CREDIT":  &domain.CreditBillingPolicy{},
		},
	}
}

func (r *BillingPolicyRegistry) Resolve(policyType string) (domain.BillingPolicy, error) {
	policy, ok := r.policies[policyType]
	if !ok {
		return nil, fmt.Errorf("unknown billing policy: %s", policyType)
	}
	return policy, nil
}

type BillingApplicationService struct {
	server   *server.Server
	Registry *BillingPolicyRegistry
}

func NewBillingApplicationService(server *server.Server) *BillingApplicationService {
	return &BillingApplicationService{
		server:   server,
		Registry: NewBillingPolicyRegistry(),
	}
}
