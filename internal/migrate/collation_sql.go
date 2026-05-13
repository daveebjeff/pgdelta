package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/internal/schema"
)

// CreateCollationSQL generates a CREATE COLLATION statement.
func CreateCollationSQL(c schema.Collation) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "CREATE COLLATION %s (\n", c.FullName())
	fmt.Fprintf(&sb, "    PROVIDER = %s,\n", c.Provider)
	fmt.Fprintf(&sb, "    LOCALE = '%s'", c.Locale)
	if !c.Deterministic {
		sb.WriteString(",\n    DETERMINISTIC = false")
	}
	sb.WriteString("\n);")
	return sb.String()
}

// DropCollationSQL generates a DROP COLLATION statement.
func DropCollationSQL(c schema.Collation) string {
	return fmt.Sprintf("DROP COLLATION %s;", c.FullName())
}

// CollationDiffSQL generates SQL statements for a CollationDiff.
func CollationDiffSQL(diff schema.CollationDiff) []string {
	var stmts []string

	for _, c := range diff.Removed {
		stmts = append(stmts, DropCollationSQL(c))
	}

	for _, c := range diff.Added {
		stmts = append(stmts, CreateCollationSQL(c))
	}

	// Collations cannot be altered in PostgreSQL; drop and recreate on change.
	for _, c := range diff.Changed {
		stmts = append(stmts, DropCollationSQL(c))
		stmts = append(stmts, CreateCollationSQL(c))
	}

	return stmts
}
