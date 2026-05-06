package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

var baseTrigger = schema.Trigger{
	Schema:    "public",
	Name:      "trg_audit",
	Table:     "orders",
	Timing:    "BEFORE",
	Event:     "INSERT OR UPDATE",
	ForEach:   "ROW",
	Function:  "audit_trigger_fn",
	Condition: "",
	Definition: "CREATE TRIGGER trg_audit BEFORE INSERT OR UPDATE ON public.orders FOR EACH ROW EXECUTE FUNCTION audit_trigger_fn()",
}

func TestCreateTriggerSQL(t *testing.T) {
	sql := migrate.CreateTriggerSQL(baseTrigger)
	expected := baseTrigger.Definition
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestDropTriggerSQL(t *testing.T) {
	sql := migrate.DropTriggerSQL(baseTrigger)
	expected := "DROP TRIGGER trg_audit ON public.orders;"
	if sql != expected {
		t.Errorf("expected: %q, got: %q", expected, sql)
	}
}

func TestTriggerDiffSQL_AddedAndRemoved(t *testing.T) {
	added := baseTrigger
	removed := schema.Trigger{
		Schema:     "public",
		Name:       "trg_old",
		Table:      "orders",
		Definition: "CREATE TRIGGER trg_old AFTER DELETE ON public.orders FOR EACH ROW EXECUTE FUNCTION old_fn()",
	}

	diff := schema.TriggerDiff{
		Added:   []schema.Trigger{added},
		Removed: []schema.Trigger{removed},
		Changed: nil,
	}

	sqls := migrate.TriggerDiffSQL(diff)
	if len(sqls) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(sqls))
	}
	if sqls[0] != migrate.DropTriggerSQL(removed) {
		t.Errorf("expected drop first, got: %s", sqls[0])
	}
	if sqls[1] != migrate.CreateTriggerSQL(added) {
		t.Errorf("expected create second, got: %s", sqls[1])
	}
}

func TestTriggerDiffSQL_Changed(t *testing.T) {
	old := baseTrigger
	new := baseTrigger
	new.Definition = "CREATE TRIGGER trg_audit BEFORE INSERT ON public.orders FOR EACH ROW EXECUTE FUNCTION audit_trigger_fn()"

	diff := schema.TriggerDiff{
		Changed: []schema.TriggerChange{{Old: old, New: new}},
	}

	sqls := migrate.TriggerDiffSQL(diff)
	if len(sqls) != 2 {
		t.Fatalf("expected 2 statements for changed trigger, got %d", len(sqls))
	}
	if sqls[0] != migrate.DropTriggerSQL(old) {
		t.Errorf("expected drop old trigger, got: %s", sqls[0])
	}
	if sqls[1] != migrate.CreateTriggerSQL(new) {
		t.Errorf("expected create new trigger, got: %s", sqls[1])
	}
}
