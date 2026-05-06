package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateIndexSQL generates a CREATE INDEX statement for the given index.
// If the index is unique, the UNIQUE keyword is included.
// If the index method is not the default (btree), the USING clause is included.
func CreateIndexSQL(idx schema.Index) string {
	unique := ""
	if idx.Unique {
		unique = "UNIQUE "
	}

	method := ""
	if idx.Method != "" && idx.Method != schema.IndexMethodBTree {
		method = fmt.Sprintf(" USING %s", idx.Method)
	}

	columns := strings.Join(idx.Columns, ", ")

	return fmt.Sprintf(
		"CREATE %sINDEX %s ON %s.%s%s (%s);",
		unique,
		idx.Name,
		idx.SchemaName,
		idx.TableName,
		method,
		columns,
	)
}

// DropIndexSQL generates a DROP INDEX statement for the given index.
// The schema-qualified index name is used to avoid ambiguity.
func DropIndexSQL(idx schema.Index) string {
	return fmt.Sprintf("DROP INDEX %s.%s;", idx.SchemaName, idx.Name)
}

// RenameIndexSQL generates an ALTER INDEX ... RENAME TO statement.
func RenameIndexSQL(schemaName, oldName, newName string) string {
	return fmt.Sprintf("ALTER INDEX %s.%s RENAME TO %s;", schemaName, oldName, newName)
}

// IndexDiffSQL generates SQL statements for all index differences.
// Removed indexes are dropped first, changed indexes are dropped and recreated,
// and new indexes are created last.
func IndexDiffSQL(diff schema.IndexDiff) []string {
	var stmts []string

	for _, idx := range diff.Removed {
		stmts = append(stmts, DropIndexSQL(idx))
	}

	for _, change := range diff.Changed {
		stmts = append(stmts, DropIndexSQL(change.Old))
		stmts = append(stmts, CreateIndexSQL(change.New))
	}

	for _, idx := range diff.Added {
		stmts = append(stmts, CreateIndexSQL(idx))
	}

	return stmts
}
