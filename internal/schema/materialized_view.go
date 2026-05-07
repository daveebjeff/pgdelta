package schema

import "fmt"

// MaterializedView represents a PostgreSQL materialized view.
type MaterializedView struct {
	Schema     string
	Name       string
	Definition string
	WithData    bool
}

// FullName returns the fully qualified name of the materialized view.
func (mv MaterializedView) FullName() string {
	return fmt.Sprintf("%s.%s", mv.Schema, mv.Name)
}

// MaterializedViewDiff holds added, removed, and changed materialized views.
type MaterializedViewDiff struct {
	Added   []MaterializedView
	Removed []MaterializedView
	Changed []MaterializedView
}

// DiffMaterializedViews compares two slices of materialized views and returns the diff.
func DiffMaterializedViews(old, new []MaterializedView) MaterializedViewDiff {
	oldMap := make(map[string]MaterializedView, len(old))
	for _, mv := range old {
		oldMap[mv.FullName()] = mv
	}

	newMap := make(map[string]MaterializedView, len(new))
	for _, mv := range new {
		newMap[mv.FullName()] = mv
	}

	var diff MaterializedViewDiff

	for _, mv := range new {
		if existing, ok := oldMap[mv.FullName()]; !ok {
			diff.Added = append(diff.Added, mv)
		} else if !materializedViewsEqual(existing, mv) {
			diff.Changed = append(diff.Changed, mv)
		}
	}

	for _, mv := range old {
		if _, ok := newMap[mv.FullName()]; !ok {
			diff.Removed = append(diff.Removed, mv)
		}
	}

	return diff
}

func materializedViewsEqual(a, b MaterializedView) bool {
	return a.Definition == b.Definition && a.WithData == b.WithData
}
