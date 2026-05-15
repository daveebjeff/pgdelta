package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// EnableRowSecuritySQL returns SQL to enable row-level security on a table.
func EnableRowSecuritySQL(rs schema.RowSecurity) string {
	return fmt.Sprintf("ALTER TABLE %s.%s ENABLE ROW LEVEL SECURITY;",
		rs.Schema, rs.Table)
}

// DisableRowSecuritySQL returns SQL to disable row-level security on a table.
func DisableRowSecuritySQL(rs schema.RowSecurity) string {
	return fmt.Sprintf("ALTER TABLE %s.%s DISABLE ROW LEVEL SECURITY;",
		rs.Schema, rs.Table)
}

// ForceRowSecuritySQL returns SQL to force row-level security for table owners.
func ForceRowSecuritySQL(rs schema.RowSecurity) string {
	return fmt.Sprintf("ALTER TABLE %s.%s FORCE ROW LEVEL SECURITY;",
		rs.Schema, rs.Table)
}

// NoForceRowSecuritySQL returns SQL to remove forced row-level security.
func NoForceRowSecuritySQL(rs schema.RowSecurity) string {
	return fmt.Sprintf("ALTER TABLE %s.%s NO FORCE ROW LEVEL SECURITY;",
		rs.Schema, rs.Table)
}

// RowSecurityDiffSQL generates migration SQL from a RowSecurityDiff.
func RowSecurityDiffSQL(diff schema.RowSecurityDiff) string {
	var stmts []string

	for _, rs := range diff.Added {
		if rs.Enabled {
			stmts = append(stmts, EnableRowSecuritySQL(rs))
		}
		if rs.Forced {
			stmts = append(stmts, ForceRowSecuritySQL(rs))
		}
	}

	for _, rs := range diff.Changed {
		if rs.Enabled {
			stmts = append(stmts, EnableRowSecuritySQL(rs))
		} else {
			stmts = append(stmts, DisableRowSecuritySQL(rs))
		}
		if rs.Forced {
			stmts = append(stmts, ForceRowSecuritySQL(rs))
		} else {
			stmts = append(stmts, NoForceRowSecuritySQL(rs))
		}
	}

	for _, rs := range diff.Removed {
		stmts = append(stmts, DisableRowSecuritySQL(rs))
	}

	return strings.Join(stmts, "\n")
}
