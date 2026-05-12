package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/internal/schema"
)

// CreateTypeSQL generates a CREATE TYPE statement.
func CreateTypeSQL(t schema.Type) string {
	switch strings.ToLower(t.Kind) {
	case "composite":
		return fmt.Sprintf("CREATE TYPE %s AS %s;", t.FullName(), t.Definition)
	case "domain":
		return fmt.Sprintf("CREATE DOMAIN %s AS %s;", t.FullName(), t.Definition)
	default:
		return fmt.Sprintf("-- unsupported type kind %q for %s", t.Kind, t.FullName())
	}
}

// DropTypeSQL generates a DROP TYPE (or DROP DOMAIN) statement.
func DropTypeSQL(t schema.Type) string {
	switch strings.ToLower(t.Kind) {
	case "domain":
		return fmt.Sprintf("DROP DOMAIN IF EXISTS %s;", t.FullName())
	default:
		return fmt.Sprintf("DROP TYPE IF EXISTS %s;", t.FullName())
	}
}

// TypeDiffSQL generates SQL statements to migrate from old to new type state.
func TypeDiffSQL(diff schema.TypeDiff) []string {
	var stmts []string

	// Drop changed types first (they must be recreated).
	for _, t := range diff.Changed {
		stmts = append(stmts, DropTypeSQL(t))
	}
	for _, t := range diff.Removed {
		stmts = append(stmts, DropTypeSQL(t))
	}
	for _, t := range diff.Added {
		stmts = append(stmts, CreateTypeSQL(t))
	}
	for _, t := range diff.Changed {
		stmts = append(stmts, CreateTypeSQL(t))
	}

	return stmts
}
