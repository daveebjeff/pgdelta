package schema

import "fmt"

// Constraint represents a table constraint in a PostgreSQL schema.
type Constraint struct {
	Schema     string
	Table      string
	Name       string
	Type       string // CHECK, UNIQUE, PRIMARY KEY, FOREIGN KEY
	Definition string
}

// FullName returns the fully qualified constraint name.
func (c Constraint) FullName() string {
	return fmt.Sprintf("%s.%s.%s", c.Schema, c.Table, c.Name)
}

// ConstraintDiff holds added, removed, and changed constraints.
type ConstraintDiff struct {
	Added   []Constraint
	Removed []Constraint
	Changed []Constraint
}

// DiffConstraints compares two slices of constraints and returns a ConstraintDiff.
func DiffConstraints(old, new []Constraint) ConstraintDiff {
	oldMap := make(map[string]Constraint, len(old))
	for _, c := range old {
		oldMap[c.FullName()] = c
	}

	newMap := make(map[string]Constraint, len(new))
	for _, c := range new {
		newMap[c.FullName()] = c
	}

	var diff ConstraintDiff

	for key, nc := range newMap {
		if oc, exists := oldMap[key]; !exists {
			diff.Added = append(diff.Added, nc)
		} else if !constraintsEqual(oc, nc) {
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

func constraintsEqual(a, b Constraint) bool {
	return a.Type == b.Type && a.Definition == b.Definition
}
