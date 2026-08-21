package postgres

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/golangnigeria/curexal/internal/modules/audit/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/shared/logger"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

func parseUserAgent(ua string) (os, browser, device string) {
	os = "Unknown"
	browser = "Unknown"
	device = "Desktop"

	uaLower := strings.ToLower(ua)

	if strings.Contains(uaLower, "windows") {
		os = "Windows"
	} else if strings.Contains(uaLower, "macintosh") || strings.Contains(uaLower, "mac os x") {
		os = "macOS"
	} else if strings.Contains(uaLower, "linux") {
		os = "Linux"
	} else if strings.Contains(uaLower, "android") {
		os = "Android"
	} else if strings.Contains(uaLower, "iphone") || strings.Contains(uaLower, "ipad") {
		os = "iOS"
	}

	if strings.Contains(uaLower, "firefox") {
		browser = "Firefox"
	} else if strings.Contains(uaLower, "edg/") || strings.Contains(uaLower, "edge") {
		browser = "Edge"
	} else if strings.Contains(uaLower, "chrome") {
		browser = "Chrome"
	} else if strings.Contains(uaLower, "safari") {
		browser = "Safari"
	} else if strings.Contains(uaLower, "postman") {
		browser = "Postman"
	} else if strings.Contains(uaLower, "curl") {
		browser = "curl"
	}

	if strings.Contains(uaLower, "mobi") {
		device = "Mobile"
	} else if strings.Contains(uaLower, "tablet") || strings.Contains(uaLower, "ipad") {
		device = "Tablet"
	} else if strings.Contains(uaLower, "curl") || strings.Contains(uaLower, "postman") || strings.Contains(uaLower, "go-http") {
		device = "API Client"
	}

	return os, browser, device
}

func generateAuditSignature(tenantID, actorID *string, action, severity, status string, resourceID *string) string {
	actID := ""
	if actorID != nil {
		actID = *actorID
	}
	tID := ""
	if tenantID != nil {
		tID = *tenantID
	}
	resID := ""
	if resourceID != nil {
		resID = *resourceID
	}

	raw := fmt.Sprintf("%s|%s|%s|%s|%s|%s", tID, actID, action, resID, severity, status)
	hash := sha256.Sum256([]byte(raw))
	return "sha256:" + hex.EncodeToString(hash[:])
}

type AuditRepository struct {
	server *server.Server
}

func NewAuditRepository(server *server.Server) *AuditRepository {
	return &AuditRepository{server: server}
}

