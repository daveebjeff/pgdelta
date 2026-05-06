package migrate

import (
	"fmt"
	"strings"

	"github.com/pgdelta/pgdelta/internal/schema"
)

// CreateFunctionSQL generates a CREATE OR REPLACE FUNCTION statement.
func CreateFunctionSQL(f schema.Function) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("CREATE OR REPLACE FUNCTION %s.%s(%s)\n", f.Schema, f.Name, f.Arguments))
	sb.WriteString(fmt.Sprintf("RETURNS %s\n", f.ReturnType))
	sb.WriteString(fmt.Sprintf("LANGUAGE %s\n", strings.ToLower(f.Language)))
	if f.Volatility != "" {
		sb.WriteString(fmt.Sprintf("%s\n", f.Volatility))
	}
	sb.WriteString("AS $$\n")
	sb.WriteString(f.Body)
	sb.WriteString("\n$$;")
	return sb.String()
}

// DropFunctionSQL generates a DROP FUNCTION statement.
func DropFunctionSQL(f schema.Function) string {
	return fmt.Sprintf("DROP FUNCTION IF EXISTS %s.%s(%s);", f.Schema, f.Name, f.Arguments)
}

// FunctionDiffSQL generates SQL statements for a FunctionDiff.
func FunctionDiffSQL(diff schema.FunctionDiff) []string {
	var stmts []string

	for _, f := range diff.Removed {
		stmts = append(stmts, DropFunctionSQL(f))
	}

	for _, f := range diff.Added {
		stmts = append(stmts, CreateFunctionSQL(f))
	}

	// Changed functions are re-created using CREATE OR REPLACE.
	for _, f := range diff.Changed {
		stmts = append(stmts, CreateFunctionSQL(f))
	}

	return stmts
}
