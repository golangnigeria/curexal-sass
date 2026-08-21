-- +goose Up
-- ==============================================================================
-- MIGRATION 000042: Update platform dashboard navigation title and icon
-- ==============================================================================

UPDATE navigation_item 
SET title = 'Platform Dashboard', 
    icon = 'LayoutDashboard'
WHERE id = 'nav_plat_dashboard' OR (context_scope = 'platform' AND path = '/platform/dashboard');

-- +goose Down
-- ==============================================================================
UPDATE navigation_item 
SET title = 'Control Center', 
    icon = 'Activity'
WHERE id = 'nav_plat_dashboard' OR (context_scope = 'platform' AND path = '/platform/dashboard');
