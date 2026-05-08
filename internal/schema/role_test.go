package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func baseRole() Role {
	return Role{
		Name:            "app_user",
		Inherit:         true,
		Login:           true,
		ConnectionLimit: -1,
	}
}

func TestRoleFullName(t *testing.T) {
	r := baseRole()
	assert.Equal(t, `"app_user"`, r.FullName())
}

func TestDiffRoles_NoChanges(t *testing.T) {
	r := baseRole()
	diff := DiffRoles([]Role{r}, []Role{r})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffRoles_AddedRole(t *testing.T) {
	r := baseRole()
	diff := DiffRoles([]Role{}, []Role{r})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, r.Name, diff.Added[0].Name)
}

func TestDiffRoles_RemovedRole(t *testing.T) {
	r := baseRole()
	diff := DiffRoles([]Role{r}, []Role{})
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, r.Name, diff.Removed[0].Name)
}

func TestDiffRoles_ChangedRole(t *testing.T) {
	old := baseRole()
	new := baseRole()
	new.CreateDB = true
	diff := DiffRoles([]Role{old}, []Role{new})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Len(t, diff.Changed, 1)
	assert.True(t, diff.Changed[0].CreateDB)
}

func TestDiffRoles_ValidUntilChanged(t *testing.T) {
	old := baseRole()
	new := baseRole()
	v := "2030-01-01"
	new.ValidUntil = &v
	diff := DiffRoles([]Role{old}, []Role{new})
	assert.Len(t, diff.Changed, 1)
}
