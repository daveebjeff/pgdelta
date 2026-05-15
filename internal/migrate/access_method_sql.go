package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateAccessMethodSQL generates a CREATE ACCESS METHOD statement.
func CreateAccessMethodSQL(a schema.AccessMethod) string {
	return fmt.Sprintf(
		"CREATE ACCESS METHOD %s TYPE %s HANDLER %s;",
		a.Name,
		strings.ToUpper(a.Type),
		a.Handler,
	)
}

// DropAccessMethodSQL generates a DROP ACCESS METHOD statement.
func DropAccessMethodSQL(a schema.AccessMethod) string {
	return fmt.Sprintf("DROP ACCESS METHOD %s;", a.Name)
}

// AccessMethodDiffSQL generates SQL statements for all access method changes.
func AccessMethodDiffSQL(diff schema.AccessMethodDiff) []string {
	var stmts []string

	for _, a := range diff.Removed {
		stmts = append(stmts, DropAccessMethodSQL(a))
	}

	for _, a := range diff.Added {
		stmts = append(stmts, CreateAccessMethodSQL(a))
	}

	// Access methods cannot be altered; drop and recreate on change.
	for _, a := range diff.Changed {
		stmts = append(stmts, DropAccessMethodSQL(a))
		stmts = append(stmts, CreateAccessMethodSQL(a))
	}

	return stmts
}
