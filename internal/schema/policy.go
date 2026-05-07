package schema

import "fmt"

// Policy represents a PostgreSQL row-level security policy.
type Policy struct {
	Schema     string
	Table      string
	Name       string
	Command    string // ALL, SELECT, INSERT, UPDATE, DELETE
	Permissive bool   // true = PERMISSIVE, false = RESTRICTIVE
	Roles      []string
	Using      string
	WithCheck  string
}

func (p Policy) FullName() string {
	return fmt.Sprintf("%s.%s.%s", p.Schema, p.Table, p.Name)
}

// PolicyDiff holds added, removed, and changed policies.
type PolicyDiff struct {
	Added   []Policy
	Removed []Policy
	Changed []Policy
}

// DiffPolicies computes the diff between two slices of policies.
func DiffPolicies(old, new []Policy) PolicyDiff {
	oldMap := make(map[string]Policy, len(old))
	for _, p := range old {
		oldMap[p.FullName()] = p
	}

	newMap := make(map[string]Policy, len(new))
	for _, p := range new {
		newMap[p.FullName()] = p
	}

	var diff PolicyDiff

	for _, p := range new {
		if old, exists := oldMap[p.FullName()]; !exists {
			diff.Added = append(diff.Added, p)
		} else if !policiesEqual(old, p) {
			diff.Changed = append(diff.Changed, p)
		}
	}

	for _, p := range old {
		if _, exists := newMap[p.FullName()]; !exists {
			diff.Removed = append(diff.Removed, p)
		}
	}

	return diff
}

func policiesEqual(a, b Policy) bool {
	if a.Command != b.Command ||
		a.Permissive != b.Permissive ||
		a.Using != b.Using ||
		a.WithCheck != b.WithCheck {
		return false
	}
	if len(a.Roles) != len(b.Roles) {
		return false
	}
	for i := range a.Roles {
		if a.Roles[i] != b.Roles[i] {
			return false
		}
	}
	return true
}
