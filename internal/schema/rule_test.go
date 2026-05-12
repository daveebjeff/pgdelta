package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseRule = Rule{
	Schema:     "public",
	Table:      "orders",
	Name:       "no_delete",
	Event:      "DELETE",
	Condition:  "",
	Instead:    true,
	Definition: "INSTEAD NOTHING",
}

func TestRuleFullName(t *testing.T) {
	assert.Equal(t, "public.orders.no_delete", baseRule.FullName())
}

func TestDiffRules_NoChanges(t *testing.T) {
	old := []Rule{baseRule}
	new := []Rule{baseRule}
	diff := DiffRules(old, new)
	assert.True(t, diff.IsEmpty())
}

func TestDiffRules_AddedRule(t *testing.T) {
	old := []Rule{}
	new := []Rule{baseRule}
	diff := DiffRules(old, new)
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, baseRule, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffRules_RemovedRule(t *testing.T) {
	old := []Rule{baseRule}
	new := []Rule{}
	diff := DiffRules(old, new)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, baseRule, diff.Removed[0])
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffRules_ChangedRule(t *testing.T) {
	old := []Rule{baseRule}
	modified := baseRule
	modified.Instead = false
	modified.Definition = "DO NOTHING"
	new := []Rule{modified}
	diff := DiffRules(old, new)
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, modified, diff.Changed[0])
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}
