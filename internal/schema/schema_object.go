package schema

// SchemaObject represents a named object within a PostgreSQL schema.
type SchemaObject struct {
	Schema string
	Name   string
}

// FullName returns the fully qualified name of the schema object.
func (s SchemaObject) FullName() string {
	if s.Schema == "" {
		return s.Name
	}
	return s.Schema + "." + s.Name
}

// Privilege represents a PostgreSQL privilege granted on a schema object.
type Privilege struct {
	Grantee    string
	ObjectType string // e.g. TABLE, SEQUENCE, FUNCTION
	Schema     string
	ObjectName string
	Privileges []string // e.g. SELECT, INSERT, UPDATE, DELETE, EXECUTE
	WithGrant  bool
}

// FullName returns the fully qualified object name for this privilege.
func (p Privilege) FullName() string {
	if p.Schema == "" {
		return p.ObjectName
	}
	return p.Schema + "." + p.ObjectName
}

// DiffPrivileges computes added and removed privileges between two snapshots.
func DiffPrivileges(old, new []Privilege) (added, removed []Privilege) {
	oldMap := make(map[string]Privilege)
	for _, p := range old {
		oldMap[privilegeKey(p)] = p
	}
	newMap := make(map[string]Privilege)
	for _, p := range new {
		newMap[privilegeKey(p)] = p
	}
	for k, p := range newMap {
		if _, exists := oldMap[k]; !exists {
			added = append(added, p)
		}
	}
	for k, p := range oldMap {
		if _, exists := newMap[k]; !exists {
			removed = append(removed, p)
		}
	}
	return added, removed
}

func privilegeKey(p Privilege) string {
	return p.Grantee + "|" + p.ObjectType + "|" + p.Schema + "|" + p.ObjectName + "|" + joinPrivileges(p.Privileges)
}

func joinPrivileges(privs []string) string {
	result := ""
	for i, pr := range privs {
		if i > 0 {
			result += ","
		}
		result += pr
	}
	return result
}
