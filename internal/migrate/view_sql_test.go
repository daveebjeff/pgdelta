package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
)

func TestCreateViewSQL(t *testing.T) {
	v := schema.View{
		Schema:     "public",
		Name:       "active_users",
		Definition: "SELECT id, name FROM users WHERE active = true",
	}
	got := migrate.CreateViewSQL(v)
	want := "CREATE OR REPLACE VIEW public.active_users AS SELECT id, name FROM users WHERE active = true;"
	if got != want {
		t.Errorf("CreateViewSQL() = %q, want %q", got, want)
	}
}

func TestDropViewSQL(t *testing.T) {
	v := schema.View{
		Schema: "public",
		Name:   "active_users",
	}
	got := migrate.DropViewSQL(v)
	want := "DROP VIEW IF EXISTS public.active_users;"
	if got != want {
		t.Errorf("DropViewSQL() = %q, want %q", got, want)
	}
}

func TestViewDiffSQL_AddedAndRemoved(t *testing.T) {
	added := []schema.View{
		{Schema: "public", Name: "new_view", Definition: "SELECT 1"},
	}
	removed := []schema.View{
		{Schema: "public", Name: "old_view", Definition: "SELECT 2"},
	}
	changed := []schema.View{}

	sqls := migrate.ViewDiffSQL(schema.ViewDiff{
		Added:   added,
		Removed: removed,
		Changed: changed,
	})

	if len(sqls) != 2 {
		t.Fatalf("expected 2 SQL statements, got %d", len(sqls))
	}
	if sqls[0] != "DROP VIEW IF EXISTS public.old_view;" {
		t.Errorf("unexpected drop SQL: %q", sqls[0])
	}
	if sqls[1] != "CREATE OR REPLACE VIEW public.new_view AS SELECT 1;" {
		t.Errorf("unexpected create SQL: %q", sqls[1])
	}
}

func TestViewDiffSQL_Changed(t *testing.T) {
	changed := []schema.View{
		{Schema: "reporting", Name: "summary", Definition: "SELECT count(*) FROM orders"},
	}

	sqls := migrate.ViewDiffSQL(schema.ViewDiff{
		Added:   []schema.View{},
		Removed: []schema.View{},
		Changed: changed,
	})

	if len(sqls) != 1 {
		t.Fatalf("expected 1 SQL statement, got %d", len(sqls))
	}
	want := "CREATE OR REPLACE VIEW reporting.summary AS SELECT count(*) FROM orders;"
	if sqls[0] != want {
		t.Errorf("unexpected changed SQL: %q", sqls[0])
	}
}
