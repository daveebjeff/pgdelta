package migrate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"pgdelta/internal/migrate"
	"pgdelta/internal/schema"
)

func TestTablespaceMigration_FullCycle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	base := schema.Tablespace{
		Name:  "fastdisk",
		Owner: "postgres",
	}

	// Simulate adding a tablespace
	added := schema.DiffTablespaces([]schema.Tablespace{}, []schema.Tablespace{base})
	require.Len(t, added.Added, 1)
	require.Empty(t, added.Removed)
	require.Empty(t, added.Changed)

	addSQL := migrate.TablespaceDiffSQL(added)
	assert.Contains(t, addSQL, "CREATE TABLESPACE\"fastdisk\"")
	assert.Contains(t, addSQL, "OWNER\"postgres\"")

	// Simulate removing a tablespace
	removed := schema.DiffTablespaces([]schema.Tablespace{base}, []schema.Tablespace{})
	require.Len(t, removed.Removed, 1)
	require.Empty(t, removed.Added)

	dropSQL := migrate.TablespaceDiffSQL(removed)
	assert.Contains(t, dropSQL, "DROP TABLESPACE\"fastdisk\"")

	// Simulate changing owner
	updated := schema.Tablespace{
		Name:  "fastdisk",
		Owner: "newowner",
	}
	changed := schema.DiffTablespaces([]schema.Tablespace{base}, []schema.Tablespace{updated})
	require.Len(t, changed.Changed, 1)

	alterSQL := migrate.TablespaceDiffSQL(changed)
	assert.Contains(t, alterSQL, "ALTER TABLESPACE")
	assert.Contains(t, alterSQL, "newowner")
}
