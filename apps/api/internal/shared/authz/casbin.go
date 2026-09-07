package authz

import (
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
)

const defaultModelText = `
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g = _, _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && (r.dom == p.dom || p.dom == "*") && (r.obj == p.obj || p.obj == "*" || keyMatch2(r.obj, p.obj)) && (r.act == p.act || p.act == "*")
`

type CasbinEngine struct {
	enforcer *casbin.Enforcer
	mu       sync.RWMutex
}

// NewCasbinEngine initializes a new domain-aware Casbin enforcer.
func NewCasbinEngine() (*CasbinEngine, error) {
	m, err := model.NewModelFromString(defaultModelText)
	if err != nil {
		return nil, fmt.Errorf("failed to load casbin model: %w", err)
	}

	e, err := casbin.NewEnforcer(m)
	if err != nil {
		return nil, fmt.Errorf("failed to create casbin enforcer: %w", err)
	}

	engine := &CasbinEngine{
		enforcer: e,
	}

	// Seed canonical policies
	if err := engine.seedCanonicalPolicies(); err != nil {
		return nil, fmt.Errorf("failed to seed casbin policies: %w", err)
	}

	return engine, nil
}

func (e *CasbinEngine) seedCanonicalPolicies() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	policies := [][]string{
		// ---------------------------------------------------------------------
		// Platform Domain Policies (dom: "platform")
		// ---------------------------------------------------------------------
		{"super_admin", "platform", "platform.*", "*"},
		{"platform_admin", "platform", "platform.*", "*"},
		{"platform_staff", "platform", "platform.diagnostics", "read"},

		// ---------------------------------------------------------------------
		// Organization Domain Policies (dom: "*")
		// ---------------------------------------------------------------------
		{"owner", "*", "organization.governance", "*"},
		{"owner", "*", "organization.billing", "*"},
		{"owner", "*", "organization.branches", "*"},
		{"owner", "*", "organization.staff", "*"},
		{"owner", "*", "organization.compliance", "*"},
		{"owner", "*", "branch.operations.overview", "read"},

		{"org_admin", "*", "organization.branches", "*"},
		{"org_admin", "*", "organization.staff", "*"},
		{"org_admin", "*", "organization.compliance", "*"},
		{"org_admin", "*", "branch.operations.overview", "read"},

		{"org_finance_manager", "*", "organization.billing", "*"},
		{"org_regional_manager", "*", "organization.branches", "read"},
		{"org_quality_manager", "*", "organization.compliance", "*"},
		{"org_hr_manager", "*", "organization.staff", "*"},

		// ---------------------------------------------------------------------
		// Shared Core / MPI Policies across Branch Domains (dom: "*")
		// ---------------------------------------------------------------------
		{"receptionist", "*", "core.patient", "create"},
		{"receptionist", "*", "core.patient", "read"},
		{"receptionist", "*", "core.patient", "search"},
		{"receptionist", "*", "core.appointment", "*"},
		{"receptionist", "*", "hms.checkin", "create"},

		{"scientist", "*", "core.patient", "create"},
		{"scientist", "*", "core.patient", "read"},
		{"scientist", "*", "core.patient", "search"},

		{"phlebotomist", "*", "core.patient", "read"},
		{"phlebotomist", "*", "core.patient", "search"},

		{"pharmacist", "*", "core.patient", "read"},
		{"pharmacist", "*", "core.patient", "search"},

		{"cashier", "*", "core.patient", "read"},
		{"cashier", "*", "core.patient", "search"},
		{"cashier", "*", "billing.payment", "create"},
		{"cashier", "*", "billing.invoice", "read"},

		{"branch_manager", "*", "branch.operations", "*"},
		{"branch_manager", "*", "branch.staff", "read"},
		{"branch_manager", "*", "branch.reports", "read"},
		{"branch_manager", "*", "billing.branch", "read"},

		// ---------------------------------------------------------------------
		// HMS / Clinical Doctor & Nurse Policies
		// ---------------------------------------------------------------------
		{"doctor", "*", "core.patient", "read"},
		{"doctor", "*", "core.patient", "search"},
		{"doctor", "*", "hms.consultation", "*"},
		{"doctor", "*", "hms.vitals", "*"},
		{"doctor", "*", "hms.diagnosis", "*"},
		{"doctor", "*", "hms.prescription", "*"},
		{"doctor", "*", "hms.lab_order", "*"},
		{"doctor", "*", "lis.order", "read"},
		{"doctor", "*", "lis.result", "read"},

		{"nurse", "*", "core.patient", "read"},
		{"nurse", "*", "core.patient", "search"},
		{"nurse", "*", "hms.vitals", "*"},
		{"nurse", "*", "hms.triage", "*"},
		{"nurse", "*", "hms.nursing_notes", "*"},
		{"nurse", "*", "hms.care_plan", "*"},

		// ---------------------------------------------------------------------
		// LIS / Laboratory Pathology Policies
		// ---------------------------------------------------------------------
		{"scientist", "*", "lis.order", "read"},
		{"scientist", "*", "lis.sample", "receive"},
		{"scientist", "*", "lis.sample", "accession"},
		{"scientist", "*", "lis.result", "create"},
		{"scientist", "*", "lis.result", "read"},
		{"scientist", "*", "lis.qc", "read"},

		{"phlebotomist", "*", "lis.sample", "collect"},
		{"phlebotomist", "*", "lis.sample", "label"},
		{"phlebotomist", "*", "lis.sample", "track"},

		{"lab_manager", "*", "lis.order", "*"},
		{"lab_manager", "*", "lis.sample", "*"},
		{"lab_manager", "*", "lis.result", "*"},
		{"lab_manager", "*", "lis.qc", "*"},
		{"lab_manager", "*", "lis.instrument", "*"},
		{"lab_manager", "*", "lis.settings", "*"},

		// ---------------------------------------------------------------------
		// Pharmacy Dispensary Policies
		// ---------------------------------------------------------------------
		{"pharmacist", "*", "pharmacy.prescription", "read"},
		{"pharmacist", "*", "pharmacy.dispense", "create"},
		{"pharmacist", "*", "pharmacy.inventory", "*"},
		{"pharmacist", "*", "pharmacy.medication", "read"},

		// ---------------------------------------------------------------------
		// Radiology RIS / PACS Policies
		// ---------------------------------------------------------------------
		{"radiologist", "*", "radiology.worklist", "read"},
		{"radiologist", "*", "radiology.study", "read"},
		{"radiologist", "*", "radiology.report", "*"},
		{"radiologist", "*", "radiology.pacs", "view"},
	}

	for _, p := range policies {
		if _, err := e.enforcer.AddPolicy(p[0], p[1], p[2], p[3]); err != nil {
			return err
		}
	}

	return nil
}

// Enforce evaluates a subject, domain, object, and action.
func (e *CasbinEngine) Enforce(sub, dom, obj, act string) (bool, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.enforcer.Enforce(sub, dom, obj, act)
}

// AddRoleForUserInDomain maps a user to a role within a specific domain.
func (e *CasbinEngine) AddRoleForUserInDomain(user, role, domain string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	_, err := e.enforcer.AddRoleForUserInDomain(user, role, domain)
	return err
}
