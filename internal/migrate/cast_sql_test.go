package migrate_test

import (
	"testing"

	"github.com/andyatkinson/pgdelta/internal/migrate"
	"github.com/andyatkinson/pgdelta/internal/schema"
)

func TestCreateCastSQL_Basic(t *testing.T) {
	c := schema.Cast{
		SourceType: "integer",
		TargetType: "text",
		FunctionName: "int4_to_text",
		CastContext: "explicit",
	}
	got := migrate.CreateCastSQL(c)
	want := "CREATE CAST (integer AS text) WITH FUNCTION int4_to_text AS EXPLICIT;"
	if got != want {
		t.Errorf("CreateCastSQL() = %q, want %q", got, want)
	}
}

func TestCreateCastSQL_Implicit(t *testing.T) {
	c := schema.Cast{
		SourceType: "varchar",
		TargetType: "text",
		FunctionName: "varchar_to_text",
		CastContext: "implicit",
	}
	got := migrate.CreateCastSQL(c)
	want := "CREATE CAST (varchar AS text) WITH FUNCTION varchar_to_text AS IMPLICIT;"
	if got != want {
		t.Errorf("CreateCastSQL() = %q, want %q", got, want)
	}
}

func TestDropCastSQL(t *testing.T) {
	c := schema.Cast{
		SourceType: "integer",
		TargetType: "text",
	}
	got := migrate.DropCastSQL(c)
	want := "DROP CAST (integer AS text);"
	if got != want {
		t.Errorf("DropCastSQL() = %q, want %q", got, want)
	}
}

func TestCastDiffSQL_AddedAndRemoved(t *testing.T) {
	added := []schema.Cast{
		{SourceType: "integer", TargetType: "text", FunctionName: "int4_to_text", CastContext: "explicit"},
	}
	removed := []schema.Cast{
		{SourceType: "bigint", TargetType: "text", FunctionName: "int8_to_text", CastContext: "explicit"},
	}
	diff := schema.CastDiff{Added: added, Removed: removed}
	statements := migrate.CastDiffSQL(diff)
	if len(statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(statements))
	}
	if statements[0] != "CREATE CAST (integer AS text) WITH FUNCTION int4_to_text AS EXPLICIT;" {
		t.Errorf("unexpected create statement: %q", statements[0])
	}
	if statements[1] != "DROP CAST (bigint AS text);" {
		t.Errorf("unexpected drop statement: %q", statements[1])
	}
}

func TestCastDiffSQL_Changed(t *testing.T) {
	changed := []schema.Cast{
		{SourceType: "integer", TargetType: "numeric", FunctionName: "int4_to_numeric", CastContext: "assignment"},
	}
	diff := schema.CastDiff{Changed: changed}
	statements := migrate.CastDiffSQL(diff)
	if len(statements) != 2 {
		t.Fatalf("expected 2 statements for changed cast, got %d", len(statements))
	}
}
