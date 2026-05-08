package migrate

import (
	"fmt"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// AddConstraintSQL generates SQL to add a constraint to a table.
func AddConstraintSQL(c schema.Constraint) string {
	return fmt.Sprintf(
		"ALTER TABLE %s ADD CONSTRAINT %s %s;",
		c.TableFullName(),
		c.Name,
		c.Definition,
	)
}

// DropConstraintSQL generates SQL to drop a constraint from a table.
func DropConstraintSQL(c schema.Constraint) string {
	return fmt.Sprintf(
		"ALTER TABLE %s DROP CONSTRAINT %s;",
		c.TableFullName(),
		c.Name,
	)
}

// ConstraintDiffSQL generates SQL statements to migrate from one set of
// constraints to another.
func ConstraintDiffSQL(from, to []schema.Constraint) []string {
	diff := schema.DiffConstraints(from, to)
	var stmts []string

	for _, c := range diff.Removed {
		stmts = append(stmts, DropConstraintSQL(c))
	}

	for _, c := range diff.Added {
		stmts = append(stmts, AddConstraintSQL(c))
	}

	for _, c := range diff.Changed {
		// Constraints cannot be altered in-place; drop and recreate.
		stmts = append(stmts, DropConstraintSQL(c.From))
		stmts = append(stmts, AddConstraintSQL(c.To))
	}

	return stmts
}
