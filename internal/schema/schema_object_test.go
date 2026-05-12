package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func basePrivilege() Privilege {
	return Privilege{
		Grantee:    "app_user",
		ObjectType: "TABLE",
		Schema:     "public",
		ObjectName: "users",
		Privileges: []string{"SELECT", "INSERT"},
		WithGrant:  false,
	}
}

func TestPrivilegeFullName(t *testing.T) {
	p := basePrivilege()
	assert.Equal(t, "public.users", p.FullName())

	p.Schema = ""
	assert.Equal(t, "users", p.FullName())
}

func TestDiffPrivileges_NoChanges(t *testing.T) {
	p := basePrivilege()
	added, removed := DiffPrivileges([]Privilege{p}, []Privilege{p})
	assert.Empty(t, added)
	assert.Empty(t, removed)
}

func TestDiffPrivileges_AddedPrivilege(t *testing.T) {
	p := basePrivilege()
	added, removed := DiffPrivileges([]Privilege{}, []Privilege{p})
	assert.Len(t, added, 1)
	assert.Empty(t, removed)
	assert.Equal(t, p, added[0])
}

func TestDiffPrivileges_RemovedPrivilege(t *testing.T) {
	p := basePrivilege()
	added, removed := DiffPrivileges([]Privilege{p}, []Privilege{})
	assert.Empty(t, added)
	assert.Len(t, removed, 1)
	assert.Equal(t, p, removed[0])
}

func TestDiffPrivileges_ChangedPrivilege(t *testing.T) {
	old := basePrivilege()
	new := basePrivilege()
	new.Privileges = []string{"SELECT", "INSERT", "UPDATE"}
	added, removed := DiffPrivileges([]Privilege{old}, []Privilege{new})
	assert.Len(t, added, 1)
	assert.Len(t, removed, 1)
}
