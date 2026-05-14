package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func baseEventTrigger() EventTrigger {
	return EventTrigger{
		Name:     "audit_ddl",
		Event:    "ddl_command_end",
		FuncName: "public.audit_ddl_func",
		Enabled:  "ENABLE",
		Tags:     []string{"CREATE TABLE", "DROP TABLE"},
	}
}

func TestEventTriggerFullName(t *testing.T) {
	et := baseEventTrigger()
	assert.Equal(t, "event_trigger.audit_ddl", et.FullName())
}

func TestDiffEventTriggers_NoChanges(t *testing.T) {
	et := baseEventTrigger()
	added, removed, changed := DiffEventTriggers([]EventTrigger{et}, []EventTrigger{et})
	assert.Empty(t, added)
	assert.Empty(t, removed)
	assert.Empty(t, changed)
}

func TestDiffEventTriggers_AddedEventTrigger(t *testing.T) {
	et := baseEventTrigger()
	added, removed, changed := DiffEventTriggers(nil, []EventTrigger{et})
	assert.Equal(t, []EventTrigger{et}, added)
	assert.Empty(t, removed)
	assert.Empty(t, changed)
}

func TestDiffEventTriggers_RemovedEventTrigger(t *testing.T) {
	et := baseEventTrigger()
	added, removed, changed := DiffEventTriggers([]EventTrigger{et}, nil)
	assert.Empty(t, added)
	assert.Equal(t, []EventTrigger{et}, removed)
	assert.Empty(t, changed)
}

func TestDiffEventTriggers_ChangedEventTrigger(t *testing.T) {
	old := baseEventTrigger()
	new := baseEventTrigger()
	new.Enabled = "DISABLE"
	added, removed, changed := DiffEventTriggers([]EventTrigger{old}, []EventTrigger{new})
	assert.Empty(t, added)
	assert.Empty(t, removed)
	assert.Equal(t, []EventTrigger{new}, changed)
}

func TestDiffEventTriggers_ChangedTags(t *testing.T) {
	old := baseEventTrigger()
	new := baseEventTrigger()
	new.Tags = []string{"CREATE TABLE"}
	added, removed, changed := DiffEventTriggers([]EventTrigger{old}, []EventTrigger{new})
	assert.Empty(t, added)
	assert.Empty(t, removed)
	assert.Equal(t, []EventTrigger{new}, changed)
}
