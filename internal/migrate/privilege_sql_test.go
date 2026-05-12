package migrate

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

func basePriv() schema.Privilege {
	return schema.Privilege{
		Grantee:    "app_user",
		ObjectType: "TABLE",
		Schema:     "public",
		ObjectName: "orders",
		Privileges: []string{"SELECT", "INSERT"},
		WithGrant:  false,
	}
}

func TestGrantPrivilegeSQL_Basic(t *testing.T) {
	p := basePriv()
	sql := GrantPrivilegeSQL(p)
	assert.Equal(t, "GRANT SELECT, INSERT ON TABLE public.orders TO app_user;", sql)
}

func TestGrantPrivilegeSQL_WithGrantOption(t *testing.T) {
	p := basePriv()
	p.WithGrant = true
	sql := GrantPrivilegeSQL(p)
	assert.Equal(t, "GRANT SELECT, INSERT ON TABLE public.orders TO app_user WITH GRANT OPTION;", sql)
}

func TestRevokePrivilegeSQL(t *testing.T) {
	p := basePriv()
	sql := RevokePrivilegeSQL(p)
	assert.Equal(t, "REVOKE SELECT, INSERT ON TABLE public.orders FROM app_user;", sql)
}

func TestPrivilegeDiffSQL_AddedAndRemoved(t *testing.T) {
	old := basePriv()
	newPriv := basePriv()
	newPriv.ObjectName = "invoices"
	stmts := PrivilegeDiffSQL([]schema.Privilege{old}, []schema.Privilege{newPriv})
	assert.Len(t, stmts, 2)
	assert.Contains(t, stmts[0], "REVOKE")
	assert.Contains(t, stmts[1], "GRANT")
}

func TestPrivilegeDiffSQL_NoChanges(t *testing.T) {
	p := basePriv()
	stmts := PrivilegeDiffSQL([]schema.Privilege{p}, []schema.Privilege{p})
	assert.Empty(t, stmts)
}

func TestPrivilegeDiffSQL_ExecuteFunction(t *testing.T) {
	p := schema.Privilege{
		Grantee:    "readonly",
		ObjectType: "FUNCTION",
		Schema:     "public",
		ObjectName: "calculate_total",
		Privileges: []string{"EXECUTE"},
	}
	sql := GrantPrivilegeSQL(p)
	assert.Equal(t, "GRANT EXECUTE ON FUNCTION public.calculate_total TO readonly;", sql)
}
