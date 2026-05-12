package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// GrantPrivilegeSQL generates a GRANT statement for the given privilege.
func GrantPrivilegeSQL(p schema.Privilege) string {
	privList := strings.Join(p.Privileges, ", ")
	stmt := fmt.Sprintf("GRANT %s ON %s %s TO %s",
		privList,
		p.ObjectType,
		p.FullName(),
		p.Grantee,
	)
	if p.WithGrant {
		stmt += " WITH GRANT OPTION"
	}
	return stmt + ";"
}

// RevokePrivilegeSQL generates a REVOKE statement for the given privilege.
func RevokePrivilegeSQL(p schema.Privilege) string {
	privList := strings.Join(p.Privileges, ", ")
	return fmt.Sprintf("REVOKE %s ON %s %s FROM %s;",
		privList,
		p.ObjectType,
		p.FullName(),
		p.Grantee,
	)
}

// PrivilegeDiffSQL generates SQL statements to migrate from old to new privileges.
func PrivilegeDiffSQL(old, new []schema.Privilege) []string {
	added, removed := schema.DiffPrivileges(old, new)
	var stmts []string
	for _, p := range removed {
		stmts = append(stmts, RevokePrivilegeSQL(p))
	}
	for _, p := range added {
		stmts = append(stmts, GrantPrivilegeSQL(p))
	}
	return stmts
}
