package migrate_test

import (
	"strings"
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

// TestSchemaMigration_FullCycle verifies that a full schema diff produces
// a coherent, ordered migration script covering multiple object types.
func TestSchemaMigration_FullCycle(t *testing.T) {
	old := schema.Schema{
		Enums: []schema.Enum{
			{Schema: "public", Name: "old_status", Values: []string{"a", "b"}},
		},
		Sequences: []schema.Sequence{
			{Schema: "public", Name: "legacy_seq", IncrementBy: 1, MinValue: 1, MaxValue: 9999, StartValue: 1, Cache: 1},
		},
	}

	new := schema.Schema{
		Extensions: []schema.Extension{
			{Schema: "public", Name: "pgcrypto", Version: "1.3"},
		},
		Enums: []schema.Enum{
			{Schema: "public", Name: "new_status", Values: []string{"active", "inactive"}},
		},
		Tables: []schema.Table{
			{
				Schema: "public",
				Name:   "users",
				Columns: []schema.Column{
					{Schema: "public", TableName: "users", Name: "id", DataType: "integer", IsNullable: false},
					{Schema: "public", TableName: "users", Name: "email", DataType: "text", IsNullable: false},
				},
			},
		},
	}

	diff := schema.DiffSchemas(old, new)
	sql := migrate.SchemaDiffSQL(diff)

	if sql == "" {
		t.Fatal("expected non-empty migration SQL")
	}

	expected := []string{
		"CREATE EXTENSION",
		"pgcrypto",
		"DROP TYPE",
		"old_status",
		"CREATE TYPE",
		"new_status",
		"DROP SEQUENCE",
		"legacy_seq",
		"CREATE TABLE",
		"users",
	}

	for _, want := range expected {
		if !strings.Contains(sql, want) {
			t.Errorf("expected %q in migration SQL\nFull output:\n%s", want, sql)
		}
	}

	// Verify ordering: extensions before enums before tables
	extIdx := strings.Index(sql, "CREATE EXTENSION")
	enumIdx := strings.Index(sql, "CREATE TYPE")
	tableIdx := strings.Index(sql, "CREATE TABLE")

	if extIdx > enumIdx {
		t.Error("expected CREATE EXTENSION before CREATE TYPE")
	}
	if enumIdx > tableIdx {
		t.Error("expected CREATE TYPE before CREATE TABLE")
	}
}
