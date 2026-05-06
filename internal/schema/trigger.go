package schema

import "fmt"

// Trigger represents a PostgreSQL trigger definition.
type Trigger struct {
	Schema    string
	Table     string
	Name      string
	Timing    string // BEFORE, AFTER, INSTEAD OF
	Events    []string // INSERT, UPDATE, DELETE, TRUNCATE
	ForEach   string // ROW, STATEMENT
	Condition string // optional WHEN clause
	Function  string // function to call
}

// FullName returns the qualified trigger name as schema.table.trigger.
func (t Trigger) FullName() string {
	return fmt.Sprintf("%s.%s.%s", t.Schema, t.Table, t.Name)
}

// TriggerDiff holds added, removed, and changed triggers.
type TriggerDiff struct {
	Added   []Trigger
	Removed []Trigger
	Changed []Trigger
}

// DiffTriggers compares two slices of triggers and returns the diff.
func DiffTriggers(old, new []Trigger) TriggerDiff {
	oldMap := make(map[string]Trigger, len(old))
	for _, t := range old {
		oldMap[t.FullName()] = t
	}

	newMap := make(map[string]Trigger, len(new))
	for _, t := range new {
		newMap[t.FullName()] = t
	}

	var diff TriggerDiff

	for _, t := range new {
		if _, exists := oldMap[t.FullName()]; !exists {
			diff.Added = append(diff.Added, t)
		} else if !triggersEqual(oldMap[t.FullName()], t) {
			diff.Changed = append(diff.Changed, t)
		}
	}

	for _, t := range old {
		if _, exists := newMap[t.FullName()]; !exists {
			diff.Removed = append(diff.Removed, t)
		}
	}

	return diff
}

func triggersEqual(a, b Trigger) bool {
	if a.Timing != b.Timing || a.ForEach != b.ForEach ||
		a.Condition != b.Condition || a.Function != b.Function {
		return false
	}
	if len(a.Events) != len(b.Events) {
		return false
	}
	for i := range a.Events {
		if a.Events[i] != b.Events[i] {
			return false
		}
	}
	return true
}
