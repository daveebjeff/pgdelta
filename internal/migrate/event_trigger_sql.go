package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateEventTriggerSQL generates SQL to create an event trigger.
func CreateEventTriggerSQL(et schema.EventTrigger) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE EVENT TRIGGER %s ON %s", et.Name, et.Event))
	if len(et.Tags) > 0 {
		quoted := make([]string, len(et.Tags))
		for i, tag := range et.Tags {
			quoted[i] = fmt.Sprintf("'%s'", tag)
		}
		sb.WriteString(fmt.Sprintf(" WHEN TAG IN (%s)", strings.Join(quoted, ", ")))
	}
	sb.WriteString(fmt.Sprintf(" EXECUTE FUNCTION %s();", et.FuncName))
	return sb.String()
}

// DropEventTriggerSQL generates SQL to drop an event trigger.
func DropEventTriggerSQL(et schema.EventTrigger) string {
	return fmt.Sprintf("DROP EVENT TRIGGER %s;", et.Name)
}

// AlterEventTriggerSQL generates SQL to enable/disable an event trigger.
func AlterEventTriggerSQL(et schema.EventTrigger) string {
	return fmt.Sprintf("ALTER EVENT TRIGGER %s %s;", et.Name, et.Enabled)
}

// EventTriggerDiffSQL returns SQL statements to migrate event triggers from old to new state.
func EventTriggerDiffSQL(old, new []schema.EventTrigger) []string {
	added, removed, changed := schema.DiffEventTriggers(old, new)
	var stmts []string
	for _, et := range removed {
		stmts = append(stmts, DropEventTriggerSQL(et))
	}
	for _, et := range added {
		stmts = append(stmts, CreateEventTriggerSQL(et))
	}
	for _, et := range changed {
		stmts = append(stmts, DropEventTriggerSQL(et))
		stmts = append(stmts, CreateEventTriggerSQL(et))
		if et.Enabled != "ENABLE" {
			stmts = append(stmts, AlterEventTriggerSQL(et))
		}
	}
	return stmts
}
