package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/golangnigeria/curexal/internal/modules/platform/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type LaunchGateRepository struct {
	server *server.Server
}

func NewLaunchGateRepository(server *server.Server) *LaunchGateRepository {
	return &LaunchGateRepository{server: server}
}

func (r *LaunchGateRepository) SaveLaunchGateAudit(ctx context.Context, audit *domain.LaunchGateAudit, actorID *uuid.UUID) (*domain.LaunchGateAudit, error) {
	dbExec := r.server.DB.Conn(ctx)
	if audit.ID == uuid.Nil {
		audit.ID = uuid.New()
	}

	checksJSON := []byte("[]")
	if len(audit.CheckResults) > 0 {
		checksJSON = audit.CheckResults
	}

	var actorStr *string
	if actorID != nil {
		str := actorID.String()
		actorStr = &str
	}

	stmt := `
		INSERT INTO platform.launch_gate_audits (id, gate_name, status, check_results, evaluated_by)
		VALUES ($1, $2, $3, $4::jsonb, $5)
		RETURNING evaluated_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		audit.ID, audit.GateName, audit.Status, string(checksJSON), actorStr,
	).Scan(&audit.EvaluatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to save launch gate audit: %w", err)
	}

	audit.EvaluatedBy = actorID
	return audit, nil
}

func (r *LaunchGateRepository) GetLatestLaunchGateAudit(ctx context.Context, gateName string) (*domain.LaunchGateAudit, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, gate_name, status, check_results, evaluated_at, evaluated_by
		FROM platform.launch_gate_audits
		WHERE gate_name = $1
		ORDER BY evaluated_at DESC
		LIMIT 1
	`

	var (
		a            domain.LaunchGateAudit
		evalByString *string
	)
	err := dbExec.QueryRow(ctx, stmt, gateName).Scan(
		&a.ID, &a.GateName, &a.Status, &a.CheckResults, &a.EvaluatedAt, &evalByString,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query latest launch gate audit: %w", err)
	}

	if evalByString != nil && *evalByString != "" {
		if parsed, pErr := uuid.Parse(*evalByString); pErr == nil {
			a.EvaluatedBy = &parsed
		}
	}

	return &a, nil
}

func (r *LaunchGateRepository) SaveSystemHealthMetric(ctx context.Context, metric *domain.SystemHealthMetric) (*domain.SystemHealthMetric, error) {
	dbExec := r.server.DB.Conn(ctx)
	if metric.ID == uuid.Nil {
		metric.ID = uuid.New()
	}

	metricsJSON := []byte("{}")
	if len(metric.Metrics) > 0 {
		metricsJSON = metric.Metrics
	}

	stmt := `
		INSERT INTO platform.system_health_metrics (id, component_name, status, metrics)
		VALUES ($1, $2, $3, $4::jsonb)
		RETURNING checked_at
	`

	err := dbExec.QueryRow(ctx, stmt,
		metric.ID, metric.ComponentName, metric.Status, string(metricsJSON),
	).Scan(&metric.CheckedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to save system health metric: %w", err)
	}

	return metric, nil
}

func (r *LaunchGateRepository) ListSystemHealthMetrics(ctx context.Context) ([]domain.SystemHealthMetric, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT DISTINCT ON (component_name) id, component_name, status, metrics, checked_at
		FROM platform.system_health_metrics
		ORDER BY component_name ASC, checked_at DESC
	`

	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to list system health metrics: %w", err)
	}
	defer rows.Close()

	var list []domain.SystemHealthMetric
	for rows.Next() {
		var m domain.SystemHealthMetric
		err := rows.Scan(&m.ID, &m.ComponentName, &m.Status, &m.Metrics, &m.CheckedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan system health metric row: %w", err)
		}
		list = append(list, m)
	}

	return list, nil
}
