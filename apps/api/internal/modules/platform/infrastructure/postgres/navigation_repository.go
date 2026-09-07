package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/golangnigeria/curexal/internal/kernel/server"
	"github.com/golangnigeria/curexal/internal/modules/platform/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type NavigationRepository struct {
	server *server.Server
	pool   *pgxpool.Pool
}

func NewNavigationRepository(s *server.Server) *NavigationRepository {
	var pool *pgxpool.Pool
	if s != nil && s.DB != nil {
		pool = s.DB.Pool
	}
	return &NavigationRepository{
		server: s,
		pool:   pool,
	}
}

func (r *NavigationRepository) SetPool(pool *pgxpool.Pool) {
	r.pool = pool
}

func (r *NavigationRepository) GetNavigationItemsByScope(
	ctx context.Context,
	scope string,
	enabledModules []string,
	userPermissions []string,
	isSuperAdminOrOwner bool,
) ([]domain.NavigationItem, error) {
	if r.pool == nil {
		return nil, errors.New("database pool is not initialized")
	}

	query := `
		SELECT 
			id,
			COALESCE(key, id) AS key,
			context_scope,
			module_code,
			title,
			description,
			icon,
			path,
			sort_order,
			parent_id,
			required_permission,
			required_capability,
			COALESCE(status, 'active') AS status,
			COALESCE(is_visible, true) AS is_visible,
			COALESCE(is_active, true) AS is_active,
			badge_key,
			created_at
		FROM navigation_item
		WHERE context_scope = $1
		  AND is_active = true
		  AND is_visible = true
		  AND (module_code IS NULL OR module_code = ANY($2))
		  AND (required_permission IS NULL OR required_permission = ANY($3) OR $4 = TRUE)
		ORDER BY sort_order ASC
	`

	rows, err := r.pool.Query(ctx, query, scope, enabledModules, userPermissions, isSuperAdminOrOwner)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.NavigationItem
	itemMap := make(map[string]*domain.NavigationItem)

	for rows.Next() {
		var item domain.NavigationItem
		var key sql.NullString
		var moduleCode sql.NullString
		var description sql.NullString
		var parentID sql.NullString
		var requiredPermission sql.NullString
		var requiredCapability sql.NullString
		var statusStr string
		var badgeKey sql.NullString

		err := rows.Scan(
			&item.ID,
			&key,
			&item.ContextScope,
			&moduleCode,
			&item.Title,
			&description,
			&item.Icon,
			&item.Path,
			&item.Order,
			&parentID,
			&requiredPermission,
			&requiredCapability,
			&statusStr,
			&item.IsVisible,
			&item.IsActive,
			&badgeKey,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if key.Valid {
			item.Key = key.String
		} else {
			item.Key = item.ID
		}
		if moduleCode.Valid {
			item.ModuleCode = &moduleCode.String
		}
		if description.Valid {
			item.Description = &description.String
		}
		if parentID.Valid {
			item.ParentID = &parentID.String
		}
		if requiredPermission.Valid {
			item.RequiredPermission = &requiredPermission.String
		}
		if requiredCapability.Valid {
			item.RequiredCapability = &requiredCapability.String
		}
		if badgeKey.Valid {
			item.BadgeKey = &badgeKey.String
		}
		item.Status = domain.NavigationStatus(statusStr)
		item.Children = []domain.NavigationItem{}

		itemCopy := item
		items = append(items, itemCopy)
	}

	// Build hierarchy if parents exist
	var rootItems []domain.NavigationItem
	for i := range items {
		itemMap[items[i].ID] = &items[i]
	}

	for _, item := range items {
		if item.ParentID != nil && *item.ParentID != "" {
			if parent, ok := itemMap[*item.ParentID]; ok {
				parent.Children = append(parent.Children, item)
			} else {
				rootItems = append(rootItems, item)
			}
		} else {
			rootItems = append(rootItems, item)
		}
	}

	return rootItems, nil
}

func (r *NavigationRepository) GetAllActiveNavigationItems(ctx context.Context) ([]domain.NavigationItem, error) {
	if r.pool == nil {
		return nil, errors.New("database pool is not initialized")
	}

	query := `
		SELECT 
			id,
			COALESCE(key, id) AS key,
			context_scope,
			module_code,
			title,
			description,
			icon,
			path,
			sort_order,
			parent_id,
			required_permission,
			required_capability,
			COALESCE(status, 'active') AS status,
			COALESCE(is_visible, true) AS is_visible,
			COALESCE(is_active, true) AS is_active,
			badge_key,
			created_at
		FROM navigation_item
		WHERE is_active = true
		ORDER BY context_scope ASC, sort_order ASC
	`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.NavigationItem
	for rows.Next() {
		var item domain.NavigationItem
		var key sql.NullString
		var moduleCode sql.NullString
		var description sql.NullString
		var parentID sql.NullString
		var requiredPermission sql.NullString
		var requiredCapability sql.NullString
		var statusStr string
		var badgeKey sql.NullString

		err := rows.Scan(
			&item.ID,
			&key,
			&item.ContextScope,
			&moduleCode,
			&item.Title,
			&description,
			&item.Icon,
			&item.Path,
			&item.Order,
			&parentID,
			&requiredPermission,
			&requiredCapability,
			&statusStr,
			&item.IsVisible,
			&item.IsActive,
			&badgeKey,
			&item.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if key.Valid {
			item.Key = key.String
		}
		if moduleCode.Valid {
			item.ModuleCode = &moduleCode.String
		}
		if description.Valid {
			item.Description = &description.String
		}
		if parentID.Valid {
			item.ParentID = &parentID.String
		}
		if requiredPermission.Valid {
			item.RequiredPermission = &requiredPermission.String
		}
		if requiredCapability.Valid {
			item.RequiredCapability = &requiredCapability.String
		}
		if badgeKey.Valid {
			item.BadgeKey = &badgeKey.String
		}
		item.Status = domain.NavigationStatus(statusStr)

		items = append(items, item)
	}

	return items, nil
}
