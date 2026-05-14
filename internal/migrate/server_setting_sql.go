package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// AlterSystemSetSQL generates an ALTER SYSTEM SET statement for a server setting.
func AlterSystemSetSQL(s schema.ServerSetting) string {
	return fmt.Sprintf("ALTER SYSTEM SET %s = '%s';", s.Name, escapeSingleQuotes(s.Value))
}

// AlterSystemResetSQL generates an ALTER SYSTEM RESET statement to remove a setting.
func AlterSystemResetSQL(s schema.ServerSetting) string {
	return fmt.Sprintf("ALTER SYSTEM RESET %s;", s.Name)
}

// ServerSettingDiffSQL generates SQL statements for all server setting changes.
func ServerSettingDiffSQL(diff schema.ServerSettingDiff) []string {
	var stmts []string

	for _, s := range diff.Added {
		stmts = append(stmts, AlterSystemSetSQL(s))
	}

	for _, c := range diff.Changed {
		stmts = append(stmts, AlterSystemSetSQL(c.New))
	}

	for _, s := range diff.Removed {
		stmts = append(stmts, AlterSystemResetSQL(s))
	}

	return stmts
}

func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}
