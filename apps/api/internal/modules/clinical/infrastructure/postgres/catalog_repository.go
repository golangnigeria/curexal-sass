package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/clinical/domain"
	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/google/uuid"
)

type CatalogRepository struct {
	server *server.Server
}

func NewCatalogRepository(server *server.Server) *CatalogRepository {
	return &CatalogRepository{server: server}
}

func (r *CatalogRepository) ListCatalog(ctx context.Context) ([]domain.CatalogItem, error) {
	dbExec := r.server.DB.Conn(ctx)

	stmt := `
		SELECT id, name, code, base_price, type, urgency_price, commission_value, commission_percentage, 
		       discount_amount, discount_percentage, display_name, short_name, recovery_time, department_id, 
		       test_group, test_category, tat_hours, created_at, updated_at
		FROM services
		ORDER BY created_at DESC
	`

	rows, err := dbExec.Query(ctx, stmt)
	if err != nil {
		return nil, fmt.Errorf("failed to list services: %w", err)
	}
	defer rows.Close()

	items := []domain.CatalogItem{}
	for rows.Next() {
		var item domain.CatalogItem
		err := rows.Scan(
			&item.ID,
			&item.Name,
			&item.Code,
			&item.BasePrice,
			&item.Type,
			&item.UrgencyPrice,
			&item.CommissionValue,
			&item.CommissionPercentage,
			&item.DiscountAmount,
			&item.DiscountPercentage,
			&item.DisplayName,
			&item.ShortName,
			&item.RecoveryTime,
			&item.DepartmentID,
			&item.TestGroup,
			&item.TestCategory,
			&item.TatHours,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan service row: %w", err)
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading services rows: %w", err)
	}

	rowsChildren, err := dbExec.Query(ctx, "SELECT profile_id, service_id FROM profile_services")
	if err == nil {
		defer rowsChildren.Close()
		childrenMap := make(map[uuid.UUID][]uuid.UUID)
		for rowsChildren.Next() {
			var pid, sid uuid.UUID
			if errScan := rowsChildren.Scan(&pid, &sid); errScan == nil {
				childrenMap[pid] = append(childrenMap[pid], sid)
			}
		}

		for i := range items {
			if items[i].Type == "profile" {
				if kids, ok := childrenMap[items[i].ID]; ok {
					items[i].ChildServiceIDs = kids
				} else {
					items[i].ChildServiceIDs = []uuid.UUID{}
				}
			}
		}
	}

	return items, nil
}

func (r *CatalogRepository) GetCatalogItemByID(ctx context.Context, id uuid.UUID) (*domain.CatalogItem, error) {
	dbExec := r.server.DB.Conn(ctx)

	stmt := `
		SELECT id, name, code, base_price, type, urgency_price, commission_value, commission_percentage, 
		       discount_amount, discount_percentage, display_name, short_name, recovery_time, department_id, 
		       test_group, test_category, tat_hours, created_at, updated_at
		FROM services
		WHERE id = $1
	`

	var item domain.CatalogItem
	err := dbExec.QueryRow(ctx, stmt, id).Scan(
		&item.ID,
		&item.Name,
		&item.Code,
		&item.BasePrice,
		&item.Type,
		&item.UrgencyPrice,
		&item.CommissionValue,
		&item.CommissionPercentage,
		&item.DiscountAmount,
		&item.DiscountPercentage,
		&item.DisplayName,
		&item.ShortName,
		&item.RecoveryTime,
		&item.DepartmentID,
		&item.TestGroup,
		&item.TestCategory,
		&item.TatHours,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if item.Type == "profile" {
		item.ChildServiceIDs = []uuid.UUID{}
		rows, err := dbExec.Query(ctx, "SELECT service_id FROM profile_services WHERE profile_id = $1", id)
		if err == nil {
			defer rows.Close()
			for rows.Next() {
				var sid uuid.UUID
				if errScan := rows.Scan(&sid); errScan == nil {
					item.ChildServiceIDs = append(item.ChildServiceIDs, sid)
				}
			}
		}
	}

	return &item, nil
}

func (r *CatalogRepository) CreateCatalogItem(ctx context.Context, payload *domain.CreateCatalogItemPayload) (*domain.CatalogItem, error) {
	db := r.server.DB
	var createdItem *domain.CatalogItem

	err := db.Tx(ctx, func(txCtx context.Context) error {
		tx := db.Conn(txCtx)

		newID := uuid.New()
		now := time.Now()

		stmt := `
			INSERT INTO services (id, name, code, base_price, type, urgency_price, commission_value, commission_percentage, 
			       discount_amount, discount_percentage, display_name, short_name, recovery_time, department_id, 
			       test_group, test_category, tat_hours, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
		`

		code := payload.Code
		if len(code) > 64 {
			code = code[:64]
		}

		_, err := tx.Exec(txCtx, stmt,
			newID,
			payload.Name,
			code,
			payload.BasePrice,
			payload.Type,
			payload.UrgencyPrice,
			payload.CommissionValue,
			payload.CommissionPercentage,
			payload.DiscountAmount,
			payload.DiscountPercentage,
			payload.DisplayName,
			payload.ShortName,
			payload.RecoveryTime,
			payload.DepartmentID,
			payload.TestGroup,
			payload.TestCategory,
			payload.TatHours,
			now,
			now,
		)
		if err != nil {
			return fmt.Errorf("failed to insert service: %w", err)
		}

		if payload.Type == "profile" && len(payload.ChildServiceIDs) > 0 {
			for _, childID := range payload.ChildServiceIDs {
				_, err = tx.Exec(txCtx, `
					INSERT INTO profile_services (profile_id, service_id)
					VALUES ($1, $2)
				`, newID, childID)
				if err != nil {
					return fmt.Errorf("failed to insert profile service child mapping: %w", err)
				}
			}
		}

		var item domain.CatalogItem
		getStmt := `
			SELECT id, name, code, base_price, type, urgency_price, commission_value, commission_percentage, 
			       discount_amount, discount_percentage, display_name, short_name, recovery_time, department_id, 
			       test_group, test_category, tat_hours, created_at, updated_at
			FROM services WHERE id = $1
		`
		err = tx.QueryRow(txCtx, getStmt, newID).Scan(
			&item.ID,
			&item.Name,
			&item.Code,
			&item.BasePrice,
			&item.Type,
			&item.UrgencyPrice,
			&item.CommissionValue,
			&item.CommissionPercentage,
			&item.DiscountAmount,
			&item.DiscountPercentage,
			&item.DisplayName,
			&item.ShortName,
			&item.RecoveryTime,
			&item.DepartmentID,
			&item.TestGroup,
			&item.TestCategory,
			&item.TatHours,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to select created service: %w", err)
		}

		if item.Type == "profile" {
			item.ChildServiceIDs = payload.ChildServiceIDs
		}
		createdItem = &item
		return nil
	})

	if err != nil {
		return nil, err
	}
	return createdItem, nil
}

func (r *CatalogRepository) UpdateCatalogMetadata(ctx context.Context, id uuid.UUID, payload *domain.CreateCatalogItemPayload) (*domain.CatalogItem, error) {
	db := r.server.DB
	var updatedItem *domain.CatalogItem

	err := db.Tx(ctx, func(txCtx context.Context) error {
		tx := db.Conn(txCtx)

		stmt := `
			UPDATE services
			SET name = $1,
			    code = $2,
			    type = $3,
			    display_name = $4,
			    short_name = $5,
			    recovery_time = $6,
			    department_id = $7,
			    test_group = $8,
			    test_category = $9,
			    tat_hours = $10,
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $11
		`

		_, err := tx.Exec(txCtx, stmt,
			payload.Name,
			payload.Code,
			payload.Type,
			payload.DisplayName,
			payload.ShortName,
			payload.RecoveryTime,
			payload.DepartmentID,
			payload.TestGroup,
			payload.TestCategory,
			payload.TatHours,
			id,
		)
		if err != nil {
			return fmt.Errorf("failed to update service metadata: %w", err)
		}

		_, err = tx.Exec(txCtx, "DELETE FROM profile_services WHERE profile_id = $1", id)
		if err != nil {
			return fmt.Errorf("failed to clear old profile service mappings: %w", err)
		}

		if payload.Type == "profile" && len(payload.ChildServiceIDs) > 0 {
			for _, childID := range payload.ChildServiceIDs {
				_, err = tx.Exec(txCtx, `
					INSERT INTO profile_services (profile_id, service_id)
					VALUES ($1, $2)
				`, id, childID)
				if err != nil {
					return fmt.Errorf("failed to insert profile service child mapping: %w", err)
				}
			}
		}

		var item domain.CatalogItem
		getStmt := `
			SELECT id, name, code, base_price, type, urgency_price, commission_value, commission_percentage, 
			       discount_amount, discount_percentage, display_name, short_name, recovery_time, department_id, 
			       test_group, test_category, tat_hours, created_at, updated_at
			FROM services WHERE id = $1
		`
		err = tx.QueryRow(txCtx, getStmt, id).Scan(
			&item.ID,
			&item.Name,
			&item.Code,
			&item.BasePrice,
			&item.Type,
			&item.UrgencyPrice,
			&item.CommissionValue,
			&item.CommissionPercentage,
			&item.DiscountAmount,
			&item.DiscountPercentage,
			&item.DisplayName,
			&item.ShortName,
			&item.RecoveryTime,
			&item.DepartmentID,
			&item.TestGroup,
			&item.TestCategory,
			&item.TatHours,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to select updated service: %w", err)
		}

		if item.Type == "profile" {
			item.ChildServiceIDs = payload.ChildServiceIDs
		}
		updatedItem = &item
		return nil
	})

	if err != nil {
		return nil, err
	}
	return updatedItem, nil
}

func (r *CatalogRepository) UpdateCatalogPricing(ctx context.Context, id uuid.UUID, payload *domain.UpdatePricingPayload) (*domain.CatalogItem, error) {
	db := r.server.DB
	var updatedItem *domain.CatalogItem

	err := db.Tx(ctx, func(txCtx context.Context) error {
		tx := db.Conn(txCtx)

		stmt := `
			UPDATE services
			SET base_price = COALESCE($1, base_price),
			    urgency_price = COALESCE($2, urgency_price),
			    commission_value = COALESCE($3, commission_value),
			    commission_percentage = COALESCE($4, commission_percentage),
			    discount_amount = COALESCE($5, discount_amount),
			    discount_percentage = COALESCE($6, discount_percentage),
			    updated_at = CURRENT_TIMESTAMP
			WHERE id = $7
		`

		_, err := tx.Exec(txCtx, stmt,
			payload.BasePrice,
			payload.UrgencyPrice,
			payload.CommissionValue,
			payload.CommissionPercentage,
			payload.DiscountAmount,
			payload.DiscountPercentage,
			id,
		)
		if err != nil {
			return fmt.Errorf("failed to update service pricing: %w", err)
		}

		var item domain.CatalogItem
		getStmt := `
			SELECT id, name, code, base_price, type, urgency_price, commission_value, commission_percentage, 
			       discount_amount, discount_percentage, display_name, short_name, recovery_time, department_id, 
			       test_group, test_category, tat_hours, created_at, updated_at
			FROM services WHERE id = $1
		`
		err = tx.QueryRow(txCtx, getStmt, id).Scan(
			&item.ID,
			&item.Name,
			&item.Code,
			&item.BasePrice,
			&item.Type,
			&item.UrgencyPrice,
			&item.CommissionValue,
			&item.CommissionPercentage,
			&item.DiscountAmount,
			&item.DiscountPercentage,
			&item.DisplayName,
			&item.ShortName,
			&item.RecoveryTime,
			&item.DepartmentID,
			&item.TestGroup,
			&item.TestCategory,
			&item.TatHours,
			&item.CreatedAt,
			&item.UpdatedAt,
		)
		if err != nil {
			return fmt.Errorf("failed to select updated service: %w", err)
		}

		if item.Type == "profile" {
			item.ChildServiceIDs = []uuid.UUID{}
			rows, err := tx.Query(txCtx, "SELECT service_id FROM profile_services WHERE profile_id = $1", id)
			if err == nil {
				defer rows.Close()
				for rows.Next() {
					var sid uuid.UUID
					if errScan := rows.Scan(&sid); errScan == nil {
						item.ChildServiceIDs = append(item.ChildServiceIDs, sid)
					}
				}
			}
		}

		updatedItem = &item
		return nil
	})

	if err != nil {
		return nil, err
	}
	return updatedItem, nil
}

func (r *CatalogRepository) ImportCatalog(ctx context.Context, items []domain.CreateCatalogItemPayload) (int, error) {
	db := r.server.DB
	importedCount := 0

	err := db.Tx(ctx, func(txCtx context.Context) error {
		tx := db.Conn(txCtx)
		now := time.Now()

		for _, item := range items {
			if len(item.Code) > 64 {
				item.Code = item.Code[:64]
			}
			if item.ShortName != nil && len(*item.ShortName) > 100 {
				truncated := (*item.ShortName)[:100]
				item.ShortName = &truncated
			}
			if item.DisplayName != nil && len(*item.DisplayName) > 255 {
				truncated := (*item.DisplayName)[:255]
				item.DisplayName = &truncated
			}

			var exists bool
			err := tx.QueryRow(txCtx, "SELECT EXISTS(SELECT 1 FROM services WHERE code = $1)", item.Code).Scan(&exists)
			if err != nil {
				return fmt.Errorf("failed duplicate code check for %s: %w", item.Code, err)
			}

			if exists {
				updateStmt := `
					UPDATE services
					SET name = $1,
					    base_price = $2,
					    type = $3,
					    urgency_price = $4,
					    commission_value = $5,
					    commission_percentage = $6,
					    discount_amount = $7,
					    discount_percentage = $8,
					    display_name = $9,
					    short_name = $10,
					    recovery_time = $11,
					    department_id = $12,
					    test_group = $13,
					    test_category = $14,
					    tat_hours = $15,
					    updated_at = CURRENT_TIMESTAMP
					WHERE code = $16
					RETURNING id
				`
				var existingID uuid.UUID
				err = tx.QueryRow(txCtx, updateStmt,
					item.Name,
					item.BasePrice,
					item.Type,
					item.UrgencyPrice,
					item.CommissionValue,
					item.CommissionPercentage,
					item.DiscountAmount,
					item.DiscountPercentage,
					item.DisplayName,
					item.ShortName,
					item.RecoveryTime,
					item.DepartmentID,
					item.TestGroup,
					item.TestCategory,
					item.TatHours,
					item.Code,
				).Scan(&existingID)
				if err != nil {
					return fmt.Errorf("failed to update duplicate service: %w", err)
				}

				if item.Type == "profile" && len(item.ChildServiceIDs) > 0 {
					_, err = tx.Exec(txCtx, "DELETE FROM profile_services WHERE profile_id = $1", existingID)
					if err != nil {
						return fmt.Errorf("failed to delete child services for duplicate check: %w", err)
					}
					for _, childID := range item.ChildServiceIDs {
						_, err = tx.Exec(txCtx, "INSERT INTO profile_services (profile_id, service_id) VALUES ($1, $2)", existingID, childID)
						if err != nil {
							return fmt.Errorf("failed to insert child service for duplicate check: %w", err)
						}
					}
				}

				importedCount++
				continue
			}

			newID := uuid.New()
			insertStmt := `
				INSERT INTO services (id, name, code, base_price, type, urgency_price, commission_value, commission_percentage, 
				       discount_amount, discount_percentage, display_name, short_name, recovery_time, department_id, 
				       test_group, test_category, tat_hours, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
			`
			_, err = tx.Exec(txCtx, insertStmt,
				newID,
				item.Name,
				item.Code,
				item.BasePrice,
				item.Type,
				item.UrgencyPrice,
				item.CommissionValue,
				item.CommissionPercentage,
				item.DiscountAmount,
				item.DiscountPercentage,
				item.DisplayName,
				item.ShortName,
				item.RecoveryTime,
				item.DepartmentID,
				item.TestGroup,
				item.TestCategory,
				item.TatHours,
				now,
				now,
			)
			if err != nil {
				return fmt.Errorf("failed to insert imported service: %w", err)
			}

			if item.Type == "profile" && len(item.ChildServiceIDs) > 0 {
				for _, childID := range item.ChildServiceIDs {
					_, err = tx.Exec(txCtx, "INSERT INTO profile_services (profile_id, service_id) VALUES ($1, $2)", newID, childID)
					if err != nil {
						return fmt.Errorf("failed to insert imported profile mapping: %w", err)
					}
				}
			}

			importedCount++
		}
		return nil
	})

	if err != nil {
		return 0, err
	}
	return importedCount, nil
}
