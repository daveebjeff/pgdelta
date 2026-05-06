package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateEnumSQL generates a CREATE TYPE ... AS ENUM statement.
func CreateEnumSQL(e schema.Enum) string {
	quoted := make([]string, len(e.Values))
	for i, v := range e.Values {
		quoted[i] = fmt.Sprintf("'%s'", v)
	}
	return fmt.Sprintf(
		"CREATE TYPE %s.%s AS ENUM (%s);",
		e.Schema, e.Name, strings.Join(quoted, ", "),
	)
}

// DropEnumSQL generates a DROP TYPE statement.
func DropEnumSQL(e schema.Enum) string {
	return fmt.Sprintf("DROP TYPE %s.%s;", e.Schema, e.Name)
}

// EnumDiffSQL generates SQL statements for all enum changes.
func EnumDiffSQL(diff schema.EnumDiff) []string {
	var stmts []string

	for _, e := range diff.Removed {
		stmts = append(stmts, DropEnumSQL(e))
	}

	for _, e := range diff.Added {
		stmts = append(stmts, CreateEnumSQL(e))
	}

	for _, c := range diff.Changed {
		stmts = append(stmts, alterEnumSQL(c.Old, c.New)...)
	}

	return stmts
}

// alterEnumSQL generates ALTER TYPE ... ADD VALUE statements for new enum values.
// Note: PostgreSQL does not support removing enum values without recreating the type.
func alterEnumSQL(old, new schema.Enum) []string {
	oldSet := make(map[string]struct{}, len(old.Values))
	for _, v := range old.Values {
		oldSet[v] = struct{}{}
	}

	var stmts []string
	for _, v := range new.Values {
		if _, exists := oldSet[v]; !exists {
			stmts = append(stmts, fmt.Sprintf(
				"ALTER TYPE %s.%s ADD VALUE '%s';",
				new.Schema, new.Name, v,
			))
		}
	}

	return stmts
}
