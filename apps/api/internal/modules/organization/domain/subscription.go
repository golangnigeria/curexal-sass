package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	TenantID                 uuid.UUID       `json:"tenantId"                 db:"tenant_id"`
	Plan                     string          `json:"plan"                     db:"plan"`
	PlanID                   *string         `json:"planId"                   db:"plan_id"`
	Status                   string          `json:"status"                   db:"status"`
	FeatureSnapshot          json.RawMessage `json:"featureSnapshot"          db:"feature_snapshot"`
	PaystackSubscriptionCode *string         `json:"paystackSubscriptionCode" db:"paystack_subscription_code"`
	PaystackCustomerCode     *string         `json:"paystackCustomerCode"     db:"paystack_customer_code"`
	TrialStartedAt           time.Time       `json:"trialStartedAt"           db:"trial_started_at"`
	TrialEndsAt              time.Time       `json:"trialEndsAt"              db:"trial_ends_at"`
	CurrentPeriodStart       time.Time       `json:"currentPeriodStart"       db:"current_period_start"`
	CurrentPeriodEnd         *time.Time      `json:"currentPeriodEnd"         db:"current_period_end"`
	GracePeriodEndsAt        *time.Time      `json:"gracePeriodEndsAt"        db:"grace_period_ends_at"`
	CanceledAt               *time.Time      `json:"canceledAt"               db:"canceled_at"`
}
