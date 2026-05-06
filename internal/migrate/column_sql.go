package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/internal/schema"
)

// AddColumnSQL generates an ALTER TABLE ... ADD COLUMN statement.
func AddColumnSQL(tableName string, col schema.Column) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "ALTER TABLE %s ADD COLUMN %s %s", tableName, col.Name, col.DataType)
	if !col.Nullable {
		sb.WriteString(" NOT NULL")
	}
	if col.Default != nil {
		fmt.Fprintf(&sb, " DEFAULT %s", *col.Default)
	}
	sb.WriteString(";")
	return sb.String()
}

// DropColumnSQL generates an ALTER TABLE ... DROP COLUMN statement.
func DropColumnSQL(tableName string, col schema.Column) string {
	return fmt.Sprintf("ALTER TABLE %s DROP COLUMN %s;", tableName, col.Name)
}

// ColumnDiffSQL generates SQL statements to migrate from old columns to new columns.
func ColumnDiffSQL(tableName string, diff schema.ColumnDiff) []string {
	var stmts []string

	for _, col := range diff.Added {
		stmts = append(stmts, AddColumnSQL(tableName, col))
	}

	for _, col := range diff.Removed {
		stmts = append(stmts, DropColumnSQL(tableName, col))
	}

	for _, change := range diff.Modified {
		// Handle type change
		if change.Old.DataType != change.New.DataType {
			stmts = append(stmts, fmt.Sprintf(
				"ALTER TABLE %s ALTER COLUMN %s TYPE %s;",
				tableName, change.New.Name, change.New.DataType,
			))
		}
		// Handle nullability change
		if change.Old.Nullable != change.New.Nullable {
			if change.New.Nullable {
				stmts = append(stmts, fmt.Sprintf(
					"ALTER TABLE %s ALTER COLUMN %s DROP NOT NULL;",
					tableName, change.New.Name,
				))
			} else {
				stmts = append(stmts, fmt.Sprintf(
					"ALTER TABLE %s ALTER COLUMN %s SET NOT NULL;",
					tableName, change.New.Name,
				))
			}
		}
		// Handle default change
		if defaultChanged(change.Old.Default, change.New.Default) {
			if change.New.Default == nil {
				stmts = append(stmts, fmt.Sprintf(
					"ALTER TABLE %s ALTER COLUMN %s DROP DEFAULT;",
					tableName, change.New.Name,
				))
			} else {
				stmts = append(stmts, fmt.Sprintf(
					"ALTER TABLE %s ALTER COLUMN %s SET DEFAULT %s;",
					tableName, change.New.Name, *change.New.Default,
				))
			}
		}
	}

	return stmts
}

func defaultChanged(old, new *string) bool {
	if old == nil && new == nil {
		return false
	}
	if old == nil || new == nil {
		return true
	}
	return *old != *new
}
