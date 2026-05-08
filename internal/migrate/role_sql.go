package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateRoleSQL generates a CREATE ROLE statement.
func CreateRoleSQL(r schema.Role) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE ROLE %q", r.Name))
	sb.WriteString(boolAttr(r.Superuser, "SUPERUSER", "NOSUPERUSER"))
	sb.WriteString(boolAttr(r.Inherit, "INHERIT", "NOINHERIT"))
	sb.WriteString(boolAttr(r.CreateRole, "CREATEROLE", "NOCREATEROLE"))
	sb.WriteString(boolAttr(r.CreateDB, "CREATEDB", "NOCREATEDB"))
	sb.WriteString(boolAttr(r.Login, "LOGIN", "NOLOGIN"))
	sb.WriteString(boolAttr(r.Replication, "REPLICATION", "NOREPLICATION"))
	sb.WriteString(boolAttr(r.BypassRLS, "BYPASSRLS", "NOBYPASSRLS"))
	if r.ConnectionLimit >= 0 {
		sb.WriteString(fmt.Sprintf(" CONNECTION LIMIT %d", r.ConnectionLimit))
	}
	if r.ValidUntil != nil {
		sb.WriteString(fmt.Sprintf(" VALID UNTIL '%s'", *r.ValidUntil))
	}
	sb.WriteString(";")
	return sb.String()
}

// DropRoleSQL generates a DROP ROLE statement.
func DropRoleSQL(r schema.Role) string {
	return fmt.Sprintf("DROP ROLE %q;", r.Name)
}

// AlterRoleSQL generates an ALTER ROLE statement.
func AlterRoleSQL(r schema.Role) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ALTER ROLE %q", r.Name))
	sb.WriteString(boolAttr(r.Superuser, "SUPERUSER", "NOSUPERUSER"))
	sb.WriteString(boolAttr(r.Inherit, "INHERIT", "NOINHERIT"))
	sb.WriteString(boolAttr(r.CreateRole, "CREATEROLE", "NOCREATEROLE"))
	sb.WriteString(boolAttr(r.CreateDB, "CREATEDB", "NOCREATEDB"))
	sb.WriteString(boolAttr(r.Login, "LOGIN", "NOLOGIN"))
	sb.WriteString(boolAttr(r.Replication, "REPLICATION", "NOREPLICATION"))
	sb.WriteString(boolAttr(r.BypassRLS, "BYPASSRLS", "NOBYPASSRLS"))
	if r.ConnectionLimit >= 0 {
		sb.WriteString(fmt.Sprintf(" CONNECTION LIMIT %d", r.ConnectionLimit))
	}
	if r.ValidUntil != nil {
		sb.WriteString(fmt.Sprintf(" VALID UNTIL '%s'", *r.ValidUntil))
	}
	sb.WriteString(";")
	return sb.String()
}

// RoleDiffSQL generates migration SQL for a RoleDiff.
func RoleDiffSQL(diff schema.RoleDiff) []string {
	var stmts []string
	for _, r := range diff.Removed {
		stmts = append(stmts, DropRoleSQL(r))
	}
	for _, r := range diff.Added {
		stmts = append(stmts, CreateRoleSQL(r))
	}
	for _, r := range diff.Changed {
		stmts = append(stmts, AlterRoleSQL(r))
	}
	return stmts
}

func boolAttr(val bool, trueStr, falseStr string) string {
	if val {
		return " " + trueStr
	}
	return " " + falseStr
}
