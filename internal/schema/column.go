package schema

import "fmt"

// Column represents a PostgreSQL table column.
type Column struct {
	Name       string
	DataType   string
	Nullable   bool
	Default    *string
	Position   int
}

// FullName returns a qualified column identifier.
func (c Column) FullName(tableName string) string {
	return fmt.Sprintf("%s.%s", tableName, c.Name)
}

// ColumnChange describes a change to a column between two schema snapshots.
type ColumnChange struct {
	ColumnName string
	ChangeType string // "type_changed", "nullable_changed", "default_changed"
	OldValue   string
	NewValue   string
}

// DiffColumns compares two slices of columns and returns added, removed, and modified columns.
func DiffColumns(old, new []Column) (added []Column, removed []Column, changed []ColumnChange) {
	oldMap := make(map[string]Column, len(old))
	for _, c := range old {
		oldMap[c.Name] = c
	}

	newMap := make(map[string]Column, len(new))
	for _, c := range new {
		newMap[c.Name] = c
	}

	for _, c := range new {
		if _, exists := oldMap[c.Name]; !exists {
			added = append(added, c)
		}
	}

	for _, c := range old {
		if _, exists := newMap[c.Name]; !exists {
			removed = append(removed, c)
		}
	}

	for _, newCol := range new {
		oldCol, exists := oldMap[newCol.Name]
		if !exists {
			continue
		}
		if oldCol.DataType != newCol.DataType {
			changed = append(changed, ColumnChange{
				ColumnName: newCol.Name,
				ChangeType: "type_changed",
				OldValue:   oldCol.DataType,
				NewValue:   newCol.DataType,
			})
		}
		if oldCol.Nullable != newCol.Nullable {
			changed = append(changed, ColumnChange{
				ColumnName: newCol.Name,
				ChangeType: "nullable_changed",
				OldValue:   fmt.Sprintf("%v", oldCol.Nullable),
				NewValue:   fmt.Sprintf("%v", newCol.Nullable),
			})
		}
	}

	return added, removed, changed
}
