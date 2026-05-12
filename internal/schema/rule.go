package schema

import "fmt"

// Rule represents a PostgreSQL rewrite rule on a table.
type Rule struct {
	Schema    string
	Table     string
	Name      string
	Event     string // SELECT, INSERT, UPDATE, DELETE
	Condition string
	Instead   bool
	Definition string
}

func (r Rule) FullName() string {
	return fmt.Sprintf("%s.%s.%s", r.Schema, r.Table, r.Name)
}

// RuleDiff holds added, removed, and changed rules.
type RuleDiff struct {
	Added   []Rule
	Removed []Rule
	Changed []Rule
}

func (d RuleDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffRules computes the difference between two slices of rules.
func DiffRules(old, new []Rule) RuleDiff {
	diff := RuleDiff{}

	oldMap := make(map[string]Rule, len(old))
	for _, r := range old {
		oldMap[r.FullName()] = r
	}

	newMap := make(map[string]Rule, len(new))
	for _, r := range new {
		newMap[r.FullName()] = r
	}

	for _, r := range new {
		if existing, ok := oldMap[r.FullName()]; !ok {
			diff.Added = append(diff.Added, r)
		} else if !rulesEqual(existing, r) {
			diff.Changed = append(diff.Changed, r)
		}
	}

	for _, r := range old {
		if _, ok := newMap[r.FullName()]; !ok {
			diff.Removed = append(diff.Removed, r)
		}
	}

	return diff
}

func rulesEqual(a, b Rule) bool {
	return a.Event == b.Event &&
		a.Condition == b.Condition &&
		a.Instead == b.Instead &&
		a.Definition == b.Definition
}
