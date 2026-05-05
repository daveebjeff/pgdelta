package schema

import "fmt"

// IndexMethod represents the access method for an index.
type IndexMethod string

const (
	IndexMethodBTree IndexMethod = "btree"
	IndexMethodHash  IndexMethod = "hash"
	IndexMethodGiST  IndexMethod = "gist"
	IndexMethodGIN   IndexMethod = "gin"
)

// Index represents a PostgreSQL index.
type Index struct {
	SchemaName string
	TableName  string
	Name       string
	Columns    []string
	Unique     bool
	Method     IndexMethod
}

// FullName returns the fully qualified index name.
func (i Index) FullName() string {
	return fmt.Sprintf("%s.%s", i.SchemaName, i.Name)
}

// IndexDiff represents the difference between two sets of indexes.
type IndexDiff struct {
	Added   []Index
	Removed []Index
	Changed []IndexChange
}

// IndexChange represents a change to an existing index.
type IndexChange struct {
	Old Index
	New Index
}

// IsEmpty returns true if there are no differences.
func (d IndexDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffIndexes compares two slices of indexes and returns the differences.
func DiffIndexes(old, new []Index) IndexDiff {
	diff := IndexDiff{}

	oldMap := make(map[string]Index, len(old))
	for _, idx := range old {
		oldMap[idx.FullName()] = idx
	}

	newMap := make(map[string]Index, len(new))
	for _, idx := range new {
		newMap[idx.FullName()] = idx
	}

	for key, newIdx := range newMap {
		if oldIdx, exists := oldMap[key]; !exists {
			diff.Added = append(diff.Added, newIdx)
		} else if !indexesEqual(oldIdx, newIdx) {
			diff.Changed = append(diff.Changed, IndexChange{Old: oldIdx, New: newIdx})
		}
	}

	for key, oldIdx := range oldMap {
		if _, exists := newMap[key]; !exists {
			diff.Removed = append(diff.Removed, oldIdx)
		}
	}

	return diff
}

func indexesEqual(a, b Index) bool {
	if a.Unique != b.Unique || a.Method != b.Method || len(a.Columns) != len(b.Columns) {
		return false
	}
	for i := range a.Columns {
		if a.Columns[i] != b.Columns[i] {
			return false
		}
	}
	return true
}
