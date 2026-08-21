package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type SubscriptionRepository struct {
	server *server.Server
}

func NewSubscriptionRepository(server *server.Server) *SubscriptionRepository {
	return &SubscriptionRepository{server: server}
}

const subscriptionSelectFields = `
	id, tenant_id, plan, plan_id, status, feature_snapshot, 
	paystack_subscription_code, paystack_customer_code, 
	trial_started_at, trial_ends_at, current_period_start, 
	current_period_end, grace_period_ends_at, canceled_at, 
	created_at, updated_at
`

func (r *SubscriptionRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := fmt.Sprintf(`SELECT %s FROM subscriptions WHERE id = @id`, subscriptionSelectFields)

	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get subscription query for id=%s: %w", id, err)
	}

	sub, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Subscription])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:subscriptions for id=%s: %w", id, err)
	}

	return &sub, nil
}

func (r *SubscriptionRepository) GetByTenantID(ctx context.Context, tenantID uuid.UUID) (*domain.Subscription, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := fmt.Sprintf(`SELECT %s FROM subscriptions WHERE tenant_id = @tenant_id ORDER BY created_at DESC LIMIT 1`, subscriptionSelectFields)

	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"tenant_id": tenantID})
	if err != nil {
		return nil, fmt.Errorf("failed to execute get subscription query for tenant_id=%s: %w", tenantID, err)
	}

	sub, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Subscription])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:subscriptions for tenant_id=%s: %w", tenantID, err)
	}

	return &sub, nil
}

func (r *SubscriptionRepository) Create(ctx context.Context, sub *domain.Subscription) (*domain.Subscription, error) {
	dbExec := r.server.DB.Conn(ctx)

	snapshot := []byte("{}")
	if len(sub.FeatureSnapshot) > 0 {
		snapshot = sub.FeatureSnapshot
	}

	stmt := fmt.Sprintf(`
		INSERT INTO subscriptions (
			tenant_id, plan, plan_id, status, feature_snapshot, 
			paystack_subscription_code, paystack_customer_code, 
			trial_started_at, trial_ends_at, current_period_start, 
			current_period_end, grace_period_ends_at
		)
		VALUES (
			@tenant_id, @plan, @plan_id, @status, @feature_snapshot, 
			@paystack_subscription_code, @paystack_customer_code, 
			@trial_started_at, @trial_ends_at, @current_period_start, 
			@current_period_end, @grace_period_ends_at
		)
		RETURNING %s
	`, subscriptionSelectFields)

	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{
		"tenant_id":                  sub.TenantID,
		"plan":                       sub.Plan,
		"plan_id":                    sub.PlanID,
		"status":                     sub.Status,
		"feature_snapshot":           snapshot,
		"paystack_subscription_code": sub.PaystackSubscriptionCode,
		"paystack_customer_code":     sub.PaystackCustomerCode,
		"trial_started_at":           sub.TrialStartedAt,
		"trial_ends_at":              sub.TrialEndsAt,
		"current_period_start":       sub.CurrentPeriodStart,
		"current_period_end":         sub.CurrentPeriodEnd,
		"grace_period_ends_at":       sub.GracePeriodEndsAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create subscription query: %w", err)
	}

	created, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Subscription])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:subscriptions: %w", err)
	}

	return &created, nil
}

func (r *SubscriptionRepository) Update(ctx context.Context, id uuid.UUID, sub *domain.Subscription) (*domain.Subscription, error) {
	dbExec := r.server.DB.Conn(ctx)
	setClauses := []string{}
	args := pgx.NamedArgs{
		"id": id,
	}

	if sub.Plan != "" {
		setClauses = append(setClauses, "plan = @plan")
		args["plan"] = sub.Plan
	}
	if sub.PlanID != nil {
		setClauses = append(setClauses, "plan_id = @plan_id")
		args["plan_id"] = sub.PlanID
	}
	if sub.Status != "" {
		setClauses = append(setClauses, "status = @status")
		args["status"] = sub.Status
	}

	if len(setClauses) == 0 {
		return nil, errors.New("no fields to update")
	}

	stmt := fmt.Sprintf("UPDATE subscriptions SET %s, updated_at = CURRENT_TIMESTAMP WHERE id = @id RETURNING %s", strings.Join(setClauses, ", "), subscriptionSelectFields)

	rows, err := dbExec.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to execute update subscription query for id=%s: %w", id, err)
	}

	updated, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Subscription])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:subscriptions for id=%s: %w", id, err)
	}

	return &updated, nil
}

func (r *SubscriptionRepository) Delete(ctx context.Context, id uuid.UUID) error {
	dbExec := r.server.DB.Conn(ctx)
	stmt := "DELETE FROM subscriptions WHERE id = @id"

	_, err := dbExec.Exec(ctx, stmt, pgx.NamedArgs{"id": id})
	if err != nil {
		return fmt.Errorf("failed to execute delete subscription query for id=%s: %w", id, err)
	}

	return nil
}
