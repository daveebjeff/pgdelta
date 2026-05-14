package schema

import "fmt"

// Cast represents a PostgreSQL cast definition.
type Cast struct {
	SourceType string
	TargetType string
	FunctionName string
	Schema       string
	CastContext  string // 'e' = explicit, 'a' = assignment, 'i' = implicit
}

// FullName returns a human-readable identifier for the cast.
func (c Cast) FullName() string {
	return fmt.Sprintf("(%s AS %s)", c.SourceType, c.TargetType)
}

func castsEqual(a, b Cast) bool {
	return a.FunctionName == b.FunctionName &&
		a.Schema == b.Schema &&
		a.CastContext == b.CastContext
}

// CastDiff represents added, removed, or changed casts.
type CastDiff struct {
	Added   []Cast
	Removed []Cast
	Changed []Cast
}

// IsEmpty returns true if there are no cast changes.
func (d CastDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffCasts compares two slices of casts and returns the diff.
func DiffCasts(old, new []Cast) CastDiff {
	oldMap := make(map[string]Cast, len(old))
	for _, c := range old {
		oldMap[c.FullName()] = c
	}

	newMap := make(map[string]Cast, len(new))
	for _, c := range new {
		newMap[c.FullName()] = c
	}

	var diff CastDiff

	for key, nc := range newMap {
		if oc, exists := oldMap[key]; !exists {
			diff.Added = append(diff.Added, nc)
		} else if !castsEqual(oc, nc) {
			diff.Changed = append(diff.Changed, nc)
		}
	}

	for key, oc := range oldMap {
		if _, exists := newMap[key]; !exists {
			diff.Removed = append(diff.Removed, oc)
		}
	}

	return diff
}
