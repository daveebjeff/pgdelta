package schema

// RowSecurity represents a row-level security enablement on a table.
type RowSecurity struct {
	Schema  string
	Table   string
	Enabled bool
	Forced  bool
}

// FullName returns a qualified identifier for the row security setting.
func (r RowSecurity) FullName() string {
	return r.Schema + "." + r.Table + ".row_security"
}

// RowSecurityDiff holds added, removed, and changed row security settings.
type RowSecurityDiff struct {
	Added   []RowSecurity
	Removed []RowSecurity
	Changed []RowSecurity
}

// IsEmpty returns true when there are no row security changes.
func (d RowSecurityDiff) IsEmpty() bool {
	return len(d.Added) == 0 && len(d.Removed) == 0 && len(d.Changed) == 0
}

// DiffRowSecurity computes the diff between two slices of RowSecurity settings.
func DiffRowSecurity(before, after []RowSecurity) RowSecurityDiff {
	var diff RowSecurityDiff

	beforeMap := make(map[string]RowSecurity, len(before))
	for _, rs := range before {
		beforeMap[rs.FullName()] = rs
	}

	afterMap := make(map[string]RowSecurity, len(after))
	for _, rs := range after {
		afterMap[rs.FullName()] = rs
	}

	for _, rs := range after {
		if prev, ok := beforeMap[rs.FullName()]; !ok {
			diff.Added = append(diff.Added, rs)
		} else if !rowSecurityEqual(prev, rs) {
			diff.Changed = append(diff.Changed, rs)
		}
	}

	for _, rs := range before {
		if _, ok := afterMap[rs.FullName()]; !ok {
			diff.Removed = append(diff.Removed, rs)
		}
	}

	return diff
}

func rowSecurityEqual(a, b RowSecurity) bool {
	return a.Enabled == b.Enabled && a.Forced == b.Forced
}
