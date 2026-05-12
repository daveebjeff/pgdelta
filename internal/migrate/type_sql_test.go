package migrate

import (
	"testing"

	"github.com/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

var baseTyp = schema.Type{
	Schema:     "public",
	Name:       "address",
	Kind:       "composite",
	Definition: "(street text, city text)",
}

func TestCreateTypeSQL_Composite(t *testing.T) {
	sql := CreateTypeSQL(baseTyp)
	assert.Equal(t, "CREATE TYPE public.address AS (street text, city text);", sql)
}

func TestCreateTypeSQL_Domain(t *testing.T) {
	domainType := schema.Type{
		Schema:     "public",
		Name:       "positive_int",
		Kind:       "domain",
		Definition: "integer CHECK (VALUE > 0)",
	}
	sql := CreateTypeSQL(domainType)
	assert.Equal(t, "CREATE DOMAIN public.positive_int AS integer CHECK (VALUE > 0);", sql)
}

func TestDropTypeSQL_Composite(t *testing.T) {
	sql := DropTypeSQL(baseTyp)
	assert.Equal(t, "DROP TYPE IF EXISTS public.address;", sql)
}

func TestDropTypeSQL_Domain(t *testing.T) {
	domainType := schema.Type{Schema: "public", Name: "positive_int", Kind: "domain", Definition: "integer"}
	sql := DropTypeSQL(domainType)
	assert.Equal(t, "DROP DOMAIN IF EXISTS public.positive_int;", sql)
}

func TestTypeDiffSQL_AddedAndRemoved(t *testing.T) {
	newType := schema.Type{Schema: "public", Name: "color", Kind: "composite", Definition: "(r int, g int, b int)"}
	diff := schema.TypeDiff{
		Added:   []schema.Type{newType},
		Removed: []schema.Type{baseTyp},
	}
	stmts := TypeDiffSQL(diff)
	assert.Len(t, stmts, 2)
	assert.Contains(t, stmts[0], "DROP TYPE")
	assert.Contains(t, stmts[1], "CREATE TYPE")
}

func TestTypeDiffSQL_Changed(t *testing.T) {
	modified := baseTyp
	modified.Definition = "(street text, city text, country text)"
	diff := schema.TypeDiff{Changed: []schema.Type{modified}}
	stmts := TypeDiffSQL(diff)
	assert.Len(t, stmts, 2)
	assert.Contains(t, stmts[0], "DROP TYPE IF EXISTS public.address")
	assert.Contains(t, stmts[1], "CREATE TYPE public.address AS (street text, city text, country text)")
}
