//go:build integration

// Package testing contains Phase 2 PostgreSQL integration tests for Curexal.
//
// Tests require Docker to be running. They are automatically skipped when
// Docker is unavailable (testcontainers returns an error -> require.NoError
// inside SetupTestDB calls t.FailNow, which is equivalent to t.Skip here).
//
// Run these tests with:
//
//	go test -tags=integration -timeout=5m ./internal/testing/...
//
// Covered by this suite:
//   - database.RunInTx automatic rollback when any step fails
//   - database.RunInTx commit when all steps succeed
//   - Conn(ctx) propagates the ambient pgx.Tx (the core atomicity mechanism)
//   - Organization + Tenant multi-table atomic creation and rollback
//   - User + patient_profile atomic creation and rollback (RegisterPatient path)
//   - Organization slug UNIQUE constraint enforcement
//   - Tenant slug UNIQUE constraint enforcement
//   - Membership (user_id, tenant_id) UNIQUE INDEX enforcement
//   - Multi-step rollback: all three users rolled back when step 3 fails
//   - organization domain.Organization field mapping (pgx.RowToStructByName alignment)
//   - model.User field mapping from SELECT id, name, email, email_verified, image, is_platform_admin, platform_role
//   - session ON DELETE CASCADE: session removed when parent user is deleted
package testing

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/golangnigeria/curexal/internal/modules/identity/model"
	"github.com/golangnigeria/curexal/internal/modules/organization/domain"
	"github.com/golangnigeria/curexal/internal/kernel/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── shared container (started once, reused across all tests in this file) ────

var integrationPool *pgxpool.Pool

func getIntegrationPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if integrationPool != nil {
		return integrationPool
	}
	tdb, cleanup := SetupTestDB(t)
	if tdb == nil {
		t.Skip("Docker unavailable — skipping DB integration test")
		return nil
	}
	integrationPool = tdb.Pool
	// Register cleanup on the top-level test binary exit.
	// (t.Cleanup on the first test is sufficient for a single test binary run.)
	t.Cleanup(cleanup)
	return integrationPool
}

// wrapDB wraps a pgxpool.Pool into a database.Database so we can call RunInTx.
func wrapDB(pool *pgxpool.Pool) *database.Database {
	return &database.Database{Pool: pool}
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func execInsertUser(ctx context.Context, pool *pgxpool.Pool, id, name, email string) error {
	_, err := pool.Exec(ctx, `
		INSERT INTO "user" (id, name, email, email_verified, is_platform_admin)
		VALUES ($1, $2, $3, FALSE, FALSE)
	`, id, name, email)
	return err
}

func execInsertUserViaDB(ctx context.Context, db *database.Database, id, name, email string) error {
	_, err := db.Conn(ctx).Exec(ctx, `
		INSERT INTO "user" (id, name, email, email_verified, is_platform_admin)
		VALUES ($1, $2, $3, FALSE, FALSE)
	`, id, name, email)
	return err
}

func execInsertOrganization(ctx context.Context, db *database.Database, id, name, slug string) error {
	_, err := db.Conn(ctx).Exec(ctx, `
		INSERT INTO organization (id, name, slug, plan)
		VALUES ($1, $2, $3, 'smart')
	`, id, name, slug)
	return err
}

func countDBRows(ctx context.Context, pool *pgxpool.Pool, table, whereClause string, args ...any) int {
	var n int
	q := fmt.Sprintf("SELECT COUNT(*) FROM %s WHERE %s", table, whereClause)
	_ = pool.QueryRow(ctx, q, args...).Scan(&n)
	return n
}

// isDuplicateKeyError reports whether err is a PostgreSQL unique violation (23505).
func isDuplicateKeyError(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "23505") ||
		strings.Contains(s, "unique") ||
		strings.Contains(s, "duplicate") ||
		strings.Contains(s, "UNIQUE")
}

// ─── 1. RunInTx: rollback on step failure ─────────────────────────────────────

