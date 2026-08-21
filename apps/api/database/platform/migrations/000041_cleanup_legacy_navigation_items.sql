-- +goose Up
-- ==============================================================================
-- MIGRATION 000041: Clean up legacy duplicate platform navigation items
-- ==============================================================================

DELETE FROM navigation_item 
WHERE context_scope = 'platform' 
  AND id IN ('nav_plt_dashboard', 'nav_plt_organizations', 'nav_plt_users', 'nav_plt_settings');

-- +goose Down
-- ==============================================================================
-- No-op: legacy duplicate items should not be restored
