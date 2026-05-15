package schema

// AccessMethod represents a PostgreSQL access method (e.g. btree, hash, gist).
type AccessMethod struct {
	Name    string
	Type    string // 'index' or 'table'
	Handler string
}

// FullName returns the unique identifier for the access method.
func (a AccessMethod) FullName() string {
	return a.Name
}

// AccessMethodDiff holds added and removed access methods between two schema snapshots.
type AccessMethodDiff struct {
	Added   []AccessMethod
	Removed []AccessMethod
	Changed []AccessMethod
}

// IsEmpty returns true when there are no access method changes.
func (d AccessMethodDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffAccessMethods computes the difference between two slices of AccessMethod.
func DiffAccessMethods(old, new []AccessMethod) AccessMethodDiff {
	oldMap := make(map[string]AccessMethod, len(old))
	for _, a := range old {
		oldMap[a.FullName()] = a
	}

	newMap := make(map[string]AccessMethod, len(new))
	for _, a := range new {
		newMap[a.FullName()] = a
	}

	var diff AccessMethodDiff

	for _, a := range new {
		if existing, ok := oldMap[a.FullName()]; !ok {
			diff.Added = append(diff.Added, a)
		} else if !accessMethodsEqual(existing, a) {
			diff.Changed = append(diff.Changed, a)
		}
	}

	for _, a := range old {
		if _, ok := newMap[a.FullName()]; !ok {
			diff.Removed = append(diff.Removed, a)
		}
	}

	return diff
}

func accessMethodsEqual(a, b AccessMethod) bool {
	return a.Type == b.Type && a.Handler == b.Handler
}
