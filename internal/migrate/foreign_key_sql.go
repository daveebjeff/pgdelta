package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

func AddForeignKeySQL(fk schema.ForeignKey) string {
	cols := strings.Join(fk.Columns, ", ")
	refCols := strings.Join(fk.RefColumns, ", ")

	sql := fmt.Sprintf(
		"ALTER TABLE %s.%s ADD CONSTRAINT %s FOREIGN KEY (%s) REFERENCES %s.%s (%s)",
		fk.Schema, fk.Table, fk.Name, cols, fk.RefSchema, fk.RefTable, refCols,
	)

	if fk.OnDelete != "" && fk.OnDelete != "NO ACTION" {
		sql += fmt.Sprintf(" ON DELETE %s", fk.OnDelete)
	}
	if fk.OnUpdate != "" && fk.OnUpdate != "NO ACTION" {
		sql += fmt.Sprintf(" ON UPDATE %s", fk.OnUpdate)
	}
	if fk.Deferrable {
		sql += " DEFERRABLE"
		if fk.InitiallyDeferred {
			sql += " INITIALLY DEFERRED"
		}
	}

	return sql + ";"
}

func DropForeignKeySQL(fk schema.ForeignKey) string {
	return fmt.Sprintf(
		"ALTER TABLE %s.%s DROP CONSTRAINT %s;",
		fk.Schema, fk.Table, fk.Name,
	)
}

func ForeignKeyDiffSQL(diff schema.ForeignKeyDiff) []string {
	var stmts []string

	for _, fk := range diff.Removed {
		stmts = append(stmts, DropForeignKeySQL(fk))
	}
	for _, fk := range diff.Changed {
		stmts = append(stmts, DropForeignKeySQL(fk))
		stmts = append(stmts, AddForeignKeySQL(fk))
	}
	for _, fk := range diff.Added {
		stmts = append(stmts, AddForeignKeySQL(fk))
	}

	return stmts
}
