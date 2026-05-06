package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func strPtr(s string) *string { return &s }

func TestAddColumnSQL_NotNull(t *testing.T) {
	col := schema.Column{Name: "email", DataType: "text", Nullable: false}
	got := migrate.AddColumnSQL("public.users", col)
	want := "ALTER TABLE public.users ADD COLUMN email text NOT NULL;"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestAddColumnSQL_WithDefault(t *testing.T) {
	col := schema.Column{Name: "active", DataType: "boolean", Nullable: true, Default: strPtr("true")}
	got := migrate.AddColumnSQL("public.users", col)
	want := "ALTER TABLE public.users ADD COLUMN active boolean DEFAULT true;"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestDropColumnSQL(t *testing.T) {
	col := schema.Column{Name: "legacy_field", DataType: "text"}
	got := migrate.DropColumnSQL("public.users", col)
	want := "ALTER TABLE public.users DROP COLUMN legacy_field;"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestColumnDiffSQL_AddedAndRemoved(t *testing.T) {
	diff := schema.ColumnDiff{
		Added:   []schema.Column{{Name: "phone", DataType: "text", Nullable: true}},
		Removed: []schema.Column{{Name: "fax", DataType: "text", Nullable: true}},
	}
	stmts := migrate.ColumnDiffSQL("public.contacts", diff)
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(stmts))
	}
	if stmts[0] != "ALTER TABLE public.contacts ADD COLUMN phone text;" {
		t.Errorf("unexpected add stmt: %q", stmts[0])
	}
	if stmts[1] != "ALTER TABLE public.contacts DROP COLUMN fax;" {
		t.Errorf("unexpected drop stmt: %q", stmts[1])
	}
}

func TestColumnDiffSQL_TypeChange(t *testing.T) {
	diff := schema.ColumnDiff{
		Modified: []schema.ColumnChange{
			{Old: schema.Column{Name: "score", DataType: "integer"}, New: schema.Column{Name: "score", DataType: "bigint"}},
		},
	}
	stmts := migrate.ColumnDiffSQL("public.results", diff)
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	want := "ALTER TABLE public.results ALTER COLUMN score TYPE bigint;"
	if stmts[0] != want {
		t.Errorf("got %q, want %q", stmts[0], want)
	}
}

func TestColumnDiffSQL_SetDefault(t *testing.T) {
	diff := schema.ColumnDiff{
		Modified: []schema.ColumnChange{
			{
				Old: schema.Column{Name: "status", DataType: "text", Default: nil},
				New: schema.Column{Name: "status", DataType: "text", Default: strPtr("'active'")},
			},
		},
	}
	stmts := migrate.ColumnDiffSQL("public.users", diff)
	if len(stmts) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(stmts))
	}
	want := "ALTER TABLE public.users ALTER COLUMN status SET DEFAULT 'active';"
	if stmts[0] != want {
		t.Errorf("got %q, want %q", stmts[0], want)
	}
}
