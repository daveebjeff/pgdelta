package schema

import (
	"testing"
)

func baseColumn(name, dataType string, nullable bool) Column {
	return Column{
		Name:     name,
		DataType: dataType,
		Nullable: nullable,
		Position: 1,
	}
}

func TestColumnFullName(t *testing.T) {
	c := baseColumn("email", "text", true)
	got := c.FullName("users")
	if got != "users.email" {
		t.Errorf("expected users.email, got %s", got)
	}
}

func TestDiffColumns_NoChanges(t *testing.T) {
	cols := []Column{baseColumn("id", "int", false), baseColumn("name", "text", true)}
	added, removed, changed := DiffColumns(cols, cols)
	if len(added) != 0 || len(removed) != 0 || len(changed) != 0 {
		t.Errorf("expected no changes, got added=%d removed=%d changed=%d", len(added), len(removed), len(changed))
	}
}

func TestDiffColumns_AddedColumn(t *testing.T) {
	old := []Column{baseColumn("id", "int", false)}
	new := []Column{baseColumn("id", "int", false), baseColumn("email", "text", true)}
	added, removed, changed := DiffColumns(old, new)
	if len(added) != 1 || added[0].Name != "email" {
		t.Errorf("expected 1 added column 'email', got %v", added)
	}
	if len(removed) != 0 || len(changed) != 0 {
		t.Errorf("unexpected removed or changed columns")
	}
}

func TestDiffColumns_RemovedColumn(t *testing.T) {
	old := []Column{baseColumn("id", "int", false), baseColumn("email", "text", true)}
	new := []Column{baseColumn("id", "int", false)}
	added, removed, changed := DiffColumns(old, new)
	if len(removed) != 1 || removed[0].Name != "email" {
		t.Errorf("expected 1 removed column 'email', got %v", removed)
	}
	if len(added) != 0 || len(changed) != 0 {
		t.Errorf("unexpected added or changed columns")
	}
}

func TestDiffColumns_TypeChanged(t *testing.T) {
	old := []Column{baseColumn("age", "int", false)}
	new := []Column{baseColumn("age", "bigint", false)}
	_, _, changed := DiffColumns(old, new)
	if len(changed) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changed))
	}
	if changed[0].ChangeType != "type_changed" || changed[0].OldValue != "int" || changed[0].NewValue != "bigint" {
		t.Errorf("unexpected change: %+v", changed[0])
	}
}

func TestDiffColumns_NullableChanged(t *testing.T) {
	old := []Column{baseColumn("name", "text", true)}
	new := []Column{baseColumn("name", "text", false)}
	_, _, changed := DiffColumns(old, new)
	if len(changed) != 1 || changed[0].ChangeType != "nullable_changed" {
		t.Errorf("expected nullable_changed, got %+v", changed)
	}
}
