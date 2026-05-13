package schema

import "fmt"

// Collation represents a PostgreSQL collation object.
type Collation struct {
	Schema   string
	Name     string
	Provider string // icu, libc, pg_default
	Locale   string
	Deterministic bool
}

// FullName returns the fully qualified collation name.
func (c Collation) FullName() string {
	return fmt.Sprintf("%s.%s", c.Schema, c.Name)
}

// CollationDiff holds added, removed, and changed collations.
type CollationDiff struct {
	Added   []Collation
	Removed []Collation
	Changed []Collation
}

// IsEmpty returns true if there are no collation changes.
func (d CollationDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffCollations computes the diff between two slices of Collation.
func DiffCollations(old, new []Collation) CollationDiff {
	diff := CollationDiff{}

	oldMap := make(map[string]Collation, len(old))
	for _, c := range old {
		oldMap[c.FullName()] = c
	}

	newMap := make(map[string]Collation, len(new))
	for _, c := range new {
		newMap[c.FullName()] = c
	}

	for _, c := range new {
		if o, exists := oldMap[c.FullName()]; !exists {
			diff.Added = append(diff.Added, c)
		} else if !collationsEqual(o, c) {
			diff.Changed = append(diff.Changed, c)
		}
	}

	for _, c := range old {
		if _, exists := newMap[c.FullName()]; !exists {
			diff.Removed = append(diff.Removed, c)
		}
	}

	return diff
}

func collationsEqual(a, b Collation) bool {
	return a.Provider == b.Provider &&
		a.Locale == b.Locale &&
		a.Deterministic == b.Deterministic
}
