package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreatePublicationSQL generates a CREATE PUBLICATION statement.
func CreatePublicationSQL(p schema.Publication) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE PUBLICATION %s", p.Name))
	if p.AllTables {
		sb.WriteString(" FOR ALL TABLES")
	} else if len(p.Tables) > 0 {
		sb.WriteString(fmt.Sprintf(" FOR TABLE %s", strings.Join(p.Tables, ", ")))
	}
	ops := publicationOps(p)
	if len(ops) > 0 {
		sb.WriteString(fmt.Sprintf(" WITH (publish = '%s')", strings.Join(ops, ", ")))
	}
	sb.WriteString(";")
	return sb.String()
}

// DropPublicationSQL generates a DROP PUBLICATION statement.
func DropPublicationSQL(p schema.Publication) string {
	return fmt.Sprintf("DROP PUBLICATION %s;", p.Name)
}

// AlterPublicationSQL generates an ALTER PUBLICATION statement for a changed publication.
func AlterPublicationSQL(p schema.Publication) string {
	var sb strings.Builder
	if len(p.Tables) > 0 {
		sb.WriteString(fmt.Sprintf("ALTER PUBLICATION %s SET TABLE %s;", p.Name, strings.Join(p.Tables, ", ")))
	} else if p.AllTables {
		sb.WriteString(fmt.Sprintf("ALTER PUBLICATION %s FOR ALL TABLES;", p.Name))
	}
	ops := publicationOps(p)
	if len(ops) > 0 {
		if sb.Len() > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(fmt.Sprintf("ALTER PUBLICATION %s SET (publish = '%s');", p.Name, strings.Join(ops, ", ")))
	}
	return sb.String()
}

// PublicationDiffSQL generates SQL statements for a PublicationDiff.
func PublicationDiffSQL(diff schema.PublicationDiff) []string {
	var stmts []string
	for _, p := range diff.Removed {
		stmts = append(stmts, DropPublicationSQL(p))
	}
	for _, p := range diff.Added {
		stmts = append(stmts, CreatePublicationSQL(p))
	}
	for _, p := range diff.Changed {
		stmts = append(stmts, AlterPublicationSQL(p))
	}
	return stmts
}

func publicationOps(p schema.Publication) []string {
	var ops []string
	if p.Insert {
		ops = append(ops, "insert")
	}
	if p.Update {
		ops = append(ops, "update")
	}
	if p.Delete {
		ops = append(ops, "delete")
	}
	if p.Truncate {
		ops = append(ops, "truncate")
	}
	return ops
}
