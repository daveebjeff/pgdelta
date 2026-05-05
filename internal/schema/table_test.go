package schema

import (
	"testing"
)

func baseTable() *Table {
	return &Table{
		Schema: "public",
		Name:   "users",
		Columns: []Column{
			{Name: "id", DataType: "integer", Nullable: false, Ordinal: 1},
			{Name: "email", DataType: "text", Nullable: false, Ordinal: 2},
			{Name: "created_at", DataType: "timestamptz", Nullable: true, Ordinal: 3},
		},
	}
}

func TestFullName(t *testing.T) {
	tbl := baseTable()
	if got := tbl.FullName(); got != "public.users" {
		t.Errorf("expected public.users, got %s", got)
	}
}

func TestDiffTables_NoChanges(t *testing.T) {
	before := baseTable()
	after := baseTable()
	diff := DiffTables(before, after)
	if !diff.IsEmpty() {
		t.Errorf("expected empty diff, got %+v", diff)
	}
}

func TestDiffTables_AddedColumn(t *testing.T) {
	before := baseTable()
	after := baseTable()
	after.Columns = append(after.Columns, Column{Name: "bio", DataType: "text", Nullable: true, Ordinal: 4})

	diff := DiffTables(before, after)
	if len(diff.Added) != 1 || diff.Added[0].Name != "bio" {
		t.Errorf("expected 1 added column 'bio', got %+v", diff.Added)
	}
	if len(diff.Removed) != 0 || len(diff.Modified) != 0 {
		t.Errorf("unexpected removed or modified columns")
	}
}

func TestDiffTables_RemovedColumn(t *testing.T) {
	before := baseTable()
	after := baseTable()
	after.Columns = after.Columns[:2] // remove created_at

	diff := DiffTables(before, after)
	if len(diff.Removed) != 1 || diff.Removed[0].Name != "created_at" {
		t.Errorf("expected 1 removed column 'created_at', got %+v", diff.Removed)
	}
}

func TestDiffTables_ModifiedColumn(t *testing.T) {
	before := baseTable()
	after := baseTable()
	after.Columns[1].DataType = "varchar(255)" // change email type

	diff := DiffTables(before, after)
	if len(diff.Modified) != 1 || diff.Modified[0].Name != "email" {
		t.Errorf("expected 1 modified column 'email', got %+v", diff.Modified)
	}
	if diff.Modified[0].Before.DataType != "text" || diff.Modified[0].After.DataType != "varchar(255)" {
		t.Errorf("unexpected before/after types: %+v", diff.Modified[0])
	}
}
