package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSchemaNamespaceMigration_FullCycle(t *testing.T) {
	db := getTestDB(t)

	ns := schema.SchemaNamespace{Name: "pgdelta_test_ns", Owner: "postgres"}

	// Create
	createSQL := migrate.CreateSchemaSQL(ns)
	_, err := db.Exec(createSQL)
	require.NoError(t, err, "create schema should succeed")

	// Verify created
	var count int
	err = db.QueryRow(
		"SELECT count(*) FROM information_schema.schemata WHERE schema_name = $1",
		ns.Name,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 1, count, "schema should exist after creation")

	// Drop
	dropSQL := migrate.DropSchemaSQL(ns)
	_, err = db.Exec(dropSQL)
	require.NoError(t, err, "drop schema should succeed")

	// Verify dropped
	err = db.QueryRow(
		"SELECT count(*) FROM information_schema.schemata WHERE schema_name = $1",
		ns.Name,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "schema should not exist after drop")
}