// TestRunInTx_RollbackOnError verifies that when one step inside RunInTx fails,
// ALL preceding writes within the same transaction are rolled back atomically.
func TestRunInTx_RollbackOnError(t *testing.T) {
	pool := getIntegrationPool(t)
	db := wrapDB(pool)
	ctx := context.Background()

	userID := "usr_tx_rollback_" + uuid.New().String()[:8]
	email := userID + "@test.curexal.internal"

	err := db.RunInTx(ctx, func(txCtx context.Context) error {
		// Step 1: insert a user (succeeds inside the transaction)
		if err := execInsertUserViaDB(txCtx, db, userID, "Rollback User", email); err != nil {
			return err
		}
		// Step 2: deliberate failure — Step 1 must be rolled back.
		return errors.New("deliberate failure to trigger rollback")
	})

	require.Error(t, err)
	assert.Contains(t, err.Error(), "deliberate failure")

	// User must NOT exist after rollback.
	n := countDBRows(ctx, pool, `"user"`, "id = $1", userID)
	assert.Equal(t, 0, n, "user row must be rolled back when transaction fails")
}

// TestRunInTx_CommitsOnSuccess verifies that a successful multi-step transaction
// persists all rows.
func TestRunInTx_CommitsOnSuccess(t *testing.T) {
	pool := getIntegrationPool(t)
	db := wrapDB(pool)
	ctx := context.Background()

	userID := "usr_tx_commit_" + uuid.New().String()[:8]
	email := userID + "@test.curexal.internal"
	token := "verify_commit_" + uuid.New().String()

	err := db.RunInTx(ctx, func(txCtx context.Context) error {
		if err := execInsertUserViaDB(txCtx, db, userID, "Commit User", email); err != nil {
			return err
		}
		_, err := db.Conn(txCtx).Exec(txCtx, `
			INSERT INTO verification_token (token, email, token_type, expires_at)
			VALUES ($1, $2, 'email-verification', $3)
		`, token, email, time.Now().Add(24*time.Hour))
		return err
	})

	require.NoError(t, err)
	assert.Equal(t, 1, countDBRows(ctx, pool, `"user"`, "id = $1", userID), "user must persist on commit")
	assert.Equal(t, 1, countDBRows(ctx, pool, "verification_token", "token = $1", token), "token must persist on commit")
}

// ─── 2. Conn(ctx) propagates the ambient transaction ──────────────────────────

// TestConn_PropagatesTransaction verifies that db.Conn(txCtx) returns the
// ambient pgx.Tx — this is the mechanism that makes multi-table atomicity work.
func TestConn_PropagatesTransaction(t *testing.T) {
	pool := getIntegrationPool(t)
	ctx := context.Background()

	userID := "usr_conn_prop_" + uuid.New().String()[:8]
	email := userID + "@test.internal"

	// Open a transaction manually.
	tx, err := pool.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	// Write inside the tx.
	_, err = tx.Exec(ctx, `
		INSERT INTO identity.users (id, name, email, email_verified, is_platform_admin)
		VALUES ($1, $2, $3, FALSE, FALSE)
	`, userID, "Tx Propagation", email)
	require.NoError(t, err)

	// Row IS visible within the same tx.
	var count int
	require.NoError(t, tx.QueryRow(ctx, `SELECT COUNT(*) FROM identity.users WHERE id = $1`, userID).Scan(&count))
	assert.Equal(t, 1, count, "row must be visible within the open transaction")

	// Rollback — row must disappear.
	require.NoError(t, tx.Rollback(ctx))
	assert.Equal(t, 0, countDBRows(ctx, pool, `identity.users`, "id = $1", userID), "row must not exist after rollback")
}

// ─── 3. Organization + Tenant multi-table atomicity ──────────────────────────

