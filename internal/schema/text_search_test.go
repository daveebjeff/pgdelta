package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var baseTextSearch = TextSearchConfig{
	Schema: "public",
	Name:   "english_config",
	Parser: "pg_catalog.default",
}

func TestTextSearchConfigFullName(t *testing.T) {
	assert.Equal(t, "public.english_config", baseTextSearch.FullName())
}

func TestDiffTextSearchConfigs_NoChanges(t *testing.T) {
	old := []TextSearchConfig{baseTextSearch}
	new := []TextSearchConfig{baseTextSearch}
	diff := DiffTextSearchConfigs(old, new)
	assert.True(t, diff.IsEmpty())
}

func TestDiffTextSearchConfigs_AddedConfig(t *testing.T) {
	old := []TextSearchConfig{}
	new := []TextSearchConfig{baseTextSearch}
	diff := DiffTextSearchConfigs(old, new)
	assert.Len(t, diff.Added, 1)
	assert.Equal(t, baseTextSearch, diff.Added[0])
	assert.Empty(t, diff.Removed)
	assert.Empty(t, diff.Changed)
}

func TestDiffTextSearchConfigs_RemovedConfig(t *testing.T) {
	old := []TextSearchConfig{baseTextSearch}
	new := []TextSearchConfig{}
	diff := DiffTextSearchConfigs(old, new)
	assert.Len(t, diff.Removed, 1)
	assert.Equal(t, baseTextSearch, diff.Removed[0])
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Changed)
}

func TestDiffTextSearchConfigs_ChangedConfig(t *testing.T) {
	old := []TextSearchConfig{baseTextSearch}
	new := []TextSearchConfig{
		{Schema: "public", Name: "english_config", Parser: "pg_catalog.simple"},
	}
	diff := DiffTextSearchConfigs(old, new)
	assert.Len(t, diff.Changed, 1)
	assert.Empty(t, diff.Added)
	assert.Empty(t, diff.Removed)
}
