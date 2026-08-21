package domain

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type LaunchGateAudit struct {
	ID           uuid.UUID       `json:"id"`
	GateName     string          `json:"gateName"`
	Status       string          `json:"status"` // 'PASSED', 'FAILED'
	CheckResults json.RawMessage `json:"checkResults"`
	EvaluatedAt  time.Time       `json:"evaluatedAt"`
	EvaluatedBy  *uuid.UUID      `json:"evaluatedBy,omitempty"`
}

type LaunchGateCheckResult struct {
	CheckName string `json:"checkName"`
	Status    string `json:"status"` // 'PASSED', 'FAILED'
	Details   string `json:"details"`
}

type SystemHealthMetric struct {
	ID            uuid.UUID       `json:"id"`
	ComponentName string          `json:"componentName"`
	Status        string          `json:"status"` // 'HEALTHY', 'DEGRADED', 'UNHEALTHY'
	Metrics       json.RawMessage `json:"metrics"`
	CheckedAt     time.Time       `json:"checkedAt"`
}

type LaunchGateRepository interface {
	SaveLaunchGateAudit(ctx context.Context, audit *LaunchGateAudit, actorID *uuid.UUID) (*LaunchGateAudit, error)
	GetLatestLaunchGateAudit(ctx context.Context, gateName string) (*LaunchGateAudit, error)
	SaveSystemHealthMetric(ctx context.Context, metric *SystemHealthMetric) (*SystemHealthMetric, error)
	ListSystemHealthMetrics(ctx context.Context) ([]SystemHealthMetric, error)
}
