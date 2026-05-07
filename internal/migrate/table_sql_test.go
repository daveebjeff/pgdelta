package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

var baseTableForSQL = schema.Table{
	Schema: "public",
	Name:   "users",
	Columns: []schema.Column{
		{Name: "id", DataType: "integer", Nullable: false},
		{Name: "email", DataType: "text", Nullable: false},
	},
}

func TestCreateTableSQL(t *testing.T) {
	sql := migrate.CreateTableSQL(baseTableForSQL)
	expected := "CREATE TABLE public.users (\n\tid integer NOT NULL,\n\temail text NOT NULL\n);"
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestCreateTableSQL_WithDefault(t *testing.T) {
	def := "0"
	tbl := schema.Table{
		Schema: "public",
		Name:   "orders",
		Columns: []schema.Column{
			{Name: "amount", DataType: "integer", Nullable: true, Default: &def},
		},
	}
	sql := migrate.CreateTableSQL(tbl)
	expected := "CREATE TABLE public.orders (\n\tamount integer DEFAULT 0\n);"
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestDropTableSQL(t *testing.T) {
	sql := migrate.DropTableSQL(baseTableForSQL)
	expected := "DROP TABLE public.users;"
	if sql != expected {
		t.Errorf("expected %q, got %q", expected, sql)
	}
}

func TestTableDiffSQL_AddedAndRemoved(t *testing.T) {
	added := schema.Table{
		Schema: "public", Name: "products",
		Columns: []schema.Column{{Name: "id", DataType: "serial", Nullable: false}},
	}
	removed := schema.Table{
		Schema: "public", Name: "legacy",
		Columns: []schema.Column{{Name: "id", DataType: "integer", Nullable: false}},
	}
	diff := schema.TableDiff{
		Added:   []schema.Table{added},
		Removed: []schema.Table{removed},
	}
	stmts := migrate.TableDiffSQL(diff)
	if len(stmts) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(stmts))
	}
	if stmts[0] != "CREATE TABLE public.products (\n\tid serial NOT NULL\n);" {
		t.Errorf("unexpected create statement: %s", stmts[0])
	}
	if stmts[1] != "DROP TABLE public.legacy;" {
		t.Errorf("unexpected drop statement: %s", stmts[1])
	}
}
