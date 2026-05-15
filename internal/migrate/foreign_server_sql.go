package migrate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateForeignServerSQL generates a CREATE SERVER statement.
func CreateForeignServerSQL(fs schema.ForeignServer) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE SERVER %s", fs.Name))
	if fs.Type != "" {
		sb.WriteString(fmt.Sprintf(" TYPE '%s'", fs.Type))
	}
	if fs.Version != "" {
		sb.WriteString(fmt.Sprintf(" VERSION '%s'", fs.Version))
	}
	sb.WriteString(fmt.Sprintf(" FOREIGN DATA WRAPPER %s", fs.FDWName))
	if opts := foreignServerOptionsClause(fs.Options); opts != "" {
		sb.WriteString(fmt.Sprintf(" OPTIONS (%s)", opts))
	}
	sb.WriteString(";")
	if fs.Owner != "" {
		sb.WriteString(fmt.Sprintf("\nALTER SERVER %s OWNER TO %s;", fs.Name, fs.Owner))
	}
	return sb.String()
}

// DropForeignServerSQL generates a DROP SERVER statement.
func DropForeignServerSQL(fs schema.ForeignServer) string {
	return fmt.Sprintf("DROP SERVER IF EXISTS %s CASCADE;", fs.Name)
}

// AlterForeignServerSQL generates ALTER SERVER statements for a changed server.
func AlterForeignServerSQL(fs schema.ForeignServer) string {
	var sb strings.Builder
	if fs.Version != "" {
		sb.WriteString(fmt.Sprintf("ALTER SERVER %s VERSION '%s';", fs.Name, fs.Version))
	}
	if opts := foreignServerOptionsClause(fs.Options); opts != "" {
		sb.WriteString(fmt.Sprintf("\nALTER SERVER %s OPTIONS (%s);", fs.Name, opts))
	}
	if fs.Owner != "" {
		sb.WriteString(fmt.Sprintf("\nALTER SERVER %s OWNER TO %s;", fs.Name, fs.Owner))
	}
	return sb.String()
}

// ForeignServerDiffSQL returns SQL statements for added, removed, and changed foreign servers.
func ForeignServerDiffSQL(added, removed, changed []schema.ForeignServer) []string {
	var stmts []string
	for _, fs := range removed {
		stmts = append(stmts, DropForeignServerSQL(fs))
	}
	for _, fs := range added {
		stmts = append(stmts, CreateForeignServerSQL(fs))
	}
	for _, fs := range changed {
		stmts = append(stmts, AlterForeignServerSQL(fs))
	}
	return stmts
}

func foreignServerOptionsClause(opts map[string]string) string {
	if len(opts) == 0 {
		return ""
	}
	keys := make([]string, 0, len(opts))
	for k := range opts {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	parts := make([]string, 0, len(keys))
	for _, k := range keys {
		parts = append(parts, fmt.Sprintf("%s '%s'", k, opts[k]))
	}
	return strings.Join(parts, ", ")
}
