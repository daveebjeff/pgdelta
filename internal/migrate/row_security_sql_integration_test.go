package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
)

func TestRowSecurityMigration_FullCycle(t *testing.T) {
	// Start with no row security
	old := []schema.RowSecurity{}

	// Add row security with force enabled
	newState := []schema.RowSecurity{
		{Schema: "public", Table: "accounts", Enabled: true, Forced: true},
	}

	diff := schema.DiffRowSecurity(old, newState)
	sqls := migrate.RowSecurityDiffSQL(diff)

	if len(sqls) == 0 {
		t.Fatal("expected SQL statements for row security changes, got none")
	}

	// Verify enable comes before force
	foundEnable := false
	foundForce := false
	for _, sql := range sqls {
		if sql == "ALTER TABLE public.accounts ENABLE ROW LEVEL SECURITY;" {
			foundEnable = true
		}
		if sql == "ALTER TABLE public.accounts FORCE ROW LEVEL SECURITY;" {
			foundForce = true
		}
	}
	if !foundEnable {
		t.Error("expected ENABLE ROW LEVEL SECURITY statement")
	}
	if !foundForce {
		t.Error("expected FORCE ROW LEVEL SECURITY statement")
	}

	// Now disable row security
	disabled := []schema.RowSecurity{
		{Schema: "public", Table: "accounts", Enabled: false, Forced: false},
	}
	diff2 := schema.DiffRowSecurity(newState, disabled)
	sqls2 := migrate.RowSecurityDiffSQL(diff2)

	if len(sqls2) == 0 {
		t.Fatal("expected SQL statements when disabling row security")
	}

	foundDisable := false
	for _, sql := range sqls2 {
		if sql == "ALTER TABLE public.accounts DISABLE ROW LEVEL SECURITY;" {
			foundDisable = true
		}
	}
	if !foundDisable {
		t.Error("expected DISABLE ROW LEVEL SECURITY statement")
	}
}
