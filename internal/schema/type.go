package schema

import "fmt"

// Type represents a PostgreSQL composite or domain type.
type Type struct {
	Schema     string
	Name       string
	Kind       string // 'composite' or 'domain'
	Definition string // e.g. column list for composite, base type for domain
}

func (t Type) FullName() string {
	return fmt.Sprintf("%s.%s", t.Schema, t.Name)
}

// TypeDiff holds added, removed, and changed types.
type TypeDiff struct {
	Added   []Type
	Removed []Type
	Changed []Type
}

func (d TypeDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffTypes computes the difference between two slices of Type.
func DiffTypes(old, new []Type) TypeDiff {
	oldMap := make(map[string]Type, len(old))
	for _, t := range old {
		oldMap[t.FullName()] = t
	}

	newMap := make(map[string]Type, len(new))
	for _, t := range new {
		newMap[t.FullName()] = t
	}

	var diff TypeDiff

	for _, t := range new {
		if o, exists := oldMap[t.FullName()]; !exists {
			diff.Added = append(diff.Added, t)
		} else if !typesEqual(o, t) {
			diff.Changed = append(diff.Changed, t)
		}
	}

	for _, t := range old {
		if _, exists := newMap[t.FullName()]; !exists {
			diff.Removed = append(diff.Removed, t)
		}
	}

	return diff
}

func typesEqual(a, b Type) bool {
	return a.Kind == b.Kind && a.Definition == b.Definition
}
