package schema_test

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
)

var baseEnum = schema.Enum{
	Schema: "public",
	Name:   "mood",
	Values: []string{"happy", "sad", "neutral"},
}

func TestEnumFullName(t *testing.T) {
	if got := baseEnum.FullName(); got != "public.mood" {
		t.Errorf("expected public.mood, got %s", got)
	}
}

func TestDiffEnums_NoChanges(t *testing.T) {
	old := []schema.Enum{baseEnum}
	new := []schema.Enum{baseEnum}
	diff := schema.DiffEnums(old, new)
	if len(diff.Added) != 0 || len(diff.Removed) != 0 || len(diff.Changed) != 0 {
		t.Errorf("expected no diff, got %+v", diff)
	}
}

func TestDiffEnums_AddedEnum(t *testing.T) {
	newEnum := schema.Enum{Schema: "public", Name: "status", Values: []string{"active", "inactive"}}
	diff := schema.DiffEnums([]schema.Enum{}, []schema.Enum{newEnum})
	if len(diff.Added) != 1 || diff.Added[0].FullName() != "public.status" {
		t.Errorf("expected one added enum, got %+v", diff)
	}
}

func TestDiffEnums_RemovedEnum(t *testing.T) {
	diff := schema.DiffEnums([]schema.Enum{baseEnum}, []schema.Enum{})
	if len(diff.Removed) != 1 || diff.Removed[0].FullName() != "public.mood" {
		t.Errorf("expected one removed enum, got %+v", diff)
	}
}

func TestDiffEnums_ChangedEnum(t *testing.T) {
	updated := schema.Enum{
		Schema: "public",
		Name:   "mood",
		Values: []string{"happy", "sad", "neutral", "excited"},
	}
	diff := schema.DiffEnums([]schema.Enum{baseEnum}, []schema.Enum{updated})
	if len(diff.Changed) != 1 {
		t.Fatalf("expected one changed enum, got %+v", diff)
	}
	if diff.Changed[0].New.FullName() != "public.mood" {
		t.Errorf("unexpected changed enum: %+v", diff.Changed[0])
	}
	if len(diff.Changed[0].New.Values) != 4 {
		t.Errorf("expected 4 values in new enum, got %d", len(diff.Changed[0].New.Values))
	}
}
