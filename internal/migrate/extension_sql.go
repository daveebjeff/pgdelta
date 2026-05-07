package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateExtensionSQL returns the SQL to create an extension.
func CreateExtensionSQL(e schema.Extension) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE EXTENSION IF NOT EXISTS %q", e.Name))
	if e.Schema != "" {
		sb.WriteString(fmt.Sprintf(" WITH SCHEMA %q", e.Schema))
	}
	if e.Version != "" {
		sb.WriteString(fmt.Sprintf(" VERSION '%s'", e.Version))
	}
	sb.WriteString(";")
	return sb.String()
}

// DropExtensionSQL returns the SQL to drop an extension.
func DropExtensionSQL(e schema.Extension) string {
	return fmt.Sprintf("DROP EXTENSION IF EXISTS %q;", e.Name)
}

// AlterExtensionSQL returns the SQL to update an extension to a new version.
func AlterExtensionSQL(e schema.Extension) string {
	if e.Version == "" {
		return fmt.Sprintf("ALTER EXTENSION %q UPDATE;", e.Name)
	}
	return fmt.Sprintf("ALTER EXTENSION %q UPDATE TO '%s';", e.Name, e.Version)
}

// ExtensionDiffSQL generates migration SQL for an ExtensionDiff.
func ExtensionDiffSQL(diff schema.ExtensionDiff) []string {
	var stmts []string
	for _, e := range diff.Removed {
		stmts = append(stmts, DropExtensionSQL(e))
	}
	for _, e := range diff.Added {
		stmts = append(stmts, CreateExtensionSQL(e))
	}
	for _, e := range diff.Changed {
		stmts = append(stmts, AlterExtensionSQL(e))
	}
	return stmts
}
