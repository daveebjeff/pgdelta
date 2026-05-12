package schema

// SchemaNamespace represents a PostgreSQL schema (namespace) object.
type SchemaNamespace struct {
	Name  string
	Owner string
}

// FullName returns the schema name.
func (s SchemaNamespace) FullName() string {
	return s.Name
}

// SchemaNSDiff holds added and removed schema namespaces.
type SchemaNSDiff struct {
	Added   []SchemaNamespace
	Removed []SchemaNamespace
	Changed []SchemaNamespace
}

// IsEmpty returns true if there are no namespace changes.
func (d SchemaNSDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffSchemaNamespaces computes the diff between two lists of schema namespaces.
func DiffSchemaNamespaces(old, new []SchemaNamespace) SchemaNSDiff {
	oldMap := make(map[string]SchemaNamespace, len(old))
	for _, s := range old {
		oldMap[s.Name] = s
	}
	newMap := make(map[string]SchemaNamespace, len(new))
	for _, s := range new {
		newMap[s.Name] = s
	}

	var diff SchemaNSDiff
	for _, s := range new {
		if o, exists := oldMap[s.Name]; !exists {
			diff.Added = append(diff.Added, s)
		} else if o.Owner != s.Owner {
			diff.Changed = append(diff.Changed, s)
		}
	}
	for _, s := range old {
		if _, exists := newMap[s.Name]; !exists {
			diff.Removed = append(diff.Removed, s)
		}
	}
	return diff
}
