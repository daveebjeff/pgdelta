package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/internal/schema"
)

// CreateTablespaceSQL generates SQL to create a tablespace.
func CreateTablespaceSQL(ts schema.Tablespace) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE TABLESPACE %s", ts.Name))
	if ts.Owner != "" {
		sb.WriteString(fmt.Sprintf(" OWNER %s", ts.Owner))
	}
	sb.WriteString(fmt.Sprintf(" LOCATION '%s'", ts.Location))
	sb.WriteString(";")
	return sb.String()
}

// DropTablespaceSQL generates SQL to drop a tablespace.
func DropTablespaceSQL(ts schema.Tablespace) string {
	return fmt.Sprintf("DROP TABLESPACE %s;", ts.Name)
}

// AlterTablespaceOwnerSQL generates SQL to alter the owner of a tablespace.
func AlterTablespaceOwnerSQL(ts schema.Tablespace) string {
	return fmt.Sprintf("ALTER TABLESPACE %s OWNER TO %s;", ts.Name, ts.Owner)
}

// TablespaceDiffSQL generates migration SQL statements for tablespace changes.
func TablespaceDiffSQL(diff schema.SchemaDiff) []string {
	var stmts []string

	for _, ts := range diff.AddedTablespaces {
		stmts = append(stmts, CreateTablespaceSQL(ts))
	}

	for _, ts := range diff.RemovedTablespaces {
		stmts = append(stmts, DropTablespaceSQL(ts))
	}

	for _, ts := range diff.ChangedTablespaces {
		stmts = append(stmts, AlterTablespaceOwnerSQL(ts))
	}

	return stmts
}
