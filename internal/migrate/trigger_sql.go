package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateTriggerSQL generates a CREATE TRIGGER statement.
func CreateTriggerSQL(t schema.Trigger) string {
	events := strings.Join(t.Events, " OR ")
	sql := fmt.Sprintf(
		"CREATE TRIGGER %s\n  %s %s ON %s.%s\n  FOR EACH %s",
		t.Name, t.Timing, events, t.Schema, t.Table, t.ForEach,
	)
	if t.Condition != "" {
		sql += fmt.Sprintf("\n  WHEN (%s)", t.Condition)
	}
	sql += fmt.Sprintf("\n  EXECUTE FUNCTION %s;", t.Function)
	return sql
}

// DropTriggerSQL generates a DROP TRIGGER statement.
func DropTriggerSQL(t schema.Trigger) string {
	return fmt.Sprintf("DROP TRIGGER IF EXISTS %s ON %s.%s;", t.Name, t.Schema, t.Table)
}

// TriggerDiffSQL generates SQL statements for a TriggerDiff.
// Changed triggers are handled as DROP + CREATE since ALTER TRIGGER is limited.
func TriggerDiffSQL(diff schema.TriggerDiff) []string {
	var stmts []string

	for _, t := range diff.Removed {
		stmts = append(stmts, DropTriggerSQL(t))
	}

	for _, t := range diff.Changed {
		stmts = append(stmts, DropTriggerSQL(t))
		stmts = append(stmts, CreateTriggerSQL(t))
	}

	for _, t := range diff.Added {
		stmts = append(stmts, CreateTriggerSQL(t))
	}

	return stmts
}
