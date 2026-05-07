package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/internal/schema"
)

// CreateTableSQL generates a CREATE TABLE statement for the given table.
func CreateTableSQL(t schema.Table) string {
	var cols []string
	for _, col := range t.Columns {
		cols = append(cols, "\t"+columnDefinition(col))
	}
	return fmt.Sprintf("CREATE TABLE %s (\n%s\n);", t.FullName(), strings.Join(cols, ",\n"))
}

// DropTableSQL generates a DROP TABLE statement for the given table.
func DropTableSQL(t schema.Table) string {
	return fmt.Sprintf("DROP TABLE %s;", t.FullName())
}

// TableDiffSQL generates the SQL statements needed to migrate from old to new table state.
func TableDiffSQL(diff schema.TableDiff) []string {
	var stmts []string

	for _, t := range diff.Added {
		stmts = append(stmts, CreateTableSQL(t))
	}

	for _, t := range diff.Removed {
		stmts = append(stmts, DropTableSQL(t))
	}

	for _, td := range diff.Modified {
		colDiff := schema.DiffColumns(td.Old.Columns, td.New.Columns)
		for _, s := range ColumnDiffSQL(td.Old.FullName(), colDiff) {
			stmts = append(stmts, s)
		}
	}

	return stmts
}

// columnDefinition returns the SQL column definition string for use in CREATE TABLE.
func columnDefinition(col schema.Column) string {
	def := fmt.Sprintf("%s %s", col.Name, col.DataType)
	if !col.Nullable {
		def += " NOT NULL"
	}
	if col.Default != nil {
		def += fmt.Sprintf(" DEFAULT %s", *col.Default)
	}
	return def
}
