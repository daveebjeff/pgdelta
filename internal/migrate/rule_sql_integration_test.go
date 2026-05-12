package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func TestRuleMigration_FullCycle(t *testing.T) {
	old := []schema.Rule{
		{
			Schema:     "public",
			Table:      "orders",
			Name:       "no_update",
			Event:      "UPDATE",
			Definition: "CREATE RULE no_update AS ON UPDATE TO public.orders DO INSTEAD NOTHING",
		},
		{
			Schema:     "public",
			Table:      "orders",
			Name:       "log_insert",
			Event:      "INSERT",
			Definition: "CREATE RULE log_insert AS ON INSERT TO public.orders DO ALSO NOTIFY orders",
		},
	}

	new := []schema.Rule{
		{
			Schema:     "public",
			Table:      "orders",
			Name:       "no_update",
			Event:      "UPDATE",
			Definition: "CREATE RULE no_update AS ON UPDATE TO public.orders WHERE (status = 'closed') DO INSTEAD NOTHING",
		},
		{
			Schema:     "public",
			Table:      "orders",
			Name:       "no_delete",
			Event:      "DELETE",
			Definition: "CREATE RULE no_delete AS ON DELETE TO public.orders DO INSTEAD NOTHING",
		},
	}

	diff := schema.DiffRules(old, new)

	if len(diff.Added) != 1 {
		t.Fatalf("expected 1 added rule, got %d", len(diff.Added))
	}
	if diff.Added[0].Name != "no_delete" {
		t.Errorf("expected added rule 'no_delete', got %q", diff.Added[0].Name)
	}

	if len(diff.Removed) != 1 {
		t.Fatalf("expected 1 removed rule, got %d", len(diff.Removed))
	}
	if diff.Removed[0].Name != "log_insert" {
		t.Errorf("expected removed rule 'log_insert', got %q", diff.Removed[0].Name)
	}

	if len(diff.Changed) != 1 {
		t.Fatalf("expected 1 changed rule, got %d", len(diff.Changed))
	}
	if diff.Changed[0].Name != "no_update" {
		t.Errorf("expected changed rule 'no_update', got %q", diff.Changed[0].Name)
	}

	statements := migrate.RuleDiffSQL(diff)

	// removed + changed(drop+create) + added = 1 + 2 + 1 = 4
	if len(statements) != 4 {
		t.Fatalf("expected 4 SQL statements, got %d: %v", len(statements), statements)
	}
}
