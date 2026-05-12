package migrate

import (
	"testing"

	"github.com/pgdelta/pgdelta/internal/schema"
	"github.com/stretchr/testify/assert"
)

func TestCreateSchemaSQL_WithOwner(t *testing.T) {
	s := schema.SchemaNamespace{Name: "analytics", Owner: "postgres"}
	result := CreateSchemaSQL(s)
	assert.Equal(t, "CREATE SCHEMA analytics AUTHORIZATION postgres;", result)
}

func TestCreateSchemaSQL_NoOwner(t *testing.T) {
	s := schema.SchemaNamespace{Name: "analytics", Owner: ""}
	result := CreateSchemaSQL(s)
	assert.Equal(t, "CREATE SCHEMA analytics;", result)
}

func TestDropSchemaSQL(t *testing.T) {
	s := schema.SchemaNamespace{Name: "analytics", Owner: "postgres"}
	result := DropSchemaSQL(s)
	assert.Equal(t, "DROP SCHEMA IF EXISTS analytics;", result)
}

func TestAlterSchemaOwnerSQL(t *testing.T) {
	s := schema.SchemaNamespace{Name: "analytics", Owner: "newowner"}
	result := AlterSchemaOwnerSQL(s)
	assert.Equal(t, "ALTER SCHEMA analytics OWNER TO newowner;", result)
}

func TestSchemaNamespaceDiffSQL_AddedAndRemoved(t *testing.T) {
	diff := schema.SchemaNSDiff{
		Added:   []schema.SchemaNamespace{{Name: "newschema", Owner: "postgres"}},
		Removed: []schema.SchemaNamespace{{Name: "oldschema", Owner: "postgres"}},
	}
	result := SchemaNamespaceDiffSQL(diff)
	assert.Contains(t, result, "CREATE SCHEMA newschema")
	assert.Contains(t, result, "DROP SCHEMA IF EXISTS oldschema")
}

func TestSchemaNamespaceDiffSQL_Changed(t *testing.T) {
	diff := schema.SchemaNSDiff{
		Changed: []schema.SchemaNamespace{{Name: "myschema", Owner: "newowner"}},
	}
	result := SchemaNamespaceDiffSQL(diff)
	assert.Equal(t, "ALTER SCHEMA myschema OWNER TO newowner;", result)
}
