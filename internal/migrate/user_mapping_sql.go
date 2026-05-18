package migrate

import (
	"fmt"
	"sort"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

func CreateUserMappingSQL(m schema.UserMapping) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "CREATE USER MAPPING FOR %s SERVER %s", m.User, m.Server)
	if len(m.Options) > 0 {
		fmt.Fprintf(&sb, " OPTIONS (%s)", userMappingOptionsString(m.Options))
	}
	sb.WriteString(";")
	return sb.String()
}

func DropUserMappingSQL(m schema.UserMapping) string {
	return fmt.Sprintf("DROP USER MAPPING IF EXISTS FOR %s SERVER %s;", m.User, m.Server)
}

func AlterUserMappingSQL(m schema.UserMapping) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "ALTER USER MAPPING FOR %s SERVER %s", m.User, m.Server)
	if len(m.Options) > 0 {
		fmt.Fprintf(&sb, " OPTIONS (SET %s)", userMappingOptionsString(m.Options))
	}
	sb.WriteString(";")
	return sb.String()
}

func UserMappingDiffSQL(diff schema.SchemaDiff) []string {
	var stmts []string
	for _, item := range diff.Removed {
		m := item.(schema.UserMapping)
		stmts = append(stmts, DropUserMappingSQL(m))
	}
	for _, item := range diff.Added {
		m := item.(schema.UserMapping)
		stmts = append(stmts, CreateUserMappingSQL(m))
	}
	for _, ch := range diff.Changed {
		m := ch.New.(schema.UserMapping)
		stmts = append(stmts, AlterUserMappingSQL(m))
	}
	return stmts
}

func userMappingOptionsString(opts map[string]string) string {
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
