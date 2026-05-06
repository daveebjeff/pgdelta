package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
)

var baseEnum = schema.Enum{
	Schema: "public",
	Name:   "mood",
	Values: []string{"happy", "sad", "neutral"},
}

func TestCreateEnumSQL(t *testing.T) {
	got := migrate.CreateEnumSQL(baseEnum)
	want := "CREATE TYPE public.mood AS ENUM ('happy', 'sad', 'neutral');"
	if got != want {
		t.Errorf("CreateEnumSQL() = %q, want %q", got, want)
	}
}

func TestDropEnumSQL(t *testing.T) {
	got := migrate.DropEnumSQL(baseEnum)
	want := "DROP TYPE public.mood;"
	if got != want {
		t.Errorf("DropEnumSQL() = %q, want %q", got, want)
	}
}

func TestEnumDiffSQL_AddedAndRemoved(t *testing.T) {
	added := schema.Enum{Schema: "public", Name: "status", Values: []string{"active", "inactive"}}
	diff := schema.EnumDiff{
		Added:   []schema.Enum{added},
		Removed: []schema.Enum{baseEnum},
	}
	stmts := migrate.EnumDiffSQL(diff)
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(stmts))
	}
	if stmts[0] != "DROP TYPE public.mood;" {
		t.Errorf("unexpected drop stmt: %s", stmts[0])
	}
	if stmts[1] != "CREATE TYPE public.status AS ENUM ('active', 'inactive');" {
		t.Errorf("unexpected create stmt: %s", stmts[1])
	}
}

func TestEnumDiffSQL_Changed(t *testing.T) {
	newEnum := schema.Enum{
		Schema: "public",
		Name:   "mood",
		Values: []string{"happy", "sad", "neutral", "excited"},
	}
	diff := schema.EnumDiff{
		Changed: []schema.EnumChange{{Old: baseEnum, New: newEnum}},
	}
	stmts := migrate.EnumDiffSQL(diff)
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	want := "ALTER TYPE public.mood ADD VALUE 'excited';"
	if stmts[0] != want {
		t.Errorf("EnumDiffSQL() = %q, want %q", stmts[0], want)
	}
}
