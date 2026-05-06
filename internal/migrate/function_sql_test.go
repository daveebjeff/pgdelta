package migrate

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

func baseFunc() schema.Function {
	return schema.Function{
		Schema:     "public",
		Name:       "add",
		Arguments:  "a integer, b integer",
		ReturnType: "integer",
		Language:   "plpgsql",
		Body:       "BEGIN RETURN a + b; END;",
		Volatility: "IMMUTABLE",
	}
}

func TestCreateFunctionSQL(t *testing.T) {
	f := baseFunc()
	sql := CreateFunctionSQL(f)
	assert.Contains(t, sql, "CREATE OR REPLACE FUNCTION public.add(a integer, b integer)")
	assert.Contains(t, sql, "RETURNS integer")
	assert.Contains(t, sql, "LANGUAGE plpgsql")
	assert.Contains(t, sql, "IMMUTABLE")
	assert.Contains(t, sql, "BEGIN RETURN a + b; END;")
	assert.Contains(t, sql, "$$;")
}

func TestDropFunctionSQL(t *testing.T) {
	f := baseFunc()
	sql := DropFunctionSQL(f)
	assert.Equal(t, "DROP FUNCTION IF EXISTS public.add(a integer, b integer);", sql)
}

func TestFunctionDiffSQL_AddedAndRemoved(t *testing.T) {
	added := baseFunc()
	added.Name = "multiply"
	removed := baseFunc()

	diff := schema.FunctionDiff{
		Added:   []schema.Function{added},
		Removed: []schema.Function{removed},
	}

	stmts := FunctionDiffSQL(diff)
	assert.Len(t, stmts, 2)
	assert.Contains(t, stmts[0], "DROP FUNCTION")
	assert.Contains(t, stmts[1], "CREATE OR REPLACE FUNCTION")
}

func TestFunctionDiffSQL_Changed(t *testing.T) {
	f := baseFunc()
	f.Body = "BEGIN RETURN a * b; END;"

	diff := schema.FunctionDiff{
		Changed: []schema.Function{f},
	}

	stmts := FunctionDiffSQL(diff)
	assert.Len(t, stmts, 1)
	assert.Contains(t, stmts[0], "CREATE OR REPLACE FUNCTION")
	assert.Contains(t, stmts[0], "BEGIN RETURN a * b; END;")
}
