package migrate

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

var baseTSConfig = schema.TextSearchConfig{
	Schema: "public",
	Name:   "english_config",
	Parser: "pg_catalog.default",
}

func TestCreateTextSearchConfigSQL(t *testing.T) {
	sql := CreateTextSearchConfigSQL(baseTSConfig)
	assert.Equal(t,
		"CREATE TEXT SEARCH CONFIGURATION public.english_config (PARSER = pg_catalog.default);",
		sql,
	)
}

func TestDropTextSearchConfigSQL(t *testing.T) {
	sql := DropTextSearchConfigSQL(baseTSConfig)
	assert.Equal(t,
		"DROP TEXT SEARCH CONFIGURATION public.english_config;",
		sql,
	)
}

func TestTextSearchDiffSQL_AddedAndRemoved(t *testing.T) {
	added := schema.TextSearchConfig{Schema: "public", Name: "simple_config", Parser: "pg_catalog.simple"}
	diff := schema.TextSearchDiff{
		Added:   []schema.TextSearchConfig{added},
		Removed: []schema.TextSearchConfig{baseTSConfig},
	}
	sql := TextSearchDiffSQL(diff)
	assert.Contains(t, sql, "DROP TEXT SEARCH CONFIGURATION public.english_config;")
	assert.Contains(t, sql, "CREATE TEXT SEARCH CONFIGURATION public.simple_config (PARSER = pg_catalog.simple);")
}

func TestTextSearchDiffSQL_Changed(t *testing.T) {
	changed := schema.TextSearchConfig{Schema: "public", Name: "english_config", Parser: "pg_catalog.simple"}
	diff := schema.TextSearchDiff{
		Changed: []schema.TextSearchConfig{changed},
	}
	sql := TextSearchDiffSQL(diff)
	assert.Contains(t, sql, "RENAME TO english_config_old")
	assert.Contains(t, sql, "DROP TEXT SEARCH CONFIGURATION public.english_config_old;")
	assert.Contains(t, sql, "CREATE TEXT SEARCH CONFIGURATION public.english_config (PARSER = pg_catalog.simple);")
}

func TestTextSearchDiffSQL_Empty(t *testing.T) {
	diff := schema.TextSearchDiff{}
	sql := TextSearchDiffSQL(diff)
	assert.Equal(t, "", sql)
}
