package migrate_test

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoleMigration_FullCycle(t *testing.T) {
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("TEST_DATABASE_URL not set")
	}

	db, err := sql.Open("postgres", dsn)
	require.NoError(t, err)
	defer db.Close()

	roleName := "pgdelta_test_role"

	// Cleanup before and after
	_, _ = db.Exec(fmt.Sprintf(`DROP ROLE IF EXISTS %q`, roleName))
	defer db.Exec(fmt.Sprintf(`DROP ROLE IF EXISTS %q`, roleName))

	r := schema.Role{
		Name:            roleName,
		Inherit:         true,
		Login:           true,
		ConnectionLimit: 5,
	}

	// Create
	createSQL := migrate.CreateRoleSQL(r)
	_, err = db.Exec(createSQL)
	require.NoError(t, err, "CreateRoleSQL failed: %s", createSQL)

	// Verify exists
	var exists bool
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = $1)`, roleName).Scan(&exists)
	require.NoError(t, err)
	assert.True(t, exists, "role should exist after creation")

	// Alter
	r.CreateDB = true
	alterSQL := migrate.AlterRoleSQL(r)
	_, err = db.Exec(alterSQL)
	require.NoError(t, err, "AlterRoleSQL failed: %s", alterSQL)

	// Verify altered
	var createdb bool
	err = db.QueryRow(`SELECT rolcreatedb FROM pg_roles WHERE rolname = $1`, roleName).Scan(&createdb)
	require.NoError(t, err)
	assert.True(t, createdb, "role should have CREATEDB after alter")

	// Drop
	dropSQL := migrate.DropRoleSQL(r)
	_, err = db.Exec(dropSQL)
	require.NoError(t, err, "DropRoleSQL failed: %s", dropSQL)

	// Verify removed
	err = db.QueryRow(`SELECT EXISTS(SELECT 1 FROM pg_roles WHERE rolname = $1)`, roleName).Scan(&exists)
	require.NoError(t, err)
	assert.False(t, exists, "role should not exist after drop")
}
