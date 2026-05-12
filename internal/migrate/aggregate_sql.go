package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateAggregateSQL generates a CREATE AGGREGATE statement.
func CreateAggregateSQL(a schema.Aggregate) string {
	args := strings.Join(a.ArgTypes, ", ")
	var sb strings.Builder
	fmt.Fprintf(&sb, "CREATE AGGREGATE %s.%s(%s) (\n", a.Schema, a.Name, args)
	fmt.Fprintf(&sb, "    SFUNC = %s,\n", a.SFuncName)
	fmt.Fprintf(&sb, "    STYPE = %s", a.SType)
	if a.InitCond != nil {
		fmt.Fprintf(&sb, ",\n    INITCOND = '%s'", *a.InitCond)
	}
	if a.FinalFunc != nil {
		fmt.Fprintf(&sb, ",\n    FINALFUNC = %s", *a.FinalFunc)
	}
	sb.WriteString("\n);")
	return sb.String()
}

// DropAggregateSQL generates a DROP AGGREGATE statement.
func DropAggregateSQL(a schema.Aggregate) string {
	args := strings.Join(a.ArgTypes, ", ")
	return fmt.Sprintf("DROP AGGREGATE %s.%s(%s);", a.Schema, a.Name, args)
}

// AggregateDiffSQL generates SQL statements for all aggregate changes.
func AggregateDiffSQL(diff schema.AggregateDiff) []string {
	var stmts []string
	for _, a := range diff.Removed {
		stmts = append(stmts, DropAggregateSQL(a))
	}
	for _, a := range diff.Changed {
		stmts = append(stmts, DropAggregateSQL(a))
		stmts = append(stmts, CreateAggregateSQL(a))
	}
	for _, a := range diff.Added {
		stmts = append(stmts, CreateAggregateSQL(a))
	}
	return stmts
}