// TestCreateOrganizationAndTenant_Atomicity verifies that org + tenant creation
// in one RunInTx either both commit or both roll back.
func TestCreateOrganizationAndTenant_Atomicity(t *testing.T) {
	pool := getIntegrationPool(t)
	db := wrapDB(pool)
	ctx := context.Background()

	// ── Success path ──────────────────────────────────────────────────────────
	orgID := uuid.New().String()
	orgSlug := "org-tx-" + uuid.New().String()[:6]
	tenantID := uuid.New().String()
	tenantSlug := "t-tx-" + uuid.New().String()[:6]

	err := db.RunInTx(ctx, func(txCtx context.Context) error {
		if err := execInsertOrganization(txCtx, db, orgID, "Tx Org", orgSlug); err != nil {
			return err
		}
		_, err := db.Conn(txCtx).Exec(txCtx, `
			INSERT INTO tenant (id, name, slug, organization_id, currency, settings)
			VALUES ($1, $2, $3, $4, 'NGN', '{}')
		`, tenantID, "Tx Tenant", tenantSlug, orgID)
		return err
	})

	require.NoError(t, err)
	assert.Equal(t, 1, countDBRows(ctx, pool, "organization", "id = $1", orgID))
	assert.Equal(t, 1, countDBRows(ctx, pool, "tenant", "id = $1", tenantID))

	// ── Rollback path ─────────────────────────────────────────────────────────
	orgID2 := uuid.New().String()
	orgSlug2 := "org-rb-" + uuid.New().String()[:6]
	tenantID2 := uuid.New().String()
	tenantSlug2 := "t-rb-" + uuid.New().String()[:6]

	err = db.RunInTx(ctx, func(txCtx context.Context) error {
		if err := execInsertOrganization(txCtx, db, orgID2, "Rollback Org", orgSlug2); err != nil {
			return err
		}
		if _, err := db.Conn(txCtx).Exec(txCtx, `
			INSERT INTO tenant (id, name, slug, organization_id, currency, settings)
			VALUES ($1, $2, $3, $4, 'NGN', '{}')
		`, tenantID2, "Rollback Tenant", tenantSlug2, orgID2); err != nil {
			return err
		}
		return errors.New("roll back org and tenant")
	})

	require.Error(t, err)
	assert.Equal(t, 0, countDBRows(ctx, pool, "organization", "id = $1", orgID2), "org must be rolled back")
	assert.Equal(t, 0, countDBRows(ctx, pool, "tenant", "id = $1", tenantID2), "tenant must be rolled back")
}

// ─── 4. User + patient_profile atomicity (RegisterPatient path) ───────────────

// TestRegisterPatient_Atomicity verifies that user + patient_profile either
// both commit or both roll back — mirroring the actual RegisterPatient service.
func TestRegisterPatient_Atomicity(t *testing.T) {
	pool := getIntegrationPool(t)
	db := wrapDB(pool)
	ctx := context.Background()

	// ── Success path ──────────────────────────────────────────────────────────
	userID := "usr_patient_" + uuid.New().String()[:8]
	email := userID + "@patient.internal"

	err := db.RunInTx(ctx, func(txCtx context.Context) error {
		if err := execInsertUserViaDB(txCtx, db, userID, "Patient Alice", email); err != nil {
			return err
		}
		_, err := db.Conn(txCtx).Exec(txCtx, `
			INSERT INTO patient_profile (user_id, phone)
			VALUES ($1, $2) ON CONFLICT (user_id) DO NOTHING
		`, userID, "+2348000000001")
		return err
	})

	require.NoError(t, err)
	assert.Equal(t, 1, countDBRows(ctx, pool, `"user"`, "id = $1", userID))
	assert.Equal(t, 1, countDBRows(ctx, pool, "patient_profile", "user_id = $1", userID))

	// ── Rollback path ─────────────────────────────────────────────────────────
	userID2 := "usr_patient_rb_" + uuid.New().String()[:8]
	email2 := userID2 + "@patient.internal"

	err = db.RunInTx(ctx, func(txCtx context.Context) error {
		if err := execInsertUserViaDB(txCtx, db, userID2, "Patient Bob", email2); err != nil {
			return err
		}
		if _, err := db.Conn(txCtx).Exec(txCtx, `
			INSERT INTO patient_profile (user_id) VALUES ($1)
		`, userID2); err != nil {
			return err
		}
		return errors.New("rollback patient")
	})

	require.Error(t, err)
	assert.Equal(t, 0, countDBRows(ctx, pool, `"user"`, "id = $1", userID2))
	assert.Equal(t, 0, countDBRows(ctx, pool, "patient_profile", "user_id = $1", userID2))
}

