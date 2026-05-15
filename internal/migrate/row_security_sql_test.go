package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
)

func TestEnableRowSecuritySQL(t *testing.T) {
	rs := schema.RowSecurity{
		Schema: "public",
		Table:  "orders",
		Enabled: true,
		Forced:  false,
	}
	got := migrate.EnableRowSecuritySQL(rs)
	want := "ALTER TABLE public.orders ENABLE ROW LEVEL SECURITY;"
	if got != want {
		t.Errorf("EnableRowSecuritySQL() = %q, want %q", got, want)
	}
}

func TestDisableRowSecuritySQL(t *testing.T) {
	rs := schema.RowSecurity{
		Schema: "public",
		Table:  "orders",
		Enabled: false,
		Forced:  false,
	}
	got := migrate.DisableRowSecuritySQL(rs)
	want := "ALTER TABLE public.orders DISABLE ROW LEVEL SECURITY;"
	if got != want {
		t.Errorf("DisableRowSecuritySQL() = %q, want %q", got, want)
	}
}

func TestForceRowSecuritySQL(t *testing.T) {
	rs := schema.RowSecurity{
		Schema: "public",
		Table:  "orders",
		Enabled: true,
		Forced:  true,
	}
	got := migrate.ForceRowSecuritySQL(rs)
	want := "ALTER TABLE public.orders FORCE ROW LEVEL SECURITY;"
	if got != want {
		t.Errorf("ForceRowSecuritySQL() = %q, want %q", got, want)
	}
}

func TestNoForceRowSecuritySQL(t *testing.T) {
	rs := schema.RowSecurity{
		Schema: "public",
		Table:  "orders",
		Enabled: true,
		Forced:  false,
	}
	got := migrate.NoForceRowSecuritySQL(rs)
	want := "ALTER TABLE public.orders NO FORCE ROW LEVEL SECURITY;"
	if got != want {
		t.Errorf("NoForceRowSecuritySQL() = %q, want %q", got, want)
	}
}

func TestRowSecurityDiffSQL_EnabledAndForced(t *testing.T) {
	old := []schema.RowSecurity{}
	new := []schema.RowSecurity{
		{Schema: "public", Table: "orders", Enabled: true, Forced: true},
	}
	diff := schema.DiffRowSecurity(old, new)
	sqls := migrate.RowSecurityDiffSQL(diff)
	if len(sqls) != 2 {
		t.Fatalf("expected 2 SQL statements, got %d", len(sqls))
	}
}

func TestRowSecurityDiffSQL_NoChanges(t *testing.T) {
	rs := []schema.RowSecurity{
		{Schema: "public", Table: "orders", Enabled: true, Forced: false},
	}
	diff := schema.DiffRowSecurity(rs, rs)
	sqls := migrate.RowSecurityDiffSQL(diff)
	if len(sqls) != 0 {
		t.Errorf("expected 0 SQL statements, got %d", len(sqls))
	}
}
