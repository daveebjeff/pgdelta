package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func TestCreateRuleSQL(t *testing.T) {
	rule := schema.Rule{
		Schema:    "public",
		Table:     "orders",
		Name:      "no_delete",
		Event:     "DELETE",
		Condition: "",
		Definition: "CREATE RULE no_delete AS ON DELETE TO public.orders DO INSTEAD NOTHING",
	}

	got := migrate.CreateRuleSQL(rule)
	want := "CREATE RULE no_delete AS ON DELETE TO public.orders DO INSTEAD NOTHING;"

	if got != want {
		t.Errorf("CreateRuleSQL() = %q, want %q", got, want)
	}
}

func TestDropRuleSQL(t *testing.T) {
	rule := schema.Rule{
		Schema: "public",
		Table:  "orders",
		Name:   "no_delete",
	}

	got := migrate.DropRuleSQL(rule)
	want := "DROP RULE no_delete ON public.orders;"

	if got != want {
		t.Errorf("DropRuleSQL() = %q, want %q", got, want)
	}
}

func TestRuleDiffSQL_AddedAndRemoved(t *testing.T) {
	added := []schema.Rule{
		{
			Schema:     "public",
			Table:      "orders",
			Name:       "no_delete",
			Event:      "DELETE",
			Definition: "CREATE RULE no_delete AS ON DELETE TO public.orders DO INSTEAD NOTHING",
		},
	}
	removed := []schema.Rule{
		{
			Schema: "public",
			Table:  "orders",
			Name:   "old_rule",
		},
	}

	diff := schema.RuleDiff{Added: added, Removed: removed}
	statements := migrate.RuleDiffSQL(diff)

	if len(statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(statements))
	}
	if statements[0] != "DROP RULE old_rule ON public.orders;" {
		t.Errorf("unexpected drop statement: %q", statements[0])
	}
	if statements[1] != "CREATE RULE no_delete AS ON DELETE TO public.orders DO INSTEAD NOTHING;" {
		t.Errorf("unexpected create statement: %q", statements[1])
	}
}

func TestRuleDiffSQL_Changed(t *testing.T) {
	changed := []schema.Rule{
		{
			Schema:     "public",
			Table:      "products",
			Name:       "redirect_insert",
			Event:      "INSERT",
			Definition: "CREATE RULE redirect_insert AS ON INSERT TO public.products DO INSTEAD NOTHING",
		},
	}

	diff := schema.RuleDiff{Changed: changed}
	statements := migrate.RuleDiffSQL(diff)

	if len(statements) != 2 {
		t.Fatalf("expected 2 statements for changed rule, got %d", len(statements))
	}
	if statements[0] != "DROP RULE redirect_insert ON public.products;" {
		t.Errorf("unexpected drop statement: %q", statements[0])
	}
	if statements[1] != "CREATE RULE redirect_insert AS ON INSERT TO public.products DO INSTEAD NOTHING;" {
		t.Errorf("unexpected create statement: %q", statements[1])
	}
}
