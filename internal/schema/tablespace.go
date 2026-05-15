package schema

// Tablespace represents a PostgreSQL tablespace object.
type Tablespace struct {
	Name     string
	Owner    string
	Location string
}

// FullName returns the identifier for the tablespace.
func (t Tablespace) FullName() string {
	return t.Name
}

// TablespaceDiff holds added, removed, and changed tablespaces.
type TablespaceDiff struct {
	Added   []Tablespace
	Removed []Tablespace
	Changed []Tablespace
}

// IsEmpty returns true when there are no tablespace differences.
func (d TablespaceDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffTablespaces computes the difference between two sets of tablespaces.
func DiffTablespaces(old, new []Tablespace) TablespaceDiff {
	oldMap := make(map[string]Tablespace, len(old))
	for _, ts := range old {
		oldMap[ts.Name] = ts
	}

	newMap := make(map[string]Tablespace, len(new))
	for _, ts := range new {
		newMap[ts.Name] = ts
	}

	var diff TablespaceDiff

	for _, ts := range new {
		if existing, ok := oldMap[ts.Name]; !ok {
			diff.Added = append(diff.Added, ts)
		} else if !tablespacesEqual(existing, ts) {
			diff.Changed = append(diff.Changed, ts)
		}
	}

	for _, ts := range old {
		if _, ok := newMap[ts.Name]; !ok {
			diff.Removed = append(diff.Removed, ts)
		}
	}

	return diff
}

func tablespacesEqual(a, b Tablespace) bool {
	return a.Owner == b.Owner && a.Location == b.Location
}
