package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateDomainSQL generates a CREATE DOMAIN statement.
func CreateDomainSQL(d schema.Domain) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "CREATE DOMAIN %s AS %s", d.FullName(), d.BaseType)
	if d.Default != nil {
		fmt.Fprintf(&sb, " DEFAULT %s", *d.Default)
	}
	if d.NotNull {
		sb.WriteString(" NOT NULL")
	}
	if d.CheckClause != nil {
		name := ""
		if d.CheckName != nil {
			name = fmt.Sprintf("CONSTRAINT %s ", *d.CheckName)
		}
		fmt.Fprintf(&sb, " %sCHECK (%s)", name, *d.CheckClause)
	}
	sb.WriteString(";")
	return sb.String()
}

// DropDomainSQL generates a DROP DOMAIN statement.
func DropDomainSQL(d schema.Domain) string {
	return fmt.Sprintf("DROP DOMAIN %s;", d.FullName())
}

// DomainDiffSQL generates SQL statements for domain differences.
func DomainDiffSQL(diff schema.DomainDiff) []string {
	var stmts []string
	for _, d := range diff.Removed {
		stmts = append(stmts, DropDomainSQL(d))
	}
	for _, d := range diff.Added {
		stmts = append(stmts, CreateDomainSQL(d))
	}
	for _, d := range diff.Changed {
		stmts = append(stmts, DropDomainSQL(d))
		stmts = append(stmts, CreateDomainSQL(d))
	}
	return stmts
}
