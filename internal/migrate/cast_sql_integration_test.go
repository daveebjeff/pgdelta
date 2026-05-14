package migrate_test

import (
	"testing"

	"github.com/andyatkinson/pgdelta/internal/migrate"
	"github.com/andyatkinson/pgdelta/internal/schema"
)

func TestCastMigration_FullCycle(t *testing.T) {
	// Simulate a full cycle: add a cast, then remove it
	added := schema.Cast{
		SourceType:   "integer",
		TargetType:   "text",
		FunctionName: "int4_to_text",
		CastContext:  "explicit",
	}

	createSQL := migrate.CreateCastSQL(added)
	expectedCreate := "CREATE CAST (integer AS text) WITH FUNCTION int4_to_text AS EXPLICIT;"
	if createSQL != expectedCreate {
		t.Errorf("CreateCastSQL() = %q, want %q", createSQL, expectedCreate)
	}

	dropSQL := migrate.DropCastSQL(added)
	expectedDrop := "DROP CAST (integer AS text);"
	if dropSQL != expectedDrop {
		t.Errorf("DropCastSQL() = %q, want %q", dropSQL, expectedDrop)
	}

	// Simulate diff with both added and removed
	diff := schema.CastDiff{
		Added:   []schema.Cast{added},
		Removed: []schema.Cast{{SourceType: "bigint", TargetType: "text", FunctionName: "int8_to_text", CastContext: "explicit"}},
	}
	statements := migrate.CastDiffSQL(diff)
	if len(statements) != 2 {
		t.Fatalf("expected 2 migration statements, got %d", len(statements))
	}

	// Verify ordering: creates before drops
	if statements[0] != expectedCreate {
		t.Errorf("first statement should be CREATE, got %q", statements[0])
	}
	if statements[1] != "DROP CAST (bigint AS text);" {
		t.Errorf("second statement should be DROP, got %q", statements[1])
	}
}
