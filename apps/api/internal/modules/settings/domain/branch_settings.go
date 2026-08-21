package domain

import (
	"encoding/json"
)

type BranchSettings struct {
	TenantID             string          `json:"tenantId"             db:"tenant_id"`
	GeneralConfig        json.RawMessage `json:"generalConfig"        db:"general_config"`
	FinancialConfig      json.RawMessage `json:"financialConfig"      db:"financial_config"`
	InventoryConfig      json.RawMessage `json:"inventoryConfig"      db:"inventory_config"`
	IntegrationsConfig   json.RawMessage `json:"integrationsConfig"   db:"integrations_config"`
	NotificationsConfig  json.RawMessage `json:"notificationsConfig"  db:"notifications_config"`
	DocumentHeaderConfig json.RawMessage `json:"documentHeaderConfig" db:"document_header_config"`
	PatientConfig        json.RawMessage `json:"patientConfig"        db:"patient_config"`
	LimsConfig           json.RawMessage `json:"limsConfig"           db:"lims_config"`
	ConsultationConfig   json.RawMessage `json:"consultationConfig"   db:"consultation_config"`
	StaffConfig          json.RawMessage `json:"staffConfig"          db:"staff_config"`
}
