package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseType = Type{
	Schema:     "public",
	Name:       "address",
	Kind:       "composite",
	Definition: "(street text, city text, zip text)",
}

func TestTypeFullName(t *testing.T) {
	assert.Equal(t, "public.address", baseType.FullName())
}

func TestDiffTypes_NoChanges(t *testing.T) {
	types := []Type{baseType}
	diff := DiffTypes(types, types)
	assert.True(t, diff.IsEmpty())
}

func TestDiffTypes_AddedType(t *testing.T) {
	newType := Type{Schema: "public", Name: "color", Kind: "composite", Definition: "(r int, g int, b int)"}
	diff := DiffTypes([]Type{baseType}, []Type{baseType, newType})
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, "public.color", diff.Added[0].FullName())
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffTypes_RemovedType(t *testing.T) {
	diff := DiffTypes([]Type{baseType}, []Type{})
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, "public.address", diff.Removed[0].FullName())
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffTypes_ChangedType(t *testing.T) {
	modified := baseType
	modified.Definition = "(street text, city text, zip text, country text)"
	diff := DiffTypes([]Type{baseType}, []Type{modified})
	assert.Len(t, diff.Changed, 1)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}
