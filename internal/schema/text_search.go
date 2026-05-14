package schema

import "fmt"

// TextSearchConfig represents a PostgreSQL text search configuration.
type TextSearchConfig struct {
	Schema string
	Name   string
	Parser string
}

func (t TextSearchConfig) FullName() string {
	return fmt.Sprintf("%s.%s", t.Schema, t.Name)
}

// TextSearchDiff holds added, removed, and changed text search configurations.
type TextSearchDiff struct {
	Added   []TextSearchConfig
	Removed []TextSearchConfig
	Changed []TextSearchConfig
}

func (d TextSearchDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffTextSearchConfigs computes the diff between two slices of TextSearchConfig.
func DiffTextSearchConfigs(old, new []TextSearchConfig) TextSearchDiff {
	oldMap := make(map[string]TextSearchConfig, len(old))
	for _, t := range old {
		oldMap[t.FullName()] = t
	}

	newMap := make(map[string]TextSearchConfig, len(new))
	for _, t := range new {
		newMap[t.FullName()] = t
	}

	var diff TextSearchDiff

	for _, t := range new {
		if existing, ok := oldMap[t.FullName()]; !ok {
			diff.Added = append(diff.Added, t)
		} else if !textSearchConfigsEqual(existing, t) {
			diff.Changed = append(diff.Changed, t)
		}
	}

	for _, t := range old {
		if _, ok := newMap[t.FullName()]; !ok {
			diff.Removed = append(diff.Removed, t)
		}
	}

	return diff
}

func textSearchConfigsEqual(a, b TextSearchConfig) bool {
	return a.Schema == b.Schema && a.Name == b.Name && a.Parser == b.Parser
}
