package migrate

import (
	"strings"

	"github.com/pgdelta/internal/schema"
)

// SchemaDiffSQL generates a full migration SQL script from a SchemaDiff.
// Statements are ordered to respect dependency constraints:
// extensions → roles → enums → sequences → tables → columns →
// constraints → foreign keys → indexes → views → materialized views →
// functions → triggers → policies
func SchemaDiffSQL(diff schema.SchemaDiff) string {
	var parts []string

	append := func(s string) {
		if s != "" {
			parts = append(parts, s)
		}
	}

	append(ExtensionDiffSQL(diff.ExtensionDiff))
	append(RoleDiffSQL(diff.RoleDiff))
	append(EnumDiffSQL(diff.EnumDiff))
	append(SequenceDiffSQL(diff.SequenceDiff))
	append(TableDiffSQL(diff.TableDiff))
	append(ColumnDiffSQL(diff.ColumnDiff))
	append(ConstraintDiffSQL(diff.ConstraintDiff))
	append(ForeignKeyDiffSQL(diff.ForeignKeyDiff))
	append(IndexDiffSQL(diff.IndexDiff))
	append(ViewDiffSQL(diff.ViewDiff))
	append(MaterializedViewDiffSQL(diff.MaterializedViewDiff))
	append(FunctionDiffSQL(diff.FunctionDiff))
	append(TriggerDiffSQL(diff.TriggerDiff))
	append(PolicyDiffSQL(diff.PolicyDiff))

	return strings.Join(parts, "\n")
}
