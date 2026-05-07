package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseExtension = Extension{
	Name:    "uuid-ossp",
	Schema:  "public",
	Version: "1.1",
}

func TestExtensionFullName(t *testing.T) {
	e := baseExtension
	assert.Equal(t, "uuid-ossp", e.FullName())
}

func TestDiffExtensions_NoChanges(t *testing.T) {
	exts := []Extension{baseExtension}
	diff := DiffExtensions(exts, exts)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffExtensions_AddedExtension(t *testing.T) {
	old := []Extension{}
	new := []Extension{baseExtension}
	diff := DiffExtensions(old, new)
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, "uuid-ossp", diff.Added[0].Name)
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffExtensions_RemovedExtension(t *testing.T) {
	old := []Extension{baseExtension}
	new := []Extension{}
	diff := DiffExtensions(old, new)
	assert.Empty(t, diff.Added)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, "uuid-ossp", diff.Removed[0].Name)
	assert.Empty(t, diff.Changed)
}

func TestDiffExtensions_ChangedExtension(t *testing.T) {
	old := []Extension{baseExtension}
	new := []Extension{{Name: "uuid-ossp", Schema: "public", Version: "1.2"}}
	diff := DiffExtensions(old, new)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
	assert.Len(t, diff.Changed, 1)
	assert.Equal(t, "1.2", diff.Changed[0].Version)
}
