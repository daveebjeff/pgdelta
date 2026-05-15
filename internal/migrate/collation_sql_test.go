package migrate_test

import (
	"testing"

	"github.com/pgdelta/internal/migrate"
	"github.com/pgdelta/internal/schema"
)

var baseCollation = schema.Collation{
	Schema:   "public",
	Name:     "my_collation",
	Provider: "icu",
	Locale:   "en-US",
	Deterministic: true,
}

func TestCreateCollationSQL_Basic(t *testing.T) {
	c := baseCollation
	sql := migrate.CreateCollationSQL(c)
	expected := `CREATE COLLATION public.my_collation (PROVIDER = icu, LOCALE = 'en-US', DETERMINISTIC = true);`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestCreateCollationSQL_NonDeterministic(t *testing.T) {
	c := baseCollation
	c.Deterministic = false
	sql := migrate.CreateCollationSQL(c)
	expected := `CREATE COLLATION public.my_collation (PROVIDER = icu, LOCALE = 'en-US', DETERMINISTIC = false);`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestDropCollationSQL(t *testing.T) {
	c := baseCollation
	sql := migrate.DropCollationSQL(c)
	expected := `DROP COLLATION public.my_collation;`
	if sql != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, sql)
	}
}

func TestCollationDiffSQL_AddedAndRemoved(t *testing.T) {
	added := []schema.Collation{baseCollation}
	removed := []schema.Collation{
		{Schema: "public", Name: "old_collation", Provider: "libc", Locale: "en_US.utf8", Deterministic: true},
	}

	diff := schema.SchemaDiff{
		AddedCollations:   added,
		RemovedCollations: removed,
	}

	sqls := migrate.CollationDiffSQL(diff)
	if len(sqls) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(sqls))
	}
	if sqls[0] != `DROP COLLATION public.old_collation;` {
		t.Errorf("unexpected drop statement: %s", sqls[0])
	}
	if sqls[1] != `CREATE COLLATION public.my_collation (PROVIDER = icu, LOCALE = 'en-US', DETERMINISTIC = true);` {
		t.Errorf("unexpected create statement: %s", sqls[1])
	}
}

func TestCollationDiffSQL_Changed(t *testing.T) {
	old := baseCollation
	new := baseCollation
	new.Locale = "fr-FR"

	diff := schema.SchemaDiff{
		ChangedCollations: []schema.CollationDiff{{Old: old, New: new}},
	}

	sqls := migrate.CollationDiffSQL(diff)
	if len(sqls) != 2 {
		t.Fatalf("expected 2 statements for changed collation, got %d", len(sqls))
	}
}
