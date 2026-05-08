package migrate_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/migrate"
	"github.com/pgdelta/pgdelta/internal/schema"
)

var baseConstraint = schema.Constraint{
	Schema:     "public",
	TableName:  "orders",
	Name:       "chk_amount_positive",
	Definition: "CHECK (amount > 0)",
}

func TestAddConstraintSQL(t *testing.T) {
	got := migrate.AddConstraintSQL(baseConstraint)
	want := "ALTER TABLE public.orders ADD CONSTRAINT chk_amount_positive CHECK (amount > 0);"
	if got != want {
		t.Errorf("AddConstraintSQL() = %q, want %q", got, want)
	}
}

func TestDropConstraintSQL(t *testing.T) {
	got := migrate.DropConstraintSQL(baseConstraint)
	want := "ALTER TABLE public.orders DROP CONSTRAINT chk_amount_positive;"
	if got != want {
		t.Errorf("DropConstraintSQL() = %q, want %q", got, want)
	}
}

func TestConstraintDiffSQL_AddedAndRemoved(t *testing.T) {
	added := schema.Constraint{
		Schema:     "public",
		TableName:  "orders",
		Name:       "chk_qty_positive",
		Definition: "CHECK (qty > 0)",
	}

	stmts := migrate.ConstraintDiffSQL(
		[]schema.Constraint{baseConstraint},
		[]schema.Constraint{added},
	)

	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(stmts))
	}
	if stmts[0] != "ALTER TABLE public.orders DROP CONSTRAINT chk_amount_positive;" {
		t.Errorf("unexpected drop statement: %q", stmts[0])
	}
	if stmts[1] != "ALTER TABLE public.orders ADD CONSTRAINT chk_qty_positive CHECK (qty > 0);" {
		t.Errorf("unexpected add statement: %q", stmts[1])
	}
}

func TestConstraintDiffSQL_Changed(t *testing.T) {
	changed := schema.Constraint{
		Schema:     "public",
		TableName:  "orders",
		Name:       "chk_amount_positive",
		Definition: "CHECK (amount >= 0)",
	}

	stmts := migrate.ConstraintDiffSQL(
		[]schema.Constraint{baseConstraint},
		[]schema.Constraint{changed},
	)

	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements for changed constraint, got %d", len(stmts))
	}
	if stmts[0] != "ALTER TABLE public.orders DROP CONSTRAINT chk_amount_positive;" {
		t.Errorf("unexpected drop statement: %q", stmts[0])
	}
	if stmts[1] != "ALTER TABLE public.orders ADD CONSTRAINT chk_amount_positive CHECK (amount >= 0);" {
		t.Errorf("unexpected add statement: %q", stmts[1])
	}
}

func TestConstraintDiffSQL_NoChanges(t *testing.T) {
	stmts := migrate.ConstraintDiffSQL(
		[]schema.Constraint{baseConstraint},
		[]schema.Constraint{baseConstraint},
	)
	if len(stmts) != 0 {
		t.Errorf("expected no statements, got %d", len(stmts))
	}
}
