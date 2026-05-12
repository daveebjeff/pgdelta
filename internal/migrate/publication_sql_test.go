package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

var basePub = schema.Publication{
	Schema: "public",
	Name:   "my_pub",
	Tables: []string{"orders", "users"},
	AllTables: false,
	Insert:    true,
	Update:    true,
	Delete:    true,
	Truncate:  false,
}

func TestCreatePublicationSQL_WithTables(t *testing.T) {
	sql := migrate.CreatePublicationSQL(basePub)
	expected := `CREATE PUBLICATION my_pub FOR TABLE orders, users WITH (publish = 'insert,update,delete');`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestCreatePublicationSQL_AllTables(t *testing.T) {
	pub := basePub
	pub.AllTables = true
	pub.Tables = nil
	sql := migrate.CreatePublicationSQL(pub)
	expected := `CREATE PUBLICATION my_pub FOR ALL TABLES WITH (publish = 'insert,update,delete');`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestDropPublicationSQL(t *testing.T) {
	sql := migrate.DropPublicationSQL(basePub)
	expected := `DROP PUBLICATION my_pub;`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestAlterPublicationSQL_ChangedTables(t *testing.T) {
	newPub := basePub
	newPub.Tables = []string{"orders", "users", "products"}
	sql := migrate.AlterPublicationSQL(basePub, newPub)
	expected := `ALTER PUBLICATION my_pub SET TABLE orders, users, products;`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestPublicationDiffSQL_AddedAndRemoved(t *testing.T) {
	added := []schema.Publication{basePub}
	removed := []schema.Publication{
		{Schema: "public", Name: "old_pub", AllTables: true, Insert: true},
	}
	changed := []schema.PublicationDiff{}

	sqls := migrate.PublicationDiffSQL(added, removed, changed)
	if len(sqls) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(sqls))
	}
}

func TestPublicationDiffSQL_Changed(t *testing.T) {
	newPub := basePub
	newPub.Tables = []string{"orders"}
	diffs := []schema.PublicationDiff{{Old: basePub, New: newPub}}

	sqls := migrate.PublicationDiffSQL(nil, nil, diffs)
	if len(sqls) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(sqls))
	}
}
