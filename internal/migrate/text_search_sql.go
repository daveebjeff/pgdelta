package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateTextSearchConfigSQL generates SQL to create a text search configuration.
func CreateTextSearchConfigSQL(t schema.TextSearchConfig) string {
	return fmt.Sprintf(
		"CREATE TEXT SEARCH CONFIGURATION %s.%s (PARSER = %s);",
		t.Schema, t.Name, t.Parser,
	)
}

// DropTextSearchConfigSQL generates SQL to drop a text search configuration.
func DropTextSearchConfigSQL(t schema.TextSearchConfig) string {
	return fmt.Sprintf(
		"DROP TEXT SEARCH CONFIGURATION %s.%s;",
		t.Schema, t.Name,
	)
}

// AlterTextSearchConfigSQL generates SQL to alter a text search configuration parser.
func AlterTextSearchConfigSQL(t schema.TextSearchConfig) string {
	return fmt.Sprintf(
		"ALTER TEXT SEARCH CONFIGURATION %s.%s RENAME TO %s_old;\nDROP TEXT SEARCH CONFIGURATION %s.%s_old;\n%s",
		t.Schema, t.Name, t.Name,
		t.Schema, t.Name,
		CreateTextSearchConfigSQL(t),
	)
}

// TextSearchDiffSQL generates migration SQL for a TextSearchDiff.
func TextSearchDiffSQL(diff schema.TextSearchDiff) string {
	if diff.IsEmpty() {
		return ""
	}

	var stmts []string

	for _, t := range diff.Removed {
		stmts = append(stmts, DropTextSearchConfigSQL(t))
	}

	for _, t := range diff.Added {
		stmts = append(stmts, CreateTextSearchConfigSQL(t))
	}

	for _, t := range diff.Changed {
		stmts = append(stmts, AlterTextSearchConfigSQL(t))
	}

	return strings.Join(stmts, "\n")
}