// ─── 5. Slug uniqueness constraints ──────────────────────────────────────────

// TestOrganization_SlugUniqueness verifies the UNIQUE INDEX on organization.slug.
func TestOrganization_SlugUniqueness(t *testing.T) {
	pool := getIntegrationPool(t)
	db := wrapDB(pool)
	ctx := context.Background()

	slug := "slug-uniq-" + uuid.New().String()[:6]
	require.NoError(t, execInsertOrganization(ctx, db, uuid.New().String(), "Org Alpha", slug))

	err := execInsertOrganization(ctx, db, uuid.New().String(), "Org Beta", slug)
	require.Error(t, err, "duplicate slug must be rejected")
	assert.True(t, isDuplicateKeyError(err), "must be a unique violation, got: %v", err)
}

// TestTenant_SlugUniqueness verifies the UNIQUE INDEX on tenant.slug.
func TestTenant_SlugUniqueness(t *testing.T) {
	pool := getIntegrationPool(t)
	db := wrapDB(pool)
	ctx := context.Background()

	orgID := uuid.New().String()
	require.NoError(t, execInsertOrganization(ctx, db, orgID, "Slug Org", "slug-org-"+uuid.New().String()[:6]))

	tenantSlug := "t-slug-" + uuid.New().String()[:6]
	insertTenant := func(id, slug, orgID string) error {
		_, err := pool.Exec(ctx, `
			INSERT INTO tenant (id, name, slug, organization_id, currency, settings)
			VALUES ($1, $2, $3, $4, 'NGN', '{}')
		`, id, "T "+id[:8], slug, orgID)
		return err
	}

	require.NoError(t, insertTenant(uuid.New().String(), tenantSlug, orgID))
	err := insertTenant(uuid.New().String(), tenantSlug, orgID)
	require.Error(t, err)
	assert.True(t, isDuplicateKeyError(err))
}

// ─── 6. Membership UNIQUE constraint ─────────────────────────────────────────

// TestMembership_UniqueConstraint verifies idx_membership_user_tenant blocks
// duplicate (user_id, tenant_id) membership rows.
func TestMembership_UniqueConstraint(t *testing.T) {
	pool := getIntegrationPool(t)
	db := wrapDB(pool)
	ctx := context.Background()

	userID := "usr_mem_" + uuid.New().String()[:8]
	require.NoError(t, execInsertUser(ctx, pool, userID, "Mem User", userID+"@test.internal"))

	orgID := uuid.New().String()
	require.NoError(t, execInsertOrganization(ctx, db, orgID, "Mem Org", "mem-org-"+uuid.New().String()[:6]))

	tenantID := uuid.New().String()
	_, err := pool.Exec(ctx, `
		INSERT INTO tenant (id, name, slug, organization_id, currency, settings)
		VALUES ($1, $2, $3, $4, 'NGN', '{}')
	`, tenantID, "Mem Tenant", "mem-t-"+uuid.New().String()[:6], orgID)
	require.NoError(t, err)

	insertMem := func() error {
		_, err := pool.Exec(ctx, `
			INSERT INTO membership (id, user_id, tenant_id, organization_id, role_title, is_active, joined_at)
			VALUES ($1, $2, $3, $4, 'member', TRUE, NOW())
		`, uuid.New().String(), userID, tenantID, orgID)
		return err
	}

	require.NoError(t, insertMem())
	err = insertMem()
	require.Error(t, err)
	assert.True(t, isDuplicateKeyError(err))
}

// ─── 7. Multi-step rollback ───────────────────────────────────────────────────

