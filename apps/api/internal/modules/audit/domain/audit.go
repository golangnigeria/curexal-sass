package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditLog struct {
	ID                uuid.UUID        `json:"id"                db:"id"`
	OccurredAt        time.Time        `json:"occurredAt"        db:"occurred_at"`
	OrganizationID    *uuid.UUID       `json:"organizationId"    db:"organization_id"`
	TenantID          *uuid.UUID       `json:"tenantId"          db:"tenant_id"`
	ActorID           *string          `json:"actorId"           db:"actor_id"`
	ActorName         *string          `json:"actorName"         db:"actor_name"`
	ActorRole         *string          `json:"actorRole"         db:"actor_role"`
	Action            string           `json:"action"            db:"action"`
	ResourceType      *string          `json:"resourceType"      db:"resource_type"`
	ResourceID        *string          `json:"resourceId"        db:"resource_id"`
	ResourceName      *string          `json:"resourceName"      db:"resource_name"`
	EventCategory     *string          `json:"eventCategory"     db:"event_category"`
	Severity          string           `json:"severity"          db:"severity"`
	Status            string           `json:"status"            db:"status"`
	IPAddress         *string          `json:"ipAddress"         db:"ip_address"`
	Device            *string          `json:"device"            db:"device"`
	OperatingSystem   *string          `json:"operatingSystem"   db:"operating_system"`
	Browser           *string          `json:"browser"           db:"browser"`
	UserAgent         *string          `json:"userAgent"         db:"user_agent"`
	Hostname          *string          `json:"hostname"          db:"hostname"`
	RequestID         *string          `json:"requestId"         db:"request_id"`
	SessionID         *string          `json:"sessionId"         db:"session_id"`
	TraceID           *string          `json:"traceId"           db:"trace_id"`
	BeforeState       *json.RawMessage `json:"beforeState"       db:"before_state"`
	AfterState        *json.RawMessage `json:"afterState"        db:"after_state"`
	Reason            *string          `json:"reason"            db:"reason"`
	ApprovalReference *string          `json:"approvalReference" db:"approval_reference"`
	DigitalSignature  *string          `json:"digitalSignature"  db:"digital_signature"`
}

type AdminOrgStats struct {
	OrgID         string `json:"orgId"`
	OrgName       string `json:"orgName"`
	TenantCount   int    `json:"tenantCount"`
	MemberCount   int    `json:"memberCount"`
	AuditLogCount int    `json:"auditLogCount"`
}

type AdminActivityTrend struct {
	Date  string `json:"date"`
	Count int    `json:"count"`
}

type AdminStats struct {
	TotalOrganizations   int                  `json:"totalOrganizations"`
	TotalTenants         int                  `json:"totalTenants"`
	TotalMemberships     int                  `json:"totalMemberships"`
	TotalAuditLogs       int                  `json:"totalAuditLogs"`
	SeverityDistribution map[string]int       `json:"severityDistribution"`
	ActionDistribution   map[string]int       `json:"actionDistribution"`
	OrganizationStats    []AdminOrgStats      `json:"organizationStats"`
	ActivityTrend        []AdminActivityTrend `json:"activityTrend"`
}

type CreateAuditLogPayload struct {
	IsPlatform        bool    `json:"isPlatform"`
	TenantID          *string `json:"tenantId"`
	ActorID           *string `json:"actorId"`
	ActorName         *string `json:"actorName"`
	ActorRole         *string `json:"actorRole"`
	Action            string  `json:"action"`
	ResourceType      *string `json:"resourceType"`
	ResourceID        *string `json:"resourceId"`
	ResourceName      *string `json:"resourceName"`
	EventCategory     *string `json:"eventCategory"`
	Severity          string  `json:"severity"`
	Status            string  `json:"status"`
	IPAddress         *string `json:"ipAddress"`
	Device            *string `json:"device"`
	OperatingSystem   *string `json:"operatingSystem"`
	Browser           *string `json:"browser"`
	UserAgent         *string `json:"userAgent"`
	Hostname          *string `json:"hostname"`
	RequestID         *string `json:"requestId"`
	SessionID         *string `json:"sessionId"`
	TraceID           *string `json:"traceId"`
	BeforeState       *string `json:"beforeState"`
	AfterState        *string `json:"afterState"`
	Reason            *string `json:"reason"`
	ApprovalReference *string `json:"approvalReference"`
	DigitalSignature  *string `json:"digitalSignature"`
}
