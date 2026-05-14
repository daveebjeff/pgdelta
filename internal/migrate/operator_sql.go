package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/internal/schema"
)

// CreateOperatorSQL generates a CREATE OPERATOR statement.
func CreateOperatorSQL(o schema.Operator) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("PROCEDURE = %s", o.Procedure))
	if o.LeftType != "" {
		parts = append(parts, fmt.Sprintf("LEFTARG = %s", o.LeftType))
	}
	if o.RightType != "" {
		parts = append(parts, fmt.Sprintf("RIGHTARG = %s", o.RightType))
	}
	if o.Commutator != "" {
		parts = append(parts, fmt.Sprintf("COMMUTATOR = %s", o.Commutator))
	}
	if o.Negator != "" {
		parts = append(parts, fmt.Sprintf("NEGATOR = %s", o.Negator))
	}
	return fmt.Sprintf(
		"CREATE OPERATOR %s.%s (\n    %s\n);",
		o.Schema, o.Name,
		strings.Join(parts, ",\n    "),
	)
}

// DropOperatorSQL generates a DROP OPERATOR statement.
func DropOperatorSQL(o schema.Operator) string {
	return fmt.Sprintf(
		"DROP OPERATOR %s.%s (%s, %s);",
		o.Schema, o.Name, o.LeftType, o.RightType,
	)
}

// OperatorDiffSQL generates SQL statements for operator differences.
func OperatorDiffSQL(diff schema.OperatorDiff) []string {
	var stmts []string
	for _, o := range diff.Removed {
		stmts = append(stmts, DropOperatorSQL(o))
	}
	for _, o := range diff.Changed {
		stmts = append(stmts, DropOperatorSQL(o))
		stmts = append(stmts, CreateOperatorSQL(o))
	}
	for _, o := range diff.Added {
		stmts = append(stmts, CreateOperatorSQL(o))
	}
	return stmts
}
