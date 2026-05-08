package schema

import "fmt"

// Role represents a PostgreSQL role/user.
type Role struct {
	Name            string
	Superuser       bool
	Inherit         bool
	CreateRole      bool
	CreateDB        bool
	Login           bool
	Replication     bool
	BypassRLS       bool
	ConnectionLimit int
	ValidUntil      *string
}

// FullName returns the quoted role name.
func (r Role) FullName() string {
	return fmt.Sprintf("%q", r.Name)
}

// RoleDiff holds added, removed, and changed roles.
type RoleDiff struct {
	Added   []Role
	Removed []Role
	Changed []Role
}

// DiffRoles computes the diff between two slices of roles.
func DiffRoles(old, new []Role) RoleDiff {
	oldMap := make(map[string]Role, len(old))
	for _, r := range old {
		oldMap[r.Name] = r
	}
	newMap := make(map[string]Role, len(new))
	for _, r := range new {
		newMap[r.Name] = r
	}

	var diff RoleDiff
	for _, r := range new {
		if o, exists := oldMap[r.Name]; !exists {
			diff.Added = append(diff.Added, r)
		} else if !rolesEqual(o, r) {
			diff.Changed = append(diff.Changed, r)
		}
	}
	for _, r := range old {
		if _, exists := newMap[r.Name]; !exists {
			diff.Removed = append(diff.Removed, r)
		}
	}
	return diff
}

func rolesEqual(a, b Role) bool {
	if a.Superuser != b.Superuser ||
		a.Inherit != b.Inherit ||
		a.CreateRole != b.CreateRole ||
		a.CreateDB != b.CreateDB ||
		a.Login != b.Login ||
		a.Replication != b.Replication ||
		a.BypassRLS != b.BypassRLS ||
		a.ConnectionLimit != b.ConnectionLimit {
		return false
	}
	if (a.ValidUntil == nil) != (b.ValidUntil == nil) {
		return false
	}
	if a.ValidUntil != nil && *a.ValidUntil != *b.ValidUntil {
		return false
	}
	return true
}
