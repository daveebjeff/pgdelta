package schema_test

import (
	"testing"

	"github.com/pgdelta/internal/schema"
)

var baseOperator = schema.Operator{
	Schema:     "public",
	Name:       "+",
	LeftType:   "integer",
	RightType:  "integer",
	ResultType: "integer",
	Procedure:  "int4pl",
}

func TestOperatorFullName(t *testing.T) {
	o := baseOperator
	if got := o.FullName(); got != "public.+(integer,integer)" {
		t.Errorf("expected public.+(integer,integer), got %s", got)
	}
}

func TestDiffOperators_NoChanges(t *testing.T) {
	old := []schema.Operator{baseOperator}
	new := []schema.Operator{baseOperator}
	diff := schema.DiffOperators(old, new)
	if !diff.IsEmpty() {
		t.Errorf("expected no diff, got %+v", diff)
	}
}

func TestDiffOperators_AddedOperator(t *testing.T) {
	old := []schema.Operator{}
	new := []schema.Operator{baseOperator}
	diff := schema.DiffOperators(old, new)
	if len(diff.Added) != 1 {
		t.Errorf("expected 1 added, got %d", len(diff.Added))
	}
}

func TestDiffOperators_RemovedOperator(t *testing.T) {
	old := []schema.Operator{baseOperator}
	new := []schema.Operator{}
	diff := schema.DiffOperators(old, new)
	if len(diff.Removed) != 1 {
		t.Errorf("expected 1 removed, got %d", len(diff.Removed))
	}
}

func TestDiffOperators_ChangedOperator(t *testing.T) {
	old := []schema.Operator{baseOperator}
	changed := baseOperator
	changed.Procedure = "custom_add"
	new := []schema.Operator{changed}
	diff := schema.DiffOperators(old, new)
	if len(diff.Changed) != 1 {
		t.Errorf("expected 1 changed, got %d", len(diff.Changed))
	}
}
