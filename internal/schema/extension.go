package schema

// Extension represents a PostgreSQL extension (e.g. uuid-ossp, pgcrypto).
type Extension struct {
	Name    string
	Schema  string
	Version string
}

// FullName returns the extension name (extensions are not schema-qualified in the
// same way objects are, but we keep Schema for informational purposes).
func (e Extension) FullName() string {
	return e.Name
}

// ExtensionDiff holds added and removed extensions between two snapshots.
type ExtensionDiff struct {
	Added   []Extension
	Removed []Extension
	Changed []Extension // version changed
}

// DiffExtensions compares two slices of extensions and returns the diff.
func DiffExtensions(old, new []Extension) ExtensionDiff {
	oldMap := make(map[string]Extension, len(old))
	for _, e := range old {
		oldMap[e.Name] = e
	}

	newMap := make(map[string]Extension, len(new))
	for _, e := range new {
		newMap[e.Name] = e
	}

	var diff ExtensionDiff

	for _, e := range new {
		if o, exists := oldMap[e.Name]; !exists {
			diff.Added = append(diff.Added, e)
		} else if !extensionsEqual(o, e) {
			diff.Changed = append(diff.Changed, e)
		}
	}

	for _, e := range old {
		if _, exists := newMap[e.Name]; !exists {
			diff.Removed = append(diff.Removed, e)
		}
	}

	return diff
}

func extensionsEqual(a, b Extension) bool {
	return a.Name == b.Name && a.Schema == b.Schema && a.Version == b.Version
}
