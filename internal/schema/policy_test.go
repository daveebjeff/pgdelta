package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var basePolicy = Policy{
	Schema:     "public",
	Table:      "orders",
	Name:       "user_policy",
	Command:    "SELECT",
	Permissive: true,
	Roles:      []string{"app_user"},
	Using:      "(user_id = current_user_id())",
	WithCheck:  "",
}

func TestPolicyFullName(t *testing.T) {
	assert.Equal(t, "public.orders.user_policy", basePolicy.FullName())
}

func TestDiffPolicies_NoChanges(t *testing.T) {
	diff := DiffPolicies([]Policy{basePolicy}, []Policy{basePolicy})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffPolicies_AddedPolicy(t *testing.T) {
	diff := DiffPolicies(nil, []Policy{basePolicy})
	assert.Len(t, diff.Added, 1)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffPolicies_RemovedPolicy(t *testing.T) {
	diff := DiffPolicies([]Policy{basePolicy}, nil)
	assert.Empty(t, diff.Added)
	assert.Len(t, diff.Removed, 1)
	assert.Empty(t, diff.Changed)
}

func TestDiffPolicies_ChangedPolicy(t *testing.T) {
	updated := basePolicy
	updated.Using = "(user_id = auth.uid())"

	diff := DiffPolicies([]Policy{basePolicy}, []Policy{updated})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Len(t, diff.Changed, 1)
}

func TestDiffPolicies_ChangedRoles(t *testing.T) {
	updated := basePolicy
	updated.Roles = []string{"app_user", "admin"}

	diff := DiffPolicies([]Policy{basePolicy}, []Policy{updated})
	assert.Len(t, diff.Changed, 1)
}
