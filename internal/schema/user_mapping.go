package schema

import "fmt"

// UserMapping represents a PostgreSQL user mapping for a foreign server.
type UserMapping struct {
	User       string
	Server     string
	Options    map[string]string
}

func (u UserMapping) FullName() string {
	return fmt.Sprintf("%s@%s", u.User, u.Server)
}

func DiffUserMappings(old, new []UserMapping) SchemaDiff {
	diff := SchemaDiff{}

	oldMap := make(map[string]UserMapping, len(old))
	for _, m := range old {
		oldMap[m.FullName()] = m
	}

	newMap := make(map[string]UserMapping, len(new))
	for _, m := range new {
		newMap[m.FullName()] = m
	}

	for key, nm := range newMap {
		if _, exists := oldMap[key]; !exists {
			diff.Added = append(diff.Added, nm)
		}
	}

	for key, om := range oldMap {
		if _, exists := newMap[key]; !exists {
			diff.Removed = append(diff.Removed, om)
		} else if !userMappingsEqual(om, newMap[key]) {
			diff.Changed = append(diff.Changed, Change{Old: om, New: newMap[key]})
		}
	}

	return diff
}

func userMappingsEqual(a, b UserMapping) bool {
	if len(a.Options) != len(b.Options) {
		return false
	}
	for k, v := range a.Options {
		if b.Options[k] != v {
			return false
		}
	}
	return true
}
