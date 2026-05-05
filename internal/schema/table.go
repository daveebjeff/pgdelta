package schema

import "fmt"

// Column represents a single column in a PostgreSQL table.
type Column struct {
	Name       string
	DataType   string
	Nullable   bool
	Default    string
	Ordinal    int
}

// Table represents a PostgreSQL table within a schema snapshot.
type Table struct {
	Schema  string
	Name    string
	Columns []Column
}

// FullName returns the fully qualified table name (schema.table).
func (t *Table) FullName() string {
	return fmt.Sprintf("%s.%s", t.Schema, t.Name)
}

// ColumnMap returns a map of column name to Column for fast lookup.
func (t *Table) ColumnMap() map[string]Column {
	m := make(map[string]Column, len(t.Columns))
	for _, c := range t.Columns {
		m[c.Name] = c
	}
	return m
}

// TableDiff holds the differences between two table snapshots.
type TableDiff struct {
	Table      string
	Added      []Column
	Removed    []Column
	Modified   []ColumnChange
}

// ColumnChange captures a before/after change for a single column.
type ColumnChange struct {
	Name   string
	Before Column
	After  Column
}

// IsEmpty returns true when there are no differences.
func (d *TableDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Modified) == 0
}

// DiffTables compares two Table snapshots and returns a TableDiff.
func DiffTables(before, after *Table) TableDiff {
	diff := TableDiff{Table: after.FullName()}

	beforeMap := before.ColumnMap()
	afterMap := after.ColumnMap()

	for name, afterCol := range afterMap {
		beforeCol, exists := beforeMap[name]
		if !exists {
			diff.Added = append(diff.Added, afterCol)
			continue
		}
		if beforeCol.DataType != afterCol.DataType ||
			beforeCol.Nullable != afterCol.Nullable ||
			beforeCol.Default != afterCol.Default {
			diff.Modified = append(diff.Modified, ColumnChange{
				Name:   name,
				Before: beforeCol,
				After:  afterCol,
			})
		}
	}

	for name, beforeCol := range beforeMap {
		if _, exists := afterMap[name]; !exists {
			diff.Removed = append(diff.Removed, beforeCol)
		}
	}

	return diff
}
