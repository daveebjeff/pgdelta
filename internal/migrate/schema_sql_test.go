package migrate_test

import (
	"strings"
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

func TestSchemaDiffSQL_Empty(t *testing.T) {
	diff := schema.SchemaDiff{}
	sql := migrate.SchemaDiffSQL(diff)
	if sql != "" {
		t.Errorf("expected empty SQL for empty diff, got: %q", sql)
	}
}

func TestSchemaDiffSQL_AddedTable(t *testing.T) {
	diff := schema.SchemaDiff{
		TableDiff: schema.TableDiff{
			Added: []schema.Table{
				{Schema: "public", Name: "orders", Columns: []schema.Column{
					{Schema: "public", TableName: "orders", Name: "id", DataType: "integer", IsNullable: false},
				}},
			},
		},
	}
	sql := migrate.SchemaDiffSQL(diff)
	if !strings.Contains(sql, "CREATE TABLE") {
		t.Errorf("expected CREATE TABLE in output, got: %s", sql)
	}
	if !strings.Contains(sql, "orders") {
		t.Errorf("expected table name 'orders' in output, got: %s", sql)
	}
}

func TestSchemaDiffSQL_OrderingExtensionBeforeTable(t *testing.T) {
	diff := schema.SchemaDiff{
		TableDiff: schema.TableDiff{
			Added: []schema.Table{
				{Schema: "public", Name: "items", Columns: []schema.Column{
					{Schema: "public", TableName: "items", Name: "id", DataType: "uuid", IsNullable: false},
				}},
			},
		},
		ExtensionDiff: schema.ExtensionDiff{
			Added: []schema.Extension{
				{Schema: "public", Name: "uuid-ossp", Version: "1.1"},
			},
		},
	}
	sql := migrate.SchemaDiffSQL(diff)
	extIdx := strings.Index(sql, "CREATE EXTENSION")
	tableIdx := strings.Index(sql, "CREATE TABLE")
	if extIdx == -1 || tableIdx == -1 {
		t.Fatal("expected both CREATE EXTENSION and CREATE TABLE in output")
	}
	if extIdx > tableIdx {
		t.Error("expected CREATE EXTENSION to appear before CREATE TABLE")
	}
}

func TestSchemaDiffSQL_RemovedEnum(t *testing.T) {
	diff := schema.SchemaDiff{
		EnumDiff: schema.EnumDiff{
			Removed: []schema.Enum{
				{Schema: "public", Name: "status", Values: []string{"active", "inactive"}},
			},
		},
	}
	sql := migrate.SchemaDiffSQL(diff)
	if !strings.Contains(sql, "DROP TYPE") {
		t.Errorf("expected DROP TYPE in output, got: %s", sql)
	}
}
