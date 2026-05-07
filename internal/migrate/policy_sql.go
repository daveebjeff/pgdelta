package migrate

import (
	"fmt"
	"strings"

	"github.com/your-org/pgdelta/internal/schema"
)

// CreatePolicySQL generates a CREATE POLICY statement.
func CreatePolicySQL(p schema.Policy) string {
	var sb strings.Builder

	permissive := "PERMISSIVE"
	if !p.Permissive {
		permissive = "RESTRICTIVE"
	}

	sb.WriteString(fmt.Sprintf("CREATE POLICY %s ON %s.%s AS %s FOR %s",
		p.Name, p.Schema, p.Table, permissive, p.Command))

	if len(p.Roles) > 0 {
		sb.WriteString(fmt.Sprintf(" TO %s", strings.Join(p.Roles, ", ")))
	}

	if p.Using != "" {
		sb.WriteString(fmt.Sprintf(" USING %s", p.Using))
	}

	if p.WithCheck != "" {
		sb.WriteString(fmt.Sprintf(" WITH CHECK %s", p.WithCheck))
	}

	sb.WriteString(";")
	return sb.String()
}

// DropPolicySQL generates a DROP POLICY statement.
func DropPolicySQL(p schema.Policy) string {
	return fmt.Sprintf("DROP POLICY %s ON %s.%s;", p.Name, p.Schema, p.Table)
}

// AlterPolicySQL generates an ALTER POLICY statement for USING and WITH CHECK changes.
func AlterPolicySQL(p schema.Policy) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("ALTER POLICY %s ON %s.%s", p.Name, p.Schema, p.Table))

	if len(p.Roles) > 0 {
		sb.WriteString(fmt.Sprintf(" TO %s", strings.Join(p.Roles, ", ")))
	}

	if p.Using != "" {
		sb.WriteString(fmt.Sprintf(" USING %s", p.Using))
	}

	if p.WithCheck != "" {
		sb.WriteString(fmt.Sprintf(" WITH CHECK %s", p.WithCheck))
	}

	sb.WriteString(";")
	return sb.String()
}

// PolicyDiffSQL generates SQL statements for a PolicyDiff.
func PolicyDiffSQL(diff schema.PolicyDiff) []string {
	var stmts []string

	for _, p := range diff.Removed {
		stmts = append(stmts, DropPolicySQL(p))
	}

	for _, p := range diff.Added {
		stmts = append(stmts, CreatePolicySQL(p))
	}

	for _, p := range diff.Changed {
		stmts = append(stmts, AlterPolicySQL(p))
	}

	return stmts
}
