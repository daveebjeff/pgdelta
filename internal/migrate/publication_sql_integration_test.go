package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func TestPublicationMigration_FullCycle(t *testing.T) {
	old := schema.Schema{
		Publications: []schema.Publication{
			{
				Schema:    "public",
				Name:      "events_pub",
				Tables:    []string{"events"},
				AllTables: false,
				Insert:    true,
				Update:    false,
				Delete:    false,
				Truncate:  false,
			},
		},
	}

	new := schema.Schema{
		Publications: []schema.Publication{
			{
				Schema:    "public",
				Name:      "events_pub",
				Tables:    []string{"events", "logs"},
				AllTables: false,
				Insert:    true,
				Update:    true,
				Delete:    false,
				Truncate:  false,
			},
			{
				Schema:    "public",
				Name:      "all_pub",
				AllTables: true,
				Insert:    true,
				Update:    true,
				Delete:    true,
				Truncate:  true,
			},
		},
	}

	diff := schema.DiffSchemas(old, new)
	sqls := migrate.SchemaDiffSQL(diff)

	if len(sqls) == 0 {
		t.Fatal("expected migration SQL statements, got none")
	}

	foundAlter := false
	foundCreate := false
	for _, sql := range sqls {
		if len(sql) > 5 && sql[:5] == "ALTER" {
			foundAlter = true
		}
		if len(sql) > 6 && sql[:6] == "CREATE" {
			foundCreate = true
		}
	}

	if !foundAlter {
		t.Error("expected an ALTER PUBLICATION statement")
	}
	if !foundCreate {
		t.Error("expected a CREATE PUBLICATION statement")
	}
}
