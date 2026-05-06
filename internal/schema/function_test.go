package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func baseFunction() Function {
	return Function{
		Schema:     "public",
		Name:       "add",
		Arguments:  "a integer, b integer",
		ReturnType: "integer",
		Language:   "plpgsql",
		Body:       "BEGIN RETURN a + b; END;",
		Volatility: "IMMUTABLE",
	}
}

func TestFunctionFullName(t *testing.T) {
	f := baseFunction()
	assert.Equal(t, "public.add(a integer, b integer)", f.FullName())
}

func TestDiffFunctions_NoChanges(t *testing.T) {
	f := baseFunction()
	diff := DiffFunctions([]Function{f}, []Function{f})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffFunctions_AddedFunction(t *testing.T) {
	f := baseFunction()
	diff := DiffFunctions([]Function{}, []Function{f})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, f, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffFunctions_RemovedFunction(t *testing.T) {
	f := baseFunction()
	diff := DiffFunctions([]Function{f}, []Function{})
	assert.Empty(t, diff.Added)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, f, diff.Removed[0])
	assert.Empty(t, diff.Changed)
}

func TestDiffFunctions_ChangedFunction(t *testing.T) {
	old := baseFunction()
	newF := baseFunction()
	newF.Body = "BEGIN RETURN a + b + 1; END;"
	diff := DiffFunctions([]Function{old}, []Function{newF})
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, newF.Body, diff.Changed[0].Body)
}
