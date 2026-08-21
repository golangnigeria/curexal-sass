package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type DemoRepository struct {
	server *server.Server
}

func NewDemoRepository(server *server.Server) *DemoRepository {
	return &DemoRepository{server: server}
}

func (r *DemoRepository) Create(ctx context.Context, labName, contactName, email string, phone, dailyVolume, notes *string) (*domain.DemoRequest, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		INSERT INTO organization.demo_requests (laboratory_name, contact_name, email, phone, daily_volume, notes)
		VALUES (@laboratory_name, @contact_name, @email, @phone, @daily_volume, @notes)
		RETURNING id, laboratory_name, contact_name, email, phone, daily_volume, notes, status, created_at, updated_at
	`

	rows, err := dbExec.Query(ctx, stmt, pgx.NamedArgs{
		"laboratory_name": labName,
		"contact_name":    contactName,
		"email":           email,
		"phone":           phone,
		"daily_volume":    dailyVolume,
		"notes":           notes,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute create demo request query: %w", err)
	}

	dr, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.DemoRequest])
	if err != nil {
		return nil, fmt.Errorf("failed to collect row from table:organization.demo_requests: %w", err)
	}

	return &dr, nil
}

func (r *DemoRepository) List(ctx context.Context) ([]domain.DemoRequest, error) {
	dbExec := r.server.DB.Conn(ctx)
	stmt := `
		SELECT id, laboratory_name, contact_name, email, phone, daily_volume, notes, status, created_at, updated_at
		FROM organization.demo_requests
		ORDER BY created_at DESC
	`

	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to list demo requests: %w", err)
	}
	defer rows.Close()

	drs, err := pgx.CollectRows(rows, pgx.RowToStructByName[domain.DemoRequest])
	if err != nil {
		return nil, fmt.Errorf("failed to collect demo request rows: %w", err)
	}

	return drs, nil
}

func (r *DemoRepository) Update(ctx context.Context, id uuid.UUID, status, notes *string) (*domain.DemoRequest, error) {
	dbExec := r.server.DB.Conn(ctx)
	setClauses := []string{"updated_at = CURRENT_TIMESTAMP"}
	args := pgx.NamedArgs{
		"id": id,
	}

	if status != nil {
		setClauses = append(setClauses, "status = @status")
		args["status"] = *status
	}
	if notes != nil {
		setClauses = append(setClauses, "notes = @notes")
		args["notes"] = *notes
	}

	stmt := fmt.Sprintf(`
		UPDATE organization.demo_requests
		SET %s
		WHERE id = @id
		RETURNING id, laboratory_name, contact_name, email, phone, daily_volume, notes, status, created_at, updated_at
	`, strings.Join(setClauses, ", "))

	rows, err := dbExec.Query(ctx, stmt, args)
	if err != nil {
		return nil, fmt.Errorf("failed to update demo request query for id=%s: %w", id, err)
	}

	dr, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.DemoRequest])
	if err != nil {
		return nil, fmt.Errorf("failed to collect updated row: %w", err)
	}

	return &dr, nil
}
