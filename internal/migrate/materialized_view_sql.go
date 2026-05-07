package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateMaterializedViewSQL generates a CREATE MATERIALIZED VIEW statement.
func CreateMaterializedViewSQL(mv schema.MaterializedView) string {
	withData := "WITH DATA"
	if !mv.WithData {
		withData = "WITH NO DATA"
	}
	return fmt.Sprintf(
		"CREATE MATERIALIZED VIEW %s AS\n%s\n%s;",
		mv.FullName(),
		mv.Definition,
		withData,
	)
}

// DropMaterializedViewSQL generates a DROP MATERIALIZED VIEW statement.
func DropMaterializedViewSQL(mv schema.MaterializedView) string {
	return fmt.Sprintf("DROP MATERIALIZED VIEW %s;", mv.FullName())
}

// RefreshMaterializedViewSQL generates a REFRESH MATERIALIZED VIEW statement.
func RefreshMaterializedViewSQL(mv schema.MaterializedView) string {
	withData := "WITH DATA"
	if !mv.WithData {
		withData = "WITH NO DATA"
	}
	return fmt.Sprintf("REFRESH MATERIALIZED VIEW %s %s;", mv.FullName(), withData)
}

// MaterializedViewDiffSQL generates SQL statements for a materialized view diff.
// Changed views are recreated (drop + create) since PostgreSQL does not support
// ALTER MATERIALIZED VIEW for definition changes.
func MaterializedViewDiffSQL(diff schema.MaterializedViewDiff) []string {
	var stmts []string

	for _, mv := range diff.Removed {
		stmts = append(stmts, DropMaterializedViewSQL(mv))
	}

	for _, mv := range diff.Changed {
		stmts = append(stmts, DropMaterializedViewSQL(mv))
		stmts = append(stmts, CreateMaterializedViewSQL(mv))
	}

	for _, mv := range diff.Added {
		stmts = append(stmts, CreateMaterializedViewSQL(mv))
	}

	return stmts
}

// joinMaterializedViewSQL is a helper to join statements for testing.
func joinMaterializedViewSQL(stmts []string) string {
	return strings.Join(stmts, "\n")
}

// suppress unused warning
var _ = joinMaterializedViewSQL