func (r *AuditRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.AuditLog, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT 
			al.id, 
			al.occurred_at, 
			al.organization_id, 
			al.tenant_id, 
			al.actor_id, 
			COALESCE(u.name, al.actor_name) AS actor_name,
			COALESCE(al.actor_role, 'User') AS actor_role,
			al.action, 
			al.resource_type, 
			al.resource_id, 
			al.resource_name, 
			al.event_category, 
			al.severity, 
			al.status, 
			al.ip_address, 
			al.device, 
			al.operating_system, 
			al.browser, 
			al.user_agent, 
			al.hostname, 
			al.request_id, 
			al.session_id, 
			al.trace_id, 
			al.before_state, 
			al.after_state, 
			al.reason, 
			al.approval_reference, 
			al.digital_signature
		FROM audit.audit_events al
		LEFT JOIN identity.users u ON u.id = al.actor_id
		WHERE al.id = @id
	`

	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{"id": id})
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log: %w", err)
	}

	log, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.AuditLog])
	if err != nil {
		return nil, fmt.Errorf("audit log not found for id=%s: %w", id, err)
	}

	return &log, nil
}

func (r *AuditRepository) Create(
	ctx context.Context,
	payload *domain.CreateAuditLogPayload,
) (*domain.AuditLog, error) {
	dbExec := r.server.DB.Conn(ctx)

	tenantID := payload.TenantID
	actorID := payload.ActorID
	actorName := payload.ActorName
	actorRole := payload.ActorRole
	action := payload.Action
	resourceType := payload.ResourceType
	resourceID := payload.ResourceID
	resourceName := payload.ResourceName
	eventCategory := payload.EventCategory
	severity := payload.Severity
	status := payload.Status
	ipAddress := payload.IPAddress
	device := payload.Device
	operatingSystem := payload.OperatingSystem
	browser := payload.Browser
	userAgent := payload.UserAgent
	hostname := payload.Hostname
	requestID := payload.RequestID
	sessionID := payload.SessionID
	traceID := payload.TraceID
	beforeState := payload.BeforeState
	afterState := payload.AfterState
	reason := payload.Reason
	approvalReference := payload.ApprovalReference
	digitalSignature := payload.DigitalSignature

	sig := generateAuditSignature(tenantID, actorID, action, severity, status, resourceID)
	if digitalSignature == nil || *digitalSignature == "" {
		digitalSignature = &sig
	}

	var osVal, browserVal, deviceVal *string
	if userAgent != nil && *userAgent != "" {
		osStr, brStr, devStr := parseUserAgent(*userAgent)
		if operatingSystem == nil || *operatingSystem == "" {
			osVal = &osStr
		} else {
			osVal = operatingSystem
		}
		if browser == nil || *browser == "" {
			browserVal = &brStr
		} else {
			browserVal = browser
		}
		if device == nil || *device == "" {
			deviceVal = &devStr
		} else {
			deviceVal = device
		}
	} else {
		osVal = operatingSystem
		browserVal = browser
		deviceVal = device
	}

	var beforeJSON *json.RawMessage
	if beforeState != nil && *beforeState != "" {
		raw := json.RawMessage(*beforeState)
		beforeJSON = &raw
	}

	var afterJSON *json.RawMessage
	if afterState != nil && *afterState != "" {
		raw := json.RawMessage(*afterState)
		afterJSON = &raw
	}

	var resolvedOrgID *uuid.UUID
	var resolvedTenantID *uuid.UUID

	if tenantID != nil && *tenantID != "" {
		if parsedTID, err := uuid.Parse(*tenantID); err == nil {
			resolvedTenantID = &parsedTID
			var orgIDStr string
			errOrg := dbExec.QueryRow(ctx, `SELECT organization_id FROM organization.tenants WHERE id = $1`, parsedTID).Scan(&orgIDStr)
			if errOrg == nil && orgIDStr != "" {
				if parsedOID, err2 := uuid.Parse(orgIDStr); err2 == nil {
					resolvedOrgID = &parsedOID
				}
			}
		}
	}

	stmt := `
		INSERT INTO audit.audit_events (
			organization_id, tenant_id, actor_id, actor_name, actor_role,
			action, resource_type, resource_id, resource_name, event_category,
			severity, status, ip_address, device, operating_system,
			browser, user_agent, hostname, request_id, session_id,
			trace_id, before_state, after_state, reason, approval_reference, digital_signature
		) VALUES (
			@organization_id, @tenant_id, @actor_id, @actor_name, @actor_role,
			@action, @resource_type, @resource_id, @resource_name, @event_category,
			@severity, @status, @ip_address, @device, @operating_system,
			@browser, @user_agent, @hostname, @request_id, @session_id,
			@trace_id, @before_state, @after_state, @reason, @approval_reference, @digital_signature
		)
		RETURNING 
			id, occurred_at, organization_id, tenant_id, actor_id, actor_name, actor_role,
			action, resource_type, resource_id, resource_name, event_category, severity, status,
			ip_address, device, operating_system, browser, user_agent, hostname, request_id,
			session_id, trace_id, before_state, after_state, reason, approval_reference, digital_signature
	`

	args := pgx.NamedArgs{
		"organization_id":    resolvedOrgID,
		"tenant_id":          resolvedTenantID,
		"actor_id":           actorID,
		"actor_name":         actorName,
		"actor_role":         actorRole,
		"action":             action,
		"resource_type":      resourceType,
		"resource_id":        resourceID,
		"resource_name":      resourceName,
		"event_category":     eventCategory,
		"severity":           severity,
		"status":             status,
		"ip_address":         ipAddress,
		"device":             deviceVal,
		"operating_system":   osVal,
		"browser":            browserVal,
		"user_agent":         userAgent,
		"hostname":           hostname,
		"request_id":         requestID,
		"session_id":         sessionID,
		"trace_id":           traceID,
		"before_state":       beforeJSON,
		"after_state":        afterJSON,
		"reason":             reason,
		"approval_reference": approvalReference,
		"digital_signature":  digitalSignature,
	}

	rows, err := dbExec.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to create audit log entry: %w", err)
	}

	createdLog, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.AuditLog])
	if err != nil {
		return nil, fmt.Errorf("failed to collect inserted audit log: %w", err)
	}

	logger.LogBusinessAudit(ctx, action, fmt.Sprintf("Audit log entry registered with ID: %s", createdLog.ID.String()))

	return &createdLog, nil
}

func (r *AuditRepository) ListTenantLogs(ctx context.Context, tenantID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]domain.AuditLog, error) {
	dbExec := r.server.DB.Conn(ctx)

	if limit <= 0 {
		limit = 50
	}

	whereClauses := []string{"1=1"}
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	if tenantID != nil {
		whereClauses = append(whereClauses, "al.tenant_id = @tenant_id")
		args["tenant_id"] = *tenantID
	}
	if category != nil && *category != "" {
		whereClauses = append(whereClauses, "al.event_category = @category")
		args["category"] = *category
	}
	if severity != nil && *severity != "" {
		whereClauses = append(whereClauses, "al.severity = @severity")
		args["severity"] = *severity
	}
	if status != nil && *status != "" {
		whereClauses = append(whereClauses, "al.status = @status")
		args["status"] = *status
	}
	if actorID != nil && *actorID != "" {
		whereClauses = append(whereClauses, "al.actor_id = @actor_id")
		args["actor_id"] = *actorID
	}
	if action != nil && *action != "" {
		whereClauses = append(whereClauses, "al.action ILIKE @action")
		args["action"] = "%" + *action + "%"
	}
	if resourceType != nil && *resourceType != "" {
		whereClauses = append(whereClauses, "al.resource_type = @resource_type")
		args["resource_type"] = *resourceType
	}
	if resourceID != nil && *resourceID != "" {
		whereClauses = append(whereClauses, "al.resource_id = @resource_id")
		args["resource_id"] = *resourceID
	}
	if startDate != nil && *startDate != "" {
		whereClauses = append(whereClauses, "al.occurred_at >= @start_date::timestamptz")
		args["start_date"] = *startDate
	}
	if endDate != nil && *endDate != "" {
		whereClauses = append(whereClauses, "al.occurred_at <= @end_date::timestamptz")
		args["end_date"] = *endDate
	}
	if search != nil && *search != "" {
		whereClauses = append(whereClauses, "(al.action ILIKE @search OR al.resource_name ILIKE @search OR al.reason ILIKE @search OR u.name ILIKE @search)")
		args["search"] = "%" + *search + "%"
	}

	stmt := fmt.Sprintf(`
		SELECT 
			al.id, 
			al.occurred_at, 
			al.organization_id, 
			al.tenant_id, 
			al.actor_id, 
			COALESCE(u.name, al.actor_name) AS actor_name,
			COALESCE(al.actor_role, 'User') AS actor_role,
			al.action, 
			al.resource_type, 
			al.resource_id, 
			al.resource_name, 
			al.event_category, 
			al.severity, 
			al.status, 
			al.ip_address, 
			al.device, 
			al.operating_system, 
			al.browser, 
			al.user_agent, 
			al.hostname, 
			al.request_id, 
			al.session_id, 
			al.trace_id, 
			al.before_state, 
			al.after_state, 
			al.reason, 
			al.approval_reference, 
			al.digital_signature
		FROM audit.audit_events al
		LEFT JOIN identity.users u ON u.id = al.actor_id
		WHERE %s
		ORDER BY al.occurred_at DESC
		LIMIT @limit OFFSET @offset
	`, strings.Join(whereClauses, " AND "))

	rows, err := dbExec.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to query tenant audit logs: %w", err)
	}
	defer rows.Close()

	logs, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.AuditLog])
	if err != nil {
		return nil, fmt.Errorf("failed to collect tenant audit logs: %w", err)
	}

	return logs, nil
}

func (r *AuditRepository) ListPlatformLogs(ctx context.Context, orgID *uuid.UUID, category, severity, status, actorID, action, resourceType, resourceID, startDate, endDate, search *string, limit, offset int) ([]domain.AuditLog, error) {
	dbExec := r.server.DB.Conn(ctx)

	if limit <= 0 {
		limit = 50
	}

	whereClauses := []string{"1=1"}
	args := pgx.NamedArgs{
		"limit":  limit,
		"offset": offset,
	}

	if orgID != nil {
		whereClauses = append(whereClauses, "organization_id = @organization_id")
		args["organization_id"] = *orgID
	}
	if category != nil && *category != "" {
		whereClauses = append(whereClauses, "event_category = @category")
		args["category"] = *category
	}
	if severity != nil && *severity != "" {
		whereClauses = append(whereClauses, "severity = @severity")
		args["severity"] = *severity
	}
	if status != nil && *status != "" {
		whereClauses = append(whereClauses, "status = @status")
		args["status"] = *status
	}
	if actorID != nil && *actorID != "" {
		whereClauses = append(whereClauses, "actor_id = @actor_id")
		args["actor_id"] = *actorID
	}
	if action != nil && *action != "" {
		whereClauses = append(whereClauses, "action ILIKE @action")
		args["action"] = "%" + *action + "%"
	}
	if resourceType != nil && *resourceType != "" {
		whereClauses = append(whereClauses, "resource_type = @resource_type")
		args["resource_type"] = *resourceType
	}
	if resourceID != nil && *resourceID != "" {
		whereClauses = append(whereClauses, "resource_id = @resource_id")
		args["resource_id"] = *resourceID
	}
	if startDate != nil && *startDate != "" {
		whereClauses = append(whereClauses, "occurred_at >= @start_date::timestamptz")
		args["start_date"] = *startDate
	}
	if endDate != nil && *endDate != "" {
		whereClauses = append(whereClauses, "occurred_at <= @end_date::timestamptz")
		args["end_date"] = *endDate
	}
	if search != nil && *search != "" {
		whereClauses = append(whereClauses, "(action ILIKE @search OR resource_name ILIKE @search OR reason ILIKE @search OR actor_name ILIKE @search)")
		args["search"] = "%" + *search + "%"
	}

	stmt := fmt.Sprintf(`
		SELECT 
			id, occurred_at, organization_id, tenant_id, actor_id, actor_name, actor_role,
			action, resource_type, resource_id, resource_name, event_category, severity, status,
			ip_address, device, operating_system, browser, user_agent, hostname, request_id,
			session_id, trace_id, before_state, after_state, reason, approval_reference, digital_signature
		FROM audit.audit_events
		WHERE %s
		ORDER BY occurred_at DESC
		LIMIT @limit OFFSET @offset
	`, strings.Join(whereClauses, " AND "))

	rows, err := dbExec.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to query platform audit logs: %w", err)
	}
	defer rows.Close()

	logs, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.AuditLog])
	if err != nil {
		return nil, fmt.Errorf("failed to collect platform audit logs: %w", err)
	}

	return logs, nil
}

func (r *AuditRepository) GetStats(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID) (*domain.AdminStats, error) {
	dbExec := r.server.DB.Conn(ctx)
	stats := &domain.AdminStats{
		SeverityDistribution: make(map[string]int),
		ActionDistribution:   make(map[string]int),
		OrganizationStats:    []domain.AdminOrgStats{},
		ActivityTrend:        []domain.AdminActivityTrend{},
	}

	var err error

	err = dbExec.QueryRow(ctx, `SELECT COUNT(*) FROM organization.organizations`).Scan(&stats.TotalOrganizations)
	if err != nil {
		stats.TotalOrganizations = 0
	}

	err = dbExec.QueryRow(ctx, `SELECT COUNT(*) FROM organization.tenants`).Scan(&stats.TotalTenants)
	if err != nil {
		stats.TotalTenants = 0
	}

	err = dbExec.QueryRow(ctx, `SELECT COUNT(*) FROM organization.memberships WHERE is_active = TRUE`).Scan(&stats.TotalMemberships)
	if err != nil {
		stats.TotalMemberships = 0
	}

	var countLogs int
	_ = dbExec.QueryRow(ctx, `SELECT COUNT(*) FROM audit.audit_events`).Scan(&countLogs)
	stats.TotalAuditLogs = countLogs

	rowsSev, errSev := dbExec.Query(ctx, `
		SELECT severity, COUNT(*) 
		FROM audit.audit_events
		GROUP BY severity
	`)
	if errSev == nil {
		defer rowsSev.Close()
		for rowsSev.Next() {
			var sev string
			var cnt int
			if err := rowsSev.Scan(&sev, &cnt); err == nil {
				stats.SeverityDistribution[sev] = cnt
			}
		}
	}

	rowsAct, errAct := dbExec.Query(ctx, `
		SELECT action, COUNT(*) 
		FROM audit.audit_events
		GROUP BY action
		ORDER BY COUNT(*) DESC
		LIMIT 10
	`)
	if errAct == nil {
		defer rowsAct.Close()
		for rowsAct.Next() {
			var act string
			var cnt int
			if err := rowsAct.Scan(&act, &cnt); err == nil {
				stats.ActionDistribution[act] = cnt
			}
		}
	}

	rowsOrg, errOrg := dbExec.Query(ctx, `
		SELECT 
			o.id::text, 
			o.name, 
			(SELECT COUNT(*) FROM organization.tenants t WHERE t.organization_id = o.id) AS tenant_count,
			(SELECT COUNT(DISTINCT m.user_id) FROM organization.memberships m WHERE m.organization_id = o.id) AS member_count,
			(SELECT COUNT(*) FROM audit.audit_events al WHERE al.organization_id = o.id) AS audit_count
		FROM organization.organizations o
		ORDER BY o.created_at DESC
		LIMIT 20
	`)
	if errOrg == nil {
		defer rowsOrg.Close()
		for rowsOrg.Next() {
			var oStats domain.AdminOrgStats
			if err := rowsOrg.Scan(&oStats.OrgID, &oStats.OrgName, &oStats.TenantCount, &oStats.MemberCount, &oStats.AuditLogCount); err == nil {
				stats.OrganizationStats = append(stats.OrganizationStats, oStats)
			}
		}
	}

	return stats, nil
}

func (r *AuditRepository) ListAll(ctx context.Context, tenantID *uuid.UUID, orgID *uuid.UUID, limit int, offset int) ([]domain.AuditLog, error) {
	if tenantID != nil {
		return r.ListTenantLogs(ctx, tenantID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, limit, offset)
	}
	return r.ListPlatformLogs(ctx, orgID, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, limit, offset)
}

