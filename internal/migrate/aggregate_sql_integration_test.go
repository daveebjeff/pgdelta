package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

func TestAggregateMigration_FullCycle(t *testing.T) {
	initCond := "0"
	agg := schema.Aggregate{
		Schema:    "public",
		Name:      "total_count",
		ArgTypes:  []string{"integer"},
		SFuncName: "int4pl",
		SType:     "integer",
		InitCond:  &initCond,
	}

	// Simulate: aggregate does not exist in old schema, added in new
	oldSchema := schema.Schema{}
	newSchema := schema.Schema{Aggregates: []schema.Aggregate{agg}}

	diff := schema.DiffSchemas(oldSchema, newSchema)
	stmts := migrate.SchemaDiffSQL(diff)

	assert.NotEmpty(t, stmts)
	found := false
	for _, s := range stmts {
		if contains(s, "CREATE AGGREGATE") && contains(s, "total_count") {
			found = true
			break
		}
	}
	assert.True(t, found, "expected CREATE AGGREGATE statement in migration output")

	// Simulate: aggregate removed
	diff2 := schema.DiffSchemas(newSchema, oldSchema)
	stmts2 := migrate.SchemaDiffSQL(diff2)

	assert.NotEmpty(t, stmts2)
	foundDrop := false
	for _, s := range stmts2 {
		if contains(s, "DROP AGGREGATE") && contains(s, "total_count") {
			foundDrop = true
			break
		}
	}
	assert.True(t, foundDrop, "expected DROP AGGREGATE statement in migration output")
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		})())
}
