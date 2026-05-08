package migrate

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

func baseRole() schema.Role {
	return schema.Role{
		Name:            "app_user",
		Inherit:         true,
		Login:           true,
		ConnectionLimit: -1,
	}
}

func TestCreateRoleSQL_Basic(t *testing.T) {
	r := baseRole()
	sql := CreateRoleSQL(r)
	assert.Contains(t, sql, `CREATE ROLE "app_user"`)
	assert.Contains(t, sql, "LOGIN")
	assert.Contains(t, sql, "INHERIT")
	assert.Contains(t, sql, "NOSUPERUSER")
}

func TestCreateRoleSQL_Superuser(t *testing.T) {
	r := baseRole()
	r.Superuser = true
	sql := CreateRoleSQL(r)
	assert.Contains(t, sql, "SUPERUSER")
	assert.NotContains(t, sql, "NOSUPERUSER")
}

func TestCreateRoleSQL_WithValidUntil(t *testing.T) {
	r := baseRole()
	v := "2030-12-31"
	r.ValidUntil = &v
	sql := CreateRoleSQL(r)
	assert.Contains(t, sql, "VALID UNTIL '2030-12-31'")
}

func TestCreateRoleSQL_WithConnectionLimit(t *testing.T) {
	r := baseRole()
	r.ConnectionLimit = 10
	sql := CreateRoleSQL(r)
	assert.Contains(t, sql, "CONNECTION LIMIT 10")
}

func TestDropRoleSQL(t *testing.T) {
	r := baseRole()
	sql := DropRoleSQL(r)
	assert.Equal(t, `DROP ROLE "app_user";`, sql)
}

func TestAlterRoleSQL_CreateDB(t *testing.T) {
	r := baseRole()
	r.CreateDB = true
	sql := AlterRoleSQL(r)
	assert.Contains(t, sql, `ALTER ROLE "app_user"`)
	assert.Contains(t, sql, "CREATEDB")
}

func TestRoleDiffSQL_AddedAndRemoved(t *testing.T) {
	added := baseRole()
	added.Name = "new_role"
	removed := baseRole()
	removed.Name = "old_role"
	diff := schema.RoleDiff{
		Added:   []schema.Role{added},
		Removed: []schema.Role{removed},
	}
	stmts := RoleDiffSQL(diff)
	assert.Len(t, stmts, 2)
	assert.Contains(t, stmts[0], `DROP ROLE "old_role"`)
	assert.Contains(t, stmts[1], `CREATE ROLE "new_role"`)
}

func TestRoleDiffSQL_Changed(t *testing.T) {
	r := baseRole()
	r.CreateDB = true
	diff := schema.RoleDiff{Changed: []schema.Role{r}}
	stmts := RoleDiffSQL(diff)
	assert.Len(t, stmts, 1)
	assert.Contains(t, stmts[0], "ALTER ROLE")
}
