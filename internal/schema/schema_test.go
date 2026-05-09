package schema_test

import (
	"testing"

	"github.com/pgdelta/internal/schema"
)

func TestDiffSchemas_Empty(t *testing.T) {
	old := schema.Schema{}
	new := schema.Schema{}
	diff := schema.DiffSchemas(old, new)
	if !diff.IsEmpty() {
		t.Error("expected empty diff for identical empty schemas")
	}
}

func TestDiffSchemas_AddedTable(t *testing.T) {
	old := schema.Schema{}
	new := schema.Schema{
		Tables: []schema.Table{
			{Schema: "public", Name: "users"},
		},
	}
	diff := schema.DiffSchemas(old, new)
	if diff.IsEmpty() {
		t.Error("expected non-empty diff when table is added")
	}
	if len(diff.TableDiff.Added) != 1 {
		t.Errorf("expected 1 added table, got %d", len(diff.TableDiff.Added))
	}
}

func TestDiffSchemas_RemovedEnum(t *testing.T) {
	old := schema.Schema{
		Enums: []schema.Enum{
			{Schema: "public", Name: "status", Values: []string{"active", "inactive"}},
		},
	}
	new := schema.Schema{}
	diff := schema.DiffSchemas(old, new)
	if diff.IsEmpty() {
		t.Error("expected non-empty diff when enum is removed")
	}
	if len(diff.EnumDiff.Removed) != 1 {
		t.Errorf("expected 1 removed enum, got %d", len(diff.EnumDiff.Removed))
	}
}

func TestDiffSchemas_MultipleChanges(t *testing.T) {
	old := schema.Schema{
		Sequences: []schema.Sequence{
			{Schema: "public", Name: "old_seq", IncrementBy: 1, MinValue: 1, MaxValue: 9999, StartValue: 1, Cache: 1},
		},
	}
	new := schema.Schema{
		Tables: []schema.Table{
			{Schema: "public", Name: "orders"},
		},
	}
	diff := schema.DiffSchemas(old, new)
	if diff.IsEmpty() {
		t.Error("expected non-empty diff")
	}
	if len(diff.TableDiff.Added) != 1 {
		t.Errorf("expected 1 added table, got %d", len(diff.TableDiff.Added))
	}
	if len(diff.SequenceDiff.Removed) != 1 {
		t.Errorf("expected 1 removed sequence, got %d", len(diff.SequenceDiff.Removed))
	}
}

func TestDiffSchemas_IsEmpty_NoChanges(t *testing.T) {
	base := schema.Schema{
		Tables: []schema.Table{{Schema: "public", Name: "users"}},
		Extensions: []schema.Extension{{Schema: "public", Name: "uuid-ossp", Version: "1.1"}},
	}
	diff := schema.DiffSchemas(base, base)
	if !diff.IsEmpty() {
		t.Error("expected empty diff for identical schemas")
	}
}
