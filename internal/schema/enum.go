package schema

import "fmt"

// Enum represents a PostgreSQL enum type.
type Enum struct {
	Schema string
	Name   string
	Values []string
}

// FullName returns the fully qualified enum name.
func (e Enum) FullName() string {
	return fmt.Sprintf("%s.%s", e.Schema, e.Name)
}

// EnumDiff holds added and removed enum types between two schema snapshots.
type EnumDiff struct {
	Added   []Enum
	Removed []Enum
	Changed []EnumChange
}

// EnumChange represents a change in enum values for an existing enum.
type EnumChange struct {
	Old Enum
	New Enum
}

// DiffEnums computes the difference between two slices of enums.
func DiffEnums(old, new []Enum) EnumDiff {
	var diff EnumDiff

	oldMap := make(map[string]Enum, len(old))
	for _, e := range old {
		oldMap[e.FullName()] = e
	}

	newMap := make(map[string]Enum, len(new))
	for _, e := range new {
		newMap[e.FullName()] = e
	}

	for _, e := range new {
		if oldEnum, exists := oldMap[e.FullName()]; !exists {
			diff.Added = append(diff.Added, e)
		} else if !enumsEqual(oldEnum, e) {
			diff.Changed = append(diff.Changed, EnumChange{Old: oldEnum, New: e})
		}
	}

	for _, e := range old {
		if _, exists := newMap[e.FullName()]; !exists {
			diff.Removed = append(diff.Removed, e)
		}
	}

	return diff
}

func enumsEqual(a, b Enum) bool {
	if len(a.Values) != len(b.Values) {
		return false
	}
	for i := range a.Values {
		if a.Values[i] != b.Values[i] {
			return false
		}
	}
	return true
}
