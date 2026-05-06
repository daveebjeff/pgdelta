package schema

import "fmt"

// View represents a PostgreSQL view.
type View struct {
	Schema     string
	Name       string
	Definition string
}

// FullName returns the fully qualified view name.
func (v View) FullName() string {
	return fmt.Sprintf("%s.%s", v.Schema, v.Name)
}

// ViewDiff represents a change to a view.
type ViewDiff struct {
	Added   []View
	Removed []View
	Changed []ViewChange
}

// ViewChange holds the old and new version of a changed view.
type ViewChange struct {
	Old View
	New View
}

// DiffViews compares two slices of views and returns the differences.
func DiffViews(old, new []View) ViewDiff {
	diff := ViewDiff{}

	oldMap := make(map[string]View, len(old))
	for _, v := range old {
		oldMap[v.FullName()] = v
	}

	newMap := make(map[string]View, len(new))
	for _, v := range new {
		newMap[v.FullName()] = v
	}

	for _, v := range new {
		if oldView, exists := oldMap[v.FullName()]; !exists {
			diff.Added = append(diff.Added, v)
		} else if oldView.Definition != v.Definition {
			diff.Changed = append(diff.Changed, ViewChange{Old: oldView, New: v})
		}
	}

	for _, v := range old {
		if _, exists := newMap[v.FullName()]; !exists {
			diff.Removed = append(diff.Removed, v)
		}
	}

	return diff
}
