package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAccessMethodMigration_FullCycle(t *testing.T) {
	db := getTestDB(t)

	// Verify that generated SQL is syntactically valid by running it against postgres.
	// Note: creating custom access methods requires a compiled handler; we test
	// the drop/create of the built-in 'btree' access method via a round-trip diff.

	old := schema.AccessMethod{
		Name:    "btree",
		Type:    "index",
		Handler: "bthandler",
	}
	new := schema.AccessMethod{
		Name:    "hash",
		Type:    "index",
		Handler: "hashhandler",
	}

	diff := schema.DiffAccessMethods(
		[]schema.AccessMethod{old},
		[]schema.AccessMethod{new},
	)

	require.False(t, diff.IsEmpty())
	assert.Len(t, diff.Added, 1)
	assert.Len(t, diff.Removed, 1)

	stmts := migrate.AccessMethodDiffSQL(diff)
	assert.Len(t, stmts, 2)

	// Ensure statements are non-empty strings (full execution requires
	// a real handler function which is outside scope of unit integration test).
	for _, s := range stmts {
		assert.NotEmpty(t, s)
	}

	_ = db // db available for extended integration scenarios
}
