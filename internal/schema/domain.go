package schema

import "fmt"

// Domain represents a PostgreSQL domain type.
type Domain struct {
	Schema      string
	Name        string
	BaseType    string
	NotNull     bool
	Default     *string
	CheckName   *string
	CheckClause *string
}

// FullName returns the fully qualified domain name.
func (d Domain) FullName() string {
	return fmt.Sprintf("%s.%s", d.Schema, d.Name)
}

// DomainDiff holds added, removed, and changed domains.
type DomainDiff struct {
	Added   []Domain
	Removed []Domain
	Changed []Domain
}

// IsEmpty returns true if there are no domain changes.
func (d DomainDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffDomains computes the difference between two slices of domains.
func DiffDomains(old, new []Domain) DomainDiff {
	oldMap := make(map[string]Domain, len(old))
	for _, d := range old {
		oldMap[d.FullName()] = d
	}

	newMap := make(map[string]Domain, len(new))
	for _, d := range new {
		newMap[d.FullName()] = d
	}

	var diff DomainDiff

	for _, d := range new {
		if o, exists := oldMap[d.FullName()]; !exists {
			diff.Added = append(diff.Added, d)
		} else if !domainsEqual(o, d) {
			diff.Changed = append(diff.Changed, d)
		}
	}

	for _, d := range old {
		if _, exists := newMap[d.FullName()]; !exists {
			diff.Removed = append(diff.Removed, d)
		}
	}

	return diff
}

func domainsEqual(a, b Domain) bool {
	if a.BaseType != b.BaseType || a.NotNull != b.NotNull {
		return false
	}
	if (a.Default == nil) != (b.Default == nil) {
		return false
	}
	if a.Default != nil && *a.Default != *b.Default {
		return false
	}
	if (a.CheckClause == nil) != (b.CheckClause == nil) {
		return false
	}
	if a.CheckClause != nil && *a.CheckClause != *b.CheckClause {
		return false
	}
	return true
}
