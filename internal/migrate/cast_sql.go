package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateCastSQL generates a CREATE CAST statement.
func CreateCastSQL(c schema.Cast) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE CAST (%s AS %s)", c.SourceType, c.TargetType))
	if c.FunctionName != "" {
		sb.WriteString(fmt.Sprintf(" WITH FUNCTION %s", c.FunctionName))
	} else {
		sb.WriteString(" WITHOUT FUNCTION")
	}
	switch c.CastContext {
	case "a":
		sb.WriteString(" AS ASSIGNMENT")
	case "i":
		sb.WriteString(" AS IMPLICIT")
	// "e" is the default (explicit), no clause needed
	}
	sb.WriteString(";")
	return sb.String()
}

// DropCastSQL generates a DROP CAST statement.
func DropCastSQL(c schema.Cast) string {
	return fmt.Sprintf("DROP CAST (%s AS %s);", c.SourceType, c.TargetType)
}

// CastDiffSQL generates SQL statements for cast differences.
func CastDiffSQL(diff schema.CastDiff) []string {
	if diff.IsEmpty() {
		return nil
	}

	var stmts []string

	for _, c := range diff.Removed {
		stmts = append(stmts, DropCastSQL(c))
	}
	for _, c := range diff.Changed {
		stmts = append(stmts, DropCastSQL(c))
		stmts = append(stmts, CreateCastSQL(c))
	}
	for _, c := range diff.Added {
		stmts = append(stmts, CreateCastSQL(c))
	}

	return stmts
}
