package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseTrigger = Trigger{
	Schema:   "public",
	Table:    "orders",
	Name:     "trg_audit",
	Timing:   "AFTER",
	Events:   []string{"INSERT", "UPDATE"},
	ForEach:  "ROW",
	Function: "audit_log()",
}

func TestTriggerFullName(t *testing.T) {
	assert.Equal(t, "public.orders.trg_audit", baseTrigger.FullName())
}

func TestDiffTriggers_NoChanges(t *testing.T) {
	old := []Trigger{baseTrigger}
	new := []Trigger{baseTrigger}
	diff := DiffTriggers(old, new)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffTriggers_AddedTrigger(t *testing.T) {
	newTrigger := Trigger{
		Schema:   "public",
		Table:    "orders",
		Name:     "trg_notify",
		Timing:   "AFTER",
		Events:   []string{"INSERT"},
		ForEach:  "ROW",
		Function: "notify_user()",
	}
	diff := DiffTriggers([]Trigger{baseTrigger}, []Trigger{baseTrigger, newTrigger})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, newTrigger.FullName(), diff.Added[0].FullName())
	assert.Empty(t, diff.Removed)
}

func TestDiffTriggers_RemovedTrigger(t *testing.T) {
	diff := DiffTriggers([]Trigger{baseTrigger}, []Trigger{})
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, baseTrigger.FullName(), diff.Removed[0].FullName())
	assert.Empty(t, diff.Added)
}

func TestDiffTriggers_ChangedTrigger(t *testing.T) {
	modified := baseTrigger
	modified.Timing = "BEFORE"
	diff := DiffTriggers([]Trigger{baseTrigger}, []Trigger{modified})
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, "BEFORE", diff.Changed[0].Timing)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}
