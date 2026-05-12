package migrate_test

import (
	"context"
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeMigration_FullCycle(t *testing.T) {
	ctx := context.Background()
	db := openTestDB(t)
	defer db.Close()

	_, err := db.ExecContext(ctx, "CREATE SCHEMA IF NOT EXISTS pgdelta_test")
	require.NoError(t, err)
	t.Cleanup(func() { db.ExecContext(ctx, "DROP SCHEMA pgdelta_test CASCADE") })

	compositeType := schema.Type{
		Schema:     "pgdelta_test",
		Name:       "address",
		Kind:       "composite",
		Definition: "(street text, city text, zip text)",
	}

	// Create the type.
	createSQL := migrate.CreateTypeSQL(compositeType)
	_, err = db.ExecContext(ctx, createSQL)
	require.NoError(t, err, "CreateTypeSQL should execute without error")

	// Verify it exists.
	var count int
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM pg_type t JOIN pg_namespace n ON n.oid = t.typnamespace WHERE n.nspname = $1 AND t.typname = $2",
		"pgdelta_test", "address",
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "type should exist after creation")

	// Drop the type.
	dropSQL := migrate.DropTypeSQL(compositeType)
	_, err = db.ExecContext(ctx, dropSQL)
	require.NoError(t, err, "DropTypeSQL should execute without error")

	// Verify it no longer exists.
	err = db.QueryRowContext(ctx,
		"SELECT COUNT(*) FROM pg_type t JOIN pg_namespace n ON n.oid = t.typnamespace WHERE n.nspname = $1 AND t.typname = $2",
		"pgdelta_test", "address",
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "type should not exist after drop")
}
