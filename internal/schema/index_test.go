package schema

import (
	"testing"
)

var baseIndex = Index{
	SchemaName: "public",
	TableName:  "users",
	Name:       "idx_users_email",
	Columns:    []string{"email"},
	Unique:     true,
	Method:     IndexMethodBTree,
}

func TestIndexFullName(t *testing.T) {
	expected := "public.idx_users_email"
	if got := baseIndex.FullName(); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}
}

func TestDiffIndexes_NoChanges(t *testing.T) {
	old := []Index{baseIndex}
	new := []Index{baseIndex}
	diff := DiffIndexes(old, new)
	if !diff.IsEmpty() {
		t.Errorf("expected empty diff, got %+v", diff)
	}
}

func TestDiffIndexes_AddedIndex(t *testing.T) {
	newIdx := Index{
		SchemaName: "public",
		TableName:  "users",
		Name:       "idx_users_name",
		Columns:    []string{"name"},
		Unique:     false,
		Method:     IndexMethodBTree,
	}
	diff := DiffIndexes([]Index{baseIndex}, []Index{baseIndex, newIdx})
	if len(diff.Added) != 1 {
		t.Errorf("expected 1 added index, got %d", len(diff.Added))
	}
	if diff.Added[0].Name != newIdx.Name {
		t.Errorf("expected added index %q, got %q", newIdx.Name, diff.Added[0].Name)
	}
}

func TestDiffIndexes_RemovedIndex(t *testing.T) {
	diff := DiffIndexes([]Index{baseIndex}, []Index{})
	if len(diff.Removed) != 1 {
		t.Errorf("expected 1 removed index, got %d", len(diff.Removed))
	}
}

func TestDiffIndexes_ChangedIndex(t *testing.T) {
	modified := baseIndex
	modified.Unique = false
	diff := DiffIndexes([]Index{baseIndex}, []Index{modified})
	if len(diff.Changed) != 1 {
		t.Errorf("expected 1 changed index, got %d", len(diff.Changed))
	}
	if diff.Changed[0].Old.Unique == diff.Changed[0].New.Unique {
		t.Errorf("expected unique flag to differ")
	}
}

func TestDiffIndexes_ChangedColumns(t *testing.T) {
	modified := baseIndex
	modified.Columns = []string{"email", "created_at"}
	diff := DiffIndexes([]Index{baseIndex}, []Index{modified})
	if len(diff.Changed) != 1 {
		t.Errorf("expected 1 changed index, got %d", len(diff.Changed))
	}
}
