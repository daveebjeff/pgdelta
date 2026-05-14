package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateStatisticsSQL generates a CREATE STATISTICS statement.
func CreateStatisticsSQL(s schema.Statistic) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE STATISTICS %s.%s", s.Schema, s.Name))

	if len(s.Kinds) > 0 {
		sb.WriteString(fmt.Sprintf(" (%s)", strings.Join(s.Kinds, ", ")))
	}

	cols := make([]string, len(s.Columns))
	for i, c := range s.Columns {
		cols[i] = c
	}
	sb.WriteString(fmt.Sprintf(" ON %s FROM %s.%s;",
		strings.Join(cols, ", "),
		s.Schema,
		s.TableName,
	))

	return sb.String()
}

// DropStatisticsSQL generates a DROP STATISTICS statement.
func DropStatisticsSQL(s schema.Statistic) string {
	return fmt.Sprintf("DROP STATISTICS %s.%s;", s.Schema, s.Name)
}

// StatisticsDiffSQL generates SQL statements for a StatisticDiff.
func StatisticsDiffSQL(diff schema.StatisticDiff) []string {
	var stmts []string

	for _, s := range diff.Removed {
		stmts = append(stmts, DropStatisticsSQL(s))
	}

	for _, s := range diff.Changed {
		stmts = append(stmts, DropStatisticsSQL(s))
		stmts = append(stmts, CreateStatisticsSQL(s))
	}

	for _, s := range diff.Added {
		stmts = append(stmts, CreateStatisticsSQL(s))
	}

	return stmts
}
