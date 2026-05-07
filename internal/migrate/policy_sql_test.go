package migrate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/your-org/pgdelta/internal/schema"
)

var basePolicy = schema.Policy{
	Schema:     "public",
	Table:      "orders",
	Name:       "user_policy",
	Command:    "SELECT",
	Permissive: true,
	Roles:      []string{"app_user"},
	Using:      "(user_id = current_user_id())",
	WithCheck:  "",
}

func TestCreatePolicySQL(t *testing.T) {
	sql := CreatePolicySQL(basePolicy)
	assert.Equal(t,
		"CREATE POLICY user_policy ON public.orders AS PERMISSIVE FOR SELECT TO app_user USING (user_id = current_user_id());",
		sql)
}

func TestCreatePolicySQL_Restrictive(t *testing.T) {
	p := basePolicy
	p.Permissive = false
	p.Command = "ALL"
	p.Roles = nil
	p.Using = ""
	p.WithCheck = "(true)"

	sql := CreatePolicySQL(p)
	assert.Equal(t,
		"CREATE POLICY user_policy ON public.orders AS RESTRICTIVE FOR ALL WITH CHECK (true);",
		sql)
}

func TestDropPolicySQL(t *testing.T) {
	sql := DropPolicySQL(basePolicy)
	assert.Equal(t, "DROP POLICY user_policy ON public.orders;", sql)
}

func TestAlterPolicySQL(t *testing.T) {
	updated := basePolicy
	updated.Using = "(user_id = auth.uid())"

	sql := AlterPolicySQL(updated)
	assert.Equal(t,
		"ALTER POLICY user_policy ON public.orders TO app_user USING (user_id = auth.uid());",
		sql)
}

func TestPolicyDiffSQL_AddedAndRemoved(t *testing.T) {
	diff := schema.PolicyDiff{
		Added:   []schema.Policy{basePolicy},
		Removed: []schema.Policy{basePolicy},
	}

	stmts := PolicyDiffSQL(diff)
	assert.Len(t, stmts, 2)
	assert.Contains(t, stmts[0], "DROP POLICY")
	assert.Contains(t, stmts[1], "CREATE POLICY")
}

func TestPolicyDiffSQL_Changed(t *testing.T) {
	updated := basePolicy
	updated.Using = "(user_id = auth.uid())"

	diff := schema.PolicyDiff{
		Changed: []schema.Policy{updated},
	}

	stmts := PolicyDiffSQL(diff)
	assert.Len(t, stmts, 1)
	assert.Contains(t, stmts[0], "ALTER POLICY")
}
