package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateViewSQL generates a CREATE OR REPLACE VIEW statement.
func CreateViewSQL(v schema.View) string {
	return fmt.Sprintf(
		"CREATE OR REPLACE VIEW %s AS\n%s;",
		v.FullName(),
		strings.TrimRight(v.Definition, ";"),
	)
}

// DropViewSQL generates a DROP VIEW statement.
func DropViewSQL(v schema.View) string {
	return fmt.Sprintf("DROP VIEW IF EXISTS %s;", v.FullName())
}

// ViewDiffSQL generates SQL statements for all view changes in a ViewDiff.
func ViewDiffSQL(diff schema.ViewDiff) []string {
	var stmts []string

	for _, v := range diff.Removed {
		stmts = append(stmts, DropViewSQL(v))
	}

	for _, v := range diff.Added {
		stmts = append(stmts, CreateViewSQL(v))
	}

	for _, change := range diff.Changed {
		// Use CREATE OR REPLACE to update an existing view definition.
		stmts = append(stmts, CreateViewSQL(change.New))
	}

	return stmts
}
