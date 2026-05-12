package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateSchemaSQL generates a CREATE SCHEMA statement.
func CreateSchemaSQL(s schema.SchemaNamespace) string {
	sql := fmt.Sprintf("CREATE SCHEMA %s", s.Name)
	if s.Owner != "" {
		sql += fmt.Sprintf(" AUTHORIZATION %s", s.Owner)
	}
	return sql + ";"
}

// DropSchemaSQL generates a DROP SCHEMA statement.
func DropSchemaSQL(s schema.SchemaNamespace) string {
	return fmt.Sprintf("DROP SCHEMA IF EXISTS %s;", s.Name)
}

// AlterSchemaOwnerSQL generates an ALTER SCHEMA ... OWNER TO statement.
func AlterSchemaOwnerSQL(s schema.SchemaNamespace) string {
	return fmt.Sprintf("ALTER SCHEMA %s OWNER TO %s;", s.Name, s.Owner)
}

// SchemaNamespaceDiffSQL generates SQL statements for a SchemaNSDiff.
func SchemaNamespaceDiffSQL(diff schema.SchemaNSDiff) string {
	var stmts []string
	for _, s := range diff.Added {
		stmts = append(stmts, CreateSchemaSQL(s))
	}
	for _, s := range diff.Changed {
		stmts = append(stmts, AlterSchemaOwnerSQL(s))
	}
	for _, s := range diff.Removed {
		stmts = append(stmts, DropSchemaSQL(s))
	}
	return strings.Join(stmts, "\n")
}