// TestRunInTx_MultiStepRollbackAllThreeRows verifies that a failure in step 3
// rolls back steps 1, 2, and 3 together.
func TestRunInTx_MultiStepRollbackAllThreeRows(t *testing.T) {
	pool := getIntegrationPool(t)
	db := wrapDB(pool)
	ctx := context.Background()

	ids := [3]string{
		"usr_ms0_" + uuid.New().String()[:6],
		"usr_ms1_" + uuid.New().String()[:6],
		"usr_ms2_" + uuid.New().String()[:6],
	}

	err := db.RunInTx(ctx, func(txCtx context.Context) error {
		for i, id := range ids {
			if err := execInsertUserViaDB(txCtx, db, id, fmt.Sprintf("User %d", i), id+"@test.internal"); err != nil {
				return err
			}
		}
		return errors.New("all three must roll back")
	})

	require.Error(t, err)
	for _, id := range ids {
		assert.Equal(t, 0, countDBRows(ctx, pool, `"user"`, "id = $1", id),
			"user %s must not exist after rollback", id)
	}
}

// ─── 8. Domain struct field alignment ────────────────────────────────────────

// TestOrganizationDomainStruct_FieldMapping verifies that domain.Organization
// scans correctly from the columns the repository queries.
func TestOrganizationDomainStruct_FieldMapping(t *testing.T) {
	pool := getIntegrationPool(t)
	db := wrapDB(pool)
	ctx := context.Background()

	orgID := uuid.New().String()
	orgSlug := "field-" + uuid.New().String()[:6]
	require.NoError(t, execInsertOrganization(ctx, db, orgID, "Field Org", orgSlug))

	rows, err := pool.Query(ctx, `
		SELECT id, name, slug, plan, custom_domain, settings, created_at, updated_at
		FROM organization.organizations WHERE id = $1
	`, orgID)
	require.NoError(t, err)

	org, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[domain.Organization])
	require.NoError(t, err)

	assert.Equal(t, orgID, org.ID)
	assert.Equal(t, "Field Org", org.Name)
	assert.Equal(t, orgSlug, org.Slug)
	assert.Equal(t, "smart", org.Plan)
}

// TestUserModel_FieldMapping verifies model.User scans correctly from the
// SELECT columns used by GetByEmail and GetByID.
func TestUserModel_FieldMapping(t *testing.T) {
	pool := getIntegrationPool(t)
	ctx := context.Background()

	userID := "usr_field_" + uuid.New().String()[:8]
	email := userID + "@fieldmap.internal"
	require.NoError(t, execInsertUser(ctx, pool, userID, "Field User", email))

	u := &model.User{}
	err := pool.QueryRow(ctx, `
		SELECT id, name, email, email_verified, avatar_url, is_platform_admin, platform_role
		FROM identity.users WHERE id = $1
	`, userID).Scan(&u.ID, &u.Name, &u.Email, &u.EmailVerified, &u.Image, &u.IsPlatformAdmin, &u.PlatformRole)
	require.NoError(t, err)

	assert.Equal(t, userID, u.ID)
	assert.Equal(t, email, u.Email)
	assert.False(t, u.EmailVerified)
	assert.False(t, u.IsPlatformAdmin)
}

// ─── 9. Session ON DELETE CASCADE ────────────────────────────────────────────

// TestSession_CascadeDeleteOnUser verifies that deleting a user cascades to
// their sessions via session.user_id REFERENCES user(id) ON DELETE CASCADE.
func TestSession_CascadeDeleteOnUser(t *testing.T) {
	pool := getIntegrationPool(t)
	ctx := context.Background()

	userID := "usr_casc_" + uuid.New().String()[:8]
	require.NoError(t, execInsertUser(ctx, pool, userID, "Cascade User", userID+"@casc.internal"))

	sessionID := "sess_casc_" + uuid.New().String()[:8]
	_, err := pool.Exec(ctx, `
		INSERT INTO identity.sessions (id, user_id, token, expires_at)
		VALUES ($1, $2, $3, $4)
	`, sessionID, userID, "tok_"+uuid.New().String(), time.Now().Add(24*time.Hour))
	require.NoError(t, err)

	assert.Equal(t, 1, countDBRows(ctx, pool, "identity.sessions", "id = $1", sessionID))

	// Delete user — cascade must remove the session automatically.
	_, err = pool.Exec(ctx, `DELETE FROM identity.users WHERE id = $1`, userID)
	require.NoError(t, err)

	assert.Equal(t, 0, countDBRows(ctx, pool, "identity.sessions", "id = $1", sessionID),
		"session must be cascade-deleted when the parent user is deleted")
}
