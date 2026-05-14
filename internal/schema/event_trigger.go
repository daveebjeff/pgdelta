package schema

import "fmt"

// EventTrigger represents a PostgreSQL event trigger.
type EventTrigger struct {
	Name     string
	Event    string // e.g. ddl_command_start, ddl_command_end, sql_drop, table_rewrite
	FuncName string
	Enabled  string // ENABLE, DISABLE, ENABLE REPLICA, ENABLE ALWAYS
	Tags     []string
}

func (e EventTrigger) FullName() string {
	return fmt.Sprintf("event_trigger.%s", e.Name)
}

// DiffEventTriggers returns added, removed, and changed event triggers.
func DiffEventTriggers(old, new []EventTrigger) (added, removed, changed []EventTrigger) {
	oldMap := make(map[string]EventTrigger, len(old))
	for _, et := range old {
		oldMap[et.Name] = et
	}
	newMap := make(map[string]EventTrigger, len(new))
	for _, et := range new {
		newMap[et.Name] = et
	}

	for _, et := range new {
		if prev, ok := oldMap[et.Name]; !ok {
			added = append(added, et)
		} else if !eventTriggersEqual(prev, et) {
			changed = append(changed, et)
		}
	}
	for _, et := range old {
		if _, ok := newMap[et.Name]; !ok {
			removed = append(removed, et)
		}
	}
	return
}

func eventTriggersEqual(a, b EventTrigger) bool {
	if a.Event != b.Event || a.FuncName != b.FuncName || a.Enabled != b.Enabled {
		return false
	}
	if len(a.Tags) != len(b.Tags) {
		return false
	}
	for i := range a.Tags {
		if a.Tags[i] != b.Tags[i] {
			return false
		}
	}
	return true
}
