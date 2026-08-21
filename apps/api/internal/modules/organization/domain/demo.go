package domain

import (
	"time"

	"github.com/google/uuid"
)

type DemoRequest struct {
	ID             uuid.UUID `json:"id"             db:"id"`
	LaboratoryName string    `json:"laboratoryName" db:"laboratory_name"`
	ContactName    string    `json:"contactName"    db:"contact_name"`
	Email          string    `json:"email"          db:"email"`
	Phone          *string   `json:"phone"          db:"phone"`
	DailyVolume    *string   `json:"dailyVolume"    db:"daily_volume"`
	Notes          *string   `json:"notes"          db:"notes"`
	Status         string    `json:"status"         db:"status"`
	CreatedAt      time.Time `json:"createdAt"      db:"created_at"`
	UpdatedAt      time.Time `json:"updatedAt"      db:"updated_at"`
}
