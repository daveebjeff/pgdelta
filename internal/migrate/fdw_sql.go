package migrate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateForeignDataWrapperSQL generates a CREATE FOREIGN DATA WRAPPER statement.
func CreateForeignDataWrapperSQL(f schema.ForeignDataWrapper) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "CREATE FOREIGN DATA WRAPPER %s", f.Name)
	if f.Handler != "" {
		fmt.Fprintf(&sb, "\n  HANDLER %s", f.Handler)
	} else {
		sb.WriteString("\n  NO HANDLER")
	}
	if f.Validator != "" {
		fmt.Fprintf(&sb, "\n  VALIDATOR %s", f.Validator)
	} else {
		sb.WriteString("\n  NO VALIDATOR")
	}
	if len(f.Options) > 0 {
		sb.WriteString("\n  OPTIONS (")
		sb.WriteString(fdwOptionsString(f.Options))
		sb.WriteString(")")
	}
	sb.WriteString(";")
	return sb.String()
}

// DropForeignDataWrapperSQL generates a DROP FOREIGN DATA WRAPPER statement.
func DropForeignDataWrapperSQL(f schema.ForeignDataWrapper) string {
	return fmt.Sprintf("DROP FOREIGN DATA WRAPPER %s;", f.Name)
}

// AlterForeignDataWrapperSQL generates an ALTER FOREIGN DATA WRAPPER statement for changed FDWs.
func AlterForeignDataWrapperSQL(f schema.ForeignDataWrapper) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "ALTER FOREIGN DATA WRAPPER %s", f.Name)
	if f.Handler != "" {
		fmt.Fprintf(&sb, "\n  HANDLER %s", f.Handler)
	}
	if f.Validator != "" {
		fmt.Fprintf(&sb, "\n  VALIDATOR %s", f.Validator)
	}
	if len(f.Options) > 0 {
		sb.WriteString("\n  OPTIONS (SET ")
		sb.WriteString(fdwOptionsString(f.Options))
		sb.WriteString(")")
	}
	sb.WriteString(";")
	return sb.String()
}

// FDWDiffSQL returns SQL statements representing the diff between old and new FDWs.
func FDWDiffSQL(old, new []schema.ForeignDataWrapper) []string {
	added, removed, changed := schema.DiffForeignDataWrappers(old, new)
	var stmts []string
	for _, f := range removed {
		stmts = append(stmts, DropForeignDataWrapperSQL(f))
	}
	for _, f := range added {
		stmts = append(stmts, CreateForeignDataWrapperSQL(f))
	}
	for _, f := range changed {
		stmts = append(stmts, AlterForeignDataWrapperSQL(f))
	}
	return stmts
}

func fdwOptionsString(opts map[string]string) string {
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
