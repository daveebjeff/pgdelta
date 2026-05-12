package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateRuleSQL generates a CREATE RULE statement.
func CreateRuleSQL(r schema.Rule) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE RULE %s AS ON %s TO %s.%s",
		r.Name, r.Event, r.Schema, r.Table))
	if r.Condition != "" {
		sb.WriteString(fmt.Sprintf(" WHERE (%s)", r.Condition))
	}
	if r.Instead {
		sb.WriteString(" DO INSTEAD")
	} else {
		sb.WriteString(" DO ALSO")
	}
	sb.WriteString(fmt.Sprintf(" %s;", r.Definition))
	return sb.String()
}

// DropRuleSQL generates a DROP RULE statement.
func DropRuleSQL(r schema.Rule) string {
	return fmt.Sprintf("DROP RULE %s ON %s.%s;", r.Name, r.Schema, r.Table)
}

// RuleDiffSQL generates SQL statements for all rule changes.
func RuleDiffSQL(diff schema.RuleDiff) []string {
	var stmts []string

	for _, r := range diff.Removed {
		stmts = append(stmts, DropRuleSQL(r))
	}

	for _, r := range diff.Changed {
		// Rules cannot be altered; drop and recreate.
		stmts = append(stmts, DropRuleSQL(r))
		stmts = append(stmts, CreateRuleSQL(r))
	}

	for _, r := range diff.Added {
		stmts = append(stmts, CreateRuleSQL(r))
	}

	return stmts
}
