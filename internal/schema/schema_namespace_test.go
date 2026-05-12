package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseNamespace = SchemaNamespace{Name: "myschema", Owner: "postgres"}

func TestSchemaNamespaceFullName(t *testing.T) {
	assert.Equal(t, "myschema", baseNamespace.FullName())
}

func TestDiffSchemaNamespaces_NoChanges(t *testing.T) {
	old := []SchemaNamespace{baseNamespace}
	new := []SchemaNamespace{baseNamespace}
	diff := DiffSchemaNamespaces(old, new)
	assert.True(t, diff.IsEmpty())
}

func TestDiffSchemaNamespaces_AddedNamespace(t *testing.T) {
	old := []SchemaNamespace{}
	new := []SchemaNamespace{baseNamespace}
	diff := DiffSchemaNamespaces(old, new)
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, baseNamespace, diff.Added[0])
	assert.Empty(t, diff.Removed)
}

func TestDiffSchemaNamespaces_RemovedNamespace(t *testing.T) {
	old := []SchemaNamespace{baseNamespace}
	new := []SchemaNamespace{}
	diff := DiffSchemaNamespaces(old, new)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, baseNamespace, diff.Removed[0])
	assert.Empty(t, diff.Added)
}

func TestDiffSchemaNamespaces_ChangedOwner(t *testing.T) {
	old := []SchemaNamespace{baseNamespace}
	new := []SchemaNamespace{{Name: "myschema", Owner: "newowner"}}
	diff := DiffSchemaNamespaces(old, new)
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, "newowner", diff.Changed[0].Owner)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}
